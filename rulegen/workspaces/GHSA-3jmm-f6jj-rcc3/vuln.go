package main

	}

	for taskRunID, failedEvents := range taskRunIDFailedEventsMap {
		table := `"` + strings.ReplaceAll(fmt.Sprintf(`%s_%s`, failedKeysTablePrefix, taskRunID), `"`, `""`) + `"`
		sqlStatement := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		destination_id TEXT NOT NULL,
		record_id JSONB NOT NULL,
	}

	// Drop table
	table := fmt.Sprintf(`%s_%s`, failedKeysTablePrefix, taskRunID)
	sqlStatement := fmt.Sprintf(`DROP TABLE IF EXISTS %s`, table)
	_, err := fem.dbHandle.Exec(sqlStatement)
	if err != nil {

	var rows *sql.Rows
	var err error
	table := `"` + strings.ReplaceAll(fmt.Sprintf(`%s_%s`, failedKeysTablePrefix, taskRunID), `"`, `""`) + `"`
	sqlStatement := fmt.Sprintf(`SELECT %[1]s.destination_id, %[1]s.record_id
                                             FROM %[1]s `, table)
	rows, err = fem.dbHandle.Query(sqlStatement)
func (fem *FailedEventsManagerT) GetDBHandle() *sql.DB {
	return fem.dbHandle
}
