package main

			return nil
		}

		if msg.AuthenticationFailureParameter == nil {
			return fmt.Errorf("missing AuthenticationFailureParameter IE for SynchFailure")
		}

		auts := msg.GetAuthenticationFailureParameter()
		resynchronizationInfo := &models.ResynchronizationInfo{
			Auts: hex.EncodeToString(auts[:]),
		return fmt.Errorf("ue Authentication Context is nil")
	}

	if msg.AuthenticationResponseParameter == nil {
		return fmt.Errorf("missing AuthenticationResponseParameter IE")
	}

	resStar := msg.GetRES()

	// Calculate HRES* (TS 33.501 Annex A.5)
