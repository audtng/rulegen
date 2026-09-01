package main

		last:     getLast(database, table),
	}
	err := store.database.Apply(func(db *sql.DB) error {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS '%s'('key' INTEGER PRIMARY KEY, 'val' BLOB);", table)
		stmt, err := db.Prepare(query)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(query)
		return err
	})
	if err != nil {

func getLast(d Database, table string) int64 {
	var last int64 = 0
	_ = d.Apply(func(db *sql.DB) error {
		query := fmt.Sprintf("SELECT key FROM %s Order by key DESC Limit 1;", table)
		stmt, err := db.Prepare(query)
		if err != nil {
			return err
		}
		row := stmt.QueryRow(query)
		return row.Scan(&last)
	})
	return last
