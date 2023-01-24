package filtering

import (
	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

type TfStateArtifact struct {
}

func Filter(_ *terraform.State) (TfStateArtifact, error) {
	return TfStateArtifact{}, nil
}
