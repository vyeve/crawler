package crawler

import (
	"io"
)

// Crawler scans link and returns Printer
type Crawler interface {
	Crawl() Printer
}

// Printer represents scanned links to io.Writer
type Printer interface {
	Print(io.Writer) error
}
