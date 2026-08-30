package main


import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	safetemp "github.com/hashicorp/go-safetemp"
)

// Client is a client for downloading things.
//
// Top-level functions such as Get are shortcuts for interacting with a client.
	// Getters is the list of protocols supported by this client. If this
	// is nil, then the default Getters variable will be used.
	Getters []Getter
}

// GetResult is the result of a Client.Get
		return nil, err
	}

	// Store this locally since there are cases we swap this
	if req.GetMode == ModeInvalid {
		req.GetMode = ModeAny
	}

	// If there is a subdir component, then we download the root separately
	// and then copy over the proper subdir.
	req.Src, req.subDir = SourceDirSubdir(req.Src)
	if req.subDir != "" {
		td, tdcloser, err := safetemp.Dir("", "getter")
		if err != nil {
			return nil, err
	// Determine if we have an archive type
	archiveV := q.Get("archive")
	if archiveV != "" {
		// Delete the paramter since it is a magic parameter we don't
		// want to pass on to the Getter
		q.Del("archive")
		req.u.RawQuery = q.Encode()
				filename = v
			}

			req.Dst = filepath.Join(req.Dst, filename)
		}
	}
			return nil, &getError{true, err}
		}

		err = copyDir(ctx, req.realDst, subDir, false, req.umask())
		if err != nil {
			return nil, &getError{false, err}
		}
package getter

// configure configures a client with options.
func (c *Client) configure() error {
	// Default decompressor values
func main() {
	modeRaw := flag.String("mode", "any", "get mode (any, file, dir)")
	progress := flag.Bool("progress", false, "display terminal progress")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
	if *progress {
		req.ProgressListener = defaultProgressBar
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	client := getter.DefaultClient

	getters := getter.Getters
	getters = append(getters, new(gcs.Getter))
	getters = append(getters, new(s3.Getter))
// should already exist.
//
// If ignoreDot is set to true, then dot-prefixed files/folders are ignored.
func copyDir(ctx context.Context, dst string, src string, ignoreDot bool, umask os.FileMode) error {
	src, err := filepath.EvalSymlinks(src)
	if err != nil {
		return err
			}
		}

		// The "path" has the src prefixed to it. We need to join our
		// destination with the path without the src on it.
		dstPath := filepath.Join(dst, path[len(src):])
		}

		// If we have a file, copy the contents.
		_, err = copyFile(ctx, dstPath, path, info.Mode(), umask)
		return err
	}

)

