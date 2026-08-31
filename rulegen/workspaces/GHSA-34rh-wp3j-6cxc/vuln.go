package main

	if err != nil {
		return err
	}
	family := bin.Uint16(value[0:2])
	if family != familyIPv6 && family != familyIPv4 {
		return newDecodeErr("xor-mapped address", "family",
		}
	}

	if len(value) <= 4 {
		return io.ErrUnexpectedEOF
	}
	if err := CheckOverflow(attr, len(value[4:]), len(a.IP)); err != nil {
		return err
	}
