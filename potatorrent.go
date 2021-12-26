package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/downloader"
	"github.com/nilbelec/potatorrent/pkg/folders"
	"github.com/nilbelec/potatorrent/pkg/github"
	"github.com/nilbelec/potatorrent/pkg/scheduler"
	"github.com/nilbelec/potatorrent/pkg/web"
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
	openbrowser(fmt.Sprintf("http://localhost:%s", port))
	server.Start(fmt.Sprintf(":%s", port))
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
