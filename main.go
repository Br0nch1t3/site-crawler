package main

import (
	"crawler/crawler"
	utilsurl "crawler/utils/url"
	"fmt"
	"log"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Provide a url like so: ./crawler {url}")
		os.Exit(1)
	}

	uri, err := url.ParseRequestURI(args[0])

	if err != nil || !utilsurl.IsHttp(uri) {
		fmt.Fprintln(os.Stderr, "Please provide a valid url")
		os.Exit(1)
	}

	res, err := crawler.SitemapGenerator(uri)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(res))
}
