package main


import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	safetemp "github.com/hashicorp/go-safetemp"
)

// ErrSymlinkCopy means that a copy of a symlink was encountered on a request with DisableSymlinks enabled.
var ErrSymlinkCopy = errors.New("copying of symlinks has been disabled")

// Client is a client for downloading things.
//
// Top-level functions such as Get are shortcuts for interacting with a client.
	// Getters is the list of protocols supported by this client. If this
	// is nil, then the default Getters variable will be used.
	Getters []Getter

	// Disable symlinks is used to prevent copying or writing files through symlinks for Get requests.
	// When set to true any copying or writing through symlinks will result in a ErrSymlinkCopy error.
	DisableSymlinks bool
}

// GetResult is the result of a Client.Get
		return nil, err
	}

	// Pass along the configured Getter client in the context for usage with the X-Terraform-Get feature.
	ctx = NewContextWithClient(ctx, c)

	// Store this locally since there are cases we swap this
	if req.GetMode == ModeInvalid {
		req.GetMode = ModeAny
	}

	// Client setting takes precedence for all requests
	if c.DisableSymlinks {
		req.DisableSymlinks = true
	}

	// If there is a subdir component, then we download the root separately
	// and then copy over the proper subdir.
	req.Src, req.subDir = SourceDirSubdir(req.Src)

	if req.subDir != "" {
		// Check if the subdirectory is attempting to traverse upwards, outside of
		// the cloned repository path.
		req.subDir = filepath.Clean(req.subDir)
		if containsDotDot(req.subDir) {
			return nil, fmt.Errorf("subdirectory component contain path traversal out of the repository")
		}

		// Prevent absolute paths, remove a leading path separator from the subdirectory
		if req.subDir[0] == os.PathSeparator {
			req.subDir = req.subDir[1:]
		}

		td, tdcloser, err := safetemp.Dir("", "getter")
		if err != nil {
			return nil, err
	// Determine if we have an archive type
	archiveV := q.Get("archive")
	if archiveV != "" {
		// Delete the parameter since it is a magic parameter we don't
		// want to pass on to the Getter
		q.Del("archive")
		req.u.RawQuery = q.Encode()
				filename = v
			}

			if containsDotDot(filename) {
				return nil, &getError{true, fmt.Errorf("filename query parameter contain path traversal")}
			}

			req.Dst = filepath.Join(req.Dst, filename)
		}
	}
			return nil, &getError{true, err}
		}

		err = copyDir(ctx, req.realDst, subDir, false, req.DisableSymlinks, req.umask())
		if err != nil {
			return nil, &getError{false, err}
		}
package getter

import (
	"context"
)

type clientContextKey int

const clientContextValue clientContextKey = 0

func NewContextWithClient(ctx context.Context, client *Client) context.Context {
	return context.WithValue(ctx, clientContextValue, client)
}

func ClientFromContext(ctx context.Context) *Client {
	// ctx.Value returns nil if ctx has no value for the key;
	client, ok := ctx.Value(clientContextValue).(*Client)
	if !ok {
		return nil
	}
	return client
}

// configure configures a client with options.
func (c *Client) configure() error {
	// Default decompressor values
func main() {
	modeRaw := flag.String("mode", "any", "get mode (any, file, dir)")
	progress := flag.Bool("progress", false, "display terminal progress")
	noSymlinks := flag.Bool("disable-symlinks", false, "prevent copying or writing files through symlinks")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
	if *progress {
		req.ProgressListener = defaultProgressBar
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	client := getter.DefaultClient

	// Disable symlinks for all client requests
	if *noSymlinks {
		client.DisableSymlinks = true
	}

	getters := getter.Getters
	getters = append(getters, new(gcs.Getter))
	getters = append(getters, new(s3.Getter))
// should already exist.
//
// If ignoreDot is set to true, then dot-prefixed files/folders are ignored.
func copyDir(ctx context.Context, dst string, src string, ignoreDot bool, disableSymlinks bool, umask os.FileMode) error {
	src, err := filepath.EvalSymlinks(src)
	if err != nil {
		return err
			}
		}

		if disableSymlinks {
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				return ErrSymlinkCopy
			}
		}

		// The "path" has the src prefixed to it. We need to join our
		// destination with the path without the src on it.
		dstPath := filepath.Join(dst, path[len(src):])
		}

		// If we have a file, copy the contents.
		_, err = copyFile(ctx, dstPath, path, disableSymlinks, info.Mode(), umask)
		return err
	}

)

