package configuration

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the download requests
type Handler struct {
	cfg *config.ConfigFile
}

// NewHandler creates a new configuration handler
func NewHandler(cfg *config.ConfigFile) *Handler {
	return &Handler{cfg}
}

// Routes return the routes the configuration handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/api/configuration", Method: "GET", HandlerFunc: h.getConfiguration},
		router.Route{Path: "/api/configuration", Method: "POST", HandlerFunc: h.saveConfiguration},
	}
}

func (h *Handler) getConfiguration(w http.ResponseWriter, r *http.Request) {
	configuration := h.cfg.GetConfiguration()
	response, err := json.Marshal(configuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) saveConfiguration(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	c := &config.Configuration{}
	err := d.Decode(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.cfg.SaveConfiguration(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
