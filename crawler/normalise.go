package crawler

import (
	"net/url"
	"strings"
)

func NormaliseURL(u *url.URL) *url.URL {
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.ToLower(u.Path)
	u.Scheme = strings.ToLower(u.Scheme)
	u.Fragment = ""

	if u.Scheme == "http" {
		u.Scheme = "https"
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u
}
