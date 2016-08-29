package token

// IsLetter returns whether the character is a alphabet that matches to [a-zA-Z]
// in regular expression.
func IsLetter(c byte) bool {
	return IsUppercase(c) || IsLowercase(c)
}

// IsWhiteSpace returns whether the character is a white space.
func IsWhiteSpace(c byte) bool {
	return c == 0x20 || c == 0x09 || c == 0x0b || c == 0x0c || c == 0x0d
}

// IsUppercase returns whether the character is a uppercase letter.
func IsUppercase(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// IsLowercase returns whether the character is a lowercase letter.
func IsLowercase(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// IsDecimal returns wheter the character is a decimal number.
func IsDecimal(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsNonZeroDecimal returns wheter the character is a decimal number without
// zero.
func IsNonZeroDecimal(c byte) bool {
	return '1' <= c && c <= '9'
}

// IsOctadecimal returns wheter the character is a hexadecimal number.
func IsOctadecimal(c byte) bool {
	return '0' <= c && c <= '7'
}

// IsHexadecimal returns wheter the character is a hexadecimal number.
func IsHexadecimal(c byte) bool {
	return IsDecimal(c) || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}
