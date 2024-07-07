// Copyright 2024 Bill Nixon. All rights reserved.
// Use of this source code is governed by the license found in the LICENSE file.

/*
Package sitestatus can be used to check if websites are up and responding.

There are options to ignore certificates and redirects.
*/
package sitestatus

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"
)

// IsValidURL determines if s is a valid URL
func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// HTTPClientOptions defines options for HTTP client configuration.
type HTTPClientOptions struct {
	IgnoreCerts     bool
	IgnoreRedirects bool
	Timeout         time.Duration
}

// NewHTTPClient creates an http.Client with the given options.
func NewHTTPClient(opt HTTPClientOptions) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opt.IgnoreCerts,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   opt.Timeout,
	}

	if opt.IgnoreRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

// Check determines if a site is up by performing a simple GET request.
// It returns the status of the site.
func Check(site string, opts HTTPClientOptions) string {
	client := NewHTTPClient(opts)

	resp, err := client.Get(site)
	if err != nil {
		// Check if the error is a timeout
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			return "DOWN timeout after " + opts.Timeout.String()
		}

		return "DOWN " + err.Error()
	}
	defer resp.Body.Close()

	// Read and discard the response body to ensure connection reuse
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return "DOWN " + err.Error()
	}

	return "UP " + resp.Status
}
