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
