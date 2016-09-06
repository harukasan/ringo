package scanner

import (
	"fmt"
	"io"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

// Scanner implements a scanner for Ruby lex.
type Scanner struct {
	src []byte // source buffer
	err error  //

	char    byte // current read character
	offset  int  // current offset
	begin   int  // offset of begin of the token
	nospace bool // whether the previous is not a space
}

// New returns a initiazlied scanner to scan script source src.
func New(src []byte) *Scanner {
	s := new(Scanner)
	s.src = src
	s.offset = -1
	s.next()
	return s
}

func (s *Scanner) next() {
	s.offset++
	if s.offset >= len(s.src) {
		s.err = io.EOF
		s.char = 0
		debug.Printf("next: len=%v, offset=%v, char=%v, err=%v", len(s.src), s.offset, s.char, s.err)
		return
	}
	s.char = s.src[s.offset]
	debug.Printf("next: len=%v, offset=%v, char=%v", len(s.src), s.offset, s.char)
}

func (s *Scanner) failf(format string, v ...interface{}) {
	if s.err != nil {
		return
	}
	s.err = fmt.Errorf(format, v...)
	debug.Printf("failf: %v", s.err)
}

// Scan reads and returns a parsed token position, type, and its literal.
func (s *Scanner) Scan() (pos int, t token.Token, literal []byte) {
StartScan:
	ch := s.char
	switch {
	case token.IsLetter(ch):
		pos, t, literal = s.scanIdent()
		s.nospace = true
		return
	}
	if s.err == io.EOF {
		pos, t, literal = s.offset, token.EOF, nil
		return
	}
	s.begin = s.offset
	s.next()
	if scan := tokenScanners[ch]; scan != nil {
		t, literal = scan(s)
		if t == token.Continue {
			goto StartScan
		}
		if t != token.NewLine {
			s.nospace = true
		}
		return s.begin, t, literal
	}
	return s.offset, token.Illegal, nil
}

func (s *Scanner) skipLine() {
	for {
		if s.err != nil {
			return
		}
		if s.char == '\n' {
			return
		}
		s.next()
	}
}

type scanFunc func(s *Scanner) (token.Token, []byte)

var tokenScanners = [127]scanFunc{
	0x09: skipWhiteSpaces,
	'\n': scanNewLine,
	0x0b: skipWhiteSpaces,
	0x0c: skipWhiteSpaces,
	0x0d: skipWhiteSpaces,
	'!':  scanNot,
	' ':  skipWhiteSpaces,
	'"':  scanDoubleQuoteString,
	'#':  scanComment,
	'%':  scanMod,
	'&':  scanAnd,
	'\'': scanSingleQuoteString,
	'(':  scanOne(token.LParen),
	')':  scanOne(token.RParen),
	'*':  scanAsterisk,
	'+':  scanPlus,
	'-':  scanMinus,
	'/':  scanDiv,
	'0':  scanZero,
	'1':  scanNonZero,
	'2':  scanNonZero,
	'3':  scanNonZero,
	'4':  scanNonZero,
	'5':  scanNonZero,
	'6':  scanNonZero,
	'7':  scanNonZero,
	'8':  scanNonZero,
	'9':  scanNonZero,
	'<':  scanLt,
	'=':  scanEq,
	'>':  scanGt,
	'[':  scanBracket,
	']':  scanOne(token.RBracket),
	'^':  scanXor,
	'{':  scanOne(token.LBrace),
	'|':  scanOr,
	'}':  scanOne(token.RBrace),
	'~':  scanOne(token.Invert),
}

func scanOne(tk token.Token) scanFunc {
	return func(s *Scanner) (token.Token, []byte) {
		return tk, nil
	}
}

func skipWhiteSpaces(s *Scanner) (token.Token, []byte) {
	s.nospace = false
	for token.IsWhiteSpace(s.char) {
		s.next()
	}
	return token.Continue, nil
}

func scanNewLine(s *Scanner) (token.Token, []byte) {
	s.nospace = false
	return token.NewLine, nil
}

func scanNot(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '=':
		s.next()
		return token.NotEqual, nil
	case '~':
		s.next()
		return token.NotMatch, nil
	}
	return token.Not, nil
}

func scanComment(s *Scanner) (token.Token, []byte) {
	s.skipLine()
	return token.Continue, nil
}

func scanDoubleQuoteString(s *Scanner) (token.Token, []byte) {
	for s.char != '"' && s.err == nil {
		s.next()
	}
	s.next()
	return token.StringPart, s.src[s.begin:s.offset]
}

func scanMod(s *Scanner) (token.Token, []byte) {
	if s.char == '=' { // %=
		s.next()
		return token.AssignMod, nil
	}
	return token.Mod, nil
}

func scanAnd(s *Scanner) (token.Token, []byte) {
	if s.char == '&' { // &&
		s.next()
		if s.char == '=' { // &&=
			s.next()
			return token.AssignAndOperator, nil
		}
		return token.AndOperator, nil
	}
	if s.char == '=' { // &=
		s.next()
		return token.AssignAnd, nil
	}
	return token.Amp, nil
}

func scanSingleQuoteString(s *Scanner) (token.Token, []byte) {
	for s.char != '\'' && s.err == nil {
		if s.char == '\\' {
			s.next()
		}
		s.next()
	}
	s.next()
	return token.StringPart, s.src[s.begin:s.offset]
}

