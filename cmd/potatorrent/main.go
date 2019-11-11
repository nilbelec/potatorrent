package main

import (
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/web"
	"github.com/pkg/browser"
)

func main() {
	crawler := crawler.NewCrawler()
	server := web.NewServer(crawler)
	browser.OpenURL("http://localhost:8080")
	server.Start(":8080")
}
