package search

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the search requests
type Handler struct {
	c *crawler.Crawler
}

// NewHandler creates a new search handler
func NewHandler(c *crawler.Crawler) *Handler {
	return &Handler{c}
}

// Routes return the routes the search handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/search", Method: "GET", Accepts: "application/json", HandlerFunc: h.searchTorrents},
		router.Route{Path: "/options", Method: "GET", Accepts: "application/json", HandlerFunc: h.searchOptions},
		router.Route{Path: "/subcategories", Method: "GET", Accepts: "application/json", HandlerFunc: h.searchSubcategories},
	}
}

func (h *Handler) searchTorrents(w http.ResponseWriter, r *http.Request) {
	params := &crawler.SearchParams{
		Categoria:    r.URL.Query().Get("categoria"),
		SubCategoria: r.URL.Query().Get("subcategoria"),
		Calidad:      r.URL.Query().Get("calidad"),
		Palabras:     r.URL.Query().Get("q"),
	}
	sr, err := h.c.Search(params, r.URL.Query().Get("pg"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) searchOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.c.SearchOptions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(options)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) searchSubcategories(w http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("categoria")
	subcategorias, err := h.c.GetSubcategories(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(subcategorias)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
