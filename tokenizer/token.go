package tokenizer

type TokenType int

const (
	EOF TokenType = iota // End of file
	EOL                  // End of line
	LOG_LEVEL
	IP
	TIMESTAMP
	STRING
	UNRECOGNIZED
)

var tokens = []string{
	EOF:          "EOF",
	EOL:          "EOL",
	LOG_LEVEL:    "LOG_LEVEL",
	IP:           "IP",
	TIMESTAMP:    "TIMESTAMP",
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
