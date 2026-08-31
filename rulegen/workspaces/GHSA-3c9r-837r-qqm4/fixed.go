package main

// FetchClient is a custom HTTP client.
type FetchClient struct {
	*http.Client
	userAgent    string
	allowedHosts map[string]struct{}
}

// NewClient creates a new FetchClient.
func NewClient(userAgent string, timeout int, reserveRedirect bool, allowedHosts map[string]struct{}) (client *FetchClient, recycle func()) {
	client = clientPool.Get().(*FetchClient)
	client.userAgent = userAgent
	client.Timeout = time.Duration(timeout) * time.Second
		if reserveRedirect && len(via) > 0 {
			return http.ErrUseLastResponse
		}
		// To avoid SSRF attacks, we check if the request URL's host is in the allowed hosts list.
		if allowedHosts != nil {
			if _, ok := allowedHosts[req.URL.Host]; !ok {
				return http.ErrUseLastResponse
			}
		}
		if len(via) >= 3 {
			return errors.New("stopped after 3 redirects")
		}

// Do sends an HTTP request and returns the response.
func (c *FetchClient) Fetch(url *url.URL, header http.Header) (resp *http.Response, err error) {
	if c.allowedHosts != nil {
		if _, ok := c.allowedHosts[url.Host]; !ok {
			return nil, errors.New("host not allowed: " + url.Host)
		}
	}
	if c.userAgent != "" {
		if header == nil {
			header = make(http.Header)
	if err != nil {
		return
	}
	client, recycle := fetch.NewClient("esmd/"+VERSION, 30, false, nil)
	defer recycle()
	res, err := client.Fetch(u, nil)
	if err != nil {
		return rex.Status(http.StatusBadRequest, "Invalid url")
	}

	client, recycle := fetch.NewClient(ctx.UserAgent(), 60, true, nil)
	defer recycle()

	res, err := client.Fetch(url, nil)
			header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(reg.User+":"+reg.Password)))
		}

		fetchClient, recycle := fetch.NewClient("esmd/"+VERSION, 15, false, nil)
		defer recycle()

		retryTimes := 0
		header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(reg.User+":"+reg.Password)))
	}

	fetchClient, recycle := fetch.NewClient("esmd/"+VERSION, 30, false, nil)
	defer recycle()

	retryTimes := 0
			indexHTML, err := withCache("index.html", time.Duration(cacheTtl)*time.Second, func() (indexHTML []byte, _ string, err error) {
				readme, err := os.ReadFile("README.md")
				if err != nil {
					fetchClient, recycle := fetch.NewClient(ctx.UserAgent(), 15, false, nil)
					defer recycle()
					readmeUrl, _ := url.Parse("https://raw.githubusercontent.com/esm-dev/esm.sh/refs/heads/main/README.md")
					var res *http.Response
			if v != "" && (!npm.Versioning.Match(v) || len(v) > 32) {
				return rex.Status(400, "Invalid Version Param")
			}
			allowedHosts := map[string]struct{}{}
			allowedHosts[modUrl.Host] = struct{}{}
			fetchClient, recycle := fetch.NewClient(ctx.UserAgent(), 15, false, allowedHosts)
			defer recycle()
			if strings.HasSuffix(modUrl.Path, "/uno.css") {
				ctxParam := query.Get("ctx")
						}
						defer res.Body.Close()
						if res.StatusCode != 200 {
							if res.StatusCode == 404 {
								return rex.Status(404, "Import map not found")
							}
							if res.StatusCode == 301 || res.StatusCode == 302 || res.StatusCode == 307 || res.StatusCode == 308 {
								return rex.Status(400, "Failed to fetch import map: redirects are not allowed")
							}
							return rex.Status(500, "Failed to fetch import map: "+res.Status)
						}
						tokenizer := html.NewTokenizer(io.LimitReader(res.Body, 5*MB))
						for {
			if err != nil {
				return rex.Err(http.StatusBadRequest, "Invalid url")
			}
			fetchClient, recycle := fetch.NewClient(ctx.UserAgent(), 15, false, nil)
			defer recycle()
			res, err := fetchClient.Fetch(url, nil)
			if err != nil {
						}
						defer res.Body.Close()
						if res.StatusCode != 200 {
							if res.StatusCode == 404 {
								return esbuild.OnLoadResult{}, errors.New("module not found: " + args.Path)
							}
							if res.StatusCode == 301 || res.StatusCode == 302 || res.StatusCode == 307 || res.StatusCode == 308 {
								return esbuild.OnLoadResult{}, errors.New("failed to fetch module " + args.Path + ": redirect not allowed")
							}
							return esbuild.OnLoadResult{}, errors.New("failed to fetch module " + args.Path + ": " + res.Status)
						}
						data, err := io.ReadAll(io.LimitReader(res.Body, 5*MB))
