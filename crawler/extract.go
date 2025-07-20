package crawler

import (
	"io"
	"log"
	"net/url"

	"golang.org/x/net/html"
)

func ExtractLinks(base *url.URL, body io.Reader) ([]*url.URL, error) {
	var result []*url.URL
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			tagName, hasAttr := tokenizer.TagName()

			if string(tagName) != "a" || !hasAttr {
				continue
			}

			for {
				key, val, more := tokenizer.TagAttr()

				if string(key) == "href" {
					href, err := base.Parse(string(val))
					if err != nil {
						log.Fatalf("Error parsing href from tag::\nkey:%+v\nvalue:%+v\n", string(key), string(val))
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
