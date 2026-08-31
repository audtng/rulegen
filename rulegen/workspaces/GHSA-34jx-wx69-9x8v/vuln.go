package main

	// remove extraneous path shortcuts - these could occur if a path contained extra "../"
	// and attempted to navigate beyond "/" in a remote filesystem
	prefix = stripPathShortcuts(prefix)
	return untarAll(reader, dest.File, prefix)
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
	return path.Clean(string(os.PathSeparator) + fileName)
}

func untarAll(reader io.Reader, destFile, prefix string) error {
	entrySeq := -1

	// TODO: use compression here?
		}
		entrySeq++
		mode := header.FileInfo().Mode()
		outFileName := path.Join(destFile, clean(header.Name[len(prefix):]))
		baseName := path.Dir(outFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
		}

		if mode&os.ModeSymlink != 0 {
			err := os.Symlink(header.Linkname, outFileName)
			if err != nil {
				return err
			}
		} else {
	}
}

func TestTarUntar(t *testing.T) {
	dir, err := ioutil.TempDir("", "input")
	dir2, err2 := ioutil.TempDir("", "output")
	if err != nil || err2 != nil {
		t.Errorf("unexpected error: %v | %v", err, err2)
		t.FailNow()
	}
	dir = dir + "/"
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("Unexpected error cleaning up: %v", err)
		}
		if err := os.RemoveAll(dir2); err != nil {
			t.Errorf("Unexpected error cleaning up: %v", err)
		}
	}()

	files := []struct {
		name     string
		nameList []string
		data     string
		fileType FileType
	}{
		{
		},
		{
			name:     "gakki",
			data:     "/tmp/gakki",
			fileType: SymLink,
		},
		{
		}
	}

	writer := &bytes.Buffer{}
	if err := makeTar(dir, dir, writer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := bytes.NewBuffer(writer.Bytes())
	if err := untarAll(reader, dir2, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

			cmpFileData(t, filePath, file.data)
		} else if file.fileType == SymLink {
			dest, err := os.Readlink(filePath)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
	}
}

// TestCopyToLocalFileOrDir tests untarAll in two cases :
// 1: copy pod file to local file
// 2: copy pod file into local directory
			}
			defer srcTarFile.Close()

			if err := untarAll(srcTarFile, destPath, getPrefix(srcFilePath)); err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}
		t.FailNow()
	}

	if err := untarAll(&buf, dir, "/prefix"); err != nil {
		t.Errorf("unexpected error: %v ", err)
		t.FailNow()
	}
