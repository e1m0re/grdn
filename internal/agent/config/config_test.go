package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	type want struct {
		cfg *Config
		err error
	}
	tests := []struct {
		mock func()
		want want
		name string
	}{
		{
			name: "successfully case (values from env)",
			mock: func() {
				os.Setenv(envServerAddrName, "127.0.0.1:8081")
				os.Setenv(envKeyName, "key")
				os.Setenv(envCryptoKeyName, "public key")
				os.Setenv(envPollInterval, "100")
				os.Setenv(envReportIntervalName, "100")
				os.Setenv(envRateLimit, "100")
			},
			want: want{
				cfg: &Config{
					Key:            "key",
					PublicKeyFile:  "public key",
					ServerAddr:     "127.0.0.1:8081",
					PollInterval:   time.Duration(100) * time.Second,
					ReportInterval: time.Duration(100) * time.Second,
					RateLimit:      100,
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			got, err := InitConfig()
			require.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.cfg, got)
		})
	}
}
