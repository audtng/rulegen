package main

					header.Name, header.Linkname)
			}

			// Ensure the destination is not through any symlinks. This prevents
			// the zipslip vulnerability.
			if fi, err := os.Lstat(dir); !os.IsNotExist(err) {
				if err != nil {
					return fmt.Errorf("Failed to stat %q: %v", dir, err)
				}
				if fi.Mode()&os.ModeSymlink != 0 {
					return fmt.Errorf("Cannot extract %q through symlink",
						header.Name)
				}
			}

			// Create the symlink.
			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("Failed creating symlink (%q -> %q): %v",
}

func TestUnpackMaliciousSymlinks(t *testing.T) {
	type link struct {
		path   string
		target string
	}

	tcases := []struct {
		desc  string
		links []link
		err   string
	}{
		{
			desc: "symlink with absolute path",
			links: []link{
				{
					path:   "l",
					target: "/etc/shadow",
				},
			},
			err: "has absolute target",
		},
		{
			desc: "symlink with external target",
			links: []link{
				{
					path:   "l",
					target: "../../../../../etc/shadow",
				},
			},
			err: "has external target",
		},
		{
			desc: "symlink with nested external target",
			links: []link{
				{
					path:   "l",
					target: "foo/bar/baz/../../../../../../../../etc/shadow",
				},
			},
			err: "has external target",
		},
		{
			desc: "zipslip vulnerability",
			links: []link{
				{
					path:   "subdir/parent",
					target: "..",
				},
				{
					path:   "subdir/parent/escapes",
					target: "..",
				},
			},
			err: `Cannot extract "subdir/parent/escapes" through symlink`,
		},
	}

			// Tar the file contents
			tarW := tar.NewWriter(gzipW)

			for _, link := range tc.links {
				var hdr tar.Header

				hdr.Typeflag = tar.TypeSymlink
				hdr.Name = link.path
				hdr.Size = int64(0)
				hdr.Linkname = link.target

				tarW.WriteHeader(&hdr)
			}

			tarW.Close()
			gzipW.Close()
slug_test.go | 18 ++++++++++++++++++
1 file changed, 18 insertions(+)
			},
			err: `Cannot extract "subdir/parent/escapes" through symlink`,
		},
		{
			desc: "multiple sequential symlinks to confuse detector",
			links: []link{
				{
					path:   "subdir/parent",
					target: "..",
				},
				{
					path:   "subdir/parent/escape",
					target: "../..",
				},
				{
					path:   "subdir/parent/escape/root",
					target: "../../..",
				},
			},
			err: `Cannot extract "subdir/parent/escape" through symlink`,
		},
	}

	for _, tc := range tcases {
slug.go      |  6 ++++
slug_test.go | 77 ++++++++++++++++++++++++++++++++++++++++++++++++++++
2 files changed, 83 insertions(+)
		}
		path = filepath.Join(dst, path)

		// Check for paths outside our directory, they are forbidden
		target := filepath.Clean(path)
		if !strings.HasPrefix(target, dst) {
			return fmt.Errorf("Invalid filename, traversal with \"..\" outside of current directory")
		}

		// Make the directories to the path.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
	}
}

func TestUnpackMaliciousFiles(t *testing.T) {
	tcases := []struct {
		desc string
		name string
		err  string
	}{
		{
			desc: "filename containing path traversal",
			name: "../../../../../../../../tmp/test",
			err:  "Invalid filename, traversal with \"..\" outside of current directory",
		},
		{
			desc: "should fail before attempting to create directories",
			name: "../../../../../../../../Users/root",
			err:  "Invalid filename, traversal with \"..\" outside of current directory",
		},
	}

	for _, tc := range tcases {
		t.Run(tc.desc, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "slug")
			if err != nil {
				t.Fatalf("err:%v", err)
			}
			defer os.RemoveAll(dir)
			in := filepath.Join(dir, "slug.tar.gz")

			// Create the output file
			wfh, err := os.Create(in)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			// Gzip compress all the output data
			gzipW := gzip.NewWriter(wfh)

			// Tar the file contents
			tarW := tar.NewWriter(gzipW)

			hdr := &tar.Header{
				Name: tc.name,
				Mode: 0600,
				Size: int64(0),
			}
			if err := tarW.WriteHeader(hdr); err != nil {
				t.Fatalf("err: %v", err)
			}
			if _, err := tarW.Write([]byte{}); err != nil {
				t.Fatalf("err: %v", err)
			}

			tarW.Close()
			gzipW.Close()
			wfh.Close()

			// Open the slug file for reading.
			fh, err := os.Open(in)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			// Create a dir to unpack into.
			dst, err := ioutil.TempDir(dir, "")
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			defer os.RemoveAll(dst)

			// Now try unpacking it, which should fail
			err = Unpack(fh, dst)
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}
}

