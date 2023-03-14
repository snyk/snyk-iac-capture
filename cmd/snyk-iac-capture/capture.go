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

package snyk_iac_capture

import (
	"fmt"
	"log"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
	"github.com/snyk/snyk-iac-capture/internal/http"
	"github.com/snyk/snyk-iac-capture/pkg/capture"
)

type Command struct {
	Logger            *log.Logger
	Org               string
	StatePath         string
	StateFromStdin    bool
	HTTPTLSSkipVerify bool
	APIURL            string
	APIToken          string
	ExtraSSlCerts     string
}

func (c *Command) Run() int {
	captured, err := c.capture()
	fmt.Printf("Captured Terraform states: %+v\n", captured)

	if err != nil {
		fmt.Printf("%+v", err)
		return 1
	}
	fmt.Println("Successfully captured all your states.")
	return 0
}

func (c *Command) capture() ([]string, error) {
	c.Logger.Println("Start capturing...")
	httpClient, err := http.NewClient(
		http.WithTLSSkipVerify(c.HTTPTLSSkipVerify),
		http.WithExtraCertificates(c.ExtraSSlCerts),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %w", err)
	}
	c.Logger.Println("Http Client created...")

	cloudApiClient, err := cloudapi.NewClient(cloudapi.ClientConfig{
		HTTPClient:     httpClient,
		URL:            c.APIURL,
		Authorization:  fmt.Sprintf("token %s", c.APIToken),
		Version:        "2022-04-13~experimental",
		OrganisationID: c.Org,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating CloudAPI client: %w", err)
	}
	c.Logger.Println("CloudApiClient created...")

	if !c.StateFromStdin {
		return capture.CaptureStatesFromPath(c.StatePath, cloudApiClient, c.Logger)
	}

	return capture.CaptureStateFromStdin(cloudApiClient, c.Logger)
}
