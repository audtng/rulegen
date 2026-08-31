package main

// Seal encrypts and authenticates the plaintext.
func (ctx *cbcAEAD) Seal(dst, nonce, plaintext, data []byte) []byte {
	// Output buffer -- must take care not to mangle plaintext input.
	ciphertext := make([]byte, uint64(len(plaintext))+uint64(ctx.Overhead()))[:len(plaintext)]
	copy(ciphertext, plaintext)
	ciphertext = padBuffer(ciphertext, ctx.blockCipher.BlockSize())

	cbc.CryptBlocks(ciphertext, ciphertext)
	authtag := ctx.computeAuthTag(data, nonce, ciphertext)

	ret, out := resize(dst, uint64(len(dst))+uint64(len(ciphertext))+uint64(len(authtag)))
	copy(out, ciphertext)
	copy(out[len(ciphertext):], authtag)

		return nil, err
	}

	ret, out := resize(dst, uint64(len(dst))+uint64(len(plaintext)))
	copy(out, plaintext)

	return ret, nil

// Compute an authentication tag
func (ctx *cbcAEAD) computeAuthTag(aad, nonce, ciphertext []byte) []byte {
	buffer := make([]byte, uint64(len(aad))+uint64(len(nonce))+uint64(len(ciphertext))+8)
	n := 0
	n += copy(buffer, aad)
	n += copy(buffer[n:], nonce)
	n += copy(buffer[n:], ciphertext)
	binary.BigEndian.PutUint64(buffer[n:], uint64(len(aad))*8)

	// According to documentation, Write() on hash.Hash never fails.
	hmac := hmac.New(ctx.hash, ctx.integrityKey)
// resize ensures the the given slice has a capacity of at least n bytes.
// If the capacity of the slice is less than n, a new slice is allocated
// and the existing data will be copied.
func resize(in []byte, n uint64) (head, tail []byte) {
	if uint64(cap(in)) >= n {
		head = in[:n]
	} else {
		head = make([]byte, n)
// Apply padding
func padBuffer(buffer []byte, blockSize int) []byte {
	missing := blockSize - (len(buffer) % blockSize)
	ret, out := resize(buffer, uint64(len(buffer))+uint64(missing))
	padding := bytes.Repeat([]byte{byte(missing)}, missing)
	copy(out, padding)
	return ret
	ct := aead.Seal(nil, nonce, data, nil)

	// Truncated ciphertext, but with correct auth tag
	truncated, tail := resize(ct[:len(ct)-ctx.authtagBytes-2], uint64(len(ct))-2)
	copy(tail, ctx.computeAuthTag(nil, nonce, truncated[:len(truncated)-ctx.authtagBytes]))

	// Open should fail
	ctx := aead.(*cbcAEAD)

	// Mutated ciphertext, but with correct auth tag
	size := uint64(len(buffer))
	ciphertext, tail := resize(buffer, size+(uint64(len(key))/2))
	copy(tail, ctx.computeAuthTag(nil, nonce, ciphertext[:size]))

	// Open should fail (b/c of invalid padding, even though tag matches)

// NewConcatKDF builds a KDF reader based on the given inputs.
func NewConcatKDF(hash crypto.Hash, z, algID, ptyUInfo, ptyVInfo, supPubInfo, supPrivInfo []byte) io.Reader {
	buffer := make([]byte, uint64(len(algID))+uint64(len(ptyUInfo))+uint64(len(ptyVInfo))+uint64(len(supPubInfo))+uint64(len(supPrivInfo)))
	n := 0
	n += copy(buffer, algID)
	n += copy(buffer[n:], ptyUInfo)

// DeriveECDHES derives a shared encryption key using ECDH/ConcatKDF as described in JWE/JWA.
// It is an error to call this function with a private/public key that are not on the same
// curve. Callers must ensure that the keys are valid before calling this function. Output
// size may be at most 1<<16 bytes (64 KiB).
func DeriveECDHES(alg string, apuData, apvData []byte, priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, size int) []byte {
	if size > 1<<16 {
		panic("ECDH-ES output size too large, must be less than 1<<16")
	}

	// algId, partyUInfo, partyVInfo inputs must be prefixed with the length
	algID := lengthPrefixed([]byte(alg))
	ptyUInfo := lengthPrefixed(apuData)
