package main

	if err != nil {
		return nil, nil, fmt.Errorf("envelope payload can't be marshalled: %w", err)
	}

	var signingAgentId string
	if opts.SigningAgent != "" {
		signingAgentId = opts.SigningAgent
			ContentType: envelope.MediaTypePayloadV1,
			Content:     payloadBytes,
		},
		Signer:        s.signer,
		SigningTime:   time.Now(),
		SigningScheme: signature.SigningSchemeX509,
		SigningAgent:  signingAgentId,
		Timestamper:   opts.Timestamper,
		TSARootCAs:    opts.TSARootCAs,
	}

	// Add expiry only if ExpiryDuration is not zero
	logger.Debugf("  Expiry:        %v", signReq.Expiry)
	logger.Debugf("  SigningScheme: %v", signReq.SigningScheme)
	logger.Debugf("  SigningAgent:  %v", signReq.SigningAgent)

	// Add ctx to the SignRequest
	signReq = signReq.WithContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	sig, err := sigEnv.Sign(signReq)
	if err != nil {
		return nil, nil, err
	}

	envContent, err := sigEnv.Verify()
	if err != nil {
		return nil, nil, fmt.Errorf("generated signature failed verification: %v", err)
	orasRegistry "oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/notaryproject/notation-core-go/signature"
	"github.com/notaryproject/notation-core-go/signature/cose"
	"github.com/notaryproject/notation-core-go/signature/jws"

	// TSARootCAs is the cert pool holding caller's TSA trust anchor
	TSARootCAs *x509.CertPool
}

// Signer is a generic interface for signing an OCI artifact.
