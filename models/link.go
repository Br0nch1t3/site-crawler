package models

import (
	"fmt"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("{\"href\": \"%s\", \"text\": \"%s\"}", l.Href, l.Text)
}

func (arr Array[Link]) String() string {
	res := make([]string, len(arr))

	for i, link := range arr {
		res[i] = link.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(res, ","))
}
