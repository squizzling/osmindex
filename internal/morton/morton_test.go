package morton

import (
	"math/rand"
	"testing"
)

func TestMorton(t *testing.T) {
	r := rand.New(rand.NewSource(12345))

	for i := 0; i < 10000000; i++ {
		evenIn := int32(r.Uint32()) // get negatives
		oddIn := int32(r.Uint32())
		evenOut, oddOut := Decode(Encode(evenIn, oddIn))
		if evenIn != int32(evenOut) {
			t.Fatalf("i=%d, evenIn=%x, evenOut=%x", i, evenIn, evenOut)
		}
		if oddIn != int32(oddOut) {
			t.Fatalf("i=%d, oddIn=%x, oddOut=%x", i, oddIn, oddOut)
		}
	}
}
