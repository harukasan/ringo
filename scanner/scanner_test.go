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
		assertScan(t, s, 0, token.IdentLocalVar, []byte("a"))
		assertScan(t, s, 3, token.EOF, nil)
	},

	// delimiters
	"::":  assertScanToken(0, token.Colon2, nil),
	",":   assertScanToken(0, token.Comma, nil),
	"..":  assertScanToken(0, token.Dot2, nil),
	"...": assertScanToken(0, token.Dot3, nil),
	"?":   assertScanToken(0, token.Question, nil),
	":":   assertScanToken(0, token.Colon, nil),
	"=>":  assertScanToken(0, token.Arrow, nil),

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

	"1 << 1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("1"))
		assertScan(t, s, 2, token.LShift, nil)
		assertScan(t, s, 5, token.DecimalInteger, []byte("1"))
	},
	"1<<1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("1"))
		assertScan(t, s, 1, token.LShift, nil)
		assertScan(t, s, 3, token.DecimalInteger, []byte("1"))
	},

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
	"+1\n-1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("+1"))
		assertScan(t, s, 2, token.NewLine, nil)
		assertScan(t, s, 3, token.DecimalInteger, []byte("-1"))
	},
	"x+1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IdentLocalVar, []byte("x"))
		assertScan(t, s, 1, token.Plus, nil)
		assertScan(t, s, 2, token.DecimalInteger, []byte("1"))
	},
	"x +1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IdentLocalVar, []byte("x"))
		assertScan(t, s, 2, token.DecimalInteger, []byte("+1"))
	},
	"x-1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IdentLocalVar, []byte("x"))
		assertScan(t, s, 1, token.Minus, nil)
		assertScan(t, s, 2, token.DecimalInteger, []byte("1"))
	},
	"x -1": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IdentLocalVar, []byte("x"))
		assertScan(t, s, 2, token.DecimalInteger, []byte("-1"))
	},

	// string literals:
	`'a'`:      assertScanToken(0, token.StringPart, []byte(`'a'`)),
	`'\''`:     assertScanToken(0, token.StringPart, []byte(`'\''`)),
	`'\a\\\''`: assertScanToken(0, token.StringPart, []byte(`'\a\\\''`)),

	// heredoc
	"<<TEXT\nabc\n\nTEXT\n": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HeredocBegin, []byte("<<TEXT"))
		assertScan(t, s, 6, token.NewLine, nil)
		assertScan(t, s, 7, token.HeredocPart, []byte("abc\n\n"))
		assertScan(t, s, 16, token.NewLine, nil)
	},
	"<<-TEXT\n  TEXT\n": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.HeredocBegin, []byte("<<-TEXT"))
		assertScan(t, s, 7, token.NewLine, nil)
		assertScan(t, s, 8, token.HeredocPart, nil)
		assertScan(t, s, 14, token.NewLine, nil)
	},
	"1 <<1\n1\n": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.DecimalInteger, []byte("1"))
		assertScan(t, s, 2, token.HeredocBegin, []byte("<<1"))
		assertScan(t, s, 5, token.NewLine, nil)
		assertScan(t, s, 6, token.HeredocPart, nil)
		assertScan(t, s, 7, token.NewLine, nil)
	},

	// ident
	"v":        assertScanToken(0, token.IdentLocalVar, []byte("v")),
	"_":        assertScanToken(0, token.IdentLocalVar, []byte("_")),
	"v?":       assertScanToken(0, token.IdentLocalMethod, []byte("v?")),
	"v!":       assertScanToken(0, token.IdentLocalMethod, []byte("v!")),
	"v=":       assertScanToken(0, token.IdentLocalMethod, []byte("v=")),
	"$v":       assertScanToken(0, token.IdentGlobalVar, []byte("$v")),
	"$v1":      assertScanToken(0, token.IdentGlobalVar, []byte("$v1")),
	"@var1":    assertScanToken(0, token.IdentInstanceVar, []byte("@var1")),
	"@@var1":   assertScanToken(0, token.IdentClassVar, []byte("@@var1")),
	"Constant": assertScanToken(0, token.IdentConst, []byte("Constant")),
	"a  b": func(t *testing.T, s *Scanner) {
		assertScan(t, s, 0, token.IdentLocalVar, []byte("a"))
		assertScan(t, s, 3, token.IdentLocalVar, []byte("b"))
	},

	// keywords
	"__LINE__":     assertScanToken(0, token.KeywordLINE, nil),
	"__ENCODING__": assertScanToken(0, token.KeywordENCODING, nil),
	"__FILE__":     assertScanToken(0, token.KeywordFILE, nil),
	"BEGIN":        assertScanToken(0, token.KeywordBEGIN, nil),
	"END":          assertScanToken(0, token.KeywordEND, nil),
	"alias":        assertScanToken(0, token.KeywordAlias, nil),
	"and":          assertScanToken(0, token.KeywordAnd, nil),
	"begin":        assertScanToken(0, token.KeywordBegin, nil),
	"break":        assertScanToken(0, token.KeywordBreak, nil),
	"case":         assertScanToken(0, token.KeywordCase, nil),
	"class":        assertScanToken(0, token.KeywordClass, nil),
	"def":          assertScanToken(0, token.KeywordDef, nil),
	"defined?":     assertScanToken(0, token.KeywordDefined, nil),
	"do":           assertScanToken(0, token.KeywordDo, nil),
	"else":         assertScanToken(0, token.KeywordElse, nil),
	"elsif":        assertScanToken(0, token.KeywordElsif, nil),
	"end":          assertScanToken(0, token.KeywordEnd, nil),
	"ensure":       assertScanToken(0, token.KeywordEnsure, nil),
	"for":          assertScanToken(0, token.KeywordFor, nil),
	"false":        assertScanToken(0, token.KeywordFalse, nil),
	"if":           assertScanToken(0, token.KeywordIf, nil),
	"in":           assertScanToken(0, token.KeywordIn, nil),
	"module":       assertScanToken(0, token.KeywordModule, nil),
	"next":         assertScanToken(0, token.KeywordNext, nil),
	"nil":          assertScanToken(0, token.KeywordNil, nil),
	"not":          assertScanToken(0, token.KeywordNot, nil),
	"or":           assertScanToken(0, token.KeywordOr, nil),
	"redo":         assertScanToken(0, token.KeywordRedo, nil),
	"rescue":       assertScanToken(0, token.KeywordRescue, nil),
	"retry":        assertScanToken(0, token.KeywordRetry, nil),
	"return":       assertScanToken(0, token.KeywordReturn, nil),
	"self":         assertScanToken(0, token.KeywordSelf, nil),
	"super":        assertScanToken(0, token.KeywordSuper, nil),
	"then":         assertScanToken(0, token.KeywordThen, nil),
	"true":         assertScanToken(0, token.KeywordTrue, nil),
	"undef":        assertScanToken(0, token.KeywordUndef, nil),
	"unless":       assertScanToken(0, token.KeywordUnless, nil),
	"until":        assertScanToken(0, token.KeywordUntil, nil),
	"when":         assertScanToken(0, token.KeywordWhen, nil),
	"while":        assertScanToken(0, token.KeywordWhile, nil),
	"yield":        assertScanToken(0, token.KeywordYield, nil),
}

func TestScanner(t *testing.T) {
	debug.Enable = true
	for input, r := range rules {
		debug.Printf("input: %q", input)
		s := New([]byte(input))
		r(t, s)
	}
}
