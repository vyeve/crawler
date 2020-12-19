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
	s := &crawlerImp{
		initLink: initLink,
		siteCh:   make(chan linkWrapper, 10),
		visited:  make(map[string]bool),
		links:    make(siteLinks),
	}
	// extract domain
	s.domain = extractDomain(initLink)
	// validate if link is alive
	resp, err := http.Get(initLink)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return s, nil
}

func (ci *crawlerImp) Crawl() Printer {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	ci.setupLink(ci.initLink)
	go ci.walkLinks(ci.initLink, wg)
	go func() {
		wg.Wait()
		close(ci.siteCh)
	}()

	for links := range ci.siteCh {
		ci.links.addLinkToSite(links.rootURL, links.linkURL)
	}
	return ci.links
}
