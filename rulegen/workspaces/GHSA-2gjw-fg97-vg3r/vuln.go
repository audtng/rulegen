package main

	*secp256k1.KoblitzCurve
}

// Marshal converts a point given as (x, y) into a byte slice.
func (curve btCurve) Marshal(x, y *big.Int) []byte {
	byteLen := (curve.Params().BitSize + 7) / 8

// IsOnCurve returns true if the given (x,y) lies on the BitCurve.
func (bitCurve *BitCurve) IsOnCurve(x, y *big.Int) bool {
	// y² = x³ + b
	y2 := new(big.Int).Mul(y, y) //y²
	y2.Mod(y2, bitCurve.P)       //y²%P
