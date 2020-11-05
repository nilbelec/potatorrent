package downloader

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/nilbelec/potatorrent/pkg/config"
	"github.com/nilbelec/potatorrent/pkg/crawler"
	"github.com/nilbelec/potatorrent/pkg/util"
)

type Downloader struct {
	c *config.ConfigFile
}

func NewDownloader(c *config.ConfigFile) *Downloader {
	return &Downloader{c}
}

func (d *Downloader) Download(id string, t *crawler.TorrentInfo) error {
	return d.DownloadOn(id, t, d.c.DownloadFolder())
}

func (d *Downloader) Reader(url string) (io.ReadCloser, error) {
	client := util.NewHTTPClient()
	resp, err := client.Get(url)
	return resp.Body, err
}

func (d *Downloader) DownloadOn(id string, t *crawler.TorrentInfo, folder string) error {
	respBody, err := d.Reader(t.URL)
	if err != nil {
		return err
	}
	defer respBody.Close()

	out, err := os.Create(folder + "/" + id + ".torrent")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, respBody)
	if err != nil {
		return err
	}
	if t.Password != "" {
		return ioutil.WriteFile(folder+"/"+id+""+"_password.txt", []byte(t.Password), 0644)
	}
	return nil
}
