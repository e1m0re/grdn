package main

import (
	"testing"

	"github.com/e1m0re/grdn/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestSendData(t *testing.T) {
	type args struct {
		data *storage.MetricsState
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendData(tt.args.data)
		})
	}
}

func TestUpdateMetrics(t *testing.T) {
	type args struct {
		data *storage.MetricsState
	}
	type want struct {
		gValue storage.GuageDateType
		cValue storage.CounterDateType
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test update Counter PollCount",
			args: args{data: storage.NewMetricsState()},
			want: want{
				gValue: 0,
				cValue: 1,
			},
		},
	}
	for _, tt := range tests {
		UpdateMetrics(tt.args.data)
		assert.Equal(t, tt.want.cValue, tt.args.data.Counters[storage.PollCount])
	}
}

func Test_sendMetric(t *testing.T) {
	type args struct {
		mType  string
		mName  string
		mValue string
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
			if err := sendMetric(tt.args.mType, tt.args.mName, tt.args.mValue); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
