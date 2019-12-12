package web

import (
	"log"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web/download"
	"github.com/nilbelec/potatorrent/pkg/web/folder"
	"github.com/nilbelec/potatorrent/pkg/web/image"
	"github.com/nilbelec/potatorrent/pkg/web/resources"
	"github.com/nilbelec/potatorrent/pkg/web/router"
	"github.com/nilbelec/potatorrent/pkg/web/schedule"
	"github.com/nilbelec/potatorrent/pkg/web/search"
	"github.com/nilbelec/potatorrent/pkg/web/version"
)

// Server web server
type Server struct {
	crawler   *crawler.Crawler
	github    *github.Client
	scheduler *scheduler.Scheduler
}

// NewServer creates a new web server
func NewServer(c *crawler.Crawler, g *github.Client, s *scheduler.Scheduler) *Server {
	return &Server{c, g, s}
}

// Start starts the web server
func (s *Server) Start(address string) {
	log.Println("Web server started at " + address)
	log.Fatal(http.ListenAndServe(address, s.router()))
}

func (s *Server) router() *router.Router {
	r := router.New()
	r.AddHandler(search.NewHandler(s.crawler))
	r.AddHandler(image.NewHandler(s.crawler))
	r.AddHandler(download.NewHandler(s.crawler))
	r.AddHandler(version.NewHandler(s.github))
	r.AddHandler(schedule.NewHandler(s.scheduler))
	r.AddHandler(folder.NewHandler())

	r.AddHandler(resources.NewHandler())
	return r
}
