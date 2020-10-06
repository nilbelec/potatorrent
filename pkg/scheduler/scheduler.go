package scheduler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/downloader"
)

type ScheduleSearch struct {
	ID                  string                `json:"id"`
	LastTorrentID       string                `json:"lastTorrentID"`
	LastTorrentName     string                `json:"lastTorrentName"`
	LastTorrentImage    string                `json:"lastTorrentImage"`
	LastTorrentDate     string                `json:"lastTorrentDate"`
	LastTorrentPassword string                `json:"lastTorrentPassword"`
	LastExecution       time.Time             `json:"lastExecutionTime"`
	Params              *crawler.SearchParams `json:"params"`
	Error               string                `json:"error"`
	Disabled            bool                  `json:"disabled"`
}

func (ss *ScheduleSearch) String() string {
	p := ss.Params
	result := ""
	if p.CategoriaTexto != "" {
		result += p.CategoriaTexto
	}
	if p.SubCategoriaTexto != "" {
		result += " - " + p.SubCategoriaTexto
	}
	if p.CalidadTexto != "" {
		result += " - " + p.CalidadTexto
	}
	if p.Palabras != "" {
		result += " - \"" + p.Palabras + "\""
	}
	if result == "" {
		return "Cualquier torrent"
	}
	return strings.Trim(result, " - ")
}

// Scheduler is an Torrent scheduler
type Scheduler struct {
	c   *crawler.Crawler
	f   *SchedulesFile
	cfg *config.ConfigFile
	d   *downloader.Downloader
}

// NewScheduler creates a new torrent scheduler
func NewScheduler(c *crawler.Crawler, f *SchedulesFile, cfg *config.ConfigFile, d *downloader.Downloader) *Scheduler {
	s := &Scheduler{c, f, cfg, d}
	go s.runAll()
	return s
}

func (s *Scheduler) GetAll() []*ScheduleSearch {
	return s.f.GetAll()
}

func (s *Scheduler) Add(ss *ScheduleSearch) error {
	return s.save(ss)
}

func (s *Scheduler) Delete(id string) error {
	return s.f.Delete(id)
}

func (s *Scheduler) Disable(id string) error {
	ss := s.f.Get(id)
	if ss == nil {
		return errors.New("No se ha encontrado la búsqueda programada")
	}
	ss.Disabled = true
	return s.f.Save(ss)
}

func (s *Scheduler) Enable(id string) error {
	ss := s.f.Get(id)
	if ss == nil {
		return errors.New("No se ha encontrado la búsqueda programada")
	}
	ss.Disabled = false
	return s.f.Save(ss)
}

func (s *Scheduler) save(ss *ScheduleSearch) error {
	id, err := randomUUID()
	if err != nil {
		return err
	}
	ss.ID = id
	ss.LastTorrentID = "NONE"
	ss.Disabled = false
	return s.f.Save(ss)
}

func (s *Scheduler) run(id string) {
	ss := s.f.Get(id)
	if ss == nil || ss.Disabled {
		return
	}
	now := time.Now()
	next := ss.LastExecution.Add(time.Duration(s.cfg.IntervalInMinutes()) * time.Minute)
	if next.After(now) {
		return
	}
	log.Printf("[%s] - Running\n", ss)
	r, err := s.c.Search(ss.Params, "1")
	if err != nil {
		ss.Error = err.Error()
		log.Printf("[%s] - Error on Search: %v\n", ss, ss.Error)
		ss.LastExecution = time.Now()
		s.f.Save(ss)
		return
	}
	ss.Error = ""
	ts := r.Data.GetTorrents()
	for _, t := range ts {
		if ss.LastTorrentID == "NONE" || t.TorrentID == ss.LastTorrentID {
			ss.LastTorrentID = ""
			log.Printf("[%s] - No new torrents found\n", ss)
			break
		}
		log.Printf("[%s] - Found new torrent \"%s\"\n", ss, t.TorrentName)
		s.downloadTorrent(t, ss)
	}
	if len(ts) > 0 {
		ss.LastTorrentID = ts[0].TorrentID
		ss.LastTorrentName = ts[0].TorrentName
		ss.LastTorrentImage = ts[0].Imagen
		ss.LastTorrentDate = ts[0].TorrentDateAdded
	}
	ss.LastExecution = time.Now()
	s.f.Save(ss)
	log.Printf("[%s] - Done for now\n", ss)
}

func (s *Scheduler) downloadTorrent(t *crawler.Torrent, ss *ScheduleSearch) error {
	result, err := s.c.SearchTorrentInfo(t.TorrentID, t.TorrentDateAdded, t.GUID)
	if err != nil {
		return err
	}
	if result == nil || result.URL == "" {
		return errors.New("No se ha encontrado la URL del torrent")
	}
	ss.LastTorrentPassword = result.Password
	s.d.DownloadOn(t.TorrentID, result, s.cfg.DownloadFolder())
	return nil
}

func (s *Scheduler) runAll() error {
	for {
		ids, err := s.f.GetAllIDs()
		if err != nil {
			return err
		}
		for _, id := range ids {
			s.run(id)
		}
		time.Sleep(time.Duration(10 * time.Second))
	}
}

func randomUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
