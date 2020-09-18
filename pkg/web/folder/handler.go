package folder

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/folders"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the folder requests
type Handler struct {
	f *folders.Folders
}

// NewHandler creates a new folder handler
func NewHandler(f *folders.Folders) *Handler {
	return &Handler{f: f}
}

// Routes return the routes the schedule handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/folder", Method: "GET", HandlerFunc: h.searchFolder},
	}
}

func (h *Handler) searchFolder(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	results := h.f.Search(q)
	response, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
