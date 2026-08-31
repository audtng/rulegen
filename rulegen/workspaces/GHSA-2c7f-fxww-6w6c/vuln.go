package main

		return errors.Join(errDownloadCaption, err)
	}

	file, err := os.Create(c.File)
	if err != nil {
		return errors.Join(errDownloadCaption, err)
	}
