package main

import (
	"crawler/crawlers"
	"crawler/logger"
	utilsio "crawler/utils/io"
	utilsurl "crawler/utils/url"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
)

type Opts struct {
	Uri     *url.URL
	Depth   int
	Output  string
	Verbose *bool
}

func main() {
	signalCh := make(chan os.Signal, 1)
	opts := parseFlags()
	errorLogger := logger.NewErrorLogger()

	signal.Notify(signalCh, os.Interrupt)

	res, err := crawlers.SiteCrawler(opts.Uri, crawlers.SiteCrawlerOpts{
		Depth:       opts.Depth,
		Verbose:     *opts.Verbose,
		InterruptCh: signalCh,
	})

	if err != nil {
		errorLogger.Fatalln(err)
	}

	if err := os.WriteFile(opts.Output+".xml", res, 0644); err != nil {
		errorLogger.Fatalln(err)
	}
	fmt.Printf("%s written to %s.xml\n", utilsio.GetHumanReadableSize(len(res)), opts.Output)
}

func parseFlags() *Opts {
	opts := &Opts{
		Verbose: flag.Bool("verbose", false, "Verbose output (default: false)"),
	}

	flag.Func("url", "Url to crawl", urlFlagBuilder(opts))
	flag.Func("output", "Name of the output file without extension (default: \"sitemap\")", outputFlagBuilder(opts))
	flag.Func("depth", "Max depth to crawl to", depthFlagBuilder(opts))

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
			return errors.New("url must be http(s)://[domain]/[path]")
		}

		opts.Uri = uri
		return nil
	}
}

func outputFlagBuilder(opts *Opts) func(string) error {
	opts.Output = "sitemap"

	return func(flagValue string) error {
		if len(flagValue) == 0 || flagValue[len(flagValue)-1] == '.' {
			return errors.New("url must be http(s)://[domain]/[path]")
		}
		opts.Output = flagValue

		return nil
	}
}

func depthFlagBuilder(opts *Opts) func(string) error {
	opts.Depth = -1

	return func(flagValue string) error {
		depth, err := strconv.Atoi(flagValue)

		if err != nil {
			return errors.New("value should be numerical")
		}

		if depth < 0 {
			opts.Depth = -1
		} else {
			opts.Depth = depth
		}

		return nil
	}
}
