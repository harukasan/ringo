package scanner

import (
	"bytes"
	"io"

	"github.com/harukasan/ringo/token"
)

func scanDoubleQuotedString(s *Scanner, term byte, head int) (token.Token, []byte) {
	t := token.String
	next, rOffset := decodeEscapes(s, term)
	switch {
	case isInsertPrefix(next):
		t = token.StringPart
		s.pushCtx(stateDoubleQuotedStringIn[term])
	case next == term:
		s.next()
	}
	off := s.offset - rOffset
	if off > s.begin+1 {
		off--
	}
	return t, s.src[s.begin+head : off]
}

var stateDoubleQuotedStringIn = [...]stateScanFunc{
	'"': stateDoubleQuotedStringInFunc('"'),
	')': stateDoubleQuotedStringInFunc(')'),
	']': stateDoubleQuotedStringInFunc(']'),
	'>': stateDoubleQuotedStringInFunc('>'),
}

func stateDoubleQuotedStringInFunc(term byte) stateScanFunc {
	return func(s *Scanner) (int, token.Token, []byte) {
		if s.char == '#' {
			p, t, lit := scanInsert(s)
			if t != token.Continue {
				return p, t, lit
			}
		}
		s.begin = s.offset
		next, nEscape := decodeEscapes(s, '"')
		if next == '@' || next == '$' || next == '{' {
			return s.begin, token.StringPart, s.src[s.begin : s.offset-nEscape]
		}
		s.next()
		s.popCtx()
		return s.begin, token.String, s.src[s.begin : s.offset-nEscape-1]
	}
}

func scanInsert(s *Scanner) (int, token.Token, []byte) {
	s.next()
	s.begin = s.offset
	c := s.char
	s.next()
	switch c {
	case '@':
		t, lit := scanAt(s)
		return s.begin - 1, t, lit
	case '$':
		t, lit := scanGlobalVar(s)
		return s.begin - 1, t, lit
	case '{':
		s.pushCtx(stateInsertStmts)
		return s.begin - 1, token.InsertBegin, nil
	}
	return 0, token.Continue, nil
}

func replace(s *Scanner, c byte, offset int) {
	s.src[s.offset-offset] = c
	s.next()
}

func decodeEscapes(s *Scanner, term byte) (byte, int) {
	var skip int
	for s.char != term && s.err == nil {
		switch s.char {
		case '#':
			next := s.peek(2)[1]
			if isInsertPrefix(next) {
				return next, skip
			}
		case '\\':
			skip = decodeEscape(s, skip)
		default:
			replace(s, s.char, skip)
		}
	}
	return s.char, skip
}

func decodeEscape(s *Scanner, skip int) int {
	skip++
	s.next()
	if v := escapes[s.char]; v != 0 {
		replace(s, v, skip)
		return skip
	}
	var n int
	c := s.char
	switch c {
	case '\n':
		n = 2
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, c = decodeOctalEsc(s)
	case 'x':
		skip++
		s.next()
		n, c = decodeHexEsc(s)
	case 'C':
		skip++
		s.next()
		if s.char != '-' {
			s.failf("invalid escape")
			return skip
		}
		skip++
		s.next()
		c = decodeCtrlEsc(s.char)
	case 'c':
		skip++
		s.next()
		c = decodeCtrlEsc(s.char)
	}
	for i := 1; i < n; i++ {
		skip++
		s.next()
	}
	replace(s, c, skip)
	return skip
}

