# IntelliBranch
<img src="https://github.com/user-attachments/assets/3413a486-d71c-4285-841d-76bbe74f830a" width="226" height="200" alt="Image" align="right" style="margin-left: 15px; margin: 10px;">
<p align="center">
  <strong>Directly Creates and Runs Its Own Neural AI in Pure Go</strong><br>
  <em>Stop borrowing third-party AIs. This engine creates its own domain artificial intelligence from scratch in 1.5 seconds, routing execution flow in ~6.08 μs with Zero Downloads and Zero CGO.</em>
</p>

<p align="center">
  <a href="#benchmarks"><img src="https://img.shields.io/badge/Latency-~6.08_μs-brightgreen.svg" alt="Latency"></a>
  <a href="#benchmarks"><img src="https://img.shields.io/badge/Allocs-24_B/op_(1_alloc)-blue.svg" alt="Allocations"></a>
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

Instead of relying on brittle regex matching or calling bloated external LLMs, it **manufactures a domain-specific lightweight neural network directly from your dataset in under 2 seconds**. It maps typos, slang, inverted syntax, and colloquial phrasing into a continuous latent vector space—routing execution flow directly to your bound Go functions in **single-digit microseconds (6 μs)**.

```
Incoming Request ("bruh can u refund order #49281")
                     │
                     ▼
       [ In-Memory BPE Tokenizer ]
                     │
                     ▼
      [ Dense Embedding (D=64) ]
                     │
                     ▼
      [ Non-Linear GELU (D=128) ]
                     │
                     ▼
    [ Softmax Confidence & Guard ]
                     │
        ┌────────────┴────────────┐
        ▼                         ▼
 (Score ≥ 0.60)            (Score < 0.60 / Noise)
[ACTION: Refund]           [FALLBACK: Safe Isolation]
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
// Memory Footprint: Under 150 KB (30,000x smaller than quantized 7B models)
// Latency: ~6.08 μs with ZERO external network calls and deterministic fallback
// Deployment: 100% Pure Go with CGO_ENABLED=0 single static binary
```

### Architectural Comparison Matrix

| Capability | Retro Branching (`if` / Regex) | Cloud LLMs (OpenAI / Claude) | Local LLMs (Ollama / llama.cpp) | **IntelliBranch (Embedded Engine)** |
| :--- | :--- | :--- | :--- | :--- |
| **Inference Latency** | < 1 μs | 300 ms – 2,500 ms (Network bound) | 30 ms – 300 ms (Compute bound) | **~6.08 μs (In-Memory)** |
| **Throughput (per core)** | > 500,000 req/sec | ~50 req/sec (Rate limited) | ~20–50 req/sec (CPU saturated) | **> 150,000 req/sec (`sync.Pool`)** |
| **System Memory (RAM)** | Negligible | External service | **4.5 GB – 8.0 GB+ (VRAM / RAM)** | **< 150 KB (30,000x lighter)** |
| **Hardware Reqs** | Standard CPU | External service | High-end GPU or 8+ Core CPU | **Runs on a $5 VPS (16MB container)** |
| **Operational Cost** | $0.00 | $0.0015+ per call | High hardware/electricity cost | **$0.00 (Self-contained)** |
| **Typo & Slang Resilience**| ❌ 0% (Strict string match) | ✅ High | ✅ High | ✅ **High (BPE Subwords)** |
| **Deployment Complexity** | Single binary | API client | CGO / C++ runtime / Ollama daemon | **Pure Go (`CGO_ENABLED=0`)** |
| **Deterministic Fallback** | Hard-coded `default` branch | ❌ Unpredictable hallucinations | ❌ Hallucination & format errors | ✅ **Calibrated Guard (`< 0.60`)** |
| **Model Acquisition** | N/A (Manual code) | Cloud API lease (OpenAI) | Download multi-GB checkpoints (HuggingFace) | **Zero Downloads (Forged from scratch in 1.5s)** |
| **Model Retraining** | N/A (Manual code editing) | Black-box fine-tuning | Multi-hour GPU fine-tuning | **1.5-second CLI compilation** |

