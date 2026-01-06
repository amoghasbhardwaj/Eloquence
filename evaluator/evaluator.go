// ==============================================================================================
// FILE: evaluator/evaluator.go
// ==============================================================================================
// PACKAGE: evaluator
// PURPOSE: Implements the runtime execution engine.
//          It traverses the AST and produces side effects (IO) or results (Objects).
//          It handles variable scoping, control flow, and error propagation.
// ==============================================================================================

package evaluator

import (
	"fmt"

	"eloquence/ast"
	"eloquence/object"
)

// Singletons for performance (avoid allocating new true/false/null objects constantly)
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval is the heart of the interpreter. It recursively evaluates AST nodes.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// --- Root ---
	case *ast.Program:
		return evalProgram(node, env)

	// --- Statements ---
	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LoopStatement:
		return evalLoopStatement(node, env)

	case *ast.PointerAssignmentStatement:
		return evalPointerAssignment(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ShowStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if val != NULL {
			fmt.Println(val.Inspect())
		}
		return NULL

	case *ast.StructDefinitionStatement:
		return evalStructDefinition(node, env)

	case *ast.TryCatchStatement:
		return evalTryCatchStatement(node, env)

	// --- Expressions ---
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.FieldAccessExpression:
		return evalFieldAccess(node, env)

	case *ast.FunctionLiteral:
		// Capture the current environment for closure support
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

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

	case *ast.PointerDereferenceExpression:
		return evalPointerDereference(node, env)

	case *ast.PointerReferenceExpression:
		return evalPointerReference(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.StructInstantiationExpression:
		return evalStructInstantiation(node, env)

	// --- Literals ---
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.BooleanLiteral:
		return nativeBool(node.Value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.MapLiteral:
		return evalMapLiteral(node, env)

	case *ast.NilLiteral:
		return NULL

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return NULL
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	// result defaults to nil interface, which we treat as NULL in sanity tests
	// but strictly speaking, we want to return the last evaluated object.

	for _, s := range p.Statements {
		result = Eval(s, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(b *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	result = NULL // Default for empty blocks

	for _, s := range b.Statements {
		result = Eval(s, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalLoopStatement(node *ast.LoopStatement, env *object.Environment) object.Object {
	// CRITICAL FIX: Loops share the parent environment scope.
	// This allows the loop body to modify variables (like counters) defined outside.

	for {
		cond := Eval(node.Condition, env)
		if isError(cond) {
			return cond
		}
		if !isTruthy(cond) {
			break
		}

		result := Eval(node.Body, env)
		if result != nil {
			// Check for interrupts (Return/Error)
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return NULL
}

func evalTryCatchStatement(node *ast.TryCatchStatement, env *object.Environment) object.Object {
	// Try block gets its own scope to contain any local definitions
	tryEnv := object.NewEnclosedEnvironment(env)
	result := evalBlockStatement(node.TryBlock, tryEnv)

	if isError(result) {
		if node.CatchBlock != nil {
			catchEnv := object.NewEnclosedEnvironment(env)
			// Future: Bind the error object to a variable here
			return evalBlockStatement(node.CatchBlock, catchEnv)
		}
		return NULL
	}

	if node.FinallyBlock != nil {
		evalBlockStatement(node.FinallyBlock, object.NewEnclosedEnvironment(env))
	}

	return result
}

func evalStructDefinition(node *ast.StructDefinitionStatement, env *object.Environment) object.Object {
	def := &object.StructDefinition{
		Name:   node.Name.Value,
		Fields: []string{},
	}
	for _, f := range node.Attributes {
		def.Fields = append(def.Fields, f.Value)
	}
	env.Set(node.Name.Value, def)
	return NULL
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}

	// Conditionals use a NEW scope to prevent variable leakage (Shadowing)
	scopedEnv := object.NewEnclosedEnvironment(env)

	if isTruthy(cond) {
		return Eval(ie.Consequence, scopedEnv)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, scopedEnv)
	}
	return NULL
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return nativeBool(!isTruthy(right))
	case "-", "minus":
		return evalMinusPrefix(right)
	}
	return newError("unknown operator: %s%s", op, right.Type())
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	// Handle NULL comparisons gracefully (e.g., node.next equals none)
	if left.Type() != right.Type() {
		if left.Type() == object.NULL_OBJ || right.Type() == object.NULL_OBJ {
			if op == "equals" {
				return FALSE
			}
			if op == "not_equals" {
				return TRUE
			}
		}
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	}

	switch left.Type() {
	case object.INTEGER_OBJ:
		return evalIntegerInfix(op, left.(*object.Integer), right.(*object.Integer))
	case object.FLOAT_OBJ:
		return evalFloatInfix(op, left.(*object.Float), right.(*object.Float))
	case object.STRING_OBJ:
		return evalStringInfix(op, left.(*object.String), right.(*object.String))
	case object.BOOLEAN_OBJ:
		return evalBooleanInfix(op, left.(*object.Boolean), right.(*object.Boolean))
	case object.NULL_OBJ:
		if op == "equals" {
			return TRUE
		}
		if op == "not_equals" {
			return FALSE
		}
	}
	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	f, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}
	// Create a new scope for the function execution, extending the closure's captured env
	env := object.NewEnclosedEnvironment(f.Env)
	for i, param := range f.Parameters {
		if i < len(args) {
			env.Set(param.Value, args[i])
		}
	}
	evaluated := Eval(f.Body, env)
	if rv, ok := evaluated.(*object.ReturnValue); ok {
		return rv.Value
	}
	return evaluated
}

func evalPointerReference(node *ast.PointerReferenceExpression, env *object.Environment) object.Object {
	ident, ok := node.Value.(*ast.Identifier)
	if !ok {
		return newError("can only point to identifier")
	}

	// Resolve exact environment where the variable lives to allow mutation
	targetEnv := env.Resolve(ident.Value)
	if targetEnv == nil {
		return newError("identifier not found: %s", ident.Value)
	}

	return &object.Pointer{Name: ident.Value, Env: targetEnv}
}

func evalPointerDereference(node *ast.PointerDereferenceExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	ptr, ok := val.(*object.Pointer)
	if !ok {
		return newError("cannot dereference non-pointer")
	}
	targetVal, ok := ptr.Env.Get(ptr.Name)
	if !ok {
		return newError("dangling pointer: %s", ptr.Name)
	}
	return targetVal
}

func evalPointerAssignment(node *ast.PointerAssignmentStatement, env *object.Environment) object.Object {
	ptrObj, ok := env.Get(node.Name.Value)
	if !ok {
		return newError("identifier not found: %s", node.Name.Value)
	}
	p, ok := ptrObj.(*object.Pointer)
	if !ok {
		return newError("'%s' is not a pointer", node.Name.Value)
	}

	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	// Mutate the original variable in its home environment
	p.Env.Set(p.Name, val)
	return val
}

func evalStructInstantiation(node *ast.StructInstantiationExpression, env *object.Environment) object.Object {
	obj, ok := env.Get(node.Name.Value)
	if !ok {
		return newError("unknown struct: %s", node.Name.Value)
	}
	def, ok := obj.(*object.StructDefinition)
	if !ok {
		return newError("%s is not a struct", node.Name.Value)
	}

	fields := make(map[string]object.Object)
	// Initialize defaults
	for _, f := range def.Fields {
		fields[f] = NULL
	}
	// Set provided values
	for _, f := range node.Fields {
		val := Eval(f.Value, env)
		if isError(val) {
			return val
		}
		fields[f.Name.Value] = val
	}
	return &object.StructInstance{Definition: def, Fields: fields}
}

func evalFieldAccess(node *ast.FieldAccessExpression, env *object.Environment) object.Object {
	left := Eval(node.Object, env)
	if isError(left) {
		return left
	}
	strct, ok := left.(*object.StructInstance)
	if !ok {
		return newError("not a struct instance: %s", left.Type())
	}
	val, ok := strct.Fields[node.Field.Value]
	if !ok {
		return newError("struct %s has no field %s", strct.Definition.Name, node.Field.Value)
	}
	return val
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	return newError("identifier not found: %s", node.Value)
}

func evalMinusPrefix(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -obj.Value}
	case *object.Float:
		return &object.Float{Value: -obj.Value}
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalIntegerInfix(op string, l, r *object.Integer) object.Object {
	switch op {
	case "adds":
		return &object.Integer{Value: l.Value + r.Value}
	case "subtracts", "minus", "-":
		return &object.Integer{Value: l.Value - r.Value}
	case "times":
		return &object.Integer{Value: l.Value * r.Value}
	case "divides":
		if r.Value == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: l.Value / r.Value}
	case "modulo":
		return &object.Integer{Value: l.Value % r.Value}
	case "equals":
		return nativeBool(l.Value == r.Value)
	case "not_equals":
		return nativeBool(l.Value != r.Value)
	case "greater":
		return nativeBool(l.Value > r.Value)
	case "less":
		return nativeBool(l.Value < r.Value)
	case "greater_equal":
		return nativeBool(l.Value >= r.Value)
	case "less_equal":
		return nativeBool(l.Value <= r.Value)
	}
	return newError("unknown operator: INTEGER %s INTEGER", op)
}

func evalFloatInfix(op string, l, r *object.Float) object.Object {
	switch op {
	case "adds":
		return &object.Float{Value: l.Value + r.Value}
	case "subtracts", "minus", "-":
		return &object.Float{Value: l.Value - r.Value}
	case "times":
		return &object.Float{Value: l.Value * r.Value}
	case "divides":
		return &object.Float{Value: l.Value / r.Value}
	case "equals":
		return nativeBool(l.Value == r.Value)
	case "not_equals":
		return nativeBool(l.Value != r.Value)
	case "greater":
		return nativeBool(l.Value > r.Value)
	case "less":
		return nativeBool(l.Value < r.Value)
	case "greater_equal":
		return nativeBool(l.Value >= r.Value)
	case "less_equal":
		return nativeBool(l.Value <= r.Value)
	}
	return newError("unknown operator: FLOAT %s FLOAT", op)
}

func evalStringInfix(op string, l, r *object.String) object.Object {
	switch op {
	case "adds":
		return &object.String{Value: l.Value + r.Value}
	case "equals":
		return nativeBool(l.Value == r.Value)
	case "not_equals":
		return nativeBool(l.Value != r.Value)
	}
	return newError("unknown operator: STRING %s STRING", op)
}

func evalBooleanInfix(op string, l, r *object.Boolean) object.Object {
	switch op {
	case "equals":
		return nativeBool(l.Value == r.Value)
	case "not_equals":
		return nativeBool(l.Value != r.Value)
	case "and":
		return nativeBool(l.Value && r.Value)
	case "or":
		return nativeBool(l.Value || r.Value)
	}
	return newError("unknown operator: BOOLEAN %s BOOLEAN", op)
}

func evalIndexExpression(left, index object.Object) object.Object {
	if left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ {
		return evalArrayIndex(left.(*object.Array), index.(*object.Integer))
	}
	if left.Type() == object.MAP_OBJ {
		return evalMapIndex(left.(*object.Map), index)
	}
	return newError("index operator not supported: %s", left.Type())
}

func evalArrayIndex(array *object.Array, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return array.Elements[idx]
}

func evalMapIndex(m *object.Map, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as map key: %s", index.Type())
	}
	pair, ok := m.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as map key: %s", key.Type())
		}
		val := Eval(valNode, env)
		if isError(val) {
			return val
		}
		pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: val}
	}
	return &object.Map{Pairs: pairs}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		val := Eval(e, env)
		if isError(val) {
			return []object.Object{val}
		}
		result = append(result, val)
	}
	return result
}

func nativeBool(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
