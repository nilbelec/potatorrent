package router

import (
	"net/http"
)

// Router handles the routing of the web server
type Router struct {
	routes Routes
}

// New creates a new Router
func New() *Router {
	return &Router{}
}

// AddHandler adds the routes handler routes to the router
func (ro *Router) AddHandler(rh RoutesHandler) {
	for _, r := range rh.Routes() {
		ro.addRoute(r)
	}
}

func (ro *Router) addRoute(r Route) {
	ro.routes = append(ro.routes, r)
}

// ServeHTTP serves any request to the corresponding handler
func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, rt := range ro.routes {
		if rt.handlesRequest(r) {
			rt.HandlerFunc.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
