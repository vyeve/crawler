package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/vyeve/crawler"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	siteName := "https://github.com/"

	tn := time.Now()
	defer func() {
		fmt.Printf("\n\tTotal time: %s\n", time.Since(tn))
	}()
	f, err := os.OpenFile("result.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = f.Truncate(0); err != nil {
		panic(err)
	}
	if _, err = f.Seek(0, 0); err != nil {
		panic(err)
	}
	st, err := crawler.New(siteName)
	if err != nil {
		panic(err)
	}
	st.Crawl().Print(f)
}
