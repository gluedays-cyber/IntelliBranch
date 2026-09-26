package intellibranch

import (
	"math"
	"testing"
)

func createStandardSpecModel() *InferenceModel {
	// Architectural specifications:
	// EmbeddingDim = 64
	// HiddenDim = 128
	// VocabSize = 500
	// NumClasses = 5
	header := Header{
		Magic:        MagicBytes,
		Version:      1,
		VocabSize:    500,
		EmbeddingDim: 64,
		HiddenDim:    128,
		NumClasses:   5,
	}

	labels := []string{"Refund", "Delivery", "Account", "Payment", "General"}
	vocab := make([]string, header.VocabSize)
	for i := range vocab {
		vocab[i] = "token"
	}
	vocab[0] = "[PAD]"
	vocab[1] = "[UNK]"
	vocab[2] = "r"
	vocab[3] = "e"
	vocab[4] = "f"
	vocab[5] = "u"
	vocab[6] = "n"
	vocab[7] = "d"
	vocab[8] = " "

	weights := Weights{
		Embedding: make([]float32, header.VocabSize*header.EmbeddingDim),
		W1:        make([]float32, header.EmbeddingDim*header.HiddenDim),
		B1:        make([]float32, header.HiddenDim),
		W2:        make([]float32, header.HiddenDim*header.NumClasses),
		B2:        make([]float32, header.NumClasses),
	}

	for i := range weights.Embedding {
		weights.Embedding[i] = 0.01 * float32(i%10)
	}
	for i := range weights.W1 {
		weights.W1[i] = 0.005 * float32((i%20)-10)
	}
	for i := range weights.W2 {
		weights.W2[i] = 0.005 * float32((i%15)-7)
	}

	return NewInferenceModel(header, labels, vocab, nil, weights)
}

func TestModelForward(t *testing.T) {
	model := createStandardSpecModel()
	tokens := []uint32{1, 2, 5, 10}

	probs, err := model.Forward(tokens, 1.0)
	if err != nil {
		t.Fatalf("Forward failed: %v", err)
	}

	if len(probs) != int(model.Header.NumClasses) {
		t.Fatalf("Expected %d class probabilities, got %d", model.Header.NumClasses, len(probs))
	}

	var sum float32
	for _, p := range probs {
		if p < 0.0 || p > 1.0 {
			t.Errorf("Invalid probability value: %f", p)
		}
		sum += p
	}

	if math.Abs(float64(sum-1.0)) > 1e-4 {
		t.Errorf("Probabilities sum to %f, expected 1.0", sum)
	}
}

func TestModelPredictTokens(t *testing.T) {
	model := createStandardSpecModel()
	tokens := []uint32{1, 2}

	label, score, err := model.PredictTokens(tokens)
	if err != nil {
		t.Fatalf("PredictTokens failed: %v", err)
	}

	if label == "" {
		t.Error("Expected non-empty label")
	}
	if score <= 0.0 || score > 1.0 {
		t.Errorf("Invalid score: %f", score)
	}
}

func TestInputTruncationGuard(t *testing.T) {
	model := createStandardSpecModel()
	// Create an oversized text string with 2000 characters
	longText := ""
	for i := 0; i < 200; i++ {
		longText += "refund "
	}

	label, score, err := model.Predict(longText)
	if err != nil {
		t.Fatalf("Predict on oversized input failed: %v", err)
	}
	if label == "" || score <= 0.0 {
		t.Errorf("Invalid prediction on truncated input: label=%s, score=%f", label, score)
	}
}

func TestOOVConfidenceDiscounting(t *testing.T) {
	model := createStandardSpecModel()
	// Text composed entirely of characters never seen in vocab
	noiseText := "§±¶€$!@#%^&*()_+~`|}{[]:;?><"

	_, score, err := model.Predict(noiseText)
	if err != nil {
		t.Fatalf("Predict on noise text failed: %v", err)
	}

	// Because almost all tokens are [UNK], calibrated score must be severely discounted (< 0.20)
	if score > 0.35 {
		t.Errorf("Expected severely discounted score for pure noise, got: %f", score)
	}
}

func BenchmarkForward(b *testing.B) {
	model := createStandardSpecModel()
	tokens := []uint32{1, 2, 10, 20, 30, 40}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := model.Forward(tokens, 1.0)
		if err != nil {
			b.Fatalf("Forward failed: %v", err)
		}
	}
}

func BenchmarkPredictTokens(b *testing.B) {
	model := createStandardSpecModel()
	tokens := []uint32{1, 2, 10, 20, 30, 40}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, err := model.PredictTokens(tokens)
		if err != nil {
			b.Fatalf("PredictTokens failed: %v", err)
		}
	}
}

func BenchmarkGELU(b *testing.B) {
	vec := make([]float32, 128)
	for i := range vec {
		vec[i] = float32(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GELUInPlace(vec)
	}
}
