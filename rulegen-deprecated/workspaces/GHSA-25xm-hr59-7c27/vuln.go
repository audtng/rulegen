package main


// readUvarint reads a uvarint from the given byte reader.
func readUvarint(r io.ByteReader) (x uint64, n int, err error) {
	var s uint
	i := 0
	for {
			return x, i, err
		}
		i++
		if b < 0x80 {
			if i > 10 || i == 10 && b > 1 {
				return x, i, errOverflowU64
			}
			return x | uint64(b)<<s, i, nil
