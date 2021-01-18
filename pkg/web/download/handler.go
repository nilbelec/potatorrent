package download

import (
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"

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
		router.Route{Path: "/api/download/info", Method: "GET", HandlerFunc: h.downloadInfo},
		router.Route{Path: "/api/download/file", Method: "GET", HandlerFunc: h.downloadFile},
		router.Route{Path: "/api/download/onFolder", Method: "GET", HandlerFunc: h.downloadFileOnFolder},
	}
}

func (h *Handler) downloadInfo(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	id := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")
	season := r.URL.Query().Get("season")
	firstEpisode := r.URL.Query().Get("firstEpisode")
	lastEpisode := r.URL.Query().Get("lastEpisode")
	result, err := h.c.SearchTorrentInfo(id, date, guid, season, firstEpisode, lastEpisode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) downloadFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	reader, err := h.dw.Reader(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Base(url)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
	w.Header().Set("Content-Type", "application/x-bittorrent")
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) downloadFileOnFolder(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	id := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")
	season := r.URL.Query().Get("season")
	firstEpisode := r.URL.Query().Get("firstEpisode")
	lastEpisode := r.URL.Query().Get("lastEpisode")
	result, err := h.c.SearchTorrentInfo(id, date, guid, season, firstEpisode, lastEpisode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.dw.Download(id, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
