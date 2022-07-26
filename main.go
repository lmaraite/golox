package main

import (
	"fmt"
	"os"

	"github.com/lmaraite/golox/interpreter"
	"github.com/lmaraite/golox/lexer"
	"github.com/lmaraite/golox/parser"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	check(err)
	err = run(string(data))
	check(err)
}

func run(source string) error {
	lexer := lexer.NewLexer(source)
	tokens, err := lexer.ScanTokens(source)
	if err != nil {
		return err
	}
	parser := parser.NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	interpreter := interpreter.NewInterpreter()
	err = interpreter.Interpret(statements)
	if err != nil {
		return err
	}

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
