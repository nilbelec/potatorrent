package crawler

import "regexp"

type Torrent struct {
	Calidad          string `json:"calidad"`
	GUID             string `json:"guid"`
	Imagen           string `json:"imagen"`
	TorrentDateAdded string `json:"torrentDateAdded"`
	TorrentID        string `json:"torrentID"`
	TorrentName      string `json:"torrentName"`
	TorrentSize      string `json:"torrentSize"`
}

func (t *Torrent) Season() string {
	re := regexp.MustCompile(`Temporada (\d+)`)
	match := re.FindStringSubmatch(t.TorrentName)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func (t *Torrent) FirstEpisode() string {
	s := t.Season()
	if s == "" {
		return ""
	}
	re := regexp.MustCompile(`Cap\.\s?` + s + `(\d+)`)
	match := re.FindStringSubmatch(t.TorrentName)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func (t *Torrent) LastEpisode() string {
	s := t.Season()
	if s == "" {
		return ""
	}
	re := regexp.MustCompile(`Cap\.\s?\d+_` + s + `(\d+)`)
	match := re.FindStringSubmatch(t.TorrentName)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
