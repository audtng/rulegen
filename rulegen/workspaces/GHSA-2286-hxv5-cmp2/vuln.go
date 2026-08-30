package main

func (w *Website) ToProtobuf(webContentDir string) *clientpb.Website {
	WebContents := map[string]*clientpb.WebContent{}
	for _, webcontent := range w.WebContents {
		contents, _ := os.ReadFile(filepath.Join(webContentDir, webcontent.Path))
		WebContents[webcontent.ID.String()] = webcontent.ToProtobuf(&contents)
	}
	return &clientpb.Website{
