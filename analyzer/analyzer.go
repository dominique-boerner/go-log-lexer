package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Analyzer
// Analyzes the file and returns some additional information about it.
type Analyzer struct {
	scanner  *bufio.Scanner
	file     *os.File
	Lines    int
	Filesize int
}

func NewAnalyzer(file *os.File) *Analyzer {
	a := &Analyzer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}

	a.Lines = a.countLines()
	a.Filesize = a.getFileSize()

	return a
}

func (a *Analyzer) PrintFileInformation() {
	fmt.Printf("=== File Information ===\n")
	fmt.Printf("Filesize: %d byte\n", a.Filesize)
	fmt.Printf("Number of lines: %d\n", a.Lines)
	fmt.Printf("========================\n\n")
}

// GetFileSize
// Returns the size of the file.
func (a *Analyzer) getFileSize() int {
	fileInfo, _ := a.file.Stat()
	return int(fileInfo.Size())
}

// CountLines
// Returns the amount of lines the file has and resets the cursor.
func (a *Analyzer) countLines() int {
	var linesCount int
	for a.scanner.Scan() {
		linesCount++
	}

	_, _ = a.file.Seek(0, io.SeekStart)

	return linesCount
}
