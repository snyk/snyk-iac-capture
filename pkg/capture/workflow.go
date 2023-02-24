package capture

import (
	"fmt"
	"log"

	"github.com/snyk/go-application-framework/pkg/workflow"
	"github.com/spf13/pflag"

	"github.com/snyk/snyk-iac-capture/internal/reader"

	"github.com/snyk/snyk-iac-capture/pkg"

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
	org := config.GetString("org")
	token := config.GetString("token")       // TODO should be removed after https://snyksec.atlassian.net/browse/HEAD-141 is fixed
	apiUrl, err := pkg.GetRestApiUrl(config) // TODO use the existing config when my PR is merged
	if err != nil {
		return nil, err
	}

	if statePath == "" {
		statePath = "." // Cannot register default to positional argument so setting default here
	}

	captured, err := capture(statePath, apiUrl, token, org, readFromStdin, ictx, logger)
	fmt.Printf("Captured Terraform states: %+v\n", captured)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully captured all your states.")
	return data, nil
}

func capture(statePath, apiUrl, token, org string, stateFromStdin bool, ictx workflow.InvocationContext, logger *log.Logger) ([]string, error) {
	cloudApiClient, err := cloudapi.NewClient(cloudapi.ClientConfig{
		HTTPClient:     ictx.GetNetworkAccess().GetHttpClient(),
		URL:            apiUrl,
		Version:        "2022-04-13~experimental",
		OrganisationID: org,
		Authorization:  fmt.Sprintf("token %s", token), // TODO remove when the client is fixed
	})
	if err != nil {
		return nil, fmt.Errorf("error creating CloudAPI client: %v", err)
	}
	logger.Println("CloudApiClient created...")

	if !stateFromStdin {
		return CaptureStatesFromPath(statePath, cloudApiClient, logger)
	}

	logger.Println("Reading state from stdin")
	state, err := reader.ReadStateFromStdin()
	if err != nil {
		return nil, err
	}
	return []string{state.Lineage}, CaptureState(state, cloudApiClient)
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