func TestDetect(t *testing.T) {
	gitGetter := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(BitBucketDetector),
			new(GitHubDetector),
		},
	}
	cases := []struct {
		Input  string
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hashicorp/go-getter/v2"

// Getter is a Getter implementation that will download a module from
// a GCS bucket.
type Getter struct {

	// Timeout sets a deadline which all GCS operations should
	// complete within. Zero value means no timeout.
	Timeout time.Duration
}

func (g *Getter) Mode(ctx context.Context, u *url.URL) (getter.Mode, error) {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	// Parse URL
	bucket, object, err := g.parseURL(u)
	if err != nil {
}

func (g *Getter) Get(ctx context.Context, req *getter.Request) error {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	// Parse URL
	bucket, object, err := g.parseURL(req.URL())
	if err != nil {
}

func (g *Getter) GetFile(ctx context.Context, req *getter.Request) error {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	// Parse URL
	bucket, object, err := g.parseURL(req.URL())
	if err != nil {
	// The order of the Getters in the list may affect the result
	// depending if the Request.Src is detected as valid by multiple getters
	Getters = []Getter{
		&GitGetter{
			Detectors: []Detector{
				new(GitHubDetector),
				new(GitDetector),
				new(BitBucketDetector),
				new(GitLabDetector),
			},
		},
		new(HgGetter),
		new(SmbClientGetter),

import (
	"context"
	"fmt"
	"io"
	"os"
)
}

// copyFile copies a file in chunks from src path to dst path, using umask to create the dst file
func copyFile(ctx context.Context, dst, src string, disableSymlinks bool, fmode, umask os.FileMode) (int64, error) {

	if disableSymlinks {
		fileInfo, err := os.Lstat(src)
		if err != nil {
			return 0, fmt.Errorf("failed to check copy file source for symlinks: %w", err)
		}

		if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			return 0, ErrSymlinkCopy
		}
	}

	srcF, err := os.Open(src)
	if err != nil {
		return 0, err
	"runtime"
	"strconv"
	"strings"
	"time"

	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"
// a git repository.
type GitGetter struct {
	Detectors []Detector

	// Timeout sets a deadline which all hg CLI operations should
	// complete within. Defaults to zero which means no timeout.
	Timeout time.Duration
}

var defaultBranchRegexp = regexp.MustCompile(`\s->\sorigin/(.*)`)
		req.u.RawQuery = q.Encode()
	}

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	var sshKeyFile string
	if sshKey != "" {
		// Check that the git version is sufficiently new.
		if err := checkGitVersion(ctx, "2.3"); err != nil {
			return fmt.Errorf("Error using ssh key: %v", err)
		}


	// Next: check out the proper tag/branch if it is specified, and checkout
	if ref != "" {
		if err := g.checkout(ctx, req.Dst, ref); err != nil {
			return err
		}
	}
	return fg.GetFile(ctx, req)
}

func (g *GitGetter) checkout(ctx context.Context, dst string, ref string) error {
	cmd := exec.CommandContext(ctx, "git", "checkout", ref)
	cmd.Dir = dst
	return getRunCommand(cmd)
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
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	os.Setenv("PATH", dir)

	// Asking for a higher version throws an error
	ctx := context.Background()
	if err := checkGitVersion(ctx, "2.3"); err == nil {
		t.Fatal("expect git version error")
	}

	// Passes when version is satisfied
	if err := checkGitVersion(ctx, "1.9"); err != nil {
		t.Fatal(err)
	}
}

		GetMode: ModeDir,
	}
	getter := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(BitBucketDetector),
			new(GitHubDetector),
		},
	}
	client := &Client{
		Getters: []Getter{getter},
	}

	pwd := "/pwd"
	f := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(BitBucketDetector),
			new(GitHubDetector),
		},
	}
	for i, tc := range cases {
		req := &Request{
	}

	pwd := "/pwd"
	getter := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(BitBucketDetector),
			new(GitHubDetector),
		},
	}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
	}
}

