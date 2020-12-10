build:
	go build ./cmd/blockcount/
	go build ./cmd/decompress/
	go build ./cmd/dumpnodeindex/
	go build ./cmd/dumpwayindex/
	go build ./cmd/entitycount/
	go build ./cmd/entityfind/
	go build ./cmd/nodeindex/
	go build ./cmd/relindex/
	go build ./cmd/rewritenodes/
	go build ./cmd/tagfilter/
	go build ./cmd/wayindex/

test:
	go test ./...
