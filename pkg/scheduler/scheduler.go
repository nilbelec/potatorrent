package scheduler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	random "math/rand"
	"os"
	"strings"
	"time"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/util"
)

type ScheduleSearch struct {
	ID                  string                `json:"id"`
	LastTorrentID       string                `json:"lastTorrentID"`
	LastTorrentName     string                `json:"lastTorrentName"`
	LastTorrentImage    string                `json:"lastTorrentImage"`
	LastTorrentDate     string                `json:"lastTorrentDate"`
	LastTorrentPassword string                `json:"lastTorrentPassword"`
	LastExecution       time.Time             `json:"lastExecutionTime"`
	Folder              string                `json:"folder"`
	Interval            int                   `json:"interval"`
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
	c *crawler.Crawler
	f *SchedulesFile
}

// NewScheduler creates a new torrent scheduler
func NewScheduler(c *crawler.Crawler, f *SchedulesFile) *Scheduler {
	s := &Scheduler{c: c, f: f}
	go s.runAll()
	return s
}

func (s *Scheduler) GetAll() []*ScheduleSearch {
	return s.f.GetAll()
}

func (s *Scheduler) Add(ss *ScheduleSearch) error {
	err := s.validate(ss)
	if err != nil {
		return err
	}
	err = s.save(ss)
	if err != nil {
		return err
	}
	go s.run(ss.ID)
	return nil
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

func (s *Scheduler) validate(ss *ScheduleSearch) error {
	ss.Folder = strings.TrimSuffix(ss.Folder, "/")
	if _, err := os.Stat(ss.Folder); os.IsNotExist(err) {
		return errors.New("El directorio introducido no existe")
	}
	return nil
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
	for {
		ss := s.f.Get(id)
		if ss == nil {
			break
		}
		if ss.Disabled {
			time.Sleep(time.Duration(ss.Interval) * time.Minute)
			continue
		}
		fmt.Printf("[%s] - Running\n", ss)
		r, err := s.c.Search(ss.Params, "1")
		if err != nil {
			ss.Error = err.Error()
			fmt.Printf("[%s] - Error on Search: %v\n", ss, ss.Error)
			ss.LastExecution = time.Now()
			s.f.Save(ss)
			time.Sleep(time.Duration(ss.Interval) * time.Minute)
			continue
		}
		ss.Error = ""
		ts := r.Data.GetTorrents()
		for _, t := range ts {
			if ss.LastTorrentID == "NONE" || t.TorrentID == ss.LastTorrentID {
				ss.LastTorrentID = ""
				fmt.Printf("[%s] - No new torrents found\n", ss)
				break
			}
			fmt.Printf("[%s] - Found new torrent \"%s\"\n", ss, t.TorrentName)
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
		fmt.Printf("[%s] - Done for now\n", ss)
		time.Sleep(time.Duration(ss.Interval) * time.Minute)
	}
}

func (s *Scheduler) downloadTorrent(t *crawler.Torrent, ss *ScheduleSearch) error {
	result, err := s.c.Download(t.TorrentID, t.GUID)
	if err != nil {
		return err
	}
	if result == nil || result.Url == "" {
		return errors.New("No se ha encontrado la URL del torrent")
	}
	ss.LastTorrentPassword = result.Password
	client := util.NewHTTPClient()
	resp, err := client.Get(result.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(ss.Folder + "/" + t.TorrentID + ".torrent")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	if result.Password != "" {
		return ioutil.WriteFile(ss.Folder+"/"+t.TorrentID+".password.txt", []byte(result.Password), 0644)
	}
	return nil
}

func (s *Scheduler) runAll() error {
	ids, err := s.f.GetAllIDs()
	if err != nil {
		return err
	}
	for _, id := range ids {
		go s.run(id)
		time.Sleep(time.Duration(random.Intn(5)+5) * time.Second)
	}
	return nil
}

func randomUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
