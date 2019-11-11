package image

import (
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the images requests
type Handler struct {
	c *crawler.Crawler
}

// NewHandler creates a new images handler
func NewHandler(c *crawler.Crawler) *Handler {
	return &Handler{c}
}

// Routes return the routes the images handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/image", Method: "GET", Accepts: "*/*", HandlerFunc: h.getImage},
	}
}

func (h *Handler) getImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	image, err := h.c.GetImage(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(image)
}
