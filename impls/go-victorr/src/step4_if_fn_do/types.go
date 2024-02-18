package main

import (
	"fmt"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalObject

type MalObject interface {
	Equals(MalObject) MalBoolean
	String() string
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalNil

type malNil struct {
	MalObject
}

var MalNil = &malNil{}

var _ MalObject = (*malNil)(nil)

func (n *malNil) Equals(other MalObject) MalBoolean {
	if n == other {
		return MalTrue
	}
	return MalFalse
}

func (n *malNil) String() string {
	return "nil"
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalBoolean

type MalBoolean interface {
	MalObject
	IsTrue() bool
	IsFalse() bool
}

type malBoolean struct {
	b bool
}

var MalTrue = &malBoolean{b: true}

var MalFalse = &malBoolean{b: false}

var _ MalBoolean = (*malBoolean)(nil)

func (b *malBoolean) IsTrue() bool {
	return b.b
}

func (b *malBoolean) IsFalse() bool {
	return !b.b
}

func (b *malBoolean) Equals(other MalObject) MalBoolean {
	if ob, ok := other.(MalBoolean); ok && b.b == ob.IsTrue() {
		return MalTrue
	}
	return MalFalse
}

func (b *malBoolean) String() string {
	if b.b {
		return "true"
	}
	return "false"
}

func IsTrue(o MalObject) bool {
	if b, ok := o.(MalBoolean); ok {
		return b.IsTrue()
	}
	return o != MalNil
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

func (n *malNumber) Equals(other MalObject) MalBoolean {
	if on, ok := other.(*malNumber); ok && n.n == on.n {
		return MalTrue
	}
	return MalFalse
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

func (s *malString) Equals(other MalObject) MalBoolean {
	if os, ok := other.(*malString); ok && s.s == os.s {
		return MalTrue
	}
	return MalFalse
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

func (s *malSymbol) Equals(other MalObject) MalBoolean {
	if os, ok := other.(*malSymbol); ok && s.s == os.s {
		return MalTrue
	}
	return MalFalse
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

func (l *malList) Equals(other MalObject) MalBoolean {
	return listEquals(l.List(), other)
}

func listEquals(left []MalObject, right MalObject) MalBoolean {
	var l2 []MalObject
	if vo, ok := right.(MalVector); ok {
		l2 = vo.Vector()
	} else if l, ok := right.(MalList); ok {
		l2 = l.List()
	} else {
		return MalFalse
	}
	fmt.Printf("l1=%s, l2=%s\n", left, l2)
	if len(left) != len(l2) {
		return MalFalse
	}
	for index := 0; index < len(left); index++ {
		if left[index].Equals(l2[index]).IsFalse() {
			return MalFalse
		}
	}
	return MalTrue
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

func (v *malVector) Equals(other MalObject) MalBoolean {
	return listEquals(v.Vector(), other)
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

func (m *malHashMap) Equals(other MalObject) MalBoolean {
	if om, ok := other.(*malHashMap); ok {
		l1 := m.HashMap()
		l2 := om.HashMap()
		if len(l1) != len(l2) {
			return MalFalse
		}
		for index := 0; index < len(l1); index++ {
			if l1[index].Equals(l2[index]).IsFalse() {
				return MalFalse
			}
		}
	}
	return MalTrue
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

type ApplyFunc func(env MalEnv, args []MalObject) (MalObject, error)

type SpecialForm uint

const (
	NotSpecialForm SpecialForm = iota
	DefSpecialForm
	LetSpecialForm
	IfSpecialForm
)

type MalFunction interface {
	MalObject
	Function() ApplyFunc
	SpecialForm() SpecialForm
}

type malFunction struct {
	f           ApplyFunc
	specialForm SpecialForm
}

var _ MalFunction = (*malFunction)(nil)

func NewMalFunction(f ApplyFunc, specialForm SpecialForm) MalObject {
	return &malFunction{
		f:           f,
		specialForm: specialForm,
	}
}

func (f *malFunction) Function() ApplyFunc {
	return f.f
}

func (f *malFunction) SpecialForm() SpecialForm {
	return f.specialForm
}

func (f *malFunction) Equals(other MalObject) MalBoolean {
	if f == other {
		return MalTrue
	}
	return MalFalse
}
func (f *malFunction) String() string {
	return "#<function>"
}
