package snyk_iac_capture

import (
	"context"
	"fmt"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
	"github.com/snyk/snyk-iac-capture/internal/filefinder"
	"github.com/snyk/snyk-iac-capture/internal/filtering"
	"github.com/snyk/snyk-iac-capture/internal/http"
	"github.com/snyk/snyk-iac-capture/internal/reader"
)

type Command struct {
	Org               string
	StatePath         string
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
	logrus.Debugf("Looking for Terraform states in '%s'", c.StatePath)
	files, err := filefinder.FindFiles(c.StatePath, path.Join("**", "*.tfstate"))
	if err != nil {
		return nil, fmt.Errorf("error looking for Terraform states in '%s': %v", c.StatePath, err)
	}
	if len(files) <= 0 {
		return nil, fmt.Errorf("could not find any Terraform state in '%s'", c.StatePath)
	}
	logrus.Debugf("Found %+v\n", files)

	httpClient, err := http.NewClient(
		http.WithTLSSkipVerify(c.HTTPTLSSkipVerify),
		http.WithExtraCertificates(c.ExtraSSlCerts),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %v", err)
	}
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

	var captured []string
	for _, file := range files {
		logrus.Debugf("Capturing '%s'", file)
		err := captureState(file, cloudApiClient)
		if err != nil {
			return nil, err
		}
		captured = append(captured, file)
	}

	return captured, nil
}

func captureState(statePath string, cloudApiClient *cloudapi.Client) error {
	// read state file
	tfState, err := reader.ReadState(statePath)
	if err != nil {
		return fmt.Errorf("error reading Terraform state %s: %v", statePath, err)
	}

	// call filter
	stateArtifact, err := filtering.FilterState(tfState)
	if err != nil {
		return fmt.Errorf("error filtering Terraform state %s: %v", statePath, err)
	}

	// send artifact to cloud api
	err = cloudApiClient.UploadTFStateArtifact(context.TODO(), stateArtifact)
	if err != nil {
		return fmt.Errorf("error uploading state artifact for '%s': %v", statePath, err)
	}
	return nil
}
