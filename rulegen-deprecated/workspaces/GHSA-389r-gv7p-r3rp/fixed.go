package main

package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	transport "github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/utils/trace"
)

// Err represents an HTTP error response.
type Err struct {
	URL    *url.URL
	Status int
	Reason string
}

// StatusCode returns the HTTP status code of the error.
func (e *Err) StatusCode() int { return e.Status }

func (e *Err) Error() string {
	format := "unexpected requesting %q status code: %d"
	if e.Reason != "" {
		return fmt.Sprintf(format+": %s", e.URL, e.Status, e.Reason)
	}
	return fmt.Sprintf(format, e.URL, e.Status)
}

// checkError maps HTTP response status codes to typed transport errors.
func checkError(r *http.Response) error {
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	var reason string
	var messageBuffer bytes.Buffer
	if r.Body != nil {
		messageLength, _ := messageBuffer.ReadFrom(r.Body)
		if messageLength > 0 {
			reason = messageBuffer.String()
		}
	}

	err := &Err{
		URL:    r.Request.URL,
		Status: r.StatusCode,
		Reason: reason,
	}

	switch r.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %w", transport.ErrAuthenticationRequired, err)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %w", transport.ErrAuthorizationFailed, err)
	case http.StatusNotFound:
		return fmt.Errorf("%w: %w", transport.ErrRepositoryNotFound, err)
	}

	return err
}

const infoRefsPath = "/info/refs"

// applyRedirect derives a new base URL from the final request URL after
// the HTTP client followed any redirects during the /info/refs GET.
//
// The logic mirrors canonical git's update_url_from_redirect(): strip
// the request-specific tail ("/info/refs") from the final URL to recover
// the new base. If the tail is missing, the redirect target is
// inconsistent and we return an error — canonical git die()s here
// because a mismatch could let a malicious server rewrite the base URL
// to an unrelated repository.
//
// Scheme is validated to prevent SSRF via unsupported protocols (e.g.
// a redirect to file:// or gopher://). Cross-scheme redirects only
// permit an upgrade from http to https; downgrades must not influence
// the session base URL used for subsequent requests.
func applyRedirect(resp *http.Response, baseURL *url.URL) (*url.URL, error) {
	if resp.Request == nil {
		return baseURL, nil
	}

	final := resp.Request.URL
	if !strings.HasSuffix(final.Path, infoRefsPath) {
		return nil, fmt.Errorf(
			"http transport: redirect target %q does not end with %s",
			final.Path, infoRefsPath,
		)
	}
	if final.Host == baseURL.Host &&
		final.Scheme == baseURL.Scheme &&
		strings.TrimSuffix(final.Path, infoRefsPath) == baseURL.Path {
		return baseURL, nil
	}

	if final.Scheme != "http" && final.Scheme != "https" {
		return nil, fmt.Errorf("http transport: redirect to unsupported scheme %q", final.Scheme)
	}
	if final.Scheme != baseURL.Scheme &&
		(baseURL.Scheme != "http" || final.Scheme != "https") {
		return nil, fmt.Errorf(
			"http transport: redirect changes scheme from %q to %q",
			baseURL.Scheme, final.Scheme,
		)
	}

	redirected := *baseURL
	redirected.Host = final.Host
	redirected.Scheme = final.Scheme
	redirected.Path = final.Path[:len(final.Path)-len(infoRefsPath)]
	return &redirected, nil
}

var safeHeaders = map[string]struct{}{
	"User-Agent":        {},
	"Host":              {},
	"Content-Encoding":  {},
}

