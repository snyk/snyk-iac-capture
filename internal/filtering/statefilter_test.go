package filtering

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snyk/snyk-iac-capture/internal/reader"
)

func TestFilter(t *testing.T) {

	tests := []struct {
		name           string
		stateFile      string
		resultJsonFile string
		wantErr        bool
	}{
		{
			name:           "FilterState a terraform state file with multiple use cases",
			stateFile:      "full.json",
			resultJsonFile: "full-filtered.json",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			read, err := reader.ReadState(filepath.Join("testdata/", tt.stateFile))
			assert.Nil(t, err)

			got, err := FilterState(read)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			read, err = reader.ReadState(filepath.Join("testdata/", tt.resultJsonFile))
			assert.Nil(t, err)

			assert.Equal(t, read, got)
		})
	}
}
