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

package cloudapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Title      string
	Detail     string
}

func NewHttpError(statusCode int, title string, detail string) *HttpError {
	return &HttpError{StatusCode: statusCode, Title: title, Detail: detail}
}

func (h *HttpError) Error() string {
	if h.Detail == "" {
		return fmt.Sprintf("%d - %s", h.StatusCode, h.Title)
	}
	return fmt.Sprintf("%d - %s: %s", h.StatusCode, h.Title, h.Detail)
}

type apiResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []apiError  `json:"errors,omitempty"`
}

type apiError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func extractError(res *http.Response) error {
	var apiResp apiResponse
	err := json.NewDecoder(res.Body).Decode(&apiResp)
	if err != nil || len(apiResp.Errors) <= 0 {
		return NewHttpError(res.StatusCode, res.Status, "")
	}

	// only extract the first error.
	apiErr := apiResp.Errors[0]
	return NewHttpError(res.StatusCode, apiErr.Title, apiErr.Detail)
}
