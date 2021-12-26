package resources

import (
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler is the resources handler
type Handler struct {
	b *packr.Box
}

// NewHandler creates a new resources handler
func NewHandler() *Handler {
	b := packr.New("p", "public")
	return &Handler{b}
}

// Routes return the routes the home handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/", Method: "GET", Accepts: "text/html", HandlerFunc: h.getHomePage},
		router.Route{Path: "/searches", Method: "GET", Accepts: "text/html", HandlerFunc: h.getHomePage},
		router.Route{Path: "/configuration", Method: "GET", Accepts: "text/html", HandlerFunc: h.getHomePage},
		router.Route{Pattern: regexp.MustCompile(`/public/.+`), Method: "GET", HandlerFunc: h.getResource},
	}
}

func (h *Handler) getHomePage(w http.ResponseWriter, r *http.Request) {
	b, err := h.b.Find("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	w.Write(b)
}

func (h *Handler) getResource(w http.ResponseWriter, r *http.Request) {
	b, err := h.b.Find(strings.Replace(r.URL.Path, "/public", "", 1))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	w.Write(b)
}
