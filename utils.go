package crawler

import (
	"strings"
)

func extractDomain(link string) string {
	link = strings.TrimPrefix(link, httpKey)
	link = strings.TrimPrefix(link, httpsKey)
	link = strings.TrimPrefix(link, wwwKey)
	link = strings.TrimSuffix(link, "/")
	return link
}

func validateLink(link string) bool {
	if strings.HasPrefix(link, httpsKey) || strings.HasPrefix(link, httpKey) {
		return true // write only links which contain http/https
	}
	return false
}
