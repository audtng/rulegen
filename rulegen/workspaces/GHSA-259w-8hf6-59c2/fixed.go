package main

	return writeManifest(ctx, store, idx, ocispec.MediaTypeImageIndex)
}

const (
	kib       = 1024
	mib       = 1024 * kib
	jsonLimit = 20 * mib
)

func onUntarJSON(r io.Reader, j interface{}) error {
	return json.NewDecoder(io.LimitReader(r, jsonLimit)).Decode(j)
}

func onUntarBlob(ctx context.Context, r io.Reader, store content.Ingester, size int64, ref string) (digest.Digest, error) {
