package eval

import (
	"errors"
	"strings"

	"github.com/Knetic/govaluate"
)

// EvaluateExpression ...
func EvaluateExpressionToFloat64(expressionInput string) (float64, error) {
	expressionInput = strings.TrimSpace(expressionInput) // remove trailing \0s
	if expressionInput == "" {
		return 0, nil
	}
	expression, err := govaluate.NewEvaluableExpression(expressionInput)
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