func TestDetect(t *testing.T) {
	gitGetter := &GitGetter{[]Detector{
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

	"cloud.google.com/go/storage"
	"github.com/hashicorp/go-getter/v2"

// Getter is a Getter implementation that will download a module from
// a GCS bucket.
type Getter struct{}

func (g *Getter) Mode(ctx context.Context, u *url.URL) (getter.Mode, error) {

	// Parse URL
	bucket, object, err := g.parseURL(u)
	if err != nil {
}

func (g *Getter) Get(ctx context.Context, req *getter.Request) error {
	// Parse URL
	bucket, object, err := g.parseURL(req.URL())
	if err != nil {
}

func (g *Getter) GetFile(ctx context.Context, req *getter.Request) error {
	// Parse URL
	bucket, object, err := g.parseURL(req.URL())
	if err != nil {
	// The order of the Getters in the list may affect the result
	// depending if the Request.Src is detected as valid by multiple getters
	Getters = []Getter{
		&GitGetter{[]Detector{
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
	"io"
	"os"
)
}

// copyFile copies a file in chunks from src path to dst path, using umask to create the dst file
func copyFile(ctx context.Context, dst, src string, fmode, umask os.FileMode) (int64, error) {
	srcF, err := os.Open(src)
	if err != nil {
		return 0, err
	"runtime"
	"strconv"
	"strings"

	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"
// a git repository.
type GitGetter struct {
	Detectors []Detector
}

var defaultBranchRegexp = regexp.MustCompile(`\s->\sorigin/(.*)`)
		req.u.RawQuery = q.Encode()
	}

	var sshKeyFile string
	if sshKey != "" {
		// Check that the git version is sufficiently new.
		if err := checkGitVersion("2.3"); err != nil {
			return fmt.Errorf("Error using ssh key: %v", err)
		}


	// Next: check out the proper tag/branch if it is specified, and checkout
	if ref != "" {
		if err := g.checkout(req.Dst, ref); err != nil {
			return err
		}
	}
	return fg.GetFile(ctx, req)
}

func (g *GitGetter) checkout(dst string, ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = dst
	return getRunCommand(cmd)
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
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"os"
	os.Setenv("PATH", dir)

	// Asking for a higher version throws an error
	if err := checkGitVersion("2.3"); err == nil {
		t.Fatal("expect git version error")
	}

	// Passes when version is satisfied
	if err := checkGitVersion("1.9"); err != nil {
		t.Fatal(err)
	}
}

		GetMode: ModeDir,
	}
	getter := &GitGetter{[]Detector{
		new(GitDetector),
		new(BitBucketDetector),
		new(GitHubDetector),
	},
	}
	client := &Client{
		Getters: []Getter{getter},
	}

	pwd := "/pwd"
	f := &GitGetter{[]Detector{
		new(GitDetector),
		new(BitBucketDetector),
		new(GitHubDetector),
	},
	}
	for i, tc := range cases {
		req := &Request{
	}

	pwd := "/pwd"
	getter := &GitGetter{[]Detector{
		new(GitDetector),
		new(BitBucketDetector),
		new(GitHubDetector),
	},
	}
	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
	}
}

// gitRepo is a helper struct which controls a single temp git repo.
type gitRepo struct {
	t   *testing.T
	"os/exec"
	"path/filepath"
	"runtime"

	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	safetemp "github.com/hashicorp/go-safetemp"

// HgGetter is a Getter implementation that will download a module from
// a Mercurial repository.
type HgGetter struct{}

func (g *HgGetter) Mode(ctx context.Context, _ *url.URL) (Mode, error) {
	return ModeDir, nil
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err != nil {
		if err := g.clone(req.Dst, newURL); err != nil {
			return err
		}
	}

	if err := g.pull(req.Dst, newURL); err != nil {
		return err
	}

	return fg.GetFile(ctx, req)
}

func (g *HgGetter) clone(dst string, u *url.URL) error {
	cmd := exec.Command("hg", "clone", "-U", u.String(), dst)
	return getRunCommand(cmd)
}

func (g *HgGetter) pull(dst string, u *url.URL) error {
	cmd := exec.Command("hg", "pull")
	cmd.Dir = dst
	return getRunCommand(cmd)
}
func (g *HgGetter) update(ctx context.Context, dst string, u *url.URL, rev string) error {
	args := []string{"update"}
	if rev != "" {
		args = append(args, rev)
	}

	cmd := exec.CommandContext(ctx, "hg", args...)

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testing_helper "github.com/hashicorp/go-getter/v2/helper/testing"
)
	}
	testing_helper.AssertContents(t, dst, "Hello\n")
}
	"os"
	"path/filepath"
	"strings"

	safetemp "github.com/hashicorp/go-safetemp"
)
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

func (g *HttpGetter) Mode(ctx context.Context, u *url.URL) (Mode, error) {
	return ModeFile, nil
}

func (g *HttpGetter) Get(ctx context.Context, req *Request) error {
	// Copy the URL so we can modify it
	var newU url.URL = *req.u
	req.u = &newU
		}
	}

	if g.Client == nil {
		g.Client = httpClient
	}

	// Add terraform-get to the parameter.
	q := req.u.Query()
	q.Add("terraform-get", "1")
	req.u.RawQuery = q.Encode()

	// Get the URL
	httpReq, err := http.NewRequestWithContext(ctx, "GET", req.u.String(), nil)
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
	}
	if source == "" {
		return fmt.Errorf("no source URL was returned")
	}

	// If there is a subdir component, then we download the root separately
	// into a temporary directory, then copy over the proper subdir.
	source, subDir := SourceDirSubdir(source)
	req = &Request{
		GetMode: ModeDir,
		Src:     source,
		Dst:     req.Dst,
	}
	if subDir == "" {
		_, err = DefaultClient.Get(ctx, req)
		return err
	}
	// We have a subdir, time to jump some hoops
	return g.getSubdir(ctx, req, source, subDir)
}

// GetFile fetches the file from src and stores it at dst.
// falsely identified as being replaced, or corrupted with extra bytes
// appended.
func (g *HttpGetter) GetFile(ctx context.Context, req *Request) error {
	if g.Netrc {
		// Add auth from netrc if we can
		if err := addAuthFromNetrc(req.u); err != nil {
	}

	var currentFileSize int64

	// We first make a HEAD request so we can check
	// if the server supports range queries. If the server/URL doesn't
	// support HEAD requests, we just fall back to GET.
	httpReq, err := http.NewRequestWithContext(ctx, "HEAD", req.u.String(), nil)
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
	httpReq.Method = "GET"

	resp, err := g.Client.Do(httpReq)
	if err != nil {

	body := resp.Body

	if req.ProgressListener != nil {
		// track download
		fn := filepath.Base(req.u.EscapedPath())
	}
	defer body.Close()

	n, err := Copy(ctx, f, body)
	if err == nil && n < resp.ContentLength {
		err = io.ErrShortWrite
	}
	return err
}

// getSubdir downloads the source into the destination, but with
// the proper subdir.
func (g *HttpGetter) getSubdir(ctx context.Context, req *Request, source, subDir string) error {
	// Create a temporary directory to store the full source. This has to be
	// a non-existent directory.
	td, tdcloser, err := safetemp.Dir("", "getter")
	}
	defer tdcloser.Close()

	// Download that into the given directory
	if _, err := Get(ctx, td, source); err != nil {
		return err
	}

		return err
	}

	return copyDir(ctx, req.Dst, sourcePath, false, req.umask())
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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	testing_helper "github.com/hashicorp/go-getter/v2/helper/testing"
)

	u.Path = "/header"

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta"

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta-subdir"

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/meta-subdir-glob"

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.User = url.UserPassword("foo", "bar")

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	defer tempEnv(t, "NETRC", path)()

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}

	u.Path = "/header"

	req := &Request{
		Dst: dst,
		u:   &u,
	}

	// Get it!
	if err := g.Get(ctx, req); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestHttpGetter__RespectsContextCanceled(t *testing.T) {
	}
}

