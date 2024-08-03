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
				os.Setenv(envRunAddrName, "127.0.0.1:8081")
				os.Setenv(envKeyName, "key")
				os.Setenv(envCryptoKeyName, "public key")
				os.Setenv(envStoreIntervalName, "100s")
				os.Setenv(envRestoreDataName, "true")
				os.Setenv(envDatabaseDSNName, "")
				os.Setenv(envFileStoragePathName, "/tmp/tmp.tmp")
				os.Setenv(envTrustedSubnet, "172.20.1.0/24")
			},
			want: want{
				cfg: &Config{
					Key:             "key",
					ServerAddr:      "127.0.0.1:8081",
					StoreInternal:   time.Duration(100) * time.Second,
					FileStoragePath: "/tmp/tmp.tmp",
					DatabaseDSN:     "",
					PrivateKeyFile:  "public key",
					RestoreData:     true,
					LoggerLevel:     "info",
					TrustedSubnet:   "172.20.1.0/24",
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
