package main

import (
	"fmt"
	"os"

	"github.com/squizzling/osmindex/internal/morton"
	"github.com/squizzling/osmindex/internal/nodeindex"
	"github.com/squizzling/osmindex/internal/system"
)

func main() {
	opts := parseArgs(os.Args[1:])

	teardown := system.Setup(&opts.SystemOpts)
	defer teardown()

	idx := nodeindex.NewNodeIndex(opts.InputFile, nil)

	var si *nodeindex.NodeSubIndex
	if opts.SubIndexIndex != nil {
		si = idx.SubIndex(*opts.SubIndexIndex, *opts.SubIndexCount)
	} else {
		si = idx.SubIndex(0, 1)
	}

	startElement := 0
	if opts.StartElement != nil {
		startElement = *opts.StartElement
	}

	endElement := si.ElementCount()
	if opts.EndElement != nil && *opts.EndElement < endElement {
		endElement = *opts.EndElement
	}

	for currentElement := startElement; currentElement < endElement; currentElement++ {
		element := si.Element(currentElement)
		if false {
			fmt.Printf("[%d] = %d+%d==%d @ %d\n",
				currentElement,
				element.FirstID(),
				element.Count(),
				element.LastID(),
				element.Offset(),
			)
		}
		if opts.SearchIndex || opts.PrintLocations {
			for idx := uint64(0); idx < element.Count(); idx++ {
				var extract, search uint64
				if opts.PrintLocations {

					extract = uint64(si.Location(element.Offset()+idx))
					if extract == 0 {
						fmt.Printf("Extract: skip\n")
					} else {
						lng, lat := morton.DecodeLocation(extract)
						fmt.Printf("Extract: %x (%f, %f) [%d+%d=%d]\n",
							extract,
							float64(lng)/10_000_000, float64(lat)/10_000_000,
							element.FirstID(),
							idx,
							element.FirstID()+idx,
						)
					}
				}
				if opts.SearchIndex {
					l, _, found := si.FindNode(element.FirstID()+idx, 0)
					search = uint64(l)
					if !found {
						fmt.Printf(" Search: skip\n")
					} else {
						lng, lat := morton.DecodeLocation(search)
						fmt.Printf(" Search: %x (%f, %f)\n",
							search,
							float64(lng)/10_000_000, float64(lat)/10_000_000,
						)
					}
				}

				if opts.PrintLocations && opts.SearchIndex {
					if extract != search {
						fmt.Printf("XX MISMATCH!\n")
					}
				}
			}
		}
	}
}
