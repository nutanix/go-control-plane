// Copyright 2018 Nutanix
//
// Author: utkarsh.simha

package producer

import (
	"log"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/gogo/protobuf/types"
)

// EndpointConfig defines the configuration for a single LB endpoint along with the corresponding metadata
type EndpointConfig struct {
	// Port number of the endpoint socket
	SocketPort uint32

	// IP address of the endpoint socket
	SocketAddress string

	// Subset load balancer key value map
	// Note: ensure keys are present in the subset selector keys provided to
	// corresponding cluster definition
	SubsetLBMap map[string]string

	// Reverese connection filter metadata
	RCMetadataMap map[string]string
}

func (endpointCfg *EndpointConfig) validateStruct() bool {
	// Check required fields
	if endpointCfg.SocketPort == 0 {
		return false
	}

	// Set default values
	if endpointCfg.SocketAddress == "" {
		endpointCfg.SocketAddress = "127.0.0.1"
	}
	return true
}

// ClusterLoadAssignment provides details of the load assignment for a given cluster, with all corresponding
// endpoints. Currently, ClusterLoadAssignment allows only LocalityLB Endpoints
type ClusterLoadAssignment struct {
	// Name of the cluster
	ClusterName string

	// List of endpoints (defined as EndpointConfig structs) that belong to the given cluster
	EndpointConfigs []*EndpointConfig
}

// MakeClusterLoadAssignment creates an endpoint using EndpointConfig
func (cla *ClusterLoadAssignment) MakeClusterLoadAssignment() *v2.ClusterLoadAssignment {
	endpointCfgs := cla.EndpointConfigs
	lbEndpoints := make([]endpoint.LbEndpoint, len(endpointCfgs))
	for _, endpointCfg := range endpointCfgs {
		if endpointCfg.SocketAddress == "" ||
			endpointCfg.SocketPort == 0 {
			// Invalid endpoint config struct
			log.Fatal("Could not validate ClusterConfig struct!")
			return nil
		}

		subsetLBFields := make(map[string]*types.Value)
		if len(endpointCfg.SubsetLBMap) != 0 {
			for key, val := range endpointCfg.SubsetLBMap {
				subsetLBFields[key] = &types.Value{Kind: &types.Value_StringValue{StringValue: val}}
			}
		}

		reverseConnFields := make(map[string]*types.Value)
		if len(endpointCfg.RCMetadataMap) != 0 {
			for key, val := range endpointCfg.RCMetadataMap {
				reverseConnFields[key] = &types.Value{Kind: &types.Value_StringValue{StringValue: val}}
			}
		}

		endpoint := endpoint.LbEndpoint{
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
					"envoy.lb":           &types.Struct{Fields: subsetLBFields},
					"envoy.reverse_conn": &types.Struct{Fields: reverseConnFields},
				},
			},
		}

		lbEndpoints = append(lbEndpoints, endpoint)

	}

	return &v2.ClusterLoadAssignment{
		ClusterName: cla.ClusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: lbEndpoints,
		}},
	}
}
