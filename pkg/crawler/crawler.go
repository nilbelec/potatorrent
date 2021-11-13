package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/nilbelec/potatorrent/pkg/config"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Crawler is an Torrent crawler
type Crawler struct {
	cfg *config.ConfigFile
}

// NewCrawler creates a new torrent crawler
func NewCrawler(cfg *config.ConfigFile) *Crawler {
	return &Crawler{cfg}
}

const searchPageURLPath = "/buscar"
const searchURLPath = "/get/result/"
const subcategoriesURLPath = "/pctn/library/include/ajax/get_subcategory.php"

func (c *Crawler) baseURLScheme() string {
	s := c.cfg.BaseURL()
	u, _ := url.Parse(s)
	return u.Scheme
}
func (c *Crawler) searchPageURL() string {
	return c.cfg.BaseURL() + searchPageURLPath
}

func (c *Crawler) searchURL() string {
	return c.cfg.BaseURL() + searchURLPath
}

func (c *Crawler) subcategoriesURL() string {
	return c.cfg.BaseURL() + subcategoriesURLPath
}

// Search returns a search result
func (c *Crawler) Search(params *SearchParams, page string) (*SearchResult, error) {
	resp, err := postForm(c.searchURL(), url.Values{
		"categoryIDR": {params.Categoria},
		"categoryID":  {params.SubCategoria},
		"idioma":      {},
		"calidad":     {params.Calidad},
		"ordenar":     {},
		"inon":        {},
		"s":           {params.Palabras},
		"pg":          {page},
	})
	if err != nil {
		log.Println("Error while requesting the search page: " + err.Error())
		return nil, errors.New("Error while requesting the search page: " + err.Error())
	}
	defer resp.Body.Close()
	r := SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Println("Error while decoding JSON search response: " + err.Error())
		return nil, errors.New("Error while decoding JSON search response: " + err.Error())
	}
	return &r, err
}

func (c *Crawler) SearchTorrentInfo(id string, date string, path string, season string, firstEpisode string, lastEpisode string) (*TorrentInfo, error) {
	path = strings.ReplaceAll(path, "--", "-")
	if season != "" && firstEpisode != "" {
		result, err := findTorrentEpisode(c.cfg.BaseURL(), id, path, season, firstEpisode, lastEpisode)
		if err == nil {
			return result, nil
		}
	}
	url := c.cfg.BaseURL() + "/" + path
	urlScheme := strings.Split(url, ":")[0]
	result, err := findTorrent(id, url, false, urlScheme, c.cfg.BaseURL(), "")
	if err == nil {
		return result, nil
	}
	return trySearchTorrent(id, date, season, firstEpisode, url, c.cfg.BaseURL())
}

func findTorrentEpisode(baseURL string, id string, path string, season string, firstEpisode string, lastEpisode string) (*TorrentInfo, error) {
	url := baseURL + "/"
	split := strings.Split(path, "/")
	if len(split) != 3 {
		log.Println("No valid path for torrent episode for " + id + " with path " + path)
		return nil, errors.New("No valid path for torrent episode")
	}
	if split[0] == "series-hd" {
		url += "descargar/serie-en-hd/"
	} else {
		url += "descargar/serie/"
	}
	url += split[1] + "/temporada-" + season + "/capitulo-"
	if len(firstEpisode) == 1 {
		url += "0"
	}
	url += firstEpisode
	if lastEpisode != "" {
		url += "-al-"
		if len(lastEpisode) == 1 {
			url += "0"
		}
		url += lastEpisode
	}
	urlScheme := strings.Split(url, ":")[0]
	result, err := findTorrent(id, url+"/", false, urlScheme, baseURL, "")
	if err == nil {
		return result, nil
	}
	return findTorrent(id, url, false, urlScheme, baseURL, "")
}

