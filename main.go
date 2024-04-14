package main

import (
	"dominiqueboerner/go-log-lexer/tokenizer"
	"fmt"
	"os"
)

// Configuration
var logFile = "examples/test-1000.log"

func main() {
	file, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}

	var fileTokenizer = tokenizer.NewTokenizer(file)

	fmt.Printf("Start tokenizing '%s'...\n", file.Name())

	tokens, err := fileTokenizer.Tokenize()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished tokenizing. %d tokens generated\n", len(tokens))

	file, err = os.Create("output_file.txt")
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	for _, token := range tokens {
		_, err := fmt.Fprintln(file, token)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Tokens have been written to output_file.txt")
}