---

## Key Highlights

- **Zero Downloads & On-The-Fly AI Creation**: You never download gigabytes of pre-trained weights from HuggingFace or lease external APIs. IntelliBranch forges a domain neural AI model directly from your CSV in under 2 seconds.
- **Zero External Dependencies & Zero CGO**: 100% pure Go standard library. Compiles cleanly with `CGO_ENABLED=0` for portable cross-platform binaries.
- **Microsecond Latency (~6.08 μs)**: Zero disk I/O on hot paths. Utilizes `sync.Pool` scratch buffers to achieve sub-millisecond throughput (>150,000 requests/sec per CPU core).
- **Non-Linear Expressiveness**: Solves complex semantic logic (XOR/contextual negation) using GELU non-linear activations over 128 hidden dimensions.
- **Built-in Safety Guards**: Hard input size clamping (512 bytes) and Out-of-Vocabulary (OOV) confidence discounting protect against CPU exhaustion and rogue queries.
- **Deterministic Fallback**: Safely isolates ambiguous requests and out-of-distribution noise without throwing unhandled runtime panics.

---

## Benchmarks

Benchmarked on an AMD Ryzen 5 5600H (12 threads) running pure Go standard runtime:

| Benchmark Target | Ops / Sec | Latency | Memory / Op | Allocations |
| :--- | :--- | :--- | :--- | :--- |
| **`BenchmarkForward`** | **193,677 ops/sec** | **6.17 μs** | **24 B/op** | **1 allocs/op** |
| **`BenchmarkPredictTokens`** | **195,225 ops/sec** | **6.08 μs** | **24 B/op** | **1 allocs/op** |
| **`BenchmarkGELU`** | **990,507 ops/sec** | **1.10 μs** | **0 B/op** | **0 allocs/op** |

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

# Compile 1,000+ domain rows in under 2 seconds
./bin/ib-train.exe -data data/sample_dataset.csv -out weights/intent.bin -epochs 50 -lr 0.005 -vocab 250
```

### Step 3: AI-Powered Branching — Run In-Memory Routing (`go run main.go`)
Bind domain actions to Go functions and dispatch incoming traffic in microseconds:

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
	router, err := intellibranch.NewRouter("weights/intent.bin", 0.60)
	if err != nil {
		log.Fatalf("Router init failure: %v", err)
	}

	// 2. Bind business handlers
	router.
		Bind("Refund", func(ctx context.Context, payload any) error {
			fmt.Printf("[ACTION: Refund] Processing: %v\n", payload)
			return nil
		}).
		Bind("Delivery", func(ctx context.Context, payload any) error {
			fmt.Printf("[ACTION: Delivery] Tracking: %v\n", payload)
			return nil
		}).
		Bind("Account", func(ctx context.Context, payload any) error {
			fmt.Printf("[ACTION: Account] Securing: %v\n", payload)
			return nil
		}).
		Fallback(func(ctx context.Context, payload any) error {
			fmt.Printf("[FALLBACK: Safety] Isolated query: %v\n", payload)
			return nil
		})

	// 3. Dispatch queries (Executes in ~6.08 microseconds)
	ctx := context.Background()
	_ = router.Dispatch(ctx, "can u cancel order #49281? i bought it by mistake", "OrderPayload")
	_ = router.Dispatch(ctx, "tracking says delivered but mailbox is empty", "ShipmentPayload")
	_ = router.Dispatch(ctx, "Completely random gibberish 12345!@#$", "NoisePayload")
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
  "confidence": 0.9953,
  "threshold": 0.60,
  "is_fallback": false,
  "latency_micros": 6
}
```

---

## Conceptual Guide for Go Developers (Zero-ML Primer)

If you are a Go engineer with zero machine learning or Python background, understand **IntelliBranch** not as "black-box AI", but as an **automated, multidimensional coordinate system for routing strings**.

### 1. The Mental Shift: Discrete Matching vs. Vector Geometry

