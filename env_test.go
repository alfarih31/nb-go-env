package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	nv, err := LoadEnv(".env.test")
	assert.Equal(t, err, nil, "Error must be nil")

	_, ok := nv.(Env)
	assert.Equal(t, ok, true, "nv must be instance of Env")
}

func TestLoadEnvFallbackToWide(t *testing.T) {
	nv, err := LoadEnv("", true)
	assert.Equal(t, err, nil, "Error must be nil")

	_, ok := nv.(Env)
	assert.Equal(t, ok, true, "nv must be instance of Env")
}
