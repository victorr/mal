package main

import (
	"fmt"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalObject

type MalObject interface {
	String() string
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalNumber

type MalNumber interface {
	MalObject
	Number() int
}

type malNumber struct {
	n int
}

var _ MalNumber = (*malNumber)(nil)

func NewMalNumber(n int) MalNumber {
	return &malNumber{n: n}
}

func (n *malNumber) Number() int {
	return n.n
}

func (n *malNumber) String() string {
	return strconv.Itoa(n.n)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalString

type MalString interface {
	MalObject
	Value() string
}

type malString struct {
	s string
}

var _ MalString = (*malString)(nil)

func NewMalString(s string) MalString {
	return &malString{s: s}
}

func (s *malString) Value() string {
	return s.s
}

func (s *malString) String() string {
	ret := strings.ReplaceAll(s.Value(), `"`, `\"`)
	ret = strings.ReplaceAll(ret, "\n", `\n`)

	return fmt.Sprintf("\"%s\"", ret)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalSymbol

type MalSymbol interface {
	MalObject
	Symbol() string
}

type malSymbol struct {
	s string
}

var _ MalSymbol = (*malSymbol)(nil)

func NewMalSymbol(s string) MalSymbol {
	return &malSymbol{s: s}
}

func (s *malSymbol) Symbol() string {
	return s.s
}

func (s *malSymbol) String() string {
	return s.Symbol()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalList

type MalList interface {
	MalObject
	List() []MalObject
}
type malList struct {
	objects []MalObject
}

var _ MalList = (*malList)(nil)

func NewMalList(objects []MalObject) MalObject {
	return &malList{
		objects: objects,
	}
}

func (l *malList) List() []MalObject {
	return l.objects
}

func (l *malList) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, o := range l.List() {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(o.String())
	}
	sb.WriteString(")")

	return sb.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalVector

type MalVector interface {
	MalObject
	Vector() []MalObject
}

type malVector struct {
	objects []MalObject
}

var _ MalVector = (*malVector)(nil)

func NewMalVector(objects []MalObject) MalObject {
	return &malVector{
		objects: objects,
	}
}

func (v *malVector) Vector() []MalObject {
	return v.objects
}

func (v *malVector) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, o := range v.Vector() {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(o.String())
	}
	sb.WriteString("]")

	return sb.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalHashMap

type MalHashMap interface {
	MalObject
	HashMap() []MalObject
}
type malHashMap struct {
	objects []MalObject
}

var _ MalHashMap = (*malHashMap)(nil)

func NewMalHashMap(objects []MalObject) MalObject {
	return &malHashMap{
		objects: objects,
	}
}

func (m *malHashMap) HashMap() []MalObject {
	return m.objects
}

func (m *malHashMap) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i < len(m.objects); i += 2 {
		key := m.objects[i]
		value := m.objects[i+1]

		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(key.String())
		sb.WriteString(" ")
		sb.WriteString(value.String())
	}
	sb.WriteString("}")

	return sb.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalFunction

type ApplyFunc func(args []MalObject) (MalObject, error)

type MalFunction interface {
	MalObject
	Value() ApplyFunc
}

type malFunction struct {
	f ApplyFunc
}

var _ MalFunction = (*malFunction)(nil)

func NewMalFunction(f ApplyFunc) MalObject {
	return &malFunction{f: f}
}

func (f *malFunction) Value() ApplyFunc {
	return f.f
}

func (f *malFunction) String() string {
	return "#function"
}
