package crawler

import (
	"io"
)

type Crawler interface {
	Crawl() Printer
}

type Printer interface {
	Print(io.Writer) error
}
