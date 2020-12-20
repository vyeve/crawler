package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	goQuery "github.com/PuerkitoBio/goquery"
)

func (cr *crawlerImp) walkLinks(path string, wg *sync.WaitGroup) {
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
	cr.parseHTML(resp.Body, wg, path)
}

func (cr *crawlerImp) parseHTML(body io.ReadCloser, wg *sync.WaitGroup, path string) {
	defer body.Close() // nolint: errcheck

	doc, err := goQuery.NewDocumentFromReader(body)
	if err != nil {
		return
	}

	for element, attributes := range htmlElements {
		for _, attr := range attributes {
			cr.lookup(doc, element, attr, path, wg)
		}
	}
}

func (cr *crawlerImp) lookup(doc *goQuery.Document, element, attr, path string, wg *sync.WaitGroup) {
	doc.Find(element).Each(func(i int, sel *goQuery.Selection) {
		name, exist := sel.Attr(attr)
		if !exist || !validateLink(name) {
			return
		}
		path = strings.TrimSuffix(path, "/")
		cr.siteCh <- linkWrapper{
			rootURL: path,
			linkURL: name,
		}
		if cr.allowedToProcess(name) {
			wg.Add(1)
			go cr.walkLinks(name, wg)
		}
	})
}

// allowedToProcess is like a lock, to prevent duplicate visiting the same link
func (cr *crawlerImp) allowedToProcess(path string) /*allowed*/ bool {
	link, err := url.Parse(path)
	if err != nil {
		return false
	}
	if strings.Index(link.Path, ".") > 0 {
		// means that path contains link to static file
		return false
	}
	path = link.Host + link.Path

	if !strings.HasPrefix(strings.TrimPrefix(path, wwwKey), cr.domain) {
		// host is not the same as domain
		return false
	}
	return cr.setupLink(path)
}

func (cr *crawlerImp) setupLink(path string) bool {
	cr.Lock()
	defer cr.Unlock()
	visited := cr.visited[path]
	if visited {
		// site already visited
		return false
	}
	cr.visited[path] = true
	if !strings.HasPrefix(path, wwwKey) {
		// to prevent visit link "google.com" in case "www.google.com" already visited
		return !cr.visited[wwwKey+path]
	}

	return true
}
