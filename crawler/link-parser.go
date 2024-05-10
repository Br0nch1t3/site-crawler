package crawler

import (
	"crawler/models"
	utilsurl "crawler/utils/url"
	"errors"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func LinkParser(uri *url.URL) (models.Array[models.Link], error) {
	links := []models.Link{}
	res, err := http.Get(uri.String())

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	tokenizer := html.NewTokenizer(res.Body)

	if tokenizer == nil {
		return nil, errors.New("tokenizer error")
	}

	isAnchorText := false

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		switch tokenType {
		case html.StartTagToken:
			if token.Data == "a" && len(token.Attr) > 0 && token.Attr[0].Key == "href" {
				href, err := url.Parse(token.Attr[0].Val)

				if err != nil || len(href.Path) == 0 {
					continue
				}

				if !href.IsAbs() {
					if href.String()[0] == '/' {
						href, err = uri.Parse(href.Path)

						if err != nil {
							log.Fatalln(err)
						}
					} else {
						href = uri.JoinPath(href.Path)
					}
				}

				if !isSupportedLink(href, uri, links) {
					continue
				}

				isAnchorText = true
				links = append(links, models.Link{
					Href: href,
				})
			}
		case html.TextToken:
			if isAnchorText {
				links[len(links)-1].Text = strings.TrimSpace(strings.ReplaceAll(token.Data, "\n", " "))
				isAnchorText = false
			}
		default:
			if isAnchorText {
				// Remove any non-text links
				links = links[:len(links)-1]
				isAnchorText = false
			}
		}
	}

	return links, nil
}

func isSupportedLink(uri *url.URL, baseUri *url.URL, visited []models.Link) bool {
	return utilsurl.SameOrigin(uri, baseUri) && !slices.ContainsFunc(visited, func(l models.Link) bool {
		return l.Href.String() == uri.String()
	})
}
