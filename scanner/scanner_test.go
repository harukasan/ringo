package scanner

import (
	"bytes"
	"testing"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

func assertScanToken(offset int, token token.Token, literal []byte) func(t *testing.T, s *Scanner) {
	return func(t *testing.T, s *Scanner) {
		assertScan(t, s, offset, token, literal)
	}
}

func assertScan(t *testing.T, s *Scanner, pos int, token token.Token, literal []byte) bool {
	gp, gt, gl := s.Scan()
	if gp != pos || gt != token || !bytes.Equal(gl, literal) {
		debug.Printf("assert: src=%v pos=%v (want=%v), token=%v (want=%v), literal=%v (want=%v)", string(s.src), gp, pos, gt, token, gl, literal)
		t.Errorf("scan: src=%v pos=%v (want=%v), token=%v (want=%v), literal=%v (want=%v)", string(s.src), gp, pos, gt, token, gl, literal)
		return true
	}
	return false
}

var rules = map[string]func(t *testing.T, s *Scanner){
	// new line
	"\n":   assertScanToken(0, token.NewLine, nil),
	"\r\n": assertScanToken(1, token.NewLine, nil),

	// white spaces
	"":    assertScanToken(0, token.EOF, nil),
	" ":   assertScanToken(1, token.EOF, nil),
	" \t": assertScanToken(2, token.EOF, nil),

	// comment
	" #": assertScanToken(2, token.EOF, nil),
	"a#b": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("a"))
		assertScan(t, s, 3, token.EOF, nil)
	},

	// operators
	"!":  assertScanToken(0, token.Not, nil),
	"!=": assertScanToken(0, token.NotEqual, nil),
	"!~": assertScanToken(0, token.NotMatch, nil),
	"&&": assertScanToken(0, token.AndOperator, nil),
	"||": assertScanToken(0, token.OrOperator, nil),

	// operator methods
	"^":   assertScanToken(0, token.Xor, nil),
	"&":   assertScanToken(0, token.Amp, nil),
	"|":   assertScanToken(0, token.Or, nil),
	"<=>": assertScanToken(0, token.Compare, nil),
	"==":  assertScanToken(0, token.Eq, nil),
	"===": assertScanToken(0, token.Eql, nil),
	"=~":  assertScanToken(0, token.Match, nil),
	">":   assertScanToken(0, token.Gt, nil),
	">=":  assertScanToken(0, token.GtEq, nil),
	"<":   assertScanToken(0, token.Lt, nil),
	"<=":  assertScanToken(0, token.LtEq, nil),
	"<<":  assertScanToken(0, token.LShift, nil),
	">>":  assertScanToken(0, token.RShift, nil),
	"+":   assertScanToken(0, token.Plus, nil),
	"-":   assertScanToken(0, token.Minus, nil),
	"*":   assertScanToken(0, token.Mul, nil),
	"/":   assertScanToken(0, token.Div, nil),
	"%":   assertScanToken(0, token.Mod, nil),
	"**":  assertScanToken(0, token.Pow, nil),
	"~":   assertScanToken(0, token.Invert, nil),
	"+@":  assertScanToken(0, token.UnaryPlus, nil),
	"-@":  assertScanToken(0, token.UnaryMinus, nil),
	"[]":  assertScanToken(0, token.ElementRef, nil),
	"[]=": assertScanToken(0, token.ElementSet, nil),

	// assign operators
	"=":   assertScanToken(0, token.Assign, nil),
	"&&=": assertScanToken(0, token.AssignAndOperator, nil),
	"||=": assertScanToken(0, token.AssignOrOperator, nil),
	"^=":  assertScanToken(0, token.AssignXor, nil),
	"&=":  assertScanToken(0, token.AssignAnd, nil),
	"|=":  assertScanToken(0, token.AssignOr, nil),
	"<<=": assertScanToken(0, token.AssignLShift, nil),
	">>=": assertScanToken(0, token.AssignRShift, nil),
	"+=":  assertScanToken(0, token.AssignPlus, nil),
	"-=":  assertScanToken(0, token.AssignMinus, nil),
	"*=":  assertScanToken(0, token.AssignMul, nil),
	"/=":  assertScanToken(0, token.AssignDiv, nil),
	"%=":  assertScanToken(0, token.AssignMod, nil),
	"**=": assertScanToken(0, token.AssignPow, nil),

	// numeric literals
	"+1":                 assertScanToken(0, token.DecimalInteger, []byte("+1")),
	"-1":                 assertScanToken(0, token.DecimalInteger, []byte("-1")),
	"0":                  assertScanToken(0, token.DecimalInteger, []byte("0")),
	"123_456_789_0":      assertScanToken(0, token.DecimalInteger, []byte("123_456_789_0")),
	"0d1234567890":       assertScanToken(0, token.DecimalInteger, []byte("0d1234567890")),
	"0b0101":             assertScanToken(0, token.BinaryInteger, []byte("0b0101")),
	"0B0101":             assertScanToken(0, token.BinaryInteger, []byte("0B0101")),
	"0_12345670":         assertScanToken(0, token.OctadecimalInteger, []byte("0_12345670")),
	"0o12345670":         assertScanToken(0, token.OctadecimalInteger, []byte("0o12345670")),
	"0O12345670":         assertScanToken(0, token.OctadecimalInteger, []byte("0O12345670")),
	"0x1234567890abcdef": assertScanToken(0, token.HexadecimalInteger, []byte("0x1234567890abcdef")),
	"0X1234567890abcdef": assertScanToken(0, token.HexadecimalInteger, []byte("0X1234567890abcdef")),
	"+0x1":               assertScanToken(0, token.HexadecimalInteger, []byte("+0x1")),
	"-0X1":               assertScanToken(0, token.HexadecimalInteger, []byte("-0X1")),
	"0.1":                assertScanToken(0, token.Float, []byte("0.1")),
	"+0.1":               assertScanToken(0, token.Float, []byte("+0.1")),
	"-0.1":               assertScanToken(0, token.Float, []byte("-0.1")),
	"123.0456789e10":     assertScanToken(0, token.Float, []byte("123.0456789e10")),
	"123.0456789E10":     assertScanToken(0, token.Float, []byte("123.0456789E10")),
	"x+1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("x"))
		assertScan(t, s, 1, token.Plus, nil)
		assertScan(t, s, 2, token.DecimalInteger, []byte("1"))
	},
	"x +1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("x"))
		assertScan(t, s, 2, token.DecimalInteger, []byte("+1"))
	},
	"x-1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("x"))
		assertScan(t, s, 1, token.Minus, nil)
		assertScan(t, s, 2, token.DecimalInteger, []byte("1"))
	},
	"x -1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("x"))
		assertScan(t, s, 2, token.DecimalInteger, []byte("-1"))
	},

	// string literals:
	`'a'`:      assertScanToken(0, token.StringPart, []byte(`'a'`)),
	`'\''`:     assertScanToken(0, token.StringPart, []byte(`'\''`)),
	`'\a\\\''`: assertScanToken(0, token.StringPart, []byte(`'\a\\\''`)),

	// ident
	"a": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("a"))
	},
	"a  b": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("a"))
		assertScan(t, s, 3, token.IDENT, []byte("b"))
	},
}

func TestScanner(t *testing.T) {
	debug.Enable = true
	for input, r := range rules {
		debug.Printf("input: %q", input)
		s := New([]byte(input))
		r(t, s)
	}
}
