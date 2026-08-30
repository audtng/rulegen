package main

func (f *Fraction) Normalize(numerator, denominator int32) {
	for numerator > MAX_FRACTION_VALUE || numerator < -MAX_FRACTION_VALUE {
		numerator /= 2
		denominator /= 2
	}
}
