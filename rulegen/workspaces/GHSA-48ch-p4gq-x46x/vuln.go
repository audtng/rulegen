package main

// GetResourcesByList fetches a list of resources from a slice of paths
func (vcls *VikunjaCaldavProjectStorage) GetResourcesByList(rpaths []string) (resources []data.Resource, err error) {

	// Parse the set of resourcepaths into usable uids
	// A path looks like this: /dav/projects/10/a6eb526d5748a5c499da202fe74f36ed1aea2aef.ics
	// So we split the url in parts, take the last one and strip the ".ics" at the end
	var uids []string
	for _, path := range rpaths {
		parts := strings.Split(path, "/")
		if len(parts) < 5 {
			continue
		}
		uids = append(uids, strings.TrimSuffix(parts[4], ".ics"))
	}

	if len(uids) == 0 {
	}

	for _, t := range tasks {
		rr := VikunjaProjectResourceAdapter{
			task: t,
		}
