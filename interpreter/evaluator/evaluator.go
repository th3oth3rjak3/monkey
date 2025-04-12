package evaluator

import (
	"fmt"
	"monkey/interpreter/ast"
	"monkey/interpreter/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval is the evaluator function that handles conversion from ast nodes to objects.
//
// Parameters:
//   - node: The input ast Node.
//   - env: The environment which contains the current state.
//
// Returns:
//   - object.Object: The evaluated object.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
}

// evalProgram evaluates the root program node.
//
// Parameters:
//   - stmts: A slice of statements to be evaluated.
//   - env: The root environment.
//
// Returns:
//   - object.Object: The evaluated object.
func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

// evalBlockStatement evaluates a block of code.
//
// Parameters:
//   - block: The block of code to evaluate.
//   - env: The environment.
//
// Returns:
//   - object.Object: The result of evaluating the body.
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
			return result
		}
	}

	return result
}

// nativeBoolToBooleanObject is a helper function that converts a go bool type to a monkey Boolean.
//
// Parameters:
//   - input: The go boolean input.
//
// Returns:
//   - *object.Boolean: The monkey boolean representation.
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

// evalPrefixExpression evaluates expressions that have a prefix operator.
//
// Parameters:
//   - operator: The prefix operator e.g. ! or -
//   - right: The object to evaluate in the prefix expression.
//
// Returns:
//   - object.Object: The result of processing the prefix operator and the expression.
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression evaluates what should happen for a given input
// object when the prefix operator is a !
//
// Parameters:
//   - right: The right hand side expression to be evaluated against the boolean negation operator.
//
// Returns:
//   - object.Object: The result of the evaluation.
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression handles the integer negation prefix operation.
//
// Parameters:
//   - right: The object to have the integer negation operator applied.
//
// Returns:
//   - object.Object: The result of evaluating the integer negation operation.
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalInfixExpression handles the evaluation for all infix expressions
//
// Parameters:
//   - operator: The infix operator.
//   - left: The left hand side expression.
//   - right: The right hand side expression.
//
// Returns:
//   - object.Object: The result of evaluating the infix expression.
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression evaluates all infix expressions for the integer type.
//
// Parameters:
//   - operator: The operator for the expression.
//   - left: The left hand side expression.
//   - right: The right hand side expression.
//
// Returns:
//   - object.Object: The result of evaluating the infix expression.
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

// evalIfExpression evaluates the results of an if/else expression.
//
// Parameters:
//   - i: The if expression to be evaluated.
//   - env: The environment with the current state.
//
// Returns:
//   - object.Object: The result of evaluating the if/else expression.
func evalIfExpression(i *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(i.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(i.Consequence, env)
	} else if i.Alternative != nil {
		return Eval(i.Alternative, env)
	} else {
		return NULL
	}
}

// isTruthy evaluates the input object and determines if is a "truthy" value.
//
// Parameters:
//   - obj: The input object.
//
// Returns:
//   - bool: False when the value is false or null, otherwise true.
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// newError creates a new Error object with the provided message.
//
// Parameters:
//   - format: The format string used to create the error message.
//   - a: Arguments to the format string.
//
// Returns:
//   - *object.Error: The error object with a formatted error message.
func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError checks to see if the input object is an error type.
//
// Parameters:
//   - obj: The input object which may be an error.
//
// Returns:
//   - bool: True when the input is an error type, otherwise false.
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

// evalIdentifier handles evaluation of an identifier from the environment.
//
// Parameters:
//   - node: The identifier node that contains the name of the identifier.
//   - env: The environment used to search for the identifier's value.
//
// Returns:
//   - object.Object: The found value or an error.
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

// evalExpressions evaluates a slice of call expressions used for calling a function.
//
// Parameters:
//   - exps: The expressions to evaluate and turn into objects.
//   - env: The environment with the current state.
//
// Returns:
//   - []object.Object: The collection of expressions evaluated and turned into objects.
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// applyFunction calls the function with the supplied args.
//
// Parameters:
//   - fn: The function to be called.
//   - args: The arguments to pass to the function call.
//
// Returns:
//   - object.Object: The result of calling the function.f
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv encloses the functions environment with a new one that has the arguments defined.
//
// Parameters:
//   - fn: The function to call.
//   - args: The args to pass to the function. These are set in the wrapping environment.
//
// Returns:
//   - *object.Environment: The extended environment.
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}

	return env
}

// unwrapReturnValue unwraps the return value when it exists
// in order to stop the return in the current scope and not
// bubble up to outer functions and blocks.
//
// Parameters:
//   - obj: The object which may contain a return value.
//
// Returns:
//   - object.Object: A return value or the input object.
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
