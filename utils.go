package crawler

import (
	"strings"
)

// extractDomain retrieves domain from the given path
func extractDomain(link string) string {
	link = strings.TrimPrefix(link, httpKey)
	link = strings.TrimPrefix(link, httpsKey)
	link = strings.TrimPrefix(link, wwwKey)
	link = strings.TrimSuffix(link, "/")
	return link
}

// validateLink checks if given link contains http/https prefix
func validateLink(link string) bool {
	if strings.HasPrefix(link, httpsKey) || strings.HasPrefix(link, httpKey) {
		return true // write only links which contain http/https
	}
	return false
}
