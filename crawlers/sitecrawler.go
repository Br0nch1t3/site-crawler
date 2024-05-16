package crawlers

import (
	"crawler/models"
	utilsxml "crawler/utils/xml"
	"encoding/xml"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
)

type SiteCrawlerOpts struct {
	Depth       int
	InterruptCh chan os.Signal
}

func SiteCrawler(uri *url.URL, opts ...SiteCrawlerOpts) ([]byte, error) {
	if len(opts) == 0 {
		opts = []SiteCrawlerOpts{{Depth: -1}}
	}

	links := &models.Array[models.Link]{models.Link{Href: uri}}
	mu := sync.Mutex{}

	crawler := buildCrawler(opts[0])

	if opts[0].InterruptCh != nil {
		go func() {
			<-opts[0].InterruptCh
			signal.Reset(os.Interrupt)
			close(opts[0].InterruptCh)
			fmt.Fprintln(os.Stderr, "Exiting...")
		}()
	}

	crawler(links, (*links)[0], &mu)

	res, err := xml.MarshalIndent(links, "", " ")

	if err != nil {
		return nil, err
	}

	return utilsxml.WithHeader(res), nil
}

func buildCrawler(opts SiteCrawlerOpts) func(*models.Array[models.Link], models.Link, *sync.Mutex) {
	return func(visited *models.Array[models.Link], baseLink models.Link, mu *sync.Mutex) {
		links, err := PageCrawler(baseLink.Href)

		if err != nil {
			return
		}

		wg := sync.WaitGroup{}
		for _, link := range links {
			select {
			case <-opts.InterruptCh:
				return
			default:
				if slices.ContainsFunc(*visited, link.SameHref) || (opts.Depth >= 0 && strings.Count(strings.Trim(link.Href.Path, "/"), "/") >= opts.Depth) {
					continue
				}

				mu.Lock()
				*visited = append(*visited, link)
				mu.Unlock()

				wg.Add(1)
				go func(_link models.Link) {
					crawl := buildCrawler(opts)
					crawl(visited, _link, mu)
					wg.Done()
				}(link)
			}
		}

		wg.Wait()
	}
}