func TestGitGetter_subdirectory_symlink(t *testing.T) {
	dst := testing_helper.TempDir(t)

	repo := testGitRepo(t, "repo-with-symlink")
	innerDir := filepath.Join(repo.dir, "this-directory-contains-a-symlink")
	if err := os.Mkdir(innerDir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(innerDir, "this-is-a-symlink")
	if err := os.Symlink("/etc/passwd", path); err != nil {
		t.Fatal(err)
	}
	repo.git("add", path)
	repo.git("commit", "-m", "Adding "+path)

	u, err := url.Parse(fmt.Sprintf("git::%s//this-directory-contains-a-symlink", repo.url.String()))
	if err != nil {
		t.Fatal(err)
	}

	req := &Request{
		Src:     u.String(),
		Dst:     dst,
		Pwd:     ".",
		GetMode: ModeDir,
	}
	getter := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(GitHubDetector),
		},
	}
	client := &Client{
		Getters:         []Getter{getter},
		DisableSymlinks: true,
	}

	ctx := context.Background()
	_, err = client.Get(ctx, req)
	if runtime.GOOS == "windows" {
		// Windows doesn't handle symlinks as one might expect with git.
		//
		// https://github.com/git-for-windows/git/wiki/Symbolic-Links
		filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, "this-is-a-symlink") {
				if info.Mode()&os.ModeSymlink == os.ModeSymlink {
					// If you see this test fail in the future, you've probably enabled
					// symlinks within git on your Windows system. Our CI/CD system does
					// not do this, so this is the only way we can make this test
					// make any sense.
					t.Fatalf("windows git should not have cloned a symlink")
				}
			}
			return nil
		})
	} else {
		// We can rely on POSIX compliant systems running git to do the right thing.
		if err == nil {
			t.Fatalf("expected client get to fail")
		}
		if !errors.Is(err, ErrSymlinkCopy) {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestGitGetter_subdirectory_traversal(t *testing.T) {
	dst := testing_helper.TempDir(t)

	repo := testGitRepo(t, "empty-repo")
	u, err := url.Parse(fmt.Sprintf("git::%s//../../../../../../etc/passwd", repo.url.String()))
	if err != nil {
		t.Fatal(err)
	}

	req := &Request{
		Src:     u.String(),
		Dst:     dst,
		Pwd:     ".",
		GetMode: ModeDir,
	}

	getter := &GitGetter{
		Detectors: []Detector{
			new(GitDetector),
			new(GitHubDetector),
		},
	}
	client := &Client{
		Getters: []Getter{getter},
	}

	ctx := context.Background()
	_, err = client.Get(ctx, req)
	if err == nil {
		t.Fatalf("expected client get to fail")
	}
	if !strings.Contains(err.Error(), "subdirectory component contain path traversal out of the repository") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// gitRepo is a helper struct which controls a single temp git repo.
type gitRepo struct {
	t   *testing.T
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"

// HgGetter is a Getter implementation that will download a module from
// a Mercurial repository.
type HgGetter struct {

	// Timeout sets a deadline which all hg CLI operations should
	// complete within. Defaults to zero which means no timeout.
	Timeout time.Duration
}

func (g *HgGetter) Mode(ctx context.Context, _ *url.URL) (Mode, error) {
	return ModeDir, nil
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	if err != nil {
		if err := g.clone(ctx, req.Dst, newURL); err != nil {
			return err
		}
	}

	if err := g.pull(ctx, req.Dst, newURL); err != nil {
		return err
	}

	return fg.GetFile(ctx, req)
}

func (g *HgGetter) clone(ctx context.Context, dst string, u *url.URL) error {
	cmd := exec.CommandContext(ctx, "hg", "clone", "-U", "--", u.String(), dst)
	return getRunCommand(cmd)
}

func (g *HgGetter) pull(ctx context.Context, dst string, u *url.URL) error {
	cmd := exec.CommandContext(ctx, "hg", "pull")
	cmd.Dir = dst
	return getRunCommand(cmd)
}
func (g *HgGetter) update(ctx context.Context, dst string, u *url.URL, rev string) error {
	args := []string{"update"}
	if rev != "" {
		args = append(args, "--", rev)
	}

	cmd := exec.CommandContext(ctx, "hg", args...)

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	testing_helper "github.com/hashicorp/go-getter/v2/helper/testing"
)
	}
	testing_helper.AssertContents(t, dst, "Hello\n")
}
func TestHgGetter_HgArgumentsNotAllowed(t *testing.T) {
	if !testHasHg {
		t.Log("hg not found, skipping")
		t.Skip()
	}
	ctx := context.Background()

	tc := []struct {
		name   string
		req    Request
		errChk func(testing.TB, error)
	}{
		{
			// If arguments are allowed in the destination, this request to Get will fail
			name: "arguments allowed in destination",
			req: Request{
				Dst: "--config=alias.clone=!touch ./TEST",
				u:   testModuleURL("basic-hg"),
			},
			errChk: func(t testing.TB, err error) {
				if err != nil {
					t.Errorf("Expected no err, got: %s", err)
				}
			},
		},
		{
			// Test arguments passed into the `rev` parameter
			// This clone call will fail regardless, but an exit code of 1 indicates
			// that the `false` command executed
			// We are expecting an hg parse error
			name: "arguments passed into rev parameter",
			req: Request{
				u: testModuleURL("basic-hg?rev=--config=alias.update=!false"),
			},
			errChk: func(t testing.TB, err error) {
				if err == nil {
					return
				}

				if !strings.Contains(err.Error(), "hg: parse error") {
					t.Errorf("Expected no err, got: %s", err)
				}
			},
		},
		{
			// Test arguments passed in the repository URL
			// This Get call will fail regardless, but it should fail
			// because the repository can't be found.
			// Other failures indicate that hg interpreted the argument passed in the URL
			name: "arguments passed in the repository URL",
			req: Request{
				u: &url.URL{Path: "--config=alias.clone=false"}},
			errChk: func(t testing.TB, err error) {
				if err == nil {
					return
				}

				if !strings.Contains(err.Error(), "repository --config=alias.clone=false not found") {
					t.Errorf("Expected no err, got: %s", err)
				}
			},
		},
	}
	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := new(HgGetter)

			if tt.req.Dst == "" {
				dst := testing_helper.TempDir(t)
				tt.req.Dst = dst
			}

			defer os.RemoveAll(tt.req.Dst)
			err := g.Get(ctx, &tt.req)
			tt.errChk(t, err)
		})
	}
}