func filterHeaders(h http.Header) http.Header {
	filtered := make(http.Header)
	for key, values := range h {
	return filtered
}

func redactedURL(u *url.URL) string {
	if u == nil {
		return ""
	if _, hasPassword := u.User.Password(); !hasPassword {
		return u.String()
	}
	redacted := *u
	redacted.User = url.UserPassword(u.User.Username(), "REDACTED")
	return redacted.String()
}

// doRequest performs an HTTP request and returns a typed error on failure.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	traceHTTP := trace.HTTP.Enabled()
	if traceHTTP {
		trace.HTTP.Printf("requesting %s %s %v", req.Method, redactedURL(req.URL), filterHeaders(req.Header))
	return res, checkError(res)
}

// applyAuth sets basic auth from URL userinfo and/or the authorizer function.
func applyAuth(httpReq *http.Request, baseURL *url.URL, authorizer func(*http.Request) error) error {
	if baseURL.User != nil {
		password, _ := baseURL.User.Password()
		httpReq.SetBasicAuth(baseURL.User.Username(), password)
	}
	if authorizer != nil {
		return authorizer(httpReq)
	}
	return nil
}
package revlist

import (
	"errors"
	"fmt"
	"sort"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
)

// objectWalk holds the state for a single Objects computation.
type objectWalk struct {
	s          storer.EncodedObjectStorer
	shallows   map[plumbing.Hash]struct{}
	wantsQueue []*object.Commit
	havesQueue []*object.Commit
	wantsSeen  map[plumbing.Hash]struct{}
	havesSeen  map[plumbing.Hash]struct{}
	seen       map[plumbing.Hash]struct{}
	result     []plumbing.Hash
}

func newObjectWalk(s storer.EncodedObjectStorer) (*objectWalk, error) {
	shallows, err := shallowSet(s)
	if err != nil {
		return nil, err
	}

	return &objectWalk{
		s:         s,
		shallows:  shallows,
		wantsSeen: make(map[plumbing.Hash]struct{}),
		havesSeen: make(map[plumbing.Hash]struct{}),
		seen:      make(map[plumbing.Hash]struct{}),
	}, nil
}

func shallowSet(s storer.EncodedObjectStorer) (map[plumbing.Hash]struct{}, error) {
	ss, ok := s.(storer.ShallowStorer)
	if !ok {
		return map[plumbing.Hash]struct{}{}, nil
	}

	hashes, err := ss.Shallow()
	if err != nil {
		return nil, err
	}

	set := make(map[plumbing.Hash]struct{}, len(hashes))
	for _, h := range hashes {
		set[h] = struct{}{}
	}

	return set, nil
}

// seedWants resolves each want hash and enqueues commits for walking.
// Non-commit objects (blobs, trees, tags) are added directly to the result.
func (w *objectWalk) seedWants(wants []plumbing.Hash) error {
	for i := 0; i < len(wants); i++ {
		h := wants[i]
		if _, ok := w.wantsSeen[h]; ok {
			continue
		}
		if _, ok := w.seen[h]; ok {
			continue
		}

		o, err := w.s.EncodedObject(plumbing.AnyObject, h)
		if err != nil {
			return fmt.Errorf("getting wanted object %s: %w", h, err)
		}

		switch o.Type() {
		case plumbing.CommitObject:
			c, err := object.DecodeCommit(w.s, o)
			if err != nil {
				return fmt.Errorf("decoding commit %s: %w", h, err)
			}
			w.wantsSeen[h] = struct{}{}
			insertSorted(&w.wantsQueue, c)
		case plumbing.TagObject:
			tag, err := object.DecodeTag(w.s, o)
			if err != nil {
				return fmt.Errorf("decoding tag %s: %w", h, err)
			}
			w.seen[tag.Hash] = struct{}{}
			w.result = append(w.result, tag.Hash)
			wants = append(wants, tag.Target)
		case plumbing.TreeObject:
			t, err := object.GetTree(w.s, h)
			if err != nil {
				return fmt.Errorf("getting tree %s: %w", h, err)
			}
			if err := collectAllTreeObjects(w.s, t, w.seen, &w.result); err != nil {
				return err
			}
		case plumbing.BlobObject:
			w.seen[h] = struct{}{}
			w.result = append(w.result, h)
		default:
			return fmt.Errorf("unsupported object type %s for %s", o.Type(), h)
		}
	}
	return nil
}

// seedHaves enqueues each have commit and pre-populates seen with all
// tree/blob objects reachable from the haves tips. Non-commit objects
// (tags, trees, blobs) are marked as seen so the diff walk skips them.
// Missing objects (ErrObjectNotFound) are tolerated since the remote
// may advertise refs we don't have locally.
func (w *objectWalk) seedHaves(haves []plumbing.Hash) error {
	for i := 0; i < len(haves); i++ {
		h := haves[i]
		if _, ok := w.havesSeen[h]; ok {
			continue
		}
		if _, ok := w.seen[h]; ok {
			continue
		}

		o, err := w.s.EncodedObject(plumbing.AnyObject, h)
		if err != nil {
			if errors.Is(err, plumbing.ErrObjectNotFound) {
				continue
			}
			return fmt.Errorf("getting haves object %s: %w", h, err)
		}

		switch o.Type() {
		case plumbing.CommitObject:
			c, err := object.DecodeCommit(w.s, o)
			if err != nil {
				return fmt.Errorf("decoding haves commit %s: %w", h, err)
			}
			w.havesSeen[h] = struct{}{}
			insertSorted(&w.havesQueue, c)
			if t, err := c.Tree(); err == nil {
				markTreeSeen(w.s, t, w.seen)
			}
		case plumbing.TagObject:
			tag, err := object.DecodeTag(w.s, o)
			if err != nil {
				return fmt.Errorf("decoding haves tag %s: %w", h, err)
			}
			w.seen[tag.Hash] = struct{}{}
			haves = append(haves, tag.Target)
		case plumbing.TreeObject:
			if t, err := object.GetTree(w.s, h); err == nil {
				markTreeSeen(w.s, t, w.seen)
			}
		case plumbing.BlobObject:
			w.seen[h] = struct{}{}
		}
	}
	return nil
}

// Paint flags for the commit walk. A commit reachable from both
// sides is a boundary; once all queue entries have both flags
// the walk can stop.
const (
	wantPaint uint8 = 1 << iota
	havePaint
)

// walk identifies new commits and collects their tree objects.
func (w *objectWalk) walk() error {
	// Fast path: no haves means we need all reachable objects.
	if len(w.havesQueue) == 0 {
		return w.walkFull()
	}

	// Phase 1: merge wants and haves into a single priority queue
	// sorted by committer time. Each commit carries paint flags that
	// propagate to parents. A commit painted from both sides is a
	// boundary. The walk stops when all queue entries are boundaries
	// (all stale), so only the overlap region is traversed.
	flags := make(map[plumbing.Hash]uint8)
	var queue []*object.Commit

	for _, c := range w.wantsQueue {
		flags[c.Hash] |= wantPaint
		insertSorted(&queue, c)
	}
	for _, c := range w.havesQueue {
		flags[c.Hash] |= havePaint
		insertSorted(&queue, c)
	}
	w.wantsQueue = nil
	w.havesQueue = nil

	var newCommits []*object.Commit
	var missing []missingParent

	for len(queue) > 0 {
		lc := queue[0]
		queue = queue[1:]

		f := flags[lc.Hash]

		// Want-only commit — tentatively new.
		if f == wantPaint {
			newCommits = append(newCommits, lc)
		}

		// Propagate this commit's flags to parents.
		if _, shallow := w.shallows[lc.Hash]; !shallow {
			if err := w.propagate(&queue, flags, &missing, lc, f); err != nil {
				return err
			}
		}

		// If all remaining queue entries have both flags, no new
		// commits can be discovered — stop early.
		if allStale(queue, flags) {
			break
		}
	}

	// Validate recorded missing parents now that all paint has settled.
	// A missing parent is tolerable only when its child has been painted
	// by haves (either directly, as a haves tip or ancestor, or via
	// later propagation that set havePaint on flags[child] before the
	// walk terminated). Otherwise the want side references history we
	// cannot traverse — match Git's behavior and error.
	for _, mp := range missing {
		if flags[mp.child]&havePaint != 0 {
			continue
		}
		return fmt.Errorf("commit %s has missing parent %s", mp.child, mp.hash)
	}

	// Phase 2: collect tree objects for new commits, skipping any
	// that were painted by haves after being added to newCommits.
	for _, lc := range newCommits {
		if flags[lc.Hash]&havePaint != 0 {
			continue
		}
		if err := w.processCommitTrees(lc); err != nil {
			return err
		}
	}
	return nil
}

// missingParent records a parent commit that could not be loaded during
// the painted walk, along with the child from which it was reached.
// Validation is deferred until the walk completes so that havePaint
// propagation from later iterations can mark the child (and its missing
// parent) as behind the haves boundary.
type missingParent struct {
	hash  plumbing.Hash
	child plumbing.Hash
}

// propagate adds the given flags to each parent commit. Parents that
// already have all the flags are skipped. Missing parents are not
// treated as fatal here; they are recorded and validated after the
// walk, where we can tell whether the child was eventually painted by
// haves (in which case the missing parent is behind the haves boundary
// and tolerable, matching Git's behavior).
func (w *objectWalk) propagate(queue *[]*object.Commit, flags map[plumbing.Hash]uint8, missing *[]missingParent, lc *object.Commit, f uint8) error {
	for _, ph := range lc.ParentHashes {
		pf := flags[ph]
		if pf|f == pf {
			continue // parent already has all our flags
		}
		flags[ph] = pf | f

		pc, err := object.GetCommit(w.s, ph)
		if err != nil {
			if errors.Is(err, plumbing.ErrObjectNotFound) {
				*missing = append(*missing, missingParent{hash: ph, child: lc.Hash})
				continue
			}
			return fmt.Errorf("getting parent commit %s: %w", ph, err)
		}
		insertSorted(queue, pc)
	}
	return nil
}

// allStale returns true when every commit in the queue has both paint
// flags, meaning all remaining commits are boundaries and no new
// commits can be discovered.
func allStale(queue []*object.Commit, flags map[plumbing.Hash]uint8) bool {
	for _, c := range queue {
		if flags[c.Hash]&wantPaint == 0 || flags[c.Hash]&havePaint == 0 {
			return false
		}
	}
	return true
}

// walkFull is the fast path when there are no haves. It walks all
// commits and collects every reachable tree/blob via a simple seen-set
// traversal — no per-commit tree diffs needed.
func (w *objectWalk) walkFull() error {
	for len(w.wantsQueue) > 0 {
		lc := w.wantsQueue[0]
		w.wantsQueue = w.wantsQueue[1:]

		if _, ok := w.seen[lc.Hash]; ok {
			continue
		}
		w.seen[lc.Hash] = struct{}{}
		w.result = append(w.result, lc.Hash)

		tree, err := lc.Tree()
		if err != nil {
			return fmt.Errorf("getting tree for %s: %w", lc.Hash, err)
		}

		if err := collectAllTreeObjects(w.s, tree, w.seen, &w.result); err != nil {
			return fmt.Errorf("collecting tree objects for %s: %w", lc.Hash, err)
		}

		if _, ok := w.shallows[lc.Hash]; ok {
			continue
		}

		for _, ph := range lc.ParentHashes {
			if _, ok := w.wantsSeen[ph]; ok {
				continue
			}
			w.wantsSeen[ph] = struct{}{}
			pc, err := object.GetCommit(w.s, ph)
			if err != nil {
				return fmt.Errorf("getting parent commit %s: %w", ph, err)
			}
			insertSorted(&w.wantsQueue, pc)
		}
	}
	return nil
}

// processCommitTrees collects new tree/blob objects for a commit by
// diffing its tree against its parents' trees.
func (w *objectWalk) processCommitTrees(lc *object.Commit) error {
	if _, ok := w.seen[lc.Hash]; !ok {
		w.seen[lc.Hash] = struct{}{}
		w.result = append(w.result, lc.Hash)
	}

	newTree, err := lc.Tree()
	if err != nil {
		return fmt.Errorf("getting tree for %s: %w", lc.Hash, err)
	}

	var oldTrees []*object.Tree
	for i := 0; i < lc.NumParents(); i++ {
		parent, err := lc.Parent(i)
		if err != nil {
			if errors.Is(err, plumbing.ErrObjectNotFound) {
				continue // parent may be beyond haves boundary
			}
			return fmt.Errorf("getting parent commit %s: %w", lc.ParentHashes[i], err)
		}
		pt, err := parent.Tree()
		if err != nil {
			return fmt.Errorf("getting parent tree for %s: %w", parent.Hash, err)
		}
		oldTrees = append(oldTrees, pt)
	}

	if err := collectChangedTreeObjects(w.s, newTree, oldTrees, w.seen, &w.result); err != nil {
		return fmt.Errorf("diffing trees for %s: %w", lc.Hash, err)
	}

	return nil
}

// insertSorted inserts a commit into a slice sorted by committer time
// descending (newest first).
func insertSorted(q *[]*object.Commit, c *object.Commit) {
	i := sort.Search(len(*q), func(i int) bool {
		return (*q)[i].Committer.When.Before(c.Committer.When)
	})
	*q = append(*q, nil)
	copy((*q)[i+1:], (*q)[i:])
	(*q)[i] = c
}

// collectChangedTreeObjects walks newTree, comparing entry hashes against
// all oldTrees. An entry is considered unchanged if any old tree contains the
// same name with the same hash. Only new or modified tree and blob hashes are
// added to result.
func collectChangedTreeObjects(
	s storer.EncodedObjectStorer,
	newTree *object.Tree,
	oldTrees []*object.Tree,
	seen map[plumbing.Hash]struct{},
	result *[]plumbing.Hash,
) error {
	// If newTree matches any old tree exactly, nothing changed.
	for _, ot := range oldTrees {
		if newTree.Hash == ot.Hash {
			return nil
		}
	}

	if _, ok := seen[newTree.Hash]; !ok {
		seen[newTree.Hash] = struct{}{}
		*result = append(*result, newTree.Hash)
	}

	// Build per-parent entry indexes for O(1) lookup.
	oldEntryMaps := make([]map[string]plumbing.Hash, len(oldTrees))
	for i, ot := range oldTrees {
		oldEntryMaps[i] = make(map[string]plumbing.Hash, len(ot.Entries))
		for _, e := range ot.Entries {
			oldEntryMaps[i][e.Name] = e.Hash
		}
	}

	for _, e := range newTree.Entries {
		// Skip blobs we've already collected. Directories are not skipped
		// here — a tree hash being "seen" means the hash itself was added
		// to the result, but a prior diff-walk may not have collected all
		// of its children. The recursive call handles dedup for trees.
		if e.Mode != filemode.Dir {
			if _, ok := seen[e.Hash]; ok {
				continue
			}
		}
		if e.Mode == filemode.Submodule {
			continue
		}

		// If same name has same hash in any parent tree, unchanged—skip.
		unchanged := false
		for _, m := range oldEntryMaps {
			if oh, ok := m[e.Name]; ok && oh == e.Hash {
				unchanged = true
				break
			}
		}
		if unchanged {
			continue
		}

		if e.Mode == filemode.Dir {
			// Recurse into changed subtree. Collect the old versions of
			// this subtree from all parents that have it.
			newSub, err := object.GetTree(s, e.Hash)
			if err != nil {
				return fmt.Errorf("getting subtree %s: %w", e.Hash, err)
			}
			var oldSubs []*object.Tree
			for _, m := range oldEntryMaps {
				if oh, ok := m[e.Name]; ok {
					if ot, err := object.GetTree(s, oh); err == nil {
						oldSubs = append(oldSubs, ot)
					}
				}
			}
			if err := collectChangedTreeObjects(s, newSub, oldSubs, seen, result); err != nil {
				return err
			}
		} else {
			seen[e.Hash] = struct{}{}
			*result = append(*result, e.Hash)
		}
	}

	return nil
}

// collectAllTreeObjects recursively walks a tree, adding all unseen
// tree and blob hashes to result. This is faster than collectChangedTreeObjects
// when we need all objects (no haves to diff against).
func collectAllTreeObjects(
	s storer.EncodedObjectStorer,
	t *object.Tree,
	seen map[plumbing.Hash]struct{},
	result *[]plumbing.Hash,
) error {
	if _, ok := seen[t.Hash]; ok {
		return nil
	}
	seen[t.Hash] = struct{}{}
	*result = append(*result, t.Hash)

	for _, e := range t.Entries {
		if e.Mode == filemode.Submodule {
			continue
		}
		if _, ok := seen[e.Hash]; ok {
			continue
		}
		if e.Mode == filemode.Dir {
			sub, err := object.GetTree(s, e.Hash)
			if err != nil {
				return fmt.Errorf("getting subtree %s: %w", e.Hash, err)
			}
			if err := collectAllTreeObjects(s, sub, seen, result); err != nil {
				return err
			}
		} else {
			seen[e.Hash] = struct{}{}
			*result = append(*result, e.Hash)
		}
	}
	return nil
}

// markTreeSeen recursively adds all tree and blob hashes in t to seen.
// Objects added to seen are not added to result — this is used to mark
// haves-reachable objects so they are skipped during diff walks.
func markTreeSeen(s storer.EncodedObjectStorer, t *object.Tree, seen map[plumbing.Hash]struct{}) {
	if _, ok := seen[t.Hash]; ok {
		return
	}
	seen[t.Hash] = struct{}{}
	for _, e := range t.Entries {
		if e.Mode == filemode.Submodule {
			continue
		}
		if _, ok := seen[e.Hash]; ok {
			continue
		}
		if e.Mode == filemode.Dir {
			if sub, err := object.GetTree(s, e.Hash); err == nil {
				markTreeSeen(s, sub, seen)
			}
		} else {
			seen[e.Hash] = struct{}{}
		}
	}
}
// Package archive provides archive generation functionality for git-upload-archive.
// It supports tar and zip formats with proper security controls.
package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/go-git/go-git/v6/storage"
)

// PAXGlobalHeader is the name used for PAX global extended headers in
// git-generated tar archives. This matches the name used by canonical git.
const PAXGlobalHeader = "pax_global_header"

// DefaultUmask is the default tar umask (002).
const DefaultUmask = 0o002

// Exported errors for archive operations.
var (
	// ErrOnlyRefNames is returned when raw hashes or relative expressions are used
	// but allowUnreachable is not set.
	ErrOnlyRefNames = errors.New("only ref names are allowed")

	// ErrRelativeExpressions is returned when relative ref expressions are used
	// but allowUnreachable is not set.
	ErrRelativeExpressions = errors.New("relative expressions are not allowed")

	// ErrObjectNotFound is returned when the specified object cannot be found.
	ErrObjectNotFound = errors.New("object not found")

	// ErrUnsupportedObjectType is returned when the object type is not supported for archiving.
	ErrUnsupportedObjectType = errors.New("unsupported object type for archive")

	// ErrPathNotFound is returned when a sub-path is not found in the tree.
	ErrPathNotFound = errors.New("path not found in tree")

	// ErrPathNotDirectory is returned when a sub-path is not a directory.
	ErrPathNotDirectory = errors.New("path is not a directory")

	// ErrPathspecNoMatch is returned when pathspec filters don't match any files.
	ErrPathspecNoMatch = errors.New("pathspec did not match any files")

	// ErrSymlinkTargetTooLarge is returned when a symlink blob exceeds the
	// maximum size accepted by the tar archive writer.
	ErrSymlinkTargetTooLarge = errors.New("symlink target too large")

	// ErrInvalidPrefix is returned when the requested archive prefix contains
	// path traversal sequences.
	ErrInvalidPrefix = errors.New("invalid archive prefix")
)

const maxTarSymlinkTargetSize = 64 * 1024

// SupportedFormats returns the list of supported archive formats.
func SupportedFormats() []string {
	return []string{"tar", "tar.gz", "tgz", "zip"}
}

// ApplyUmask applies umask to the given mode for regular files.
// Returns mode with all permission bits set, then applies umask.
func ApplyUmask(mode int64, isExecutable bool) int64 {
	if isExecutable {
		return (mode | 0o777) &^ DefaultUmask
	}
	return (mode | 0o666) &^ DefaultUmask
}

// ApplyUmaskDir applies umask to directories.
// Directories always get full permissions minus umask.
func ApplyUmaskDir(mode int64) int64 {
	return (mode | 0o777) &^ DefaultUmask
}

// ResolveTreeish resolves a tree-ish expression to a tree object.
//
// Security: By default, only direct ref names (v1.0, main) and ref:path
// sub-tree syntax (v1.0:Documentation) are allowed. Raw SHA-1 hashes and
// relative expressions (main^, HEAD~2) are rejected unless
// allowUnreachable is true. See https://git-scm.com/docs/git-upload-archive
//
// Returns the tree, commit hash (if applicable), commit time, and any error.
func ResolveTreeish(st storage.Storer, treeish string, allowUnreachable bool) (*object.Tree, *plumbing.Hash, time.Time, error) {
	var subPath string
	if idx := strings.IndexByte(treeish, ':'); idx >= 0 {
		subPath = treeish[idx+1:]
		treeish = treeish[:idx]
	}

	if !allowUnreachable {
		if plumbing.IsHash(treeish) {
			return nil, nil, time.Time{}, fmt.Errorf("%w (got %s)", ErrOnlyRefNames, treeish)
		}
		if strings.ContainsAny(treeish, "^~@{}") {
			return nil, nil, time.Time{}, fmt.Errorf("%w (got %s)", ErrRelativeExpressions, treeish)
		}
	}

	h, err := ResolveRef(st, treeish, allowUnreachable)
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	obj, err := object.GetObject(st, h)
	if err != nil {
		return nil, nil, time.Time{}, fmt.Errorf("%w: %s", ErrObjectNotFound, treeish)
	}

	for {
		tag, ok := obj.(*object.Tag)
		if !ok {
			break
		}

		obj, err = object.GetObject(st, tag.Target)
		if err != nil {
			return nil, nil, time.Time{}, fmt.Errorf("resolve annotated tag: %w", err)
		}
	}

	var commitHash *plumbing.Hash
	var commitTime time.Time
	var tree *object.Tree

	switch o := obj.(type) {
	case *object.Commit:
		commitHash = &o.Hash
		commitTime = o.Committer.When
		tree, err = o.Tree()
		if err != nil {
			return nil, nil, time.Time{}, err
		}
	case *object.Tree:
		tree = o

		// git archive behaves differently when given a tree ID versus when
		// given a commit ID or tag ID. In the first case the current time is
		// used as the modification time of each file in the archive. In the
		// latter case the commit time as recorded in the referenced commit
		// object is used instead.
		//
		// See https://git-scm.com/docs/git-archive
		commitTime = time.Now()
	default:
		return nil, nil, time.Time{}, fmt.Errorf("unsupported object type for archive: %T", obj)
	}

	if subPath != "" {
		entry, err := tree.FindEntry(subPath)
		if err != nil {
			return nil, nil, time.Time{}, fmt.Errorf("path not found in tree: %s", subPath)
		}
		if entry.Mode != filemode.Dir {
			return nil, nil, time.Time{}, fmt.Errorf("path is not a directory: %s", subPath)
		}
		tree, err = object.GetTree(st, entry.Hash)
		if err != nil {
			return nil, nil, time.Time{}, err
		}
	}

	return tree, commitHash, commitTime, nil
}

// ResolveRef resolves a ref name to a hash.
func ResolveRef(st storage.Storer, name string, allowHash bool) (plumbing.Hash, error) {
	if allowHash && plumbing.IsHash(name) {
		return plumbing.NewHash(name), nil
	}

	for _, candidate := range []plumbing.ReferenceName{
		plumbing.ReferenceName(name),
		plumbing.ReferenceName("refs/heads/" + name),
		plumbing.ReferenceName("refs/tags/" + name),
	} {
		ref, err := storer.ResolveReference(st, candidate)
		if err == nil {
			return ref.Hash(), nil
		}
	}

	return plumbing.ZeroHash, fmt.Errorf("cannot resolve %q", name)
}

// WriteTarArchive writes a tar archive from a tree.
func WriteTarArchive(st storage.Storer, w io.Writer, tree *object.Tree, commitHash *plumbing.Hash, prefix string, pathFilter []string, modTime time.Time) error {
	tw := tar.NewWriter(w)

	// Write PAX global extended header with commit ID if available.
	// This matches the behavior of git archive and allows extraction
	// via git get-tar-commit-id.
	if commitHash != nil {
		err := tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeXGlobalHeader,
			Name:     PAXGlobalHeader,
			PAXRecords: map[string]string{
				"comment": commitHash.String(),
			},
		})
		if err != nil {
			return fmt.Errorf("writing global PAX header: %w", err)
		}
	}

	if prefix != "" && strings.HasSuffix(prefix, "/") {
		if err := tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     prefix,
			Mode:     ApplyUmaskDir(0),
			ModTime:  modTime,
		}); err != nil {
			return err
		}
	}

	walker := object.NewTreeWalker(tree, true, nil)
	defer walker.Close()

	var matchedAny bool
	for {
		name, entry, err := walker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(pathFilter) > 0 && !MatchesPathFilter(name, pathFilter) {
			continue
		}
		matchedAny = true

		fullName := prefix + name

		// Extract Unix permission bits from git mode.
		unixMode := int64(entry.Mode) & 0o777

		if entry.Mode == filemode.Dir || entry.Mode == filemode.Submodule {
			if err := tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeDir,
				Name:     fullName + "/",
				Mode:     ApplyUmaskDir(unixMode),
				ModTime:  modTime,
			}); err != nil {
				return err
			}
			continue
		}

		blob, err := object.GetBlob(st, entry.Hash)
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name:    fullName,
			Size:    blob.Size,
			Mode:    unixMode,
			ModTime: modTime,
		}

		if entry.Mode == filemode.Symlink {
			rc, err := blob.Reader()
			if err != nil {
				return err
			}
			target, err := io.ReadAll(io.LimitReader(rc, maxTarSymlinkTargetSize+1))
			closeErr := rc.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
			if len(target) > maxTarSymlinkTargetSize {
				return fmt.Errorf("%w: %s (%d bytes)", ErrSymlinkTargetTooLarge, fullName, len(target))
			}
			hdr.Typeflag = tar.TypeSymlink
			hdr.Linkname = string(target)
			hdr.Size = 0
			// Symlinks always get 0777 per canonical git.
			hdr.Mode = 0o777
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			continue
		}

		isExec := entry.Mode == filemode.Executable
		hdr.Mode = ApplyUmask(unixMode, isExec)

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		rc, err := blob.Reader()
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, rc)
		closeErr := rc.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}

	if len(pathFilter) > 0 && !matchedAny {
		return fmt.Errorf("%w: '%s'", ErrPathspecNoMatch, strings.Join(pathFilter, " "))
	}

	return tw.Close()
}

