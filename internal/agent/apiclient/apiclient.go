package apiclient

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"log/slog"
	"net/http"
)

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

type APIClient struct {
	client  *http.Client
	baseURL string
	key     []byte
}

func NewAPI(baseURL string, key []byte) *APIClient {
	return &APIClient{
		client:  &http.Client{},
		baseURL: baseURL,
		key:     key,
	}
}

func (api *APIClient) DoRequest(request *http.Request) (*http.Response, error) {
	response, err := api.client.Do(request)
	if err != nil {
		slog.Warn("error while doing request",
			slog.String("url", request.URL.String()),
			slog.String("method", request.Method),
			slog.String("error", err.Error()),
		)
	}

	if err == nil && response.StatusCode < 500 {
		return response, err
	}

	if response != nil {
		response.Body.Close()
	}

	return response, request.Context().Err()
}

func (api *APIClient) SendMetricsData(data *[]byte) error {

	cBody, err := compressBody(data)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, api.baseURL+"/updates/", cBody)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Encoding", "gzip")

	if len(api.key) > 0 {
		h := hmac.New(sha256.New, api.key)
		h.Write(*data)
		request.Header.Set("HashSHA256", string(h.Sum(nil)))
	}

	response, err := api.DoRequest(request)
	if response != nil {
		defer response.Body.Close()
	}

	return err
}
