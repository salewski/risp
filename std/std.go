package std

import (
	"fmt"
	"github.com/raoulvdberge/risp/parser"
	"github.com/raoulvdberge/risp/runtime"
	"math/big"
)

var Symbols = runtime.Symtab{
	"t":       runtime.True,
	"f":       runtime.False,
	"nil":     runtime.Nil,
	"print":   runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrint, "print")),
	"println": runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdPrintln, "println")),
	"list":    runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdList, "list")),
	"+":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "+")),
	"-":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "-")),
	"*":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "*")),
	"/":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMath, "/")),
	"=":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdEquals, "=")),
	">":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">")),
	">=":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, ">=")),
	"<":       runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<")),
	"<=":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdMathCmp, "<=")),
	"cat":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdCat, "cat")),
	"and":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdAnd, "and")),
	"or":      runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdOr, "or")),
	"not":     runtime.NewFunctionValue(runtime.NewBuiltinFunction(stdNot, "not")),
}

var Macros = runtime.Mactab{
	"defun": runtime.NewMacro(stdDefun, "identifier", "list", "list"),
}

func stdDefun(macro *runtime.Macro, block *runtime.Block, nodes []parser.Node) (*runtime.Value, error) {
	name := nodes[0].(*parser.IdentifierNode).Token.Data
	callback := nodes[2].(*parser.ListNode)
	argNodes := nodes[1].(*parser.ListNode)
	var args []string

	for _, argNode := range argNodes.Nodes {
		ident, ok := argNode.(*parser.IdentifierNode)

		if !ok {
			return nil, runtime.NewRuntimeError(argNode.Pos(), "expected an identifier")
		}

		args = append(args, ident.Token.Data)
	}

	functionBlock := runtime.NewBlock([]parser.Node{callback}, runtime.NewScope(block.Scope))
	function := runtime.NewDeclaredFunction(functionBlock, name, args)

	block.Scope.SetSymbol(name, runtime.NewFunctionValue(function))

	return runtime.Nil, nil
}

func stdPrint(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Print(arg)
	}

	return runtime.Nil, nil
}

func stdPrintln(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	for _, arg := range context.Args {
		fmt.Println(arg)
	}

	return runtime.Nil, nil
}

func stdList(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	l := runtime.NewListValue()

	for _, arg := range context.Args {
		l.List = append(l.List, arg)
	}

	return l, nil
}

func stdMath(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	base := big.NewRat(0, 1)

	var callback func(*big.Rat, *big.Rat) *big.Rat

	switch context.Name {
	case "+":
		callback = base.Add
	case "-":
		callback = base.Sub
	case "*":
		callback = base.Mul
	case "/":
		if context.Args[1].Number.Cmp(base) == 0 {
			return nil, runtime.NewRuntimeError(context.Pos, "division by zero")
		}

		callback = base.Quo
	}

	return runtime.NewNumberValueFromRat(callback(context.Args[0].Number, context.Args[1].Number)), nil
}

func stdMathCmp(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.NumberValue, runtime.NumberValue)

	if err != nil {
		return nil, err
	}

	n1 := context.Args[0].Number
	n2 := context.Args[1].Number

	ok := false

	switch context.Name {
	case ">":
		ok = n1.Cmp(n2) == 1
	case ">=":
		ok = n1.Cmp(n2) >= 0
	case "<":
		ok = n1.Cmp(n2) == -1
	case "<=":
		ok = n1.Cmp(n2) <= 0
	}

	if ok {
		return runtime.True, nil
	} else {
		return runtime.False, nil
	}
}

func stdEquals(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.AnyValue, runtime.AnyValue)

	if err != nil {
		return nil, err
	}

	if context.Args[0].Equals(context.Args[1]) {
		return runtime.True, nil
	} else {
		return runtime.False, nil
	}
}

func stdCat(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	s := ""

	for _, arg := range context.Args {
		s += arg.String()
	}

	return runtime.NewStringValue(s), nil
}

func stdAnd(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	if context.Args[0].Boolean && context.Args[1].Boolean {
		return runtime.True, nil
	} else {
		return runtime.False, nil
	}
}

func stdOr(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	if context.Args[0].Boolean || context.Args[1].Boolean {
		return runtime.True, nil
	} else {
		return runtime.False, nil
	}
}

func stdNot(context *runtime.FunctionCallContext) (*runtime.Value, error) {
	err := runtime.ValidateArguments(context, runtime.BooleanValue)

	if err != nil {
		return nil, err
	}

	if context.Args[0].Boolean == true {
		return runtime.False, nil
	} else {
		return runtime.True, nil
	}
}
