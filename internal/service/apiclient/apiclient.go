// Package apiclient implements methods of http requests execution.
package apiclient

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log/slog"
	"net"
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

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=APIClient
type APIClient interface {
	// DoRequest executes HTTP request.
	DoRequest(request *http.Request) (*http.Response, error)
	// SendMetricsData sends metrics data to server.
	SendMetricsData(data *[]byte) error
}

type client struct {
	client  *http.Client
	baseURL string
	key     []byte
}

// NewAPIClient is client constructor.
func NewAPIClient(baseURL string, key []byte) APIClient {
	return &client{
		client:  &http.Client{},
		baseURL: baseURL,
		key:     key,
	}
}

// DoRequest executes HTTP request.
func (api *client) DoRequest(request *http.Request) (*http.Response, error) {
	ip, err := api.GetLocalIP()
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Real-IP", ip.String())
	response, err := api.client.Do(request)
	if err != nil {
		slog.Warn("error while doing request",
			slog.String("url", request.URL.String()),
			slog.String("method", request.Method),
			slog.String("error", err.Error()),
		)
	}

	// todo внимательно посмотреть на обработку ответов с кодом > 400
	if err == nil && response.StatusCode < 500 {
		return response, err
	}

	if response != nil {
		response.Body.Close()
	}

	return response, request.Context().Err()
}

// SendMetricsData sends metrics data to server.
func (api *client) SendMetricsData(data *[]byte) error {
	rawData := *data
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
		h.Write(rawData)
		sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
		request.Header.Set("HashSHA256", sum)
	}

	response, err := api.DoRequest(request)
	if response != nil {
		defer response.Body.Close()
	}

	return err
}

func (api *client) GetLocalIP() (net.IP, error) {
	conn, err := net.Dial("udp", api.baseURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP, nil
}
