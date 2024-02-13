package main

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalObject

type MalObject interface{}

////////////////////////////////////////////////////////////////////////////////////////////////////
// MalAtom

type MalAtom interface {
	MalObject
	Token() string
}

type malAtom struct {
	token string
}

var _ MalAtom = (*malAtom)(nil)

func NewAtom(token string) MalAtom {
	return &malAtom{
		token: token,
	}
}

func (a *malAtom) Token() string {
	return a.token
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

func NewList(objects []MalObject) MalList {
	return &malList{
		objects: objects,
	}
}

func (l *malList) Objects() []MalObject {
	return l.objects
}
