package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Mal interface {
	Read() (MalObject, error)
	Eval(MalObject) MalObject
	Print(MalObject) error
}

var (
	ErrMalEof = errors.New("MAL EOF")
)

type mal struct {
	scanner *bufio.Scanner
	eof     bool
}

var _ Mal = (*mal)(nil)

func newMal() *mal {
	return &mal{
		scanner: bufio.NewScanner(os.Stdin),
		eof:     false,
	}
}

func (m *mal) Read() (MalObject, error) {
	if m.eof {
		return nil, ErrMalEof
	}
	m.eof = !m.scanner.Scan()
	return ReadStr(m.scanner.Text())
}

func (m *mal) Eval(object MalObject) MalObject {
	return object
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
		if err == ErrMalEof {
			fmt.Printf("error: %s\n", err)
		} else if err != nil {
			return err
		}
		m.Print(m.Eval(in))
	}
}

func main() {
	err := rep()
	if err != nil {
		panic(err.Error())
	}
}
