package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetPage call the client page by HTTP request and extract the body to HTML document.
func GetPage(ctx context.Context, method, siteURL string, cookies []*http.Cookie, headers, formDatas map[string]string, timeout int) (*goquery.Document, []*http.Cookie, error) {
	// This function can handle both all methods.
	// Initiate this body variable as nil for method that doesn't required body.
	body := io.Reader(nil)
	// If the request contain form-data, add the form-data parameters to the body.
	if len(formDatas) > 0 {
		form := url.Values{}
		for k, v := range formDatas {
			form.Add(k, v)
		}

		body = strings.NewReader(form.Encode())
	}

	// Create a new HTTP request with context.
	req, err := http.NewRequestWithContext(ctx, method, siteURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request context: %w", err)
	}
	// If the request contain headers, add the header parameters.
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	// If the request contain cookies, add the cookie parameters.
	if len(cookies) > 0 {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}

	// Use the default timeout if the timeout parameter isn't configured.
	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}

	// Use default http Client.
	httpClient := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       reqTimeout,
	}

	// Execute the request.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	// Close the response body
	defer func() { _ = resp.Body.Close() }()
	// // Parsing response body to HTML document reader.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse html: %w", err)
	}

	// Return HTML doc, cookies.
	return doc, resp.Cookies(), nil
}
