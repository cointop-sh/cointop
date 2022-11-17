package eval

import (
	"errors"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/ast"
)

// AST Visitor that changes Integer to Float
type patcher struct{}

func (p *patcher) Enter(_ *ast.Node) {}
func (p *patcher) Exit(node *ast.Node) {
	n, ok := (*node).(*ast.IntegerNode)
	if ok {
		ast.Patch(node, &ast.FloatNode{Value: float64(n.Value)})
	}
}

// EvaluateExpressionToFloat64 evaulates a simple math expression string to a float64
func EvaluateExpressionToFloat64(input string, env interface{}) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, nil
	}
	program, err := expr.Compile(input, expr.Env(env), expr.Patch(&patcher{}))
	if err != nil {
		return 0, err
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return 0, err
	}
	f64, ok := result.(float64)
	if !ok {
		ival, ok := result.(int)
		if !ok {
			return 0, errors.New("expression did not return numeric type")
		}
		f64 = float64(ival)
	}
	return f64, nil
}

func EvaluateExpressionToString(input string, env interface{}) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	program, err := expr.Compile(input, expr.Env(env))
	if err != nil {
		return "", err
	}
	result, err := expr.Run(program, env)
	if err != nil {
		return "", err
	}
	s, ok := result.(string)
	if !ok {
		return "", errors.New("expression did not return string type")
	}
	return s, nil
}
