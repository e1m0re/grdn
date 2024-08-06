package grpc

const (
	defaultServerAddr = "localhost:30000"
)

type Config struct {
	ServerAddr string
}

func InitConfig() (*Config, error) {
	return &Config{
		ServerAddr: defaultServerAddr,
	}, nil
}
