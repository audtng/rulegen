package main

// TillitisKey is a serial connection to a TKey and the commands that
// the firmware supports.
type TillitisKey struct {
	speed int
	conn  serial.Port
}

// New allocates a new TillitisKey. Use the Connect() method to
	return tk
}

func WithSpeed(speed int) func(*TillitisKey) {
	return func(tk *TillitisKey) {
		tk.speed = speed
	}
}

// Connect connects to a TKey serial port using the provided port
// device and options.
func (tk *TillitisKey) Connect(port string, options ...func(*TillitisKey)) error {
	var err error

	tk.speed = SerialSpeed
	for _, opt := range options {
		opt(tk)
	}
	return tk.LoadApp(content, secretPhrase)
}

// LoadApp loads the USS (User Supplied Secret), and contents of bin
// into the TKey, running the app after verifying that the digest
// calculated on the host is the same as the digest from the TKey.
//
// The USS is a 32 bytes digest hashed from secretPhrase (which is
// provided by the user). If secretPhrase is an empty slice, 32 bytes
// of zeroes will be loaded as USS.
//
// Loading USS is always done together with loading and running an
// app, because the host program can't otherwise be sure that the
// expected USS is used.
func (tk TillitisKey) LoadApp(bin []byte, secretPhrase []byte) error {
	binLen := len(bin)
	if binLen > AppMaxSize {

	le.Printf("app size: %v, 0x%x, 0b%b\n", binLen, binLen, binLen)

	err := tk.loadApp(binLen, secretPhrase)
	if err != nil {
		return err
	}

}

// loadApp sets the size and USS of the app to be loaded into the TKey.
func (tk TillitisKey) loadApp(size int, secretPhrase []byte) error {
	id := 2
	tx, err := NewFrameBuf(cmdLoadApp, id)
	if err != nil {
		tx[6] = 1
		// Hash user's phrase as USS
		uss := blake2s.Sum256(secretPhrase)
		copy(tx[6:], uss[:])
	}

	Dump("LoadApp tx", tx)
