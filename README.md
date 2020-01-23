# potatorrent [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nilbelec/potatorrent)](https://github.com/nilbelec/potatorrent/releases) [![Build Status](https://travis-ci.com/nilbelec/potatorrent.svg?branch=master)](https://travis-ci.com/nilbelec/potatorrent) [![Go Report Card](https://goreportcard.com/badge/github.com/nilbelec/potatorrent)](https://goreportcard.com/report/github.com/nilbelec/potatorrent)

Scraper application with web interface to help you download torrent files from -spanish- websites full of annoying ads.

It also allows you to schedule your searches to automatically download new torrents.

If you like it and want to support development, you can buy me a coffee! :)

<a href="https://www.buymeacoffee.com/nilbelec" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

## Usage

Just grab one of the pre-built executable files for your OS from the [releases](https://github.com/nilbelec/potatorrent/releases) page and launch it. It should open a new web browser tab with the application.

If no tab is opened, open your web browser and navigate to:

```
http://localhost:8080
```

By default it will be launch using port 8080, but you can change it setting the PORT environment variable before launch.

## Build

Potatorrent was develop using [go](https://golang.org/). To build it just clone the code and run:

```bash
$ cd cmd/potatorrent
$ go get
$ go built
```
