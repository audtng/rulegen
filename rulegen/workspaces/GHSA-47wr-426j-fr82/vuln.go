package main

	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/datacharmer/dbdeployer/common"
	"github.com/datacharmer/dbdeployer/globals"
	"github.com/pkg/errors"
	"github.com/xi2/xz"
)

const (
	}
	// Create a tar Reader
	tr := tar.NewReader(r)
	return unpackTarFiles(tr)
}

func UnpackTar(filename string, destination string, verbosityLevel int) (err error) {
	} else {
		reader = tar.NewReader(fileReader)
	}
	return unpackTarFiles(reader)
}

func unpackTarFiles(reader *tar.Reader) (err error) {
	var header *tar.Header
	var count int = 0
	var reSlash = regexp.MustCompile(`/.*`)
			}
		case tar.TypeSymlink:
			if header.Linkname != "" {
				condPrint(fmt.Sprintf("%s -> %s", filename, header.Linkname), true, CHATTY)
				err = os.Symlink(header.Linkname, filename)
				if err != nil {
	// return nil
}

func unpackTarFile(filename string,
	reader *tar.Reader) (err error) {
	var writer *os.File
package common

// This file was generated during build. Do not edit.
// Build time: 2020-12-14 16:00

var VersionDef string = "1.58.1" // 2020-12-14

// Compatible version is the version used to mark compatible archives (templates, configuration).
// It is usually major.minor.0, except when we are at version 0.x, when
