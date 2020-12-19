package crawler

import "testing"

func TestExtractDomain(t *testing.T) {
	testCases := []struct {
		name      string
		link      string
		expResult string
	}{
		{
			name:      "test http",
			link:      "http://golang.org",
			expResult: "golang.org",
		},
		{
			name:      "test https",
			link:      "https://golang.org",
			expResult: "golang.org",
		},
		{
			name:      "test http with www",
			link:      "https://www.golang.org",
			expResult: "golang.org",
		},
		{
			name:      "test https with www",
			link:      "https://www.golang.org",
			expResult: "golang.org",
		},
		{
			name:      "test wrong url format",
			link:      "https:/www.golang.org",
			expResult: "https:/www.golang.org",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := extractDomain(tc.link)
			if res != tc.expResult {
				t.Errorf("Expected: %s, but got: %s", tc.expResult, res)
			}
		})
	}
}

func TestValidateLink(t *testing.T) {
	testCases := []struct {
		name      string
		link      string
		expResult bool
	}{
		{
			name:      "test with http",
			link:      "http://golang.org",
			expResult: true,
		},
		{
			name:      "test with https",
			link:      "https://golang.org",
			expResult: true,
		},
		{
			name:      "test with wrong link",
			link:      "/golang.org",
			expResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := validateLink(tc.link)
			if res != tc.expResult {
				t.Errorf("wrong result for %s: expected: %t, actual: %t", tc.link, tc.expResult, res)
			}
		})
	}
}
