package main

import (
	"fmt"
	"os"

	"github.com/squizzling/osmindex/internal/pb"
	"github.com/squizzling/osmindex/internal/system"
	"github.com/squizzling/osmindex/internal/wayindex"
)

func main() {
	opts := parseArgs(os.Args[1:])

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	ix := wayindex.NewWayIndex(opts.InputFile, nil)

	var si *wayindex.SubWayIndex
	if opts.SubIndexIndex != nil {
		si = ix.SubIndexWithScaledPass(*opts.SubIndexIndex, *opts.SubIndexCount)
	} else {
		si = ix.SubIndexWithScaledPass(0, 1)
	}

	for _, data := range ix.Ids {
		if opts.PrintLocations {
			printBytes := ix.Loc[data.Offset() << ix.AlignmentShift:]
			printLocs := pb.DecodeS64PackedDeltaZero(printBytes)
			fmt.Printf("%d %v\n", data.WayID(), printLocs)
		}
		if opts.SearchIndex {
			searchLocs, _, ok := si.FindLocations(data.WayID(), 0)
			if !ok {
				fmt.Printf("XX failed to find way %d\n", data.WayID())
			} else {
				fmt.Printf("%d %v\n", data.WayID(), searchLocs)
			}
		}
	}
}
