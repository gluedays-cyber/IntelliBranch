# IntelliBranch
<img src="https://github.com/user-attachments/assets/3413a486-d71c-4285-841d-76bbe74f830a" width="226" height="200" alt="Image" align="right" style="margin-left: 15px; margin: 10px;">
<p align="center">
  <strong>Directly Creates and Runs Its Own Neural AI in Pure Go</strong><br>
  <em>Stop borrowing third-party AIs. This engine creates its own domain artificial intelligence from scratch in 1.5 seconds, routing execution flow in ~30 μs with 0 B/op (Zero Allocations), Zero Downloads, and Zero CGO.</em>
</p>

<p align="center">
  <a href="#benchmarks"><img src="https://img.shields.io/badge/Latency-~30_μs-brightgreen.svg" alt="Latency"></a>
  <a href="#benchmarks"><img src="https://img.shields.io/badge/Allocs-0_B/op_(0_allocs)-blue.svg" alt="Allocations"></a>
  <img src="https://img.shields.io/badge/Wire_Format-v2_Positional-orange.svg" alt="Format v2">
  <img src="https://img.shields.io/badge/CGO-Zero_Disabled-success.svg" alt="CGO Zero">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-lightgrey.svg" alt="License">
</p>

<p align="center">
  <a href="docs/MANUAL.md"><strong>📖 Read the Full Developer Manual & Production Tutorial →</strong></a>
</p>

---

## What is IntelliBranch?

**IntelliBranch does NOT borrow, lease, or download external AI models. This engine directly creates and runs its own domain artificial intelligence from scratch.**

Instead of relying on brittle regex matching or calling bloated external LLMs, it **manufactures a domain-specific lightweight neural network directly from your dataset in under 2 seconds**. It maps typos, slang, inverted syntax, and colloquial phrasing into a continuous latent vector space—routing execution flow directly to your bound Go functions in **microseconds (~30 μs) with strictly 0 B/op heap allocation**.

```
Incoming Request ("bruh can u refund order #49281")
                     │
                     ▼
       [ In-Memory BPE Tokenizer ]
                     │
                     ▼
  [ Dense (D=64) + Positional (P=32x64) ]
                     │
                     ▼
  [ Non-Linear GELU Mean Pooling (D=64) ]
                     │
                     ▼
      [ Hidden Projection (D=128) ]
                     │
                     ▼
   [ Softmax + Shannon Entropy Calibrated Guard ]
                     │
     ┌───────────────┼───────────────┬────────────────┐
     ▼               ▼               ▼                ▼
(Score ≥ 0.75)  (Score ≥ 0.30)  (Margin < 0.15)  (Entropy > 2.0 / UNK ≥ 0.5)
[DEFINITE ROUTE] [PIPELINE]      [AMBIGUOUS]      [FALLBACK ISOLATION]
```

---

## Why IntelliBranch? (Beyond Retro Branching, Cloud LLMs, and Bloated Local Models)

Modern backends face an architectural dilemma when routing unstructured or noisy user requests:

```go
// ❌ RETRO BRANCHING: Brittle, explodes in complexity, collapses under real-world noise
if strings.Contains(input, "refund") || strings.Contains(input, "cancel") {
    // FAILS on: "sent the return box a week ago when do i get my money back"
    // FAILS on: "can u reverse the charge?" (typos, slang, synonyms)
    // MISROUTES on: "cancel shipment delay notifications" (word collision)
}

// ❌ CLOUD LLMs: Massive network latency, recurring per-token cost, third-party dependency
// Latency: 400ms – 2,500ms (Unusable in high-throughput microservices)
// Cost: $0.0015 – $0.03 per request (Bills explode under scale)
// Vulnerability: Outages, rate limits, JSON hallucination, network partitions

// ❌ LOCAL LLMs & SLMs (Ollama, llama.cpp, Mistral-7B, Phi-3): Severe host resource exhaustion
// Memory: Monopolizes 4.5 GB to 8.0 GB+ of RAM/VRAM just to pick a 4-byte enum
// CPU Starvation: Burns 100% CPU across multiple cores, starving companion microservices
// Deployment Complexity: Requires CGO, C++ shared libraries (libllama.so), or background daemons

// ✅ INTELLIBRANCH: Self-Generated Micro-AI (In-Memory Go Engine)
// Memory Footprint: Under 180 KB (25,000x smaller than quantized 7B models)
// Latency: ~30 μs with 0 B/op (0 allocs) and deterministic 3-tier fallback
// Deployment: 100% Pure Go with CGO_ENABLED=0 single static binary
```

### Architectural Comparison Matrix

