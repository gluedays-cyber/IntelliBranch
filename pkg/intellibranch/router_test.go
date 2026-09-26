package intellibranch

import (
	"context"
	"path/filepath"
	"testing"
)

func TestRouterDispatchAndFallback(t *testing.T) {
	tempDir := t.TempDir()
	modelPath := filepath.Join(tempDir, "router_test.bin")

	origModel := createSampleModel()
	if err := SaveBinaryModel(modelPath, origModel); err != nil {
		t.Fatalf("Failed to save model: %v", err)
	}

	router, err := NewRouter(modelPath, 0.1)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	var refundCalled, deliveryCalled, fallbackCalled bool

	router.
		Bind("Refund", func(ctx context.Context, payload any) error {
			refundCalled = true
			return nil
		}).
		Bind("Delivery", func(ctx context.Context, payload any) error {
			deliveryCalled = true
			return nil
		}).
		Fallback(func(ctx context.Context, payload any) error {
			fallbackCalled = true
			return nil
		})

	// Dispatch with matching label tokens
	err = router.Dispatch(context.Background(), "refund", nil)
	if err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}

	if !refundCalled && !fallbackCalled && !deliveryCalled {
		t.Fatal("No route action was executed")
	}

	// Dispatch with high threshold causing deterministic fallback
	strictRouter, err := NewRouter(modelPath, 0.999999)
	if err != nil {
		t.Fatalf("Failed to create strict router: %v", err)
	}

	fallbackTriggered := false
	strictRouter.
		Bind("Refund", func(ctx context.Context, payload any) error {
			return nil
		}).
		Fallback(func(ctx context.Context, payload any) error {
			fallbackTriggered = true
			return nil
		})

	if err := strictRouter.Dispatch(context.Background(), "refund", nil); err != nil {
		t.Fatalf("Strict dispatch failed: %v", err)
	}

	if !fallbackTriggered {
		t.Error("Expected fallback route to trigger when threshold exceeds confidence")
	}
}

func TestRouterInspectAndTrace(t *testing.T) {
	tempDir := t.TempDir()
	modelPath := filepath.Join(tempDir, "router_trace_test.bin")

	origModel := createSampleModel()
	if err := SaveBinaryModel(modelPath, origModel); err != nil {
		t.Fatalf("Failed to save model: %v", err)
	}

	router, err := NewRouter(modelPath, 0.5)
	if err != nil {
		t.Fatalf("Failed to create router: %v", err)
	}

	router.Bind("Refund", func(ctx context.Context, payload any) error {
		return nil
	})

	trace := router.Inspect("refund")
	if len(trace.TokenIDs) == 0 {
		t.Error("Expected non-empty token IDs in trace")
	}
	if len(trace.ClassProbabilities) == 0 {
		t.Error("Expected class probability breakdown in trace")
	}
	if trace.Confidence <= 0.0 {
		t.Errorf("Invalid confidence in trace: %f", trace.Confidence)
	}
}

