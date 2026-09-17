package main

	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	safetemp "github.com/hashicorp/go-safetemp"
// wish. The response must be a 2xx.
//
// First, a header is looked for "X-Terraform-Get" which should contain
// a source URL to download. This source must use one of the configured
// protocols and getters for the client, or "http"/"https" if using
// the HttpGetter directly.
//
// If the header is not present, then a meta tag is searched for named
// "terraform-get" and the content should be a source URL.
	// and as such it needs to be initialized before use, via something like
	// make(http.Header).
	Header http.Header

	// DoNotCheckHeadFirst configures the client to NOT check if the server
	// supports HEAD requests.
	DoNotCheckHeadFirst bool

	// HeadFirstTimeout configures the client to enforce a timeout when
	// the server supports HEAD requests.
	//
	// The zero value means no timeout.
	HeadFirstTimeout time.Duration

	// ReadTimeout configures the client to enforce a timeout when
	// making a request to an HTTP server and reading its response body.
	//
	// The zero value means no timeout.
	ReadTimeout time.Duration

	// MaxBytes limits the number of bytes that will be ready from an HTTP
	// response body returned from a server. The zero value means no limit.
	MaxBytes int64

	// XTerraformGetLimit configures how many times the client with follow
	// the " X-Terraform-Get" header value.
	//
	// The zero value means no limit.
	XTerraformGetLimit int

	// XTerraformGetDisabled disables the client's usage of the "X-Terraform-Get"
	// header value.
	XTerraformGetDisabled bool
}

func (g *HttpGetter) ClientMode(u *url.URL) (ClientMode, error) {
	return ClientModeFile, nil
}

type contextKey int

const (
	xTerraformGetDisable           contextKey = 0
	xTerraformGetLimit             contextKey = 1
	xTerraformGetLimitCurrentValue contextKey = 2
	httpClientValue                contextKey = 3
	httpMaxBytesValue              contextKey = 4
)

func xTerraformGetDisabled(ctx context.Context) bool {
	value, ok := ctx.Value(xTerraformGetDisable).(bool)
	if !ok {
		return false
	}
	return value
}

func xTerraformGetLimitCurrentValueFromContext(ctx context.Context) int {
	value, ok := ctx.Value(xTerraformGetLimitCurrentValue).(int)
	if !ok {
		return 1
	}
	return value
}

func xTerraformGetLimiConfiguredtFromContext(ctx context.Context) int {
	value, ok := ctx.Value(xTerraformGetLimit).(int)
	if !ok {
		return 0
	}
	return value
}

func httpClientFromContext(ctx context.Context) *http.Client {
	value, ok := ctx.Value(httpClientValue).(*http.Client)
	if !ok {
		return nil
	}
	return value
}

func httpMaxBytesFromContext(ctx context.Context) int64 {
	value, ok := ctx.Value(httpMaxBytesValue).(int64)
	if !ok {
		return 0 // no limit
	}
	return value
}

type limitedWrappedReaderCloser struct {
	underlying io.Reader
	closeFn    func() error
}

func (l *limitedWrappedReaderCloser) Read(p []byte) (n int, err error) {
	return l.underlying.Read(p)
}

func (l *limitedWrappedReaderCloser) Close() (err error) {
	return l.closeFn()
}

func newLimitedWrappedReaderCloser(r io.ReadCloser, limit int64) io.ReadCloser {
	return &limitedWrappedReaderCloser{
		underlying: io.LimitReader(r, limit),
		closeFn:    r.Close,
	}
}