func TestHgGetter_GetWithTimeout(t *testing.T) {
	if !testHasHg {
		t.Log("hg not found, skipping")
		t.Skip()
	}
	ctx := context.Background()
	g := &HgGetter{
		Timeout: 1 * time.Millisecond,
	}

	dst := testing_helper.TempDir(t)
	defer os.RemoveAll(filepath.Dir(dst))
	req := &Request{
		Dst: dst,
		u:   testModuleURL("basic-hg/foo.txt"),
	}

	if err := g.Get(ctx, req); err == nil {
		t.Fatalf("err: %s", err.Error())
	}
}
	"os"
	"path/filepath"
	"strings"
	"time"

	safetemp "github.com/hashicorp/go-safetemp"
)
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

func (g *HttpGetter) Mode(ctx context.Context, u *url.URL) (Mode, error) {
	return ModeFile, nil
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

func (g *HttpGetter) Get(ctx context.Context, req *Request) error {
	// Optionally disable any X-Terraform-Get redirects. This is recommended for usage of
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
	var newU url.URL = *req.u
	req.u = &newU
		}
	}

	// If the HTTP client is nil, check if there is one available in the context,
	// otherwise create one using cleanhttp's default transport.
	if g.Client == nil {
		if client := httpClientFromContext(ctx); client != nil {
			g.Client = client
		} else {
			g.Client = httpClient
		}
	}

	// Pass along the configured HTTP client in the context for usage with the X-Terraform-Get feature.
	ctx = context.WithValue(ctx, httpClientValue, g.Client)

	// Add terraform-get to the parameter.
	q := req.u.Query()
	q.Add("terraform-get", "1")
	req.u.RawQuery = q.Encode()

	readCtx := ctx
	if g.ReadTimeout > 0 {
		var cancel context.CancelFunc
		readCtx, cancel = context.WithTimeout(ctx, g.ReadTimeout)
		defer cancel()
	}

	// Get the URL
	httpReq, err := http.NewRequestWithContext(readCtx, "GET", req.u.String(), nil)
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

	// Get client with configured Getters from the context
	// If the client is nil, we know we're using the HttpGetter directly. In this case,
	// we don't know exactly which protocols are configured, but we can make a good guess.
	//
	// This prevents all default getters from being allowed when only using the
	// HttpGetter directly. To enable protocol switching, a client "wrapper" must
	// be used.
	var getterClient *Client
	if v := ClientFromContext(ctx); v != nil {
		getterClient = v
	} else {
		getterClient = &Client{
			Getters: []Getter{g},
		}
	}

	// Extract the source URL
	var source string
	if v := resp.Header.Get("X-Terraform-Get"); v != "" {
		source = v
	} else {
		source, err = g.parseMeta(readCtx, body)
		if err != nil {
			return err
		}
	}

	if source == "" {
		return fmt.Errorf("no source URL was returned")
	}

	return g.getXTerraformSource(ctx, req, source, getterClient)
}

