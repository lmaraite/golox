package lexer

import (
	"errors"
	"fmt"
)

type Token int

func newError(line int, message string) error {
	formattedMessage := fmt.Sprintf("[line %d] Error: %s", line, message)
	return errors.New(formattedMessage)
}

func scanTokens(source string) ([]Token, error) {
	return nil, newError(0, "unimplemented method")
}

func Run(source string) error {
	tokens, err := scanTokens(source)
	if err != nil {
		return err
	}
	for token := range tokens {
		fmt.Println(token)
	}
	return nil
}
