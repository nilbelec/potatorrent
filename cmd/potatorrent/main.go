package main

import (
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web"
	"github.com/pkg/browser"
)

func main() {
	crawler := crawler.NewCrawler()
	github := github.NewClient("nilbelec", "potatorrent")
	storage := scheduler.NewSchedulesFile("schedules.json")
	scheduler := scheduler.NewScheduler(crawler, storage)

	server := web.NewServer(crawler, github, scheduler)
	browser.OpenURL("http://localhost:8080")
	server.Start(":8080")
}
