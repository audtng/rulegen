package main

	model.FlushTxQueue()
}

func getBackmentionDoc(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	defID := arg["defID"].(string)
	refTreeID := arg["refTreeID"].(string)
	keyword := arg["keyword"].(string)
	var notebook string
	if val, ok := arg["notebook"]; ok {

	defID := arg["defID"].(string)
	refTreeID := arg["refTreeID"].(string)
	keyword := arg["keyword"].(string)
	var notebook string
	if val, ok := arg["notebook"]; ok {
}

func CheckBlockIdAccessableByPublishAccess(c *gin.Context, publishAccess PublishAccess, blockID string) bool {
	publishIgnore := GetDisablePublishAccess(publishAccess)
	bt := treenode.GetBlockTree(blockID)
	if bt == nil {
		return false
	}
	passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
	return CheckPathAccessableByPublishIgnore(bt.BoxID, bt.Path, publishIgnore) && (password == "" || CheckPublishAuthCookie(c, passwordID, password))
}
