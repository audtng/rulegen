package main

// limitations under the License.

/*
	cleanup-index checks what index entries are in the MySQL table and deletes those entries from the Redis database.
	It does not go the other way

	To run:
		}

		if err := s.Verify(dataFile, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying signature: %v", tc.caseDesc, err)
		}
	}


	var k PublicKey
	if len(k.Subjects()) != 0 {
		t.Errorf("Subjects for uninitialized key should give empty slice")
	}
	tests := []test{
		{caseDesc: "Valid armored public key", inputFile: "testdata/valid_armored_public.pgp", subjects: []string{}, keys: 2},
		}

		if err := s.Verify(dataFile, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying signature: %v", tc.caseDesc, err)
		}
	}

		}

		if err := s.Verify(nil, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying signature: %v", tc.caseDesc, err)
		}
	}

	var result []string

	// We add the key, the hash of the overall cose envelope, and the hash of the payload itself as keys.
	if v.CoseObj.PublicKey == nil {
		return nil, errors.New("missing public key")
	}
	keyObj, err := x509.NewPublicKey(bytes.NewReader(*v.CoseObj.PublicKey))
	if err != nil {
		return nil, err
		return err
	}

	if v.CoseObj.PublicKey == nil {
		return errors.New("missing public key")
	}
	v.keyObj, err = x509.NewPublicKey(bytes.NewReader(*v.CoseObj.PublicKey))
	if err != nil {
		return err

func (v *V001Entry) Canonicalize(_ context.Context) ([]byte, error) {
	if v.keyObj == nil {
		return nil, errors.New("cannot canonicalize empty key")
	}
	if v.sign1Msg == nil {
		return nil, errors.New("signed message uninitialized")
	}
	if v.sign1Msg.Payload == nil {
		return nil, errors.New("payload empty")
	}

	pk, err := v.keyObj.CanonicalValue()
	if err != nil {
		return nil, err
		})
	}
}

func TestV001Entry_IndexKeys_MissingPublicKey(t *testing.T) {
	v := V001Entry{
		CoseObj: models.CoseV001Schema{
			Data:      &models.CoseV001SchemaData{},
			PublicKey: nil,
		},
	}
	_, err := v.IndexKeys()
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "missing public key" {
		t.Errorf("expected 'missing public key' error, got %v", err)
	}
}

func TestCanonicalizeHandlesInvalidInput(t *testing.T) {
	v := &V001Entry{}

	// 1. Missing keyObj
	_, err := v.Canonicalize(context.TODO())
	if err == nil || err.Error() != "cannot canonicalize empty key" {
		t.Fatalf("expected error 'cannot canonicalize empty key', got %v", err)
	}

	// Setup valid keyObj for subsequent tests
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	der, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pub := pem.EncodeToMemory(&pem.Block{
		Bytes: der,
		Type:  "PUBLIC KEY",
	})
	keyObj, _ := sigx509.NewPublicKey(bytes.NewReader(pub))
	v.keyObj = keyObj

	// 2. Missing sign1Msg
	_, err = v.Canonicalize(context.TODO())
	if err == nil || err.Error() != "signed message uninitialized" {
		t.Fatalf("expected error 'signed message uninitialized', got %v", err)
	}

	// 3. Missing Payload in sign1Msg
	v.sign1Msg = gocose.NewSign1Message()
	v.sign1Msg.Payload = nil
	_, err = v.Canonicalize(context.TODO())
	if err == nil || err.Error() != "payload empty" {
		t.Fatalf("expected error 'payload empty', got %v", err)
	}
}
	}

	env := &dsse.Envelope{}
	if dsseObj.ProposedContent.Envelope == nil {
		return errors.New("proposed content envelope is missing")
	}
	if err := json.Unmarshal([]byte(*dsseObj.ProposedContent.Envelope), env); err != nil {
		return err
	}
	}

	for _, s := range canonicalEntry.Signatures {
		if s == nil || s.Signature == nil {
			return nil, errors.New("canonical entry missing required signature")
		}
	}
			},
			wantErr: true,
		},
		{
			name: "missing envelope with verifiers",
			it: &models.DSSEV001Schema{
				ProposedContent: &models.DSSEV001SchemaProposedContent{
					Verifiers: []strfmt.Base64{[]byte("verifier")},
				},
			},
			wantErr: true,
		},
		{
			env:  envelope(t, key, []byte(validPayload)),
			name: "valid",
	if err == nil {
		t.Fatalf("expected error canonicalizing invalid input")
	}

	v.DSSEObj.Signatures = []*models.DSSEV001SchemaSignaturesItems0{nil}
	_, err = v.Canonicalize(context.TODO())
	if err == nil {
		t.Fatalf("expected error canonicalizing nil signature")
	}
}
