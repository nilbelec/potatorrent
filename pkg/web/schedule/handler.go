package schedule

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web/router"
)

// Handler handles the schedule requests
type Handler struct {
	s *scheduler.Scheduler
}

// NewHandler creates a new schedule handler
func NewHandler(s *scheduler.Scheduler) *Handler {
	return &Handler{s}
}

// Routes return the routes the schedule handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/api/schedules", Method: "GET", HandlerFunc: h.getSchedules},
		router.Route{Path: "/api/schedule", Method: "POST", HandlerFunc: h.scheduleSearch},
		router.Route{Path: "/api/schedule", Method: "DELETE", HandlerFunc: h.deleteSchedule},
		router.Route{Path: "/api/schedule/disable", Method: "POST", HandlerFunc: h.disableSchedule},
		router.Route{Path: "/api/schedule/enable", Method: "POST", HandlerFunc: h.enableSchedule},
	}
}

func (h *Handler) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.s.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) disableSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.s.Disable(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) enableSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := h.s.Enable(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) scheduleSearch(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	j := &scheduler.ScheduleSearch{}
	err := d.Decode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.s.Add(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getSchedules(w http.ResponseWriter, r *http.Request) {
	ss := h.s.GetAll()
	response, err := json.Marshal(ss)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
