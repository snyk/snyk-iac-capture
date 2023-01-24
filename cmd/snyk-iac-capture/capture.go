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
	// read statefile
	stateReader := reader.NewStateReader(c.StateFile)
	tfState, err := stateReader.Read()
	if err != nil {
		return 1
	}

	// call filter
	_, err = filtering.Filter(tfState)
	if err != nil {
		return 1
	}

	// send to cloud-api-service

	return 0
}
