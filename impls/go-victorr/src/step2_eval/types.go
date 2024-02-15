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
	Value() int
}

type malNumber struct {
	n int
}

var _ MalNumber = (*malNumber)(nil)

func NewMalNumber(n int) MalNumber {
	return &malNumber{n: n}
}

func (n *malNumber) Value() int {
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
	Value() string
}

type malSymbol struct {
	s string
}

var _ MalSymbol = (*malSymbol)(nil)

func NewMalSymbol(s string) MalSymbol {
	return &malSymbol{s: s}
}

func (s *malSymbol) Value() string {
	return s.s
}

func (s *malSymbol) String() string {
	return s.Value()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalList

type MalList interface {
	MalObject
	Objects() []MalObject
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

func (l *malList) Objects() []MalObject {
	return l.objects
}

func (l *malList) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, o := range l.Objects() {
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
	Objects() []MalObject
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

func (v *malVector) Objects() []MalObject {
	return v.objects
}

func (v *malVector) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, o := range v.Objects() {
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
	Objects() []MalObject
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

func (m *malHashMap) Objects() []MalObject {
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
