package main

	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
			urlpath = r.URL.Path
		}
		urlpath = strings.Trim(urlpath, "/")
		// Reject any non-canonical path, in particular one containing ".."
		// traversal elements.
		if urlpath != "" && path.Clean(urlpath) != urlpath {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
