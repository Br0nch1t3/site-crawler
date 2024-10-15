package crawlers

import (
	httpclient "crawler/http-client"
	"crawler/models"
	utilshttp "crawler/utils/http"
	utilsurl "crawler/utils/url"
	"errors"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func PageCrawler(uri *url.URL) (models.Array[models.Link], error) {
	visited := []models.Link{}
	res, redirectErr := httpclient.Instance().Get(uri.String())

	if redirectErr != nil {
		return nil, redirectErr
	}

	if res == nil {
		panic(fmt.Sprint("Request to ", uri.String(), " couldn't be resolved."))
	}

	if err := utilshttp.ExtractError(res); err != nil {
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
			if token.Data == "a" && slices.IndexFunc(token.Attr, func(attr html.Attribute) bool { return attr.Key == "href" }) != -1 {
				href, err := parseHref(token.Attr[0].Val, uri, visited)

				if err != nil {
					continue
				}

				isAnchorText = true
				visited = append(visited, models.Link{
					Href: href,
				})
			}
		case html.TextToken:
			if isAnchorText {
				visited[len(visited)-1].Text = strings.TrimSpace(strings.ReplaceAll(token.Data, "\n", " "))
				isAnchorText = false
			}
		default:
			if isAnchorText {
				// Remove any non-text links
				visited = visited[:len(visited)-1]
				isAnchorText = false
			}
		}
	}

	return visited, nil
}

func isSupportedUrl(uri *url.URL, baseUri *url.URL, visited []models.Link) bool {
	return utilsurl.SameOrigin(uri, baseUri) && !slices.ContainsFunc(visited, func(l models.Link) bool {
		return l.Href.String() == uri.String()
	})
}

func parseHref(rawHref string, uri *url.URL, visited []models.Link) (*url.URL, error) {
	if i := strings.LastIndex(rawHref, "#"); i != -1 {
		rawHref = rawHref[:i]
	}

	href, err := url.Parse(rawHref)

	if err != nil {
		return nil, err
	}

	if len(href.Path) == 0 {
		return nil, errors.New("empty href")
	}

	if !href.IsAbs() {
		if href.String()[0] == '/' {
			href, err = uri.Parse(href.EscapedPath())

			if err != nil {
				log.Println(uri)
				log.Fatalln(err)
			}
		} else {
			href = uri.JoinPath(href.EscapedPath())
		}
	}

	if !isSupportedUrl(href, uri, visited) {
		return nil, errors.New("url not supported")
	}

	return href, nil
}