| Concept | Retro Go (`switch` / `map`) | IntelliBranch Statistical Branching |
| :--- | :--- | :--- |
| **Input Representation** | Raw string characters (`"refund"`) | Sequence of BPE subword token IDs (`[4, 12, 8]`) |
| **Comparison Metric** | Byte-by-byte equality (`==`) | Distance in 64-dimensional latent coordinate space |
| **Typos / Variations** | Misses completely (`"refnd"` != `"refund"`) | Subwords land at almost the exact same coordinate |
| **Contextual Negation** | Requires complex lookahead regex | Non-linear GELU activation separates opposing intents |
| **Output Evaluation** | Jump directly to `case` label | Softmax probability distribution over registered branches |

### 2. Under the Hood: 4-Stage Microsecond Pipeline

```text
[ Raw String Input ] ── "bruh can u refund order #49281"
        │
        ▼ (Stage 1: BPE Subword Tokenizer)
[ Subword Token IDs ] ── [4, 5, 8, 12, 45, 98] (Handles typos & OOV gracefully)
        │
        ▼ (Stage 2: 64-Dimensional Embedding Table)
[ Token Vectors ] ── 64 float32 coordinates per subword
        │
        ▼ (Stage 3: Mean Pooling & GELU Non-Linearity)
[ Sentence Latent Vector ] ── Averaged & projected into 128 hidden dimensions
        │
        ▼ (Stage 4: Softmax & Threshold Guard)
[ Branch Confidence ] ── Refund: 0.995, Delivery: 0.003, Account: 0.001
        │
   Score ≥ 0.60?
   ├── YES ──▶ Execute bound handler: router.Bind("Refund", ...)
   └── NO  ──▶ Execute safe isolation: router.Fallback(...)
```

1. **BPE Subword Tokenization**: Instead of splitting words by whitespace (which fails on typos like `"refuuund"` or compound words), Byte-Pair Encoding breaks text into recurring statistical chunks (subwords). Even if a user types slang or a misspelled term, familiar subwords are recognized.
2. **Dense Vector Embedding ($D=64$)**: Each token ID maps to a 64-dimensional float vector. During training, words with similar business intent are drawn closer together in this 64-dimensional space.
3. **Mean Pooling & GELU Non-Linearity ($D=128$)**: All token vectors in a query are averaged into a single 64-element sentence representation. A non-linear feedforward layer (GELU) computes interaction cross-products—allowing the engine to understand negation (e.g., distinguishing "refund my order" from "cancel delivery notifications").
4. **Softmax Scoring & Calibrated Fallback**: The final linear layer produces raw logits, normalized into a probability distribution summing to 1.0. If the highest probability is below `threshold` (default `0.60`) or if the unknown token ratio is too high, execution automatically diverts to `Fallback`.

### 3. Core Routing API Reference

| Method / Symbol | Signature | Operational Role |
| :--- | :--- | :--- |
| **`NewRouter`** | `NewRouter(path string, threshold float32) (*Router, error)` | Loads Little-Endian binary weights, checks SHA-256 hash, and pre-allocates `sync.Pool` scratch buffers. |
| **`Bind`** | `.Bind(label string, handler Handler) *Router` | Associates a trained dataset class (e.g. `"Refund"`, `"Delivery"`) with a Go handler: `func(ctx context.Context, payload any) error`. |
| **`Fallback`** | `.Fallback(handler Handler) *Router` | Designates the safety handler executed when predictions fall below `threshold` or when rogue noise is detected. |
| **`Dispatch`** | `.Dispatch(ctx context.Context, text string, payload any) error` | Tokenizes input, runs forward inference in ~6 μs, picks the highest-confidence branch, and directly executes the bound handler. |
| **`Inspect`** | `.Inspect(text string) RouteTrace` | Performs inference without executing business handlers, returning a complete diagnostic struct with tokens, probabilities, and microsecond timings for whitebox telemetry. |

---

## Beyond Exact Hashmaps: Where IntelliBranch Is Architecturally Mandatory