| Capability | Retro Branching (`if` / Regex) | Cloud LLMs (OpenAI / Claude) | Local LLMs (Ollama / llama.cpp) | **IntelliBranch v2.0 (Embedded Engine)** |
| :--- | :--- | :--- | :--- | :--- |
| **Inference Latency** | < 1 μs | 300 ms – 2,500 ms (Network bound) | 30 ms – 300 ms (Compute bound) | **~30 μs (In-Memory)** |
| **Throughput (per core)** | > 500,000 req/sec | ~50 req/sec (Rate limited) | ~20–50 req/sec (CPU saturated) | **> 33,000 req/sec (Zero Alloc)** |
| **Runtime Allocation** | 0 B/op | High (HTTP payload) | High (CGO buffers) | **0 B/op (0 allocs/op)** |
| **System Memory (RAM)** | Negligible | External service | **4.5 GB – 8.0 GB+ (VRAM / RAM)** | **< 180 KB (Format v2)** |
| **Token Order Awareness** | Rigid regex position | ✅ Transformer Attention | ✅ Transformer Attention | ✅ **Learned Positional Embeddings** |
| **Hardware Reqs** | Standard CPU | External service | High-end GPU or 8+ Core CPU | **Runs on a $5 VPS (16MB container)** |
| **Operational Cost** | $0.00 | $0.0015+ per call | High hardware/electricity cost | **$0.00 (Self-contained)** |
| **Hot Weight Reload** | Binary recompile | API model string switch | Multi-second model reload | **Lock-free Atomic Hot-Swap (`0 ns` stop)** |
| **Active Learning Loop** | N/A | Manual logging | N/A | **Built-in Ring Buffer Telemetry** |
| **Deployment Complexity** | Single binary | API client | CGO / C++ runtime / Ollama daemon | **Pure Go (`CGO_ENABLED=0`)** |

---

## Key Highlights

- **Zero Downloads & On-The-Fly AI Creation**: You never download gigabytes of pre-trained weights from HuggingFace or lease external APIs. IntelliBranch forges a domain neural AI model directly from your CSV in under 2 seconds.
- **Zero Allocations on Hot Path (`0 B/op`)**: `PredictSlots` executes inference without triggering GC pressure, returning zero-heap stack results.
- **Semantic XOR & Word Order Disambiguation**: Format v2 embeds 32 positional vectors coupled with non-linear $GELU(E_i + P_i)$ pooling, mathematically distinguishing permutations like `"delivery refund"` from `"refund delivery"`.
- **3-Tier Decision Pipeline**: Classifies predictions into **Definite** (High confidence), **Ambiguous** (Borderline/narrow margin), or **Fallback** (Out-of-Distribution / High Shannon Entropy).
- **Multi-Intent Pipeline Support**: Automatically executes composite pipelines when secondary intent confidence meets multi-intent thresholds.
- **Lock-Free Atomic Hot-Swap & Telemetry**: Replace model weights on live traffic without locks (`sync/atomic.Pointer`), and stream drift queries into a bounded ring buffer for active learning.

---

## Benchmarks

Benchmarked on an AMD Ryzen 5 5600H (12 threads) running pure Go standard runtime (`go test -bench="." -benchmem`):

| Benchmark Target | Ops / Sec | Latency | Memory / Op | Allocations |
| :--- | :--- | :--- | :--- | :--- |
| **`BenchmarkPredictSlots`** | **33,433 ops/sec** | **29.91 μs** | **0 B/op** | **0 allocs/op** |
| **`BenchmarkPredictTokens`** | **33,126 ops/sec** | **30.18 μs** | **0 B/op** | **0 allocs/op** |
| **`BenchmarkForward`** | **33,091 ops/sec** | **30.21 μs** | **24 B/op** | **1 allocs/op** |
| **`BenchmarkGELU`** | **494,071 ops/sec** | **2.02 μs** | **0 B/op** | **0 allocs/op** |

---

## 3-Step Lifecycle

### Step 1: AI Design — Prepare Your Domain Knowledge (`data/sample_dataset.csv`)
Create a clean two-column CSV containing natural user queries and corresponding target labels:

```csv
text,label
I want to cancel my payment and request a refund,Refund
Where is my package and delivery tracking,Delivery
Forgot my account password please reset,Account
sent the return box a week ago when do i get my money back,Refund
yo i typed the wrong apt number please update address,Delivery
locked out of my account after 3 tries help pls,Account
```

### Step 2: Build Your Own AI — Compile Model Weights (`ib-train.exe`)
Train your domain vocabulary and neural weights into a compact Little-Endian binary (`intent.bin`) using the standalone CLI:

