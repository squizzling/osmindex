package t

type RelId int64

type HLRelation struct {
	Id         RelId
	Tags       map[string]string
	Roles      []string
	Members    []int64
	MemberType []int
}

type Relations []*HLRelation

func (pblk *PrimitiveBlock) ToRelations() Relations {
	var rs Relations

	for _, pg := range pblk.PrimitiveGroup {
		for _, r := range pg.Relations {
			tags := make(map[string]string)
			for idx := 0; idx < len(r.Keys); idx++ {
				tags[pblk.GetString(r.Keys[idx])] = pblk.GetString(r.Vals[idx])
			}
			roles := make([]string, 0, len(r.RolesSID))
			for _, role := range r.RolesSID {
				roles = append(roles, pblk.GetString(uint32(role)))
			}
			ms := make([]int64, len(r.MemIDs))
			copy(ms, r.MemIDs)
			ts := make([]int, 0, len(r.MemberType))
			for _, t := range r.MemberType {
				ts = append(ts, int(t))
			}
			rs = append(rs, &HLRelation{
				Id:         r.ID,
				Tags:       tags,
				Roles:      roles,
				Members:    ms,
				MemberType: ts,
			})
		}

	}
	return rs
}

func (rs Relations) AddTagsToStringTable(stb *stringTableBuilder) {
	for _, r := range rs {
		stb.AddMap(r.Tags)
		for _, role := range r.Roles {
			stb.AddString(role)
		}
	}
}

func (rs Relations) ToPrimitiveBlock(stb *stringTableBuilder) *PrimitiveGroup {
	nrs := make([]*Relation, 0, len(rs))
	for _, r := range rs {
		ks, vs := stb.ToKeyVals(r.Tags)
		roles := make([]int32, 0, len(r.Roles))
		for _, role := range r.Roles {
			roles = append(roles, int32(stb.GetIndex(role)))
		}
		ms := make([]int64, len(r.Members))
		copy(ms, r.Members)
		ts := make([]uint64, 0, len(r.MemberType))
		for _, t := range r.MemberType {
			ts = append(ts, uint64(t))
		}
		nrs = append(nrs, &Relation{
			ID:         r.Id,
			Keys:       ks,
			Vals:       vs,
			RolesSID:   roles,
			MemIDs:     ms,
			MemberType: ts,
		})
	}
	return &PrimitiveGroup{
		Relations: nrs,
	}
}
