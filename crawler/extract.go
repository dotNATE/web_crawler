package crawler

import (
	"io"
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(base *url.URL, body io.Reader) ([]*url.URL, error) {
	var result []*url.URL
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			tagName, hasAttr := tokenizer.TagName()

			if (string(tagName) != "a" && string(tagName) != "link") || !hasAttr {
				continue
			}

			for {
				key, val, more := tokenizer.TagAttr()

				if string(key) == "href" {
					hrefString := strings.TrimSpace(string(val))
					if hrefString == "" || strings.HasPrefix(hrefString, "mailto:") || strings.HasPrefix(hrefString, "javascript:") {
						break
					}

					href, err := base.Parse(hrefString)
					if err != nil {
						log.Fatalf("Error parsing href from tag::\nkey:%+v\nvalue:%+v\nerror:%+v\n", string(key), string(val), err)
					}

					result = append(result, href)
				}

				if !more {
					break
				}
			}
		}

		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return result, nil
			}

			return nil, tokenizer.Err()
		}
	}
}
