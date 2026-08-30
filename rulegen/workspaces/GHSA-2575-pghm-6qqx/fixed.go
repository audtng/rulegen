package main

}

func (d *CachedDiscoveryClient) writeCachedFile(filename string, obj runtime.Object) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
		return err
	}

		return err
	}

	err = os.Chmod(f.Name(), 0660)
	if err != nil {
		return err
	}
import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	assert.Equal(c.groupCalls, 2)
}

func TestNewCachedDiscoveryClient_PathPerm(t *testing.T) {
	assert := assert.New(t)

	d, err := ioutil.TempDir("", "")
	assert.NoError(err)
	os.RemoveAll(d)
	defer os.RemoveAll(d)

	c := fakeDiscoveryClient{}
	cdc := newCachedDiscoveryClient(&c, d, 1*time.Nanosecond)
	cdc.ServerGroups()

	err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			assert.Equal(os.FileMode(0750), info.Mode().Perm())
		} else {
			assert.Equal(os.FileMode(0660), info.Mode().Perm())
		}
		return nil
	})
	assert.NoError(err)
}

type fakeDiscoveryClient struct {
	groupCalls    int
	resourceCalls int

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gregjones/httpcache"
// corresponding requests.
func newCacheRoundTripper(cacheDir string, rt http.RoundTripper) http.RoundTripper {
	d := diskv.New(diskv.Options{
		PathPerm: os.FileMode(0750),
		FilePerm: os.FileMode(0660),
		BasePath: cacheDir,
		TempDir:  filepath.Join(cacheDir, ".diskv-temp"),
	})
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// copied from k8s.io/client-go/transport/round_trippers_test.go
		t.Errorf("Invalid content read from cache %q", string(content))
	}
}

func TestCacheRoundTripperPathPerm(t *testing.T) {
	assert := assert.New(t)

	rt := &testRoundTripper{}
	cacheDir, err := ioutil.TempDir("", "cache-rt")
	os.RemoveAll(cacheDir)
	defer os.RemoveAll(cacheDir)

	if err != nil {
		t.Fatal(err)
	}
	cache := newCacheRoundTripper(cacheDir, rt)

	// First call, caches the response
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Host: "localhost"},
	}
	rt.Response = &http.Response{
		Header:     http.Header{"ETag": []string{`"123456"`}},
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("Content"))),
		StatusCode: http.StatusOK,
	}
	resp, err := cache.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "Content" {
		t.Errorf(`Expected Body to be "Content", got %q`, string(content))
	}

	err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			assert.Equal(os.FileMode(0750), info.Mode().Perm())
		} else {
			assert.Equal(os.FileMode(0660), info.Mode().Perm())
		}
		return nil
	})
	assert.NoError(err)
}
