package server

import (
	"testing"

	"github.com/e1m0re/grdn/internal/server/config"
	"github.com/e1m0re/grdn/internal/storage/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := config.Config{}
	s := mocks.NewStore(t)

	srv := NewServer(&cfg, s)
	assert.Implements(t, (*Server)(nil), srv)
}
