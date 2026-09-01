package main

	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/xi2/xz"

	"github.com/datacharmer/dbdeployer/common"
	"github.com/datacharmer/dbdeployer/globals"
)

const (
	}
	// Create a tar Reader
	tr := tar.NewReader(r)
	return unpackTarFiles(tr, destination)
}

func UnpackTar(filename string, destination string, verbosityLevel int) (err error) {
	} else {
		reader = tar.NewReader(fileReader)
	}
	return unpackTarFiles(reader, destination)
}

func unpackTarFiles(reader *tar.Reader, extractDir string) error {
	const errLinkedDirectoryOutside = "linked directory '%s' is outside the extraction directory"
	extractAbsDir, err := filepath.Abs(extractDir)
	if err != nil {
		return fmt.Errorf("error defining the absolute path of '%s': %s", extractDir, err)
	}
	var header *tar.Header
	var count int = 0
	var reSlash = regexp.MustCompile(`/.*`)
			}
		case tar.TypeSymlink:
			if header.Linkname != "" {
				linkDepth := pathDepth(header.Linkname)
				nameDepth := pathDepth(header.Name)
				if linkDepth > nameDepth {
					fmt.Println()
					return fmt.Errorf(errLinkedDirectoryOutside, header.Linkname)
				}
				if common.FileExists(header.Linkname) {
					absFile, err := filepath.Abs(header.Linkname)
					if err != nil {
						return fmt.Errorf("error retrieving absolute path of %s: %s", header.Linkname, err)
					}
					if !common.BeginsWith(absFile, extractAbsDir) {
						return fmt.Errorf(errLinkedDirectoryOutside, header.Linkname)
					}
				} else {
					if common.BeginsWith(header.Linkname, "/") {
						if !common.BeginsWith(header.Linkname, extractAbsDir) {
							return fmt.Errorf(errLinkedDirectoryOutside, header.Linkname)
						}
					}
				}
				condPrint(fmt.Sprintf("%s -> %s", filename, header.Linkname), true, CHATTY)
				err = os.Symlink(header.Linkname, filename)
				if err != nil {
	// return nil
}

func pathDepth(s string) int {
	reSlash := regexp.MustCompilePOSIX("(/)")
	list := reSlash.FindAllStringIndex(s, -1)
	return len(list)
}

func unpackTarFile(filename string,
	reader *tar.Reader) (err error) {
	var writer *os.File
package common

// This file was generated during build. Do not edit.
// Build time: 2020-12-15 10:49

var VersionDef string = "1.58.2" // 2020-12-15

// Compatible version is the version used to mark compatible archives (templates, configuration).
// It is usually major.minor.0, except when we are at version 0.x, when
