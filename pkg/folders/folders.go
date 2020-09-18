package folders

import (
	"os"
	"path/filepath"

	"github.com/nilbelec/potatorrent/pkg/config"
)

// FoldersResults is folder search results
type FoldersResults struct {
	Results []string `json:"results"`
}

type Folders struct {
	c *config.ConfigFile
}

func NewFolders(c *config.ConfigFile) *Folders {
	return &Folders{c}
}

func (f *Folders) Search(q string) *FoldersResults {
	r := f.searchFolders(q)
	return &FoldersResults{
		Results: r,
	}
}

func (f *Folders) searchFolders(q string) []string {
	r := make([]string, 0)
	if q == "" {
		r = append(r, f.c.DownloadFolder())
		return r
	}
	matches, err := filepath.Glob(q + "*")
	if err != nil {
		return make([]string, 0)
	}

	for _, p := range matches {
		f, err := os.Stat(p)
		if err != nil {
			continue
		}
		if f.IsDir() {
			abs, err := filepath.Abs(p)
			if err != nil {
				continue
			}
			r = append(r, abs)
		}
	}
	return r
}
