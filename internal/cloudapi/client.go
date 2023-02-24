package cloudapi

import (
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientConfig struct {
	HTTPClient     HTTPClient
	URL            string
	Authorization  string
	Version        string
	OrganisationID string
}

type Client struct {
	httpClient     HTTPClient
	url            string
	authorization  string
	version        string
	organisationID string
}

func NewClient(config ClientConfig) (*Client, error) {
	httpClient := config.HTTPClient

	/*if httpClient == nil {
		httpClient = http.DefaultClient
	}*/

	if config.URL == "" {
		return nil, fmt.Errorf("no URL provided")
	}

	if config.Authorization == "" {
		return nil, fmt.Errorf("no Authorization provided")
	}

	if config.Version == "" {
		return nil, fmt.Errorf("no version provided")
	}

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	sanitizedURL := url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
	}

	client := Client{
		httpClient:     httpClient,
		url:            sanitizedURL.String(),
		authorization:  config.Authorization,
		version:        config.Version,
		organisationID: config.OrganisationID,
	}

	return &client, nil
}
