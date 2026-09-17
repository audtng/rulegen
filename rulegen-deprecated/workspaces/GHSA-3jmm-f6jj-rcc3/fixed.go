package main

}

func getPendingStagingFileCount(sourceOrDestId string, isSourceId bool) (fileCount int64, err error) {
	sourceOrDestId = pq.QuoteIdentifier(sourceOrDestId)
	sourceOrDestColumn := ""
	if isSourceId {
		sourceOrDestColumn = "source_id"
		FROM
		  %[1]s
		WHERE
		  %[2]s = $1;
`,
		warehouseutils.WarehouseUploadsTable,
		sourceOrDestColumn,
	)
	err = dbHandle.QueryRow(sqlStatement, sourceOrDestId).Scan(&lastStagingFileIDRes)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("query: %s run failed with Error : %w", sqlStatement, err)
		return
	}
	lastStagingFileID := int64(0)
		FROM
		  %[1]s
		WHERE
		  id > %[2]v
		  AND %[3]s = $1;
`,
		warehouseutils.WarehouseStagingFilesTable,
		lastStagingFileID,
		sourceOrDestColumn,
	)
	err = dbHandle.QueryRow(sqlStatement, sourceOrDestId).Scan(&fileCount)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("query: %s run failed with Error : %w", sqlStatement, err)
		return
	}

