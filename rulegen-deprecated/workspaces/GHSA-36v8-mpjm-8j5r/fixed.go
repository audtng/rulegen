package main

	model.FlushTxQueue()
}

func isBacklinkDocAccessible(c *gin.Context, refTreeID string) bool {
	if !model.IsReadOnlyRoleContext(c) {
		return true
	}

	return model.CheckBlockIdAccessableByPublishAccess(c, model.GetPublishAccess(), refTreeID)
}

func getBackmentionDoc(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	defID := arg["defID"].(string)
	refTreeID := arg["refTreeID"].(string)
	if !isBacklinkDocAccessible(c, refTreeID) {
		ret.Data = map[string]any{
			"backmentions": []*model.Backlink{},
			"keywords":     []string{},
		}
		return
	}
	keyword := arg["keyword"].(string)
	var notebook string
	if val, ok := arg["notebook"]; ok {

	defID := arg["defID"].(string)
	refTreeID := arg["refTreeID"].(string)
	if !isBacklinkDocAccessible(c, refTreeID) {
		ret.Data = map[string]any{
			"backlinks": []*model.Backlink{},
			"keywords":  []string{},
		}
		return
	}
	keyword := arg["keyword"].(string)
	var notebook string
	if val, ok := arg["notebook"]; ok {
}

func CheckBlockIdAccessableByPublishAccess(c *gin.Context, publishAccess PublishAccess, blockID string) bool {
	bt := treenode.GetBlockTree(blockID)
	return checkBlockTreeAccessableByPublishAccess(c, publishAccess, bt)
}

func checkBlockTreeAccessableByPublishAccess(c *gin.Context, publishAccess PublishAccess, bt *treenode.BlockTree) bool {
	if bt == nil {
		return false
	}

	publishIgnore := GetDisablePublishAccess(publishAccess)
	passwordID, password := GetPathPasswordByPublishAccess(bt.BoxID, bt.Path, publishAccess)
	return CheckPathAccessableByPublishIgnore(bt.BoxID, bt.Path, publishIgnore) && (password == "" || CheckPublishAuthCookie(c, passwordID, password))
}
