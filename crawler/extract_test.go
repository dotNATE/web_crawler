package crawler_test

import (
	"net/url"
	"strings"
	"testing"

	"example.com/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractLinks(t *testing.T) {
	html := `<html><body>
		<a href="https://example.com">Internal secure</a>
        <a href="/about">Internal by reference</a>
        <a href="http://external.com">External insecure</a>
        <a href="https://external.com">External secure</a>
        <a href="#main">Fragment only</a>
        <a href="https://example.com#main">Internal with fragment</a>
        <a href="https://example.com/about">Internal with path</a>
        <a href="mailto:example@example.com">Mailto link</a>
    </body></html>`

	baseURL, err := url.Parse("https://example.com")
	require.NoError(t, err)

	links, err := crawler.ExtractLinks(baseURL, strings.NewReader(html))
	require.NoError(t, err)
	require.Len(t, links, 8)

	expected := []string{
		"https://example.com",
		"https://example.com/about",
		"http://external.com",
		"https://external.com",
		"https://example.com#fragment",
		"https://example.com#main",
		"https://example.com/about",
		"mailto:example@example.com",
	}

	for i, link := range links {
		assert.Equal(t, expected[i], link.String())
	}
}
