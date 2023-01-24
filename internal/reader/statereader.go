package reader

import (
	"encoding/json"
	"os"

	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

func ReadState(path string) (*terraform.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tfState terraform.State
	if err := json.Unmarshal(data, &tfState); err != nil {
		return nil, err
	}

	return &tfState, nil
}
