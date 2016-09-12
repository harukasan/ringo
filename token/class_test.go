package token

import "testing"

type class int

const (
	Letter class = 1 << iota
	WhiteSpace
	Uppercase
	Lowercase
	Decimal
	NonZeroDecimal
	Octadecimal
	Hexadecimal
	IdentStart
	Ident
)

var tests = map[byte]class{
	0x09: WhiteSpace,
	0x0b: WhiteSpace,
	0x0c: WhiteSpace,
	0x0d: WhiteSpace,
	' ':  WhiteSpace,
	'0':  Decimal | Octadecimal | Hexadecimal | Ident,
	'1':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'2':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'3':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'4':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'5':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'6':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'7':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident,
	'8':  Decimal | NonZeroDecimal | Hexadecimal | Ident,
	'9':  Decimal | NonZeroDecimal | Hexadecimal | Ident,
	'A':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'B':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'C':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'D':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'E':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'F':  Letter | Uppercase | Hexadecimal | IdentStart | Ident,
	'G':  Letter | Uppercase | IdentStart | Ident,
	'H':  Letter | Uppercase | IdentStart | Ident,
	'I':  Letter | Uppercase | IdentStart | Ident,
	'J':  Letter | Uppercase | IdentStart | Ident,
	'K':  Letter | Uppercase | IdentStart | Ident,
	'L':  Letter | Uppercase | IdentStart | Ident,
	'M':  Letter | Uppercase | IdentStart | Ident,
	'N':  Letter | Uppercase | IdentStart | Ident,
	'O':  Letter | Uppercase | IdentStart | Ident,
	'P':  Letter | Uppercase | IdentStart | Ident,
	'Q':  Letter | Uppercase | IdentStart | Ident,
	'R':  Letter | Uppercase | IdentStart | Ident,
	'S':  Letter | Uppercase | IdentStart | Ident,
	'T':  Letter | Uppercase | IdentStart | Ident,
	'U':  Letter | Uppercase | IdentStart | Ident,
	'V':  Letter | Uppercase | IdentStart | Ident,
	'W':  Letter | Uppercase | IdentStart | Ident,
	'X':  Letter | Uppercase | IdentStart | Ident,
	'Y':  Letter | Uppercase | IdentStart | Ident,
	'Z':  Letter | Uppercase | IdentStart | Ident,
	'_':  IdentStart | Ident,
	'a':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'b':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'c':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'd':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'e':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'f':  Letter | Lowercase | Hexadecimal | IdentStart | Ident,
	'g':  Letter | Lowercase | IdentStart | Ident,
	'h':  Letter | Lowercase | IdentStart | Ident,
	'i':  Letter | Lowercase | IdentStart | Ident,
	'j':  Letter | Lowercase | IdentStart | Ident,
	'k':  Letter | Lowercase | IdentStart | Ident,
	'l':  Letter | Lowercase | IdentStart | Ident,
	'm':  Letter | Lowercase | IdentStart | Ident,
	'n':  Letter | Lowercase | IdentStart | Ident,
	'o':  Letter | Lowercase | IdentStart | Ident,
	'p':  Letter | Lowercase | IdentStart | Ident,
	'q':  Letter | Lowercase | IdentStart | Ident,
	'r':  Letter | Lowercase | IdentStart | Ident,
	's':  Letter | Lowercase | IdentStart | Ident,
	't':  Letter | Lowercase | IdentStart | Ident,
	'u':  Letter | Lowercase | IdentStart | Ident,
	'v':  Letter | Lowercase | IdentStart | Ident,
	'w':  Letter | Lowercase | IdentStart | Ident,
	'x':  Letter | Lowercase | IdentStart | Ident,
	'y':  Letter | Lowercase | IdentStart | Ident,
	'z':  Letter | Lowercase | IdentStart | Ident,
}

var funcs = map[class]func(byte) bool{
	Letter:         IsLetter,
	WhiteSpace:     IsWhiteSpace,
	Uppercase:      IsUppercase,
	Lowercase:      IsLowercase,
	Decimal:        IsDecimal,
	NonZeroDecimal: IsNonZeroDecimal,
	Octadecimal:    IsOctadecimal,
	Hexadecimal:    IsHexadecimal,
	IdentStart:     IsIdentStart,
	Ident:          IsIdent,
}

func TestClass(t *testing.T) {
	for sbj, sbjC := range tests {
		for fC, f := range funcs {
			if want := sbjC&fC > 0; want != f(sbj) {
				t.Errorf("%v: '%d' must returns %v", fC, sbj, want)
			}
		}
	}
}
