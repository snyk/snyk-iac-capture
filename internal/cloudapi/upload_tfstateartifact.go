package cloudapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

type UploadTFStateArtifactRequest struct {
	Data Data `json:"data"`
}

type Data struct {
	Attributes terraform.State `json:"attributes"`
	Type       string          `json:"type"`
}

func (c *Client) UploadTFStateArtifact(ctx context.Context, orgID string, artifact *terraform.State) (e error) {
	request := UploadTFStateArtifactRequest{
		Data: Data{
			Attributes: *artifact,
			Type:       "tfstate",
		},
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(request); err != nil {
		return err
	}

	url := fmt.Sprintf("%s/hidden/orgs/%s/cloud/mappings_artifact/tfstate", c.url, orgID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Set("version", c.version)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", c.authorization)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil && e == nil {
			e = err
		}
	}()

	if res.StatusCode != http.StatusCreated {
		return extractError(res)
	}

	return nil
}
