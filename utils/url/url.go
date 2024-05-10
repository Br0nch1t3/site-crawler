package utilsurl

import (
	"net/url"
	"regexp"
)

// Returns true if url a and url b are the same origin
func SameOrigin(a *url.URL, b *url.URL) bool {
	return a.Scheme+a.Host == b.Scheme+b.Host
}

// Returns true if url scheme is http or https
func IsHttp(uri *url.URL) bool {
	reg := regexp.MustCompile(`^https?`)

	return reg.Match([]byte(uri.Scheme))
}
