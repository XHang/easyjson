package jlexer

import "fmt"

// LexerError implements the error interface and represents all possible errors that can be
// generated during parsing the JSON data.
type LexerError struct {
	Reason string
	Offset int
	Data   string
}

func (l *LexerError) Error() string {
	// orm.Model 重置错误提示
	if l.Offset == 0 && l.Reason == "" {
		return l.Data
	}

	return fmt.Sprintf("parse error: %s near offset %d of '%s'", l.Reason, l.Offset, l.Data)
}
