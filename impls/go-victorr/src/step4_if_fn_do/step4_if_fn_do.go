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
	return &mal{
		scanner: bufio.NewScanner(os.Stdin),
		eof:     false,
		env:     newEnv(),
	}
}

func newEnv() MalEnv {
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
	addSpecial("def!", primitiveDef, DefSpecialForm)
	addSpecial("let*", primitiveLet, LetSpecialForm)

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
