package main

const maxConfigSize = 5 * 1024 * 1024    // 5 MB, should be largely enough
const maxDumpEntrySize = 500 * 1024 * 1024 // 500 MB

// Restore takes a zip file name and restores it
func Restore(filename string, overrideConfig bool) error {

		}
		if strings.HasPrefix(file.Name, "database/") {
			fname := strings.TrimPrefix(file.Name, "database/")
			if !strings.HasSuffix(fname, ".json") || len(fname) <= 5 {
				return fmt.Errorf("invalid database file name in zip archive: %q", file.Name)
			}
			dbfiles[fname[:len(fname)-5]] = file
			continue
		}
		if file.Name == ".env" {

	lastMigration := ms[len(ms)-2]

	// Start by wiping everything - only after we've validated the archive
	if err := db.WipeEverything(); err != nil {
		return fmt.Errorf("could not wipe database: %w", err)
