package router

import (
	"net/http"
	"regexp"
	"strings"
)

// Route is the mapping to a request made to the server
type Route struct {
	Path        string
	Pattern     *regexp.Regexp
	Method      string
	Accepts     string
	HandlerFunc http.HandlerFunc
}

// RoutesHandler is any object that has routes
type RoutesHandler interface {
	Routes() Routes
}

// Routes is a slice of Routes
type Routes []Route

func (rt *Route) handlesRequest(r *http.Request) bool {
	if rt.Path == "" && rt.Pattern == nil {
		return false
	}
	if rt.Path != "" && rt.Path != r.URL.Path {
		return false
	}
	if rt.Pattern != nil && !rt.Pattern.MatchString(r.URL.Path) {
		return false
	}
	if rt.Method != "" && rt.Method != r.Method {
		return false
	}
	if rt.Accepts != "" && !strings.Contains(r.Header.Get("Accept"), rt.Accepts) {
		return false
	}
	return true
}
