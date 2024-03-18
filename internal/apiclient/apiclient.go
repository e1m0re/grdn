package apiclient

import (
	"bytes"
	"net/http"
)

type API struct {
	client  *http.Client
	baseURL string
}

func NewAPI(baseURL string) *API {
	return &API{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (api *API) DoRequest(path string) ([]byte, error) {
	var buf []byte
	responseBody := bytes.NewBuffer(buf)

	resp, err := api.client.Post(api.baseURL+path, "text/plain", responseBody)
	if err == nil {
		defer resp.Body.Close()
	}

	return buf, err
}
