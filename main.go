package main

import (
	"dominiqueboerner/go-log-lexer/config"
	"dominiqueboerner/go-log-lexer/tokenizer"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
)

const (
	configFile = "config.toml"
)

func main() {
	// Read config file
	c, err := loadConfig()
	if err != nil {
		fmt.Printf("Error while reading configuration file")
		c = config.Config{}
	}

	file, err := os.Open(c.Configuration.LogFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Start tokenizing '%s'...\n\n", file.Name())
	var fileTokenizer = tokenizer.NewTokenizer(file)

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

func loadConfig() (config.Config, error) {
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		return config.Config{}, err
	}

	var c config.Config
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		return config.Config{}, err
	}

	return c, nil
}
