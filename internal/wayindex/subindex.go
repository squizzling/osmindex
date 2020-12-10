package wayindex

import (
	"sync/atomic"

	"github.com/dterei/gotsc"

	"github.com/squizzling/osmindex/internal/pb"
)

type SubWayIndex struct {
	parent *WayIndex
	Ids    []WayIndexData
	Loc    []byte
	Top    int64
}

func (si *SubWayIndex) findLocationsActual(id int64, hint int) ([]int64, int, bool) {
	if id > si.Top {
		// Due to multiple passes, early passes will be missing higher IDs.  We don't
		// need to check lower IDs, because they should already be translated.
		atomic.AddUint64(&si.parent.MissDeferred, 1)
		return nil, hint, false
	}

	lowerPoint, upperPoint := 0, len(si.Ids)

	if hint < upperPoint {
		hintId := si.Ids[hint].WayID()
		if id < hintId {
			upperPoint = hint
		} else if id == hintId {
			atomic.AddUint64(&si.parent.HitHint, 1)

			rawLocationBytes := si.Loc[si.Ids[hint].Offset()<<si.parent.AlignmentShift:]
			return pb.DecodeS64PackedDeltaZero(rawLocationBytes), hint, true
		} else if id > hintId {
			lowerPoint = hint
		}
	}

	for lowerPoint < upperPoint {
		midPoint := int(uint(lowerPoint+upperPoint) >> 1) // avoid overflow when computing midPoint
		wid := si.Ids[midPoint].WayID()
		if id < wid {
			upperPoint = midPoint
		} else if id == wid {
			lowerPoint = midPoint
			break
		} else if id > wid {
			lowerPoint = midPoint + 1
		}
	}

	if lowerPoint >= len(si.Ids) || si.Ids[lowerPoint].WayID() != id {
		atomic.AddUint64(&si.parent.MissDefault, 1)
		return nil, hint, false
	}

	atomic.AddUint64(&si.parent.HitDefault, 1)
	rawLocationBytes := si.Loc[si.Ids[lowerPoint].Offset()<<si.parent.AlignmentShift:]
	return pb.DecodeS64PackedDeltaZero(rawLocationBytes), lowerPoint, true
}

var overhead = gotsc.TSCOverhead()

func (si *SubWayIndex) FindLocations(id int64, hint int) ([]int64, int, bool) {
	if si.parent.lookupTimingCount == nil {
		return si.findLocationsActual(id, hint)
	}

	start := gotsc.BenchStart()
	l, h, ok := si.findLocationsActual(id, hint)
	end := gotsc.BenchEnd()
	delta := float64(int64((end - start) - overhead))
	if delta > 1 && delta < 2*cpuFreq {
		// A page fault should be serviced in under 2 seconds, so if the delta is higher then it's likely a core
		// jump with a different timestamp, rather than a slow page fault.
		i := 0
		for delta > 1 {
			i++
			delta /= si.parent.divFactor
		}
		atomic.AddUint64(&si.parent.lookupTimingCount[i-1], 1)
	}
	return l, h, ok
}
