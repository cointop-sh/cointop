package eval

import (
	"errors"
	"strings"

	"github.com/Knetic/govaluate"
)

// EvaluateExpression evaulates a simple math expression string to a float64
func EvaluateExpressionToFloat64(input string) (float64, error) {
	input = strings.TrimSpace(input) // remove trailing \0s
	if input == "" {
		return 0, nil
	}
	expression, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return 0, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	f64, ok := result.(float64)
	if !ok {
		return 0, errors.New("could not type assert float64")
	}
	return f64, nil
}
