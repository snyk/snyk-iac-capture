/*
 * © 2023 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if config.URL == "" {
		return nil, fmt.Errorf("no URL provided")
	}

	if config.OrganisationID == "" {
		return nil, fmt.Errorf("no OrganisationID provided")
	}

	if config.Version == "" {
		return nil, fmt.Errorf("no version provided")
	}

	parsedURL, err := url.Parse(config.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	client := Client{
		httpClient:     httpClient,
		url:            parsedURL.String(),
		authorization:  config.Authorization,
		version:        config.Version,
		organisationID: config.OrganisationID,
	}

	return &client, nil
}
