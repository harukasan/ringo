package scanner

import (
	"bytes"
	"testing"

	"github.com/harukasan/ringo/debug"
	"github.com/harukasan/ringo/token"
)

var rules = map[string][]struct {
	pos     int
	token   token.Token
	literal []byte
}{

	// new line
	"\n":   {{0, token.NewLine, nil}},
	"\r\n": {{1, token.NewLine, nil}},
	";":    {{0, token.NewLine, nil}},

	// white spaces
	"":       {{0, token.EOF, nil}},
	" ":      {{1, token.EOF, nil}},
	" \t":    {{2, token.EOF, nil}},
	"\\\n":   {{2, token.EOF, nil}},
	"\\\r\n": {{3, token.EOF, nil}},

	// comment
	" #": {{2, token.EOF, nil}},
	"a#b": {
		{0, token.IdentLocalVar, []byte("a")},
		{3, token.EOF, nil},
	},
	"=begin\nTEXT\n=end\n":          {{17, token.EOF, nil}},
	"=begin open\nTEXT\n=end close": {{27, token.EOF, nil}},

	// __END__
	"__END__":     {{0, token.EOF, nil}},
	"__END__\n1":  {{0, token.EOF, nil}},
	"__END__\r":   {{0, token.IdentLocalVar, []byte("__END__")}},
	"__END__\r\n": {{0, token.EOF, nil}},
	"\n__END__\n": {
		{0, token.NewLine, nil},
		{1, token.EOF, nil},
	},
	" __END__": {{1, token.IdentLocalVar, []byte("__END__")}},

	// delimiters
	",":   {{0, token.Comma, nil}},
	".":   {{0, token.Dot, nil}},
	"..":  {{0, token.Dot2, nil}},
	"...": {{0, token.Dot3, nil}},
	"?":   {{0, token.Question, nil}},
	":":   {{0, token.Colon, nil}},
	"::":  {{0, token.Colon2, nil}},
	"=>":  {{0, token.Arrow, nil}},

	// operators
	"!":  {{0, token.Not, nil}},
	"!=": {{0, token.NotEqual, nil}},
	"!~": {{0, token.NotMatch, nil}},
	"&&": {{0, token.AndOperator, nil}},
	"||": {{0, token.OrOperator, nil}},

	// operator methods
	"^":   {{0, token.Xor, nil}},
	"&":   {{0, token.Amp, nil}},
	"|":   {{0, token.Or, nil}},
	"<=>": {{0, token.Compare, nil}},
	"==":  {{0, token.Eq, nil}},
	"===": {{0, token.Eql, nil}},
	"=~":  {{0, token.Match, nil}},
	">":   {{0, token.Gt, nil}},
	">=":  {{0, token.GtEq, nil}},
	"<":   {{0, token.Lt, nil}},
	"<=":  {{0, token.LtEq, nil}},
	"<<":  {{0, token.LShift, nil}},
	">>":  {{0, token.RShift, nil}},
	"+":   {{0, token.Plus, nil}},
	"-":   {{0, token.Minus, nil}},
	"*":   {{0, token.Mul, nil}},
	"/":   {{0, token.Div, nil}},
	"%":   {{0, token.Mod, nil}},
	"**":  {{0, token.Pow, nil}},
	"~":   {{0, token.Invert, nil}},
	"+@":  {{0, token.UnaryPlus, nil}},
	"-@":  {{0, token.UnaryMinus, nil}},
	"[]":  {{0, token.ElementRef, nil}},
	"[]=": {{0, token.ElementSet, nil}},
	"1 << 1": {
		{0, token.DecimalInteger, []byte("1")},
		{2, token.LShift, nil},
		{5, token.DecimalInteger, []byte("1")},
	},
	"1<<1": {
		{0, token.DecimalInteger, []byte("1")},
		{1, token.LShift, nil},
		{3, token.DecimalInteger, []byte("1")},
	},

	// assign operators
	"=":   {{0, token.Assign, nil}},
	"&&=": {{0, token.AssignAndOperator, nil}},
	"||=": {{0, token.AssignOrOperator, nil}},
	"^=":  {{0, token.AssignXor, nil}},
	"&=":  {{0, token.AssignAnd, nil}},
	"|=":  {{0, token.AssignOr, nil}},
	"<<=": {{0, token.AssignLShift, nil}},
	">>=": {{0, token.AssignRShift, nil}},
	"+=":  {{0, token.AssignPlus, nil}},
	"-=":  {{0, token.AssignMinus, nil}},
	"*=":  {{0, token.AssignMul, nil}},
	"/=":  {{0, token.AssignDiv, nil}},
	"%=":  {{0, token.AssignMod, nil}},
	"**=": {{0, token.AssignPow, nil}},

	// numeric literals
	"+1":                 {{0, token.DecimalInteger, []byte("+1")}},
	"-1":                 {{0, token.DecimalInteger, []byte("-1")}},
	"0":                  {{0, token.DecimalInteger, []byte("0")}},
	"123_456_789_0":      {{0, token.DecimalInteger, []byte("123_456_789_0")}},
	"0d1234567890":       {{0, token.DecimalInteger, []byte("0d1234567890")}},
	"0b0101":             {{0, token.BinaryInteger, []byte("0b0101")}},
	"0B0101":             {{0, token.BinaryInteger, []byte("0B0101")}},
	"0_12345670":         {{0, token.OctadecimalInteger, []byte("0_12345670")}},
	"0o12345670":         {{0, token.OctadecimalInteger, []byte("0o12345670")}},
	"0O12345670":         {{0, token.OctadecimalInteger, []byte("0O12345670")}},
	"0x1234567890abcdef": {{0, token.HexadecimalInteger, []byte("0x1234567890abcdef")}},
	"0X1234567890abcdef": {{0, token.HexadecimalInteger, []byte("0X1234567890abcdef")}},
	"+0x1":               {{0, token.HexadecimalInteger, []byte("+0x1")}},
	"-0X1":               {{0, token.HexadecimalInteger, []byte("-0X1")}},
	"0.1":                {{0, token.Float, []byte("0.1")}},
	"+0.1":               {{0, token.Float, []byte("+0.1")}},
	"-0.1":               {{0, token.Float, []byte("-0.1")}},
	"123.0456789e10":     {{0, token.Float, []byte("123.0456789e10")}},
	"123.0456789E10":     {{0, token.Float, []byte("123.0456789E10")}},
	"+1\n-1": {
		{0, token.DecimalInteger, []byte("+1")},
		{2, token.NewLine, nil},
		{3, token.DecimalInteger, []byte("-1")},
	},
	"x+1": {
		{0, token.IdentLocalVar, []byte("x")},
		{1, token.Plus, nil},
		{2, token.DecimalInteger, []byte("1")},
	},
	"x +1": {
		{0, token.IdentLocalVar, []byte("x")},
		{2, token.DecimalInteger, []byte("+1")},
	},
	"x-1": {
		{0, token.IdentLocalVar, []byte("x")},
		{1, token.Minus, nil},
		{2, token.DecimalInteger, []byte("1")},
	},
	"x -1": {
		{0, token.IdentLocalVar, []byte("x")},
		{2, token.DecimalInteger, []byte("-1")},
	},

	// string literals:
	`"a"`: {{0, token.String, []byte(`a`)}},
	`"\"`: {{0, token.String, []byte(``)}},
	`"#{a}"`: {
		{0, token.StringPart, []byte("")},
		{1, token.InsertBegin, nil}, // points to '#'
		{3, token.IdentLocalVar, []byte("a")},
		{4, token.InsertEnd, nil},
	},
	`"#@@a"`: {
		{0, token.StringPart, []byte("")},
		{1, token.IdentClassVar, []byte("@@a")}, // points to '#'
	},
	`"#@a"`: {
		{0, token.StringPart, []byte("")},
		{1, token.IdentInstanceVar, []byte("@a")}, // points to '#'
	},
	`"#$a"`: {
		{0, token.StringPart, []byte("")},
		{1, token.IdentGlobalVar, []byte("$a")}, // points to '#'
	},
	`"#{"a"}"`: {
		{0, token.StringPart, []byte("")},
		{1, token.InsertBegin, nil}, // points to '#'
		{3, token.String, []byte("a")},
		{6, token.InsertEnd, nil},
		{7, token.String, []byte("")},
	},
	`""""`: {
		{0, token.String, []byte("")},
		{2, token.String, []byte("")},
	},
	`"#{}#{}"`: {
		{0, token.StringPart, []byte("")},
		{1, token.InsertBegin, []byte("")},
		{3, token.InsertEnd, []byte("")},
		{4, token.InsertBegin, []byte("")},
		{6, token.InsertEnd, []byte("")},
		{7, token.String, []byte("")},
	},
	`"".a`: {
		{0, token.String, []byte("")},
		{2, token.Dot, nil},
		{3, token.IdentLocalVar, []byte("a")},
	},
	`"#{}".a`: {
		{0, token.StringPart, []byte("")},
		{1, token.InsertBegin, []byte("")},
		{3, token.InsertEnd, []byte("")},
		{4, token.String, []byte("")},
		{5, token.Dot, nil},
		{6, token.IdentLocalVar, []byte("a")},
	},
	`"\n"`:       {{0, token.String, []byte{0x0a}}},
	`"\t"`:       {{0, token.String, []byte{0x09}}},
	`"\r"`:       {{0, token.String, []byte{0x0d}}},
	`"\f"`:       {{0, token.String, []byte{0x0c}}},
	`"\v"`:       {{0, token.String, []byte{0x0b}}},
	`"\a"`:       {{0, token.String, []byte{0x07}}},
	`"\e"`:       {{0, token.String, []byte{0x1b}}},
	`"\b"`:       {{0, token.String, []byte{0x08}}},
	`"\s"`:       {{0, token.String, []byte{0x20}}},
	`"\xf0\xFF"`: {{0, token.String, []byte{0xf0, 0xff}}},
	`"\377\377"`: {{0, token.String, []byte{0xff, 0xff}}},
	`"\ca"`:      {{0, token.String, []byte{'\a'}}},
	`"\C-a"`:     {{0, token.String, []byte{'\a'}}},
	"\"\\\n\"":   {{0, token.String, []byte(``)}},
	`%{\n}`:      {{0, token.String, []byte{0x0a}}},
	`%(\n)`:      {{0, token.String, []byte{0x0a}}},
	`%[\n]`:      {{0, token.String, []byte{0x0a}}},
	`%<\n>`:      {{0, token.String, []byte{0x0a}}},
	`%Q{\n}`:     {{0, token.String, []byte{0x0a}}},

	`'a'`:      {{0, token.String, []byte(`a`)}},
	`'\''`:     {{0, token.String, []byte(`'`)}},
	`'\\'`:     {{0, token.String, []byte(`\`)}},
	`'\a\\\''`: {{0, token.String, []byte(`\a\'`)}},
	`%q{a}`:    {{0, token.String, []byte(`a`)}},
	`%q(a)`:    {{0, token.String, []byte(`a`)}},
	`%q[a]`:    {{0, token.String, []byte(`a`)}},
	`%q<a>`:    {{0, token.String, []byte(`a`)}},
	`%q{\\}`:   {{0, token.String, []byte(`\`)}},
	`%q{\}}`:   {{0, token.String, []byte(`}`)}},
	`%q<\}>`:   {{0, token.String, []byte(`\}`)}},

	// heredoc
	"a <<TEXT, x\nabc\n\nTEXT\n": {
		{0, token.IdentLocalVar, []byte("a")},
		{2, token.HeredocBegin, []byte("<<TEXT")},
		{8, token.Comma, nil},
		{10, token.IdentLocalVar, []byte("x")},
		{11, token.NewLine, nil},
		{12, token.HeredocEnd, []byte("abc\n\n")},
		{22, token.EOF, nil},
	},
	"<<-TEXT\n  TEXT\n": {
		{0, token.HeredocBegin, []byte("<<-TEXT")},
		{7, token.NewLine, nil},
		{8, token.HeredocEnd, nil},
		{15, token.EOF, nil},
	},
	"<<-'TEXT'\n  TEXT\n": {
		{0, token.HeredocBegin, []byte("<<-'TEXT'")},
		{9, token.NewLine, nil},
		{10, token.HeredocEnd, nil},
		{17, token.EOF, nil},
	},
	"1 <<1\n1\n": {
		{0, token.DecimalInteger, []byte("1")},
		{2, token.HeredocBegin, []byte("<<1")},
		{5, token.NewLine, nil},
		{6, token.HeredocEnd, nil},
		{8, token.EOF, nil},
	},
	"<<TEXT\n#@a #$b\nTEXT\n": {
		{0, token.HeredocBegin, []byte("<<TEXT")},
		{6, token.NewLine, nil},
		{7, token.IdentInstanceVar, []byte("@a")},
		{10, token.HeredocPart, []byte(" ")},
		{11, token.IdentGlobalVar, []byte("$b")},
		{14, token.HeredocEnd, []byte("\n")},
		{20, token.EOF, nil},
	},
	"<<TEXT\n#{<<TEXT\nTEXT\n}\nTEXT\n": {
		{0, token.HeredocBegin, []byte("<<TEXT")},
		{6, token.NewLine, nil},
		{7, token.InsertBegin, nil},
		{9, token.HeredocBegin, []byte("<<TEXT")},
		{15, token.NewLine, nil},
		{16, token.HeredocEnd, nil},
		{21, token.InsertEnd, nil},
		{22, token.HeredocEnd, []byte("\n")},
		{28, token.EOF, nil},
	},

	// ident
	"v":        {{0, token.IdentLocalVar, []byte("v")}},
	"_":        {{0, token.IdentLocalVar, []byte("_")}},
	"v?":       {{0, token.IdentLocalMethod, []byte("v?")}},
	"v!":       {{0, token.IdentLocalMethod, []byte("v!")}},
	"v=":       {{0, token.IdentLocalMethod, []byte("v=")}},
	"$v":       {{0, token.IdentGlobalVar, []byte("$v")}},
	"$v1":      {{0, token.IdentGlobalVar, []byte("$v1")}},
	"@var1":    {{0, token.IdentInstanceVar, []byte("@var1")}},
	"@@var1":   {{0, token.IdentClassVar, []byte("@@var1")}},
	"Constant": {{0, token.IdentConst, []byte("Constant")}},
	"a  b": {
		{0, token.IdentLocalVar, []byte("a")},
		{3, token.IdentLocalVar, []byte("b")},
	},

	// keywords
	"__LINE__":     {{0, token.KeywordLINE, nil}},
	"__ENCODING__": {{0, token.KeywordENCODING, nil}},
	"__FILE__":     {{0, token.KeywordFILE, nil}},
	"BEGIN":        {{0, token.KeywordBEGIN, nil}},
	"END":          {{0, token.KeywordEND, nil}},
	"alias":        {{0, token.KeywordAlias, nil}},
	"and":          {{0, token.KeywordAnd, nil}},
	"begin":        {{0, token.KeywordBegin, nil}},
	"break":        {{0, token.KeywordBreak, nil}},
	"case":         {{0, token.KeywordCase, nil}},
	"class":        {{0, token.KeywordClass, nil}},
	"def":          {{0, token.KeywordDef, nil}},
	"defined?":     {{0, token.KeywordDefined, nil}},
	"do":           {{0, token.KeywordDo, nil}},
	"else":         {{0, token.KeywordElse, nil}},
	"elsif":        {{0, token.KeywordElsif, nil}},
	"end":          {{0, token.KeywordEnd, nil}},
	"ensure":       {{0, token.KeywordEnsure, nil}},
	"for":          {{0, token.KeywordFor, nil}},
	"false":        {{0, token.KeywordFalse, nil}},
	"if":           {{0, token.KeywordIf, nil}},
	"in":           {{0, token.KeywordIn, nil}},
	"module":       {{0, token.KeywordModule, nil}},
	"next":         {{0, token.KeywordNext, nil}},
	"nil":          {{0, token.KeywordNil, nil}},
	"not":          {{0, token.KeywordNot, nil}},
	"or":           {{0, token.KeywordOr, nil}},
	"redo":         {{0, token.KeywordRedo, nil}},
	"rescue":       {{0, token.KeywordRescue, nil}},
	"retry":        {{0, token.KeywordRetry, nil}},
	"return":       {{0, token.KeywordReturn, nil}},
	"self":         {{0, token.KeywordSelf, nil}},
	"super":        {{0, token.KeywordSuper, nil}},
	"then":         {{0, token.KeywordThen, nil}},
	"true":         {{0, token.KeywordTrue, nil}},
	"undef":        {{0, token.KeywordUndef, nil}},
	"unless":       {{0, token.KeywordUnless, nil}},
	"until":        {{0, token.KeywordUntil, nil}},
	"when":         {{0, token.KeywordWhen, nil}},
	"while":        {{0, token.KeywordWhile, nil}},
	"yield":        {{0, token.KeywordYield, nil}},
}

func TestScanner(t *testing.T) {
	for input, wants := range rules {
		debug.Printf("input: %q", input)
		s := New([]byte(input))

		for _, want := range wants {
			p, tk, l := s.Scan()
			if p != want.pos || tk != want.token || !bytes.Equal(l, want.literal) {
				format := "scan(src=%v): pos=%v (want=%v), token=%v (want=%v), literal=%v (want=%v)"
				debug.Printf(format, input, p, want.pos, tk, want.token, l, want.literal)
				t.Errorf(format, input, p, want.pos, tk, want.token, l, want.literal)
			}
		}
	}
}

// Note: When scanning a single quoted string, the source array of bytes will
// broken to unescape characters.
func TestScanSingleQuotedString(t *testing.T) {
	input := []byte(`'a\\\\\\\'b\c'`)
	want := []byte(`a\\\'b\c`)

	// copy to working buffer
	work := make([]byte, len(input))
	copy(work, input)

	s := New(work)
	s.next() // skip first '
	_, got := scanSingleQuote(s)

	if !bytes.Equal(got, want) {
		t.Fatalf("\ninput =%#v\nwant  =%#v\ngot   =%#v", string(input), string(want), string(got))
	}
}

func TestDecodeOctalEsc(t *testing.T) {
	rules := map[string]struct {
		n   int
		v   byte
		err error
	}{
		"7":   {n: 1, v: 07},
		"77":  {n: 2, v: 077},
		"377": {n: 3, v: 0377},
	}
	for input, want := range rules {
		n, v := decodeOctalEsc(NewString(input))
		if n != want.n || v != want.v {
			t.Errorf("decodeOctalEsc: n=%v (want=%v), v=%v (want=%v)", n, want.n, v, want.v)
		}
	}
}

func TestDecodeHexEsc(t *testing.T) {
	rules := map[string]struct {
		n int
		v byte
	}{
		"f":  {n: 1, v: 0xf},
		"ff": {n: 2, v: 0xff},
	}
	for input, want := range rules {
		n, v := decodeHexEsc(NewString(input))
		if n != want.n || v != want.v {
			t.Errorf("decodeHexEsc: n=%v (want=%v), v=%v (want=%v)", n, want.n, v, want.v)
		}
	}
}