```bash
# Build the training tool once
go build -ldflags="-s -w" -o bin/ib-train.exe ./cmd/ib-train

# Compile 1,000+ domain rows in under 2 seconds (creates format v2 with positional embeddings)
./bin/ib-train.exe -data data/sample_dataset.csv -out weights/intent.bin -epochs 50 -lr 0.005 -vocab 250
```

### Step 3: AI-Powered Branching — Run In-Memory Routing (`go run main.go`)
Bind domain actions, multi-intent pipelines, and safety guards to Go functions:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"intellibranch/pkg/intellibranch"
)

func main() {
	// 1. Initialize in-memory router with calibrated confidence threshold
	router, err := intellibranch.NewRouter("weights/intent.bin", 0.75)
	if err != nil {
		log.Fatalf("Router init failure: %v", err)
	}

	// 2. Configure 3-Tier Policy and Active Learning Telemetry Buffer
	router.SetPolicy(intellibranch.DispatchPolicy{
		HighThreshold:     0.75,
		LowThreshold:      0.40,
		MarginCutoff:      0.15,
		MaxEntropy:        2.0,
		PipelineThreshold: 0.30,
	}).EnableTelemetry(1024)

	// 3. Bind standard and multi-intent handlers
	router.
		Bind("Refund", func(ctx context.Context, payload any) error {
			fmt.Printf("[ACTION: Refund] Processing: %v\n", payload)
			return nil
		}).
		Bind("Delivery", func(ctx context.Context, payload any) error {
			fmt.Printf("[ACTION: Delivery] Tracking: %v\n", payload)
			return nil
		}).
		BindPipeline("Refund", "Delivery", func(ctx context.Context, p, s string, payload any) error {
			fmt.Printf("[PIPELINE: %s -> %s] Processing combined return & shipment: %v\n", p, s, payload)
			return nil
		}).
		Ambiguous(func(ctx context.Context, p, s string, payload any) error {
			fmt.Printf("[AMBIGUOUS: %s vs %s] Requesting user confirmation: %v\n", p, s, payload)
			return nil
		}).
		Fallback(func(ctx context.Context, payload any) error {
			fmt.Printf("[FALLBACK: Safety Isolation] Isolated OOD query: %v\n", payload)
			return nil
		})

	// 4. Dispatch queries (Executes in ~30 microseconds)
	ctx := context.Background()
	_ = router.DispatchPipeline(ctx, "can u cancel order #49281 and update delivery?", "OrderPayload")
	_ = router.Dispatch(ctx, "tracking says delivered but mailbox is empty", "ShipmentPayload")
	_ = router.Dispatch(ctx, "Completely random gibberish 12345!@#$", "NoisePayload")

	// 5. Atomic Hot-Reloading & Telemetry Feedback (Concurrent & Lock-free)
	_ = router.Reload("weights/intent_v2.bin")
	events := router.DrainTelemetry()
	fmt.Printf("Harvested %d drift events for active learning retraining.\n", len(events))
}
```

---

## Observability & Whitebox Debugging

Need to understand why a query routed to a specific branch or why it fell back? Use `Inspect`:

```go
trace := router.Inspect("can u cancel order #49281? i bought it by mistake")
```

```json
{
  "input_text": "can u cancel order #49281? i bought it by mistake",
  "token_ids": [4, 5, 8, 12, 45, 98],
  "subwords": ["can", "u", "cancel", "order", "#", "mistake"],
  "unknown_token_ratio": 0.0,
  "class_probabilities": {
    "Account": 0.0012,
    "Delivery": 0.0035,
    "Refund": 0.9953
  },
  "predicted_label": "Refund",
  "secondary_label": "Delivery",
  "confidence": 0.9953,
  "margin": 0.9918,
  "entropy": 0.0351,
  "threshold": 0.75,
  "is_ambiguous": false,
  "is_pipeline": false,
  "is_fallback": false,
  "latency_micros": 30
}
```

---

## 6-Domain Multi-Task Demonstration Suite (`ib-demo`)

IntelliBranch comes bundled with a production-grade multi-task demonstration driver (`cmd/ib-demo`) validating 6 distinct enterprise use cases running on pure Go in microseconds:

| Domain | Model Binary | Dataset | Architectural Capabilities Demonstrated |
| :--- | :--- | :--- | :--- |
| **1. E-Commerce CS Gateway** | `weights/demo_cs.bin` | `data/demo_cs.csv` | Semantic XOR order resolution & Multi-Intent Composite Pipeline (`DispatchPipeline`) |
| **2. Semantic LLM Gateway** | `weights/demo_llm.bin` | `data/demo_llm.csv` | Microsecond local Go bypass ($0.00) & Shannon Entropy OOD Cloud LLM Fallback ($0.02) |
| **3. SRE High-Throughput Triage**| `weights/demo_sre.bin` | `data/demo_sre.csv` | **0 B/op (Zero Allocations)** stack inference via `PredictSlots` on 100k+ logs/sec |
| **4. Offline Edge IoT Control** | `weights/demo_iot.bin` | `data/demo_iot.csv` | Sub-180KB binary offline voice/text command routing with slang resilience |
| **5. Automated CI/CD Remediation**| `weights/demo_cicd.bin` | `data/demo_cicd.csv`| Build failure tail triage: Auto-Retry (`Network`), Scale-Up (`OOM`), Notify (`Syntax`) |
| **6. FinTech Memo Fraud Audit** | `weights/demo_fintech.bin`| `data/demo_fintech.csv`| Real-time scam interception & Borderline Step-Up 2FA Challenge (`Ambiguous`) |

### Running the Demonstration CLI

```bash
# Compile the demonstration driver
go build -ldflags="-s -w" -o bin/ib-demo.exe ./cmd/ib-demo

