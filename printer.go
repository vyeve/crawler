package crawler

import (
	"fmt"
	"io"
	"sort"
)

type siteLinks map[string]map[string]struct{}

func (sl siteLinks) addLinkToSite(site, link string) {
	if _, ok := sl[site]; !ok {
		sl[site] = make(map[string]struct{})
	}
	sl[site][link] = struct{}{}
}

func (sl siteLinks) Print(wr io.Writer) error {
	roots := make([]string, 0, len(sl))
	for r := range sl {
		roots = append(roots, r)
	}
	sort.Strings(roots)
	var err error
	for _, root := range roots {
		if _, err = fmt.Fprintf(wr, "%s\n", root); err != nil {
			return err
		}
		for _, link := range sortLinks(sl[root]) {
			if _, err = fmt.Fprintf(wr, "\t%s\n", link); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortLinks(links map[string]struct{}) []string {
	out := make([]string, 0, len(links))
	for link := range links {
		out = append(out, link)
	}
	sort.Strings(out)
	return out
}
