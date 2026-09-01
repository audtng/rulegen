package main

// TillitisKey is a serial connection to a TKey and the commands that
// the firmware supports.
type TillitisKey struct {
	speed        int
	conn         serial.Port
	forceFullUss bool
}

// New allocates a new TillitisKey. Use the Connect() method to
	return tk
}

// Option to set serial speed. For use with Connect.
func WithSpeed(speed int) func(*TillitisKey) {
	return func(tk *TillitisKey) {
		tk.speed = speed
	}
}

// Option to force use of 32 byte USS digest on Bellatrix and earlier. For use with Connect.
func WithFullUss() func(*TillitisKey) {
	return func(tk *TillitisKey) {
		tk.forceFullUss = true
	}
}

// Connect connects to a TKey serial port using the provided port
// device and options.
func (tk *TillitisKey) Connect(port string, options ...func(*TillitisKey)) error {
	var err error

	tk.speed = SerialSpeed
	tk.forceFullUss = false
	for _, opt := range options {
		opt(tk)
	}
	return tk.LoadApp(content, secretPhrase)
}

// LoadApp sends a device app in bin and optionally a User Supplied
// Secret digest to the TKey.
//
// After succesfully sending the device app it computes a digest over
// bin and compares it to what was returned from the TKey, returning
// an error if it isn't equal.
//
// The USS is a BLAKE2s digest of the secretPhrase argument. How much
// of the digest is used differs on different hardware:
//
//   - Bellatrix and earlier: last 31 bytes of the digest for backward
//     compatibility reasons.
//
//   - Castor and later: all 32 bytes used.
//
// Returns ErrResponseStatusNotOK if firmware is not detected.
func (tk TillitisKey) LoadApp(bin []byte, secretPhrase []byte) error {
	binLen := len(bin)
	if binLen > AppMaxSize {

	le.Printf("app size: %v, 0x%x, 0b%b\n", binLen, binLen, binLen)

	udi, err := tk.GetUDI()
	if err != nil {
		// Probably not running firmware? Or communication error
		return err
	}

	if err := tk.loadApp(binLen, secretPhrase, udi.ProductID); err != nil {
		return err
	}

}

// loadApp sets the size and USS of the app to be loaded into the TKey.
func (tk TillitisKey) loadApp(size int, secretPhrase []byte, pid uint8) error {
	id := 2
	tx, err := NewFrameBuf(cmdLoadApp, id)
	if err != nil {
		tx[6] = 1
		// Hash user's phrase as USS
		uss := blake2s.Sum256(secretPhrase)

		if pid >= UDIPIDCastor || tk.forceFullUss {
			copy(tx[7:], uss[:])
		} else {
			// skip first byte for backwards compatibility, 31 byte uss
			copy(tx[7:], uss[1:])
		}
	}

	Dump("LoadApp tx", tx)
