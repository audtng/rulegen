package main

	if idsShouldMatch {
		for _, id := range params.FilterIDs {
			// for each OR filter we need to verify that lookup id column is not null to avoid failing during Find
			tx.Or("? = ? AND ? is not null", filterIDColumnName, id,
				lookupIDColumnName)
		}
	} else {
		for _, id := range params.FilterIDs {
