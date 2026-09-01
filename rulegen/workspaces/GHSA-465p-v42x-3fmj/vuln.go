package main

	switch s := object.(type) {
	case *ssv1alpha1.SealedSecret:
		// Verify metainformation is well set up in Template ObjectMeta and ObjectMeta to avoid unconsistences with the scope during the rotate.
		if !reflect.DeepEqual(s.ObjectMeta, s.Spec.Template.ObjectMeta) {
			return nil, fmt.Errorf("Sealed Secret invalid: metadata no longer matches the sealed secret")
		}

		secret, err := c.attemptUnseal(s)