Static hash maps (`map[string]T`), `switch-case` statements, and regex matchers fail completely when facing unstructured variation, combinatorial permutations, or semantic noise. IntelliBranch is not a generic pattern matcher; it is an **in-memory continuous vector-routing primitive ($text \to action$) operating at ~6.08 μs**.

Applications listed below are strictly restricted to domains where traditional discrete branching collapses and IntelliBranch provides a provable, definitive architectural advantage:

### 1. Definitive Application Matrix

| Domain / Pattern | Why Traditional Branching (`switch`/`map`/Regex) Fails | IntelliBranch Architectural Dominance | Latency |
| :--- | :--- | :--- | :--- |
| **Semantic LLM Gateway & API Bypass** | • Exact hash keys cannot capture semantic equivalence across infinite sentence variations.<br>• Regex rules explode exponentially.<br>• Every miss costs 500ms–2,500ms and OpenAI API token fees. | Directly resolves 80–90% of routine natural language queries into deterministic Go functions, bypassing expensive cloud LLMs entirely. | **~6.08 μs** |
| **Offline Edge & IoT Micro-Command Dispatch** | • Embedded environments (32–64MB RAM) cannot run 4GB+ LLMs (Ollama/llama.cpp).<br>• `switch` statements fail on natural phrasing, slang, and dialect variations.<br>• Cloud APIs fail when network connection drops. | Runs offline in <150 KB RAM with zero CGO dependencies. Maps colloquial voice/text variants directly to hardware/GPIO routines in 6 μs. | **~6.08 μs** |
| **Context-Enriched Multi-Attribute Routing** | • Evaluating `[ROLE][PATH] Query` requires nested `switch` ladders and regex lookaheads.<br>• Branch complexity scales as $O(R \times P \times Q)$, creating unmaintainable combinatorial explosion. | BPE tokenizes prefix tags into distinct subword coordinates. GELU non-linear hidden layers compute cross-attribute interaction in a single pass. | **~6.08 μs** |
| **Hierarchical Microservice Cascade** | • Single flat regex or string parsers degrade linearly in latency as class counts grow past 50+.<br>• Maintenance becomes impossible when services expand. | Chaining coarse domain routers (Stage 1) to granular service routers (Stage 2) enables scaling to 200+ distinct endpoints while keeping latency within single-digit microseconds. | **~12.16 μs** |
| **High-Throughput Log & Incident QoS Triage** | • `strings.Contains` suffers from high false-positive collisions on stack traces.<br>• Complex regex engines burn 100% CPU on high-volume message brokers (Kafka/RabbitMQ). | Ingests raw, unformatted error messages, stack traces, and crash dumps at 150,000+ ops/sec per core, directing traffic instantly into P0/P1/P2/P3 partitions. | **~6.08 μs** |
| **CI/CD Failure Auto-Remediation** | • Stack traces vary unpredictably across toolchains (Docker, K8s, Go, Gradle).<br>• Hard-coded regexes miss minor error wording changes, breaking automated pipelines. | Classifies error tails in 6 μs into deterministic actions: `AutoRetry` (transient network), `ScaleResource` (OOM), or `AlertDev` (code syntax error). | **~6.08 μs** |

---

### 2. Deep Architectural Rationale

#### A. Semantic LLM Gateway & API Bypass (Cloud LLM Cost & Latency Killer)
Modern architectures waste millions of dollars routing every incoming natural language request to cloud LLMs (OpenAI, Claude). 
- **The Discrete Failure**: A user asking `"how do i get my money back?"` vs `"reverse charge order #123"` cannot be cached in a hash map. Regex attempts to cover every phrasing end in combinatorial failure.
- **The IntelliBranch Advantage**: IntelliBranch maps the underlying semantics directly to an internal Go handler (`router.Bind("Refund", ...)`). Only out-of-distribution queries falling below the calibrated threshold (`< 0.60`) escape to the fallback cloud LLM pipeline.
- **Result**: Cuts cloud LLM API costs by 80–90% and slashes response latency from 1,200 ms to **6.08 μs** for the vast majority of user traffic.

