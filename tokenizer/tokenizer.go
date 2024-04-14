package tokenizer

import (
	"bufio"
	"dominiqueboerner/go-log-lexer/analyzer"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	reader   *bufio.Reader
	analyzer *analyzer.Analyzer
	file     *os.File
	line     string
	pos      Position
}

var (
	dateRegex     = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	timeRegex     = regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	logLevelRegex = regexp.MustCompile(`\[(INFO|WARNING|ERROR|DEBUG)\]`)
	ipRegex       = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
)

func NewTokenizer(file *os.File) *Lexer {
	lex := &Lexer{
		reader:   bufio.NewReader(file),
		analyzer: analyzer.NewAnalyzer(file),
	}
	lex.file = file
	return lex
}

// Tokenize
// Tokenizes the given file line-by-line.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	var err error

	l.analyzer.PrintFileInformation()
	linesCount := l.analyzer.Lines
	startTime := time.Now()

	for {
		l.line, err = l.reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		l.pos.Line++
		l.pos.Column = 0

		// TODO: this operation is quite expensive, i should implement some kind of normalizing...
		words := strings.Fields(l.line)
		for _, word := range words {
			// Check if the word matches any special patterns
			if dateRegex.MatchString(word) {
				tokens = append(tokens, Token{Type: TIME, Content: word, Pos: l.pos})
			} else if timeRegex.MatchString(word) {
				tokens = append(tokens, Token{Type: DATE, Content: word, Pos: l.pos})
			} else if logLevelRegex.MatchString(word) {
				tokens = append(tokens, Token{Type: LOG_LEVEL, Content: word, Pos: l.pos})
			} else if ipRegex.MatchString(word) {
				tokens = append(tokens, Token{Type: IP, Content: word, Pos: l.pos})
			} else {
				tokens = append(tokens, Token{Type: STRING, Content: word, Pos: l.pos})
			}
			l.pos.Column += len(word) + 1 // Increment column position by word length + space
		}

		tokens = append(tokens, Token{Type: EOL, Content: "", Pos: l.pos})

		logAfterLineTokenization(startTime, l, linesCount)
	}

	fmt.Printf("\n\nTotal Tokenization Time: %s\n", time.Since(startTime).Round(time.Second))

	err = l.file.Close()
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func logAfterLineTokenization(startTime time.Time, l *Lexer, linesCount int) {
	timeElapsed := time.Since(startTime)
	estimatedTotalTime := time.Duration(float64(timeElapsed) / float64(l.pos.Line) * float64(linesCount))
	timeRemaining := estimatedTotalTime - timeElapsed
	if l.pos.Line > linesCount {
		fmt.Printf("\rProgress: 100.00%%. Time remaining: 00:00:00\n")
	} else {
		percentage := (float64(l.pos.Line) / float64(linesCount)) * 100
		fmt.Printf("\rProgress: %.2f%%. Time remaining: %s", percentage, timeRemaining.Round(time.Second))
	}
}
