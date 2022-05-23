package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lmaraite/golox/lexer"
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
	err = lexer.Run(string(data))
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
		err = lexer.Run(line)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
