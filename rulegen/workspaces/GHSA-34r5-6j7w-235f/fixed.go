package main

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/inspektor-gadget/inspektor-gadget/pkg/columns"
		return ""
	}

	// Escape the string to avoid the terminal to interpret escape sequence when
	// printing and remove first and last character as they are quotes added by
	// strconv.Quote().
	s = strconv.Quote(s)
	s = s[1 : len(s)-1]
	if !tf.options.ShouldTruncate {
		return s
	}
