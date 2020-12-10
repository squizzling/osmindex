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