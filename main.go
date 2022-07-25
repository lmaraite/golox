package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lmaraite/golox/astprinter"
	"github.com/lmaraite/golox/lexer"
	"github.com/lmaraite/golox/parser"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	data, err := os.ReadFile(path)
	check(err)
	err = run(string(data))
	check(err)
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		check(err)
		if line == "\n" {
			break
		}
		err = run(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func run(source string) error {
	lexer := lexer.NewLexer(source)
	tokens, err := lexer.ScanTokens(source)
	if err != nil {
		return err
	}
	parser := parser.NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		return err
	}
	astprinter := astprinter.AstPrinter{}
	fmt.Println(astprinter.Print(expression))
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
