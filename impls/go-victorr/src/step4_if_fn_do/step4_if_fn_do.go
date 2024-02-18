package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Mal interface {
	Read() (MalObject, error)
	Eval(MalObject) (MalObject, error)
	Print(MalObject) error
}

var (
	ErrMalEof = errors.New("MAL EOF")
)

type mal struct {
	scanner *bufio.Scanner
	eof     bool
	env     MalEnv
}

var _ Mal = (*mal)(nil)

func newMal() *mal {
	ret := &mal{
		scanner: bufio.NewScanner(os.Stdin),
		eof:     false,
	}
	ret.env = ret.newEnv()

	return ret
}

func (m *mal) newEnv() MalEnv {
	env := NewMalEnv(nil)

	add := func(s string, f ApplyFunc) {
		env.Set(NewMalSymbol(s), NewMalFunction(f, NotSpecialForm))
	}
	addSpecial := func(s string, f ApplyFunc, specialForm SpecialForm) {
		env.Set(NewMalSymbol(s), NewMalFunction(f, specialForm))
	}

	add("+", primitiveAdd)
	add("-", primitiveSubtract)
	add("*", primitiveMultiply)
	add("/", primitiveDivide)
	add("<", primitiveNumberComparison("<", func(l MalNumber, r MalNumber) bool {
		return l.Number() < r.Number()
	}))
	add("<=", primitiveNumberComparison("<=", func(l MalNumber, r MalNumber) bool {
		return l.Number() <= r.Number()
	}))
	add(">", primitiveNumberComparison(">", func(l MalNumber, r MalNumber) bool {
		return l.Number() > r.Number()
	}))
	add(">=", primitiveNumberComparison(">=", func(l MalNumber, r MalNumber) bool {
		return l.Number() >= r.Number()
	}))
	add("list", primitiveList)
	add("list?", primitivePredicate(func(o MalObject) (bool, error) {
		_, ok := o.(MalList)
		return ok, nil
	}))
	add("empty?", primitivePredicate(func(o MalObject) (bool, error) {
		list, isList := o.(MalList)
		vector, isVector := o.(MalVector)
		if !(isList || isVector) {
			return false, fmt.Errorf("expected a MalList or a MalVector but got %T", o)
		}
		if isList {
			return len(list.List()) == 0, nil
		}
		return len(vector.Vector()) == 0, nil
	}))
	add("count", primitiveCount)
	add("=", primitiveEquals)
	add("prn", primitivePrn)
	add("println", primitivePrintln)
	add("pr-str", primitivePrStr)
	add("str", primitiveStr)

	addSpecial("def!", primitiveDef, DefSpecialForm)
	addSpecial("let*", primitiveLet, LetSpecialForm)
	addSpecial("if", m.ifSpecialForm, IfSpecialForm)
	addSpecial("do", m.doSpecialForm, IfSpecialForm)
	addSpecial("fn*", m.fnSpecialForm, IfSpecialForm)

	return env
}

func (m *mal) Read() (MalObject, error) {
	if m.eof {
		return nil, ErrMalEof
	}
	m.eof = !m.scanner.Scan()
	return ReadStr(m.scanner.Text())
}

func (m *mal) Eval(object MalObject) (MalObject, error) {
	// fmt.Printf("Eval(%T %s)\n", object, object)
	switch object.(type) {
	case MalList:
		// fmt.Printf("Eval() a list: %s\n", object)
		if len(object.(MalList).List()) == 0 {
			return object, nil
		}
		f, args, err := m.evalAst(object)
		if err != nil {
			return nil, err
		}
		return m.Apply(m.env, f, args)

	default:
		// fmt.Printf("Eval() not a list: %T %s\n", object, object)
		evaluated, _, err := m.evalAst(object)
		if err != nil {
			return nil, err
		}

		return evaluated, nil
	}
}

