package httpapi_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/transport/httpapi"
	"github.com/britlydev/c210cc/internal/transport/httpapi/handler"
	"github.com/britlydev/c210cc/internal/transport/httpapi/oapi"
)

func newTestServer(t *testing.T, h *handler.Handler) *httpapi.Server {
	t.Helper()

	return httpapi.NewServer(config.HTTP{}, h, slog.New(slog.DiscardHandler))
}

func TestHealthcheck(t *testing.T) {
	srv := newTestServer(t, &handler.Handler{
		HealthHandler: handler.NewHealthHandler(config.EnvDevelopment),
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/healthcheck", nil)
	resp, err := srv.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q, want application/json", ct)
	}

	var got oapi.HealthcheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	want := oapi.HealthcheckResponse{
		Status:      "available",
		Environment: config.EnvDevelopment,
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
