package main


	// 将需要导出的文件/文件夹复制到临时文件夹
	for _, resourcePath := range resourcePaths {
		resourceFullPath := filepath.Join(util.WorkspaceDir, resourcePath)    // 资源完整路径
		resourceBaseName := filepath.Base(resourceFullPath)                   // 资源名称
		resourceCopyPath := filepath.Join(exportFolderPath, resourceBaseName) // 资源副本完整路径
		if err = filelock.Copy(resourceFullPath, resourceCopyPath); err != nil {
		return
	}

	preview := false
	if previewArg := arg["preview"]; nil != previewArg {
		preview = previewArg.(bool)

func BuiltInTemplateFuncs() (ret template.FuncMap) {
	ret = sprig.TxtFuncMap()
	ret["Weekday"] = util.Weekday
	ret["WeekdayCN"] = util.WeekdayCN
	ret["WeekdayCN2"] = util.WeekdayCN2
