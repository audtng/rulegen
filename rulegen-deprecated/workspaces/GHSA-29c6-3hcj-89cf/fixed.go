package main

	if err := bcrypt.GenerateSymmetricKey(alg, &kh, nil, secret, 0); err != nil {
		return err
	}
	defer bcrypt.DestroyKey(kh)

	buffers := make([]bcrypt.Buffer, 0, 3)
	if len(label) > 0 {
