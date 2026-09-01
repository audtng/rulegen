package main

// GetResourcesByList fetches a list of resources from a slice of paths
func (vcls *VikunjaCaldavProjectStorage) GetResourcesByList(rpaths []string) (resources []data.Resource, err error) {

	// Path format: /dav/projects/{projectID}/{uid}.ics.
	// Remember the href's project ID per uid so the consistency check below
	// can drop tasks requested via the wrong project (GHSA-48ch-p4gq-x46x).
	var uids []string
	uidURLProjects := map[string]int64{}
	for _, path := range rpaths {
		parts := strings.Split(path, "/")
		if len(parts) < 5 {
			continue
		}
		// Skip malformed hrefs: without these guards an empty/non-numeric
		// project would bypass the consistency check, and an empty UID
		// would query for tasks with empty/NULL uids.
		if !strings.HasSuffix(parts[4], ".ics") {
			continue
		}
		uid := strings.TrimSuffix(parts[4], ".ics")
		if uid == "" {
			continue
		}
		urlProjectID, perr := strconv.ParseInt(parts[3], 10, 64)
		if perr != nil {
			continue
		}
		uids = append(uids, uid)
		uidURLProjects[uid] = urlProjectID
	}

	if len(uids) == 0 {
	}

	for _, t := range tasks {
		// Closes the URL-path leak for calendar-multiget REPORTs
		// (GHSA-48ch-p4gq-x46x).
		if urlProjectID, ok := uidURLProjects[t.UID]; ok && urlProjectID != t.ProjectID {
			continue
		}
		rr := VikunjaProjectResourceAdapter{
			task: t,
		}
