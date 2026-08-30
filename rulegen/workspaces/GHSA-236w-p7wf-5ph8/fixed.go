package main

		if b[i] == 'P' && b[i+1] == 'K' && b[i+2] == 0x05 && b[i+3] == 0x06 {
			// n is length of comment
			n := int(b[i+directoryEndLen-2]) | int(b[i+directoryEndLen-1])<<8
			if n+directoryEndLen+i > len(b) {
				// Truncated comment.
				// Some parsers (such as Info-ZIP) ignore the truncated comment
				// rather than treating it as a hard error.
				return -1
			}
			return i
		}
	}
	return -1
			},
		},
	},
	// Issue 66869: Don't skip over an EOCDR with a truncated comment.
	// The test file sneakily hides a second EOCDR before the first one;
	// previously we would extract one file ("file") from this archive,
	// while most other tools would reject the file or extract a different one ("FILE").
	{
		Name:  "comment-truncated.zip",
		Error: ErrFormat,
	},
}

func TestReader(t *testing.T) {
