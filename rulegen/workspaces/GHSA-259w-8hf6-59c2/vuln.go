package main

	return writeManifest(ctx, store, idx, ocispec.MediaTypeImageIndex)
}

func onUntarJSON(r io.Reader, j interface{}) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, j)
}

func onUntarBlob(ctx context.Context, r io.Reader, store content.Ingester, size int64, ref string) (digest.Digest, error) {
