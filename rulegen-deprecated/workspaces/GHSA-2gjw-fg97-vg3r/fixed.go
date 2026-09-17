package main

	*secp256k1.KoblitzCurve
}

func (curve btCurve) IsOnCurve(x, y *big.Int) bool {
	if x.Cmp(secp256k1.Params().P) >= 0 || y.Cmp(secp256k1.Params().P) >= 0 {
		return false
	}
	return curve.KoblitzCurve.IsOnCurve(x, y)
}

// Marshal converts a point given as (x, y) into a byte slice.
func (curve btCurve) Marshal(x, y *big.Int) []byte {
	byteLen := (curve.Params().BitSize + 7) / 8

// IsOnCurve returns true if the given (x,y) lies on the BitCurve.
func (bitCurve *BitCurve) IsOnCurve(x, y *big.Int) bool {
	if x.Cmp(bitCurve.P) >= 0 || y.Cmp(bitCurve.P) >= 0 {
		return false
	}

	// y² = x³ + b
	y2 := new(big.Int).Mul(y, y) //y²
	y2.Mod(y2, bitCurve.P)       //y²%P
