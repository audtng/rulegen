package main

	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	safetemp "github.com/hashicorp/go-safetemp"
// wish. The response must be a 2xx.
//
// First, a header is looked for "X-Terraform-Get" which should contain
// a source URL to download.
//
// If the header is not present, then a meta tag is searched for named
// "terraform-get" and the content should be a source URL.
	// and as such it needs to be initialized before use, via something like
	// make(http.Header).
	Header http.Header
}

func (g *HttpGetter) ClientMode(u *url.URL) (ClientMode, error) {
	return ClientModeFile, nil
}

func (g *HttpGetter) Get(dst string, u *url.URL) error {
	ctx := g.Context()
	// Copy the URL so we can modify it
	var newU url.URL = *u
	u = &newU
		}
	}

	if g.Client == nil {
		g.Client = httpClient
		if g.client != nil && g.client.Insecure {
			insecureTransport := cleanhttp.DefaultTransport()
			insecureTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			g.Client.Transport = insecureTransport
		}
	}

	// Add terraform-get to the parameter.
	q := u.Query()
	q.Add("terraform-get", "1")
	u.RawQuery = q.Encode()

	// Get the URL
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	// Extract the source URL
	var source string
	if v := resp.Header.Get("X-Terraform-Get"); v != "" {
		source = v
	} else {
		source, err = g.parseMeta(resp.Body)
		if err != nil {
			return err
		}
	source, subDir := SourceDirSubdir(source)
	if subDir == "" {
		var opts []ClientOption
		if g.client != nil {
			opts = g.client.Options
		}
		return Get(dst, source, opts...)
	}

// appended.
func (g *HttpGetter) GetFile(dst string, src *url.URL) error {
	ctx := g.Context()
	if g.Netrc {
		// Add auth from netrc if we can
		if err := addAuthFromNetrc(src); err != nil {
		}
	}

	var currentFileSize int64

	// We first make a HEAD request so we can check
	// if the server supports range queries. If the server/URL doesn't
	// support HEAD requests, we just fall back to GET.
	req, err := http.NewRequestWithContext(ctx, "HEAD", src.String(), nil)
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

	req, err = http.NewRequestWithContext(ctx, "GET", src.String(), nil)
	if err != nil {
		return err
	}

	body := resp.Body

	if g.client != nil && g.client.ProgressListener != nil {
		// track download
		fn := filepath.Base(src.EscapedPath())
	}
	defer body.Close()

	n, err := Copy(ctx, f, body)
	if err == nil && n < resp.ContentLength {
		err = io.ErrShortWrite
	}
		return err
	}

	return copyDir(ctx, dst, sourcePath, false, g.client.umask())
}

// parseMeta looks for the first meta tag in the given reader that
// will give us the source URL.
func (g *HttpGetter) parseMeta(r io.Reader) (string, error) {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	d.Strict = false
	var err error
	var t xml.Token
	for {
		t, err = d.Token()
		if err != nil {
			if err == io.EOF {
package getter

import "context"

// A ClientOption allows to configure a client
type ClientOption func(*Client) error

// Configure configures a client with options.
func (c *Client) Configure(opts ...ClientOption) error {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}
	c.Options = opts
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return err
		}
	}
	// Default decompressor values
	if c.Decompressors == nil {
		c.Decompressors = Decompressors
	}
	// Default detector values
	if c.Detectors == nil {
		c.Detectors = Detectors
	}
	// Default getter values
	if c.Getters == nil {
		c.Getters = Getters
	}

	for _, getter := range c.Getters {
		getter.SetClient(c)
	}
	return nil
}

// WithContext allows to pass a context to operation
// in order to be able to cancel a download in progress.
func WithContext(ctx context.Context) func(*Client) error {
	return func(c *Client) error {
		c.Ctx = ctx
		return nil
	}
}
	"runtime"
	"strconv"
	"strings"

	urlhelper "github.com/hashicorp/go-getter/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"
// a git repository.
type GitGetter struct {
	getter
}

var defaultBranchRegexp = regexp.MustCompile(`\s->\sorigin/(.*)`)

func (g *GitGetter) Get(dst string, u *url.URL) error {
	ctx := g.Context()
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git must be available and on the PATH")
	}
	var sshKeyFile string
	if sshKey != "" {
		// Check that the git version is sufficiently new.
		if err := checkGitVersion("2.3"); err != nil {
			return fmt.Errorf("Error using ssh key: %v", err)
		}


	// Next: check out the proper tag/branch if it is specified, and checkout
	if ref != "" {
		if err := g.checkout(dst, ref); err != nil {
			return err
		}
	}
	return fg.GetFile(dst, u)
}

func (g *GitGetter) checkout(dst string, ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = dst
	return getRunCommand(cmd)
}

	originalRef := ref // we handle an unspecified ref differently than explicitly selecting the default branch below
	if ref == "" {
		ref = findRemoteDefaultBranch(u)
	}
	if depth > 0 {
		args = append(args, "--depth", strconv.Itoa(depth))
		// If we didn't add --depth and --branch above then we will now be
		// on the remote repository's default branch, rather than the selected
		// ref, so we'll need to fix that before we return.
		return g.checkout(dst, originalRef)
	}
	return nil
}
		// Not a branch, switch to default branch. This will also catch
		// non-existent branches, in which case we want to switch to default
		// and then checkout the proper branch later.
		ref = findDefaultBranch(dst)
	}

	// We have to be on a branch to pull
	if err := g.checkout(dst, ref); err != nil {
		return err
	}

	if depth > 0 {
		cmd = exec.Command("git", "pull", "--depth", strconv.Itoa(depth), "--ff-only")
	} else {
		cmd = exec.Command("git", "pull", "--ff-only")
	}

	cmd.Dir = dst
// findDefaultBranch checks the repo's origin remote for its default branch
// (generally "master"). "master" is returned if an origin default branch
// can't be determined.
func findDefaultBranch(dst string) string {
	var stdoutbuf bytes.Buffer
	cmd := exec.Command("git", "branch", "-r", "--points-at", "refs/remotes/origin/HEAD")
	cmd.Dir = dst
	cmd.Stdout = &stdoutbuf
	err := cmd.Run()

// findRemoteDefaultBranch checks the remote repo's HEAD symref to return the remote repo's
// default branch. "master" is returned if no HEAD symref exists.
func findRemoteDefaultBranch(u *url.URL) string {
	var stdoutbuf bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--symref", u.String(), "HEAD")
	cmd.Stdout = &stdoutbuf
	err := cmd.Run()
	matches := lsRemoteSymRefRegexp.FindStringSubmatch(stdoutbuf.String())
// checkGitVersion is used to check the version of git installed on the system
// against a known minimum version. Returns an error if the installed version
// is older than the given minimum.
func checkGitVersion(min string) error {
	want, err := version.NewVersion(min)
	if err != nil {
		return err
	}

	out, err := exec.Command("git", "version").Output()
	if err != nil {
		return err
	}
