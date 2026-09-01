package main


	// Check if we need to query specific data
	if len(r.Data) > 0 {
		if err := r.ProbeWithDataVerification(db); err != nil {
			return false, err.Error()
		}
	} else {
		if err := r.ProbeWithPing(db); err != nil {
			return false, err.Error()
		}
	}

	return true, "Check MySQL Server Successfully!"

}

// ProbeWithPing do the health check with ping
func (r *MySQL) ProbeWithPing(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	row, err := db.Query("show status like \"uptime\"") // run a SQL to test
	if err != nil {
		return err
	}
	defer row.Close()
	return nil
}

// ProbeWithDataVerification do the health check with data verification
func (r *MySQL) ProbeWithDataVerification(db *sql.DB) error {
	for k, v := range r.Data {
		log.Debugf("[%s / %s / %s] - Verifying Data - [%s] : [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, k, v)
		sql, err := r.getSQL(k)
		if err != nil {
			return err
		}
		log.Debugf("[%s / %s / %s] - SQL - [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, sql)
		rows, err := db.Query(sql)
		if err != nil {
			return err
		}
		if !rows.Next() {
			rows.Close()
			return fmt.Errorf("No data found for [%s]", k)
		}
		//check the value is equal to the value in data
		var value string
		if err := rows.Scan(&value); err != nil {
			rows.Close()
			return err
		}
		if value != v {
			rows.Close()
			return fmt.Errorf("Value not match for [%s] expected [%s] got [%s] ", k, v, value)
		}
		rows.Close()
		log.Debugf("[%s / %s / %s] - Data Verified Successfully! - [%s] : [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, k, v)
	}
	return nil
}

// getSQL get the SQL statement
// input: database:table:column:key:value
// output: SELECT column FROM database.table WHERE key = value
	if len(fields) != 5 {
		return "", fmt.Errorf("Invalid SQL data - [%s]. (syntax: database:table:field:key:value)", str)
	}
	db := global.EscapeQuote(fields[0])
	table := global.EscapeQuote(fields[1])
	field := global.EscapeQuote(fields[2])
	key := global.EscapeQuote(fields[3])
	value := global.EscapeQuote(fields[4])
	//check value is int or not
	if _, err := strconv.Atoi(value); err != nil {
		return "", fmt.Errorf("Invalid SQL data - [%s], the value must be int", str)
	}

	sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` = %s", field, db, table, key, value)
	return sql, nil
}
	}
	return result
}

// EscapeQuote escape the string the single quote, double quote, and backtick
func EscapeQuote(str string) string {
	type Escape struct {
		From string
		To   string
	}
	escape := []Escape{
		{From: "`", To: ""}, // remove the backtick
		{From: `\`, To: `\\`},
		{From: `'`, To: `\'`},
		{From: `"`, To: `\"`},
	}

	for _, e := range escape {
		str = strings.ReplaceAll(str, e.From, e.To)
	}
	return str
}
	if len(fields) != 5 {
		return "", "", fmt.Errorf("Invalid SQL data - [%s]. (syntax: database:table:field:key:value)", str)
	}
	db := global.EscapeQuote(fields[0])
	table := global.EscapeQuote(fields[1])
	field := global.EscapeQuote(fields[2])
	key := global.EscapeQuote(fields[3])
	value := global.EscapeQuote(fields[4])
	//check value is int or not
	if _, err := strconv.Atoi(value); err != nil {
		return "", "", fmt.Errorf("Invalid SQL data - [%s], the value must be int", str)
	}

	sql := fmt.Sprintf(`SELECT "%s" FROM "%s" WHERE "%s" = %s`, field, table, key, value)
	return db, sql, nil
}
