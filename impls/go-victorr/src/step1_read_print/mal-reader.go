package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type MalReader interface {
	Next() string
	Peek() string
}

type malReader struct {
	tokens  []string
	current int
}

var _ MalReader = (*malReader)(nil)

func NewMalReader(tokens []string) MalReader {
	return &malReader{
		tokens:  tokens,
		current: 0,
	}
}

func (r *malReader) Next() string {
	if r.current >= len(r.tokens) {
		return ""
	}
	r.current++
	return r.tokens[r.current-1]
}

func (r *malReader) Peek() string {
	if r.current >= len(r.tokens) {
		return ""
	}

	return r.tokens[r.current]
}

func ReadStr(in string) (MalObject, error) {
	r := tokenize(in)
	return ReadForm(r)
}

func tokenize(in string) MalReader {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

	matches := re.FindAllStringSubmatch(in, -1)
	// matches is a [][]string
	// each entry is [whole-match subgroup]

	// fmt.Printf("The matches are: %v\n", matches)
	// fmt.Printf("The matches 0 are:\n")
	// for _, t := range matches[0] {
	// 	fmt.Printf("  - '%v'\n", t)
	// }

	var tokens []string
	for _, m := range matches {
		t := strings.TrimSpace(m[1])
		if t != "" {
			tokens = append(tokens, t)
		}
	}

	return NewMalReader(tokens)
}

func ReadForm(r MalReader) (MalObject, error) {
	next := r.Peek()
	switch {
	case next == "":
		return nil, nil

	case next == "(":
		return ReadList(r)

	default:
		return ReadAtom(r)
	}
}

func ReadAtom(r MalReader) (MalAtom, error) {
	token := r.Next()
	switch {
	case token == "":
		return nil, ErrMalEof
	case token == "(":
		return nil, errors.New("expected an atom but got a list")
	default:
		return NewAtom(token), nil
	}
}

func ReadList(r MalReader) (MalList, error) {
	token := r.Next() // drop the '('
	if token == "" || token != "(" {
		return nil, fmt.Errorf("expected an open paren '(' but got '%s'", token)
	}
	var forms []MalObject
	for {
		if r.Peek() == ")" || r.Peek() == "" {
			break
		}
		// fmt.Printf("ReadList peek = '%v'\n", r.Peek())
		form, err := ReadForm(r)
		if err != nil {
			return nil, err
		}
		if form == nil {
			return nil, errors.New("unexpected end of list")
		}
		forms = append(forms, form)
	}
	token = r.Next() // drop the ')'
	if token == "" {
		return nil, ErrMalEof
	}
	if token != ")" {
		return nil, fmt.Errorf("expected a close paren ')' but got '%s'", token)
	}

	return NewList(forms), nil
}
