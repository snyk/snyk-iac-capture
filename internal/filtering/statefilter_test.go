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
		{
			name:           "Filter a valid but empty state",
			stateFile:      "empty.json",
			resultJsonFile: "empty-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with security group rules",
			stateFile:      "aws_security_group_rule.json",
			resultJsonFile: "aws_security_group_rule-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with S3 bucket ACLs",
			stateFile:      "aws_s3_bucket_acl.json",
			resultJsonFile: "aws_s3_bucket_acl-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with IAM group policy attachment",
			stateFile:      "aws_iam_group_policy_attachment.json",
			resultJsonFile: "aws_iam_group_policy_attachment-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with IAM role policy attachment",
			stateFile:      "aws_iam_role_policy_attachment.json",
			resultJsonFile: "aws_iam_role_policy_attachment-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with IAM user policy attachment",
			stateFile:      "aws_iam_user_policy_attachment.json",
			resultJsonFile: "aws_iam_user_policy_attachment-filtered.json",
			wantErr:        false,
		},
		{
			name:           "Filter a state with IAM policy attachment",
			stateFile:      "aws_iam_policy_attachment.json",
			resultJsonFile: "aws_iam_policy_attachment-filtered.json",
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			read, err := reader.ReadStateFile(filepath.Join("testdata/", tt.stateFile))
			assert.Nil(t, err)

			got, err := FilterState(read)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			read, err = reader.ReadStateFile(filepath.Join("testdata/", tt.resultJsonFile))
			assert.Nil(t, err)

			assert.Equal(t, read, got)
		})
	}
}
