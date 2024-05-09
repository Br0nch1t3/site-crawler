package crawler

import (
	"crawler/models"
	"encoding/xml"
	"net/url"
	"slices"
	"sync"
)

func SitemapGenerator(href string) ([]byte, error) {
	links := &models.Array[models.Link]{}
	mu := sync.Mutex{}

	crawl(links, models.Link{Href: href}, &mu)

	res, err := xml.Marshal(links)

	if err != nil {
		return res, err
	}

	return res, nil
}

func crawl(visited *models.Array[models.Link], baseLink models.Link, mu *sync.Mutex) {
	links, err := LinkParser(baseLink.Href)

	if err != nil {
		return
	}

	wg := sync.WaitGroup{}
	for _, link := range links {
		if slices.ContainsFunc(*visited, link.SamePath) {
			continue
		}

		parsedBaseLink, err := url.Parse(baseLink.Href)
		if err != nil {
			continue
		}

		parsedLink, err := url.Parse(link.Href)
		if err != nil {
			continue
		}

		if parsedLink.Host != parsedBaseLink.Host && (len(parsedLink.Scheme) > 0 || link.Href[0] == '#') {
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