func (g *HttpGetter) Get(dst string, u *url.URL) error {
	ctx := g.Context()

	// Optionally disable any X-Terraform-Get redirects. This is reccomended for usage of
	// this client outside of Terraform's. This feature is likely not required if the
	// source server can provider normal HTTP redirects.
	if g.XTerraformGetDisabled {
		ctx = context.WithValue(ctx, xTerraformGetDisable, g.XTerraformGetDisabled)
	}

	// Optionally enforce a limit on X-Terraform-Get redirects. We check this for every
	// invocation of this function, because the value is not passed down to subsequent
	// client Get function invocations.
	if g.XTerraformGetLimit > 0 {
		ctx = context.WithValue(ctx, xTerraformGetLimit, g.XTerraformGetLimit)
	}

	// If there was a limit on X-Terraform-Get redirects, check what the current count value.
	//
	// If the value is greater than the limit, return an error. Otherwise, increment the value,
	// and include it in the the context to be passed along in all the subsequent client
	// Get function invocations.
	if limit := xTerraformGetLimiConfiguredtFromContext(ctx); limit > 0 {
		currentValue := xTerraformGetLimitCurrentValueFromContext(ctx)

		if currentValue > limit {
			return fmt.Errorf("too many X-Terraform-Get redirects: %d", currentValue)
		}

		currentValue++

		ctx = context.WithValue(ctx, xTerraformGetLimitCurrentValue, currentValue)
	}

	// Optionally enforce a maxiumum HTTP response body size.
	if g.MaxBytes > 0 {
		ctx = context.WithValue(ctx, httpMaxBytesValue, g.MaxBytes)
	}

	// Copy the URL so we can modify it
	var newU url.URL = *u
	u = &newU
		}
	}

	// If the HTTP client is nil, check if there is one available in the context,
	// otherwise create one using cleanhttp's default transport.
	if g.Client == nil {
		if client := httpClientFromContext(ctx); client != nil {
			g.Client = client
		} else {
			client := httpClient
			if g.client != nil && g.client.Insecure {
				insecureTransport := cleanhttp.DefaultTransport()
				insecureTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				client.Transport = insecureTransport
			}
			g.Client = client
		}
	}

	// Pass along the configured HTTP client in the context for usage with the X-Terraform-Get feature.
	ctx = context.WithValue(ctx, httpClientValue, g.Client)

	// Add terraform-get to the parameter.
	q := u.Query()
	q.Add("terraform-get", "1")
	u.RawQuery = q.Encode()

	readCtx := ctx

	if g.ReadTimeout > 0 {
		var cancel context.CancelFunc
		readCtx, cancel = context.WithTimeout(ctx, g.ReadTimeout)
		defer cancel()
	}

	// Get the URL
	req, err := http.NewRequestWithContext(readCtx, "GET", u.String(), nil)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body := resp.Body

	if maxBytes := httpMaxBytesFromContext(ctx); maxBytes > 0 {
		body = newLimitedWrappedReaderCloser(body, maxBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	if disabled := xTerraformGetDisabled(ctx); disabled {
		return nil
	}

	// Extract the source URL,
	var source string
	if v := resp.Header.Get("X-Terraform-Get"); v != "" {
		source = v
	} else {
		source, err = g.parseMeta(readCtx, body)
		if err != nil {
			return err
		}
	source, subDir := SourceDirSubdir(source)
	if subDir == "" {
		var opts []ClientOption

		// Check if the protocol was switched to one which was not configured.
		//
		// Otherwise, all default getters are allowed.
		if g.client != nil && g.client.Getters != nil {
			protocol := strings.Split(source, ":")[0]
			_, allowed := g.client.Getters[protocol]
			if !allowed {
				return fmt.Errorf("no getter available for X-Terraform-Get source protocol: %q", protocol)
			}
		}

		// Add any getter client options.
		if g.client != nil {
			opts = g.client.Options
		}

		// If the client is nil, we know we're using the HttpGetter directly. In this case,
		// we don't know exactly which protocols are configued, but we can make a good guess.
		//
		// This prevents all default getters from being allowed when only using the
		// HttpGetter directly. To enable protocol switching, a client "wrapper" must
		// be used.
		if g.client == nil {
			opts = append(opts, WithGetters(map[string]Getter{
				"http":  g,
				"https": g,
			}))
		}

		// Ensure we pass along the context we constructed in this function.
		//
		// This is especially important to enforce a limit on X-Terraform-Get redirects
		// which could be setup, if configured, at the top of this function.
		opts = append(opts, WithContext(ctx))

		// Note: this allows the protocol to be switched to another configured getters.
		return Get(dst, source, opts...)
	}

// appended.
func (g *HttpGetter) GetFile(dst string, src *url.URL) error {
	ctx := g.Context()

	// Optionally enforce a maxiumum HTTP response body size.
	if g.MaxBytes > 0 {
		ctx = context.WithValue(ctx, httpMaxBytesValue, g.MaxBytes)
	}

	if g.Netrc {
		// Add auth from netrc if we can
		if err := addAuthFromNetrc(src); err != nil {
		}
	}

	var (
		currentFileSize int64
		req             *http.Request
	)

	if !g.DoNotCheckHeadFirst {
		headCtx := ctx

		if g.HeadFirstTimeout > 0 {
			var cancel context.CancelFunc

			headCtx, cancel = context.WithTimeout(ctx, g.HeadFirstTimeout)
			defer cancel()
		}

		// We first make a HEAD request so we can check
		// if the server supports range queries. If the server/URL doesn't
		// support HEAD requests, we just fall back to GET.
		req, err = http.NewRequestWithContext(headCtx, "HEAD", src.String(), nil)
		if err != nil {
			return err
		}
		if g.Header != nil {
			req.Header = g.Header.Clone()
		}
		headResp, err := g.Client.Do(req)
		if err == nil {
			headResp.Body.Close()
			if headResp.StatusCode == 200 {
				// If the HEAD request succeeded, then attempt to set the range
				// query if we can.
				if headResp.Header.Get("Accept-Ranges") == "bytes" && headResp.ContentLength >= 0 {
					if fi, err := f.Stat(); err == nil {
						if _, err = f.Seek(0, io.SeekEnd); err == nil {
							currentFileSize = fi.Size()
							if currentFileSize >= headResp.ContentLength {
								// file already present
								return nil
							}
						}
					}
				}
		}
	}

	readCtx := ctx

	if g.ReadTimeout > 0 {
		var cancel context.CancelFunc
		readCtx, cancel = context.WithTimeout(ctx, g.ReadTimeout)
		defer cancel()
	}

	req, err = http.NewRequestWithContext(readCtx, "GET", src.String(), nil)
	if err != nil {
		return err
	}

	body := resp.Body

	if maxBytes := httpMaxBytesFromContext(ctx); maxBytes > 0 {
		body = newLimitedWrappedReaderCloser(body, maxBytes)
	}

	if g.client != nil && g.client.ProgressListener != nil {
		// track download
		fn := filepath.Base(src.EscapedPath())
	}
	defer body.Close()

	n, err := Copy(readCtx, f, body)
	if err == nil && n < resp.ContentLength {
		err = io.ErrShortWrite
	}
		return err
	}

	var disableSymlinks bool

	if g.client != nil && g.client.DisableSymlinks {
		disableSymlinks = true
	}

	return copyDir(ctx, dst, sourcePath, false, disableSymlinks, g.client.umask())
}

// parseMeta looks for the first meta tag in the given reader that
// will give us the source URL.
func (g *HttpGetter) parseMeta(ctx context.Context, r io.Reader) (string, error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var err error
	var t xml.Token
	for {
		if ctx.Err() != nil {
			return "", fmt.Errorf("context error while parsing meta tag: %w", ctx.Err())
		}

		t, err = d.Token()
		if err != nil {
			if err == io.EOF {
package getter

import (
	"context"
	"os"
)

// ClientOption is used to configure a client.
type ClientOption func(*Client) error

// Configure applies all of the given client options, along with any default
// behavior including context, decompressors, detectors, and getters used by
// the client.
func (c *Client) Configure(opts ...ClientOption) error {
	// If the context has not been configured use the background context.
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}

	// Store the options used to configure this client.
	c.Options = opts

	// Apply all of the client options.
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return err
		}
	}

	// If the client was not configured with any Decompressors, Detectors,
	// or Getters, use the default values for each.
	if c.Decompressors == nil {
		c.Decompressors = Decompressors
	}
	if c.Detectors == nil {
		c.Detectors = Detectors
	}
	if c.Getters == nil {
		c.Getters = Getters
	}

	// Set the client for each getter, so the top-level client can know
	// the getter-specific client functions or progress tracking.
	for _, getter := range c.Getters {
		getter.SetClient(c)
	}

	return nil
}

// WithContext allows to pass a context to operation
// in order to be able to cancel a download in progress.
func WithContext(ctx context.Context) ClientOption {
	return func(c *Client) error {
		c.Ctx = ctx
		return nil
	}
}

// WithDecompressors specifies which Decompressor are available.
func WithDecompressors(decompressors map[string]Decompressor) ClientOption {
	return func(c *Client) error {
		c.Decompressors = decompressors
		return nil
	}
}

// WithDecompressors specifies which compressors are available.
func WithDetectors(detectors []Detector) ClientOption {
	return func(c *Client) error {
		c.Detectors = detectors
		return nil
	}
}

// WithGetters specifies which getters are available.
func WithGetters(getters map[string]Getter) ClientOption {
	return func(c *Client) error {
		c.Getters = getters
		return nil
	}
}

// WithMode specifies which client mode the getters should operate in.
func WithMode(mode ClientMode) ClientOption {
	return func(c *Client) error {
		c.Mode = mode
		return nil
	}
}

// WithUmask specifies how to mask file permissions when storing local
// files or decompressing an archive.
func WithUmask(mode os.FileMode) ClientOption {
	return func(c *Client) error {
		c.Umask = mode
		return nil
	}
}
	"runtime"
	"strconv"
	"strings"
	"time"

	urlhelper "github.com/hashicorp/go-getter/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"
// a git repository.
type GitGetter struct {
	getter

	// Timeout sets a deadline which all git CLI operations should
	// complete within. Zero value means no timeout.
	Timeout time.Duration
}

var defaultBranchRegexp = regexp.MustCompile(`\s->\sorigin/(.*)`)

func (g *GitGetter) Get(dst string, u *url.URL) error {
	ctx := g.Context()

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git must be available and on the PATH")
	}
	var sshKeyFile string
	if sshKey != "" {
		// Check that the git version is sufficiently new.
		if err := checkGitVersion(ctx, "2.3"); err != nil {
			return fmt.Errorf("Error using ssh key: %v", err)
		}


	// Next: check out the proper tag/branch if it is specified, and checkout
	if ref != "" {
		if err := g.checkout(ctx, dst, ref); err != nil {
			return err
		}
	}
	return fg.GetFile(dst, u)
}

func (g *GitGetter) checkout(ctx context.Context, dst string, ref string) error {
	cmd := exec.CommandContext(ctx, "git", "checkout", ref)
	cmd.Dir = dst
	return getRunCommand(cmd)
}

	originalRef := ref // we handle an unspecified ref differently than explicitly selecting the default branch below
	if ref == "" {
		ref = findRemoteDefaultBranch(ctx, u)
	}
	if depth > 0 {
		args = append(args, "--depth", strconv.Itoa(depth))
		// If we didn't add --depth and --branch above then we will now be
		// on the remote repository's default branch, rather than the selected
		// ref, so we'll need to fix that before we return.
		return g.checkout(ctx, dst, originalRef)
	}
	return nil
}
		// Not a branch, switch to default branch. This will also catch
		// non-existent branches, in which case we want to switch to default
		// and then checkout the proper branch later.
		ref = findDefaultBranch(ctx, dst)
	}

	// We have to be on a branch to pull
	if err := g.checkout(ctx, dst, ref); err != nil {
		return err
	}

	if depth > 0 {
		cmd = exec.CommandContext(ctx, "git", "pull", "--depth", strconv.Itoa(depth), "--ff-only")
	} else {
		cmd = exec.CommandContext(ctx, "git", "pull", "--ff-only")
	}

	cmd.Dir = dst
// findDefaultBranch checks the repo's origin remote for its default branch
// (generally "master"). "master" is returned if an origin default branch
// can't be determined.
func findDefaultBranch(ctx context.Context, dst string) string {
	var stdoutbuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "branch", "-r", "--points-at", "refs/remotes/origin/HEAD")
	cmd.Dir = dst
	cmd.Stdout = &stdoutbuf
	err := cmd.Run()

// findRemoteDefaultBranch checks the remote repo's HEAD symref to return the remote repo's
// default branch. "master" is returned if no HEAD symref exists.
func findRemoteDefaultBranch(ctx context.Context, u *url.URL) string {
	var stdoutbuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--symref", u.String(), "HEAD")
	cmd.Stdout = &stdoutbuf
	err := cmd.Run()
	matches := lsRemoteSymRefRegexp.FindStringSubmatch(stdoutbuf.String())
// checkGitVersion is used to check the version of git installed on the system
// against a known minimum version. Returns an error if the installed version
// is older than the given minimum.
func checkGitVersion(ctx context.Context, min string) error {
	want, err := version.NewVersion(min)
	if err != nil {
		return err
	}

	out, err := exec.CommandContext(ctx, "git", "version").Output()
	if err != nil {
		return err
	}
