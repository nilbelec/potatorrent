package crawler

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/nilbelec/potatorrent/pkg/util"
	"golang.org/x/net/html"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36"

func getBytes(url string) ([]byte, error) {
	res, err := get(url)
	if err != nil {
		return nil, errors.New("Error while requesting the page \"" + url + "\": " + err.Error())
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func getAndParse(url string) (*html.Node, error) {
	res, err := get(url)
	if err != nil {
		return nil, errors.New("Error while requesting the page \"" + url + "\": " + err.Error())
	}
	return parse(res)
}

func getString(url string) (string, error) {
	res, err := get(url)
	if err != nil {
		return "", errors.New("Error while requesting the page \"" + url + "\": " + err.Error())
	}
	if res.StatusCode == 404 {
		return "", errors.New("Page \"" + url + "\" not found")
	}
	if res.StatusCode != 200 {
		return "", errors.New("Error while requesting the page \"" + url + "\". StatusCode: " + string(res.StatusCode))
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error while requesting the page \"" + url + "\": " + err.Error())
		return "", errors.New("Error while requesting the page \"" + url + "\": " + err.Error())
	}
	return string(bodyBytes), nil
}

func get(url string) (resp *http.Response, err error) {
	client := util.NewHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	return client.Do(req)
}

func postForm(url string, data url.Values) (resp *http.Response, err error) {
	client := util.NewHTTPClient()
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	return client.Do(req)
}
