package crawler

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Category struct
type Category struct {
	ID   string
	Name string
}

type Torrent struct {
	Calidad          string `json:"calidad"`
	GUID             string `json:"guid"`
	Imagen           string `json:"imagen"`
	TorrentDateAdded string `json:"torrentDateAdded"`
	TorrentID        string `json:"torrentID"`
	TorrentName      string `json:"torrentName"`
	TorrentSize      string `json:"torrentSize"`
}

type SearchData struct {
	All      int                           `json:"all"`
	Items    int                           `json:"items"`
	Torrents map[string]map[string]Torrent `json:"torrents"`
	Total    int                           `json:"total"`
}

func (s *SearchData) GetTorrents() []Torrent {
	var ts []Torrent
	for k := range s.Torrents {
		for _, t := range s.Torrents[k] {
			ts = append(ts, t)
		}
	}
	return ts
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

// Crawler is an Amazon crawler
type Crawler struct {
}

// NewCrawler creates a new Amazon product crawler
func NewCrawler() *Crawler {
	return &Crawler{}
}

const baseURLWithoutProtocol = "//descargas2020.org"
const baseURL = "https:" + baseURLWithoutProtocol
const searchPageURL = baseURL + "/buscar"
const searchURL = baseURL + "/get/result/"
const subcategoriesURL = baseURL + "/d20/library/include/ajax/get_subcategory.php"

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36"

// Search returns a search result
func (c *Crawler) Search(params url.Values) (*SearchResult, error) {
	resp, err := postForm(searchURL, params)
	if err != nil {
		return nil, errors.New("Error while requesting the search page: " + err.Error())
	}
	defer resp.Body.Close()
	r := SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	return &r, err
}

func (c *Crawler) Download(id string, path string) ([]byte, error) {
	url := baseURL + "/" + path
	bytes, err := findTorrent(id, url)
	if err == nil {
		return bytes, nil
	}
	return trySearchTorrent(id, url)
}

func trySearchTorrent(id string, url string) ([]byte, error) {
	pg := 1
	for {
		u := fmt.Sprint(url, "/pg/", pg)
		doc, err := getAndParse(u)
		if err != nil {
			return nil, errors.New("Error parsing request: " + err.Error())
		}
		lis := htmlquery.Find(doc, "//ul[@class=\"buscar-list\"]/li")
		if len(lis) == 0 {
			return nil, errors.New("No se ha encontrado el torrent")
		}
		for _, li := range lis {
			dp := extractLiDownloadPage(li)
			bytes, err := findTorrent(id, dp)
			if err == nil {
				return bytes, nil
			}
		}
		pg = pg + 1
	}
}

func extractLiDownloadPage(li *html.Node) string {
	a := htmlquery.FindOne(li, "/a")
	href, _ := findAttribute(a, "href")
	return href
}

func findTorrent(id string, url string) ([]byte, error) {
	println("Searching for torrent " + id + " in " + url)
	text, err := getString(url)
	if err != nil {
		return nil, errors.New("Error parsing request: " + err.Error())
	}
	re := regexp.MustCompile("\".+descargar-torrent\\/" + id + ".+\"")
	match := re.FindStringSubmatch(text)
	if len(match) == 0 {
		return nil, errors.New("Unable to find the download link")
	}
	return getBytes("https:" + strings.Trim(match[0], "\""))
}

func (c *Crawler) SearchOptions() (*SearchOptions, error) {
	doc, err := getAndParse(searchPageURL)
	if err != nil {
		return nil, errors.New("Error parsing request: " + err.Error())
	}
	return &SearchOptions{
		Categoria: extractCategorias(doc),
		Calidad:   extractCalidades(doc),
		Idioma:    extractIdiomas(doc),
		Ordenar:   extractOrden(doc),
		Inon:      extractInon(doc),
	}, nil
}

func extractCategorias(doc *html.Node) map[string]string {
	options := htmlquery.Find(doc, "//select[@id=\"categoryIDR\"]/option")
	return optionsToMap(options)
}

func extractInon(doc *html.Node) map[string]string {
	options := htmlquery.Find(doc, "//select[@id=\"inon\"]/option")
	return optionsToMap(options)
}

func extractOrden(doc *html.Node) map[string]string {
	options := htmlquery.Find(doc, "//select[@id=\"ordenar\"]/option")
	return optionsToMap(options)
}

func extractIdiomas(doc *html.Node) map[string]string {
	options := htmlquery.Find(doc, "//select[@id=\"idioma\"]/option")
	return optionsToMap(options)
}

func extractCalidades(doc *html.Node) map[string]string {
	options := htmlquery.Find(doc, "//select[@id=\"calidad\"]/option")
	return optionsToMap(options)
}

// GetImage returns a image
func (c *Crawler) GetImage(path string) ([]byte, error) {
	return getBytes(baseURL + path)
}

func getBytes(url string) ([]byte, error) {
	res, err := get(url)
	if err != nil {
		return nil, errors.New("Error while requesting the page: " + err.Error())
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// GetSubcategories returns the subcategories
func (c *Crawler) GetSubcategories(category string) (map[string]string, error) {
	res, err := postForm(subcategoriesURL, url.Values{"categoryIDR": {category}})
	if err != nil {
		return nil, errors.New("Error while requesting the page: " + err.Error())
	}
	options, err := parseFragment(res)
	if err != nil {
		return nil, errors.New("Error while parsing response: " + err.Error())
	}
	return optionsToMap(options), nil
}

func optionsToMap(options []*html.Node) map[string]string {
	m := make(map[string]string)
	for _, option := range options {
		value, ok := findAttribute(option, "value")
		if ok && value != "" {
			m[value] = option.FirstChild.Data
		}
	}
	return m
}

func findAttribute(element *html.Node, name string) (string, bool) {
	for _, a := range element.Attr {
		if a.Key == name {
			return a.Val, true
		}
	}
	return "", false
}

func getString(url string) (string, error) {
	res, err := get(url)
	if err != nil {
		return "", errors.New("Error while requesting the page: " + err.Error())
	}
	if res.StatusCode == 404 {
		return "", errors.New("No se ha encontrado la página")
	}
	if res.StatusCode != 200 {
		return "", errors.New("La página no ha podido ser cargada")
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes), nil
}

func getAndParse(url string) (*html.Node, error) {
	res, err := get(url)
	if err != nil {
		return nil, errors.New("Error while requesting the page: " + err.Error())
	}
	return parse(res)
}

func parse(resp *http.Response) (*html.Node, error) {
	defer resp.Body.Close()
	return html.Parse(resp.Body)
}

func parseFragment(resp *http.Response) ([]*html.Node, error) {
	defer resp.Body.Close()
	ctx := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "div",
	}
	return html.ParseFragment(resp.Body, ctx)
}

func get(url string) (resp *http.Response, err error) {
	client := getClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return client.Do(req)
}

func postForm(url string, data url.Values) (resp *http.Response, err error) {
	client := getClient()
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	return client.Do(req)
}

func getClient() http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return http.Client{Transport: transport}
}
