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
	env     map[string]MalObject
}

var _ Mal = (*mal)(nil)

func newMal() *mal {
	return &mal{
		scanner: bufio.NewScanner(os.Stdin),
		eof:     false,
		env:     newEnv(),
	}
}

func newEnv() map[string]MalObject {
	env := make(map[string]MalObject)

	env["+"] = NewMalFunction(primitiveAdd)
	env["-"] = NewMalFunction(primitiveSubtract)
	env["*"] = NewMalFunction(primitiveMultiply)
	env["/"] = NewMalFunction(primitiveDivide)

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
	fmt.Printf("Eval(%T %s)\n", object, object)
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
		return m.Apply(f, args)

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
		symbol := ast.(MalSymbol).Symbol()
		if strings.HasPrefix(symbol, ":") {
			return ast, nil, nil
		}
		return m.env[symbol], nil, nil

	case MalList:
		list := ast.(MalList)
		if len(list.List()) == 0 {
			return ast, nil, nil
		}
		var evaluated []MalObject
		for _, o := range list.List() {
			e, err := m.Eval(o)
			if err != nil {
				return nil, nil, err
			}
			evaluated = append(evaluated, e)
		}
		return evaluated[0], evaluated[1:], nil

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

func (m *mal) Apply(f MalObject, args []MalObject) (MalObject, error) {
	fun, ok := f.(MalFunction)
	if !ok {
		return nil, fmt.Errorf("expected a MalFunction but got a %T", f)
	}

	return fun.Value()(args)
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
