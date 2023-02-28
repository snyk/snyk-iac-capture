package pkg

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/snyk/go-application-framework/pkg/configuration"
)

const (
	SNYK_API_REST_URL = "snyk_api_rest_url"
	SNYK_API_V3_URL   = "snyk_api_v3_url"
)

func GetRestApiUrl(config configuration.Configuration) (string, error) {
	restURL := config.GetString(SNYK_API_REST_URL)
	if restURL != "" {
		return restURL, nil
	}
	v3URL := config.GetString(SNYK_API_V3_URL)
	if v3URL != "" {
		return v3URL, nil
	}

	// REST API URL should always look like this: https://api.$DOMAIN/rest
	apiURL := config.GetString(configuration.API_URL)
	parsedBaseUrl, err := url.Parse(apiURL)
	if err != nil {
		return "", fmt.Errorf("could not parse api url %s: %+v", apiURL, err)
	}
	parsedBaseUrl.Path = "/rest"
	if strings.HasPrefix(parsedBaseUrl.Host, "app") {
		// Rewrite app.snyk.io/ to api.snyk.io/rest
		parsedBaseUrl.Host = strings.Replace(parsedBaseUrl.Host, "app.", "api.", 1)
	} else if !strings.HasPrefix(parsedBaseUrl.Host, "localhost") && strings.HasPrefix(parsedBaseUrl.Host, "api") { // Ignore localhosts and URLs with api. already defined
		parsedBaseUrl.Host = fmt.Sprintf("api.%s", parsedBaseUrl.Host)
	}

	return parsedBaseUrl.String(), nil
}
