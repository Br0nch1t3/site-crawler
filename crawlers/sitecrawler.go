package crawlers

import (
	"crawler/logger"
	"crawler/models"
	utilsurl "crawler/utils/url"
	utilsxml "crawler/utils/xml"
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"sync"
)

type SiteCrawlerOpts struct {
	Depth       int
	Verbose     bool
	InterruptCh chan os.Signal
	debugLogger *log.Logger
	errorLogger *log.Logger
}

func SiteCrawler(uri *url.URL, opts ...SiteCrawlerOpts) ([]byte, error) {
	if len(opts) == 0 {
		opts = []SiteCrawlerOpts{{Depth: -1}}
	} else if opts[0].Verbose {
		opts[0].debugLogger = logger.NewDebugLogger()
		opts[0].errorLogger = logger.NewErrorLogger()
	}

	links := &models.Array[models.Link]{models.Link{Href: uri}}
	mu := sync.Mutex{}

	crawler := crawlerBuilder(opts[0])

	if opts[0].InterruptCh != nil {
		go func() {
			<-opts[0].InterruptCh
			signal.Reset(os.Interrupt)
			close(opts[0].InterruptCh)
			fmt.Fprintln(os.Stderr, "Exiting...")
		}()
	}

	if err := crawler(links, (*links)[0], &mu); err != nil {
		if opts[0].errorLogger != nil {
			opts[0].errorLogger.Println(err)
		}
		return nil, fmt.Errorf(`unable to crawl "%s"`, uri.String())
	}

	res, err := xml.MarshalIndent(links, "", " ")

	if err != nil {
		if opts[0].errorLogger != nil {
			opts[0].errorLogger.Println(err)
		}
		return nil, err
	}

	return utilsxml.WithHeader(res), nil
}

type CrawlerFn func(*models.Array[models.Link], models.Link, *sync.Mutex) error

func crawlerBuilder(opts SiteCrawlerOpts) CrawlerFn {
	return func(visited *models.Array[models.Link], baseLink models.Link, mu *sync.Mutex) error {
		if opts.debugLogger != nil {
			opts.debugLogger.Printf("crawling %s\n", baseLink.Href)
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
	return (opts.Depth == -1 || utilsurl.PathLen(link.Href) <= opts.Depth) && !slices.ContainsFunc(visited, link.SameHref)
}
