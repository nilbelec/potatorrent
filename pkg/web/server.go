package web

import (
	"log"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/downloader"
	"github.com/nilbelec/potatorrent/pkg/folders"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web/configuration"
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
	crawler    *crawler.Crawler
	github     *github.Client
	scheduler  *scheduler.Scheduler
	folders    *folders.Folders
	downloader *downloader.Downloader
	config     *config.ConfigFile
}

// NewServer creates a new web server
func NewServer(c *crawler.Crawler, g *github.Client, s *scheduler.Scheduler, f *folders.Folders, d *downloader.Downloader, cfg *config.ConfigFile) *Server {
	return &Server{c, g, s, f, d, cfg}
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
	r.AddHandler(download.NewHandler(s.crawler, s.downloader))
	r.AddHandler(version.NewHandler(s.github))
	r.AddHandler(schedule.NewHandler(s.scheduler))
	r.AddHandler(folder.NewHandler(s.folders))
	r.AddHandler(configuration.NewHandler(s.config))

	r.AddHandler(resources.NewHandler())
	return r
}
