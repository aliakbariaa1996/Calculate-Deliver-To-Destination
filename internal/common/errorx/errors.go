package errorx

import "errors"

var (
	ErrCalculate = errors.New("Calculate error")
)

var code = map[error]string{
	ErrCalculate: "CALCULATE",
}

func CodeError(err error) string {
	return code[err] //todo not implemented
}
