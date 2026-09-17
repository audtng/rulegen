package main

// pre-segwit, segwit v0, segwit v1 (taproot key spend validation), and the
// base tapscript verification.
type signatureVerifier interface {
	// Verify returns whether or not the signature verifier context deems the
	// signature to be valid for the given context.
	Verify() verifyResult
}

type verifyResult struct {
	sigValid bool
	sigMatch bool
}

// baseSigVerifier is used to verify signatures for the _base_ system, meaning
	return valid
}

// Verify returns whether or not the signature verifier context deems the
// signature to be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (b *baseSigVerifier) Verify() verifyResult {
	// Remove the signature since there is no way for a signature
	// to sign itself.
	subScript, match := removeOpcodeByData(b.subScript, b.fullSigBytes)

	sigHash := calcSignatureHash(
		subScript, b.hashType, &b.vm.tx, b.vm.txIdx,
	)

	return verifyResult{
		sigValid: b.verifySig(sigHash),
		sigMatch: match,
	}
}

// A compile-time assertion to ensure baseSigVerifier implements the
// be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (s *baseSegwitSigVerifier) Verify() verifyResult {
	var sigHashes *TxSigHashes
	if s.vm.hashCache != nil {
		sigHashes = s.vm.hashCache
		// TODO(roasbeef): this doesn't need to return an error, should
		// instead be further up the stack? this only returns an error
		// if the input index is greater than the number of inputs
		return verifyResult{}
	}

	return verifyResult{
		sigValid: s.verifySig(sigHash),
	}
}

// A compile-time assertion to ensure baseSegwitSigVerifier implements the
	return false
}

// Verify returns whether or not the signature verifier context deems the
// signature to be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (t *taprootSigVerifier) Verify() verifyResult {
	var opts []TaprootSigHashOption
	if t.annex != nil {
		opts = append(opts, WithAnnex(t.annex))
	)
	if err != nil {
		// TODO(roasbeef): propagate the error here?
		return verifyResult{}
	}

	return verifyResult{
		sigValid: t.verifySig(sigHash),
	}
}

// A compile-time assertion to ensure taprootSigVerifier implements the
	}
}

// Verify returns whether or not the signature verifier context deems the
// signature to be valid for the given context.
//
// NOTE: This is part of the baseSigVerifier interface.
func (b *baseTapscriptSigVerifier) Verify() verifyResult {
	// If the public key is blank, then that means it wasn't 0 or 32 bytes,
	// so we'll treat this as an unknown public key version and return
	// that it's valid.
	if b.pubKey == nil {
		return verifyResult{
			sigValid: true,
		}
	}

	var opts []TaprootSigHashOption
	)
	if err != nil {
		// TODO(roasbeef): propagate the error here?
		return verifyResult{}
	}

	return verifyResult{
		sigValid: b.verifySig(sigHash),
	}
}

// A compile-time assertion to ensure baseTapscriptSigVerifier implements the
// NOTE: This function is only valid for version 0 scripts.  Since the function
// does not accept a script version, the results are undefined for other script
// versions.
func removeOpcodeByData(script []byte, dataToRemove []byte) ([]byte, bool) {
	// Avoid work when possible.
	if len(script) == 0 || len(dataToRemove) == 0 {
		return script, false
	}

	// Parse through the script looking for a canonical data push that contains
	const scriptVersion = 0
	var result []byte
	var prevOffset int32
	var match bool
	tokenizer := MakeScriptTokenizer(scriptVersion, script)
	for tokenizer.Next() {
		var found bool
		result, prevOffset, found = removeOpcodeCanonical(
			&tokenizer, script, dataToRemove, prevOffset, result,
		)
		if found {
			match = true
		}
	}
	if result == nil {
		result = script
	}
	return result, match
}

func removeOpcodeCanonical(t *ScriptTokenizer, script, dataToRemove []byte,
	prevOffset int32, result []byte) ([]byte, int32, bool) {

	var found bool

	// In practice, the script will basically never actually contain the
	// data since this function is only used during signature verification
	// to remove the signature itself which would require some incredibly
	// non-standard code to create.
	//
	// Thus, as an optimization, avoid allocating a new script unless there
	// is actually a match that needs to be removed.
	op, data := t.Opcode(), t.Data()
	if isCanonicalPush(op, data) && bytes.Equal(data, dataToRemove) {
		if result == nil {
			fullPushLen := t.ByteIndex() - prevOffset
			result = make([]byte, 0, int32(len(script))-fullPushLen)
			result = append(result, script[0:prevOffset]...)
		}
		found = true
	} else if result != nil {
		result = append(result, script[prevOffset:t.ByteIndex()]...)
	}

	return result, t.ByteIndex(), found
}

// AsSmallInt returns the passed opcode, which must be true according to

	if vm.taprootCtx != nil {
		vm.taprootCtx.codeSepPos = uint32(vm.tokenizer.OpcodePosition())
	} else if vm.witnessProgram == nil &&
		vm.hasFlag(ScriptVerifyConstScriptCode) {

		// Disable OP_CODESEPARATOR for non-segwit scripts.
		str := "OP_CODESEPARATOR used in non-segwit script"
		return scriptError(ErrCodeSeparator, str)
	}

	return nil
		// TODO(roasbeef): return an error?
	}

	result := sigVerifier.Verify()
	valid := result.sigValid

	if vm.hasFlag(ScriptVerifyConstScriptCode) && result.sigMatch {
		str := "non-const script code"
		return scriptError(ErrNonConstScriptCode, str)
	}

	switch {
	// For tapscript, and prior execution with null fail active, if the
		return err
	}

	result := sigVerifier.Verify()

	// If the signature is invalid, this we fail execution, as it should
	// have been an empty signature.
	if !result.sigValid {
		str := "signature not empty on failed checksig"
		return scriptError(ErrNullFail, str)
	}
	// no way for a signature to sign itself.
	if !vm.isWitnessVersionActive(0) {
		for _, sigInfo := range signatures {
			var match bool
			script, match = removeOpcodeByData(script, sigInfo.signature)
			if vm.hasFlag(ScriptVerifyConstScriptCode) && match {
				str := fmt.Sprintf("got match of %v in %v", sigInfo.signature,
					script)
				return scriptError(ErrNonConstScriptCode, str)
			}
		}
	}

