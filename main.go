package main

import (
	"crawler/crawlers"
	"crawler/logger"
	utilsflag "crawler/utils/flag"
	utilsio "crawler/utils/io"
	utilsurl "crawler/utils/url"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
)

type Opts struct {
	Uri        *url.URL
	Depth      int
	Verbose    bool
	Expression *regexp.Regexp
}

func main() {
	opts := parseFlags()
	errorLogger := logger.NewErrorLogger()

	InterruptCh := interruptHandler()

	siteCrawlerOpts := crawlers.SiteCrawlerOpts{
		Depth:       opts.Depth,
		InterruptCh: InterruptCh,
	}

	if opts.Verbose {
		siteCrawlerOpts.DebugLogger = logger.NewDebugLogger()
		siteCrawlerOpts.ErrorLogger = errorLogger
	}

	res, err := crawlers.SiteCrawler(opts.Uri, siteCrawlerOpts)

	if err != nil {
		errorLogger.Fatalln(err)
	}

	fmt.Println(string(res))
	fmt.Fprintf(os.Stderr, "%s written\n", utilsio.GetHumanReadableSize(len(res)))
}

func interruptHandler() chan os.Signal {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	return signalCh
}

// Command line arguments Parsing
func parseFlags() *Opts {
	opts := &Opts{}

	utilsflag.Parse(
		"[OPTIONS] [URL]",
		[][2]string{
			utilsflag.NewVar(flag.BoolVar, &opts.Verbose, "verbose", "v", false, "Verbose output"),
			utilsflag.NewFunc("depth", "d", "Max depth to crawl to", depthFlagBuilder(opts)),
			utilsflag.NewFunc("expression", "e", "Regular expression urls must match", expressionFlagBuilder(opts)),
		},
		[]utilsflag.LookupFn{
			urlFlagBuilder(opts),
		},
	)

	return opts
}

func expressionFlagBuilder(opts *Opts) func(string) error {
	return func(flagValue string) error {
		re, err := regexp.Compile(flagValue)

		if err != nil {
			return err
		}

		opts.Expression = re
		return nil
	}
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
