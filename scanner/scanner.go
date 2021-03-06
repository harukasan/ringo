package scanner

import (
	"bytes"
	"fmt"
	"io"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

// Scanner implements a scanner for Ruby lex.
type Scanner struct {
	src []byte // source buffer
	err error  //

	char   byte // current read character
	offset int  // current offset
	begin  int  // offset of begin of the token

	ctx *scannerCtx // scanner context
}

type scannerCtx struct {
	nospace   bool          // whether the previous is not a space
	stateScan stateScanFunc // scanner func for the special state
	parent    *scannerCtx   // parent context
}

// scanning function for special state
type stateScanFunc func(s *Scanner) (pos int, t token.Token, literal []byte)

// New returns a initiazlied scanner to scan script source src.
func New(src []byte) *Scanner {
	s := &Scanner{
		src:    src,
		offset: -1,
		ctx: &scannerCtx{
			stateScan: stateCompStmts,
		},
	}
	s.next()
	return s
}

// NewString returns a initiazlied scanner with given string s.
func NewString(s string) *Scanner {
	return New([]byte(s))
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

func (s *Scanner) skip(n int) {
	for i := 0; i < n; i++ {
		s.next()
	}
}

func (s *Scanner) peek(n int) []byte {
	if len(s.src) < s.offset+n {
		return nil
	}
	return s.src[s.offset : s.offset+n]
}

func (s *Scanner) failf(format string, v ...interface{}) {
	if s.err != nil && s.err != io.EOF {
		return
	}
	s.err = fmt.Errorf(format, v...)
	debug.Printf("failf: %v", s.err)
}

func (s *Scanner) pushCtx(state stateScanFunc) {
	s.ctx = &scannerCtx{
		stateScan: state,
		parent:    s.ctx,
	}
	debug.Printf("push: -> %#v", s.ctx)
}

func (s *Scanner) popCtx() {
	s.ctx = s.ctx.parent
	debug.Printf("pop: -> %#v", s.ctx)
}

// Scan reads and returns a parsed token position, type, and its literal.
func (s *Scanner) Scan() (pos int, t token.Token, literal []byte) {
	t = token.Continue
	for t == token.Continue {
		if s.err == io.EOF {
			pos, t, literal = s.offset, token.EOF, nil
			break
		}
		pos, t, literal = s.ctx.stateScan(s)
	}
	return
}

func stateCompStmts(s *Scanner) (pos int, t token.Token, literal []byte) {
	s.begin = s.offset
	if scan := scanners[s.char]; scan != nil {
		s.next()
		t, literal = scan(s)
		if t != token.Continue && t != token.NewLine {
			s.ctx.nospace = true
		}
		return s.begin, t, literal
	}
	return s.offset, token.Illegal, nil
}

func (s *Scanner) skipLine() {
	for s.err == nil && s.char != '\n' {
		s.next()
	}
}

// scanFunc implements a scanner that returns a token type and its literal.
type scanFunc func(s *Scanner) (token.Token, []byte)

func scanOne(tk token.Token) scanFunc {
	return func(s *Scanner) (token.Token, []byte) {
		return tk, nil
	}
}

/*
The scanner for speicfic token is picked up by the first letter of literal of
token. Note that the scanner for 0-9, A-Z and a-z is set by below init func.
*/

/* set scanners for 0-9, A-Z, and a-z. */
var scanners [127]scanFunc

var escapes = [127]byte{
	'n': 0x0a,
	't': 0x09,
	'r': 0x0d,
	'f': 0x0c,
	'v': 0x0b,
	'a': 0x07,
	'e': 0x1b,
	'b': 0x08,
	's': 0x20,
}

func init() {
	scanners = [...]scanFunc{
		0x09: skipWhiteSpaces,
		'\n': scanNewLine,
		0x0b: skipWhiteSpaces,
		0x0c: skipWhiteSpaces,
		0x0d: skipWhiteSpaces,
		'!':  scanNot,
		' ':  skipWhiteSpaces,
		'"':  scanDoubleQuote,
		'#':  scanComment,
		'$':  scanDollar,
		'%':  scanPercent,
		'&':  scanAmp,
		'\'': scanSingleQuote,
		'(':  scanOne(token.LParen),
		')':  scanOne(token.RParen),
		'*':  scanAsterisk,
		'+':  scanPlus,
		',':  scanOne(token.Comma),
		'-':  scanMinus,
		'.':  scanDot,
		'/':  scanDiv,
		':':  scanColon,
		';':  scanNewLine,
		'<':  scanLt,
		'=':  scanEq,
		'>':  scanGt,
		'?':  scanOne(token.Question),
		'@':  scanAt,
		'[':  scanBracket,
		'\\': scanEscSeq,
		']':  scanOne(token.RBracket),
		'^':  scanXor,
		'_':  scanUnderscore,
		'{':  scanOne(token.LBrace),
		'|':  scanOr,
		'}':  scanOne(token.RBrace),
		'~':  scanOne(token.Invert),
	}
	/* set scanners for 0-9, A-Z, and a-z. */
	scanners['0'] = scanZero
	for i := '1'; i <= '9'; i++ {
		scanners[i] = scanNonZero
	}
	for i := 'A'; i <= 'Z'; i++ {
		scanners[i] = scanUppercase
	}
	for i := 'a'; i <= 'z'; i++ {
		scanners[i] = scanLowercase
	}
}

func skipWhiteSpaces(s *Scanner) (token.Token, []byte) {
	s.ctx.nospace = false
	for token.IsWhiteSpace(s.char) {
		s.next()
	}
	return token.Continue, nil
}

func scanNewLine(s *Scanner) (token.Token, []byte) {
	s.ctx.nospace = false
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

func isInsertPrefix(c byte) bool {
	return c == '@' || c == '$' || c == '{'
}

func scanDoubleQuote(s *Scanner) (token.Token, []byte) {
	return scanDoubleQuotedString(s, '"', 1)
}

func scanComment(s *Scanner) (token.Token, []byte) {
	s.skipLine()
	return token.Continue, nil
}

func scanDollar(s *Scanner) (token.Token, []byte) {
	return scanGlobalVar(s)
}

func scanGlobalVar(s *Scanner) (token.Token, []byte) {
	if !token.IsIdentStart(s.char) {
		return token.Illegal, s.src[s.begin:s.offset]
	}
	s.next()
	for token.IsIdent(s.char) {
		s.next()
	}
	return token.IdentGlobalVar, s.src[s.begin:s.offset]
}

func closeBracket(c byte) byte {
	switch c {
	case '{':
		return '}'
	case '(':
		return ')'
	case '[':
		return ']'
	case '<':
		return '>'
	}
	return c
}

func scanPercent(s *Scanner) (token.Token, []byte) {
	switch s.char {
	case 'Q':
		p := s.peek(2)
		if p != nil && !token.IsAlnum(p[1]) { // %Q!...!
			s.next()
			term := closeBracket(s.char)
			s.next()
			return scanDoubleQuotedString(s, term, 3)
		}
	case 'q':
		p := s.peek(2)
		if p != nil && !token.IsAlnum(p[1]) { // %q!...!
			s.next()
			term := closeBracket(s.char)
			s.next()
			return scanSingleQuotedString(s, term, 3)
		}
	case '=': // %=
		s.next()
		return token.AssignMod, nil
	}
	if s.err == nil && !token.IsAlnum(s.char) { // %!...!
		term := closeBracket(s.char)
		s.next()
		return scanDoubleQuotedString(s, term, 2)
	}
	return token.Mod, nil
}

func scanAmp(s *Scanner) (token.Token, []byte) {
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

func scanSingleQuote(s *Scanner) (token.Token, []byte) {
	return scanSingleQuotedString(s, '\'', 1)
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
	if !s.ctx.nospace {
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
	if !s.ctx.nospace {
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

func scanDot(s *Scanner) (token.Token, []byte) {
	if s.char == '.' {
		s.next()
		if s.char == '.' {
			return token.Dot3, nil
		}
		return token.Dot2, nil
	}
	return token.Dot, nil
}

func scanDiv(s *Scanner) (token.Token, []byte) {
	if s.char == '=' { // /=
		s.next()
		return token.AssignDiv, nil
	}
	return token.Div, nil
}

func scanColon(s *Scanner) (token.Token, []byte) {
	if s.char == ':' {
		s.next()
		return token.Colon2, nil
	}
	return token.Colon, nil
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
		switch s.char {
		case '=': // <<=
			s.next()
			return token.AssignLShift, nil
		case '-': // <<-
			t, l := scanHeredocBegin(s)
			if t != token.Continue {
				return t, l
			}
		default:
			t, l := scanHeredocBegin(s)
			if t != token.Continue {
				return t, l
			}
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
	case '>': // =>
		s.next()
		return token.Arrow, nil
	case 'b': // =begin
		if s.offset < 2 || s.src[s.offset-2] == '\n' {
			p := s.peek(6)
			if p != nil && bytes.HasPrefix(p, []byte("begin")) {
				if token.IsWhiteSpace(p[5]) || p[5] == '\n' {
					skipMultiLineComment(s)
					return token.Continue, nil
				}
			}
		}
	case '~': // =~
		s.next()
		return token.Match, nil
	}
	return token.Assign, nil
}

func skipMultiLineComment(s *Scanner) {
	for {
		if s.char == '=' {
			s.next()
			if bytes.HasPrefix(s.src[s.offset:], []byte("end")) {
				s.skip(3)
				if token.IsWhiteSpace(s.char) {
					s.skipLine()
				}
				if s.char == '\n' {
					s.next()
					return
				}
				if s.err == io.EOF {
					return
				}
			}
		}
		s.next()
		if s.err != nil {
			if s.err == io.EOF {
				s.failf("multi-line comment must be closed")
			}
			return
		}
	}
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

func scanAt(s *Scanner) (token.Token, []byte) {
	t := token.IdentInstanceVar
	if s.char == '@' {
		t = token.IdentClassVar
		s.next()
	}
	if !token.IsIdentStart(s.char) {
		return token.Illegal, s.src[s.begin:s.offset]
	}
	s.next()
	for token.IsIdent(s.char) {
		s.next()
	}
	return t, s.src[s.begin:s.offset]
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

func scanEscSeq(s *Scanner) (token.Token, []byte) {
	switch s.char {
	case '\r':
		s.next()
		if s.char == '\n' {
			s.next()
			s.ctx.nospace = false
			return token.Continue, nil
		}
	case '\n':
		s.next()
		s.ctx.nospace = false
		return token.Continue, nil
	}
	s.failf("escape character must be at end of line")
	return token.Illegal, nil
}

func scanXor(s *Scanner) (token.Token, []byte) {
	if s.char == '=' { // ^=
		s.next()
		return token.AssignXor, nil
	}
	return token.Xor, nil
}

func scanUnderscore(s *Scanner) (token.Token, []byte) {
	if s.offset < 2 || s.src[s.offset-2] == '\n' {
		if bytes.HasPrefix(s.src[s.offset:], []byte("_END__")) {
			s.skip(6) // skip "_END__"
			if s.char == '\r' {
				if p := s.peek(2); p != nil && p[1] == '\n' {
					s.next()
				}
			}
			if s.char == '\n' || s.err == io.EOF {
				s.err = io.EOF
				return token.EOF, nil
			}
		}
	}
	return scanLowercase(s)
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

func scanUppercase(s *Scanner) (token.Token, []byte) {
	for token.IsIdent(s.char) {
		s.next()
	}
	lit := s.src[s.begin:s.offset]
	if t := token.KeywordToken(lit); t != token.None {
		return t, nil
	}
	return token.IdentConst, lit
}

func scanLowercase(s *Scanner) (token.Token, []byte) {
	t := token.IdentLocalVar
	for token.IsIdent(s.char) {
		s.next()
	}
	if s.char == '?' || s.char == '!' || s.char == '=' {
		t = token.IdentLocalMethod
		s.next()
	}
	lit := s.src[s.begin:s.offset]
	if kt := token.KeywordToken(lit); kt != token.None {
		return kt, nil
	}
	return t, lit
}
