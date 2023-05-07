package auth

import (
	"testing"

	"github.com/fi3te/monthly-planning-backend/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthorizedTrue(t *testing.T) {
	input := map[string]string{"authorization": "Basic YWRtaW46cGFzc3dvcmQ="}
	cfg := config.Config{Username: "admin", Password: "password"}

	result := IsAuthorized(&cfg, input)

	assert.True(t, result)
}

func TestIsAuthorizedFalse(t *testing.T) {
	input := map[string]string{"authorization": "Basic YWRtaW46cGFzc3dvcmQ="}
	cfg := config.Config{Username: "admin", Password: "other"}

	result := IsAuthorized(&cfg, input)

	assert.False(t, result)
}

func TestIsAuthorizedFalseInvalidHeader(t *testing.T) {
	input := map[string]string{"authorization": "Basic invalid"}
	cfg := config.Config{Username: "admin", Password: "password"}

	result := IsAuthorized(&cfg, input)

	assert.False(t, result)
}
