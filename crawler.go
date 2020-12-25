package crawler

import (
	"fmt"
	"net/http"
	"sync"
)

// crawlerImp is a struct which implements Crawler interface
type crawlerImp struct {
	sync.RWMutex
	initLink  string
	scheme    string
	domain    string
	visited   map[string]bool
	links     siteLinks
	siteCh    chan linkWrapper
	semaphore chan struct{}
}

// linkWrapper is a container to communicate scanned links through channel
type linkWrapper struct {
	rootURL string
	linkURL string
}

// New initialize Crawler
func New(initLink string) (Crawler, error) {
	cr := &crawlerImp{
		initLink:  initLink,
		siteCh:    make(chan linkWrapper, 10),
		visited:   make(map[string]bool),
		links:     make(siteLinks),
		semaphore: make(chan struct{}, maxRoutines),
	}
	var err error
	// extract domain
	cr.scheme, cr.domain, err = extractDomain(initLink)
	if err != nil {
		return nil, err
	}
	// validate if link is alive
	resp, err := http.Get(initLink)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return cr, nil
}

// Crawl scans given link
func (cr *crawlerImp) Crawl() Printer {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	cr.setupLink(cr.initLink)
	go cr.walkLinks(cr.initLink, wg)
	go func() {
		wg.Wait()
		close(cr.siteCh)
		close(cr.semaphore)
	}()

	for links := range cr.siteCh {
		cr.links.addLinkToSite(links.rootURL, links.linkURL)
	}
	return cr.links
}
