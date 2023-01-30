package snyk_iac_capture

import (
	"github.com/snyk/snyk-iac-capture/internal/filtering"
	"github.com/snyk/snyk-iac-capture/internal/reader"
)

type Command struct {
	Org       string
	StateFile string
}

func (c Command) Run() int {
	// read state file
	tfState, err := reader.ReadState(c.StateFile)
	if err != nil {
		return 1
	}

	// call filter
	stateArtifact, err := filtering.FilterState(tfState)
	if err != nil {
		return 1
	}

	// TODO send to cloud-api-service
	_ = stateArtifact
	return 0
}
