package http

import (
	"testing"

	"github.com/e1m0re/grdn/internal/listeners"
	"github.com/e1m0re/grdn/internal/listeners/http/config"
	"github.com/e1m0re/grdn/internal/storage/store/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := config.Config{}
	s := mocks.NewStore(t)

	httpListener := NewHTTPListener(&cfg, s)
	assert.Implements(t, (*listeners.Listener)(nil), httpListener)
}
