package listeners

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Listener
type Listener interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
