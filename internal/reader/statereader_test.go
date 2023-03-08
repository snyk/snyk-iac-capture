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

package reader

import (
	"os"
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

func TestReadStateFromStdin(t *testing.T) {
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

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }() // Restore original Stdin

			file, err := os.Open(tt.path)
			assert.Nil(t, err)
			os.Stdin = file
			tfState, err := ReadStateFromStdin()

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
