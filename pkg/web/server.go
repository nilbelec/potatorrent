package web

import (
	"log"
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/web/download"
	"github.com/nilbelec/potatorrent/pkg/web/home"
	"github.com/nilbelec/potatorrent/pkg/web/image"
	"github.com/nilbelec/potatorrent/pkg/web/router"
	"github.com/nilbelec/potatorrent/pkg/web/search"
)

// Server web server
type Server struct {
	crawler *crawler.Crawler
}

// NewServer creates a new web server
func NewServer(c *crawler.Crawler) *Server {
	return &Server{c}
}

// Start starts the web server
func (s *Server) Start(address string) {
	log.Println("Web server started at " + address)
	log.Fatal(http.ListenAndServe(address, s.router()))
}

func (s *Server) router() *router.Router {
	r := router.New()
	r.AddHandler(home.NewHandler())
	r.AddHandler(search.NewHandler(s.crawler))
	r.AddHandler(image.NewHandler(s.crawler))
	r.AddHandler(download.NewHandler(s.crawler))
	return r
}
