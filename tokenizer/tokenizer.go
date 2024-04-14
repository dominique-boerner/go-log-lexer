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
	timestampRegex   = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	logLevelRegex    = regexp.MustCompile(`\[(INFO|WARNING|ERROR|DEBUG)\]`)
	ipRegex          = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	singleWordRegexp = regexp.MustCompile(`[A-Za-z]+\s`)
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

		// Extract timestamp
		if match := timestampRegex.FindString(l.line); match != "" {
			tokens = append(tokens, Token{Type: TIMESTAMP, Content: match, Pos: l.pos})
			l.line = strings.TrimPrefix(l.line, match)
			l.pos.Column += len(match)
		}

		// Extract log level
		if match := logLevelRegex.FindString(l.line); match != "" {
			tokens = append(tokens, Token{Type: LOG_LEVEL, Content: match, Pos: l.pos})
			l.line = strings.TrimPrefix(l.line, match)
			l.pos.Column += len(match)
		}

		// Extract IP address if present
		if match := ipRegex.FindString(l.line); match != "" {
			tokens = append(tokens, Token{Type: IP, Content: match, Pos: l.pos})
			l.line = strings.TrimPrefix(l.line, match)
			l.pos.Column += len(match)
		}

		// Extract string if present
		if match := singleWordRegexp.FindString(l.line); match != "" {
			matchedString := strings.TrimSpace(match)
			tokens = append(tokens, Token{Type: STRING, Content: matchedString, Pos: l.pos})
			l.line = strings.TrimPrefix(l.line, matchedString)
			l.pos.Column += len(matchedString)
		}

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
