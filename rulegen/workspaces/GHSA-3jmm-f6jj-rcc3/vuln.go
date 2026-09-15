package main

}

func getPendingStagingFileCount(sourceOrDestId string, isSourceId bool) (fileCount int64, err error) {
	sourceOrDestColumn := ""
	if isSourceId {
		sourceOrDestColumn = "source_id"
		FROM
		  %[1]s
		WHERE
		  %[1]s.%[3]s = '%[2]s';
`,
		warehouseutils.WarehouseUploadsTable,
		sourceOrDestId,
		sourceOrDestColumn,
	)

	err = dbHandle.QueryRow(sqlStatement).Scan(&lastStagingFileIDRes)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("query: %s failed with Error : %w", sqlStatement, err)
		return
	}
	lastStagingFileID := int64(0)
		FROM
		  %[1]s
		WHERE
		  %[1]s.id > %[2]v
		  AND %[1]s.%[4]s = '%[3]s';
`,
		warehouseutils.WarehouseStagingFilesTable,
		lastStagingFileID,
		sourceOrDestId,
		sourceOrDestColumn,
	)

	err = dbHandle.QueryRow(sqlStatement).Scan(&fileCount)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("query: %s failed with Error : %w", sqlStatement, err)
		return
	}

