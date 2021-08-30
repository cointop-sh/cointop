package eval

import (
	"errors"
	"strings"

	"github.com/antonmedv/expr"
)

// EvaluateExpression evaulates a simple math expression string to a float64
func EvaluateExpressionToFloat64(input string) (float64, error) {
	input = strings.TrimSpace(input) // remove trailing \0s
	if input == "" {
		return 0, nil
	}
	result, err := expr.Eval(input, nil)
	if err != nil {
		return 0, err
	}
	f64, ok := result.(float64)
	if !ok {
		ival, ok := result.(int)
		if !ok {
			return 0, errors.New("could not type assert float64")
		}
		f64 = float64(ival)
	}
	return f64, nil
}
