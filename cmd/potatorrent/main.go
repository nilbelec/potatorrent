package main

import (
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/web"
	"github.com/pkg/browser"
)

func main() {
	crawler := crawler.NewCrawler()
	github := github.NewClient("nilbelec", "potatorrent")
	server := web.NewServer(crawler, github)
	browser.OpenURL("http://localhost:8080")
	server.Start(":8080")
}
