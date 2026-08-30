package main

		return
	}

	if !util.IsAbsPathInWorkspace(p) {
		ret.Code = -1
		ret.Msg = "Path [" + p + "] is not in workspace"
		return
	}

	preview := false
	if previewArg := arg["preview"]; nil != previewArg {
		preview = previewArg.(bool)

	// 将需要导出的文件/文件夹复制到临时文件夹
	for _, resourcePath := range resourcePaths {
		resourceFullPath := filepath.Join(util.WorkspaceDir, resourcePath) // 资源完整路径
		if !util.IsAbsPathInWorkspace(resourceFullPath) {
			logging.LogErrorf("resource path [%s] is not in workspace", resourceFullPath)
			err = errors.New("resource path [" + resourcePath + "] is not in workspace")
			return
		}

		resourceBaseName := filepath.Base(resourceFullPath)                   // 资源名称
		resourceCopyPath := filepath.Join(exportFolderPath, resourceBaseName) // 资源副本完整路径
		if err = filelock.Copy(resourceFullPath, resourceCopyPath); err != nil {
	if nil != form.Value["assetsDirPath"] {
		relAssetsDirPath = form.Value["assetsDirPath"][0]
		assetsDirPath = filepath.Join(util.DataDir, relAssetsDirPath)
		if !util.IsAbsPathInWorkspace(assetsDirPath) {
			ret.Code = -1
			ret.Msg = "Path [" + assetsDirPath + "] is not in workspace"
			return
		}
	}
	if !gulu.File.IsExist(assetsDirPath) {
		if err = os.MkdirAll(assetsDirPath, 0755); err != nil {

func BuiltInTemplateFuncs() (ret template.FuncMap) {
	ret = sprig.TxtFuncMap()

	// 因为安全原因移除一些函数 https://github.com/siyuan-note/siyuan/issues/13426
	delete(ret, "env")
	delete(ret, "expandenv")
	delete(ret, "getHostByName")

	ret["Weekday"] = util.Weekday
	ret["WeekdayCN"] = util.WeekdayCN
	ret["WeekdayCN2"] = util.WeekdayCN2
	}
	return "", os.ErrPermission
}

func IsAbsPathInWorkspace(absPath string) bool {
	return IsSubPath(WorkspaceDir, absPath)
}
