package crawler

import (
	"net/url"
	"strings"
)

// extractDomain retrieves domain from the given path
func extractDomain(link string) ( /*scheme*/ string /*domain*/, string, error) {
	path, err := url.Parse(link)
	if err != nil {
		return "", "", err
	}
	link = strings.TrimPrefix(link, httpKey)
	link = strings.TrimPrefix(link, httpsKey)
	link = strings.TrimPrefix(link, wwwKey)
	link = strings.TrimSuffix(link, "/")
	return path.Scheme + "://", link, nil
}

// validateLink checks if given link contains http/https prefix
func validateLink(link string) bool {
	if strings.HasPrefix(link, httpsKey) ||
		strings.HasPrefix(link, httpKey) ||
		strings.HasPrefix(link, "/") {
		return true // write only links which contain http/https or starts with /
	}
	return false
}
