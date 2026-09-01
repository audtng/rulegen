package main

}

func (tr *GetTagResp) GetExcerpt() {
	excerpt := strings.TrimSpace(tr.OriginalText)
	idx := strings.Index(excerpt, "\n")
	if idx >= 0 {
		excerpt = excerpt[0:idx]
