package t

import (
	"math"
	"sort"
	"strings"
)

type stringTableBuilder struct {
	stringFrequency map[string]uint32
	stringIndex     map[string]uint32
	stringList      []string
}

func NewStringTableBuilder() *stringTableBuilder {
	return &stringTableBuilder{
		stringFrequency: map[string]uint32{"": math.MaxUint32},
		stringList:      []string{""},
	}
}

func (stb *stringTableBuilder) AddMap(m map[string]string) {
	for k, v := range m {
		stb.AddString(k)
		stb.AddString(v)
	}
}

func (stb *stringTableBuilder) AddString(s string) {
	if _, ok := stb.stringFrequency[s]; !ok {
		stb.stringFrequency[s] = 1
		stb.stringList = append(stb.stringList, s)
	} else {
		stb.stringFrequency[s] = stb.stringFrequency[s] + 1
	}
}

func (stb *stringTableBuilder) Finalize() {
	sort.Slice(stb.stringList, func(i, j int) bool {
		freqI := stb.stringFrequency[stb.stringList[i]]
		freqJ := stb.stringFrequency[stb.stringList[j]]
		if freqI != freqJ {
			return freqI > freqJ
		}

		// Tie breaker to make sort deterministic
		return strings.Compare(stb.stringList[i], stb.stringList[j]) < 0
	})

	stb.stringIndex = stb.stringFrequency
	stb.stringFrequency = nil
	for idx, s := range stb.stringList {
		stb.stringIndex[s] = uint32(idx)
	}
}

func (stb *stringTableBuilder) GetIndex(s string) uint32 {
	id, ok := stb.stringIndex[s]
	if !ok {
		panic(s)
	}
	return id
}

func (stb *stringTableBuilder) ToStringTable() *StringTable {
	return &StringTable{
		S: stb.stringList,
	}
}

func (stb *stringTableBuilder) ToKeyVals(m map[string]string) ([]uint32, []uint32) {
	ks := make([]uint32, 0, len(m))
	vs := make([]uint32, 0, len(m))

	sortedKeys := make([]string, 0, len(m))
	for k, _ := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		ks = append(ks, stb.GetIndex(k))
		vs = append(vs, stb.GetIndex(m[k]))
	}
	return ks, vs

}
