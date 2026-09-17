package main

		}
		// strip tarball root dir
		_, name := utils.SplitByFirstByte(h.Name, '/')
		filename := path.Join(pkgDir, name)
		if h.Typeflag != tar.TypeReg {
			continue
		}
