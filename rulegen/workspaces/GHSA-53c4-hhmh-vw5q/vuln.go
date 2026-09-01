package main

// MustAdd adds a file to the index
// This can leave the index in an unsorted state
func (i IndexFile) MustAdd(md *chart.Metadata, filename, baseURL, digest string) error {
	if md.APIVersion == "" {
		md.APIVersion = chart.APIVersionV1
	}

	for name, cvs := range i.Entries {
		for idx := len(cvs) - 1; idx >= 0; idx-- {
			if cvs[idx].APIVersion == "" {
				cvs[idx].APIVersion = chart.APIVersionV1
			}
	cp := []*Entry{}
	found := false
	for _, rf := range r.Repositories {
		if rf.Name == name {
			found = true
			continue
