package snyk_iac_capture

import (
	"fmt"
	"log"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
	"github.com/snyk/snyk-iac-capture/internal/http"
	"github.com/snyk/snyk-iac-capture/internal/reader"
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
		fmt.Printf("An error occured: %+v\n", err)
		return 1
	}
	fmt.Println("Successfully captured all your states.")
	return 0
}

func (c *Command) capture() ([]string, error) {
	c.Logger.Println("Start capture...")
	httpClient, err := http.NewClient(
		http.WithTLSSkipVerify(c.HTTPTLSSkipVerify),
		http.WithExtraCertificates(c.ExtraSSlCerts),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %v", err)
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
		return nil, fmt.Errorf("error creating CloudAPI client: %v", err)
	}
	c.Logger.Println("CloudApiClient created...")

	if !c.StateFromStdin {
		return capture.CaptureStatesFromPath(c.StatePath, cloudApiClient, c.Logger)
	}

	c.Logger.Println("Reading state from stdin")
	state, err := reader.ReadStateFromStdin()
	if err != nil {
		return nil, err
	}
	return []string{state.Lineage}, capture.CaptureState(state, cloudApiClient)
}
