package main

import (
	"fmt"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/downloader"
	"github.com/nilbelec/potatorrent/pkg/folders"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web"
	"github.com/pkg/browser"
)

func main() {
	config := config.NewConfigFile("config.json")
	port := config.Port()
	crawler := crawler.NewCrawler(config)
	github := github.NewClient("nilbelec", "potatorrent")
	storage := scheduler.NewSchedulesFile("schedules.json")
	downloader := downloader.NewDownloader(config)
	scheduler := scheduler.NewScheduler(crawler, storage, config, downloader)
	folders := folders.NewFolders(config)

	server := web.NewServer(crawler, github, scheduler, folders, downloader, config)
	browser.OpenURL(fmt.Sprintf("http://localhost:%s", port))
	server.Start(fmt.Sprintf(":%s", port))
}
