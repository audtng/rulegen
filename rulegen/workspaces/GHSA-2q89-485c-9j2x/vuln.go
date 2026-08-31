package main


func EncryptCCA(rand io.Reader, public *PublicParams, policy *Policy, msg []byte) ([]byte, error) {
	seed := make([]byte, macKeySeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, err
	}
	}

	salt := make([]byte, v.hash.Size())
	_, err := random.Read(salt)
	if err != nil {
		return nil, nil, err
	}
func (pk *PublicKey) EncapsulateTo(ct []byte, ss []byte, seed []byte) {
	if seed == nil {
		seed = make([]byte, EncapsulationSeedSize)
		_, _ = cryptoRand.Read(seed[:])
	}
	if len(seed) != EncapsulationSeedSize {
		panic("seed must be of length EncapsulationSeedSize")
func (pk *PublicKey) EncapsulateTo(ct, ss []byte, seed []byte) {
	if seed == nil {
		seed = make([]byte, EncapsulationSeedSize)
		cryptoRand.Read(seed[:])
	} else {
		if len(seed) != EncapsulationSeedSize {
			panic("seed must be of length EncapsulationSeedSize")
func (pk *PublicKey) EncapsulateTo(ct, ss []byte, seed []byte) {
	if seed == nil {
		seed = make([]byte, EncapsulationSeedSize)
		cryptoRand.Read(seed[:])
	} else {
		if len(seed) != EncapsulationSeedSize {
			panic("seed must be of length EncapsulationSeedSize")
func (pk *PublicKey) EncapsulateTo(ct, ss []byte, seed []byte) {
	if seed == nil {
		seed = make([]byte, EncapsulationSeedSize)
		cryptoRand.Read(seed[:])
	} else {
		if len(seed) != EncapsulationSeedSize {
			panic("seed must be of length EncapsulationSeedSize")
func (pk *PublicKey) EncapsulateTo(ct, ss []byte, seed []byte) {
	if seed == nil {
		seed = make([]byte, EncapsulationSeedSize)
		cryptoRand.Read(seed[:])
	} else {
		if len(seed) != EncapsulationSeedSize {
			panic("seed must be of length EncapsulationSeedSize")

func (sch *scheme) Encapsulate(pk kem.PublicKey) (ct []byte, ss []byte, err error) {
	var seed [EncapsulationSeedSize]byte
	cryptoRand.Read(seed[:])
	return sch.EncapsulateDeterministically(pk, seed[:])
}


func (sch *scheme) Encapsulate(pk kem.PublicKey) (ct []byte, ss []byte, err error) {
	var seed [EncapsulationSeedSize]byte
	cryptoRand.Read(seed[:])
	return sch.EncapsulateDeterministically(pk, seed[:])
}


func (sch *scheme) Encapsulate(pk kem.PublicKey) (ct []byte, ss []byte, err error) {
	var seed [EncapsulationSeedSize]byte
	cryptoRand.Read(seed[:])
	return sch.EncapsulateDeterministically(pk, seed[:])
}


func (sch *scheme) Encapsulate(pk kem.PublicKey) (ct []byte, ss []byte, err error) {
	var seed [EncapsulationSeedSize]byte
	cryptoRand.Read(seed[:])
	return sch.EncapsulateDeterministically(pk, seed[:])
}

