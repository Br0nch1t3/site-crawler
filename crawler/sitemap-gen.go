package crawler

import (
	"crawler/models"
	utilsxml "crawler/utils/xml"
	"encoding/xml"
	"net/url"
	"slices"
	"sync"
)

func SitemapGenerator(uri *url.URL) ([]byte, error) {
	links := &models.Array[models.Link]{}
	mu := sync.Mutex{}

	crawl(links, models.Link{Href: uri}, &mu)

	res, err := xml.MarshalIndent(links, "", " ")

	if err != nil {
		return nil, err
	}

	return utilsxml.WithHeader(res), nil
}

func crawl(visited *models.Array[models.Link], baseLink models.Link, mu *sync.Mutex) {
	links, err := LinkParser(baseLink.Href)

	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	for _, link := range links {
		if slices.ContainsFunc(*visited, link.SameHref) {
			continue
		}

		mu.Lock()
		*visited = append(*visited, link)
		mu.Unlock()

		wg.Add(1)
		go func(_link models.Link) {
			crawl(visited, _link, mu)
			wg.Done()
		}(link)
	}

	wg.Wait()
}
