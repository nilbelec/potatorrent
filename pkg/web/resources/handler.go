package resources

import (
	"mime"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/gobuffalo/packr"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler is the public handler
type Handler struct {
	b packr.Box
}

// NewHandler creates a new public handler
func NewHandler() *Handler {
	b := packr.NewBox(".")
	return &Handler{b}
}

// Routes return the routes the home handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Pattern: regexp.MustCompile(`/public/.+`), Method: "GET", HandlerFunc: h.getResource},
	}
}

func (h *Handler) getResource(w http.ResponseWriter, r *http.Request) {
	b, err := h.b.Find(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	w.Write(b)
}
