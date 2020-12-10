package tagfilterer

import (
	"context"
	"strings"

	"github.com/edsrzf/mmap-go"

	"github.com/squizzling/osmindex/cmd/tagfilter/tagfilterer/state"
	"github.com/squizzling/osmindex/internal/pbf"
	"github.com/squizzling/osmindex/internal/t"
)

var deleteTags = []string{
	"note",
	"source",
	"source_ref",
	"attribution",
	"comment",
	"fixme",
	"created_by",
	"odbl",
	"odbl:note",
	"project:eurosha_2012",
	"ref:UrbIS",
	"accuracy:meters",
	"sub_sea:type",
	"waterway:type",
	"statscan:rbuid",
	"ref:ruian:addr",
	"ref:ruian",
	"building:ruian:type",
	"dibavod:id",
	"uir_adr:ADRESA_KOD",
	"gst:feat_id",
	"maaamet:ETAK",
	"ref:FR:FANTOIR",
	"3dshapes:ggmodelk",
	"AND_nosr_r",
	"OPPDATERIN",
	"addr:city:simc",
	"addr:street:sym_ul",
	"building:usage:pl",
	"building:use:pl",
	"teryt:simc",
	"raba:id",
	"dcgis:gis_id",
	"nycdoitt:bin",
	"chicago:building_id",
	"lojic:bgnum",
	"massgis:way_id",
	"import",
	"import_uuid",
	"OBJTYPE",
	"SK53_bulk:load",
}

var deletePrefixes = []string{
	"note:",
	"source:",
	"CLC:",
	"geobase:",
	"canvec:",
	"geobase:",
	"osak:",
	"kms:",
	"ngbe:",
	"it:fvg:",
	"KSJ2:",
	"yh:",
	"LINZ2OSM:",
	"linz2osm:",
	"LINZ:",
	"WroclawGIS:",
	"naptan:",
	"tiger:",
	"gnis:",
	"NHD:",
	"nhd:",
	"mvdgis:",
}

const nodesToBuffer = 0

func keepTag(s string) bool {
	for _, f := range deleteTags {
		if f == s {
			return false
		}
	}
	for _, p := range deletePrefixes {
		if strings.HasPrefix(s, p) {
			return false
		}
	}
	return true
}

func sendBlock(sw *state.Worker, out chan<- pbf.Block, msg pbf.Block) {
	sw.SetCurrentState(state.WorkerWriting)
	out <- msg
	sw.SetCurrentState(state.WorkerReading)
}

func sendFiller(sw *state.Worker, outBlocks chan<- pbf.Block, index uint64) {
	sendBlock(sw, outBlocks, pbf.Block{
		Index:    index,
		Data:     nil,
		BlobType: "Filler",
	})
}

func (tf *TagFilterer) Worker() pbf.WorkFunc {
	return func(ctx context.Context, inBlocks <-chan pbf.Block, outBlocks chan<- pbf.Block) {
		sw := &state.Worker{}
		tf.TrackWorker(sw)
		defer tf.UntrackWorker(sw)

		pbr := t.PBReader{
			SkipStringTable: false,
			SkipDenseNodes:  false,
			SkipWays:        false,
			SkipRelations:   false,
		}

		var pendingNodes t.Nodes

		for data := range inBlocks {
			if data.BlobType != "OSMData" {
				sendBlock(sw, outBlocks, data)
				continue
			}
			sw.SetStateDecoding(data.Index)

			buffer := data.Data.(mmap.MMap)

			var b t.Blob
			pbr.ReadBlob(buffer, &b)

			var pb t.PrimitiveBlock
			pbr.ReadPrimitiveBlock(b.GetRawData().Buffer, &pb)

			sw.SetCurrentState(state.WorkerWorking)

			nodes := pb.ToNodes()
			ws := pb.ToWays()
			rs := pb.ToRelations()

			// filter
			idx := 0
			localKept := uint64(0)
			localDropped := uint64(0)
			for _, node := range nodes {
				for k, _ := range node.Tags {
					if !keepTag(k) {
						delete(node.Tags, k)
					}
				}
				if len(node.Tags) > 0 {
					nodes[idx] = node
					idx++
					localKept++
				} else {
					localDropped++
				}
			}
			nodes = nodes[:idx]

			sw.AddKept(localKept)
			sw.AddDropped(localDropped)

			if len(nodes)+len(pendingNodes) > nodesToBuffer || len(ws) > 0 || len(rs) > 0 {
				pbOut := t.MakePrimitiveBlock(pendingNodes, ws, rs)

				sw.SetCurrentState(state.WorkerEncoding)
				outputData := (&t.Blob{RawFunc: pbOut.Write}).Write(make([]byte, 0, 2*1048576))
				sendBlock(sw, outBlocks, pbf.Block{
					Index:    data.Index,
					Data:     outputData,
					BlobType: "OSMData",
				})

				pendingNodes = nodes
			} else {
				pendingNodes = append(pendingNodes, nodes...)
				sendFiller(sw, outBlocks, data.Index)
			}
		}

		if len(pendingNodes) > 0 {
			// This should only ever happen if the last block processed by this worker
			// contains only nodes (if it contains ways or relations it would have been
			// flushed).  It's possible with a small extract, but very unlikely.  If it
			// does happen, then it's small enough that the worker count can be set to
			// 1, which will prevent this case.
			//
			// It's also impossible to recover from, as the re-ordering logic is not
			// expecting an extra block.
			panic("pending nodes")
		}
	}
}
