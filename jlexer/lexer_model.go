package jlexer

func (r *Lexer) ResetError(e error) {
	r.fatalError = e
}
