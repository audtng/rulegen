package main

	f.fileData[name] = data
}

func (f *fakeConnection) deleteFile(name string) {
	f.mut.Lock()
	defer f.mut.Unlock()

	for i, fi := range f.files {
		if fi.Name == name {
			fi.Deleted = true
			fi.ModifiedS = time.Now().Unix()
			fi.Version = fi.Version.Update(f.id.Short())
			fi.Sequence = time.Now().UnixNano()
			fi.Blocks = nil

			f.files = append(append(f.files[:i], f.files[i+1:]...), fi)
			return
		}
	}
}

func (f *fakeConnection) sendIndexUpdate() {
	f.model.IndexUpdate(f.id, f.folder, f.files)
}
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	}
}

func TestRequestVersioningSymlinkAttack(t *testing.T) {
	// Sets up a folder with trashcan versioning and tries to use a
	// deleted symlink to escape

	cfg := defaultConfig.RawCopy()
	cfg.Folders[0] = config.NewFolderConfiguration("default", "_tmpfolder")
	cfg.Folders[0].PullerSleepS = 1
	cfg.Folders[0].Devices = []config.FolderDeviceConfiguration{
		{DeviceID: device1},
		{DeviceID: device2},
	}
	cfg.Folders[0].Versioning = config.VersioningConfiguration{
		Type: "trashcan",
	}
	w := config.Wrap("/tmp/cfg", cfg)

	db := db.OpenMemory()
	m := NewModel(w, device1, "syncthing", "dev", db, nil)
	m.AddFolder(cfg.Folders[0])
	m.ServeBackground()
	m.StartFolder("default")
	defer m.Stop()

	defer os.RemoveAll("_tmpfolder")

	fc := addFakeConn(m, device2)
	fc.folder = "default"

	// Create a temporary directory that we will use as target to see if
	// we can escape to it
	tmpdir, err := ioutil.TempDir("", "syncthing-test")
	if err != nil {
		t.Fatal(err)
	}

	// We listen for incoming index updates and trigger when we see one for
	// the expected test file.
	idx := make(chan int)
	fc.mut.Lock()
	fc.indexFn = func(folder string, fs []protocol.FileInfo) {
		idx <- len(fs)
	}
	fc.mut.Unlock()

	// Send an update for the test file, wait for it to sync and be reported back.
	fc.addFile("foo", 0644, protocol.FileInfoTypeSymlink, []byte(tmpdir))
	fc.sendIndexUpdate()

	for updates := 0; updates < 1; updates += <-idx {
	}

	// Delete the symlink, hoping for it to get versioned
	fc.deleteFile("foo")
	fc.sendIndexUpdate()
	for updates := 0; updates < 1; updates += <-idx {
	}

	// Recreate foo and a file in it with some data
	fc.addFile("foo", 0755, protocol.FileInfoTypeDirectory, nil)
	fc.addFile("foo/test", 0644, protocol.FileInfoTypeFile, []byte("testtesttest"))
	fc.sendIndexUpdate()
	for updates := 0; updates < 1; updates += <-idx {
	}

	// Remove the test file and see if it escaped
	fc.deleteFile("foo/test")
	fc.sendIndexUpdate()
	for updates := 0; updates < 1; updates += <-idx {
	}

	path := filepath.Join(tmpdir, "test")
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatal("File escaped to", path)
	}
}

func setupModelWithConnection() (*Model, *fakeConnection) {
	cfg := defaultConfig.RawCopy()
	cfg.Folders[0] = config.NewFolderConfiguration("default", "_tmpfolder")
		err = osutil.InWritableDir(func(name string) error {
			return f.moveForConflict(name, file.ModifiedBy.String())
		}, realName)
	} else if f.versioner != nil && !cur.IsSymlink() {
		err = osutil.InWritableDir(f.versioner.Archive, realName)
	} else {
		err = osutil.InWritableDir(os.Remove, realName)
				return err
			}

		case f.versioner != nil && !state.file.IsSymlink():
			// If we should use versioning, let the versioner archive the old
			// file before we replace it. Archiving a non-existent file is not
			// an error.
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

func NewSimple(folderID, folderPath string, params map[string]string) Versioner {
	cleanSymlinks(folderPath)

	keep, err := strconv.Atoi(params["keep"])
	if err != nil {
		keep = 5 // A reasonable default
	} else if err != nil {
		return err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		panic("bug: attempting to version a symlink")
	}

	versionsDir := filepath.Join(v.folderPath, ".stversions")
	_, err = os.Stat(versionsDir)
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
}

func NewTrashcan(folderID, folderPath string, params map[string]string) Versioner {
	cleanSymlinks(folderPath)

	cleanoutDays, _ := strconv.Atoi(params["cleanoutDays"])
	// On error we default to 0, "do not clean out the trash can"

// Archive moves the named file away to a version archive. If this function
// returns nil, the named file does not exist any more (has been archived).
func (t *Trashcan) Archive(filePath string) error {
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

	versionsDir := filepath.Join(t.folderPath, ".stversions")
	if _, err := os.Stat(versionsDir); err != nil {
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
