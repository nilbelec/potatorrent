package main

import (
	"fmt"
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
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
		_ = fmt.Errorf("unsupported platform")
	}

}
