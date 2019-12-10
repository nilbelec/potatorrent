package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web"
	"github.com/pkg/browser"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("No PORT environment variable found. Defaulting to port %s\n", port)
	}
	crawler := crawler.NewCrawler()
	github := github.NewClient("nilbelec", "potatorrent")
	storage := scheduler.NewSchedulesFile("schedules.json")
	scheduler := scheduler.NewScheduler(crawler, storage)

	server := web.NewServer(crawler, github, scheduler)
	browser.OpenURL(fmt.Sprintf("http://localhost:%s", port))
	server.Start(fmt.Sprintf(":%s", port))
}
