package main

	// remove extraneous path shortcuts - these could occur if a path contained extra "../"
	// and attempted to navigate beyond "/" in a remote filesystem
	prefix = stripPathShortcuts(prefix)
	return o.untarAll(reader, dest.File, prefix)
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
	return path.Clean(string(os.PathSeparator) + fileName)
}

func (o *CopyOptions) untarAll(reader io.Reader, destFile, prefix string) error {
	entrySeq := -1

	// TODO: use compression here?
		}
		entrySeq++
		mode := header.FileInfo().Mode()
		// all the files will start with the prefix, which is the directory where
		// they were located on the pod, we need to strip down that prefix, but
		// if the prefix is missing it means the tar was tempered with
		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}
		outFileName := path.Join(destFile, clean(header.Name[len(prefix):]))
		baseName := path.Dir(outFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname
			// error is returned if linkname can't be made relative to destFile,
			// but relative can end up being ../dir that's why we also need to
			// verify if relative path is the same after Clean-ing
			relative, err := filepath.Rel(destFile, linkname)
			if path.IsAbs(linkname) && (err != nil || relative != stripPathShortcuts(relative)) {
				fmt.Fprintf(o.IOStreams.ErrOut, "warning: link %q is pointing to %q which is outside target destination, skipping\n", outFileName, header.Linkname)
				continue
			}
			if err := os.Symlink(linkname, outFileName); err != nil {
				return err
			}
		} else {
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		t.FailNow()
	}
}

func TestTarUntar(t *testing.T) {
	dir, err := ioutil.TempDir("", "input")
	checkErr(t, err)
	dir2, err := ioutil.TempDir("", "output")
	checkErr(t, err)
	dir3, err := ioutil.TempDir("", "dir")
	checkErr(t, err)

	dir = dir + "/"
	defer func() {
		os.RemoveAll(dir)
		os.RemoveAll(dir2)
		os.RemoveAll(dir3)
	}()

	files := []struct {
		name     string
		nameList []string
		data     string
		omitted  bool
		fileType FileType
	}{
		{
		},
		{
			name:     "gakki",
			data:     "tmp/gakki",
			fileType: SymLink,
		},
		{
			name:     "relative_to_dest",
			data:     path.Join(dir2, "foo"),
			fileType: SymLink,
		},
		{
			name:     "tricky_relative",
			data:     path.Join(dir3, "xyz"),
			omitted:  true,
			fileType: SymLink,
		},
		{
			name:     "absolute_path",
			data:     "/tmp/gakki",
			omitted:  true,
			fileType: SymLink,
		},
		{
		}
	}

	opts := NewCopyOptions(genericclioptions.NewTestIOStreamsDiscard())

	writer := &bytes.Buffer{}
	if err := makeTar(dir, dir, writer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := bytes.NewBuffer(writer.Bytes())
	if err := opts.untarAll(reader, dir2, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

			cmpFileData(t, filePath, file.data)
		} else if file.fileType == SymLink {
			dest, err := os.Readlink(filePath)
			if file.omitted {
				if err != nil && strings.Contains(err.Error(), "no such file or directory") {
					continue
				}
				t.Fatalf("expected to omit symlink for %s", filePath)
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
	}
}

func TestTarUntarWrongPrefix(t *testing.T) {
	dir, err := ioutil.TempDir("", "input")
	checkErr(t, err)
	dir2, err := ioutil.TempDir("", "output")
	checkErr(t, err)

	dir = dir + "/"
	defer func() {
		os.RemoveAll(dir)
		os.RemoveAll(dir2)
	}()

	filepath := path.Join(dir, "foo")
	if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	createTmpFile(t, filepath, "sample data")

	opts := NewCopyOptions(genericclioptions.NewTestIOStreamsDiscard())

	writer := &bytes.Buffer{}
	if err := makeTar(dir, dir, writer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := bytes.NewBuffer(writer.Bytes())
	err = opts.untarAll(reader, dir2, "verylongprefix-showing-the-tar-was-tempered-with")
	if err == nil || !strings.Contains(err.Error(), "tar contents corrupted") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestCopyToLocalFileOrDir tests untarAll in two cases :
// 1: copy pod file to local file
// 2: copy pod file into local directory
			}
			defer srcTarFile.Close()

			opts := NewCopyOptions(genericclioptions.NewTestIOStreamsDiscard())
			if err := opts.untarAll(srcTarFile, destPath, getPrefix(srcFilePath)); err != nil {
				t.Errorf("unexpected error: %v", err)
				t.FailNow()
			}
		t.FailNow()
	}

	opts := NewCopyOptions(genericclioptions.NewTestIOStreamsDiscard())
	if err := opts.untarAll(&buf, dir, "/prefix"); err != nil {
		t.Errorf("unexpected error: %v ", err)
		t.FailNow()
	}
