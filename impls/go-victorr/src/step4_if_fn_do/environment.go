package main

import (
	"fmt"
	"strings"
)

type MalEnv interface {
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

func (e *malEnv) Equals(other MalObject) MalBoolean {
	// oe, ok := other.(MalEnv)
	// if !ok {
	// 	return MalFalse
	// }
	return MalFalse
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
	var sb strings.Builder
	sb.WriteString("#<environment\n")
	for k, v := range e.env {
		sb.WriteString(fmt.Sprintf("\t%s => %s\n", k, v))
	}
	sb.WriteString(">")
	//return "#<environment>"
	return sb.String()
}
