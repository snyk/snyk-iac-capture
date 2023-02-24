package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

func ReadStateFile(path string) (*terraform.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file, please check permissions: %+v", err)
	}

	return readState(data)
}

func ReadStateFromStdin() (*terraform.State, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error reading standard input: %v", err)
	}
	state, err := readState(data)
	if err != nil {
		return nil, fmt.Errorf("error reading Terraform state from standard input: %v", err)
	}
	return state, nil
}

func readState(data []byte) (*terraform.State, error) {
	var tfState terraform.State
	if err := json.Unmarshal(data, &tfState); err != nil {
		return nil, fmt.Errorf("invalid format, please check that the state is in correct json format: %+v", err)
	}
	return &tfState, nil
}
