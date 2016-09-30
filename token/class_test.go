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
	Alnum
)

var tests = map[byte]class{
	0x09: WhiteSpace,
	0x0b: WhiteSpace,
	0x0c: WhiteSpace,
	0x0d: WhiteSpace,
	' ':  WhiteSpace,
	'0':  Decimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'1':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'2':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'3':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'4':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'5':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'6':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'7':  Decimal | NonZeroDecimal | Octadecimal | Hexadecimal | Ident | Alnum,
	'8':  Decimal | NonZeroDecimal | Hexadecimal | Ident | Alnum,
	'9':  Decimal | NonZeroDecimal | Hexadecimal | Ident | Alnum,
	'A':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'B':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'C':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'D':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'E':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'F':  Letter | Uppercase | Hexadecimal | IdentStart | Ident | Alnum,
	'G':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'H':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'I':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'J':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'K':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'L':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'M':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'N':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'O':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'P':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'Q':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'R':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'S':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'T':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'U':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'V':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'W':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'X':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'Y':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'Z':  Letter | Uppercase | IdentStart | Ident | Alnum,
	'_':  IdentStart | Ident,
	'a':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'b':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'c':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'd':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'e':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'f':  Letter | Lowercase | Hexadecimal | IdentStart | Ident | Alnum,
	'g':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'h':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'i':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'j':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'k':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'l':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'm':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'n':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'o':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'p':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'q':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'r':  Letter | Lowercase | IdentStart | Ident | Alnum,
	's':  Letter | Lowercase | IdentStart | Ident | Alnum,
	't':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'u':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'v':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'w':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'x':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'y':  Letter | Lowercase | IdentStart | Ident | Alnum,
	'z':  Letter | Lowercase | IdentStart | Ident | Alnum,
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
	Alnum:          IsAlnum,
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
