package crawler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	goQuery "github.com/PuerkitoBio/goquery"
)

// walkLinks requests given path, validates response and parse HTML
func (cr *crawlerImp) walkLinks(path string, wg *sync.WaitGroup) {
	defer wg.Done()
	cr.semaphore <- struct{}{}
	defer func() {
		<-cr.semaphore
	}()
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

// parseHTML parses given response body to lookup HTML elements, which contain external URL
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

// lookup tries to find HTML's elements which contain URL,  according to the W3C's list of
// HTML attributes. In case link is found, it launch walkLinks in separate goroutine
func (cr *crawlerImp) lookup(doc *goQuery.Document, element, attr, path string, wg *sync.WaitGroup) {
	doc.Find(element).Each(func(i int, sel *goQuery.Selection) {
		name, exist := sel.Attr(attr)
		if !exist {
			return
		}
		link, err := url.Parse(name)
		if err != nil {
			return
		}
		var domain string
		if len(link.Host) == 0 {
			domain = cr.domain
			name = cr.scheme + cr.domain + link.Path
		} else {
			domain = link.Host
			var scheme string
			if len(link.Scheme) == 0 {
				scheme = cr.scheme
			} else {
				scheme = link.Scheme + "://"
			}
			name = scheme + link.Host + link.Path
		}

		name = strings.TrimSuffix(name, "/")
		cr.siteCh <- linkWrapper{
			rootURL: path,
			linkURL: name,
		}
		if cr.allowedToProcess(domain, link.Path) {
			wg.Add(1)
			go cr.walkLinks(name, wg)
		}
	})
}

// allowedToProcess is like a lock, to prevent duplicate visiting the same link
func (cr *crawlerImp) allowedToProcess(domain, path string) /*allowed*/ bool {
	if strings.Index(path, ".") > 0 {
		// means that path contains link to static file
		return false
	}

	if !strings.HasPrefix(strings.TrimPrefix(domain, wwwKey), cr.domain) {
		// host is not the same as domain
		return false
	}
	return cr.setupLink(domain + path)
}

// setupLink sets given link to prevent duplication
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
