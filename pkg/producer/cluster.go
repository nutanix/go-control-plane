// Copyright 2018 Nutanix
//
// Author: utkarsh.simha@nutanix.com

package producer

import (
	"log"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
)

// ClusterConfig provides details of the configuration of a cluster definition
type ClusterConfig struct {
	// Name of the cluster
	ClusterName string

	// Name of (static) cluster for getting corresponding EDS
	EdsClusterName string

	// Timeout for connection
	ConnectTimeout time.Duration

	// Subset load balancing selector keys
	SubsetSelectorKeys []string
}

func (clusterCfg *ClusterConfig) validateStruct() bool {
	// Check for required fields
	if clusterCfg.ClusterName == "" ||
		clusterCfg.EdsClusterName == "" {
		return false
	}

	// Set default fields
	if clusterCfg.ConnectTimeout.String() == "" {
		clusterCfg.ConnectTimeout = 2 * time.Second
	}
	return true
}

// MakeCluster creates a cluster using config seeded by ClusterConfig
func (clusterCfg *ClusterConfig) MakeCluster() *v2.Cluster {
	if ok := clusterCfg.validateStruct(); !ok {
		log.Fatal("Could not validate ClusterConfig struct!")
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
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: clusterCfg.EdsClusterName},
					},
				}},
			},
		},
	}

	return &v2.Cluster{
		Name:           clusterCfg.ClusterName,
		ConnectTimeout: 2 * time.Second,
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