func (m *mal) evalAst(ast MalObject) (MalObject, []MalObject, error) {
	switch ast.(type) {

	case MalSymbol:
		//fmt.Printf("evalAst symbol %s, ret %s", ast, m.env[ast])
		symbol := ast.(MalSymbol)
		if strings.HasPrefix(symbol.Symbol(), ":") {
			return ast, nil, nil
		}
		value, err := m.env.Get(symbol)
		if err != nil {
			return nil, nil, err
		}
		return value, nil, nil

	case MalList:
		// Special forms:
		// (def! symbol value) -- (setq symbol value)
		// (let* (s1 v1 s2 v2) body)

		list := ast.(MalList)
		if len(list.List()) == 0 {
			return ast, nil, nil
		}
		forms := list.List()

		funObject, err := m.Eval(forms[0])
		if err != nil {
			return nil, nil, err
		}
		fun, isFunction := funObject.(MalFunction)
		if !isFunction {
			return nil, nil, fmt.Errorf("expected a function but got %T", funObject)
		}

		formIndex := 1
		var args []MalObject

		switch fun.SpecialForm() {
		case NotSpecialForm:
			// nothing to do

		case DefSpecialForm:
			args = append(args, forms[formIndex])
			formIndex++

		case LetSpecialForm:
			// Need to evaluate form at index 1 to make a new environment
			// Need to evaluate rest of the forms with the new environment
			args = append(args, forms[formIndex:]...)
			formIndex = len(forms)

			fun = NewMalFunction(m.letSpecialForm, NotSpecialForm).(MalFunction)

		case IfSpecialForm:
			args = append(args, forms[formIndex:]...)
			formIndex = len(forms)

		default:
			return nil, nil, fmt.Errorf("unsupported special form: %v", fun.SpecialForm())
		}

		for index := formIndex; index < len(forms); index++ {
			object, err := m.Eval(forms[index])
			if err != nil {
				return nil, nil, err
			}
			args = append(args, object)
		}

		return fun, args, nil

	case MalVector:
		vector := ast.(MalVector)
		var evaluated []MalObject
		for _, o := range vector.Vector() {
			// fmt.Printf("evalAst vector o=%T %s\n", o, o)
			e, err := m.Eval(o)
			if err != nil {
				return nil, nil, err
			}
			evaluated = append(evaluated, e)
		}
		return NewMalVector(evaluated), nil, nil

	case MalHashMap:
		hashmap := ast.(MalHashMap)
		var evaluated []MalObject
		for _, o := range hashmap.HashMap() {
			e, err := m.Eval(o)
			if err != nil {
				return nil, nil, err
			}
			// fmt.Printf("evalAst hashmap object %T %s\n", e, e)
			evaluated = append(evaluated, e)
		}
		return NewMalHashMap(evaluated), nil, nil

	default:
		return ast, nil, nil
	}
}

func (m *mal) ifSpecialForm(_ MalEnv, args []MalObject) (MalObject, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of elements for if: %d", len(args))
	}
	cond := args[0]
	ifTrue := args[1]
	ifFalse := MalObject(MalNil)
	if len(args) >= 3 {
		ifFalse = args[2]
	}

	result, err := m.Eval(cond)
	if err != nil {
		return nil, err
	}
	if IsTrue(result) {
		return m.Eval(ifTrue)
	}
	return m.Eval(ifFalse)
}

