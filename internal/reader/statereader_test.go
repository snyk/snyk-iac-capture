package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadState(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "successfully read a valid state file",
			path: "testdata/valid-terraform.tfstate",
		},
		{
			name:    "error empty state",
			path:    "testdata/empty-terraform.tfstate",
			wantErr: true,
		},
		{
			name:    "error invalid state",
			path:    "testdata/invalid.tfstate",
			wantErr: true,
		},
		{
			name:    "error while reading a directory",
			path:    "testdata/directory-terraform.tfstate",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfState, err := ReadStateFile(tt.path)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, tfState)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, tfState)
		})
	}
}
