// Copyright 2018 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package main contains the test driver for testing xDS manually.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/nutanix/go-control-plane/pkg/producer"

	"github.com/envoyproxy/go-control-plane/pkg/server"
)

var (
	debug bool

	port         uint
	gatewayPort  uint
	upstreamPort uint
	basePort     uint
	alsPort      uint

	delay    time.Duration
	requests int

	clusters      int
	httpListeners int
	tcpListeners  int

	nodeID string

	xdsCluster string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Use debug logging")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&upstreamPort, "upstream", 18080, "Upstream HTTP/1.1 port")
	flag.UintVar(&basePort, "base", 9999, "Listener port")
	flag.DurationVar(&delay, "delay", 500*time.Millisecond, "Interval between request batch retries")
	flag.IntVar(&requests, "r", 5, "Number of requests for each snapshot")
	flag.IntVar(&clusters, "clusters", 4, "Number of clusters")
	flag.IntVar(&httpListeners, "http", 2, "Number of HTTP listeners (and RDS configs)")
	flag.IntVar(&tcpListeners, "tcp", 2, "Number of TCP pass-through listeners")
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
	flag.StringVar(&xdsCluster, "xdsCluster", "xds_cluster", "Name of the xDS cluster")
}

// GenerateConfig creates a configuration for sending via xDS
func GenerateConfig(version string) cache.Snapshot {
	clusters := make([]cache.Resource, 0)
	endpoints := make([]cache.Resource, 0)
	routes := make([]cache.Resource, 0)
	listeners := make([]cache.Resource, 0)

	// Create a cluster
	clusterName := fmt.Sprintf("cluster-%s", version)
	clusterDef := producer.ClusterConfig{
		ClusterName:        clusterName,
		EdsClusterName:     xdsCluster,
		SubsetSelectorKeys: []string{"node_id"},
	}
	clusters = append(clusters, clusterDef.MakeCluster())

	// Create endpoints for a given cluster
	// Two endpoints are created:
	// ep1: localhost:upstreamPort with node_id: node0
	// ep2: localhost:upstreamPort+1 with node_id: node1
	clusterLoadAssignment := producer.ClusterLoadAssignment{
		ClusterName: clusterName,
		EndpointConfigs: []*producer.EndpointConfig{
			&producer.EndpointConfig{
				SocketAddress: "127.0.0.1",
				SocketPort:    uint32(upstreamPort),
				SubsetLBMap: map[string]string{
					"node_id": "node0",
				},
			},
			&producer.EndpointConfig{
				SocketAddress: "127.0.0.1",
				SocketPort:    uint32(upstreamPort) + 1,
				SubsetLBMap: map[string]string{
					"node_id": "node1",
				},
			},
		},
	}
	endpoints = append(endpoints, clusterLoadAssignment.MakeClusterLoadAssignment())

	// Create a route
	routeName := fmt.Sprintf("route-%s", version)
	routeDef := producer.RouteConfig{
		RouteName: routeName,
		Cluster:   clusterName,
	}
	routes = append(routes, routeDef.MakeRoute())

	// Create a HTTP listener
	httpListenerName := fmt.Sprintf("listener-%d", basePort)
	httpListenerDef := producer.HTTPListenerConfig{
		ListenerName:   httpListenerName,
		ListenerPort:   uint32(basePort),
		RouteName:      routeName,
		LdsClusterName: xdsCluster,
	}
	listeners = append(listeners, httpListenerDef.MakeHTTPListener())

	// Create a snapshot of the above config
	return cache.NewSnapshot(version, endpoints, clusters, routes, listeners)
}

// main returns code 1 if any of the batches failed to pass all requests
func main() {
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	ctx := context.Background()

	// create a cache
	signal := make(chan struct{})
	cb := &callbacks{signal: signal}
	config := cache.NewSnapshotCache(false /*is ADS*/, producer.Hasher{}, logger{})
	srv := server.NewServer(config, cb)

	// start the xDS server
	go producer.RunManagementServer(ctx, srv, port)

	log.Infof("waiting for the first request...")
	<-signal
	log.WithFields(log.Fields{"requests": requests}).Info("executing sequence")

	snapshot := GenerateConfig("v0")
	if err := snapshot.Consistent(); err != nil {
		log.Errorf("snapshot inconsistency: %+v", snapshot)
	}
	err := config.SetSnapshot(nodeID, snapshot)
	if err != nil {
		log.Errorf("snapshot error %q for %+v", err, snapshot)
		os.Exit(1)
	}
	pass := false
	for j := 0; j < requests; j++ {
		ok := callEcho()
		if ok == true {
			pass = true
		}
		log.WithFields(log.Fields{"batch": j, "ok": ok}).Info("request batch")
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return
		}
	}

	if pass == true {
		log.Infof("Test passed!")
	} else {
		log.Infof("Test failed!")
	}
}

// callEcho calls upstream echo service on all listener ports and returns an error
// if any of the listeners returned an error.
func callEcho() bool {
	ch := make(chan error)

	// spawn requests
	go func() {
		client := http.Client{
			Timeout: 100 * time.Millisecond,
		}
		req, err := client.Get(fmt.Sprintf("http://localhost:%d", basePort))
		if err != nil {
			log.Error(err)
			ch <- err
			return
		}
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ch <- err
			return
		}
		if string(body) != producer.Hello {
			ch <- fmt.Errorf("unexpected return %q", string(body))
			return
		}
		ch <- nil
	}()

	out := <-ch
	if out != nil {
		return false
	}
	return true
}

type logger struct{}

func (logger logger) Infof(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
func (logger logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

type callbacks struct {
	signal   chan struct{}
	fetches  int
	requests int
	mu       sync.Mutex
}

func (cb *callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.WithFields(log.Fields{"fetches": cb.fetches, "requests": cb.requests}).Info("server callbacks")
}
func (cb *callbacks) OnStreamOpen(_ context.Context, id int64, typ string) error {
	log.Debugf("stream %d open for %s", id, typ)
	return nil
}
func (cb *callbacks) OnStreamClosed(id int64) {
	log.Debugf("stream %d closed", id)
}
func (cb *callbacks) OnStreamRequest(int64, *v2.DiscoveryRequest) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
}
func (cb *callbacks) OnStreamResponse(int64, *v2.DiscoveryRequest, *v2.DiscoveryResponse) {}
func (cb *callbacks) OnFetchRequest(_ context.Context, req *v2.DiscoveryRequest) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {}
