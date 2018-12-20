package producer

import (
	"log"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

// RouteConfig is used to configure a particular route to a cluster
type RouteConfig struct {
	// Name of the route - referenced in listener config
	RouteName string

	// Name of the cluster to route to, or name of cluster header to use
	Cluster string

	// If the given cluster is a cluster header
	IsClusterHeader bool

	// A prefix to match with the URI or the regex to match with the URI
	PathMatch string

	// If the given path match is a regex
	IsPathMatchRegex bool
}

func (routeCfg *RouteConfig) validateStruct() bool {
	// Check required fields
	if routeCfg.RouteName == "" ||
		routeCfg.Cluster == "" {
		return false
	}

	// Set default values
	if routeCfg.PathMatch == "" {
		routeCfg.PathMatch = "/"
	}
	return true
}

// MakeRoute creates an HTTP route that routes to a given cluster.
func (routeCfg *RouteConfig) MakeRoute() *v2.RouteConfiguration {
	if ok := routeCfg.validateStruct(); !ok {
		log.Fatal("Could not validate RouteConfig struct!")
		return nil
	}

	// RouteAction can be `cluster_header` or `cluster`
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

	//PathMatch can be `prefix` or `regex`
	var routeMatch route.RouteMatch
	if routeCfg.IsPathMatchRegex {
		routeMatch = route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Regex{
				Regex: routeCfg.PathMatch,
			},
		}
	} else {
		routeMatch = route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: routeCfg.PathMatch,
			},
		}
	}

	return &v2.RouteConfiguration{
		Name: routeCfg.RouteName,
		VirtualHosts: []route.VirtualHost{{
			Name:    routeCfg.RouteName,
			Domains: []string{"*"},
			Routes: []route.Route{{
				Match: routeMatch,
				Action: &route.Route_Route{
					Route: routeAction,
				},
			}},
		}},
	}
}
