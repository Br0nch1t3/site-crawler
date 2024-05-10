package models

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

type Link struct {
	XMLName xml.Name
	Href    *url.URL
	Text    string
}

type linkXmlAdapter struct {
	XMLName xml.Name `xml:"url"`
	Href    string   `xml:"loc"`
}

func (l Link) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	adapter := linkXmlAdapter{
		Href: l.Href.String(),
	}

	return e.Encode(adapter)
}

func (l Link) String() string {
	return fmt.Sprintf("{\"href\": \"%s\", \"text\": \"%s\"}", l.Href.String(), l.Text)
}

func (l Link) SameHref(link Link) bool {
	return strings.Trim(l.Href.String(), "/") == strings.Trim(link.Href.String(), "/")
}

type linkArrayXmlAdapter[T fmt.Stringer] struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	Content []T
}

func (r Array[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	adapter := linkArrayXmlAdapter[T]{Content: []T(r)}

	return e.Encode(adapter)
}

func (arr Array[Link]) String() string {
	res := make([]string, len(arr))

	for i, link := range arr {
		res[i] = link.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(res, ","))
}
