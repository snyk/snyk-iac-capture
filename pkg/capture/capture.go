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

package capture

import (
	"context"
	"fmt"
	"log"
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
		return nil, fmt.Errorf("error looking for Terraform states in '%s': %w", statePath, err)
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
		return fmt.Errorf("error reading Terraform state %s: %w", statePath, err)
	}

	if err = CaptureState(tfState, cloudApiClient); err != nil {
		return fmt.Errorf("error capturing Terraform state '%s': %+v", statePath, err)
	}
	return nil
}

func CaptureStateFromStdin(cloudApiClient *cloudapi.Client, logger *log.Logger) ([]string, error) {
	logger.Println("Reading state from stdin")
	state, err := reader.ReadStateFromStdin()
	if err != nil {
		return nil, err
	}
	return []string{state.Lineage}, CaptureState(state, cloudApiClient)
}

func CaptureState(tfState *terraform.State, cloudApiClient *cloudapi.Client) error {
	// call filter
	stateArtifact, err := filtering.FilterState(tfState)
	if err != nil {
		return fmt.Errorf("unable to filter: %w", err)
	}

	// send artifact to cloud api
	err = cloudApiClient.UploadTFStateArtifact(context.Background(), stateArtifact)
	if err != nil {
		return fmt.Errorf("unable to upload artifact: %w", err)
	}
	return nil
}
