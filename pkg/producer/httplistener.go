package producer

import (
	"log"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
)

// HTTPListenerConfig provides configuration for the HTTP listener
type HTTPListenerConfig struct {
	// ---- Required ----

	// Name of the HTTP listener
	ListenerName string

	// Port to bind the HTTP listener to
	ListenerPort uint32

	// Name of the route associated with the listener
	RouteName string

	// Name of (static) LDS cluster
	LdsClusterName string

	// ---- Optional ----

	// Address to bind the listener to
	// Defaults to "0.0.0.0"
	ListenerAddress string

	// Prefix for statistics
	// Defaults to "http"
	StatPrefix string
}

func (httpListenerCfg *HTTPListenerConfig) validateStruct() bool {
	// Check required fields
	if httpListenerCfg.ListenerName == "" ||
		httpListenerCfg.ListenerPort == 0 ||
		httpListenerCfg.RouteName == "" ||
		httpListenerCfg.LdsClusterName == "" {
		return false
	}

	// Set defaults if the fields are empty
	if httpListenerCfg.ListenerAddress == "" {
		httpListenerCfg.ListenerAddress = "0.0.0.0"
	}

	if httpListenerCfg.StatPrefix == "" {
		httpListenerCfg.StatPrefix = "http"
	}
	return true
}

// MakeHTTPListener creates a listener using either ADS or RDS for the route.
func (httpListenerCfg *HTTPListenerConfig) MakeHTTPListener() *v2.Listener {
	if ok := httpListenerCfg.validateStruct(); !ok {
		log.Fatal("Could not validate HTTPListenerConfig struct!")
		return nil
	}
	// #TODO(utkarsh.simha) Add TLS context
	// data source configuration
	rdsSource := core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			ApiType: core.ApiConfigSource_GRPC,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: httpListenerCfg.LdsClusterName},
				},
			}},
		},
	}

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: httpListenerCfg.StatPrefix,
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: httpListenerCfg.RouteName,
			},
		},
		// #TODO(utkarsh.simha) Allow for more filters + configs
		// Maybe make a struct?
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
					Address:  httpListenerCfg.ListenerAddress,
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
