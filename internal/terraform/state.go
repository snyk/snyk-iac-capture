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

package terraform

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type State struct {
	Version          int        `json:"version"`
	TerraformVersion string     `json:"terraform_version"`
	Resources        []Resource `json:"resources"`
	Lineage          string     `json:"lineage"`
}

// UnmarshalJSON Custom deserializer to Throw error when a field not marked as omitempty is missing
func (s *State) UnmarshalJSON(bytes []byte) error {
	type tfState State

	if err := json.Unmarshal(bytes, (*tfState)(s)); err != nil {
		return err
	}

	fields := reflect.ValueOf(s).Elem()
	for i := 0; i < fields.NumField(); i++ {

		jsonTags := fields.Type().Field(i).Tag.Get("json")
		if !strings.Contains(jsonTags, "omitempty") && fields.Field(i).IsZero() {
			return fmt.Errorf("invalid state: required field '%s' is missing", fields.Type().Field(i).Name)
		}

	}
	return nil
}

type Resource struct {
	Module    *string            `json:"module,omitempty"`
	Mode      string             `json:"mode"`
	Type      string             `json:"type"`
	Name      string             `json:"name"`
	Instances []ResourceInstance `json:"instances"`
}

type ResourceInstance struct {
	IndexKey   any                    `json:"index_key,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}