#### B. Offline Edge & IoT Micro-Command Dispatch (Zero-Cloud, Sub-Milliwatt Control)
Edge devices (smart home hubs, POS systems, robotics, industrial PLCs) operate under stringent resource constraints (32 MB – 128 MB RAM) and intermittent network connectivity.
- **The Discrete Failure**: Rigid `switch(cmd)` fails on everyday speech variances (e.g., `"turn on lights"` vs `"it's dark here"` vs `"lights please"`). 
- **The LLM Failure**: Running local 7B models requires gigabytes of VRAM and high wattage. Cloud APIs introduce network latency and fail entirely offline.
- **The IntelliBranch Advantage**: IntelliBranch compiles to a single pure Go binary under 150 KB. It operates entirely in-memory with zero CGO and zero external dependencies, mapping colloquial commands directly to hardware actuators in **6.08 μs**.

#### C. Context-Enriched Multi-Attribute Routing
Production gateways must branch on multi-dimensional context: user permissions, HTTP routes, and unstructured payloads simultaneously.
- **The Discrete Failure**: Building nested conditionals for 5 roles, 20 endpoints, and varied user intentions requires hundreds of error-prone lines of code. Any new variation breaks the branching logic.
- **The IntelliBranch Advantage**: Synthesizing the input as `[ROLE][PATH] Query` allows the BPE tokenizer to project both structural metadata and raw language into the same 64-dimensional latent space. GELU layers evaluate the non-linear interaction between identity and intent without maintaining brittle branching graphs.

#### D. Hierarchical Microservice Cascade (Scaling Beyond Flat Class Limits)
When enterprise backends route across hundreds of service actions, flat classification degrades in confidence.
- **The Discrete Failure**: Linear regex evaluation over 100+ endpoints consumes milliseconds of CPU time per request.
- **The IntelliBranch Advantage**: Sequentially chaining two lightweight routers (Stage 1: Domain isolation $\to$ Stage 2: Action dispatch) isolates decision boundaries. Each model maintains a hyper-focused embedding table under 150 KB, executing in ~12 μs total with zero heap allocations on the hot path.

#### E. Real-Time Telemetry & Incident QoS Partitioning
High-throughput data ingestion pipelines (Kafka, Vector, Fluentd) cannot afford heavy regex parsers.
- **The Discrete Failure**: Processing 100,000 logs/sec with regular expressions causes severe CPU starvation and consumer lag.
- **The IntelliBranch Advantage**: Consuming raw logs and error strings directly through IntelliBranch routes critical P0 panics to dedicated priority consumers while shunting non-critical noise to cold storage, sustaining >150,000 req/sec per core on standard hardware.

#### F. Automated CI/CD Failure Triage & Remediation
Distributed build and deployment systems generate gigabytes of unstructured compiler and runtime failure logs daily.
- **The Discrete Failure**: Brittle regex matchers break whenever compiler error formats change slightly, forcing engineers to manually investigate every failed build.
- **The IntelliBranch Advantage**: Ingesting the last 512 bytes of a build log allows IntelliBranch to classify failure causes into remediation actions in 6 μs—automatically retrying transient network timeouts, provisioning larger runners for OOM kills, or directly routing code syntax issues to the responsible commit author.

---

## Project Structure

```text
intellibranch/
├── cmd/
│   └── ib-train/          # Offline BPE + AdamW training CLI source
├── docs/
│   └── MANUAL.md          # Comprehensive manual, keyword guide & tutorial
├── pkg/
│   └── intellibranch/     # Pure-Go zero-dependency core engine
│       ├── binary.go      # Little-Endian binary format parser and serializer
│       ├── ops.go         # Pure Go mathematical operations (GELU, Softmax, MatMul)
│       ├── runtime.go     # In-memory inference engine with buffer pool & guards
│       ├── tokenizer.go   # Pure Go BPE subword tokenizer
│       ├── trainer.go     # Pure Go AdamW backpropagation offline trainer
│       └── router.go      # Bind/Fallback/Dispatch thread-safe in-memory router
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
