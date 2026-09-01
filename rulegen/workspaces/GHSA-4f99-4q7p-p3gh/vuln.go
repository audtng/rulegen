package main

	"bufio"
	"io"
	"runtime"
)

// Writer at INFO level. See WriterLevel for details.
	return NewEntry(logger).WriterLevel(level)
}

func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

func (entry *Entry) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})

	switch level {
	case TraceLevel:
		printFunc = entry.Trace
		printFunc = entry.Print
	}

	go entry.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		entry.Errorf("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
writer.go | 3 ++-
1 file changed, 2 insertions(+), 1 deletion(-)
		if len(data) > chunkSize {
			return chunkSize, data[:chunkSize], nil
		}
		return 0, nil, nil
	}

	//Use the custom split function to split the input
