package crawler_test

import (
	"net/url"
	"testing"

	"example.com/crawler"
	"github.com/stretchr/testify/assert"
)

func TestNormaliseURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "http converted to https",
			inputURL: "http://example.com",
			expected: "https://example.com/",
		},
		{
			name:     "fragment removed",
			inputURL: "https://example.com/page#section1",
			expected: "https://example.com/page/",
		},
		{
			name:     "trailing slash added",
			inputURL: "https://example.com/page",
			expected: "https://example.com/page/",
		},
		{
			name:     "host, path and scheme normalized to lowercase",
			inputURL: "HTTPS://EXAMPLE.COM/PAGE",
			expected: "https://example.com/page/",
		},
	}

	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			u, err := url.Parse(testData.inputURL)
			assert.NoError(t, err)

			normalized := crawler.NormaliseURL(u)
			assert.Equal(t, testData.expected, normalized.String())
		})
	}
}
