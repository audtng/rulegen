package main


	// Check if we need to query specific data
	if len(r.Data) > 0 {
		for k, v := range r.Data {
			log.Debugf("[%s / %s / %s] - Verifying Data - [%s] : [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, k, v)
			sql, err := r.getSQL(k)
			if err != nil {
				return false, err.Error()
			}
			log.Debugf("[%s / %s / %s] - SQL - [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, sql)
			rows, err := db.Query(sql)
			if err != nil {
				return false, err.Error()
			}
			if !rows.Next() {
				rows.Close()
				return false, fmt.Sprintf("No data found for [%s]", k)
			}
			//check the value is equal to the value in data
			var value string
			if err := rows.Scan(&value); err != nil {
				rows.Close()
				return false, err.Error()
			}
			if value != v {
				rows.Close()
				return false, fmt.Sprintf("Value not match for [%s] expected [%s] got [%s] ", k, v, value)
			}
			rows.Close()
			log.Debugf("[%s / %s / %s] - Data Verified Successfully! - [%s] : [%s]", r.ProbeKind, r.ProbeName, r.ProbeTag, k, v)
		}
	} else {
		err = db.Ping()
		if err != nil {
			return false, err.Error()
		}
		row, err := db.Query("show status like \"uptime\"") // run a SQL to test
		if err != nil {
			return false, err.Error()
		}
		defer row.Close()
	}

	return true, "Check MySQL Server Successfully!"

}

// getSQL get the SQL statement
// input: database:table:column:key:value
// output: SELECT column FROM database.table WHERE key = value
	if len(fields) != 5 {
		return "", fmt.Errorf("Invalid SQL data - [%s]. (syntax: database:table:field:key:value)", str)
	}
	db := fields[0]
	table := fields[1]
	field := fields[2]
	key := fields[3]
	value := fields[4]
	//check value is int or not
	if _, err := strconv.Atoi(value); err != nil {
		return "", fmt.Errorf("Invalid SQL data - [%s], the value must be int", str)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s = %s", field, db, table, key, value)
	return sql, nil
}
	}
	return result
}
	if len(fields) != 5 {
		return "", "", fmt.Errorf("Invalid SQL data - [%s]. (syntax: database:table:field:key:value)", str)
	}
	db := fields[0]
	table := fields[1]
	field := fields[2]
	key := fields[3]
	value := fields[4]
	//check value is int or not
	if _, err := strconv.Atoi(value); err != nil {
		return "", "", fmt.Errorf("Invalid SQL data - [%s], the value must be int", str)
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s", field, table, key, value)
	return db, sql, nil
}
