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

// EvaluateExpression evaulates a simple math expression string to a float64
func EvaluateExpressionToFloat64(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, nil
	}
	program, err := expr.Compile(input, expr.Env(nil), expr.Patch(&patcher{}))
	if err != nil {
		return 0, err
	}
	result, err := expr.Run(program, nil)
	if err != nil {
		return 0, err
	}
	f64, ok := result.(float64)
	if !ok {
		return 0, errors.New("could not type assert float64")
	}
	return f64, nil
}
