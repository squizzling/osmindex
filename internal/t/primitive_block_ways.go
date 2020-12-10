package t

type HLWay struct {
	Id    WayId
	Tags  map[string]string
	Nodes []NodeId
}

type Ways []*HLWay

func (pblk *PrimitiveBlock) ToWays() Ways {
	var ws Ways

	for _, pg := range pblk.PrimitiveGroup {
		for _, w := range pg.Ways {
			tags := make(map[string]string)
			for idx := 0; idx < len(w.Keys); idx++ {
				tags[pblk.GetString(w.Keys[idx])] = pblk.GetString(w.Vals[idx])
			}
			ns := make([]NodeId, 0, len(w.Refs))
			for _, id := range w.Refs {
				ns = append(ns, NodeId(id))
			}
			ws = append(ws, &HLWay{
				Id:    w.ID,
				Tags:  tags,
				Nodes: ns,
			})
		}

	}
	return ws
}

func (ws Ways) AddTagsToStringTable(stb *stringTableBuilder) {
	for _, w := range ws {
		stb.AddMap(w.Tags)
	}
}

func (ws Ways) ToPrimitiveGroup(stb *stringTableBuilder) *PrimitiveGroup {
	newWays := make([]*Way, 0, len(ws))
	for _, w := range ws {
		ks, vs := stb.ToKeyVals(w.Tags)
		refs := make([]int64, 0, len(w.Nodes))
		for _, id := range w.Nodes {
			refs = append(refs, int64(id))
		}
		newWays = append(newWays, &Way{
			ID:   w.Id,
			Keys: ks,
			Vals: vs,
			Refs: refs,
		})
	}
	return &PrimitiveGroup{
		Ways: newWays,
	}
}
