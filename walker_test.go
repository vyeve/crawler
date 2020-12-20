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
		initLink: link,
		siteCh:   make(chan linkWrapper, 10),
		visited:  make(map[string]bool),
		links:    make(siteLinks),
	}
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, link,
		httpmock.NewStringResponder(http.StatusBadRequest, ""))
	cr.walkLinks(link, &wg)
	go func() {
		wg.Wait()
		close(cr.siteCh)
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
		name string
		path string
		exp  bool
	}{
		{
			name: "test with parse URL error",
			path: string([]byte{0x7f}),
			exp:  false,
		},
		{
			name: "test with link to file",
			path: "https://test.com/foo/bar/index.html",
			exp:  false,
		},
		{
			name: "test with another domain",
			path: "https://www.test.test/foo/bar/",
			exp:  false,
		},
		{
			name: "test with already visited site",
			path: "https://www.test.com/foo/",
			exp:  false,
		},
		{
			name: "test with already visited www site",
			path: "https://test.com/foo/",
			exp:  false,
		},
		{
			name: "test with OK",
			path: "https://www.test.com/golang/",
			exp:  true,
		},
	}
	cr := &crawlerImp{
		visited: make(map[string]bool),
		domain:  "test.com",
	}
	cr.visited["www.test.com/foo/"] = true
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := cr.allowedToProcess(tc.path)
			if tc.exp != res {
				t.Errorf("wrong result for %s. Expected: %t, actual: %t", tc.path, tc.exp, res)
			}
		})
	}
}
