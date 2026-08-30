package main


// readUvarint reads a uvarint from the given byte reader.
func readUvarint(r io.ByteReader) (x uint64, n int, err error) {
	const maxUvarintLen = 10

	var s uint
	i := 0
	for {
			return x, i, err
		}
		i++
		if i > maxUvarintLen {
			return x, i, errOverflowU64
		}
		if b < 0x80 {
			if i == maxUvarintLen && b > 1 {
				return x, i, errOverflowU64
			}
			return x | uint64(b)<<s, i, nil
		}
	}
}

func TestUvarIntCVE_2020_16845(t *testing.T) {
	var a = []byte{0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x8b}

	r := bytes.NewReader(a)
	_, _, err := readUvarint(r)
	if err != errOverflowU64 {
		t.Fatalf("readUvarint overflow not detected")
	}
}
