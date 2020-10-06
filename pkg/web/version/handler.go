package version

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/version"

	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

type appVersion struct {
	Version string `json:"version"`
}

// Handler handles the search requests
type Handler struct {
	c *github.Client
}

// NewHandler creates a new search handler
func NewHandler(c *github.Client) *Handler {
	return &Handler{c}
}

// Routes return the routes the search handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/api/version", Method: "GET", HandlerFunc: h.getVersion},
		router.Route{Path: "/api/latest", Method: "GET", HandlerFunc: h.getLatest},
	}
}

func (h *Handler) getLatest(w http.ResponseWriter, r *http.Request) {
	l, err := h.c.LatestVersion()
	if err != nil {
		http.Error(w, "Error while checking latest version", http.StatusInternalServerError)
		return
	}
	version := appVersion{
		Version: l,
	}
	data, err := json.Marshal(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) getVersion(w http.ResponseWriter, r *http.Request) {
	version := appVersion{
		Version: version.Current,
	}
	data, err := json.Marshal(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
