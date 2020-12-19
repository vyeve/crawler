package crawler

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNew(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	path := "https://golang.org"
	testCases := []struct {
		name    string
		path    string
		code    int
		needErr bool
	}{
		{
			name:    "test OK",
			path:    path,
			code:    http.StatusOK,
			needErr: false,
		},
		{
			name:    "test OK",
			code:    http.StatusOK,
			needErr: false,
		},
		{
			name:    "test with bad request",
			path:    path,
			code:    http.StatusNotFound,
			needErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, tc.path,
				httpmock.NewStringResponder(tc.code, ""))
			_, err := New(path)
			if tc.needErr {
				if err == nil {
					t.Error("expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}

		})
	}
}

func TestNewWithGetError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Error("expected not <nil> error")
	}
}

func TestCrawlerCrawl(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	link := "https://foo.bar"
	cr := &crawlerImp{
		initLink: link,
		siteCh:   make(chan linkWrapper, 10),
		visited:  make(map[string]bool),
		links:    make(siteLinks),
	}

	doc1 := `
	<a class="Button tour" href="https://test.foo.bar/"
	title="Playground Go from your browser">Tour</a>
	`
	httpmock.Reset()
	httpmock.RegisterResponder(http.MethodGet, link,
		httpmock.NewStringResponder(http.StatusOK, doc1))
	cr.Crawl()
	visited, ok := cr.links["https://foo.bar"]
	if !ok {
		t.Fatalf("unexpected result: %v", cr.links)
	}
	_, ok = visited["https://test.foo.bar/"]
	if !ok {
		t.Errorf("unexpected result: %v", visited)
	}

}