func trySearchTorrent(id string, date string, season string, firstEpisode string, url string, baseURL string) (*TorrentInfo, error) {
	pg := 1
	urlScheme := strings.Split(url, ":")[0]
	for {
		u := fmt.Sprint(url, "/pg/", pg)
		log.Println("Searching for torrent " + id + " in " + u)
		doc, err := getAndParse(u)
		if err != nil {
			log.Println("Error parsing request: " + err.Error())
			return nil, errors.New("Error parsing request: " + err.Error())
		}
		any := htmlquery.Find(doc, "//ul[@class=\"buscar-list\"]/li")
		if len(any) == 0 {
			log.Println("No episodes found for " + id + " in " + u)
			return &TorrentInfo{}, errors.New("No episodes found for " + id + " in " + u)
		}
		lis := htmlquery.Find(doc, "//ul[@class=\"buscar-list\"]/li[.//span[contains(text(),'"+strings.ReplaceAll(date, "/", "-")+"')]]")
		for _, li := range lis {
			dp := extractLiDownloadPage(li)
			bytes, err := findTorrent(id, dp, true, urlScheme, baseURL, "")
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

func prepareTorrentLinkRegexes(id string, strict bool) []string {
	var regexes []string
	if strict {
		regexes = append(regexes, "\".+/download\\/"+id+".+\"")
		regexes = append(regexes, "\".+/descargar-torrent\\/"+id+".+\"")
	} else {
		regexes = append(regexes, "\".+/download\\/.+\"")
		regexes = append(regexes, "\".+/descargar-torrent\\/.+\"")
	}
	return regexes
}

func matchAnyRegexes(regexes []string, text string) (string, error) {
	for _, rex := range regexes {
		re := regexp.MustCompile(rex)
		match := re.FindStringSubmatch(text)
		if len(match) > 0 {
			return match[0], nil
		}
	}
	return "", errors.New("No match found")
}

func findTorrent(id string, url string, strict bool, URLScheme string, baseURL string, password string) (*TorrentInfo, error) {
	log.Println("Searching for torrent " + id + " in " + url)
	text, err := getString(url)
	if err != nil {
		return nil, errors.New("Error parsing request: " + err.Error())
	}
	rex := "\"http.*?/download-link\\/" + id + ".*?\""
	re := regexp.MustCompile(rex)
	log.Println("Searching torrent link for " + id + " in " + url + " using the regex " + rex)
	match := re.FindStringSubmatch(text)
	if len(match) > 0 {
		r := &TorrentInfo{
			URL:      strings.Trim(match[0], "\""),
			Password: password,
		}
		log.Println("Found torrent file for " + id + ": " + r.URL)
		return r, nil
	}
	re = regexp.MustCompile("\"/descargar\\/torrent\\/.+?\"")
	match = re.FindStringSubmatch(text)
	if len(match) > 0 {
		log.Println(match)
		newURL := baseURL + strings.ReplaceAll(match[0], "\"", "")
		return findTorrent(id, newURL, strict, URLScheme, baseURL, "")
	}
	regexes := prepareTorrentLinkRegexes(id, strict)
	for _, rex := range regexes {
		log.Println("Searching torrent link for " + id + " in " + url + " using the regex " + rex)
		match, err := matchAnyRegexes(regexes, text)
		if err != nil {
			log.Println("Failed finding torrent link for " + id + " in " + url + " using the regex " + rex)
			continue
		}
		newURL := URLScheme + ":" + strings.Trim(match, "\"")
		p := findPassword(url)
		return findTorrent(id, newURL, strict, URLScheme, baseURL, p)
	}
	return nil, errors.New("Unable to find the download link")
}

func findPassword(url string) string {
	doc, err := getAndParse(url)
	if err != nil {
		return ""
	}
	input := htmlquery.FindOne(doc, "//input[@id=\"txt_password\"]")
	if input == nil {
		return ""
	}
	password, _ := findAttribute(input, "value")
	return password
}

func (c *Crawler) SearchOptions() (*SearchOptions, error) {
	doc, err := getAndParse(c.searchPageURL())
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

// GetImage returns an image
func (c *Crawler) GetImage(path string) ([]byte, error) {
	return getBytes(c.cfg.BaseURL() + path)
}

// GetSubcategories returns the subcategories
func (c *Crawler) GetSubcategories(category string) (map[string]string, error) {
	res, err := postForm(c.subcategoriesURL(), url.Values{"categoryIDR": {category}})
	if err != nil {
		return nil, errors.New("Error while requesting subcategories: " + err.Error())
	}
	options, err := parseFragment(res)
	if err != nil {
		return nil, errors.New("Error while parsing subcategories response: " + err.Error())
	}
	return optionsToMap(options), nil
}
