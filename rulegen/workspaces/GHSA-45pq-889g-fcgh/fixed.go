package main

	"encoding/json"
	"errors"
	"fmt"
	iofs "io/fs"
	"net"
	"net/http"
	"os"
			urlpath = r.URL.Path
		}
		urlpath = strings.Trim(urlpath, "/")
		// Reject anything which isn't a canonical relative path free of "."
		// and ".." elements. The backends join the path with the Fs root, so
		// such elements could otherwise address objects outside it.
		//
		// The empty path is the root of the API so is allowed. "." is not a
		// valid object name here, even though iofs.ValidPath accepts it as
		// the root of an FS.
		if urlpath != "" && (urlpath == "." || !iofs.ValidPath(urlpath)) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
