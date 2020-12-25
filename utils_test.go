package crawler

import "testing"

func TestExtractDomain(t *testing.T) {
	testCases := []struct {
		name      string
		link      string
		expResult string
		expScheme string
		needErr   bool
	}{
		{
			name:      "test http",
			link:      "http://golang.org",
			expResult: "golang.org",
			expScheme: "http",
			needErr:   false,
		},
		{
			name:      "test https",
			link:      "https://golang.org",
			expResult: "golang.org",
			expScheme: "https",
			needErr:   false,
		},
		{
			name:      "test http with www",
			link:      "https://www.golang.org",
			expResult: "golang.org",
			expScheme: "https",
			needErr:   false,
		},
		{
			name:      "test https with www",
			link:      "https://www.golang.org",
			expResult: "golang.org",
			expScheme: "https",
			needErr:   false,
		},
		{
			name:      "test wrong url format",
			link:      "https:/www.golang.org",
			expResult: "https:/www.golang.org",
			needErr:   false,
		},
		{
			name:    "test wrong url format",
			link:    string([]byte{0x7f}),
			needErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme, res, err := extractDomain(tc.link)
			if tc.needErr {
				if err == nil {
					t.Error("expected not <nil> error")
				}
			} else {
				if err != nil {
					t.Error(err)
				}
				if res != tc.expResult && scheme != tc.expScheme {
					t.Errorf("Expected: %s, but got: %s", tc.expResult, res)
				}
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
			link:      "golang.org",
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
