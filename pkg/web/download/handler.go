package download

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/downloader"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the download requests
type Handler struct {
	c  *crawler.Crawler
	dw *downloader.Downloader
}

// NewHandler creates a new downloads handler
func NewHandler(c *crawler.Crawler, dw *downloader.Downloader) *Handler {
	return &Handler{c, dw}
}

// Routes return the routes the download handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/api/download", Method: "GET", HandlerFunc: h.downloadTorrent},
	}
}

func (h *Handler) downloadTorrent(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	id := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")
	folder := r.URL.Query().Get("folder") == "true"
	result, err := h.c.SearchTorrentInfo(id, date, guid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if folder {
		err := h.dw.Download(id, result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
