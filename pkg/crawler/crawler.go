package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// Crawler is an Torrent crawler
type Crawler struct {
}

// NewCrawler creates a new torrent crawler
func NewCrawler() *Crawler {
	return &Crawler{}
}

const protocol = "https:"
const baseURLWithoutProtocol = "//" + "pc" + "tmi" + "x.com"
const baseURL = protocol + baseURLWithoutProtocol
const searchPageURL = baseURL + "/buscar"
const searchURL = baseURL + "/get/result/"
const subcategoriesURL = baseURL + "/pctn/library/include/ajax/get_subcategory.php"

// Search returns a search result
func (c *Crawler) Search(params *SearchParams, page string) (*SearchResult, error) {
	resp, err := postForm(searchURL, url.Values{
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

func (c *Crawler) Download(id string, date string, path string) (*SearchTorrentResult, error) {
	url := baseURL + "/" + path
	result, err := findTorrent(id, url, false)
	if err == nil {
		return result, nil
	}
	return trySearchTorrent(id, date, url)
}

func trySearchTorrent(id string, date string, url string) (*SearchTorrentResult, error) {
	pg := 1
	for {
		u := fmt.Sprint(url, "/pg/", pg)
		log.Println("Searching for torrent " + id + " in " + u)
		doc, err := getAndParse(u)
		if err != nil {
			return nil, errors.New("Error parsing request: " + err.Error())
		}
		any := htmlquery.Find(doc, "//ul[@class=\"buscar-list\"]/li")
		if len(any) == 0 {
			return &SearchTorrentResult{}, nil
		}
		lis := htmlquery.Find(doc, "//ul[@class=\"buscar-list\"]/li[.//span[contains(text(),'"+strings.ReplaceAll(date, "/", "-")+"')]]")
		for _, li := range lis {
			dp := extractLiDownloadPage(li)
			bytes, err := findTorrent(id, dp, true)
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

func findTorrent(id string, url string, strict bool) (*SearchTorrentResult, error) {
	log.Println("Searching for torrent " + id + " in " + url)
	text, err := getString(url)
	if err != nil {
		return nil, errors.New("Error parsing request: " + err.Error())
	}
	rex := "\".+descargar-torrent\\/.+\""
	if strict {
		rex = "\".+descargar-torrent\\/" + id + ".+\""
	}
	re := regexp.MustCompile(rex)
	match := re.FindStringSubmatch(text)
	if len(match) == 0 {
		return nil, errors.New("Unable to find the download link")
	}
	r := &SearchTorrentResult{
		Url:      protocol + strings.Trim(match[0], "\""),
		Password: findPassword(url),
	}
	log.Println("Found torrent file for " + id + ": " + r.Url)
	return r, nil
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

// GetSubcategories returns the subcategories
func (c *Crawler) GetSubcategories(category string) (map[string]string, error) {
	res, err := postForm(subcategoriesURL, url.Values{"categoryIDR": {category}})
	if err != nil {
		return nil, errors.New("Error while requesting subcategories: " + err.Error())
	}
	options, err := parseFragment(res)
	if err != nil {
		return nil, errors.New("Error while parsing subcategories response: " + err.Error())
	}
	return optionsToMap(options), nil
}