func (m *mal) doSpecialForm(env MalEnv, args []MalObject) (MalObject, error) {
	var ret MalObject
	var err error
	for _, arg := range args {
		ret, err = m.Eval(arg)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// fnSpecialForm, fn*: Return a new function closure. The body of that closure does the following:
//
// Create a new environment using env (closed over from outer scope) as the outer parameter, the
// first parameter (second list element of ast from the outer scope) as the binds parameter, and the
// parameters to the closure as the exprs parameter.
//
// Call EVAL on the second parameter (third list element of ast from outer scope), using the new
// environment. Use the result as the return value of the closure.
func (m *mal) fnSpecialForm(env MalEnv, args []MalObject) (MalObject, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments for fn*: %d", len(args))
	}
	var binds []MalObject
	if bindList, ok := args[0].(MalList); ok {
		binds = bindList.List()
	} else if bindVector, ok := args[0].(MalVector); ok {
		binds = bindVector.Vector()
	} else {
		return nil, errWrongType("MalList or MalVector", args[0])
	}
	for _, bind := range binds {
		if _, isSymbol := bind.(MalSymbol); !isSymbol {
			return nil, errWrongType("MalSymbol", bind)
		}
	}
	body := args[1:]

	closure := m.env
	fn := func(env MalEnv, args []MalObject) (MalObject, error) {
		return m.fnClosure(closure, env, binds, body, args)
	}
	//fmt.Printf("closure: %s\n", closure)
	return NewMalFunction(fn, NotSpecialForm), nil
}

func (m *mal) fnClosure(closure, env MalEnv, binds, body, args []MalObject) (MalObject, error) {
	// fmt.Printf("closure: %s\n", closure)
	// fmt.Printf("env: %s\n", env)

	// Build the apply environment
	applyEnv := NewMalEnv(closure)
	bindRest := false
	for index := 0; index < len(binds); index++ {
		bind := binds[index].(MalSymbol) // caller verified it is a symbol
		switch {
		case bind.Symbol() == "&":
			bindRest = true
			continue

		case bindRest:
			var rest []MalObject
			if index-1 < len(args) {
				rest = args[index-1:]
			}
			applyEnv.Set(bind, NewMalList(rest))
			break

		case index < len(args):
			applyEnv.Set(bind, args[index])

		default:
			applyEnv.Set(bind, MalNil)
		}
	}

	m.env = applyEnv
	defer func() {
		m.env = env
	}()
	var ret MalObject
	var err error
	for _, form := range body {
		ret, err = m.Eval(form)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (m *mal) letSpecialForm(env MalEnv, args []MalObject) (MalObject, error) {
	m.env = NewMalEnv(env)
	defer func() {
		m.env = env
	}()

	err := m.makeLetEnv(args[0])
	if err != nil {
		return nil, err
	}

	var ret MalObject
	for index := 1; index < len(args); index++ {
		var err error
		ret, err = m.Eval(args[index])
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (m *mal) makeLetEnv(form MalObject) error {
	list, isList := form.(MalList)
	vector, isVector := form.(MalVector)
	if !(isList || isVector) {
		return fmt.Errorf("expected a list of bindings for let*, but got %T", form)
	}
	var objects []MalObject
	if isList {
		objects = list.List()
	}
	if isVector {
		objects = vector.Vector()
	}

	if len(objects)%2 == 1 {
		return fmt.Errorf("expected a list of bindings for let*, but got %d forms", len(objects))
	}

	for index := 0; index < len(objects); index += 2 {
		key := objects[index]
		value := objects[index+1]

		symbol, ok := key.(MalSymbol)
		if !ok {
			return fmt.Errorf("expected a MalSymbol in let* binding, bug got %T", key)
		}
		evaluated, err := m.Eval(value)
		if err != nil {
			return err
		}
		m.env.Set(symbol, evaluated)
	}
	return nil
}

func (m *mal) Apply(env MalEnv, f MalObject, args []MalObject) (MalObject, error) {
	fun, ok := f.(MalFunction)
	if !ok {
		return nil, fmt.Errorf("expected a MalFunction but got a %T", f)
	}

	return fun.Function()(env, args)
}

func (m *mal) Print(object MalObject) error {
	s, err := PrintString(object)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

func rep() error {
	m := newMal()

	{
		not, err := ReadStr("(def! not (fn* (a) (if a false true)))")
		if err == nil {
			_, err = m.Eval(not)
		}
		if err != nil {
			return err
		}
	}
	for {
		fmt.Print("user> ")
		in, err := m.Read()
		// if err == ErrMalEof {
		// 	fmt.Printf("error: %s\n", err)
		// } else
		var result MalObject
		if err == nil {
			result, err = m.Eval(in)
		}
		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			m.Print(result)
		}
	}
}

func main() {
	err := rep()
	if err != nil {
		panic(err.Error())
	}
}