// GetFile fetches the file from src and stores it at dst.
// falsely identified as being replaced, or corrupted with extra bytes
// appended.
func (g *HttpGetter) GetFile(ctx context.Context, req *Request) error {
	// Optionally enforce a maxiumum HTTP response body size.
	if g.MaxBytes > 0 {
		ctx = context.WithValue(ctx, httpMaxBytesValue, g.MaxBytes)
	}

	if g.Netrc {
		// Add auth from netrc if we can
		if err := addAuthFromNetrc(req.u); err != nil {
	}

	var currentFileSize int64
	var httpReq *http.Request

	if g.DoNotCheckHeadFirst == false {
		headCtx := ctx

		if g.HeadFirstTimeout > 0 {
			var cancel context.CancelFunc

			headCtx, cancel = context.WithTimeout(ctx, g.HeadFirstTimeout)
			defer cancel()
		}

		// We first make a HEAD request so we can check
		// if the server supports range queries. If the server/URL doesn't
		// support HEAD requests, we just fall back to GET.
		httpReq, err = http.NewRequestWithContext(headCtx, "HEAD", req.u.String(), nil)
		if err != nil {
			return err
		}
		if g.Header != nil {
			httpReq.Header = g.Header.Clone()
		}
		headResp, err := g.Client.Do(httpReq)
		if err == nil {
			headResp.Body.Close()
			if headResp.StatusCode == 200 {
				// If the HEAD request succeeded, then attempt to set the range
				// query if we can.
				if headResp.Header.Get("Accept-Ranges") == "bytes" && headResp.ContentLength >= 0 {
					if fi, err := f.Stat(); err == nil {
						if _, err = f.Seek(0, io.SeekEnd); err == nil {
							currentFileSize = fi.Size()
							httpReq.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentFileSize))
							if currentFileSize >= headResp.ContentLength {
								// file already present
								return nil
							}
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

	httpReq, err = http.NewRequestWithContext(readCtx, "GET", req.u.String(), nil)
	if err != nil {
		return err
	}
	if g.Header != nil {
		httpReq.Header = g.Header.Clone()
	}
	if currentFileSize > 0 {
		httpReq.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentFileSize))
	}

	resp, err := g.Client.Do(httpReq)
	if err != nil {

	body := resp.Body

	if maxBytes := httpMaxBytesFromContext(readCtx); maxBytes > 0 {
		body = newLimitedWrappedReaderCloser(body, maxBytes)
	}

	if req.ProgressListener != nil {
		// track download
		fn := filepath.Base(req.u.EscapedPath())
	}
	defer body.Close()

	n, err := Copy(readCtx, f, body)
	if err == nil && n < resp.ContentLength {
		err = io.ErrShortWrite
	}
	return err
}

// getXTerraformSource downloads the source into the destination
// using a protocol switching capable client.
func (g *HttpGetter) getXTerraformSource(ctx context.Context, req *Request, source string, client *Client) error {

	// If there is a subdir component, then we download the root separately
	// into a temporary directory, then copy over the proper subdir.
	source, subDir := SourceDirSubdir(source)
	req = &Request{
		GetMode:         ModeDir,
		Src:             source,
		Dst:             req.Dst,
		DisableSymlinks: req.DisableSymlinks,
	}

	if subDir == "" {
		// We have a X-Terraform-Get source lets check for supported Getters
		var allowed bool
		for _, getter := range client.Getters {
			shouldDownload, err := Detect(req, getter)
			if err != nil {
				return fmt.Errorf("failed to detect the proper Getter to handle %s: %w", source, err)
			}
			if !shouldDownload {
				// the request should not be processed by that getter
				continue
			}
			allowed = true
		}

		if !allowed {
			protocol := strings.Split(source, ":")[0]
			return fmt.Errorf("download not supported for scheme %q", protocol)
		}

		_, err := client.Get(ctx, req)
		return err
	}

	// We have a subdir, time to jump some hoops
	return g.getSubdir(ctx, req, source, subDir, client)

}

// getSubdir downloads the source into the destination, but with
// the proper subdir.
func (g *HttpGetter) getSubdir(ctx context.Context, req *Request, source, subDir string, client *Client) error {
	// Create a temporary directory to store the full source. This has to be
	// a non-existent directory.
	td, tdcloser, err := safetemp.Dir("", "getter")
	}
	defer tdcloser.Close()

	tdReq := &Request{
		Src:             source,
		Dst:             td,
		GetMode:         ModeDir,
		DisableSymlinks: req.DisableSymlinks,
	}
	if _, err := client.Get(ctx, tdReq); err != nil {
		return err
	}

		return err
	}

	return copyDir(ctx, req.Dst, sourcePath, false, req.DisableSymlinks, req.umask())
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
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	testing_helper "github.com/hashicorp/go-getter/v2/helper/testing"
)

	u.Path = "/header"

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "download not supported for scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta"

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "download not supported for scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta-subdir"

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "error downloading") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta-subdir-glob"

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "error downloading") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.User = url.UserPassword("foo", "bar")

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "download not supported for scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	defer tempEnv(t, "NETRC", path)()

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "download not supported for scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/header"

	req := &Request{
		Dst:     dst,
		Src:     u.String(),
		u:       &u,
		GetMode: ModeDir,
	}

	// Get it, which should error because it uses the file protocol.
	err := g.Get(ctx, req)
	if !strings.Contains(err.Error(), "download not supported for scheme") {
		t.Fatalf("unexpected error: %v", err)
	}
	// But, using a wrapper client with a file getter will work.
	c := &Client{
		Getters: []Getter{
			g,
			new(FileGetter),
		},
	}

	if _, err = c.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestHttpGetter__RespectsContextCanceled(t *testing.T) {
	}
}