# Run all 6 domains sequentially
./bin/ib-demo.exe -domain all

# Or inspect an isolated domain
./bin/ib-demo.exe -domain sre
./bin/ib-demo.exe -domain llm
```

---

## Core Routing API Reference

| Method / Struct | Signature | Operational Role |
| :--- | :--- | :--- |
| **`NewRouter`** | `NewRouter(path string, threshold float64) (*Router, error)` | Loads v2 binary weights, initializes atomic model pointer, and builds 3-tier router. |
| **`SetPolicy`** | `.SetPolicy(policy DispatchPolicy) *Router` | Configures high/low thresholds, top-1/top-2 margin cutoff, OOD max entropy, and pipeline boundaries. |
| **`Bind`** | `.Bind(label string, handler RouteAction) *Router` | Associates a trained class with a Go handler: `func(ctx context.Context, payload any) error`. |
| **`BindPipeline`** | `.BindPipeline(p, s string, handler PipelineAction) *Router` | Registers composite handler triggered when primary and secondary intents are both eligible. |
| **`Ambiguous`** | `.Ambiguous(handler AmbiguousAction) *Router` | Intercepts borderline confidence or narrow margin queries to prompt user confirmation. |
| **`Fallback`** | `.Fallback(handler RouteAction) *Router` | Designates safety handler for low confidence, high unknown token ratio, or OOD entropy. |
| **`Dispatch`** | `.Dispatch(ctx context.Context, text string, payload any) error` | Evaluates 3-tier routing and executes bound branch in ~30 μs. |
| **`DispatchPipeline`**| `.DispatchPipeline(ctx context.Context, text string, payload any) error` | Executes multi-intent pipeline handlers if eligible, falling back to 3-tier routing. |
| **`Reload`** | `.Reload(path string) error` | Atomically swaps weights on live traffic without locks (`0 ns` stop-the-world). |
| **`EnableTelemetry`**| `.EnableTelemetry(capacity int) *Router` | Allocates thread-safe ring buffer capturing ambiguous, OOD, and pipeline requests. |
| **`DrainTelemetry`** | `.DrainTelemetry() []TelemetryEvent` | Extracts collected drift events in FIFO order for active learning retraining. |
| **`PredictSlots`** | `model.PredictSlots(text string, out *StaticInferenceResult) error` | Stack-allocated inference primitive achieving strictly **`0 B/op, 0 allocs/op`**. |

---

## Project Structure

```text
intellibranch/
├── cmd/
│   └── ib-train/          # Offline BPE + AdamW training CLI source
├── docs/
│   ├── MASTER_PLAN.md     # 6-stage architectural hardening & DoD specification
│   └── MANUAL.md          # Comprehensive manual, keyword guide & tutorial
├── pkg/
│   └── intellibranch/     # Pure-Go zero-dependency core engine
│       ├── binary.go      # Format v1 & v2 Little-Endian parser and serializer
│       ├── ops.go         # SafeClamp, GELU, Softmax, MatMul, and Non-Linear Pooling
│       ├── runtime.go     # Zero-alloc PredictSlots, Shannon Entropy & In-memory model
│       ├── telemetry.go   # Thread-safe ring buffer for active learning feedback
│       ├── tokenizer.go   # Pure Go BPE subword tokenizer
│       ├── trainer.go     # AdamW backprop trainer with positional embedding learning
│       └── router.go      # 3-Tier router, atomic reload, and pipeline dispatch
├── weights/
│   └── .gitkeep           # Directory placeholder for serialized weights
├── data/
│   └── sample_dataset.csv # 1,000+ domain training rows
├── main.go                # Server entrypoint with auto-train bootstrap
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

---

## License

This project is licensed under the MIT License.
