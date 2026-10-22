package handler

import (
	"context"
	"errors"

	"github.com/britlydev/c210cc/internal/service/calc"
	"github.com/britlydev/c210cc/internal/transport/httpapi/oapi"
)

type Calculator interface {
	Calculate(op1, op2 float64, op calc.Op) (float64, error)
}

type CalcHandler struct {
	calc Calculator
}

func NewCalcHandler(c Calculator) *CalcHandler {
	return &CalcHandler{calc: c}
}

func (h *CalcHandler) Calculate(ctx context.Context, request oapi.CalculateRequestObject) (oapi.CalculateResponseObject, error) {
	res, err := h.calc.Calculate(request.Params.Op1, request.Params.Op2, calc.Op(request.Params.Op))
	if err != nil {
		if errors.Is(err, calc.ErrDivisionByZero) || errors.Is(err, calc.ErrOverflow) || errors.Is(err, calc.ErrUnknownOp) {
			return oapi.Calculate400JSONResponse{Error: err.Error()}, nil
		}
		return nil, err
	}

	return oapi.Calculate200JSONResponse{Result: res}, nil
}
