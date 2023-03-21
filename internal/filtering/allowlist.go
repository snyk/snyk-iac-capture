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

type ResourceAllowlist map[string][]string

var globalAllowlist = ResourceAllowlist{
	"aws_iam_group_policy_attachment": []string{"policy_arn", "group"},
	"aws_iam_policy_attachment":       []string{"policy_arn", "users", "groups", "roles"},
	"aws_iam_role_policy_attachment":  []string{"policy_arn", "role"},
	"aws_iam_user_policy_attachment":  []string{"policy_arn", "user"},
	"aws_s3_bucket_acl":               []string{"bucket"},
	"aws_security_group_rule":         []string{"security_group_id"},
}

func (a ResourceAllowlist) GetAllowedAttributes(ty string) []string {
	allowedAttributes := []string{"id"}
	if attributes, ok := a[ty]; ok {
		allowedAttributes = append(allowedAttributes, attributes...)
	}
	return allowedAttributes
}
