package intellibranch

import (
	"errors"
	"math"
)

const (
	// Sqrt2OverPi is sqrt(2 / pi) for the GELU tanh approximation.
	Sqrt2OverPi float32 = 0.7978845608

	// GeluCoeff is the polynomial coefficient for GELU tanh approximation.
	GeluCoeff float32 = 0.044715
)

var (
	ErrZeroLengthTokens = errors.New("cannot pool over zero tokens")
	ErrDimensionMismatch = errors.New("tensor dimension mismatch during linear operation")
)

// GELU calculates the Gaussian Error Linear Unit activation using the standard tanh approximation.
// GELU(x) = 0.5 * x * (1 + tanh(sqrt(2 / pi) * (x + 0.044715 * x^3)))
func GELU(x float32) float32 {
	cube := x * x * x
	inner := Sqrt2OverPi * (x + GeluCoeff*cube)
	tanhVal := float32(math.Tanh(float64(inner)))
	return 0.5 * x * (1.0 + tanhVal)
}

// GELUInPlace applies the GELU non-linear activation function across a float32 slice in-place.
func GELUInPlace(vec []float32) {
	for i := 0; i < len(vec); i++ {
		x := vec[i]
		cube := x * x * x
		inner := Sqrt2OverPi * (x + GeluCoeff*cube)
		tanhVal := float32(math.Tanh(float64(inner)))
		vec[i] = 0.5 * x * (1.0 + tanhVal)
	}
}

// MeanPooling computes the average embedding vector across the given token IDs.
// out must have a length of at least embDim.
func MeanPooling(tokenIDs []uint32, embeddingTable []float32, embDim int, out []float32) error {
	seqLen := len(tokenIDs)
	if seqLen == 0 {
		return ErrZeroLengthTokens
	}

	// Zero out target buffer
	for i := 0; i < embDim; i++ {
		out[i] = 0.0
	}

	// Accumulate embeddings
	for _, id := range tokenIDs {
		offset := int(id) * embDim
		if offset+embDim > len(embeddingTable) {
			return errors.New("token ID exceeds embedding table bounds")
		}
		for d := 0; d < embDim; d++ {
			out[d] += embeddingTable[offset+d]
		}
	}

	// Scale by 1 / L
	invLen := 1.0 / float32(seqLen)
	for d := 0; d < embDim; d++ {
		out[d] *= invLen
	}

	return nil
}

// MatMulVecAdd computes out = vec * weights + bias where vec is [1 x inDim], weights is [inDim x outDim],
// bias is [outDim], and out is [outDim].
func MatMulVecAdd(vec []float32, weights []float32, bias []float32, inDim int, outDim int, out []float32) error {
	if len(vec) < inDim || len(bias) < outDim || len(out) < outDim || len(weights) < inDim*outDim {
		return ErrDimensionMismatch
	}

	// Initialize with bias values
	copy(out[:outDim], bias[:outDim])

	// Perform vector-matrix product with cache-efficient layout: weights is inDim x outDim row-major
	for i := 0; i < inDim; i++ {
		v := vec[i]
		if v == 0.0 {
			continue
		}
		rowOffset := i * outDim
		for j := 0; j < outDim; j++ {
			out[j] += v * weights[rowOffset+j]
		}
	}

	return nil
}

// Softmax computes the numerically stable softmax probabilities over logits with temperature scaling.
// out must have at least len(logits).
func Softmax(logits []float32, temperature float32, out []float32) error {
	n := len(logits)
	if n == 0 {
		return errors.New("empty logits")
	}
	if temperature <= 0.0 {
		temperature = 1.0
	}

	invTemp := 1.0 / temperature

	// Find max logit for numerical stability
	maxLogit := logits[0] * invTemp
	for i := 1; i < n; i++ {
		scaled := logits[i] * invTemp
		if scaled > maxLogit {
			maxLogit = scaled
		}
	}

	// Compute exp and sum
	var sumExp float32
	for i := 0; i < n; i++ {
		e := float32(math.Exp(float64(logits[i]*invTemp - maxLogit)))
		out[i] = e
		sumExp += e
	}

	// Normalize
	invSum := 1.0 / sumExp
	for i := 0; i < n; i++ {
		out[i] *= invSum
	}

	return nil
}
