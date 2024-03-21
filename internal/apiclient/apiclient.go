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

func (api *API) DoRequest(path string, body *[]byte) error {

	resp, err := api.client.Post(api.baseURL+path, "text/plain", bytes.NewBuffer(*body))
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}
