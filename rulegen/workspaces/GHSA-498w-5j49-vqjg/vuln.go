package main

	"github.com/consensys/gnark/debug"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/internal/expr"
	"github.com/consensys/gnark/internal/utils"
	"github.com/consensys/gnark/std/math/bits"
)


// AssertIsLessOrEqual fails if  v > bound
func (builder *builder) AssertIsLessOrEqual(v frontend.Variable, bound frontend.Variable) {
	switch b := bound.(type) {
	case expr.Term:
		builder.mustBeLessOrEqVar(v, b)
	default:
		builder.mustBeLessOrEqCst(v, utils.FromInterface(b))
	}
}


	nbBits := builder.cs.FieldBitLen()

	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())
	boundBits := builder.ToBinary(bound, nbBits)

	p := make([]frontend.Variable, nbBits+1)
	p[nbBits] = 1

}

func (builder *builder) mustBeLessOrEqCst(a frontend.Variable, bound big.Int) {

	nbBits := builder.cs.FieldBitLen()

	// ensure the bound is positive, it's bit-len doesn't matter
	if bound.Sign() == -1 {
		panic("AssertIsLessOrEqual: bound is too large, constraint will never be satisfied")
	}

	if ca, ok := builder.constantValue(a); ok {
		// a is constant, compare the big int values
		ba := builder.cs.ToBigInt(ca)
		if ba.Cmp(&bound) == 1 {
			panic(fmt.Sprintf("AssertIsLessOrEqual: %s > %s", ba.String(), bound.String()))
		}
	}

	// debug info
	debug := builder.newDebugInfo("mustBeLessOrEq", a, " <= ", bound)

	// note that at this stage, we didn't boolean-constraint these new variables yet
	// (as opposed to ToBinary)
	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())

	// t trailing bits in the bound
	t := 0
	NbDigits             int
	UnconstrainedOutputs bool
	UnconstrainedInputs  bool
}

// BaseConversionOption configures the behaviour of scalar decomposition.
type BaseConversionOption func(opt *baseConversionConfig) error

// WithNbDigits sets the resulting number of digits (nbDigits) to be used in the base conversion.
// nbDigits must be > 0. If nbDigits is lower than the length of full decomposition and
// WithUnconstrainedOutputs option is not used, then the conversion functions will generate an
// unsatisfiable constraint. If WithNbDigits option is not set, then the full decomposition is
// returned.
func WithNbDigits(nbDigits int) BaseConversionOption {
	return func(opt *baseConversionConfig) error {
		if nbDigits <= 0 {
		return nil
	}
}
		}
	}

	// bound is constant
	if bConst {
		vv := builder.toVariable(v)
		builder.mustBeLessOrEqCst(vv, builder.cs.ToBigInt(cb))
		return
	}


	nbBits := builder.cs.FieldBitLen()

	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())
	boundBits := builder.ToBinary(bound, nbBits)

	// constraint added
	added := make([]int, 0, nbBits)

}

func (builder *builder) mustBeLessOrEqCst(a expr.LinearExpression, bound *big.Int) {

	nbBits := builder.cs.FieldBitLen()

	// ensure the bound is positive, it's bit-len doesn't matter
	if bound.Sign() == -1 {
	}

	// debug info
	debug := builder.newDebugInfo("mustBeLessOrEq", a, " <= ", builder.toVariable(bound))

	// note that at this stage, we didn't boolean-constraint these new variables yet
	// (as opposed to ToBinary)
	aBits := bits.ToBinary(builder, a, bits.WithNbDigits(nbBits), bits.WithUnconstrainedOutputs())

	// t trailing bits in the bound
	t := 0
