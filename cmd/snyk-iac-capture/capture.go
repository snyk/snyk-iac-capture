package snyk_iac_capture

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
	"github.com/snyk/snyk-iac-capture/internal/filtering"
	"github.com/snyk/snyk-iac-capture/internal/http"
	"github.com/snyk/snyk-iac-capture/internal/reader"
)

type Command struct {
	Org               string
	StateFile         string
	HTTPTLSSkipVerify bool
	APIURL            string
	APIToken          string
}

func (c Command) Run() int {
	if c.StateFile == "" {
		log.Println("error: missing required statefile argument")
		return 1
	}

	httpClient, err := http.NewClient(
		http.WithTLSSkipVerify(c.HTTPTLSSkipVerify),
		http.WithExtraCertificates(os.Getenv("NODE_EXTRA_CA_CERTS")),
	)
	if err != nil {
		log.Printf("error: create HTTP client: %v\n", err)
		return 1
	}
	cloudapiClient, err := cloudapi.NewClient(cloudapi.ClientConfig{
		HTTPClient:    httpClient,
		URL:           c.APIURL,
		Authorization: fmt.Sprintf("token %s", c.APIToken),
		Version:       "2022-04-13~experimental",
	})
	if err != nil {
		log.Printf("error: create CloudAPI client: %v\n", err)
		return 1
	}

	// read state file
	tfState, err := reader.ReadState(c.StateFile)
	if err != nil {
		log.Printf("error: cannot read state %s: %v\n", c.StateFile, err)
		return 1
	}

	// call filter
	stateArtifact, err := filtering.FilterState(tfState)
	if err != nil {
		log.Printf("error: cannot filter state %s: %v\n", c.StateFile, err)
		return 1
	}

	// send artifact to cloud api
	err = cloudapiClient.UploadTFStateArtifact(context.TODO(), c.Org, stateArtifact)
	if err != nil {
		log.Printf("error: there was an error uploading state artifact: %v\n", err)
		return 1
	}
	return 0
}
