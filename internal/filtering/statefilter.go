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
	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

func FilterState(state *terraform.State) (*terraform.State, error) {
	artifact := terraform.State{
		Version:          state.Version,
		TerraformVersion: state.TerraformVersion,
		Lineage:          state.Lineage,
		Resources:        []terraform.Resource{},
	}

	for _, resource := range state.Resources {
		if resource.Mode != "managed" {
			continue
		}
		allowedAttributes := globalAllowlist.GetAllowedAttributes(resource.Type)
		var instances []terraform.ResourceInstance
		for _, instance := range resource.Instances {
			attributes := map[string]interface{}{}
			for _, attr := range allowedAttributes {
				if _, exists := instance.Attributes[attr]; exists {
					attributes[attr] = instance.Attributes[attr]
				}
			}
			instance.Attributes = attributes
			instances = append(instances, instance)
		}
		resource.Instances = instances
		artifact.Resources = append(artifact.Resources, resource)
	}

	return &artifact, nil
}