func scanAsterisk(s *Scanner) (token.Token, []byte) {
	if s.char == '*' { // **
		s.next()
		if s.char == '=' { // **=
			s.next()
			return token.AssignPow, nil
		}
		return token.Pow, nil
	}
	if s.char == '=' { // *=
		s.next()
		return token.AssignMul, nil
	}
	return token.Mul, nil
}

func scanPlus(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '@': // +@
		s.next()
		return token.UnaryPlus, nil
	case '=': // +=
		s.next()
		return token.AssignPlus, nil
	}
	if !s.nospace {
		if '0' <= ch && ch <= '9' {
			s.next()
			if ch == '0' {
				return scanZero(s)
			}
			return scanNonZero(s)
		}
	}
	return token.Plus, nil
}

func scanMinus(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '@': // -@
		s.next()
		return token.UnaryMinus, nil
	case '=': // -=
		s.next()
		return token.AssignMinus, nil
	}
	if !s.nospace {
		if '0' <= ch && ch <= '9' {
			s.next()
			if ch == '0' {
				return scanZero(s)
			}
			return scanNonZero(s)
		}
	}
	return token.Minus, nil
}

func scanDiv(s *Scanner) (token.Token, []byte) {
	if s.char == '=' { // /=
		s.next()
		return token.AssignDiv, nil
	}
	return token.Div, nil
}

func scanZero(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '.':
		s.next()
		return scanFloatDecimal(s)
	case 'd', 'D':
		s.next()
		return scanInt(s)
	case 'b', 'B':
		s.next()
		return scanBinInt(s)
	case '_', 'o', 'O':
		s.next()
		return scanOctInt(s)
	case 'x', 'X':
		s.next()
		return scanHexInt(s)
	}
	return token.DecimalInteger, s.src[s.begin:s.offset]
}

func scanBinInt(s *Scanner) (token.Token, []byte) {
	for s.char == '0' || s.char == '1' || s.char == '_' {
		s.next()
	}
	return token.BinaryInteger, s.src[s.begin:s.offset]
}

func scanOctInt(s *Scanner) (token.Token, []byte) {
	for token.IsOctadecimal(s.char) || s.char == '_' {
		s.next()
	}
	return token.OctadecimalInteger, s.src[s.begin:s.offset]
}

func scanHexInt(s *Scanner) (token.Token, []byte) {
	for token.IsHexadecimal(s.char) || s.char == '_' {
		s.next()
	}
	return token.HexadecimalInteger, s.src[s.begin:s.offset]
}

func scanInt(s *Scanner) (token.Token, []byte) {
	for token.IsDecimal(s.char) || s.char == '_' {
		s.next()
	}
	return token.DecimalInteger, s.src[s.begin:s.offset]
}

func scanNonZero(s *Scanner) (token.Token, []byte) {
	for token.IsDecimal(s.char) || s.char == '_' {
		s.next()
	}
	if s.char == '.' {
		s.next()
		return scanFloatDecimal(s)
	}
	return token.DecimalInteger, s.src[s.begin:s.offset]
}

func scanFloatDecimal(s *Scanner) (token.Token, []byte) {
	for token.IsDecimal(s.char) || s.char == '_' {
		s.next()
	}
	if s.char == 'e' || s.char == 'E' {
		s.next()
		for token.IsDecimal(s.char) || s.char == '_' {
			s.next()
		}
	}
	return token.Float, s.src[s.begin:s.offset]
}

func scanLt(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '=': // <=
		s.next()
		if s.char == '>' { // <=>
			s.next()
			return token.Compare, nil
		}
		return token.LtEq, nil
	case '<': // <<
		s.next()
		if s.char == '=' { // <<=
			s.next()
			return token.AssignLShift, nil
		}
		return token.LShift, nil
	}
	return token.Lt, nil
}

func scanEq(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '=': // ==
		s.next()
		if s.char == '=' { // ===
			s.next()
			return token.Eql, nil
		}
		return token.Eq, nil
	case '~': // =~
		s.next()
		return token.Match, nil
	}
	return token.Assign, nil
}

func scanGt(s *Scanner) (token.Token, []byte) {
	ch := s.char
	switch ch {
	case '=': // >=
		s.next()
		return token.GtEq, nil
	case '>': // >>
		s.next()
		if s.char == '=' { // >>=
			s.next()
			return token.AssignRShift, nil
		}
		return token.RShift, nil
	}
	return token.Gt, nil
}

func scanBracket(s *Scanner) (token.Token, []byte) {
	if s.char == ']' {
		s.next()
		if s.char == '=' {
			s.next()
			return token.ElementSet, nil
		}
		return token.ElementRef, nil
	}
	return token.LBracket, nil
}

func scanXor(s *Scanner) (token.Token, []byte) {
	if s.char == '=' { // ^=
		s.next()
		return token.AssignXor, nil
	}
	return token.Xor, nil
}

func scanOr(s *Scanner) (token.Token, []byte) {
	if s.char == '|' { // ||
		s.next()
		if s.char == '=' { // ||=
			s.next()
			return token.AssignOrOperator, nil
		}
		return token.OrOperator, nil
	}
	if s.char == '=' { // |=
		s.next()
		return token.AssignOr, nil
	}
	return token.Or, nil
}

////////

func (s *Scanner) scanIdent() (int, token.Token, []byte) {
	begin := s.offset
	ch := s.char
	for token.IsLetter(ch) || token.IsDecimal(ch) {
		s.next()
		ch = s.char
	}

	return begin, token.IDENT, s.src[begin:s.offset]
}

func (s *Scanner) skipWhiteSpace() {
	for token.IsWhiteSpace(s.char) {
		s.next()
	}
}
