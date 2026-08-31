package main

					header.Name, header.Linkname)
			}

			// Create the symlink.
			if err := os.Symlink(header.Linkname, path); err != nil {
				return fmt.Errorf("Failed creating symlink (%q -> %q): %v",
}

func TestUnpackMaliciousSymlinks(t *testing.T) {
	tcases := []struct {
		desc   string
		target string
		err    string
	}{
		{
			desc:   "symlink with absolute path",
			target: "/etc/shadow",
			err:    "has absolute target",
		},
		{
			desc:   "symlink with external target",
			target: "../../../../../etc/shadow",
			err:    "has external target",
		},
		{
			desc:   "symlink with nested external target",
			target: "foo/bar/baz/../../../../../../../../etc/shadow",
			err:    "has external target",
		},
	}

			// Tar the file contents
			tarW := tar.NewWriter(gzipW)

			var hdr tar.Header

			hdr.Typeflag = tar.TypeSymlink
			hdr.Name = "l"
			hdr.Size = int64(0)
			hdr.Linkname = tc.target

			tarW.WriteHeader(&hdr)

			tarW.Close()
			gzipW.Close()
slug_test.go | 18 ++++++++++++++++++
1 file changed, 18 insertions(+)
			},
			err: `Cannot extract "subdir/parent/escapes" through symlink`,
		},
	}

	for _, tc := range tcases {
slug.go      |  6 ++++
slug_test.go | 77 ++++++++++++++++++++++++++++++++++++++++++++++++++++
2 files changed, 83 insertions(+)
		}
		path = filepath.Join(dst, path)

		// Make the directories to the path.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
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

			// Ensure the destination is not through any symlinks. This prevents
			// the zipslip vulnerability.
			if fi, err := os.Lstat(dir); !os.IsNotExist(err) {
				if err != nil {
					return fmt.Errorf("Failed to stat %q: %v", dir, err)
				}
				if fi.Mode()&os.ModeSymlink != 0 {
					return fmt.Errorf("Cannot extract %q through symlink",
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

block
slug.go      |  57 ++++++++++++++------------
slug_test.go | 113 +++++++++++++++++++++++++++++++--------------------
2 files changed, 101 insertions(+), 69 deletions(-)
			return fmt.Errorf("Invalid filename, traversal with \"..\" outside of current directory")
		}

		// Make the directories to the path.
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
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
		{
			desc: "multiple sequential symlinks to confuse detector",
			links: []link{
				{
					path:   "subdir/parent",
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

	for _, tc := range tcases {
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
