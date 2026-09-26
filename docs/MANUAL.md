# IntelliBranch: Comprehensive Manual & Tutorial for Go Developers

This guide provides pure Go engineers with a deep-dive technical manual and hands-on tutorial for **IntelliBranch**. It covers architectural workflows, exact keyword/API semantics, idiomatic usage patterns, and real-world production recipes.

---

## Table of Contents

1. [Architectural Mental Model for Go Engineers](#1-architectural-mental-model-for-go-engineers)
2. [The 4-Step Operational Workflow](#2-the-4-step-operational-workflow)
3. [Keyword & API Reference Manual](#3-keyword--api-reference-manual)
   - [Constructor: `NewRouter`](#31-newrouter)
   - [Branch Binding: `Bind`](#32-bind)
   - [Safety Guard: `Fallback`](#33-fallback)
   - [Execution: `Dispatch`](#34-dispatch)
   - [Observability: `Inspect`](#35-inspect)
   - [Diagnostic Struct: `RouteTrace`](#36-routetrace)
4. [End-to-End Production Tutorial](#4-end-to-end-production-tutorial)
   - [Step 1: Dataset Authoring](#step-1-authoring-the-domain-dataset)
   - [Step 2: Binary Weight Compilation](#step-2-compiling-binary-weights)
   - [Step 3: Building an HTTP Microservice Router](#step-3-building-an-http-microservice-router)
5. [Advanced Production Recipes](#5-advanced-production-recipes)
   - [Context Propagation & Timeouts](#51-context-propagation--timeouts)
   - [Atomic Zero-Downtime Weight Hot-Reloading](#52-atomic-zero-downtime-weight-hot-reloading)
   - [Whitebox Telemetry & Structured Logging](#53-whitebox-telemetry--structured-logging)

---

## 1. Architectural Mental Model for Go Engineers

In standard Go, control flow branching over strings relies on discrete equality:

```go
// Standard Go: Discrete String Equality
switch input {
case "refund":
    return processRefund()
}
```

This works if and only if `input` precisely equals `"refund"`. If the caller sends `"refnd"`, `"I need my money back"`, or `"reverse charge"`, the statement falls through.

**IntelliBranch** replaces discrete byte comparison with **continuous vector coordinate proximity**:

```text
[ Input Text ]
     │
     ▼
[ BPE Tokenizer ] ────── Maps characters to statistical subword chunks (robust to typos)
     │
     ▼
[ Dense Latent Space ] ── Words with identical intent share neighboring 64-D coordinates
     │
     ▼
[ Non-Linear Hyperplane ] ─ 128-D GELU separates opposing meanings (e.g. negation)
     │
     ▼
[ Softmax Distribution ] ── Probabilities summing to 1.0 (e.g. Refund: 0.98, Delivery: 0.01)
     │
     ▼
[ Branch Dispatch ] ───── Direct Go function execution in ~6.08 microseconds
```

---

## 2. The 4-Step Operational Workflow

```text
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│ 1. Data Design  │ ────▶ │ 2. Compilation  │ ────▶ │ 3. Wire Handlers│ ────▶ │ 4. Live Traffic │
│ (CSV Dataset)   │       │ (ib-train CLI)  │       │ (router.Bind)   │       │ (Dispatch / 6μs)│
└─────────────────┘       └─────────────────┘       └─────────────────┘       └─────────────────┘
```

1. **Data Design (`sample_dataset.csv`)**: Define target classes and author 20–50 representative real-world phrasing examples per class.
2. **Offline Compilation (`ib-train`)**: The CLI engine builds subword merges and trains neural weights into a compact Little-Endian binary (`.bin`) with SHA-256 integrity verification.
3. **Handler Binding (`main.go`)**: Initialize `Router`, bind target labels to standard Go functions, and set a fallback handler.
4. **Live In-Memory Dispatch (`router.Dispatch`)**: Incoming requests are evaluated and dispatched within 6 microseconds with single-digit memory allocations (`sync.Pool`).

---

## 3. Keyword & API Reference Manual

### 3.1. `NewRouter`

Initializes an in-memory `Router` from a pre-compiled binary weight file.

```go
func NewRouter(weightsPath string, threshold float32) (*Router, error)
```

- **Parameters**:
  - `weightsPath` (`string`): Absolute or relative filesystem path to the compiled `.bin` file.
  - `threshold` (`float32`): Minimum Softmax confidence score (recommended: `0.55` – `0.70`). Predictions with confidence below this threshold are rejected and routed to `Fallback`.
- **Guarantees**:
  - Validates the `IBRN` 4-byte magic header.
  - Validates format version compatibility (`0x0001`).
  - Verifies the SHA-256 binary integrity checksum against tampering.
  - Pre-allocates `sync.Pool` scratch buffers to ensure zero dynamic heap allocations on the hot path.

---

### 3.2. `Bind`

Registers a domain label to a target Go business handler.

```go
func (r *Router) Bind(label string, handler Handler) *Router
```

- **Handler Type**:
  ```go
  type Handler func(ctx context.Context, payload any) error
  ```
- **Semantics**:
  - Matches the exact class string defined in your CSV dataset (e.g. `"Refund"`, `"Delivery"`).
  - Supports fluent method chaining.
  - Thread-safe after startup.

---

### 3.3. `Fallback`

Designates the safety handler executed when statistical confidence is insufficient or when input is unclassifiable.

```go
func (r *Router) Fallback(handler Handler) *Router
```

- **Triggers**:
  1. Top Softmax confidence score is strictly less than `threshold`.
  2. Input consists predominantly of Out-of-Vocabulary (OOV) tokens (e.g., gibberish, random punctuation, binary noise).
  3. No specific `Bind` handler was registered for the predicted label.

---

### 3.4. `Dispatch`

Performs subword tokenization, forward neural inference, confidence evaluation, and executes the appropriate handler.

```go
func (r *Router) Dispatch(ctx context.Context, text string, payload any) error
```

- **Execution Latency**: ~6.08 microseconds (AMD Ryzen 5600H).
- **Concurrency**: Safe for concurrent execution across hundreds of goroutines without locks on the read path.
- **Safety Protections**:
  - Truncates raw text at 512 bytes to defend against memory amplification.
  - Clamps subword token sequences to 128 tokens to cap quadratic loop limits.
  - Propagates standard `context.Context` to allow cancellation and deadline enforcement.

---

### 3.5. `Inspect`

Inspects the decision process without invoking business handlers. Returns whitebox diagnostic metadata.

```go
func (r *Router) Inspect(text string) RouteTrace
```

- **Use Cases**:
  - Real-time logging of ambiguous requests.
  - Visualizing tokenization splits in developer staging environments.
  - Continuous evaluation of confidence distributions across incoming traffic.

---

### 3.6. `RouteTrace`

The structured diagnostic record returned by `Inspect`.

```go
type RouteTrace struct {
    InputText          string             `json:"input_text"`
    TokenIDs           []int              `json:"token_ids"`
    Subwords           []string           `json:"subwords"`
    UnknownTokenRatio  float32            `json:"unknown_token_ratio"`
    ClassProbabilities map[string]float32 `json:"class_probabilities"`
    PredictedLabel     string             `json:"predicted_label"`
    Confidence         float32            `json:"confidence"`
    Threshold          float32            `json:"threshold"`
    IsFallback         bool               `json:"is_fallback"`
    LatencyMicros      int64              `json:"latency_micros"`
}
```

---

## 4. End-to-End Production Tutorial

### Step 1: Authoring the Domain Dataset (`dataset.csv`)

The intelligence of the routing engine directly reflects the quality and variety of your dataset. Below are the mandatory structural specifications and data engineering principles:

#### 1. File Format & Schema Specifications

- **Header**: The first row must strictly be `text,label`.
- **Encoding**: UTF-8 without BOM.
- **Delimiter**: Comma (`,`). If an input text contains commas, wrap the text in standard double quotes (`"`):
  ```csv
  text,label
  "hey, where is my order?",Delivery
  ```
- **Label Consistency**: Labels are case-sensitive strings and must exactly match the string literals passed to `.Bind("Label", ...)` in your Go code.

#### 2. Golden Rules for High-Accuracy Datasets

| Rule | Specification | Engineering Rationale |
| :--- | :--- | :--- |
| **Minimum Sample Count** | **30 – 150 samples per class** | Guarantees enough subword co-occurrences for BPE and AdamW convergence. |
| **Class Balance** | Keep sample ratios within **1:1 to 2:1** | Prevents the model from biasing predictions toward over-represented classes. |
| **Linguistic Entropy** | Vary syntax, length, and vocabulary | Mix short queries (`"refund plz"`), full sentences, questions, and commands. |
| **Slang & Typos** | Deliberately include common mistakes | Expose the BPE tokenizer to misspellings (`"refnd"`, `"delivry"`, `"pasword"`). |
| **Boundary Disambiguation**| Include shared-word contrastive samples | Disambiguate `"cancel delivery alerts"` (Delivery) from `"cancel my charge"` (Refund). |
| **Noise Exclusion** | **Do NOT add random noise rows** | The engine's linear OOV penalty and `< 0.60` threshold automatically isolate noise. |

#### 3. Dataset Example: DOs vs. DONTs

```csv
text,label
# ✅ DO: Realistic phrasing, abbreviations, and sentence variety
can u cancel order #49281? i bought it by mistake,Refund
got charged twice on my card refund the extra charge asap,Refund
tracking says delivered but mailbox is empty where is my stuff,Delivery
sent back the return box 3 days ago when do i see money,Refund
locked out of my account after 3 failed tries,Account

# ❌ DONT: Robotic, repetitive keywords with zero variation
refund,Refund
refund please,Refund
refund now,Refund
delivery,Delivery
```

### Step 2: Compiling Binary Weights

Build the training tool and train the dataset:

```bash
# 1. Compile the training tool
go build -ldflags="-s -w" -o bin/ib-train.exe ./cmd/ib-train

# 2. Compile model weights into Little-Endian binary
./bin/ib-train.exe -data data/support_intents.csv -out weights/support.bin -epochs 50 -lr 0.005 -vocab 250
```

### Step 3: Building an HTTP Microservice Router

Create a high-performance HTTP service routing incoming support requests in microseconds:

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"intellibranch/pkg/intellibranch"
)

type RequestPayload struct {
	Query  string `json:"query"`
	UserID string `json:"user_id"`
}

type ResponsePayload struct {
	Status  string  `json:"status"`
	Action  string  `json:"action"`
	Score   float32 `json:"confidence"`
	Latency string  `json:"latency"`
}

func main() {
	// 1. Initialize Router with a 0.60 calibrated confidence threshold
	router, err := intellibranch.NewRouter("weights/support.bin", 0.60)
	if err != nil {
		log.Fatalf("Failed to initialize IntelliBranch: %v", err)
	}

	// 2. Bind business domain handlers
	router.
		Bind("Refund", func(ctx context.Context, payload any) error {
			req := payload.(*RequestPayload)
			log.Printf("[ACTION: Refund] Initiating refund process for user: %s", req.UserID)
			return nil
		}).
		Bind("Delivery", func(ctx context.Context, payload any) error {
			req := payload.(*RequestPayload)
			log.Printf("[ACTION: Delivery] Querying carrier API for user: %s", req.UserID)
			return nil
		}).
		Bind("Account", func(ctx context.Context, payload any) error {
			req := payload.(*RequestPayload)
			log.Printf("[ACTION: Account] Triggering MFA reset for user: %s", req.UserID)
			return nil
		}).
		Fallback(func(ctx context.Context, payload any) error {
			req := payload.(*RequestPayload)
			log.Printf("[FALLBACK: Triage] Diverting ambiguous query to human queue. User: %s, Text: %s", req.UserID, req.Query)
			return nil
		})

	// 3. Expose high-throughput HTTP handler
	http.HandleFunc("/api/v1/route", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var body RequestPayload
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 50*time.Millisecond)
		defer cancel()

		// Microsecond in-memory dispatch
		trace := router.Inspect(body.Query)
		_ = router.Dispatch(ctx, body.Query, &body)

		resp := ResponsePayload{
			Status:  "OK",
			Action:  trace.PredictedLabel,
			Score:   trace.Confidence,
			Latency: time.Since(start).String(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	log.Println("IntelliBranch HTTP Router listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 5. Advanced Production Recipes

### 5.1. Context Propagation & Timeouts

The `Handler` signature strictly adheres to Go's standard `context.Context`:

```go
router.Bind("Refund", func(ctx context.Context, payload any) error {
    select {
    case <-ctx.Done():
        return ctx.Err() // Gracefully abort if downstream client disconnected
    default:
        // Execute business logic
        return nil
    }
})
```

---

### 5.2. Atomic Zero-Downtime Weight Hot-Reloading

When updated weights are trained, replace the active router in production without dropping a single in-flight request:

```go
package main

import (
	"sync/atomic"

	"intellibranch/pkg/intellibranch"
)

type DynamicEngine struct {
	router atomic.Pointer[intellibranch.Router]
}

func (e *DynamicEngine) Reload(newWeightPath string) error {
	freshRouter, err := intellibranch.NewRouter(newWeightPath, 0.60)
	if err != nil {
		return err
	}
	// Atomic pointer swap: instantaneous, lock-free, zero downtime
	e.router.Store(freshRouter)
	return nil
}
```

---

### 5.3. Whitebox Telemetry & Structured Logging

Integrate `Inspect` directly with your logging pipeline (e.g. `uber-go/zap` or standard library `log/slog`):

```go
trace := router.Inspect(userMessage)

slog.Info("IntelliBranch dispatch complete",
    "query", trace.InputText,
    "selected_branch", trace.PredictedLabel,
    "confidence", trace.Confidence,
    "is_fallback", trace.IsFallback,
    "unknown_tokens", trace.UnknownTokenRatio,
    "duration_micros", trace.LatencyMicros,
)
```
