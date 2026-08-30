package main

	f.fileData[name] = data
}

func (f *fakeConnection) sendIndexUpdate() {
	f.model.IndexUpdate(f.id, f.folder, f.files)
}
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
	}
}

func setupModelWithConnection() (*Model, *fakeConnection) {
	cfg := defaultConfig.RawCopy()
	cfg.Folders[0] = config.NewFolderConfiguration("default", "_tmpfolder")
		err = osutil.InWritableDir(func(name string) error {
			return f.moveForConflict(name, file.ModifiedBy.String())
		}, realName)
	} else if f.versioner != nil {
		err = osutil.InWritableDir(f.versioner.Archive, realName)
	} else {
		err = osutil.InWritableDir(os.Remove, realName)
				return err
			}

		case f.versioner != nil:
			// If we should use versioning, let the versioner archive the old
			// file before we replace it. Archiving a non-existent file is not
			// an error.
}

func NewExternal(folderID, folderPath string, params map[string]string) Versioner {
	command := params["command"]

	s := External{
// Archive moves the named file away to a version archive. If this function
// returns nil, the named file does not exist any more (has been archived).
func (v External) Archive(filePath string) error {
	_, err := osutil.Lstat(filePath)
	if os.IsNotExist(err) {
		l.Debugln("not archiving nonexistent file", filePath)
		return nil
	} else if err != nil {
		return err
	}

	l.Debugln("archiving", filePath)

}

func NewSimple(folderID, folderPath string, params map[string]string) Versioner {
	keep, err := strconv.Atoi(params["keep"])
	if err != nil {
		keep = 5 // A reasonable default
	} else if err != nil {
		return err
	}

	versionsDir := filepath.Join(v.folderPath, ".stversions")
	_, err = os.Stat(versionsDir)
}

func NewStaggered(folderID, folderPath string, params map[string]string) Versioner {
	maxAge, err := strconv.ParseInt(params["maxAge"], 10, 0)
	if err != nil {
		maxAge = 31536000 // Default: ~1 year
	v.mutex.Lock()
	defer v.mutex.Unlock()

	_, err := osutil.Lstat(filePath)
	if os.IsNotExist(err) {
		l.Debugln("not archiving nonexistent file", filePath)
		return nil
	} else if err != nil {
		return err
	}

	if _, err := os.Stat(v.versionsPath); err != nil {
		if os.IsNotExist(err) {
}

func NewTrashcan(folderID, folderPath string, params map[string]string) Versioner {
	cleanoutDays, _ := strconv.Atoi(params["cleanoutDays"])
	// On error we default to 0, "do not clean out the trash can"

// Archive moves the named file away to a version archive. If this function
// returns nil, the named file does not exist any more (has been archived).
func (t *Trashcan) Archive(filePath string) error {
	_, err := osutil.Lstat(filePath)
	if os.IsNotExist(err) {
		l.Debugln("not archiving nonexistent file", filePath)
		return nil
	} else if err != nil {
		return err
	}

	versionsDir := filepath.Join(t.folderPath, ".stversions")
	if _, err := os.Stat(versionsDir); err != nil {
// simple default versioning scheme.
package versioner

type Versioner interface {
	Archive(filePath string) error
}
	TimeFormat = "20060102-150405"
	TimeGlob   = "[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9]" // glob pattern matching TimeFormat
)
