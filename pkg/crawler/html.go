package crawler

import (
	"net/http"

	"github.com/nilbelec/potatorrent/pkg/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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

func findAttribute(element *html.Node, name string) (string, bool) {
	for _, a := range element.Attr {
		if a.Key == name {
			return util.ToUtf8(a.Val), true
		}
	}
	return "", false
}

func optionsToMap(options []*html.Node) map[string]string {
	m := make(map[string]string)
	for _, option := range options {
		value, ok := findAttribute(option, "value")
		if ok && value != "" {
			m[value] = util.ToUtf8(option.FirstChild.Data)
		}
	}
	return m
}