// WriteZipArchive writes a zip archive from a tree.
func WriteZipArchive(st storage.Storer, w io.Writer, tree *object.Tree, commitHash *plumbing.Hash, prefix string, pathFilter []string, modTime time.Time) error {
	zw := zip.NewWriter(w)

	walker := object.NewTreeWalker(tree, true, nil)
	defer walker.Close()

	var matchedAny bool
	for {
		name, entry, err := walker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(pathFilter) > 0 && !MatchesPathFilter(name, pathFilter) {
			continue
		}
		matchedAny = true

		if entry.Mode == filemode.Dir || entry.Mode == filemode.Submodule {
			continue
		}

		fullName := prefix + name
		blob, err := object.GetBlob(st, entry.Hash)
		if err != nil {
			return err
		}

		// Extract Unix permission bits from git mode and apply default umask.
		unixMode := int64(entry.Mode) & 0o777

		fh := &zip.FileHeader{
			Name:     fullName,
			Method:   zip.Deflate,
			Modified: modTime,
		}
		switch entry.Mode {
		case filemode.Executable:
			fh.SetMode(fs.FileMode(ApplyUmask(unixMode, true)))
		case filemode.Symlink:
			// Zip stores symlinks with mode 0o120000 + permissions.
			fh.SetMode(fs.FileMode(0o120000 | (ApplyUmask(unixMode, true) & 0o777)))
		default:
			fh.SetMode(fs.FileMode(ApplyUmask(unixMode, false)))
		}

		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		rc, err := blob.Reader()
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, rc)
		closeErr := rc.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}

	if len(pathFilter) > 0 && !matchedAny {
		return fmt.Errorf("pathspec '%s' did not match any files", strings.Join(pathFilter, " "))
	}

	// Store commit ID as ZIP file comment if available.
	// This matches the behavior of git archive.
	if commitHash != nil {
		if err := zw.SetComment(commitHash.String()); err != nil {
			return err
		}
	}

	return zw.Close()
}

