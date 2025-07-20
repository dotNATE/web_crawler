package crawler_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"example.com/crawler"
	"github.com/stretchr/testify/mock"
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
	}, nil).Once()

	c, err := crawler.NewCrawler("https://example.com", mockHttpClient)
	if err != nil {
		t.Errorf("Failed to initialise Crawler: error: %+v", err)
	}

	c.Crawl(c.Base)
}
