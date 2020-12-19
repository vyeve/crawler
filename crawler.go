package crawler

import (
	"fmt"
	"net/http"
	"sync"
)

type crawlerImp struct {
	sync.RWMutex
	initLink string
	domain   string
	visited  map[string]bool
	links    siteLinks
	siteCh   chan linkWrapper
}

type linkWrapper struct {
	rootURL string
	linkURL string
}

func New(initLink string) (Crawler, error) {
	cr := &crawlerImp{
		initLink: initLink,
		siteCh:   make(chan linkWrapper, 10),
		visited:  make(map[string]bool),
		links:    make(siteLinks),
	}
	// extract domain
	cr.domain = extractDomain(initLink)
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

func (cr *crawlerImp) Crawl() Printer {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	cr.setupLink(cr.initLink)
	go cr.walkLinks(cr.initLink, wg)
	go func() {
		wg.Wait()
		close(cr.siteCh)
	}()

	for links := range cr.siteCh {
		cr.links.addLinkToSite(links.rootURL, links.linkURL)
	}
	return cr.links
}