func TestCheckFileMode(t *testing.T) {
	for _, tc := range []struct {
		desc string
through other symlinks
slug.go      | 23 ++++++++++++++++++-----
slug_test.go |  8 ++------
2 files changed, 20 insertions(+), 11 deletions(-)
					header.Name, header.Linkname)
			}

			// Ensure the destination is not through any symlinks. This
			// prevents the zipslip vulnerability.
			//
			// The strategy is to Lstat each path  component from dst up to the
			// immediate parent directory of the file name in the tarball,
			// checking the mode on each to ensure we wouldn't be passing
			// through any symlinks.
			currentPath := dst // Start at the root of the unpacked tarball.
			components := strings.Split(header.Name, "/")

			for i := 0; i < len(components)-1; i++ {
				currentPath = filepath.Join(currentPath, components[i])
				fi, err := os.Lstat(currentPath)
				if os.IsNotExist(err) {
					continue
				}
				if err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("Failed to evaluate symlink: %v", err)
				}
				if fi.Mode()&os.ModeSymlink != 0 {
					return fmt.Errorf("Cannot extract %q through symlink",
					target: "..",
				},
				{
					path:   "subdir/parent/otherdir/escapes",
					target: "../..",
				},
			},
			err: `Cannot extract "subdir/parent/otherdir/escapes" through symlink`,
		},
	}

block
slug.go      |  57 ++++++++++++++------------
slug_test.go | 113 +++++++++++++++++++++++++++++++--------------------
2 files changed, 101 insertions(+), 69 deletions(-)
			return fmt.Errorf("Invalid filename, traversal with \"..\" outside of current directory")
		}

		// Ensure the destination is not through any symlinks. This prevents
		// any files from being deployed through symlinks defined in the slug.
		// There are malicious cases where this could be used to escape the
		// slug's boundaries (zipslip), and any legitimate use is questionable
		// and likely indicates a hand-crafted tar file, which we are not in
		// the business of supporting here.
		//
		// The strategy is to Lstat each path  component from dst up to the
		// immediate parent directory of the file name in the tarball, checking
		// the mode on each to ensure we wouldn't be passing through any
		// symlinks.
		currentPath := dst // Start at the root of the unpacked tarball.
		components := strings.Split(header.Name, "/")

		for i := 0; i < len(components)-1; i++ {
			currentPath = filepath.Join(currentPath, components[i])
			fi, err := os.Lstat(currentPath)
			if os.IsNotExist(err) {
				// Parent directory structure is incomplete. Technically this
				// means from here upward cannot be a symlink, so we cancel the
				// remaining path tests.
				break
			}
			if err != nil {
				return fmt.Errorf("Failed to evaluate path %q: %v", header.Name, err)
			}
			if fi.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("Cannot extract %q through symlink",
					header.Name)
			}
		}

		// Make the directories to the path.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
					header.Name, header.Linkname)
			}

			// Create the symlink.
			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("Failed creating symlink (%q -> %q): %v",
}

func TestUnpackMaliciousSymlinks(t *testing.T) {
	tcases := []struct {
		desc    string
		headers []*tar.Header
		err     string
	}{
		{
			desc: "symlink with absolute path",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "l",
					Linkname: "/etc/shadow",
					Typeflag: tar.TypeSymlink,
				},
			},
			err: "has absolute target",
		},
		{
			desc: "symlink with external target",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "l",
					Linkname: "../../../../../etc/shadow",
					Typeflag: tar.TypeSymlink,
				},
			},
			err: "has external target",
		},
		{
			desc: "symlink with nested external target",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "l",
					Linkname: "foo/bar/baz/../../../../../../../../etc/shadow",
					Typeflag: tar.TypeSymlink,
				},
			},
			err: "has external target",
		},
		{
			desc: "zipslip vulnerability",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "subdir/parent",
					Linkname: "..",
					Typeflag: tar.TypeSymlink,
				},
				&tar.Header{
					Name:     "subdir/parent/escapes",
					Linkname: "..",
					Typeflag: tar.TypeSymlink,
				},
			},
			err: `Cannot extract "subdir/parent/escapes" through symlink`,
		},
		{
			desc: "nested symlinks within symlinked dir",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "subdir/parent",
					Linkname: "..",
					Typeflag: tar.TypeSymlink,
				},
				&tar.Header{
					Name:     "subdir/parent/otherdir/escapes",
					Linkname: "../..",
					Typeflag: tar.TypeSymlink,
				},
			},
			err: `Cannot extract "subdir/parent/otherdir/escapes" through symlink`,
		},
		{
			desc: "regular file through symlink",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "subdir/parent",
					Linkname: "..",
					Typeflag: tar.TypeSymlink,
				},
				&tar.Header{
					Name:     "subdir/parent/file",
					Typeflag: tar.TypeReg,
				},
			},
			err: `Cannot extract "subdir/parent/file" through symlink`,
		},
		{
			desc: "directory through symlink",
			headers: []*tar.Header{
				&tar.Header{
					Name:     "subdir/parent",
					Linkname: "..",
					Typeflag: tar.TypeSymlink,
				},
				&tar.Header{
					Name:     "subdir/parent/dir",
					Typeflag: tar.TypeDir,
				},
			},
			err: `Cannot extract "subdir/parent/dir" through symlink`,
		},
	}

	for _, tc := range tcases {
			// Tar the file contents
			tarW := tar.NewWriter(gzipW)

			for _, hdr := range tc.headers {
				tarW.WriteHeader(hdr)
			}

			tarW.Close()
