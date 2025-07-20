package crawler_test

import (
	"net/url"
	"strings"
	"testing"

	"example.com/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractLinks_GetATags(t *testing.T) {
	html := `<html><body>
		<a href="https://example.com">Internal secure</a>
		<div>Should be ignored</div>
		<div href="/error">Should be ignored</div>
        <a href="/about">Internal by reference</a>
        <a href="http://external.com">External insecure</a>
        <a href="https://external.com">External secure</a>
        <a href="#main">Fragment only</a>
        <a href="https://example.com#main">Internal with fragment</a>
        <a href="https://example.com/about">Internal with path</a>
        <a href="mailto:example@example.com">Mailto link</a>
        <a attribute="first" href="https://example.com">Multiple attributes internal</a>
    </body></html>`

	baseURL, err := url.Parse("https://example.com")
	require.NoError(t, err)

	links, err := crawler.ExtractLinks(baseURL, strings.NewReader(html))
	require.NoError(t, err)

	expected := []string{
		"https://example.com",
		"https://example.com/about",
		"http://external.com",
		"https://external.com",
		"https://example.com#main",
		"https://example.com#main",
		"https://example.com/about",
		"mailto:example@example.com",
		"https://example.com",
	}

	require.Len(t, links, len(expected))
	for i, link := range links {
		assert.Equal(t, expected[i], link.String())
	}
}
