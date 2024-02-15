package main

import "fmt"

func primitiveAdd(args []MalObject) (MalObject, error) {
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

func primitiveSubtract(args []MalObject) (MalObject, error) {
	// ret := 0
	// switch len(args) {
	// case 0:
	// 	// nothing to do
	// case 1:
	// 	number, ok := args[0].(MalNumber)
	// 	if !ok {
	// 		return nil, fmt.Errorf("expected a MalNumber but got %T", arg)
	// 	}
	// 	ret = number.Number()

	// }
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

func primitiveMultiply(args []MalObject) (MalObject, error) {
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

func primitiveDivide(args []MalObject) (MalObject, error) {
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
