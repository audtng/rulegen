package main

func renderPackageREADME(linkBase string, mdData []byte) (ret string) {
	mdData = bytes.TrimPrefix(mdData, []byte("\xef\xbb\xbf")) // 移除文件开头的 BOM
	luteEngine := lute.New()
	luteEngine.SetSanitize(true)
	luteEngine.SetSoftBreak2HardBreak(false)
	luteEngine.SetCodeSyntaxHighlight(false)
	luteEngine.SetLinkBase(linkBase)
