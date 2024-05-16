package main

import (
	"crawler/crawlers"
	utilsio "crawler/utils/io"
	utilsurl "crawler/utils/url"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type Opts struct {
	Uri    *url.URL
	Depth  *int
	Output string
}

func main() {
	signalCh := make(chan os.Signal, 1)
	opts := parseFlags()

	signal.Notify(signalCh, os.Interrupt)

	res, err := crawlers.SiteCrawler(opts.Uri, crawlers.SiteCrawlerOpts{
		Depth:       *opts.Depth,
		InterruptCh: signalCh,
	})

	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(opts.Output+".xml", res, 0644); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s written to %s.xml\n", utilsio.GetHumanReadableSize(len(res)), opts.Output)
}

func parseFlags() *Opts {
	opts := &Opts{
		Depth: flag.Int("depth", -1, "Max depth to crawl to"),
	}

	flag.Func("url", "Url to crawl", urlFlagBuilder(opts))
	flag.Func("output", "Name of the output file without extension (default: \"sitemap\")", outputFlagBuilder(opts))

	flag.Parse()

	if opts.Uri == nil {
		flag.Usage()
		os.Exit(1)
	}

	return opts
}

func urlFlagBuilder(opts *Opts) func(string) error {
	return func(flagValue string) error {
		uri, err := url.ParseRequestURI(flagValue)

		if err != nil || !utilsurl.IsHttp(uri) {
			return err
		}

		opts.Uri = uri
		return nil
	}
}

func outputFlagBuilder(opts *Opts) func(string) error {
	opts.Output = "sitemap"

	return func(flagValue string) error {
		if flagValue[len(flagValue)-1] == '.' || len(flagValue) == 0 {
			return errors.New("invalid filename")
		}
		opts.Output = flagValue

		return nil
	}
}
