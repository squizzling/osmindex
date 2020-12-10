package iunsafe

import (
	"testing"
)

type structTest struct {
	A int64
	B int64
}

func TestByteSliceAsArbSlice(t *testing.T) {
	b := Int64SliceAsByteSlice([]int64{0, 1, 2, 3, 4, 5, 6})
	var out []structTest
	ByteSliceAsArbSlice(b, &out)
	if len(out) != 3 { // ensure the 6 is ignored, because it's only a partial value
		t.Fatalf("length should be 3, is %d", len(out))
	}
	if cap(out) != len(out) {
		t.Fatalf("length shold match capacity, len is %d, cap is %d", len(out), cap(out))
	}
	if out[0].A != 0 {
		t.Fatalf("out[0].A != 0, is %x\n", out[0].A)
	}
	if out[0].B != 1 {
		t.Fatalf("out[0].B != 1, is %x\n", out[0].B)
	}
	if out[1].A != 2 {
		t.Fatalf("out[1].A != 2, is %x\n", out[1].A)
	}
	if out[1].B != 3 {
		t.Fatalf("out[1].B != 3, is %x\n", out[1].B)
	}
	if out[2].A != 4 {
		t.Fatalf("out[2].A != 4, is %x\n", out[2].A)
	}
	if out[2].B != 5 {
		t.Fatalf("out[2].B != 5, is %x\n", out[2].B)
	}
}
