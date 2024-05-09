package models

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

type Link struct {
	XMLName xml.Name `xml:"url"`
	Href    string   `xml:"loc"`
	Text    string   `xml:"-"`
}

func (l Link) String() string {
	return fmt.Sprintf("{\"href\": \"%s\", \"text\": \"%s\"}", l.Href, l.Text)
}

func (l Link) SamePath(link Link) bool {
	parsedBaseLink, err := url.Parse(l.Href)

	if err != nil {
		return false
	}

	parsedLink, err := url.Parse(link.Href)

	if err != nil {
		return false
	}

	return strings.Trim(parsedBaseLink.Path, "/") == strings.Trim(parsedLink.Path, "/")
}

type xmlArrayWrapper[T fmt.Stringer] struct {
	XMLName xml.Name `xml:"urlset"`
	Content []T
}

func (arr Array[Link]) String() string {
	res := make([]string, len(arr))

	for i, link := range arr {
		res[i] = link.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(res, ","))
}

func (r Array[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	wrapper := xmlArrayWrapper[T]{Content: []T(r)}

	return e.Encode(wrapper)
}
