package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/SonDo580/bookmark-manager/internal/config"
	"github.com/stretchr/testify/require"
)

func getDefaultConfig() *config.Config {
	return &config.Config{
		ServiceName: "bookmark_service",
		InstanceID:  uuid.New().String(),
		Port:        8080,
	}
}

func newJSONRequest(t *testing.T, method string, target string, payload any) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if payload != nil {
		err := json.NewEncoder(&buf).Encode(payload)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, target, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}
