// Package grpc implements logic of GRPC services.
package grpc

import "google.golang.org/grpc"

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=Server
type Server interface {
	Register(gRPCServer *grpc.Server)
}
