package crawler

import "sort"

type SearchData struct {
	All      int                            `json:"all"`
	Items    int                            `json:"items"`
	Torrents map[string]map[string]*Torrent `json:"torrents"`
	Total    int                            `json:"total"`
}

type SearchParams struct {
	Categoria         string `json:"categoria"`
	CategoriaTexto    string `json:"categoriaTexto"`
	SubCategoria      string `json:"subcategoria"`
	SubCategoriaTexto string `json:"subcategoriaTexto"`
	Calidad           string `json:"calidad"`
	CalidadTexto      string `json:"calidadTexto"`
	Palabras          string `json:"q"`
}

type SearchResult struct {
	Data    SearchData `json:"data"`
	Success bool       `json:"success"`
}

type SearchOptions struct {
	Categoria map[string]string `json:"categorias"`
	Idioma    map[string]string `json:"idiomas"`
	Calidad   map[string]string `json:"calidades"`
	Ordenar   map[string]string `json:"ordenar"`
	Inon      map[string]string `json:"inon"`
}

type TorrentInfo struct {
	URL      string `json:"url"`
	Password string `json:"password"`
}

func (s *SearchData) GetTorrents() []*Torrent {
	return sortTorrents(s.Torrents)
}

func sortTorrents(m map[string]map[string]*Torrent) []*Torrent {
	result := make([]*Torrent, 0)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		keys2 := make([]string, 0, len(m[k]))
		for k2 := range m[k] {
			keys2 = append(keys2, k2)
		}
		sort.Strings(keys2)
		for _, k2 := range keys2 {
			result = append(result, m[k][k2])
		}
	}
	return result
}
