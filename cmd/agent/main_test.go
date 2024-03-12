package main

import (
	"github.com/e1m0re/grdn/internal/monitor"
	"testing"
)

func TestSendData(t *testing.T) {
	type args struct {
		monitor1 *monitor.MetricsMonitor
	}

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendData(tt.args.monitor1)
		})
	}
}

func Test_doRequest(t *testing.T) {
	type args struct {
		uriPath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := doRequest(tt.args.uriPath); (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
