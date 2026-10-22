package handler

import "github.com/britlydev/c210cc/internal/transport/httpapi/oapi"

var _ oapi.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	*HealthHandler
}
