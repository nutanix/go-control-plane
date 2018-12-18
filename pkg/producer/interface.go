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

// Package producer creates test xDS resources
package producer

import (
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	tcp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/tcp_proxy/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
)

const (
	localhost = "127.0.0.1"

	// XdsCluster is the cluster name for the control server (used by non-ADS set-up)
	XdsCluster = "xds_cluster"
)

var (
	// RefreshDelay for the polling config source
	RefreshDelay = 500 * time.Millisecond
)

// EndpointConfig provides details of the configuration of an endpoint definition
type EndpointConfig struct {
	// Name of the cluster the endpoint belongs to
	ClusterName string

	// IP address of the endpoint socket
	SocketAddress string

	// Port number of the endpoint socket
	SocketPort uint32

	// Subset load balancer key value map
	// Note: ensure keys are present in the subset selector keys provided to
	// corresponding cluster definition
	SubsetLBMap map[string]string
}

// MakeEndpoint creates an endpoint using EndpointConfig
func (endpointCfg *EndpointConfig) MakeEndpoint() *v2.ClusterLoadAssignment {
	if endpointCfg.ClusterName == "" ||
		endpointCfg.SocketAddress == "" ||
		endpointCfg.SocketPort == 0 {
		// Invalid endpoint config struct
		return nil
	}

	subsetLBFields := make(map[string]*types.Value)
	if len(endpointCfg.SubsetLBMap) != 0 {
		for key, val := range endpointCfg.SubsetLBMap {
			subsetLBFields[key] = &types.Value{Kind: &types.Value_StringValue{StringValue: val}}
		}
	}

	return &v2.ClusterLoadAssignment{
		ClusterName: endpointCfg.ClusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: []endpoint.LbEndpoint{{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Protocol: core.TCP,
								Address:  endpointCfg.SocketAddress,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: endpointCfg.SocketPort,
								},
							},
						},
					},
				},
				Metadata: &core.Metadata{
					FilterMetadata: map[string]*types.Struct{
						"envoy.lb": &types.Struct{Fields: subsetLBFields},
					},
				},
			}},
		}},
	}
}

// ClusterConfig provides details of the configuration of a cluster definition
type ClusterConfig struct {
	// Name of the cluster
	ClusterName string

	// Subset load balancing selector keys
	SubsetSelectorKeys []string
}

// MakeCluster creates a cluster using config seeded by ClusterConfig
func (clusterCfg *ClusterConfig) MakeCluster() *v2.Cluster {
	if clusterCfg.ClusterName == "" {
		// Invalid config
		return nil
	}
	// #TODO(utkarsh.simha) Add TLS context
	var edsSource *core.ConfigSource
	edsSource = &core.ConfigSource{
		ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
			ApiConfigSource: &core.ApiConfigSource{
				ApiType: core.ApiConfigSource_GRPC,
				GrpcServices: []*core.GrpcService{{
					TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
					},
				}},
			},
		},
	}

	return &v2.Cluster{
		Name:           clusterCfg.ClusterName,
		ConnectTimeout: 5 * time.Second,
		Type:           v2.Cluster_EDS,
		LbPolicy:       v2.Cluster_ROUND_ROBIN,
		LbSubsetConfig: &v2.Cluster_LbSubsetConfig{
			FallbackPolicy: v2.Cluster_LbSubsetConfig_ANY_ENDPOINT,
			SubsetSelectors: []*v2.Cluster_LbSubsetConfig_LbSubsetSelector{{
				Keys: clusterCfg.SubsetSelectorKeys,
			}},
		},
		EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
	}
}

// RouteConfig is used to configure a particular route to a cluster
type RouteConfig struct {
	// Name of the route - referenced in listener config
	RouteName string

	// Name of the cluster to route to, or name of cluster header to use
	Cluster string

	// If the given cluster is a cluster header
	IsClusterHeader bool

	// URI prefix to match
	RoutePrefix string
}

// MakeRoute creates an HTTP route that routes to a given cluster.
func (routeCfg *RouteConfig) MakeRoute() *v2.RouteConfiguration {
	var routeAction *route.RouteAction
	if routeCfg.IsClusterHeader {
		// cluster is specified dynamically by setting a header with the given key
		routeAction = &route.RouteAction{
			ClusterSpecifier: &route.RouteAction_ClusterHeader{
				ClusterHeader: routeCfg.Cluster,
			},
		}
	} else {
		routeAction = &route.RouteAction{
			ClusterSpecifier: &route.RouteAction_Cluster{
				Cluster: routeCfg.Cluster,
			},
		}
	}
	return &v2.RouteConfiguration{
		Name: routeCfg.RouteName,
		VirtualHosts: []route.VirtualHost{{
			Name:    routeCfg.RouteName,
			Domains: []string{"*"},
			Routes: []route.Route{{
				Match: route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: routeCfg.RoutePrefix,
					},
				},
				Action: &route.Route_Route{
					Route: routeAction,
				},
			}},
		}},
	}
}

// HTTPListenerConfig provides configuration for the HTTP listener
type HTTPListenerConfig struct {
	// Name of the HTTP listener
	ListenerName string

	// Port to bind the HTTP listener to
	ListenerPort uint32

	// Name of the route associated with the listener
	RouteName string
}

// MakeHTTPListener creates a listener using either ADS or RDS for the route.
func (httpListenerCfg *HTTPListenerConfig) MakeHTTPListener() *v2.Listener {
	// #TODO(utkarsh.simha) Add TLS context
	// data source configuration
	rdsSource := core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			ApiType: core.ApiConfigSource_GRPC,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
				},
			}},
		},
	}

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: httpListenerCfg.RouteName,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name: util.Router,
		}},
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	return &v2.Listener{
		Name: httpListenerCfg.ListenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  localhost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: httpListenerCfg.ListenerPort,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name: util.HTTPConnectionManager,
				ConfigType: &listener.Filter_Config{
					Config: pbst,
				},
			}},
		}},
	}
}

// MakeTCPListener creates a TCP listener for a cluster.
func MakeTCPListener(listenerName string, port uint32, clusterName string) *v2.Listener {
	// TCP filter configuration
	config := &tcp.TcpProxy{
		StatPrefix: "tcp",
		ClusterSpecifier: &tcp.TcpProxy_Cluster{
			Cluster: clusterName,
		},
	}
	pbst, err := util.MessageToStruct(config)
	if err != nil {
		panic(err)
	}
	return &v2.Listener{
		Name: listenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  localhost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name: util.TCPProxy,
				ConfigType: &listener.Filter_Config{
					Config: pbst,
				},
			}},
		}},
	}
}
