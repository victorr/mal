package main

import "fmt"

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
