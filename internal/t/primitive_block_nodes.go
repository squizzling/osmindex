package t

type NodeId uint64

type HLNode struct {
	Id   NodeId
	Tags map[string]string
	Lat  int64
	Lon  int64
}

type Nodes []*HLNode

func (pblk *PrimitiveBlock) ToNodes() Nodes {
	var ns Nodes

	for _, pg := range pblk.PrimitiveGroup {
		if pg.Dense != nil {
			for i := 0; i < len(pg.Dense.Id); i++ {
				tags := make(map[string]string)
				ks, vs := pg.Dense.GetKeyVals(i)
				for idx, kId := range ks {
					tags[pblk.GetString(kId)] = pblk.GetString(vs[idx])
				}
				ns = append(ns, &HLNode{
					Id:   NodeId(pg.Dense.Id[i]),
					Tags: tags,
					Lat:  pblk.LatOffset + (int64(pblk.Granularity) * pg.Dense.Lat[i]),
					Lon:  pblk.LonOffset + (int64(pblk.Granularity) * pg.Dense.Lon[i]),
				})
			}
		}

	}
	return ns
}

func (ns Nodes) AddTagsToStringTable(stb *stringTableBuilder) {
	for _, n := range ns {
		stb.AddMap(n.Tags)
	}
}

func (ns Nodes) ToPrimitiveGroup(stb *stringTableBuilder) *PrimitiveGroup {
	var dn DenseNodes

	dn.Id = make([]int64, 0, len(ns))
	dn.Lat = make([]int64, 0, len(ns))
	dn.Lon = make([]int64, 0, len(ns))

	for _, n := range ns {
		dn.Id = append(dn.Id, int64(n.Id))
		dn.Lat = append(dn.Lat, n.Lat/100)
		dn.Lon = append(dn.Lon, n.Lon/100)
		dn.KeyVals = append(dn.KeyVals, stb.ToKeyValLinear(n.Tags)...)
		dn.KeyVals = append(dn.KeyVals, 0)
	}

	if len(dn.KeyVals) == len(dn.Id) {
		// Every id must have a 0 terminating it. If the number of keyvals is exactly the number
		// of ids, then there was no actual keyvals in the entire set, so we can drop it entirely.
		dn.KeyVals = nil
	}

	return &PrimitiveGroup{
		Dense: &dn,
	}
}
