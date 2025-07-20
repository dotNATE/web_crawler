package crawler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"example.com/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)

	html := `
		<html>
			<body>
				<a href="https://example.com/page1">Page 1</a>
				<a href="https://example.com/page2">Page 2</a>
			</body>
		</html>
	`
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(html)),
	}

	return resp, args.Error(1)
}

func TestCrawl(t *testing.T) {
	mockHttpClient := new(MockHTTPClient)
	mockHttpClient.On("Get", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body: io.NopCloser(bytes.NewBufferString(`
			<html>
				<body>
					<a href="https://example.com/page1">Page 1</a>
					<a href="https://example.com/page2">Page 2</a>
				</body>
			</html>
		`)),
	}, nil)

	c, err := crawler.NewCrawler("https://example.com", mockHttpClient, 5, 0)
	if err != nil {
		t.Errorf("Failed to initialise Crawler: error: %+v", err)
	}

	c.Crawl(c.Base, 0)
	c.WaitGroup.Wait()

	require.Len(t, c.Results, 3)
	assert.Equal(t, c.Results["https://example.com/"][0], "https://example.com/page1")
	assert.Equal(t, c.Results["https://example.com/"][1], "https://example.com/page2")
}

func TestShouldVisit(t *testing.T) {
	base := "https://monzo.com"

	tests := []struct {
		name     string
		link     string
		expected bool
	}{
		{
			name:     "Same domain",
			link:     "https://monzo.com/about",
			expected: true,
		},
		{
			name:     "Subdomain (should skip)",
			link:     "https://community.monzo.com",
			expected: false,
		},
		{
			name:     "External domain",
			link:     "https://example.com",
			expected: false,
		},
		{
			name:     "Mailto only link",
			link:     "mailto:support@monzo.com",
			expected: false,
		},
		{
			name:     "Javascript link",
			link:     "javascript:void(0)",
			expected: false,
		},
		{
			name:     "Fragment-only link",
			link:     "#section1",
			expected: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockHttpClient := new(MockHTTPClient)
			c, err := crawler.NewCrawler(base, mockHttpClient, 2, 0)
			if err != nil {
				t.Errorf("Failed to initialise Crawler: error: %+v", err)
			}

			linkURL, err := url.Parse(testCase.link)
			assert.NoError(t, err, "should parse URL")

			result := c.ShouldVisit(linkURL)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
