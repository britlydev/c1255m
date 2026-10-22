package httpapi_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/service/calc"
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

func TestCalc(t *testing.T) {
	srv := newTestServer(t, &handler.Handler{
		CalcHandler: handler.NewCalcHandler(calc.New()),
	})

	tests := []struct {
		name           string
		op1            string
		op2            string
		op             calc.Op
		wantStatusCode int
		wantResult     float64
		wantErr        string
	}{
		{
			name:           "test mul op",
			op1:            "2",
			op2:            "2",
			op:             calc.OpMul,
			wantStatusCode: http.StatusOK,
			wantResult:     4.0,
			wantErr:        "",
		},
		{
			name:           "div by zero error",
			op1:            "2",
			op2:            "0",
			op:             calc.OpDiv,
			wantStatusCode: http.StatusBadRequest,
			wantResult:     0.0,
			wantErr:        calc.ErrDivisionByZero.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf("/api/v1/calc?op1=%s&op=%s&op2=%s", tt.op1, string(tt.op), tt.op2)
			req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, nil)
			resp, err := srv.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatusCode)
			}
			if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
				t.Fatalf("content-type = %q, want application/json", ct)
			}

			if tt.wantErr != "" {
				var respBody oapi.Calculate400JSONResponse
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatal(err)
				}

				if tt.wantErr != respBody.Error {
					t.Fatalf("got err = %s, want err = %s", respBody.Error, tt.wantErr)
				}
			} else {
				var respBody oapi.Calculate200JSONResponse
				if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
					t.Fatal(err)
				}

				if respBody.Result != tt.wantResult {
					t.Fatalf("got = %.2f, want = %.2f", respBody.Result, tt.wantResult)
				}
			}
		})
	}
}
