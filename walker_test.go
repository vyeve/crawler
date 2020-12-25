package crawler

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/vyeve/crawler/mocks"
)

func TestCrawlerWalkGetBadStatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	link := "https://foo.bar"
	wg := sync.WaitGroup{}
	wg.Add(1)
	cr := &crawlerImp{
		initLink:  link,
		siteCh:    make(chan linkWrapper, 10),
		visited:   make(map[string]bool),
		links:     make(siteLinks),
		semaphore: make(chan struct{}, 1),
		scheme:    "https://",
	}
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, link,
		httpmock.NewStringResponder(http.StatusBadRequest, ""))
	cr.walkLinks(link, &wg)
	go func() {
		wg.Wait()
		close(cr.siteCh)
		close(cr.semaphore)
	}()
	_, opened := <-cr.siteCh
	if opened {
		t.Error("unexpected result")
	}
}

func TestCrawler_parseHTML(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockReadCloser(ctrl)
	testErr := errors.New("test")
	reader.EXPECT().Read(gomock.Any()).Return(0, testErr)
	reader.EXPECT().Close().Return(nil).AnyTimes()
	var wg sync.WaitGroup
	cr := crawlerImp{}
	cr.parseHTML(reader, &wg, "")
}

func TestCrawler_allowedToProcess(t *testing.T) {
	testCases := []struct {
		name   string
		domain string
		path   string
		exp    bool
	}{
		{
			name:   "test with link to file",
			domain: "test.com",
			path:   "/foo/bar/index.html",
			exp:    false,
		},
		{
			name:   "test with another domain",
			domain: "www.test.test",
			path:   "/foo/bar/",
			exp:    false,
		},
		{
			name:   "test with already visited site",
			domain: "www.test.com",
			path:   "/foo/",
			exp:    false,
		},
		{
			name:   "test with already visited www site",
			domain: "test.com",
			path:   "/foo/",
			exp:    false,
		},
		{
			name:   "test with OK",
			domain: "www.test.com",
			path:   "/golang/",
			exp:    true,
		},
	}
	cr := &crawlerImp{
		visited: make(map[string]bool),
		domain:  "test.com",
		scheme:  "https://",
	}
	cr.visited["www.test.com/foo/"] = true
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := cr.allowedToProcess(tc.domain, tc.path)
			if tc.exp != res {
				t.Errorf("wrong result for %s. Expected: %t, actual: %t", tc.path, tc.exp, res)
			}
		})
	}
}
