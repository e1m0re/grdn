package apiclient

import (
	"bytes"
	"compress/gzip"
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
	cBody, err := compressBody(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, api.baseURL+path, cBody)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	resp, err := api.client.Do(request)
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}

func compressBody(content *[]byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	compressor := gzip.NewWriter(&buf)
	if _, err := compressor.Write(*content); err != nil {
		return nil, err
	}

	if err := compressor.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}
