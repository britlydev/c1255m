package handler

import (
	"context"

	"github.com/britlydev/c210cc/internal/transport/httpapi/oapi"
)

type HealthHandler struct {
	env string
}

func NewHealthHandler(env string) *HealthHandler {
	return &HealthHandler{env: env}
}

func (h *HealthHandler) Healthcheck(ctx context.Context, request oapi.HealthcheckRequestObject) (oapi.HealthcheckResponseObject, error) {
	return oapi.Healthcheck200JSONResponse{
		Status:      "available",
		Environment: h.env,
	}, nil
}
