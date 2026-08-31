package main

		return errors.Join(errDownloadCaption, err)
	}

	file, err := pkg.Root.OpenFile(c.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}
