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
	"fmt"
	"log"

	"github.com/snyk/go-application-framework/pkg/configuration"
	"github.com/snyk/go-application-framework/pkg/workflow"
	"github.com/spf13/pflag"

	"github.com/snyk/snyk-iac-capture/internal/cloudapi"
)

var (
	WorkflowID            = workflow.NewWorkflowIdentifier("iac.capture")
	ReadFromStdinFlag     = "stdin"
	TargetDirectoryConfig = "targetDirectory"
)

func CaptureWorkflow(
	ictx workflow.InvocationContext,
	_ []workflow.Data,
) (data []workflow.Data, err error) {
	config := ictx.GetConfiguration()
	logger := ictx.GetLogger()

	logger.Println("CaptureWorkflow start")

	readFromStdin := config.GetBool(ReadFromStdinFlag)
	statePath := config.GetString(TargetDirectoryConfig)
	org := config.GetString(configuration.ORGANIZATION)
	apiUrl := config.GetString(configuration.API_URL)

	if statePath == "" {
		statePath = "." // Cannot register default to positional argument so setting default here
	}

	captured, err := capture(statePath, apiUrl, org, readFromStdin, ictx, logger)
	fmt.Printf("Captured Terraform states: %+v\n", captured)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully captured all your states.")
	return data, nil
}

func capture(statePath, apiUrl, org string, stateFromStdin bool, ictx workflow.InvocationContext, logger *log.Logger) ([]string, error) {
	cloudApiClient, err := cloudapi.NewClient(cloudapi.ClientConfig{
		HTTPClient:     ictx.GetNetworkAccess().GetHttpClient(),
		URL:            apiUrl,
		Version:        "2022-04-13~experimental",
		OrganisationID: org,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating CloudAPI client: %w", err)
	}
	logger.Println("CloudApiClient created...")

	if !stateFromStdin {
		return CaptureStatesFromPath(statePath, cloudApiClient, logger)
	}

	return CaptureStateFromStdin(cloudApiClient, logger)
}

func Init(e workflow.Engine) error {
	flagset := pflag.NewFlagSet("snyk-cli-extension-capture", pflag.ExitOnError)

	flagset.Bool(ReadFromStdinFlag, false, "Read states from standard input instead of using target directory.")

	c := workflow.ConfigurationOptionsFromFlagset(flagset)

	if _, err := e.Register(WorkflowID, c, CaptureWorkflow); err != nil {
		return fmt.Errorf("error while registering Capture workflow: %w", err)
	}

	return nil
}
