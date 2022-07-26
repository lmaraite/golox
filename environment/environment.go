package environment

import (
	"errors"
	"fmt"

	"github.com/lmaraite/golox/token"
)

func newError(errorToken token.Token, message string) error {
	var formattedMessage string
	if errorToken.TokenType == token.EOF {
		formattedMessage = fmt.Sprintf("[line %d] Runtime error at end: %s", errorToken.Line, message)
	} else {
		formattedMessage = fmt.Sprintf("[line %d] Runtime error at '%s': %s", errorToken.Line, errorToken.Lexeme, message)
	}
	return errors.New(formattedMessage)
}

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) Assign(name token.Token, value interface{}) error {
	if value, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return newError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (interface{}, error) {
	if value, ok := e.values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, newError(name, "Undefined variable '"+name.Lexeme+"'.")
}
