package main

		}
		// strip tarball root dir
		_, name := utils.SplitByFirstByte(h.Name, '/')
		filename := path.Join(pkgDir, path.Clean(name))
		if h.Typeflag != tar.TypeReg {
			continue
		}
