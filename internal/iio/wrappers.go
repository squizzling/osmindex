package iio

import (
	"io"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/internal/iunsafe"
)

func FWriteSliceI64(f *os.File, data []int64) {
	b := iunsafe.Int64SliceAsByteSlice(data)
	if n, err := f.Write(b); err != nil {
		panic(err)
	} else if n != len(b) {
		panic("torn write")
	}
}

func FWriteSliceU64(f *os.File, data []uint64) {
	b := iunsafe.Uint64SliceAsByteSlice(data)
	if n, err := f.Write(b); err != nil {
		panic(err)
	} else if n != len(b) {
		panic("torn write")
	}
}

func FSeek(f *os.File, o int64, w int) int64 {
	if n, err := f.Seek(o, w); err != nil {
		panic(err)
	} else {
		return n
	}
}

func FClose(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func FRead(f *os.File, buf []byte) int {
	if n, err := f.Read(buf); err != nil && err != io.EOF {
		panic(err)
	} else {
		return n
	}
}

func MMapMapWrite(f *os.File) mmap.MMap {
	if m, err := mmap.Map(f, mmap.RDWR, 0); err != nil {
		panic(err)
	} else {
		return m
	}
}

func MMapMap(f *os.File) mmap.MMap {
	if m, err := mmap.Map(f, mmap.RDONLY, 0); err != nil {
		panic(err)
	} else {
		return m
	}
}

func MMapMapRegion(f *os.File, length, offset uint64) mmap.MMap {
	if m, err := mmap.MapRegion(f, int(length), mmap.RDONLY, 0, int64(offset)); err != nil {
		panic(err)
	} else {
		return m
	}
}

func MMapUnmap(m mmap.MMap) {
	if err := m.Unmap(); err != nil {
		panic(err)
	}
}

func FWrite(f *os.File, buf []byte) {
	remaining := len(buf)
	for remaining > 0 {
		n, err := f.Write(buf)
		if err != nil {
			panic(err)
		}
		if n == 0 {
			panic("zero write")
		}
		remaining -= n
		buf = buf[n:]
	}
}

func FCopy(fDst, fSrc *os.File, buf []byte) int {
	n := FRead(fSrc, buf)
	FWrite(fDst, buf[:n])
	return n
}

func IoCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		panic(err)
	}
}

func OsOpen(fn string) *os.File {
	if f, err := os.Open(fn); err != nil {
		panic(err)
	} else {
		return f
	}
}

func OsCreate(fn string) *os.File {
	if f, err := os.Create(fn); err != nil {
		panic(err)
	} else {
		return f
	}
}

func OsRemove(fn string) {
	if err := os.Remove(fn); err != nil {
		panic(err)
	}
}

func OsRemoveTry(fn string) {
	_ = os.Remove(fn)
}

func OsRename(from, to string) {
	if err := os.Rename(from, to); err != nil {
		panic(err)
	}
}
