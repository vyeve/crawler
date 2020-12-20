# crawler
A simple web crawler, implemented in Go

## Description
Concurrently implementation of a simple web crawler using Go. Given a starting URL, the crawler visits each URL it found on the same domain. It should print each URL visited, and a list of links found on that page. The crawler is limited to one subdomain. So when it starts with `https://www.github.com/`, it crawls all pages within `github.com`, but not follows external links, for example to `youtube.com` or `docs.github.com`.

## Example of usage:

```go
package main

import (
	"os"
	"runtime"

	"github.com/vyeve/crawler"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	siteName := "https://github.com/"

	st, err := crawler.New(siteName)
	if err != nil {
		panic(err)
	}
	st.Crawl().Print(os.Stdout)
}

```

### More examples can be found by [link](https://github.com/vyeve/crawler/tree/master/example "Crawler usage")

## Test coverage
To see test coverage run:
```bash
go test -v -tags=unit -coverprofile=cover.out && go tool cover -func=cover.out >> test.out && go tool cover -html=cover.out
```