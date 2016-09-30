package scanner

import (
	"bytes"
	"testing"
)

func TestScanSingleQuotedString(t *testing.T) {
	input := []byte(`'a\\\\\\\'b\c'`)
	want := []byte(`a\\\'b\c`)

	// copy to working buffer
	work := make([]byte, len(input))
	copy(work, input)

	s := New(work)
	s.next() // skip first '
	_, got := scanSingleQuotedString(s, '\'', 1)

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
