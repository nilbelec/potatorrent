package folders

import (
	"os"
	"path/filepath"
)

// FoldersResults is folder search results
type FoldersResults struct {
	Results []*FoldersResult `json:"results"`
}

// FoldersResult is folder search result
type FoldersResult struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func Search(q string) *FoldersResults {
	r := searchFolders(q)
	return &FoldersResults{
		Results: r,
	}
}

func searchFolders(q string) []*FoldersResult {
	matches, err := filepath.Glob(q + "*")
	if err != nil {
		return make([]*FoldersResult, 0)
	}

	r := make([]*FoldersResult, 0)
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
			r = append(r, &FoldersResult{
				ID:   abs,
				Text: abs,
			})
		}
	}
	return r
}
