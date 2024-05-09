package main

import (
	"crawler/crawler"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Provide a domain like so: ./crawler {domain}")
		os.Exit(1)
	}

	res, err := crawler.SitemapGenerator(args[0])

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(res))
}
