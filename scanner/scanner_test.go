package scanner

import (
	"bytes"
	"testing"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

func assertScan(t *testing.T, s *Scanner, pos int, token token.Token, literal []byte) bool {
	gp, gt, gl := s.Scan()
	if gp != pos || gt != token || !bytes.Equal(gl, literal) {
		debug.Printf("assert: pos=%v (want=%v), token=%v (want=%v), literal=%v (want=%v)", gp, pos, gt, token, gl, literal)
		t.Errorf("scan: pos=%v (want=%v), token=%v (want=%v), literal=%v (want=%v)", gp, pos, gt, token, gl, literal)
		return true
	}
	return false
}

var rules = map[string]func(t *testing.T, s *Scanner){
	// new line
	"\n": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.NewLine, nil)
	},
	"\r\n": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 1, token.NewLine, nil)
	},
	// white spaces
	"": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.EOF, nil)
	},
	" ": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 1, token.EOF, nil)
	},
	" \t": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 2, token.EOF, nil)
	},
	// comment
	" #": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 2, token.EOF, nil)
	},
	"a#b": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IDENT, []byte("a"))
		assertScan(t, s, 3, token.EOF, nil)
	},
	// operators
	"!": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Not, nil)
	},
	"!=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.NotEqual, nil)
	},
	"!~": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.NotMatch, nil)
	},
	"&&": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AndOperator, nil)
	},
	"||": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.OrOperator, nil)
	},
	// operator methods
	"^": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Xor, nil)
	},
	"&": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Amp, nil)
	},
	"|": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Or, nil)
	},
	"<=>": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Compare, nil)
	},
	"==": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Eq, nil)
	},
	"===": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Eql, nil)
	},
	"=~": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Match, nil)
	},
	">": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Gt, nil)
	},
	">=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.GtEq, nil)
	},
	"<": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Lt, nil)
	},
	"<=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.LtEq, nil)
	},
	"<<": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.LShift, nil)
	},
	">>": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.RShift, nil)
	},
	"+": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Plus, nil)
	},
	"-": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Minus, nil)
	},
	"*": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Mul, nil)
	},
	"/": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Div, nil)
	},
	"%": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Mod, nil)
	},
	"**": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Pow, nil)
	},
	"~": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Invert, nil)
	},
	"+@": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.UnaryPlus, nil)
	},
	"-@": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.UnaryMinus, nil)
	},
	"[]": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.ElementRef, nil)
	},
	"[]=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.ElementSet, nil)
	},
	// assign operators
	"=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Assign, nil)
	},
	"&&=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignAndOperator, nil)
	},
	"||=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignOrOperator, nil)
	},
	"^=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignXor, nil)
	},
	"&=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignAnd, nil)
	},
	"|=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignOr, nil)
	},
	"<<=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignLShift, nil)
	},
	">>=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignRShift, nil)
	},
	"+=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignPlus, nil)
	},
	"-=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignMinus, nil)
	},
	"*=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignMul, nil)
	},
	"/=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignDiv, nil)
	},
	"%=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignMod, nil)
	},
	"**=": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.AssignPow, nil)
	},
	// numeric literals
	"+1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("+1"))
	},
	"-1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("-1"))
	},
	"0": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("0"))
	},
	"123_456_789_0": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("123_456_789_0"))
	},
	"0d1234567890": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("0d1234567890"))
	},
	"0b0101": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.BinaryInteger, []byte("0b0101"))
	},
	"0B0101": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.BinaryInteger, []byte("0B0101"))
	},
	"0_12345670": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.OctadecimalInteger, []byte("0_12345670"))
	},
	"0o12345670": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.OctadecimalInteger, []byte("0o12345670"))
	},
	"0O12345670": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.OctadecimalInteger, []byte("0O12345670"))
	},
	"0x1234567890abcdef": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HexadecimalInteger, []byte("0x1234567890abcdef"))
	},
	"0X1234567890abcdef": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HexadecimalInteger, []byte("0X1234567890abcdef"))
	},
	"+0x1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HexadecimalInteger, []byte("+0x1"))
	},
	"-0X1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HexadecimalInteger, []byte("-0X1"))
	},
	"0.1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Float, []byte("0.1"))
	},
	"+0.1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Float, []byte("+0.1"))
	},
	"-0.1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Float, []byte("-0.1"))
	},
	"123.0456789e10": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Float, []byte("123.0456789e10"))
	},
	"123.0456789E10": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.Float, []byte("123.0456789E10"))
	},
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
