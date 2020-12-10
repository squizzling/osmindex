package nodeindex

import (
	"fmt"
	"math"
	"os"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/internal/iio"
	"github.com/squizzling/osmindex/internal/iunsafe"
)

type NodeIndex struct {
	f   *os.File
	mem mmap.MMap

	elements  []int64
	locations []int64

	lookupTimingCount []uint64
	divFactor         float64

	HitHint    uint64
	HitDefault uint64

	MissDeferred uint64
	MissDefault  uint64
}

func NewNodeIndex(fn string, timing *float64) *NodeIndex {
	f := iio.OsOpen(fn)
	m := iio.MMapMap(f)
	d := iunsafe.ByteSliceAsInt64Slice(m)
	version := d[0]
	if version != Version {
		panic("wrong version")
	}
	startOfLocations := d[1]
	d = d[2:]
	indexData := d[:startOfLocations*2]
	locationData := d[startOfLocations*2:]
	idx := &NodeIndex{
		f:         f,
		mem:       m,
		elements:  indexData,
		locations: locationData,
	}

	if timing != nil {
		idx.lookupTimingCount = make([]uint64, 200)
		idx.divFactor = *timing
	}

	return idx
}

func (idx *NodeIndex) LocationCount() uint64 {
	return uint64(len(idx.locations))
}

// SubIndex will create a NodeSubIndex where the number of locations is
// distributed evenly among the passes, but the number of elements is not.
func (idx *NodeIndex) SubIndex(curPass uint64, countPass uint64) *NodeSubIndex {
	var subIds []IndexElement
	iunsafe.ByteSliceAsArbSlice(
		iunsafe.Int64SliceAsByteSlice(idx.elements),
		&subIds,
	)

	sliceSize := uint64(len(idx.locations)) / countPass
	startTarget := sliceSize * curPass
	endTarget := startTarget + sliceSize

	runningTotal := uint64(0)
	startIndex := uint64(0)
	endIndex := uint64(len(subIds))
	for i := uint64(0); i < uint64(len(subIds)); i++ {
		if runningTotal >= startTarget {
			startIndex = i
			break
		}
		runningTotal += subIds[i].Count()
	}

	for i := startIndex; i < uint64(len(subIds)); i++ {
		if runningTotal >= endTarget {
			endIndex = i
			break
		}
		runningTotal += subIds[i].Count()
	}

	return &NodeSubIndex{
		parent:    idx,
		element:   subIds[startIndex:endIndex],
		locations: idx.locations,
		Top:       subIds[endIndex-1].LastID(),
	}
}

const cpuFreq = 3_300_000_000

func (idx *NodeIndex) PrintStats() {
	fmt.Printf("hitHint:      %d\n", idx.HitHint)
	fmt.Printf("hitDefault:   %d\n", idx.HitDefault)
	fmt.Printf("missDefault:  %d\n", idx.MissDefault)
	fmt.Printf("missDeferred: %d\n", idx.MissDeferred)

	// Assumes a 3.3GHz system
	cyclesPerNS := float64(cpuFreq) / (1000 * 1000 * 1000)

	for i, count := range idx.lookupTimingCount {
		if count != 0 {
			fmt.Printf(
				"time: %d %.0f-%.0f = %d\n",
				i,
				(math.Pow(idx.divFactor, float64(i))) / cyclesPerNS,
				(math.Pow(idx.divFactor, float64(i+1))) / cyclesPerNS,
				count,
			)
		}
	}
}

func (idx *NodeIndex) Close() {
	idx.elements = nil
	idx.locations = nil
	iio.MMapUnmap(idx.mem)
	iio.FClose(idx.f)
}
