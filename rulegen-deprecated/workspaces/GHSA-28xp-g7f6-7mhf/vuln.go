package main

// simple default versioning scheme.
package versioner

type Versioner interface {
	Archive(filePath string) error
}
	TimeFormat = "20060102-150405"
	TimeGlob   = "[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9]" // glob pattern matching TimeFormat
)
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
