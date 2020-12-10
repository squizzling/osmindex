package t

import (
	"testing"
)

func TestStringTableBuilderSorts1(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("def")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.Finalize()
	if stb.GetIndex("abc") != 1 {
		t.Errorf("abc should be first")
	}
	if stb.GetIndex("def") != 2 {
		t.Errorf("def should be second")
	}
}

func TestStringTableBuilderSorts2(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("def")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.AddString("abc")
	stb.Finalize()
	if stb.GetIndex("abc") != 1 {
		t.Errorf("abc should be first")
	}
	if stb.GetIndex("def") != 2 {
		t.Errorf("def should be second")
	}
}

func TestStringTableBuilderSorts3(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("def")
	stb.AddString("def")
	stb.AddString("def")
	stb.AddString("def")
	stb.AddString("def")
	stb.AddString("def")
	stb.AddString("abc")
	stb.Finalize()
	if stb.GetIndex("def") != 1 {
		t.Errorf("def should be first")
	}
	if stb.GetIndex("abc") != 2 {
		t.Errorf("abc should be second")
	}
}

func TestStringTableBuilderSortsEqual(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("abc")
	stb.AddString("def")
	stb.AddString("ghi")
	stb.Finalize()
	if stb.GetIndex("abc") != 1 {
		t.Errorf("abc should be first")
	}
	if stb.GetIndex("def") != 2 {
		t.Errorf("def should be second")
	}
	if stb.GetIndex("ghi") != 3 {
		t.Errorf("ghi should be third")
	}
}

func TestStringTableBuilderPanics(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("abc")
	stb.Finalize()
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("def should panic")
			}
		}()
		stb.GetIndex("def")
	}()
}

func TestStringTableBuilderToKeyVals(t *testing.T) {
	stb := NewStringTableBuilder()
	stb.AddString("c")
	stb.AddString("d")
	stb.AddString("a")
	stb.AddString("b")
	stb.AddString("2")
	stb.AddString("4")
	stb.AddString("1")
	stb.AddString("3")
	stb.Finalize()
	ks, vs := stb.ToKeyVals(map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"})
	// expected sort ordering is 1234abcd, as the keys are abcd they
	// should be 5, 6, 7, 8 and the values should be 1, 2, 3, 4.
	expectedKeys := []uint32{5, 6, 7, 8}
	expectedVals := []uint32{1, 2, 3, 4}
	for i := 0; i < 4; i++ {
		if expectedKeys[i] != ks[i] {
			t.Fatalf("keys is not as expected: %#v", ks)
		}
	}
	for i := 0; i < 4; i++ {
		if expectedVals[i] != vs[i] {
			t.Fatalf("vals is not as expected: %#v", vs)
		}
	}
}
