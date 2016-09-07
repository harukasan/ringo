package token

import "bytes"

//go:generate $GOPATH/bin/stringer -type=Token

// Token represents TODO
type Token int

// Token definitions:
const (
	None Token = iota
	Continue
	Illegal
	IDENT
	EOF
	NewLine // new line

	BinaryInteger
	DecimalInteger
	OctadecimalInteger
	HexadecimalInteger
	Float
	StringPart
	HeredocBegin
	HeredocPart

	// brackets:
	LParen   // (
	RParen   // )
	LBracket // [
	RBracket // ]
	LBrace   // {
	RBrace   // }

	// delimiters
	Colon2    // ::
	Comma     // ,
	Semicolon // ;
	Dot       // .
	Dot2      // ..
	Dot3      // ...
	Question  // ?
	Colon     // :
	Arrow     // =>

	// operators
	Not         // !
	NotEqual    // !=
	NotMatch    // !~
	AndOperator // &&
	OrOperator  // ||

	// operator methods
	Xor        // ^
	Amp        // &
	Or         // |
	Compare    // <=>
	Eq         // ==
	Eql        // ==
	Match      // =~
	Gt         // >
	GtEq       // >=
	Lt         // <
	LtEq       // <=
	LShift     // <<
	RShift     // >>
	Plus       // +
	Minus      // -
	Mul        // *
	Div        // /
	Mod        // %
	Pow        // **
	Invert     // ~
	UnaryPlus  // +@
	UnaryMinus // -@
	ElementSet // []
	ElementRef // []=

	// assign operator:
	Assign            // =
	AssignAndOperator // &&=
	AssignOrOperator  // ||=
	AssignXor         // ^=
	AssignAnd         // &=
	AssignOr          // |=
	AssignLShift      // <<=
	AssignRShift      // >>=
	AssignPlus        // +=
	AssignMinus       // -=
	AssignMul         // *=
	AssignDiv         // /=
	AssignMod         // %=
	AssignPow         // **=

	// keywords:
	KeywordLINE     // __LINE__
	KeywordENCODING // __ENCODING__
	KeywordFILE     // __FILE__
	KeywordBEGIN    // BEGIN
	KeywordEND      // END
	KeywordAlias    // alias
	KeywordAnd      // and
	KeywordBegin    // begin
	KeywordBreak    // break
	KeywordCase     // case
	KeywordClass    // class
	KeywordDef      // def
	KeywordDefined  // defined?
	KeywordDo       // do
	KeywordElse     // else
	KeywordElsif    // elsif
	KeywordEnd      // end
	KeywordEnsure   // ensure
	KeywordFor      // for
	KeywordFalse    // false
	KeywordIf       // if
	KeywordIn       // in
	KeywordModule   // module
	KeywordNext     // next
	KeywordNil      // nil
	KeywordNot      // not
	KeywordOr       // or
	KeywordRedo     // redo
	KeywordRescue   // rescue
	KeywordRetry    // retry
	KeywordReturn   // return
	KeywordSelf     // self
	KeywordSuper    // super
	KeywordThen     // then
	KeywordTrue     // true
	KeywordUndef    // undef
	KeywordUnless   // unless
	KeywordUntil    // until
	KeywordWhen     // when
	KeywordWhile    // while
	KeywordYield    // yield
)

var keywordLiterals = [127][][]byte{
	'_': [][]byte{
		[]byte("__LINE__"),
		[]byte("__ENCODING__"),
		[]byte("__FILE__"),
	},
	'B': [][]byte{
		[]byte("BEGIN"),
	},
	'E': [][]byte{
		[]byte("END"),
	},
	'a': [][]byte{
		[]byte("alias"),
		[]byte("and"),
	},
	'b': [][]byte{
		[]byte("begin"),
		[]byte("break"),
	},
	'c': [][]byte{
		[]byte("case"),
		[]byte("class"),
	},
	'd': [][]byte{
		[]byte("def"),
		[]byte("defined?"),
		[]byte("do"),
	},
	'e': [][]byte{
		[]byte("else"),
		[]byte("elsif"),
		[]byte("end"),
		[]byte("ensure"),
	},
	'f': [][]byte{
		[]byte("for"),
		[]byte("false"),
	},
	'i': [][]byte{
		[]byte("if"),
		[]byte("in"),
	},
	'm': [][]byte{
		[]byte("module"),
	},
	'n': [][]byte{
		[]byte("next"),
		[]byte("nil"),
		[]byte("not"),
	},
	'o': [][]byte{
		[]byte("or"),
	},
	'r': [][]byte{
		[]byte("redo"),
		[]byte("rescue"),
		[]byte("retry"),
		[]byte("return"),
	},
	's': [][]byte{
		[]byte("self"),
		[]byte("super"),
	},
	't': [][]byte{
		[]byte("then"),
		[]byte("true"),
	},
	'u': [][]byte{
		[]byte("undef"),
		[]byte("unless"),
		[]byte("until"),
	},
	'w': [][]byte{
		[]byte("when"),
		[]byte("while"),
	},
	'y': [][]byte{
		[]byte("yield"),
	},
}

var keywordTokens = [127][]Token{
	'_': []Token{
		KeywordLINE,     // __LINE__
		KeywordENCODING, // __ENCODING__
		KeywordFILE,     // __FILE__
	},
	'B': []Token{
		KeywordBEGIN, // BEGIN
	},
	'E': []Token{
		KeywordEND, // END
	},
	'a': []Token{
		KeywordAlias, // alias
		KeywordAnd,   // and
	},
	'b': []Token{
		KeywordBegin, // begin
		KeywordBreak, // break
	},
	'c': []Token{
		KeywordCase,  // case
		KeywordClass, // class
	},
	'd': []Token{
		KeywordDef,     // def
		KeywordDefined, // defined?
		KeywordDo,      // do
	},
	'e': []Token{
		KeywordElse,   // else
		KeywordElsif,  // elsif
		KeywordEnd,    // end
		KeywordEnsure, // ensure
	},
	'f': []Token{
		KeywordFor,   // for
		KeywordFalse, // false
	},
	'i': []Token{
		KeywordIf, // if
		KeywordIn, // in
	},
	'm': []Token{
		KeywordModule, // module
	},
	'n': []Token{
		KeywordNext, // next
		KeywordNil,  // nil
		KeywordNot,  // not
	},
	'o': []Token{
		KeywordOr, // or
	},
	'r': []Token{
		KeywordRedo,   // redo
		KeywordRescue, // rescue
		KeywordRetry,  // retry
		KeywordReturn, // return
	},
	's': []Token{
		KeywordSelf,  // self
		KeywordSuper, // super
	},
	't': []Token{
		KeywordThen, // then
		KeywordTrue, // true
	},
	'u': []Token{
		KeywordUndef,  // undef
		KeywordUnless, // unless
		KeywordUntil,  // until
	},
	'w': []Token{
		KeywordWhen,  // when
		KeywordWhile, // while
	},
	'y': []Token{
		KeywordYield, // yield
	},
}

// KeywordToken returns the token identifier that is mathced to given literal.
// If the literal is not matched to any keyword, it returns IDENT token.
func KeywordToken(literal []byte) Token {
	initial := literal[0]
	if list := keywordLiterals[initial]; list != nil {
		for i := 0; i < len(list); i++ {
			if bytes.Equal(list[i], literal) {
				return keywordTokens[initial][i]
			}
		}
	}
	return IDENT
}
