package main

// FetchClient is a custom HTTP client.
type FetchClient struct {
	*http.Client
	userAgent string
}

// NewClient creates a new FetchClient.
func NewClient(userAgent string, timeout int, reserveRedirect bool) (client *FetchClient, recycle func()) {
	client = clientPool.Get().(*FetchClient)
	client.userAgent = userAgent
	client.Timeout = time.Duration(timeout) * time.Second
		if reserveRedirect && len(via) > 0 {
			return http.ErrUseLastResponse
		}
		if len(via) >= 3 {
			return errors.New("stopped after 3 redirects")
		}

// Do sends an HTTP request and returns the response.
func (c *FetchClient) Fetch(url *url.URL, header http.Header) (resp *http.Response, err error) {
	if c.userAgent != "" {
		if header == nil {
			header = make(http.Header)
			indexHTML, err := withCache("index.html", time.Duration(cacheTtl)*time.Second, func() (indexHTML []byte, _ string, err error) {
				readme, err := os.ReadFile("README.md")
				if err != nil {
					fetchClient, recycle := fetch.NewClient(ctx.UserAgent(), 15, false)
					defer recycle()
					readmeUrl, _ := url.Parse("https://raw.githubusercontent.com/esm-dev/esm.sh/refs/heads/main/README.md")
					var res *http.Response
			if v != "" && (!npm.Versioning.Match(v) || len(v) > 32) {
				return rex.Status(400, "Invalid Version Param")
			}
			fetchClient, recycle := fetch.NewClient(ctx.UserAgent(), 15, false)
			defer recycle()
			if strings.HasSuffix(modUrl.Path, "/uno.css") {
				ctxParam := query.Get("ctx")
						}
						defer res.Body.Close()
						if res.StatusCode != 200 {
							return rex.Status(500, "Failed to fetch import map")
						}
						tokenizer := html.NewTokenizer(io.LimitReader(res.Body, 5*MB))
						for {
						}
						defer res.Body.Close()
						if res.StatusCode != 200 {
							return esbuild.OnLoadResult{}, errors.New("failed to fetch module " + args.Path + ": " + res.Status)
						}
						data, err := io.ReadAll(io.LimitReader(res.Body, 5*MB))
