package main

// limitations under the License.

/*
	cleanup-index checks what index entries are in the MySQL table and deletes those entries from the Redis databse.
	It does not go the other way

	To run:
		}

		if err := s.Verify(dataFile, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying sigature: %v", tc.caseDesc, err)
		}
	}


	var k PublicKey
	if len(k.Subjects()) != 0 {
		t.Errorf("Subjects for unitialized key should give empty slice")
	}
	tests := []test{
		{caseDesc: "Valid armored public key", inputFile: "testdata/valid_armored_public.pgp", subjects: []string{}, keys: 2},
		}

		if err := s.Verify(dataFile, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying sigature: %v", tc.caseDesc, err)
		}
	}

		}

		if err := s.Verify(nil, k); (err == nil) != tc.verified {
			t.Errorf("%v: unexpected result in verifying sigature: %v", tc.caseDesc, err)
		}
	}

	var result []string

	// We add the key, the hash of the overall cose envelope, and the hash of the payload itself as keys.
	keyObj, err := x509.NewPublicKey(bytes.NewReader(*v.CoseObj.PublicKey))
	if err != nil {
		return nil, err
		return err
	}

	v.keyObj, err = x509.NewPublicKey(bytes.NewReader(*v.CoseObj.PublicKey))
	if err != nil {
		return err

func (v *V001Entry) Canonicalize(_ context.Context) ([]byte, error) {
	if v.keyObj == nil {
		return nil, errors.New("cannot canonicalze empty key")
	}
	pk, err := v.keyObj.CanonicalValue()
	if err != nil {
		return nil, err
		})
	}
}
	}

	env := &dsse.Envelope{}
	if err := json.Unmarshal([]byte(*dsseObj.ProposedContent.Envelope), env); err != nil {
		return err
	}
	}

	for _, s := range canonicalEntry.Signatures {
		if s.Signature == nil {
			return nil, errors.New("canonical entry missing required signature")
		}
	}
			},
			wantErr: true,
		},
		{
			env:  envelope(t, key, []byte(validPayload)),
			name: "valid",
	if err == nil {
		t.Fatalf("expected error canonicalizing invalid input")
	}
}
