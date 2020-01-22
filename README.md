[![Build Status](https://travis-ci.com/nilbelec/potatorrent.svg?branch=master)](https://travis-ci.com/nilbelec/potatorrent)
[![Go Report Card](https://goreportcard.com/badge/github.com/nilbelec/potatorrent)](https://goreportcard.com/report/github.com/nilbelec/potatorrent)


<a href="https://www.buymeacoffee.com/nilbelec" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

# potatorrent
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nilbelec/potatorrent)](https://github.com/nilbelec/potatorrent/releases)

Simple web-application written in Go to download torrent files


## Build

```go
$ go get github.com/mitchellh/gox
$ go get github.com/gobuffalo/packr/packr
$ packr
$ go get -t -v ./...
$ gox -osarch="linux/arm linux/amd64 windows/amd64 darwin/amd64" -output="potatorrent.{{.OS}}.{{.Arch}}"
  -ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose ./...
```
