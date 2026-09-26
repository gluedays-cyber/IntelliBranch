package intellibranch

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RouteAction defines the execution handler signature for a matched branch.
type RouteAction func(ctx context.Context, payload any) error

// RouteTrace encapsulates comprehensive diagnostic metadata explaining a routing decision.
type RouteTrace struct {
	InputText          string             `json:"input_text"`
	TokenIDs           []uint32           `json:"token_ids"`
	Subwords           []string           `json:"subwords"`
	UnknownTokenRatio  float64            `json:"unknown_token_ratio"`
	ClassProbabilities map[string]float32 `json:"class_probabilities"`
	PredictedLabel     string             `json:"predicted_label"`
	Confidence         float64            `json:"confidence"`
	Threshold          float64            `json:"threshold"`
	IsFallback         bool               `json:"is_fallback"`
	FallbackReason     string             `json:"fallback_reason,omitempty"`
	LatencyMicros      int64              `json:"latency_micros"`
}

// Router coordinates in-memory inference routing with deterministic fallback safety.
type Router struct {
	mu        sync.RWMutex
	model     *InferenceModel
	threshold float64
	routes    map[string]RouteAction
	fallback  RouteAction
}

// NewRouter loads a binary model file into memory once and constructs an immutable routing core.
func NewRouter(modelPath string, defaultThreshold float64) (*Router, error) {
	model, err := LoadBinaryModel(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize inference model: %w", err)
	}

	return &Router{
		model:     model,
		threshold: defaultThreshold,
		routes:    make(map[string]RouteAction),
		fallback: func(ctx context.Context, payload any) error {
			return nil
		},
	}, nil
}

// Bind registers an action handler for a target class label.
func (r *Router) Bind(label string, action RouteAction) *Router {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[label] = action
	return r
}

// Fallback registers the default handler triggered when confidence is below threshold or label is unmatched.
func (r *Router) Fallback(action RouteAction) *Router {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = action
	return r
}

// Dispatch executes microsecond inference and triggers the bound handler.
func (r *Router) Dispatch(ctx context.Context, text string, payload any) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	label, score, err := r.model.Predict(text)
	if err != nil || score < r.threshold {
		return r.fallback(ctx, payload)
	}

	action, exists := r.routes[label]
	if !exists {
		return r.fallback(ctx, payload)
	}

	return action(ctx, payload)
}

// Inspect evaluates input text and generates a full diagnostic RouteTrace without executing handlers.
func (r *Router) Inspect(text string) RouteTrace {
	start := time.Now()
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Guard 1: Truncate oversized input strings
	if len(text) > MaxInputBytes {
		text = text[:MaxInputBytes]
	}

	tokens := r.model.Tokenizer.Encode(text)
	if len(tokens) == 0 && r.model.Header.VocabSize > 0 {
		tokens = []uint32{0}
	}

	// Guard 2: Clamp sequence length
	if len(tokens) > MaxSequenceTokens {
		tokens = tokens[:MaxSequenceTokens]
	}

	subwords := make([]string, len(tokens))
	unkCount := 0
	unkID, hasUnk := r.model.Tokenizer.VocabMap["[UNK]"]
	for i, id := range tokens {
		if hasUnk && id == unkID {
			unkCount++
		}
		if int(id) < len(r.model.Vocab) {
			subwords[i] = r.model.Vocab[id]
		}
	}

	var unkRatio float64
	if hasUnk && len(tokens) > 0 {
		unkRatio = float64(unkCount) / float64(len(tokens))
	}

	probs, err := r.model.Forward(tokens, r.model.Temperature)
	probMap := make(map[string]float32, len(r.model.Labels))

	var bestLabel string
	var bestScore float32 = -1.0
	for i, p := range probs {
		if i < len(r.model.Labels) {
			lbl := r.model.Labels[i]
			probMap[lbl] = p
			if p > bestScore {
				bestScore = p
				bestLabel = lbl
			}
		}
	}

	// Guard 3: Apply calibrated confidence penalty
	calibratedConfidence := float64(bestScore) * (1.0 - unkRatio)

	trace := RouteTrace{
		InputText:          text,
		TokenIDs:           tokens,
		Subwords:           subwords,
		UnknownTokenRatio:  unkRatio,
		ClassProbabilities: probMap,
		PredictedLabel:     bestLabel,
		Confidence:         calibratedConfidence,
		Threshold:          r.threshold,
		LatencyMicros:      time.Since(start).Microseconds(),
	}

	if err != nil {
		trace.IsFallback = true
		trace.FallbackReason = fmt.Sprintf("inference error: %v", err)
	} else if trace.Confidence < r.threshold {
		trace.IsFallback = true
		trace.FallbackReason = fmt.Sprintf("confidence %.4f below threshold %.4f", trace.Confidence, r.threshold)
	} else if _, exists := r.routes[bestLabel]; !exists {
		trace.IsFallback = true
		trace.FallbackReason = fmt.Sprintf("label '%s' has no bound route handler", bestLabel)
	}

	return trace
}

// DispatchWithTrace runs inference, collects full diagnostic telemetry, and dispatches to handler.
func (r *Router) DispatchWithTrace(ctx context.Context, text string, payload any) (RouteTrace, error) {
	trace := r.Inspect(text)

	r.mu.RLock()
	defer r.mu.RUnlock()

	if trace.IsFallback {
		return trace, r.fallback(ctx, payload)
	}

	action := r.routes[trace.PredictedLabel]
	return trace, action(ctx, payload)
}
