package folders

import (
	"os"
	"path/filepath"
)

// FoldersResults is folder search results
type FoldersResults struct {
	Results []string `json:"results"`
}

func Search(q string) *FoldersResults {
	r := searchFolders(q)
	return &FoldersResults{
		Results: r,
	}
}

func searchFolders(q string) []string {
	r := make([]string, 0)
	if q == "" {
		abs, err := filepath.Abs(".")
		if err != nil {
			return make([]string, 0)
		}
		r = append(r, abs)
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
