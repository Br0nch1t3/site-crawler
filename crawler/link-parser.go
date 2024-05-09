package crawler

import (
	"crawler/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func LinkParser(url string) (models.Array[models.Link], error) {
	links := []models.Link{}
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	tokenizer := html.NewTokenizer(res.Body)

	if tokenizer == nil {
		return nil, errors.New("error tokenizer")
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
				isAnchorText = true
				links = append(links, models.Link{
					Href: token.Attr[0].Val,
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

func PrintBody(r io.ReadCloser) {
	body, err := io.ReadAll(r)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