// MatchesPathFilter checks if a name matches any of the path filters.
//
// Note: This function assumes paths use forward slashes (/), which is the
// format used by Git internally. The TreeWalker produces paths with forward
// slashes regardless of the operating system, so this function is
// platform-independent.
//
// Supported patterns:
//   - Exact match: "README.md"
//   - Prefix match: "docs/" matches "docs/guide.md"
//   - Glob patterns: "*.go" matches "main.go"
//   - Parent-child: "docs/guide.md" matches parent "docs"
func MatchesPathFilter(name string, filters []string) bool {
	for _, f := range filters {
		if name == f || strings.HasPrefix(name, f+"/") || strings.HasPrefix(f, name+"/") {
			return true
		}
		matched, _ := path.Match(f, name)
		if matched {
			return true
		}
	}
	return false
}

// GetTarCommitID extracts the commit ID from a git-generated tar archive.
// It reads the PAX global extended header from the beginning of the archive
// and returns the value of the "comment" field, which contains the commit hash.
// If no global header is found or it doesn't contain a comment, it returns
// an error.
func GetTarCommitID(r io.Reader) (*plumbing.Hash, error) {
	tr := tar.NewReader(r)
	hdr, err := tr.Next()
	if err != nil {
		return nil, fmt.Errorf("reading tar header: %w", err)
	}
	// Check for PAX global extended header (typeflag 'g')
	if hdr.Typeflag != tar.TypeXGlobalHeader {
		return nil, fmt.Errorf("expected global PAX header, got typeflag %c", hdr.Typeflag)
	}
	// The PAX records should contain the commit ID in the "comment" field
	comment, ok := hdr.PAXRecords["comment"]
	if !ok {
		return nil, fmt.Errorf("global header missing comment field")
	}
	hash := plumbing.NewHash(comment)
	return &hash, nil
}

