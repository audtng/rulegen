package main

		if db.ParadeDBAvailable() {
			// ParadeDB fuzzy(1, prefix=true) on "TEST10" also matches
			// "test1", "test11", "test19", "test30" (edit distance 1), etc.
			// The recursive CTE also pulls in project 43 as a child of the
			// matched project 10 (reparent-escalation fixture).
			require.Len(t, ls, 7)
			projectIDs := make([]int64, len(ls))
			for i, p := range ls {
				projectIDs[i] = p.ID
			}
			assert.Contains(t, projectIDs, int64(10))
			assert.Contains(t, projectIDs, int64(43))
			assert.Contains(t, projectIDs, int64(-1))
		} else {
			// Expect project 10 (the search target), project 43 (its child —