func TestHttpGetter__XTerraformGetLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln := testHttpServerWithXTerraformGetLoop(t)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/loop"

	dst := testing_helper.TempDir(t)
	defer os.RemoveAll(dst)

	g := new(HttpGetter)
	g.XTerraformGetLimit = 10
	g.Client = &http.Client{}

	req := Request{
		Dst:     dst,
		u:       &u,
		GetMode: ModeDir,
	}

	err := g.Get(ctx, &req)
	if !strings.Contains(err.Error(), "too many X-Terraform-Get redirects") {
		t.Fatalf("too many X-Terraform-Get redirects, got: %v", err)
	}
}

func TestHttpGetter__XTerraformGetDisabled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln := testHttpServerWithXTerraformGetLoop(t)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/loop"
	dst := testing_helper.TempDir(t)

	g := new(HttpGetter)
	g.XTerraformGetDisabled = true
	g.Client = &http.Client{}

	req := Request{
		Dst:     dst,
		u:       &u,
		GetMode: ModeDir,
	}

	err := g.Get(ctx, &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestHttpGetter__XTerraformGetProxyBypass(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln := testHttpServerWithXTerraformGetProxyBypass(t)

	proxyLn := testHttpServerProxy(t, ln.Addr().String())

	t.Logf("starting malicious server on: %v", ln.Addr().String())
	t.Logf("starting proxy on: %v", proxyLn.Addr().String())

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/start"
	dst := testing_helper.TempDir(t)

	proxy, err := url.Parse(fmt.Sprintf("http://%s/", proxyLn.Addr().String()))
	if err != nil {
		t.Fatalf("failed to parse proxy URL: %v", err)
	}

	transport := cleanhttp.DefaultTransport()
	transport.Proxy = http.ProxyURL(proxy)

	g := new(HttpGetter)
	g.XTerraformGetLimit = 10
	g.Client = &http.Client{
		Transport: transport,
	}

	client := &Client{
		Getters: []Getter{g},
	}

	req := Request{
		Dst: dst,
		Src: u.String(),
	}

	_, err = client.Get(ctx, &req)
	if err != nil {
		t.Logf("client get error: %v", err)
	}
}

func TestHttpGetter__XTerraformGetConfiguredGettersBypass(t *testing.T) {
	tc := []struct {
		name              string
		configuredGetters []Getter
		errExpected       bool
	}{
		{name: "configured getter for git protocol switch", configuredGetters: []Getter{new(GitGetter)}, errExpected: false},
		{name: "configured getter for multiple protocol switch", configuredGetters: []Getter{new(GitGetter), new(HgGetter), new(FileGetter)}, errExpected: false},
		{name: "configured getter for file protocol switch", configuredGetters: []Getter{new(FileGetter)}, errExpected: true},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ln := testHttpServerWithXTerraformGetConfiguredGettersBypass(t)

			var u url.URL
			u.Scheme = "http"
			u.Host = ln.Addr().String()
			u.Path = "/start"

			dst := testing_helper.TempDir(t)

			rt := hookableHTTPRoundTripper{
				before: func(req *http.Request) {
					t.Logf("making request")
				},
				RoundTripper: http.DefaultTransport,
			}

			g := new(HttpGetter)
			g.XTerraformGetLimit = 10
			g.Client = &http.Client{
				Transport: &rt,
			}

			client := &Client{
				Getters: []Getter{g},
			}
			client.Getters = append(client.Getters, tt.configuredGetters...)

			t.Logf("%v", u.String())

			req := Request{
				Dst:     dst,
				Src:     u.String(),
				GetMode: ModeDir,
			}

			_, err := client.Get(ctx, &req)
			// For configured getters that support git, the git repository doesn't exist so error will not be nil.
			// If we get a nil error when we expect one other than the git error git exited with -1 we should fail.
			if tt.errExpected && err == nil {
				t.Fatalf("error expected")
			}
			// We only care about the error messages that indicate that we can download the git header URL
			if tt.errExpected && err != nil {
				if !strings.Contains(err.Error(), "download not supported for scheme") {
					t.Fatalf("expected download not supported for scheme, got: %v", err)
				}
			}
		})
	}
}

