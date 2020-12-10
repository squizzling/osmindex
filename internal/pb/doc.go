package pb

/*
There are a number of PB coding functions, all function have one of 6 types:
- U32: Unsigned integer, 32bits
- U64: Unsigned integer, 64bits
- I32: Signed integer, 32bits
- I64: Signed integer, 64bits
- S32: Zigzag signed integer, 32bits
- S64: Zigzag signed integer, 64bits

All integers are decoded as 64bit, and then cast to the appropriate type.

A "buffer" in the context of this is a length prefixed sequence of values.

All functions then take one of the following forms:
- {Encode,Decode}{type} - These will encode or decode a single integer of the suitable type,

- Decode{type}Opt - These are identical to the above, however they return a pointer to the decoded value.

- {Encode,Decode}{type}Packed - These code a buffer, to/from a slice of the provided type.  Encoding will allocate additional temporary memory, so it may be advantageous to use a *Func style function instead.

- {Encode,Decode}{type}PackedDelta - These are the same as the above, except each slice element is the delta of the previous one before/after coding.

- {Encode,Decode}{type}PackedDeltaZero - These are similar to the PackedDelta form, except it's not coded as a buffer, but rather terminated with a 0.  A slice with sequential duplicate values will not work.

- Encode{type}Packed[Delta]Func - These take a slice, and return a function which will encode the slice (using the above rules, except it doesn't encode a buffer) in to a byte slice passed to it.  Used to avoid allocation.

Not every function is implemented, only the ones required.
*/