// WriteArchive generates an archive from the repository and writes it to w.
//
// Args follow the same format as git-archive: [options...] <tree-ish> [paths...]
// Supported formats: tar, zip, tar.gz, tgz.
// The prefix is prepended to all file paths in the archive.
// The paths slice can be used to filter which files are included.
// If allowUnreachable is false, only ref names are allowed (no raw hashes).
func WriteArchive(st storage.Storer, w io.Writer, treeish, format, prefix string,
	paths []string, allowUnreachable bool,
) error {
	if treeish == "" {
		return fmt.Errorf("no tree-ish specified")
	}
	if hasInvalidArchivePrefix(prefix) {
		return fmt.Errorf("%w: %s", ErrInvalidPrefix, prefix)
	}

	tree, commitHash, commitTime, err := ResolveTreeish(st, treeish, allowUnreachable)
	if err != nil {
		return err
	}

	switch format {
	case "tar":
		return WriteTarArchive(st, w, tree, commitHash, prefix, paths, commitTime)
	case "tar.gz", "tgz":
		gw := gzip.NewWriter(w)
		if err := WriteTarArchive(st, gw, tree, commitHash, prefix, paths, commitTime); err != nil {
			return err
		}
		return gw.Close()
	case "zip":
		return WriteZipArchive(st, w, tree, commitHash, prefix, paths, commitTime)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

func hasInvalidArchivePrefix(prefix string) bool {
	if strings.HasPrefix(prefix, "/") || strings.HasPrefix(prefix, "\\") {
		return true
	}
	return slices.Contains(strings.FieldsFunc(prefix, func(r rune) bool {
		return r == '/' || r == '\\'
	}), "..")
}
