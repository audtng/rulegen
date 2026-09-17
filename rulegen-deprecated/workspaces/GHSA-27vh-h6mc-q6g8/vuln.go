package main

// pre-segwit, segwit v0, segwit v1 (taproot key spend validation), and the
// base tapscript verification.
type signatureVerifier interface {
	// Verify returns true if the signature verifier context deems the
	// signature to be valid for the given context.
	Verify() bool
}

// baseSigVerifier is used to verify signatures for the _base_ system, meaning
	return valid
}

// Verify returns true if the signature verifier context deems the signature to
// be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (b *baseSigVerifier) Verify() bool {
	// Remove the signature since there is no way for a signature
	// to sign itself.
	subScript := removeOpcodeByData(b.subScript, b.fullSigBytes)

	sigHash := calcSignatureHash(
		subScript, b.hashType, &b.vm.tx, b.vm.txIdx,
	)

	return b.verifySig(sigHash)
}

// A compile-time assertion to ensure baseSigVerifier implements the
// be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (s *baseSegwitSigVerifier) Verify() bool {
	var sigHashes *TxSigHashes
	if s.vm.hashCache != nil {
		sigHashes = s.vm.hashCache
		// TODO(roasbeef): this doesn't need to return an error, should
		// instead be further up the stack? this only returns an error
		// if the input index is greater than the number of inputs
		return false
	}

	return s.verifySig(sigHash)
}

// A compile-time assertion to ensure baseSegwitSigVerifier implements the
	return false
}

// Verify returns true if the signature verifier context deems the signature to
// be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (t *taprootSigVerifier) Verify() bool {
	var opts []TaprootSigHashOption
	if t.annex != nil {
		opts = append(opts, WithAnnex(t.annex))
	)
	if err != nil {
		// TODO(roasbeef): propagate the error here?
		return false
	}

	return t.verifySig(sigHash)
}

// A compile-time assertion to ensure taprootSigVerifier implements the
	}
}

// Verify returns true if the signature verifier context deems the signature to
// be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (b *baseTapscriptSigVerifier) Verify() bool {
	// If the public key is blank, then that means it wasn't 0 or 32 bytes,
	// so we'll treat this as an unknown public key version and return
	// true.
	if b.pubKey == nil {
		return true
	}

	var opts []TaprootSigHashOption
	)
	if err != nil {
		// TODO(roasbeef): propagate the error here?
		return false
	}

	return b.verifySig(sigHash)
}

// A compile-time assertion to ensure baseTapscriptSigVerifier implements the
// NOTE: This function is only valid for version 0 scripts.  Since the function
// does not accept a script version, the results are undefined for other script
// versions.
func removeOpcodeByData(script []byte, dataToRemove []byte) []byte {
	// Avoid work when possible.
	if len(script) == 0 || len(dataToRemove) == 0 {
		return script
	}

	// Parse through the script looking for a canonical data push that contains
	const scriptVersion = 0
	var result []byte
	var prevOffset int32
	tokenizer := MakeScriptTokenizer(scriptVersion, script)
	for tokenizer.Next() {
		// In practice, the script will basically never actually contain the
		// data since this function is only used during signature verification
		// to remove the signature itself which would require some incredibly
		// non-standard code to create.
		//
		// Thus, as an optimization, avoid allocating a new script unless there
		// is actually a match that needs to be removed.
		op, data := tokenizer.Opcode(), tokenizer.Data()
		if isCanonicalPush(op, data) && bytes.Contains(data, dataToRemove) {
			if result == nil {
				fullPushLen := tokenizer.ByteIndex() - prevOffset
				result = make([]byte, 0, int32(len(script))-fullPushLen)
				result = append(result, script[0:prevOffset]...)
			}
		} else if result != nil {
			result = append(result, script[prevOffset:tokenizer.ByteIndex()]...)
		}

		prevOffset = tokenizer.ByteIndex()
	}
	if result == nil {
		result = script
	}
	return result
}

// AsSmallInt returns the passed opcode, which must be true according to

	if vm.taprootCtx != nil {
		vm.taprootCtx.codeSepPos = uint32(vm.tokenizer.OpcodePosition())
	}

	return nil
		// TODO(roasbeef): return an error?
	}

	valid := sigVerifier.Verify()

	switch {
	// For tapscript, and prior execution with null fail active, if the
		return err
	}

	valid := sigVerifier.Verify()

	// If the signature is invalid, this we fail execution, as it should
	// have been an empty signature.
	if !valid {
		str := "signature not empty on failed checksig"
		return scriptError(ErrNullFail, str)
	}
	// no way for a signature to sign itself.
	if !vm.isWitnessVersionActive(0) {
		for _, sigInfo := range signatures {
			script = removeOpcodeByData(script, sigInfo.signature)
		}
	}

