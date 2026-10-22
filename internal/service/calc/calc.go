package calc

import (
	"errors"
	"math"
)

type Op string

const (
	OpAdd Op = "add"
	OpSub Op = "sub"
	OpMul Op = "mul"
	OpDiv Op = "div"
)

var (
	ErrUnknownOp      = errors.New("unknown operation")
	ErrDivisionByZero = errors.New("division by zero")
	ErrOverflow       = errors.New("result overflow")
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Calculate(op1, op2 float64, op Op) (float64, error) {
	var res float64
	switch op {
	case OpAdd:
		res = op1 + op2
	case OpSub:
		res = op1 - op2
	case OpMul:
		res = op1 * op2
	case OpDiv:
		if op2 == 0 {
			return 0, ErrDivisionByZero
		}
		res = op1 / op2
	default:
		return 0, ErrUnknownOp
	}

	if math.IsInf(res, 0) {
		return 0, ErrOverflow
	}

	return res, nil
}
