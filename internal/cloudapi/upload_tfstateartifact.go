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

func (c *Client) UploadTFStateArtifact(ctx context.Context, artifact *terraform.State) (e error) {
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

	url := fmt.Sprintf("%s/hidden/orgs/%s/cloud/mappings_artifact/tfstate", c.url, c.organisationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Set("version", c.version)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Content-Type", "application/vnd.api+json")
	if c.authorization != "" {
		req.Header.Set("Authorization", c.authorization)
	}

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
