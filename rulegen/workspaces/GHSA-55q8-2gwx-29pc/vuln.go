package main

			return nil
		}

		auts := msg.GetAuthenticationFailureParameter()
		resynchronizationInfo := &models.ResynchronizationInfo{
			Auts: hex.EncodeToString(auts[:]),
		return fmt.Errorf("ue Authentication Context is nil")
	}

	resStar := msg.GetRES()

	// Calculate HRES* (TS 33.501 Annex A.5)
