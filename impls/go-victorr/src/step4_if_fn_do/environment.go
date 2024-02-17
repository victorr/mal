package main

import "fmt"

type MalEnv interface {
	MalObject
	Set(MalSymbol, MalObject) MalObject
	Get(MalSymbol) (MalObject, error)
}

type malEnv struct {
	outer MalEnv
	env   map[string]MalObject
}

var _ MalEnv = (*malEnv)(nil)

func NewMalEnv(outer MalEnv) MalEnv {
	return &malEnv{
		outer: outer,
		env:   make(map[string]MalObject),
	}
}

func (e *malEnv) Set(symbol MalSymbol, value MalObject) MalObject {
	e.env[symbol.Symbol()] = value

	return value
}

func (e *malEnv) Get(symbol MalSymbol) (MalObject, error) {
	symbolStr := symbol.Symbol()
	if _, ok := e.env[symbolStr]; ok {
		return e.env[symbolStr], nil
	}
	if e.outer != nil {
		return e.outer.Get(symbol)
	}
	return nil, errSymbolNotFound(symbol.Symbol())
}

func errSymbolNotFound(symbol string) error {
	return fmt.Errorf("'%s' not found", symbol)
}

func (e *malEnv) String() string {
	return "#environment"
}
