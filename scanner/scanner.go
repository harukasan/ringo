package scanner

import (
	"fmt"
	"io"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

// Scanner implements a scanner for Ruby lex.
type Scanner struct {
	src    []byte
	offset int
	err    error
	char   byte

	lastScan token.Token
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
		s.lastScan = t
		return
	case token.IsWhiteSpace(ch):
		s.lastScan = token.None
		s.skipWhiteSpace()
		goto StartScan
	}
	if s.err == io.EOF {
		pos, t, literal = s.offset, token.EOF, nil
		s.lastScan = t
		return
	}
	s.next()
	if scan := tokenScanners[ch]; scan != nil {
		pos, t, literal = scan(s)
		pos--
		if t == token.Continue {
			goto StartScan
		}
		s.lastScan = t
		return pos, t, literal
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

type scanFunc func(s *Scanner) (int, token.Token, []byte)

var tokenScanners = [127]scanFunc{
	'\n': scanOne(token.NewLine),
	'!':  scanNot,
	'#':  scanComment,
	'%':  scanMod,
	'&':  scanAnd,
	'(':  scanOne(token.LParen),
	')':  scanOne(token.RParen),
	'*':  scanAsterisk,
	'+':  scanPlus,
	'-':  scanMinus,
	'/':  scanDiv,
	'0':  scanZero,
	'1':  scanNum,
	'2':  scanNum,
	'3':  scanNum,
	'4':  scanNum,
	'5':  scanNum,
	'6':  scanNum,
	'7':  scanNum,
	'8':  scanNum,
	'9':  scanNum,
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
	return func(s *Scanner) (int, token.Token, []byte) {
		return s.offset, tk, nil
	}
}

func scanNot(s *Scanner) (int, token.Token, []byte) {
	ch := s.char
	offset := s.offset
	switch ch {
	case '=':
		s.next()
		return offset, token.NotEqual, nil
	case '~':
		s.next()
		return offset, token.NotMatch, nil
	}
	return offset, token.Not, nil
}

func scanComment(s *Scanner) (int, token.Token, []byte) {
	s.skipLine()
	return 0, token.Continue, nil
}

func scanMod(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '=' { // %=
		s.next()
		return offset, token.AssignMod, nil
	}
	return offset, token.Mod, nil
}

func scanAnd(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '&' { // &&
		s.next()
		if s.char == '=' { // &&=
			s.next()
			return offset, token.AssignAndOperator, nil
		}
		return offset, token.AndOperator, nil
	}
	if s.char == '=' { // &=
		s.next()
		return offset, token.AssignAnd, nil
	}
	return offset, token.Amp, nil
}

func scanAsterisk(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '*' { // **
		s.next()
		if s.char == '=' { // **=
			s.next()
			return offset, token.AssignPow, nil
		}
		return offset, token.Pow, nil
	}
	if s.char == '=' { // *=
		s.next()
		return offset, token.AssignMul, nil
	}
	return offset, token.Mul, nil
}

func scanPlus(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case '@': // +@
		s.next()
		return offset, token.UnaryPlus, nil
	case '=': // +=
		s.next()
		return offset, token.AssignPlus, nil
	}
	if s.lastScan != token.IDENT && token.IsNonZeroDecimal(ch) {
		return scanDecimal(s, offset)
	}
	return offset, token.Plus, nil
}

func scanMinus(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case '@': // -@
		s.next()
		return offset, token.UnaryMinus, nil
	case '=': // -=
		s.next()
		return offset, token.AssignMinus, nil
	}
	if s.lastScan != token.IDENT && token.IsNonZeroDecimal(ch) {
		return scanDecimal(s, offset)
	}
	return offset, token.Minus, nil
}

func scanDiv(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '=' { // /=
		s.next()
		return offset, token.AssignDiv, nil
	}
	return offset, token.Div, nil
}

func scanZero(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case 'd':
		s.next()
		return scanDecimal(s, offset)
	case 'D':
		s.next()
		return scanDecimal(s, offset)
	case 'b':
		s.next()
		return scanBinNum(s, offset)
	case 'B':
		s.next()
		return scanBinNum(s, offset)
	case '_':
		s.next()
		return scanOctNum(s, offset)
	case 'o':
		s.next()
		return scanOctNum(s, offset)
	case 'O':
		s.next()
		return scanOctNum(s, offset)
	case 'x':
		s.next()
		return scanHexNum(s, offset)
	case 'X':
		s.next()
		return scanHexNum(s, offset)
	}
	return offset, token.DecimalInteger, s.src[offset-1 : offset]
}

func scanBinNum(s *Scanner, offset int) (int, token.Token, []byte) {
	for s.char == '0' || s.char == '1' || s.char == '_' {
		s.next()
	}
	return offset, token.BinaryInteger, s.src[offset-1 : s.offset]
}

func scanOctNum(s *Scanner, offset int) (int, token.Token, []byte) {
	for token.IsOctadecimal(s.char) || s.char == '_' {
		s.next()
	}
	return offset, token.OctadecimalInteger, s.src[offset-1 : s.offset]
}

func scanHexNum(s *Scanner, offset int) (int, token.Token, []byte) {
	for token.IsHexadecimal(s.char) || s.char == '_' {
		s.next()
	}
	return offset, token.HexadecimalInteger, s.src[offset-1 : s.offset]
}

func scanDecimal(s *Scanner, offset int) (int, token.Token, []byte) {
	for token.IsDecimal(s.char) || s.char == '_' {
		debug.Printf("offset=%v", s.offset)
		s.next()
	}
	return offset, token.DecimalInteger, s.src[offset-1 : s.offset]
}

func scanNum(s *Scanner) (int, token.Token, []byte) {
	return scanDecimal(s, s.offset)
}

func scanLt(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case '=': // <=
		s.next()
		if s.char == '>' { // <=>
			s.next()
			return offset, token.Compare, nil
		}
		return offset, token.LtEq, nil
	case '<': // <<
		s.next()
		if s.char == '=' { // <<=
			s.next()
			return offset, token.AssignLShift, nil
		}
		return offset, token.LShift, nil
	}
	return offset, token.Lt, nil
}

func scanEq(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case '=': // ==
		s.next()
		if s.char == '=' { // ===
			s.next()
			return offset, token.Eql, nil
		}
		return offset, token.Eq, nil
	case '~': // =~
		s.next()
		return offset, token.Match, nil
	}
	return offset, token.Assign, nil
}

func scanGt(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	ch := s.char
	switch ch {
	case '=': // >=
		s.next()
		return offset, token.GtEq, nil
	case '>': // >>
		s.next()
		if s.char == '=' { // >>=
			s.next()
			return offset, token.AssignRShift, nil
		}
		return offset, token.RShift, nil
	}
	return offset, token.Gt, nil
}

func scanBracket(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == ']' {
		s.next()
		if s.char == '=' {
			s.next()
			return offset, token.ElementSet, nil
		}
		return offset, token.ElementRef, nil
	}
	return offset, token.LBracket, nil
}

func scanXor(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '=' { // ^=
		s.next()
		return offset, token.AssignXor, nil
	}
	return offset, token.Xor, nil
}

func scanOr(s *Scanner) (int, token.Token, []byte) {
	offset := s.offset
	if s.char == '|' { // ||
		s.next()
		if s.char == '=' { // ||=
			s.next()
			return offset, token.AssignOrOperator, nil
		}
		return offset, token.OrOperator, nil
	}
	if s.char == '=' { // |=
		s.next()
		return offset, token.AssignOr, nil
	}
	return offset, token.Or, nil
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
