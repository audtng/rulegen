package main

func (w *Website) ToProtobuf(webContentDir string) *clientpb.Website {
	WebContents := map[string]*clientpb.WebContent{}
	for _, webcontent := range w.WebContents {
		contents, err := os.ReadFile(filepath.Join(webContentDir, webcontent.ID.String()))
		if err != nil {
			continue
		}
		WebContents[webcontent.ID.String()] = webcontent.ToProtobuf(&contents)
	}
	return &clientpb.Website{
