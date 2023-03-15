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
	"aws_security_group_rule": []string{"security_group_id"},
	"aws_s3_bucket_acl":       []string{"bucket"},
}

func (a ResourceAllowlist) GetAllowedAttributes(ty string) []string {
	allowedAttributes := []string{"id"}
	if attributes, ok := a[ty]; ok {
		allowedAttributes = append(allowedAttributes, attributes...)
	}
	return allowedAttributes
}
