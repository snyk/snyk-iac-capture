package capture

import (
	"context"
	"fmt"
	"path"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
	"github.com/snyk/snyk-iac-capture/internal/filefinder"
	"github.com/snyk/snyk-iac-capture/internal/filtering"
	"github.com/snyk/snyk-iac-capture/internal/reader"
	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func CaptureStatesFromPath(statePath string, cloudApiClient *cloudapi.Client, logger Logger) ([]string, error) {
	logger.Printf("Looking for Terraform states in '%s'\n", statePath)
	files, err := filefinder.FindFiles(statePath, path.Join("**", "*.tfstate"))
	if err != nil {
		return nil, fmt.Errorf("error looking for Terraform states in '%s': %v", statePath, err)
	}
	if len(files) <= 0 {
		return nil, fmt.Errorf("could not find any Terraform state in '%s'", statePath)
	}
	logger.Printf("Found %+v\n", files)

	var captured []string
	for _, file := range files {
		logger.Printf("Capturing '%s'\n", file)
		err := CaptureStateFromPath(file, cloudApiClient)
		if err != nil {
			return nil, err
		}
		captured = append(captured, file)
	}
	return captured, nil
}

func CaptureStateFromPath(statePath string, cloudApiClient *cloudapi.Client) error {
	// read state file
	tfState, err := reader.ReadStateFile(statePath)
	if err != nil {
		return fmt.Errorf("error reading Terraform state %s: %v", statePath, err)
	}

	err = CaptureState(tfState, cloudApiClient)
	if err != nil {
		return fmt.Errorf("error capturing Terraform state '%s': %+v", statePath, err)
	}
	return nil
}

func CaptureState(tfState *terraform.State, cloudApiClient *cloudapi.Client) error {
	// call filter
	stateArtifact, err := filtering.FilterState(tfState)
	if err != nil {
		return fmt.Errorf("unable to filter: %v", err)
	}

	// send artifact to cloud api
	err = cloudApiClient.UploadTFStateArtifact(context.TODO(), stateArtifact)
	if err != nil {
		return fmt.Errorf("unable to upload artifact: %v", err)
	}
	return nil
}
