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
	Values map[string]interface{}
}

func (e *Environment) define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) get(name token.Token) (interface{}, error) {
	if value, ok := e.Values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, newError(name, "Undefined variable '"+name.Lexeme+"'.")
}
