package nodeindex

type IndexElement struct {
	baseId      uint64
	countOffset CountOffset
}

func (ie IndexElement) FirstID() uint64 {
	return ie.baseId
}

func (ie IndexElement) LastID() uint64 {
	return ie.FirstID() + ie.countOffset.Count() - 1
}

func (ie IndexElement) Offset() uint64 {
	return ie.countOffset.Offset()
}

func (ie IndexElement) Count() uint64 {
	return ie.countOffset.Count()
}

func (ie IndexElement) OffsetForID(id uint64) uint64 {
	if !ie.contains(id) {
		panic("does not contain")
	}
	off := ie.Offset()
	start := ie.FirstID()
	delta := id - start
	return off + delta
}

func (ie IndexElement) compareId(id uint64) int {
	start := ie.FirstID()
	end := ie.LastID()

	if id < start {
		return -1
	} else if id > end {
		return 1
	} else {
		return 0
	}
}
func (ie IndexElement) contains(id uint64) bool {
	return id >= ie.FirstID() && id <= ie.LastID()
}
