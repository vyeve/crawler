package crawler

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	goQuery "github.com/PuerkitoBio/goquery"
)

func (ci *crawlerImp) walkLinks(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(path)
	if err != nil {
		// TODO: should we log errors???
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		// bad status
		return
	}
	defer resp.Body.Close() // nolint: errcheck

	doc, err := goQuery.NewDocumentFromReader(resp.Body)
	if err != nil {
		ci.siteCh <- linkWrapper{
			rootURL: path,
		}
		return
	}

	for element, attributes := range htmlElements {
		for _, attr := range attributes {
			ci.lookup(doc, element, attr, path, wg)
		}
	}
}

func (ci *crawlerImp) lookup(doc *goQuery.Document, element, attr, path string, wg *sync.WaitGroup) {
	doc.Find(element).Each(func(i int, sel *goQuery.Selection) {
		name, _ := sel.Attr(attr)
		if !validateLink(name) {
			return
		}
		path = strings.TrimSuffix(path, "/")
		ci.siteCh <- linkWrapper{
			rootURL: path,
			linkURL: name,
		}
		if ci.allowedToProcess(name) {
			wg.Add(1)
			go ci.walkLinks(name, wg)
		}
	})
}

// allowedToProcess is like a lock, to prevent duplicate visiting the same link
func (ci *crawlerImp) allowedToProcess(path string) /*allowed*/ bool {
	link, err := url.Parse(path)
	if err != nil {
		return false
	}
	if strings.Index(link.Path, ".") > 0 {
		// means that path contains link to static file
		return false
	}
	path = link.Host + link.Path

	if !strings.HasPrefix(strings.TrimPrefix(path, wwwKey), ci.domain) {
		// host is not the same as domain
		return false
	}
	return ci.setupLink(path)
}

func (ci *crawlerImp) setupLink(path string) bool {
	ci.Lock()
	defer ci.Unlock()
	visited := ci.visited[path]
	if visited {
		// site already visited
		return false
	}
	ci.visited[path] = true
	if !strings.HasPrefix(path, wwwKey) {
		// to prevent visit link "google.com" in case "www.google.com" already visited
		return !ci.visited[wwwKey+path]
	}

	return true
}
