package wayindex

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/iunsafe"
)

type WayIndex struct {
	f   *os.File
	mem mmap.MMap

	Ids            []WayIndexData
	Loc            []byte
	AlignmentShift int

	lookupTimingCount []uint64
	divFactor         float64

	HitHint    uint64
	HitDefault uint64

	MissDeferred uint64
	MissDefault  uint64
}

const cpuFreq = 3_300_000_000

func NewWayIndex(fn string, timing *float64) *WayIndex {
	f := iio.OsOpen(fn)
	m := iio.MMapMap(f)
	startOfLocationsInBytes := binary.LittleEndian.Uint64(m[0:8])
	startOfLocationsInElements := startOfLocationsInBytes * 8
	alignmentShift := int(binary.LittleEndian.Uint64(m[8:16]))

	var indexData []WayIndexData
	iunsafe.ByteSliceAsArbSlice(m[16:16+startOfLocationsInElements], &indexData)
	locationData := m[16+startOfLocationsInElements:]
	idx := &WayIndex{
		f:              f,
		mem:            m,
		Ids:            indexData,
		Loc:            locationData,
		AlignmentShift: alignmentShift,
	}

	if timing != nil {
		idx.lookupTimingCount = make([]uint64, 200)
		idx.divFactor = *timing
	}

	return idx
}

// SubIndexWithScaledPass will create a SubWayIndex where the number of locations
// is distributed evenly among the passes, but the number of elements is not.
// This is much more useful for managing RSS, as the location data makes up the
// bulk of the index.
func (idx *WayIndex) SubIndexWithScaledPass(curPass int, countPass int) *SubWayIndex {
	subIds := idx.Ids

	// These targets are likely to be mid way, but that's fine
	startTarget := int64((len(idx.Loc) / countPass) * curPass)
	endTarget := int64((len(idx.Loc) / countPass) * (curPass + 1))

	startIndex := 0

	endIndex := len(subIds)
	for i := 0; i < len(subIds); i++ {
		if subIds[i].Offset() << idx.AlignmentShift >= startTarget {
			startIndex = i
			break
		}
	}

	for i := startIndex; i < len(subIds); i++ {
		if subIds[i].Offset() << idx.AlignmentShift  >= endTarget {
			endIndex = i
			break
		}
	}

	return &SubWayIndex{
		parent: idx,
		Ids:    subIds[startIndex:endIndex],
		Loc:    idx.Loc,
		Top:    int64(subIds[endIndex-1].WayID()),
	}
}

func (idx *WayIndex) PrintStats() {
	fmt.Printf("hitHint:      %d\n", idx.HitHint)
	fmt.Printf("hitDefault:   %d\n", idx.HitDefault)
	fmt.Printf("missDefault:  %d\n", idx.MissDefault)
	fmt.Printf("missDeferred: %d\n", idx.MissDeferred)

	// Assumes a 3.3GHz system
	cyclesPerNS := float64(cpuFreq) / (1000.0 * 1000.0 * 1000.0)

	for i, count := range idx.lookupTimingCount {
		if count != 0 {
			fmt.Printf(
				"time: %d %.0f-%.0f = %d\n",
				i,
				(math.Pow(idx.divFactor, float64(i)))/cyclesPerNS,
				(math.Pow(idx.divFactor, float64(i+1)))/cyclesPerNS,
				count,
			)
		}
	}
}

func (idx *WayIndex) Close() {
	idx.Ids = nil
	idx.Loc = nil
	iio.MMapUnmap(idx.mem)
	iio.FClose(idx.f)
}
