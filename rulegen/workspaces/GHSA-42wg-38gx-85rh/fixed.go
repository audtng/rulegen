package main

const maxConfigSize = 5 * 1024 * 1024    // 5 MB, should be largely enough
const maxDumpEntrySize = 500 * 1024 * 1024 // 500 MB

// parseDbFileName validates and extracts the table name from a database dump filename.
// Returns the table name and true if valid, or empty string and false if invalid.
func parseDbFileName(fname string) (string, bool) {
	if !strings.HasSuffix(fname, ".json") {
		return "", false
	}
	tableName := strings.TrimSuffix(fname, ".json")
	if tableName == "" || strings.ContainsAny(tableName, "/\\") {
		return "", false
	}
	return tableName, true
}

// Restore takes a zip file name and restores it
func Restore(filename string, overrideConfig bool) error {

		}
		if strings.HasPrefix(file.Name, "database/") {
			fname := strings.TrimPrefix(file.Name, "database/")
			tableName, valid := parseDbFileName(fname)
			if !valid {
				return fmt.Errorf("invalid database file name in zip archive: %q", file.Name)
			}
			dbfiles[tableName] = file
			continue
		}
		if file.Name == ".env" {

	lastMigration := ms[len(ms)-2]

	// Pre-validate all table data JSON before wiping to avoid leaving the database
	// in a destroyed state when the archive contains corrupted table data.
	for table, d := range dbfiles {
		if table == "migration" {
			continue
		}
		rc, err := d.Open()
		if err != nil {
			return fmt.Errorf("could not open table data for %s: %w", table, err)
		}
		var bufValidate bytes.Buffer
		if _, err := bufValidate.ReadFrom(io.LimitReader(rc, maxDumpEntrySize)); err != nil {
			rc.Close()
			return fmt.Errorf("could not read table data for %s: %w", table, err)
		}
		rc.Close()
		var test []map[string]interface{}
		if err := json.Unmarshal(bufValidate.Bytes(), &test); err != nil {
			return fmt.Errorf("invalid JSON in table data for %s: %w", table, err)
		}
	}

	// Start by wiping everything - only after we've validated the archive
	if err := db.WipeEverything(); err != nil {
		return fmt.Errorf("could not wipe database: %w", err)