func decodeCtrlEsc(c byte) byte {
	if c == '?' {
		return 0x7f
	}
	if v := escapes[c]; v != 0 {
		c = v
	}
	return c & 0x9f
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func decodeOctalEsc(s *Scanner) (n int, v byte) {
	m := min(3, len(s.src)-s.offset)
	for n = 0; n < m; n++ {
		c := s.src[s.offset+n]
		if !token.IsOctadecimal(c) {
			break
		}
		v = v*8 + (c - '0')
	}
	if n == 0 {
		s.failf("invalid octal escape")
	}
	return
}

func decodeHexEsc(s *Scanner) (n int, v byte) {
	m := min(2, len(s.src)-s.offset)
	for n = 0; n < m; n++ {
		c := s.src[s.offset+n]
		var d byte
		switch {
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'A' <= c && c <= 'F':
			d = c - 'A' + 10
		case 'a' <= c && c <= 'f':
			d = c - 'a' + 10
		default:
			if n == 0 {
				s.failf("invalid hex escape")
			}
			return
		}
		v = v*16 + d
	}
	return
}

func stateInsertStmts(s *Scanner) (pos int, t token.Token, literal []byte) {
	s.begin = s.offset
	if s.char == '}' {
		s.next()
		s.popCtx()
		return s.begin, token.InsertEnd, nil
	}
	return stateCompStmts(s)
}

func scanSingleQuotedString(s *Scanner, term byte, head int) (token.Token, []byte) {
	var skip int
	for s.char != term && s.err == nil {
		if s.char == '\\' {
			s.next()
			if s.char == '\\' || s.char == term {
				skip++
			}
		}
		replace(s, s.char, skip)
	}
	return token.String, s.src[s.begin+head : s.offset-skip]
}

// TODO: cyclomatic complexity >= 12
func scanHeredocBegin(s *Scanner) (token.Token, []byte) {
	indent := false
	termBegin := s.offset
	switch c := s.char; {
	case token.IsLetter(c) || c == '_':
		s.next()
	case token.IsDecimal(c):
		if s.ctx.nospace {
			return token.Continue, nil
		}
		s.next()
	case c == '-':
		if s.ctx.nospace {
			return token.Continue, nil
		}
		indent = true
		s.next()
		termBegin = s.offset
	default:
		return token.Continue, nil
	}

	var quote byte
	if s.char == '\'' {
		quote = s.char
		s.next()
		termBegin = s.offset
	}
	for token.IsIdent(s.char) {
		s.next()
	}
	term := s.src[termBegin:s.offset]
	if quote != 0 {
		if s.char != quote {
			s.failf("invalid heredoc identifier")
		}
		s.next()
	}
	s.pushCtx(stateHeredocFirstLine(term, indent))
	return token.HeredocBegin, s.src[s.begin:s.offset]
}

func stateHeredocFirstLine(term []byte, indent bool) stateScanFunc {
	return func(s *Scanner) (int, token.Token, []byte) {
		if s.char == '\n' {
			begin := s.offset
			s.next()
			s.popCtx()
			s.pushCtx(stateInHeredoc(term, indent))
			return begin, token.NewLine, nil
		}
		return stateCompStmts(s)
	}
}

func isHeredocEndTerm(s *Scanner, term []byte, indent bool) (bool, int) {
	if s.src[s.offset-1] != '\n' {
		return false, 0
	}
	tOff := s.offset
	if indent {
		for token.IsWhiteSpace(s.char) {
			s.next()
		}
	}
	if bytes.HasPrefix(s.src[s.offset:], term) {
		s.skip(len(term))
		if s.char == '\n' || s.err == io.EOF {
			if s.char == '\n' {
				s.next()
			}
			return true, tOff
		}
	}
	return false, 0
}

func stateInHeredoc(term []byte, indent bool) stateScanFunc {
	return func(s *Scanner) (int, token.Token, []byte) {
		if s.char == '#' {
			p, t, lit := scanInsert(s)
			if t != token.Continue {
				return p, t, lit
			}
		}

		s.begin = s.offset
		for s.err == nil {
			if isEnd, off := isHeredocEndTerm(s, term, indent); isEnd {
				s.popCtx()
				return s.begin, token.HeredocEnd, s.src[s.begin:off]
			}
			next, skip := decodeEscapes(s, '\n')
			if next == '@' || next == '$' || next == '{' {
				return s.begin, token.HeredocPart, s.src[s.begin : s.offset-skip]
			}
			s.next()
		}
		return s.begin, token.Illegal, nil
	}
}
