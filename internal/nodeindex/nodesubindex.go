package nodeindex

import (
	"sync/atomic"

	"github.com/dterei/gotsc"
)

type NodeSubIndex struct {
	parent    *NodeIndex
	element   []IndexElement
	locations []int64
	Top       uint64
}

func (si *NodeSubIndex) Element(idx int) IndexElement {
	return si.element[idx]
}

func (si *NodeSubIndex) ElementCount() int {
	return len(si.element)
}

func (si *NodeSubIndex) Location(idx uint64) int64 {
	return si.locations[idx]
}

func (si *NodeSubIndex) findNodeActual(id uint64, hint int) (int64, int, bool) {
	if id > si.Top {
		// Due to multiple passes, early passes will be missing higher IDs.  We don't
		// need to check lower IDs, because they should already be translated.
		atomic.AddUint64(&si.parent.MissDeferred, 1)
		return 0, hint, false
	}

	lowerPoint, upperPoint := 0, len(si.element)

	// It can be disadvantageous to use the hint as an upper/lower, as "power 2 relative to the entire
	// range" pages are likely to be in the working set, but "power 2 relative to the hint" pages are not.
	if hint < upperPoint {
		switch si.element[hint].compareId(id) {
		case -1:
			upperPoint = hint
		case 0:
			atomic.AddUint64(&si.parent.HitHint, 1)
			loc := si.locations[si.element[hint].OffsetForID(id)]
			return loc, hint, loc != 0
		case 1:
			lowerPoint = hint
		}
	}

	for lowerPoint < upperPoint {
		midPoint := int(uint(lowerPoint+upperPoint) >> 1) // avoid overflow when computing midPoint
		cmp := si.element[midPoint].compareId(id)
		if cmp < 0 {
			upperPoint = midPoint
		} else if cmp == 0 {
			lowerPoint = midPoint
			break
		} else {
			lowerPoint = midPoint + 1
		}
	}

	if lowerPoint >= len(si.element) || !si.element[lowerPoint].contains(id) {
		atomic.AddUint64(&si.parent.MissDefault, 1)
		return 0, hint, false
	}

	atomic.AddUint64(&si.parent.HitDefault, 1)
	loc := si.locations[si.element[lowerPoint].OffsetForID(id)]
	return loc, lowerPoint, loc != 0
}

var overhead = gotsc.TSCOverhead()

func (si *NodeSubIndex) FindNode(id uint64, hint int) (int64, int, bool) {
	if si.parent.lookupTimingCount == nil {
		// gotsc overhead is a non-trivial fraction of the work on fast path, so it's worth skipping if it's not required
		return si.findNodeActual(id, hint)
	}

	start := gotsc.BenchStart()
	l, h, ok := si.findNodeActual(id, hint)
	end := gotsc.BenchEnd()
	delta := float64(int64((end - start) - overhead))
	if delta > 1 && delta < 2 * cpuFreq {
		// A page fault should be serviced in under 2 seconds, so if the delta is higher then it's likely a core
		// jump with a different timestamp, rather than a slow page fault.
		i := 0
		for delta > 1 {
			// TODO: something something log?
			i++
			delta /= si.parent.divFactor
		}
		atomic.AddUint64(&si.parent.lookupTimingCount[i-1], 1)
	}
	return l, h, ok
}
