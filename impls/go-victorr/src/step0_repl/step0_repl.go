package main

import (
	"bufio"
	"fmt"
	"os"
)

type Mal interface {
	Read() string
	Eval(in string) string
	Print(in string)
}

const malEOF = "the-end-of-file"

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

func (m *mal) Read() string {
	if m.eof {
		return malEOF
	}
	m.eof = !m.scanner.Scan()
	return m.scanner.Text()
}

func (m *mal) Eval(in string) string {
	return in
}

func (m *mal) Print(in string) {
	fmt.Println(in)
}

func rep() {
	m := newMal()
	for {
		fmt.Print("user> ")
		in := m.Read()
		if in == malEOF {
			return
		}
		m.Print(m.Eval(in))
	}
}

func main() {
	rep()
}
