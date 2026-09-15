package main

// simple default versioning scheme.
package versioner

import (
	"os"
	"path/filepath"
	"runtime"
)

type Versioner interface {
	Archive(filePath string) error
}
	TimeFormat = "20060102-150405"
	TimeGlob   = "[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9]" // glob pattern matching TimeFormat
)

func cleanSymlinks(dir string) {
	if runtime.GOOS == "windows" {
		// We don't do symlinks on Windows. Additionally, there may
		// be things that look like symlinks that are not, which we
		// should leave alone. Deduplicated files, for example.
		return
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			l.Infoln("Removing incorrectly versioned symlink", path)
			os.Remove(path)
			return filepath.SkipDir
		}
		return nil
	})
}
}

func NewExternal(folderID, folderPath string, params map[string]string) Versioner {
	cleanSymlinks(folderPath)

	command := params["command"]

	s := External{
// Archive moves the named file away to a version archive. If this function
// returns nil, the named file does not exist any more (has been archived).
func (v External) Archive(filePath string) error {
	info, err := osutil.Lstat(filePath)
	if os.IsNotExist(err) {
		l.Debugln("not archiving nonexistent file", filePath)
		return nil
	} else if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		panic("bug: attempting to version a symlink")
	}

	l.Debugln("archiving", filePath)

}

func NewStaggered(folderID, folderPath string, params map[string]string) Versioner {
	cleanSymlinks(folderPath)

	maxAge, err := strconv.ParseInt(params["maxAge"], 10, 0)
	if err != nil {
		maxAge = 31536000 // Default: ~1 year
	v.mutex.Lock()
	defer v.mutex.Unlock()

	info, err := osutil.Lstat(filePath)
	if os.IsNotExist(err) {
		l.Debugln("not archiving nonexistent file", filePath)
		return nil
	} else if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		panic("bug: attempting to version a symlink")
	}

	if _, err := os.Stat(v.versionsPath); err != nil {
		if os.IsNotExist(err) {
