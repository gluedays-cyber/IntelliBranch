package intellibranch

import (
	"errors"
	"fmt"
	"sync"
	"unicode/utf8"
)

const (
	// MaxInputBytes limits raw input byte length to defend against algorithmic CPU exhaustion.
	MaxInputBytes = 512

	// MaxSequenceTokens limits token length to preserve sub-millisecond forward pass SLAs.
	MaxSequenceTokens = 128
)

var (
	ErrModelNotInitialized = errors.New("model not properly initialized")
	ErrClassIndexOutOfRange = errors.New("predicted class index exceeds label count")
)

// inferenceBuffer holds scratch memory slices to enable zero-allocation forward passes.
type inferenceBuffer struct {
	pooled []float32
	hidden []float32
	logits []float32
	probs  []float32
}

// InferenceModel represents an in-memory embedded classifier loaded from binary format.
type InferenceModel struct {
	Header      Header
	Labels      []string
	Vocab       []string
	MergeRules  []MergeRule
	Weights     Weights
	Temperature float32
	Tokenizer   *BPETokenizer

	bufPool sync.Pool
}

// NewInferenceModel constructs and prepares an InferenceModel with an internal scratch buffer pool.
func NewInferenceModel(
	header Header,
	labels []string,
	vocab []string,
	mergeRules []MergeRule,
	weights Weights,
) *InferenceModel {
	tok := NewBPETokenizer(vocab, mergeRules)

	model := &InferenceModel{
		Header:      header,
		Labels:      labels,
		Vocab:       vocab,
		MergeRules:  mergeRules,
		Weights:     weights,
		Temperature: 1.0,
		Tokenizer:   tok,
	}

	model.bufPool = sync.Pool{
		New: func() any {
			return &inferenceBuffer{
				pooled: make([]float32, header.EmbeddingDim),
				hidden: make([]float32, header.HiddenDim),
				logits: make([]float32, header.NumClasses),
				probs:  make([]float32, header.NumClasses),
			}
		},
	}

	return model
}

// Forward executes the 2-layer MLP inference over a slice of token IDs.
// It returns a newly allocated slice of class probabilities.
func (m *InferenceModel) Forward(tokenIDs []uint32, temperature float32) ([]float32, error) {
	if len(tokenIDs) == 0 {
		return nil, ErrEmptyInput
	}

	buf := m.bufPool.Get().(*inferenceBuffer)
	defer m.bufPool.Put(buf)

	// 1. Mean Pooling: [SeqLen] -> [EmbeddingDim]
	if err := MeanPooling(tokenIDs, m.Weights.Embedding, int(m.Header.EmbeddingDim), buf.pooled); err != nil {
		return nil, fmt.Errorf("mean pooling failed: %w", err)
	}

	// 2. Layer 1 Linear: [EmbeddingDim] x [EmbeddingDim x HiddenDim] + [HiddenDim] -> [HiddenDim]
	if err := MatMulVecAdd(buf.pooled, m.Weights.W1, m.Weights.B1, int(m.Header.EmbeddingDim), int(m.Header.HiddenDim), buf.hidden); err != nil {
		return nil, fmt.Errorf("layer 1 forward failed: %w", err)
	}

	// 3. GELU Non-Linear Activation In-Place
	GELUInPlace(buf.hidden)

	// 4. Layer 2 Linear: [HiddenDim] x [HiddenDim x NumClasses] + [NumClasses] -> [NumClasses]
	if err := MatMulVecAdd(buf.hidden, m.Weights.W2, m.Weights.B2, int(m.Header.HiddenDim), int(m.Header.NumClasses), buf.logits); err != nil {
		return nil, fmt.Errorf("layer 2 forward failed: %w", err)
	}

	// 5. Softmax with Temperature Scaling
	temp := temperature
	if temp <= 0.0 {
		temp = m.Temperature
	}
	if err := Softmax(buf.logits, temp, buf.probs); err != nil {
		return nil, fmt.Errorf("softmax failed: %w", err)
	}

	result := make([]float32, m.Header.NumClasses)
	copy(result, buf.probs)
	return result, nil
}

// PredictTokens computes class probabilities and returns the top label alongside its confidence score.
func (m *InferenceModel) PredictTokens(tokenIDs []uint32) (string, float64, error) {
	probs, err := m.Forward(tokenIDs, m.Temperature)
	if err != nil {
		return "", 0.0, err
	}

	var bestIdx int
	var maxProb float32 = -1.0

	for i, p := range probs {
		if p > maxProb {
			maxProb = p
			bestIdx = i
		}
	}

	if bestIdx < 0 || bestIdx >= len(m.Labels) {
		return "", 0.0, ErrClassIndexOutOfRange
	}

	return m.Labels[bestIdx], float64(maxProb), nil
}

// Predict tokenizes raw text with subword BPE and returns predicted label and confidence score with safety guards.
func (m *InferenceModel) Predict(text string) (string, float64, error) {
	// Guard 1: Validate UTF-8 and truncate oversized input strings
	if !utf8.ValidString(text) {
		return "", 0.0, ErrEmptyInput
	}
	if len(text) > MaxInputBytes {
		text = text[:MaxInputBytes]
	}

	tokenIDs := m.Tokenizer.Encode(text)
	if len(tokenIDs) == 0 {
		if m.Header.VocabSize > 0 {
			tokenIDs = []uint32{0}
		} else {
			return "", 0.0, ErrEmptyInput
		}
	}

	// Guard 2: Clamp sequence length to prevent excessive pooling latency
	if len(tokenIDs) > MaxSequenceTokens {
		tokenIDs = tokenIDs[:MaxSequenceTokens]
	}

	label, score, err := m.PredictTokens(tokenIDs)
	if err != nil {
		return "", 0.0, err
	}

	// Guard 3: Penalize confidence proportionally if input is dominated by out-of-vocabulary [UNK] tokens
	unkID, hasUnk := m.Tokenizer.VocabMap["[UNK]"]
	if hasUnk && len(tokenIDs) > 0 {
		unkCount := 0
		for _, id := range tokenIDs {
			if id == unkID {
				unkCount++
			}
		}
		unkRatio := float64(unkCount) / float64(len(tokenIDs))
		score = score * (1.0 - unkRatio)
	}

	return label, score, nil
}


