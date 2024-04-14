package tokenizer

type TokenType int

const (
	EOF TokenType = iota // End of file
	EOL                  // End of line
	LOG_LEVEL
	IP
	TIME
	DATE
	STRING
	UNRECOGNIZED
)

var tokens = []string{
	EOF:          "EOF",
	EOL:          "EOL",
	LOG_LEVEL:    "LOG_LEVEL",
	IP:           "IP",
	TIME:         "TIME",
	DATE:         "DATE",
	STRING:       "STRING",
	UNRECOGNIZED: "UNRECOGNIZED",
}

func (t TokenType) String() string {
	return tokens[t]
}

type Token struct {
	Type    TokenType
	Content string
	Pos     Position
}
