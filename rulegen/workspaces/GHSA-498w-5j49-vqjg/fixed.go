package main

	"github.com/consensys/gnark/debug"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/internal/expr"
	"github.com/consensys/gnark/std/math/bits"
)


// AssertIsLessOrEqual fails if  v > bound
func (builder *builder) AssertIsLessOrEqual(v frontend.Variable, bound frontend.Variable) {
	cv, vConst := builder.constantValue(v)
	cb, bConst := builder.constantValue(bound)

	// both inputs are constants
	if vConst && bConst {
		bv, bb := builder.cs.ToBigInt(cv), builder.cs.ToBigInt(cb)
		if bv.Cmp(bb) == 1 {
			panic(fmt.Sprintf("AssertIsLessOrEqual: %s > %s", bv.String(), bb.String()))
		}
	}

	nbBits := builder.cs.FieldBitLen()
	vBits := bits.ToBinary(builder, v, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())

	// bound is constant
	if bConst {
		builder.MustBeLessOrEqCst(vBits, builder.cs.ToBigInt(cb), v)
		return
	}

	if b, ok := bound.(expr.Term); ok {
		builder.mustBeLessOrEqVar(v, b)
	} else {
		panic(fmt.Sprintf("expected bound type expr.Term, got %T", bound))
	}
}


	nbBits := builder.cs.FieldBitLen()

	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs(), bits.OmitModulusCheck())
	boundBits := bits.ToBinary(builder, bound, bits.WithNbDigits(nbBits)) // enforces range check against modulus

	p := make([]frontend.Variable, nbBits+1)
	p[nbBits] = 1

}

// MustBeLessOrEqCst asserts that value represented using its bit decomposition
// aBits is less or equal than constant bound. The method boolean constraints
// the bits in aBits, so the caller can provide unconstrained bits.
func (builder *builder) MustBeLessOrEqCst(aBits []frontend.Variable, bound *big.Int, aForDebug frontend.Variable) {

	nbBits := builder.cs.FieldBitLen()
	if len(aBits) > nbBits {
		panic("more input bits than field bit length")
	}
	for i := len(aBits); i < nbBits; i++ {
		aBits = append(aBits, 0)
	}

	// ensure the bound is positive, it's bit-len doesn't matter
	if bound.Sign() == -1 {
		panic("AssertIsLessOrEqual: bound is too large, constraint will never be satisfied")
	}

	// debug info
	debug := builder.newDebugInfo("mustBeLessOrEq", aForDebug, " <= ", bound)

	// t trailing bits in the bound
	t := 0
	NbDigits             int
	UnconstrainedOutputs bool
	UnconstrainedInputs  bool

	omitModulusCheck bool
}

// BaseConversionOption configures the behaviour of scalar decomposition.
type BaseConversionOption func(opt *baseConversionConfig) error

// WithNbDigits sets the resulting number of digits (nbDigits) to be used in the
// base conversion.
//
// nbDigits must be > 0. If nbDigits is lower than the length of full
// decomposition and [WithUnconstrainedOutputs] option is not used, then the
// conversion functions will generate an unsatisfiable constraint.
//
// If nbDigits is larger than the bitlength of the modulus, then the returned
// slice has length nbDigits with excess bits being 0.
//
// If WithNbDigits option is not set, then the full decomposition is returned.
func WithNbDigits(nbDigits int) BaseConversionOption {
	return func(opt *baseConversionConfig) error {
		if nbDigits <= 0 {
		return nil
	}
}

// OmitModulusCheck omits the comparison against native field modulus in
// case the bitlength of the decomposed value (if [WithNbDigits] not set or set
// to bitlength of the native modulus) eqals bitlength of the modulus.
//
// The check is otherwise required as there are possibly multiple correct binary
// decompositions. For example, when decomposing small a the decomposition could
// return the slices for both a or a+r, where r is the native modulus and the
// enforced constraints are correct due to implicit modular reduction by r.
//
// This option can be used in case the decomposed output is manually checked to
// be unique or if uniqueness is not required.
func OmitModulusCheck() BaseConversionOption {
	return func(opt *baseConversionConfig) error {
		opt.omitModulusCheck = true
		return nil
	}
}
		}
	}

	nbBits := builder.cs.FieldBitLen()
	vBits := bits.ToBinary(builder, v, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())

	// bound is constant
	if bConst {
		builder.MustBeLessOrEqCst(vBits, builder.cs.ToBigInt(cb), v)
		return
	}


	nbBits := builder.cs.FieldBitLen()

	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs(), bits.OmitModulusCheck())
	boundBits := bits.ToBinary(builder, bound, bits.WithNbDigits(nbBits))

	// constraint added
	added := make([]int, 0, nbBits)

}

// MustBeLessOrEqCst asserts that value represented using its bit decomposition
// aBits is less or equal than constant bound. The method boolean constraints
// the bits in aBits, so the caller can provide unconstrained bits.
func (builder *builder) MustBeLessOrEqCst(aBits []frontend.Variable, bound *big.Int, aForDebug frontend.Variable) {

	nbBits := builder.cs.FieldBitLen()
	if len(aBits) > nbBits {
		panic("more input bits than field bit length")
	}
	for i := len(aBits); i < nbBits; i++ {
		aBits = append(aBits, 0)
	}

	// ensure the bound is positive, it's bit-len doesn't matter
	if bound.Sign() == -1 {
	}

	// debug info
	debug := builder.newDebugInfo("mustBeLessOrEq", aForDebug, " <= ", builder.toVariable(bound))

	// t trailing bits in the bound
	t := 0
