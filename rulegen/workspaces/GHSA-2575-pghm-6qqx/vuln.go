package main

}

func (d *CachedDiscoveryClient) writeCachedFile(filename string, obj runtime.Object) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

		return err
	}

	err = os.Chmod(f.Name(), 0755)
	if err != nil {
		return err
	}
import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	assert.Equal(c.groupCalls, 2)
}

type fakeDiscoveryClient struct {
	groupCalls    int
	resourceCalls int

import (
	"net/http"
	"path/filepath"

	"github.com/gregjones/httpcache"
// corresponding requests.
func newCacheRoundTripper(cacheDir string, rt http.RoundTripper) http.RoundTripper {
	d := diskv.New(diskv.Options{
		BasePath: cacheDir,
		TempDir:  filepath.Join(cacheDir, ".diskv-temp"),
	})
	"net/http"
	"net/url"
	"os"
	"testing"
)

// copied from k8s.io/client-go/transport/round_trippers_test.go
		t.Errorf("Invalid content read from cache %q", string(content))
	}
}
