# Endian
All raw values are little endian, and varints are don't have an endianness.  There are no allowances or consideration for running on big endian architectures in the codebase.

# Morton locations
Each location is 8 bytes.  A location is a [Morton encoded](https://en.wikipedia.org/wiki/Morton_code) pair of latitude and longitude, the values are multiplied by 10,000,000, eg -123.456,789 becomes -1,234,567,890.

Before Morton encoding the latitude is also shifted right by 1, and the high bit set to 1.  This forces the resulting int64 value to always be negative, and we use this property to differentiate between an ID (positive) and encoded location (negative).

The encoded latitude is in the odd bits, and the longitude is in the even bits.

```
mortion_location = MortonEncode(
  even = longitude * 10_000_000,
  odd = ((latitude * 10_000_000) >> 1) | 0x8000_0000,
)
```

This encoding is used across all indexes.  It is used because locations are frequently near each other, and therefore a run of nearby locations can be efficiently encoded as a varint deltas on the previous location, although this varint delta encoding is only used for the `.way.idx`.  The range of latitude is also -90 to 90, which means only 31 bits will be required, allowing us to use the extra bit as a flag.

# `.nodeidx` encoding

## version 1
A `.nodeidx` has 3 parts, in order

### Part 1 (header)
- 8 bytes defines the version
- 8 bytes defines the number of index elements

### Part 2 (index elements)
Each element is 16 bytes:
- 8 bytes defines the start id of a range
- 8 bytes encodes the number of ids in a range, and the start offset of the location data
-- The high 40 bits defines the offset in to the location data (part 3)
-- The low 24 bits defines the count

The IDs must be in ascending order, however this is not enforced and lookup will fail if they are not.  The `.osm.pbf` source data must be in ascending order. 

### Part 3 (locations)
Each location is either `0`, or a Morton location.  There is no varint encoding, each location is 64 bits.

`0` indicates that the node was not present in the source `.osm.pbf`, and allows for there to be less index elements, at the expense of more locations. 

## Possible enhancements

### Smaller index elements
Index elements could be encoded with just 64bits, if the encoding is specified as:

- id:     33 bits, up to 8 billion IDs, currently the maximum is around 6 billion
- count:  18 bits, up to 262144 per block
- offset: 13 bits, each block must contain 8192 locations, even if there is only a single valid location in it

The idea of `offset` referring to a larger block is how we encode the `.relidx`, however in `.nodeidx`, it is likely unnecessary as the disk/memory required for the `planet` index elements is only 1.7GB.  The bulk of the size is in the location data (49GB), and it is this aspect which is used to manage the RSS

### varint delta encoding blocks
By varint delta encoding blocks, we may be able to trade decoding time for better memory usage.

### Copying less data 
Swapping index elements and locations would make it much faster to generate the index, as we generate them both in parallel and then copy the locations to the end of the index elements

# `.wayidx` encoding

## Part 1 (header)
- 8 bytes defines the number of index data elements
- 8 bytes for the alignment shift

## Part 2 (index data)
- 8 bytes of ID and offset

The top 31 bits define the Way ID, the bottom 33 bits define the block offset.  The block offset is the number of blocks in to the location blocks to find the start of the location data. 

This encoding only allows for a maximum Way ID of approximately 2 billion, at present the maximum is around 844 million, or 40%.

## Part 3 (location blocks)
Each location block is:
- A list of varint delta Morton encoded locations
- A varint `0` (one byte)
- Padding to align the block to 2 ** alignmentShift bytes

The maximum number of bytes in the entire location block section is 2**(33 + alignmentShift).

With no padding, up to 8GB of location data can be encoded.  With 0-7 byte padding (the default), up to 68GB of location data can be encoded.  This is sufficient for `planet`.  

## Possible enhancements

### A version
The header has plenty of free space, and a version number should be used.

### Location blocks first
Putting the location blocks ahead of the index data would make index generation faster.

### Different bit allocations for id / offset 
The ID needs to fit the maximum way ID, however the offset could be changed to require only the count of ways by making the block size large enough to fit the way with the most locations in it.  The current split allows for smaller files at the expense of tighter limits.

The largest way is 2000 entries  

# `.relidx` encoding
A sequence of RelationIndex.

Each RelationIndex is:
- 64bit RelationID
- 64bit WayCount
- Sequence of WayIndex, counted by WayCount

Each WayIndex is:
- 64bit LocationCount (may be 0)
- Sequence of 64bit Morton encoded Locations, counted by LocationCount 

This is output in the order it appears in the file, which is the order osm2pgsql will process the relations.  Some relations will be filtered by osm2pgsq, but are still present in the `.relidx`, so osm2pgsql will need to scan forward until it finds a match.  There is no encoding to reduce the file size, because it reduces the complexity of code in osm2pgsql.

