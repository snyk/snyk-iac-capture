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
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/snyk/snyk-iac-capture/internal/terraform"
)

func ReadStateFile(path string) (*terraform.State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return readState(data)
}

func ReadStateFromStdin() (*terraform.State, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error reading standard input: %w", err)
	}
	state, err := readState(data)
	if err != nil {
		return nil, fmt.Errorf("error reading Terraform state from standard input: %w", err)
	}
	return state, nil
}

func readState(data []byte) (*terraform.State, error) {
	var tfState terraform.State
	if err := json.Unmarshal(data, &tfState); err != nil {
		return nil, fmt.Errorf("invalid format, please check that the state is in correct json format: %w", err)
	}
	return &tfState, nil
}
