package main

	switch s := object.(type) {
	case *ssv1alpha1.SealedSecret:
		// Verify metainformation is well set up in Template ObjectMeta and ObjectMeta to avoid unconsistences with the scope during the rotate.
		// This is going to keep the original scope.
		if !reflect.DeepEqual(s.ObjectMeta, s.Spec.Template.ObjectMeta) {
			s.ObjectMeta.DeepCopyInto(&s.Spec.Template.ObjectMeta)
			slog.Warn("Sealed Secret metadata doesn't match. Please align your Sealed Secret metadata")
		}

		secret, err := c.attemptUnseal(s)
