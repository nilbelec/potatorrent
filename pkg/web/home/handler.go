package home

import (
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler is the home handler
type Handler struct {
	b packr.Box
}

// NewHandler creates a new home handler
func NewHandler() *Handler {
	b := packr.NewBox(".")
	return &Handler{b}
}

// Routes return the routes the home handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/", Method: "GET", Accepts: "text/html", HandlerFunc: h.getHomePage},
	}
}

func (h *Handler) getHomePage(w http.ResponseWriter, r *http.Request) {
	b, err := h.b.Find("home.html")
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	w.Write(b)
}
