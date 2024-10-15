package crawlers

import (
	"crawler/models"
	utilsurl "crawler/utils/url"
	utilsxml "crawler/utils/xml"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"sync"
)

type SiteCrawlerOpts struct {
	Depth       int
	Expression  *regexp.Regexp
	InterruptCh chan os.Signal
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
}

func SiteCrawler(uri *url.URL, opts ...SiteCrawlerOpts) ([]byte, error) {
	parsedOpts := parseOpts(opts...)

	links := &models.Array[models.Link]{models.Link{Href: uri}}
	mu := sync.Mutex{}

	crawler := crawlerBuilder(parsedOpts)

	if parsedOpts.InterruptCh != nil {
		go func() {
			<-parsedOpts.InterruptCh
			signal.Reset(os.Interrupt)
			close(parsedOpts.InterruptCh)
			fmt.Fprintln(os.Stderr, "Exiting...")
		}()
	}

	if err := crawler(links, (*links)[0], &mu); err != nil {
		if parsedOpts.ErrorLogger != nil {
			parsedOpts.ErrorLogger.Println(err)
		}
		return nil, fmt.Errorf(`unable to crawl "%s"`, uri.String())
	}

	res, err := xml.MarshalIndent(links, "", " ")

	if err != nil {
		if parsedOpts.ErrorLogger != nil {
			parsedOpts.ErrorLogger.Println(err)
		}
		return nil, err
	}

	return utilsxml.WithHeader(res), nil
}

type CrawlerFn func(*models.Array[models.Link], models.Link, *sync.Mutex) error

func crawlerBuilder(opts SiteCrawlerOpts) CrawlerFn {
	return func(visited *models.Array[models.Link], baseLink models.Link, mu *sync.Mutex) error {
		if opts.DebugLogger != nil {
			opts.DebugLogger.Printf("crawling %s\n", baseLink.Href)
		}
		links, err := PageCrawler(baseLink.Href)

		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		for _, link := range links {
			select {
			case <-opts.InterruptCh:
				return nil
			default:
				if !isCrawlable(*visited, link, opts) {
					continue
				}

				mu.Lock()
				*visited = append(*visited, link)
				mu.Unlock()

				wg.Add(1)
				go func(_link models.Link) {
					crawler := crawlerBuilder(opts)
					crawler(visited, _link, mu)
					wg.Done()
				}(link)
			}
		}

		wg.Wait()
		return nil
	}
}

func isCrawlable(visited models.Array[models.Link], link models.Link, opts SiteCrawlerOpts) bool {
	matches := opts.Expression == nil || opts.Expression.Match([]byte(link.Href.String()))

	return matches && (opts.Depth == -1 || utilsurl.PathLen(link.Href) <= opts.Depth) && !slices.ContainsFunc(visited, link.SameHref)
}

func parseOpts(opts ...SiteCrawlerOpts) SiteCrawlerOpts {
	if len(opts) == 0 {
		return SiteCrawlerOpts{Depth: -1}
	}
	return opts[0]
}
