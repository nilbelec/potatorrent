package crawler

type Torrent struct {
	Calidad          string `json:"calidad"`
	GUID             string `json:"guid"`
	Imagen           string `json:"imagen"`
	TorrentDateAdded string `json:"torrentDateAdded"`
	TorrentID        string `json:"torrentID"`
	TorrentName      string `json:"torrentName"`
	TorrentSize      string `json:"torrentSize"`
}