func testHttpServer(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
	return ln
}

func testHttpHandlerExpectHeader(w http.ResponseWriter, r *http.Request) {
	if expected, ok := r.URL.Query()["expected"]; ok {
		if r.Header.Get(expected[0]) != "" {
	}
}

const testHttpMetaStr = `
<html>
<head>
	// Verify the main file exists
	testing_helper.AssertContents(t, dst, "Hello\n")
}

func TestGetFile_archiveChecksum(t *testing.T) {
	ctx := context.Background()
	}
}

func TestgetForcedGetter(t *testing.T) {
	type args struct {
		src string
	}
	}

	if !reflect.DeepEqual(data, []byte(contents)) {
		t.Fatalf("bad. expected:\n\n%s\n\nGot:\n\n%s", contents, string(data))
	}
}

	// By default a no op progress listener is used.
	ProgressListener ProgressTracker

	u               *url.URL
	subDir, realDst string
}
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

// Getter is a Getter implementation that will download a module from
// a S3 bucket.
type Getter struct{}

func (g *Getter) Mode(ctx context.Context, u *url.URL) (getter.Mode, error) {
	// Parse URL
	region, bucket, path, _, creds, err := g.parseUrl(u)
	if err != nil {
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	}
	resp, err := client.ListObjects(req)
	if err != nil {
		return 0, err
	}

func (g *Getter) Get(ctx context.Context, req *getter.Request) error {

	// Parse URL
	region, bucket, path, _, creds, err := g.parseUrl(req.URL())
	if err != nil {
			s3Req.Marker = aws.String(lastMarker)
		}

		resp, err := client.ListObjects(s3Req)
		if err != nil {
			return err
		}
}

func (g *Getter) GetFile(ctx context.Context, req *getter.Request) error {
	region, bucket, path, version, creds, err := g.parseUrl(req.URL())
	if err != nil {
		return err
		s3req.VersionId = aws.String(version)
	}

	resp, err := client.GetObject(s3req)
	if err != nil {
		return err
	}
//
// The returned path is the full absolute path.
func SubdirGlob(dst, subDir string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(dst, subDir))
	if err != nil {
		return "", err
	}
