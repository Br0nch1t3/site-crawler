package main

import (
	"errors"
	"fmt"
	"io"
	"link-parser/models"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

var prevDuration time.Duration

func main() {
	links, err := WithTimer(func() ([]models.Link, error) { return LinkParser("https://hugo-mialhe.vercel.app/fr") }, "LOGIC")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(links)
}

func LinkParser(url string) (models.Array[models.Link], error) {
	links := []models.Link{}
	res, err := WithTimer(func() (*http.Response, error) { return http.Get(url) }, "GET")

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
				links[len(links)-1].Text = token.Data
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

func WithTimer[T any](fn func() (T, error), prefix ...string) (T, error) {
	now := time.Now()
	res, err := fn()
	duration := time.Since(now) - prevDuration
	prevDuration = duration
	text := duration.String()
	if len(prefix) > 0 {
		text = prefix[0] + ": " + text
	}
	fmt.Println(text)
	return res, err
}
