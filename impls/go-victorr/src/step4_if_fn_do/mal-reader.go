package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

	var tokens []string
	for _, m := range matches {
		t := strings.TrimSpace(m[1])
		if t == "" {
			continue
		}

		tokens = append(tokens, t)
	}

	return NewMalReader(tokens)
}

func ReadForm(r MalReader) (MalObject, error) {
	next := r.Peek()
	switch {
	case next == "":
		return nil, nil

	case next == "(":
		return ReadList("(", ")", NewMalList, r)

	case next == "[":
		return ReadList("[", "]", NewMalVector, r)

	case next == "{":
		return ReadHashMap(r)

	case next == "'":
		return ReadQuote("'", "quote", r)

	case next == "~":
		return ReadQuote("~", "unquote", r)

	case next == "~@":
		return ReadQuote("~@", "splice-unquote", r)

	case next == "`":
		return ReadQuote("`", "quasiquote", r)

	case next == "@":
		return ReadQuote("@", "deref", r)

	case next == "^":
		return ReadWithMeta(r)

	default:
		return ReadAtom(r)
	}
}

func ReadAtom(r MalReader) (MalObject, error) {
	token := r.Next()

	numberRe := regexp.MustCompile("^[-+]{0,1}[0-9]+$")

	switch {
	case token == "":
		return nil, ErrMalEof
	case token == "(":
		return nil, errors.New("expected an atom but got a list")
	case token == "[":
		return nil, errors.New("expected an atom but got a vector")
	case numberRe.MatchString(token):
		// Numbers are only ints for now.

		n, err := strconv.ParseInt(token, 10, 0)
		if err != nil {
			return nil, err
		}
		return NewMalNumber(int(n)), nil

	case strings.HasPrefix(token, "\""): //stringRe.MatchString(token):
		return ReadMalString(token)

	case token == "nil":
		return MalNil, nil

	case token == "true":
		return MalTrue, nil

	case token == "false":
		return MalFalse, nil

	default:
		// symbols
		return NewMalSymbol(token), nil
	}
}

func ReadMalString(token string) (MalString, error) {
	re := regexp.MustCompile(`^"(?:\\.|[^\\"])*"$`)
	if !re.MatchString(token) {
		return nil, errors.New("unbalanced '\"', missing at end of string")
	}

	const quote = "\""
	s := strings.TrimPrefix(token, quote)
	s = strings.TrimSuffix(s, quote)

	s = strings.ReplaceAll(s, `\\`, "\u029e") // lifted from the JS implementation (or was it the Java one?)
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, "\u029e", `\`)

	return NewMalString(s), nil
}

type listFactory func([]MalObject) MalObject

func ReadList(open, close string, newList listFactory, r MalReader) (MalObject, error) {
	token := r.Next() // drop the '('
	if token == "" || token != open {
		return nil, fmt.Errorf("expected '%s' but got '%s'", open, token)
	}
	var forms []MalObject
	for {
		if r.Peek() == close || r.Peek() == "" {
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
	if token != close {
		return nil, fmt.Errorf("expected a '%s' but got '%s'", close, token)
	}

	return newList(forms), nil
}

func ReadQuote(prefix, symbol string, r MalReader) (MalObject, error) {
	token := r.Next() // drop the "'"
	if token == "" || token != prefix {
		return nil, fmt.Errorf("expected %s but got '%s'", prefix, token)
	}
	var forms []MalObject
	forms = append(forms, NewMalSymbol(symbol))
	form, err := ReadForm(r)
	if err != nil {
		return nil, err
	}
	if form == nil {
		return nil, errors.New("unexpected end of list")
	}
	forms = append(forms, form)

	return NewMalList(forms), nil
}

func ReadHashMap(r MalReader) (MalObject, error) {
	token := r.Next()
	if token == "" || token != "{" {
		return nil, fmt.Errorf("expected '{' but got %s", token)
	}
	var values []MalObject
	for {
		if r.Peek() == "}" {
			break
		}

		key, err := ReadAtom(r)
		if err != nil {
			return nil, err
		}
		value, err := ReadForm(r)
		if err != nil {
			return nil, err
		}
		values = append(values, key, value)
	}
	token = r.Next()
	if token == "" || token != "}" {
		return nil, fmt.Errorf("expected '}' but got %s", token)
	}

	return NewMalHashMap(values), nil
}

func ReadWithMeta(r MalReader) (MalObject, error) {
	token := r.Next()
	if token == "" || token != "^" {
		return nil, fmt.Errorf("expected '^' but got %s", token)
	}
	meta, err := ReadForm(r)
	if err != nil {
		return nil, err
	}
	object, err := ReadForm(r)
	if err != nil {
		return nil, err
	}
	objects := []MalObject{
		NewMalSymbol("with-meta"),
		object,
		meta,
	}
	return NewMalList(objects), nil
}
