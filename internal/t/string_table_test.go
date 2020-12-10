package t

import (
	"bytes"
	"testing"
)

func TestStringTableRead(t *testing.T) {
	// Deliberately not testing anything that would require a multi-byte varint encoding
	buf := []byte{
		1<<3 | 2, // id 1, type 2
		2,        // length
		'a', 'b',
		1<<3 | 2, // id 1, type 2
		3,        // length
		'a', 'b', 'c',
	}

	var st StringTable
	ReadStringTable(buf, &st)

	if len(st.S) != 2 {
		t.Fatal("expected length 2")
	}

	if st.S[0] != "ab" {
		t.Fatal("expected item 0 to be 'ab'")
	}
	if st.S[1] != "abc" {
		t.Fatal("expected item 1 to be 'abc'")
	}
}

func TestStringTableWrite(t *testing.T) {
	var st StringTable

	st.S = append(st.S, "ab")
	st.S = append(st.S, "abc")

	buf := []byte{0xff} // marker to ensure we're actually appending
	buf = st.Write(buf)

	expected := []byte{
		0xff,     // marker
		1<<3 | 2, // id 1, type 2
		2,        // length
		'a', 'b',
		1<<3 | 2, // id 1, type 2
		3,        // length
		'a', 'b', 'c',
	}

	if bytes.Compare(buf, expected) != 0 {
		t.Fatal("encode data mismatch")
	}
}