func TestHttpGetter__endless_body(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln := testHttpServerWithEndlessBody(t)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/"
	dst := testing_helper.TempDir(t)

	g := new(HttpGetter)
	g.MaxBytes = 10
	g.DoNotCheckHeadFirst = true

	client := &Client{
		Getters: []Getter{g},
	}

	t.Logf("%v", u.String())

	req := Request{
		Dst:     dst,
		Src:     u.String(),
		GetMode: ModeFile,
	}

	_, err := client.Get(ctx, &req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHttpGetter_subdirLink(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ln := testHttpServerSubDir(t)
	defer ln.Close()

	dst, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	t.Logf("dst: %q", dst)

	var u url.URL
	u.Scheme = "http"
	u.Host = ln.Addr().String()
	u.Path = "/regular-subdir//meta-subdir"

	g := new(HttpGetter)
	client := &Client{
		Getters: []Getter{g},
	}

	t.Logf("url: %q", u.String())

	req := Request{
		Dst:     dst,
		Src:     u.String(),
		GetMode: ModeAny,
	}

	_, err = client.Get(ctx, &req)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
}

func testHttpServerWithXTerraformGetLoop(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	header := fmt.Sprintf("http://%v:%v", ln.Addr().String(), "/loop")

	mux := http.NewServeMux()
	mux.HandleFunc("/loop", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Terraform-Get", header)
		t.Logf("serving loop")
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func testHttpServerWithXTerraformGetProxyBypass(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	header := fmt.Sprintf("http://%v/bypass", ln.Addr().String())

	mux := http.NewServeMux()
	mux.HandleFunc("/start/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Terraform-Get", header)
		t.Logf("serving start")
	})

	mux.HandleFunc("/bypass", func(w http.ResponseWriter, r *http.Request) {
		t.Fail()
		t.Logf("bypassed proxy")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("serving HTTP server path: %v", r.URL.Path)
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func testHttpServerWithXTerraformGetConfiguredGettersBypass(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	header := fmt.Sprintf("git::http://%v/some/repository.git", ln.Addr().String())

	mux := http.NewServeMux()
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Terraform-Get", header)
		t.Logf("serving start")
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("serving git HTTP server path: %v", r.URL.Path)
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func TestHttpGetter_XTerraformWithClientFromContext(t *testing.T) {
	tc := []struct {
		name        string
		client      *Client
		errExpected bool
	}{
		{
			name: "default getters",
			client: &Client{
				Getters: Getters,
			},
			errExpected: false,
		},
		{
			name: "client configured with needed getters",
			client: &Client{
				Getters: []Getter{
					new(HttpGetter),
					new(FileGetter),
				},
			},
			errExpected: false,
		},
		{
			name:        "nil client",
			errExpected: true,
		},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ln := testHttpServer(t)
			defer ln.Close()
			ctx := context.Background()

			g := new(HttpGetter)
			dst := testing_helper.TempDir(t)
			defer os.RemoveAll(dst)

			var u url.URL
			u.Scheme = "http"
			u.Host = ln.Addr().String()
			u.Path = "/header"

			req := &Request{
				Dst:     dst,
				Src:     u.String(),
				u:       &u,
				GetMode: ModeDir,
			}

			// Using a client stored in the ctx with a file getter should work
			ctx = NewContextWithClient(ctx, tt.client)

			err := g.Get(ctx, req)
			if tt.errExpected && err == nil {
				t.Fatalf("error expected")
			}

			if err != nil {
				if !strings.Contains(err.Error(), "download not supported for scheme") {
					t.Fatalf("expected download not supported for scheme, got: %v", err)
				}
				return
			}

			// Verify the main file exists
			mainPath := filepath.Join(dst, "main.tf")
			if _, err := os.Stat(mainPath); err != nil {
				t.Fatalf("err: %s", err)
			}
		})
	}
}

func testHttpServerProxy(t *testing.T, upstreamHost string) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Logf("serving proxy: %v: %#+v", r.URL.Path, r.Header)
		// create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(r.URL)
		// Note that ServeHttp is non blocking & uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func testHttpServer(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
	return ln
}

func testHttpServerWithEndlessBody(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		for {
			w.Write([]byte(".\n"))
		}
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

func testHttpHandlerExpectHeader(w http.ResponseWriter, r *http.Request) {
	if expected, ok := r.URL.Query()["expected"]; ok {
		if r.Header.Get(expected[0]) != "" {
	}
}

func testHttpServerSubDir(t *testing.T) net.Listener {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			t.Logf("serving: %v: %v: %#+[1]v", r.Method, r.URL.String(), r.Header)
		}
	})

	var server http.Server
	server.Handler = mux
	go server.Serve(ln)

	return ln
}

const testHttpMetaStr = `
<html>
<head>
	// Verify the main file exists
	testing_helper.AssertContents(t, dst, "Hello\n")
}
func TestGetFile_filename_path_traversal(t *testing.T) {
	dst := testing_helper.TempDir(t)
	u := testModule("basic-file/foo.txt")

	u += "?filename=../../../../../../../../../../../../../tmp/bar.txt"

	ctx := context.Background()
	op, err := GetAny(ctx, dst, u)

	if op != nil {
		t.Fatalf("unexpected op: %v", op)
	}

	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "filename query parameter contain path traversal") {
		t.Fatalf("unexpected err: %s", err)
	}
}

func TestGetFile_archiveChecksum(t *testing.T) {
	ctx := context.Background()
	}
}

func TestGetForcedGetter(t *testing.T) {
	type args struct {
		src string
	}
	}

	if !reflect.DeepEqual(data, []byte(contents)) {
		t.Fatalf("bad. expected:\n\n%q\n\nGot:\n\n%q", contents, string(data))
	}
}

	// By default a no op progress listener is used.
	ProgressListener ProgressTracker

	// Disable symlinks is used to prevent copying or writing files through symlinks.
	// When set to true any copying or writing through symlinks will result in a ErrSymlinkCopy error.
	DisableSymlinks bool

	u               *url.URL
	subDir, realDst string
}
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

// Getter is a Getter implementation that will download a module from
// a S3 bucket.
type Getter struct {

	// Timeout sets a deadline which all S3 operations should
	// complete within. Zero value means no timeout.
	Timeout time.Duration
}

func (g *Getter) Mode(ctx context.Context, u *url.URL) (getter.Mode, error) {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	// Parse URL
	region, bucket, path, _, creds, err := g.parseUrl(u)
	if err != nil {
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	}
	resp, err := client.ListObjectsWithContext(ctx, req)
	if err != nil {
		return 0, err
	}

func (g *Getter) Get(ctx context.Context, req *getter.Request) error {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	// Parse URL
	region, bucket, path, _, creds, err := g.parseUrl(req.URL())
	if err != nil {
			s3Req.Marker = aws.String(lastMarker)
		}

		resp, err := client.ListObjectsWithContext(ctx, s3Req)
		if err != nil {
			return err
		}
}

func (g *Getter) GetFile(ctx context.Context, req *getter.Request) error {

	if g.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.Timeout)
		defer cancel()
	}

	region, bucket, path, version, creds, err := g.parseUrl(req.URL())
	if err != nil {
		return err
		s3req.VersionId = aws.String(version)
	}

	resp, err := client.GetObjectWithContext(ctx, s3req)
	if err != nil {
		return err
	}
//
// The returned path is the full absolute path.
func SubdirGlob(dst, subDir string) (string, error) {
	pattern := filepath.Join(dst, subDir)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
