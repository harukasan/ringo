package token

import "testing"

func TestKeywordToken(t *testing.T) {
	for i, l := range keywordLiterals {
		for j, literal := range l {
			if got, want := KeywordToken(literal), keywordTokens[i][j]; got != want {
				t.Errorf("KeywordToken(%v)=%v (want=%v)", literal, got, want)
			}
		}
	}
}
