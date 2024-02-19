package main

import (
	"fmt"
	"strings"
)

func primitiveAdd(_ MalEnv, args []MalObject) (MalObject, error) {
	ret := 0
	for _, arg := range args {
		number, ok := arg.(MalNumber)
		if !ok {
			return nil, fmt.Errorf("expected a MalNumber but got %T", arg)
		}
		ret += number.Number()
	}
	return NewMalNumber(ret), nil
}

func primitiveSubtract(_ MalEnv, args []MalObject) (MalObject, error) {
	ret := 0
	for index, arg := range args {
		number, ok := arg.(MalNumber)
		if !ok {
			return nil, fmt.Errorf("expected a MalNumber but got %T", arg)
		}
		if index == 0 {
			ret = number.Number()
		} else {
			ret -= number.Number()
		}
	}
	if len(args) == 1 {
		ret = -ret
	}
	return NewMalNumber(ret), nil
}

func primitiveMultiply(_ MalEnv, args []MalObject) (MalObject, error) {
	ret := 1
	for _, arg := range args {
		number, ok := arg.(MalNumber)
		if !ok {
			return nil, fmt.Errorf("expected a MalNumber but got %T", arg)
		}
		ret *= number.Number()
	}
	return NewMalNumber(ret), nil
}

func primitiveDivide(_ MalEnv, args []MalObject) (MalObject, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("zero arguments for /")
	}
	ret := 0
	for index, arg := range args {
		number, ok := arg.(MalNumber)
		if !ok {
			return nil, fmt.Errorf("expected a MalNumber but got %T", arg)
		}
		if number.Number() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		if index == 0 {
			ret = number.Number()
		} else {
			ret /= number.Number()
		}
	}
	if len(args) == 1 {
		ret = 0
	}
	return NewMalNumber(ret), nil
}

func primitiveDef(env MalEnv, args []MalObject) (MalObject, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected two arguments for def!, but got %d", len(args))
	}

	symbol, ok := args[0].(MalSymbol)
	if !ok {
		return nil, fmt.Errorf("expected a MalSymbol, but got %T", symbol)
	}
	value := args[1].(MalObject)

	return env.Set(symbol, value), nil
}

func primitiveLet(env MalEnv, args []MalObject) (MalObject, error) {
	return nil, nil
}

func primitiveNumberComparison(symbol string, comp func(MalNumber, MalNumber) bool) ApplyFunc {
	return func(_ MalEnv, args []MalObject) (MalObject, error) {
		switch len(args) {
		case 0:
			return nil, fmt.Errorf("wrong number of arguments for %s", symbol)

		case 1:
			return MalTrue, nil
		}
		var left MalNumber
		for _, arg := range args {
			right, ok := arg.(MalNumber)
			if !ok {
				return nil, fmt.Errorf("expected a MalNumber for %s but got %T", symbol, arg)
			}
			if left != nil {
				if !comp(left, right) {
					return MalFalse, nil
				}
			}
			left = right
		}
		return MalTrue, nil
	}
}

func primitiveList(_ MalEnv, args []MalObject) (MalObject, error) {
	return NewMalList(args), nil
}

func primitivePredicate(predicate func(o MalObject) (bool, error)) ApplyFunc {
	return func(_ MalEnv, args []MalObject) (MalObject, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("wrong number of arguments for predicate")
		}

		for _, arg := range args {
			result, err := predicate(arg)
			if err != nil {
				return nil, err
			}
			if !result {
				return MalFalse, nil
			}
		}
		return MalTrue, nil
	}
}

func primitiveCount(_ MalEnv, args []MalObject) (MalObject, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments for count, expected 1 but got %d", len(args))
	}
	if args[0] == MalNil {
		return NewMalNumber(0), nil
	}
	if list, ok := args[0].(MalList); ok {
		return NewMalNumber(len(list.List())), nil
	}
	if vector, ok := args[0].(MalVector); ok {
		return NewMalNumber(len(vector.Vector())), nil
	}
	return nil, errWrongType("MalList", args[0])
}

func primitiveEquals(_ MalEnv, args []MalObject) (MalObject, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments for =")
	}
	first := args[0]
	for _, arg := range args[1:] {
		if first.Equals(arg).IsFalse() {
			return MalFalse, nil
		}
	}
	return MalTrue, nil
}

func primitivePrn(env MalEnv, args []MalObject) (MalObject, error) {
	s, err := primitivePrStr(env, args)
	if err != nil {
		return nil, err
	}
	fmt.Println(s.(MalString).Value())
	return MalNil, nil
}

func primitivePrintln(env MalEnv, args []MalObject) (MalObject, error) {
	var sb strings.Builder
	for index, arg := range args {
		if index > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(arg.Print(false))
	}
	fmt.Println(sb.String())
	return MalNil, nil
}

func primitiveStr(env MalEnv, args []MalObject) (MalObject, error) {
	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(arg.Print(false))
	}

	return NewMalString(sb.String()), nil
}

func primitivePrStr(_ MalEnv, args []MalObject) (MalObject, error) {
	var s []string
	for _, arg := range args {
		s = append(s, arg.Print(true))
	}

	return NewMalString(strings.Join(s, " ")), nil
}

func errWrongType(expected string, got any) error {
	return fmt.Errorf("wrong argument, expected a %s but got %T", expected, got)
}
