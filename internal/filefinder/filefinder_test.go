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

package filefinder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFiles(t *testing.T) {
	type args struct {
		p          string
		endPattern string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"Directory as input",
			args{
				p:          "test/",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate", "test/sub/terraform.tfstate", "test/sub/subsub/terraform.tfstate", "test/terraform.tfstate/terraform.tfstate"},
			false,
		},
		{
			"**/*.tfstate as input",
			args{
				p:          "test/**/*.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate", "test/sub/terraform.tfstate", "test/sub/subsub/terraform.tfstate", "test/terraform.tfstate/terraform.tfstate"},
			false,
		},
		{
			"*.tfstate as input",
			args{
				p:          "test/*.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate"},
			false,
		},
		{
			"test.tfstate as input",
			args{
				p:          "test/test.tfstate",
				endPattern: "**/*.tfstate",
			},
			[]string{"test/test.tfstate"},
			false,
		},
		{
			"**/*.notexist as input",
			args{
				p:          "test/**/*.notexist",
				endPattern: "**/*.tfstate",
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindFiles(tt.args.p, tt.args.endPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
