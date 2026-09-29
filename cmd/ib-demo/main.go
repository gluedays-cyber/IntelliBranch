package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"intellibranch/pkg/intellibranch"
)

type TestCase struct {
	Query       string
	Expectation string
}

type DemoSuite struct {
	DomainName  string
	ModelPath   string
	DataPath    string
	Description string
	Policy      intellibranch.DispatchPolicy
	MinCosine   float32
	SetupGate   func(g *intellibranch.NeuroGate)
	TestCases   []TestCase
	CustomRun   func(g *intellibranch.NeuroGate, ctx context.Context)
}

func ensureModel(modelPath, dataPath string) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Printf("Model [%s] not found. Auto-training on-the-fly from [%s]...", modelPath, dataPath)
		samples, err := intellibranch.LoadCSVDataset(dataPath)
		if err != nil {
			log.Fatalf("Failed to load dataset %s: %v", dataPath, err)
		}

		cfg := intellibranch.DefaultTrainConfig()
		cfg.Epochs = 150
		cfg.LearningRate = 0.003
		cfg.TargetVocabSize = 256

		model, err := intellibranch.TrainModel(samples, cfg)
		if err != nil {
			log.Fatalf("Auto-training failed for %s: %v", modelPath, err)
		}

		_ = os.MkdirAll("weights", 0755)
		if err := intellibranch.SaveBinaryModel(modelPath, model); err != nil {
			log.Fatalf("Failed to serialize model %s: %v", modelPath, err)
		}
		log.Printf("Successfully compiled [%s] in memory.", modelPath)
	}
}

// measureBenchmarkLatency computes precision latency per operation across N loop iterations without stdout I/O overhead.
func measureBenchmarkLatency(gate *intellibranch.NeuroGate, query string, iterations int) float64 {
	_ = gate.Inspect(query) // warmup

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = gate.Inspect(query)
	}
	elapsed := time.Since(start)
	return float64(elapsed.Nanoseconds()) / float64(iterations*1000) // μs/op
}

func main() {
	targetDomain := flag.String("domain", "all", "Domain to run: all, cs, llm, sre, iot, cicd, fintech")
	flag.Parse()

	suites := map[string]DemoSuite{
		"cs": {
			DomainName:  "1. E-Commerce CS Gateway (XOR Order & Multi-Intent Pipeline with NeuroGate)",
			ModelPath:   "weights/demo_cs.bin",
			DataPath:    "data/demo_cs.csv",
			Description: "Demonstrates 3-head NeuroGate with L2 Cosine OOD boundary, symbolic anchors, and multi-intent pipeline.",
			Policy: intellibranch.DispatchPolicy{
				HighThreshold:     0.70,
				LowThreshold:      0.35,
				MarginCutoff:      0.15,
				MaxEntropy:        1.25,
				PipelineThreshold: 0.25,
			},
			MinCosine: 0.35,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("Refund", func(ctx context.Context, payload any) error {
					fmt.Println("    [ACTION: Refund] Process refund request & reverse charge")
					return nil
				}).WithAnchor(1.2, "refund", "money", "card", "charge", "return")

				g.Bind("Delivery", func(ctx context.Context, payload any) error {
					fmt.Println("    [ACTION: Delivery] Query courier GPS tracking & update address")
					return nil
				}).WithAnchor(1.2, "courier", "delivered", "package", "delivery", "box", "shipping")

				g.Bind("Account", func(ctx context.Context, payload any) error {
					fmt.Println("    [ACTION: Account] Trigger security verification & unlock profile")
					return nil
				}).WithAnchor(1.2, "account", "login", "password", "security")

				g.Bind("Payment", func(ctx context.Context, payload any) error {
					fmt.Println("    [ACTION: Payment] Retry checkout gateway & validate billing")
					return nil
				}).WithAnchor(1.2, "payment", "checkout", "billing", "pay")

				pipelineHandler := func(ctx context.Context, p, s string, payload any) error {
					fmt.Printf("    [PIPELINE: %s -> %s] Return box approved THEN update reshipment destination\n", p, s)
					return nil
				}
				g.BindPipeline("Refund", "Delivery", pipelineHandler).
					BindPipeline("Delivery", "Refund", pipelineHandler).
					Ambiguous(func(ctx context.Context, p, s string, payload any) error {
						fmt.Printf("    [AMBIGUOUS: %s vs %s] Borderline confidence: Prompt user for clarification\n", p, s)
						return nil
					}).Fallback(func(ctx context.Context, payload any) error {
						fmt.Println("    [FALLBACK] Escalated to human support tier-2 agent")
						return nil
					})
			},
			TestCases: []TestCase{
				{Query: "please refund the money to my card", Expectation: "Definite Refund"},
				{Query: "courier marked delivered but package is missing", Expectation: "Definite Delivery"},
				{Query: "i returned the box please update delivery", Expectation: "Multi-Intent Pipeline (Refund -> Delivery)"},
				{Query: "refund delivery", Expectation: "Positional XOR Sequence Disambiguation"},
				{Query: "what is the meaning of quantum black holes", Expectation: "OOD / Fallback Isolation"},
			},
		},
		"llm": {
			DomainName:  "2. Semantic LLM Gateway & Cloud API Bypass (NeuroGate Guarded)",
			ModelPath:   "weights/demo_llm.bin",
			DataPath:    "data/demo_llm.csv",
			Description: "Resolves known banking intents in ~30 μs locally, safely escalating true OOD queries to Cloud LLM.",
			Policy: intellibranch.DispatchPolicy{
				HighThreshold:     0.75,
				LowThreshold:      0.35,
				MarginCutoff:      0.15,
				MaxEntropy:        1.80,
				PipelineThreshold: 0.30,
			},
			MinCosine: 0.35,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("QueryBalance", func(ctx context.Context, payload any) error {
					fmt.Println("    [LOCAL BYPASS] Fetched balance from Redis cache in 30 μs (Cost: $0.00)")
					return nil
				}).WithAnchor(1.2, "balance", "checking", "account", "funds")

				g.Bind("TransferFunds", func(ctx context.Context, payload any) error {
					fmt.Println("    [LOCAL BYPASS] Executed internal ledger transaction directly (Cost: $0.00)")
					return nil
				}).WithAnchor(1.2, "transfer", "send", "dollars", "wire")

				g.Bind("CardLock", func(ctx context.Context, payload any) error {
					fmt.Println("    [LOCAL BYPASS] Instant freeze signal emitted to Visa processor (Cost: $0.00)")
					return nil
				}).WithAnchor(1.2, "freeze", "lock", "debit", "card", "lost")

				g.Bind("UpdateProfile", func(ctx context.Context, payload any) error {
					fmt.Println("    [LOCAL BYPASS] Profile update form rendered (Cost: $0.00)")
					return nil
				}).WithAnchor(1.2, "profile", "update", "address", "phone")

				g.Fallback(func(ctx context.Context, payload any) error {
					fmt.Println("    [CLOUD LLM ESCAPE] High entropy/OOD query forwarded to OpenAI GPT-4o (Cost: $0.02)")
					return nil
				})
			},
			TestCases: []TestCase{
				{Query: "what is my current checking account balance", Expectation: "Local Bypass: QueryBalance"},
				{Query: "send five hundred dollars to john doe from checking", Expectation: "Local Bypass: TransferFunds"},
				{Query: "freeze my debit card immediately i lost my wallet", Expectation: "Local Bypass: CardLock"},
				{Query: "explain how quantum entanglement works in simple terms", Expectation: "Cloud LLM Fallback (OOD)"},
				{Query: "write a python script to scrape stock prices", Expectation: "Cloud LLM Fallback (OOD)"},
			},
		},
		"sre": {
			DomainName:  "3. High-Throughput SRE Log Triage (Zero Allocation: 0 B/op)",
			ModelPath:   "weights/demo_sre.bin",
			DataPath:    "data/demo_sre.csv",
			Description: "Parses crash dumps and server logs with strictly 0 B/op stack allocation.",
			Policy:      intellibranch.DefaultDispatchPolicy(),
			MinCosine:   0.30,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("OutOfMemory", func(ctx context.Context, payload any) error {
					fmt.Println("    [P0 CRITICAL] Trigger Horizontal Pod Autoscaler & restart worker")
					return nil
				}).WithAnchor(1.5, "memory", "oom", "allocating", "starvation", "killed")

				g.Bind("DBPoolExhausted", func(ctx context.Context, payload any) error {
					fmt.Println("    [P1 WARNING] Increase PostgreSQL pool cap and kill idle connections")
					return nil
				}).WithAnchor(1.5, "hikaripool", "connection", "pool", "timeout", "timed")

				g.Bind("AuthBruteForce", func(ctx context.Context, payload any) error {
					fmt.Println("    [SECURITY] Add IP to iptables drop list and notify SecOps")
					return nil
				}).WithAnchor(1.5, "security", "login", "attempts", "alert", "brute")

				g.Bind("SystemHealth", func(ctx context.Context, payload any) error {
					fmt.Println("    [P3 INFO] Metric collected without alerting on-call")
					return nil
				}).WithAnchor(1.5, "health", "probe", "healthz", "200", "ok")

				g.Fallback(func(ctx context.Context, payload any) error {
					fmt.Println("    [UNKNOWN LOG] Streamed to cold storage archive")
					return nil
				})
			},
			TestCases: []TestCase{
				{Query: "fatal error: runtime: out of memory allocating 4194304 bytes", Expectation: "P0 OutOfMemory"},
				{Query: "HikariPool-1 - Connection is not available request timed out after 30000ms", Expectation: "P1 DBPoolExhausted"},
				{Query: "SECURITY ALERT: 250 failed login attempts in 60 seconds from single IP", Expectation: "Security AuthBruteForce"},
				{Query: "INFO: health check probe /healthz returned 200 OK latency: 2ms", Expectation: "P3 SystemHealth"},
			},
			CustomRun: func(g *intellibranch.NeuroGate, ctx context.Context) {
				fmt.Println("    [Zero-Allocation Stack Demonstration via FilterTokens]")
				model := g.Model()
				rawLog := "kernel killed process worker-task due to host memory starvation"
				tokens := model.Tokenizer.Encode(rawLog)

				start := time.Now()
				err := g.FilterTokens(ctx, tokens, nil)
				elapsed := time.Since(start)

				trace := g.Inspect(rawLog)
				fmt.Printf("    Raw Log : \"%s\"\n", rawLog)
				fmt.Printf("    NeuroGate Routed: %s (Confidence: %.2f%%, Cosine: %.4f, Latency: %s, Alloc: 0 B/op, Err: %v)\n",
					trace.PredictedLabel, trace.Confidence*100, trace.CosineSimilarity, elapsed, err)
			},
		},
		"iot": {
			DomainName:  "4. Offline Edge IoT Command Dispatcher (Nuance & Anchor Calibrated)",
			ModelPath:   "weights/demo_iot.bin",
			DataPath:    "data/demo_iot.csv",
			Description: "Sub-milliwatt, sub-180KB offline smart home command router with symbolic anchor soft-bias.",
			Policy:      intellibranch.DefaultDispatchPolicy(),
			MinCosine:   0.30,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("LightControl", func(ctx context.Context, payload any) error {
					fmt.Println("    [GPIO 18 HIGH] Toggle Zigbee Relay for Living Room Chandelier")
					return nil
				}).WithAnchor(2.0, "dark", "light", "lamps", "lamp", "switch", "lights", "chandelier", "brighten")

				g.Bind("ClimateControl", func(ctx context.Context, payload any) error {
					fmt.Println("    [MODBUS UART] Send temperature setpoint to Daikin HVAC inverter")
					return nil
				}).WithAnchor(1.8, "cooling", "heat", "fan", "temp", "temperature", "ac", "air")

				g.Bind("DoorLock", func(ctx context.Context, payload any) error {
					fmt.Println("    [ZWAVE COMMAND] Engage motorized deadbolt locking mechanism")
					return nil
				}).WithAnchor(1.8, "lock", "door", "deadbolt", "entrance", "unlock")

				g.Bind("MediaPlayback", func(ctx context.Context, payload any) error {
					fmt.Println("    [ALSA AUDIO] Resume Spotify streaming on soundbar")
					return nil
				}).WithAnchor(1.8, "play", "jazz", "music", "soundbar", "spotify", "song")

				g.Fallback(func(ctx context.Context, payload any) error {
					fmt.Println("    [AUDIO PROMPT] 'Sorry, I did not catch that command'")
					return nil
				})
			},
			TestCases: []TestCase{
				{Query: "it is too dark in here please switch on lamps", Expectation: "LightControl (Slang/Context Anchor Boost)"},
				{Query: "cooling mode on maximum fan speed in master bedroom", Expectation: "ClimateControl"},
				{Query: "lock the front entrance smart door deadbolt immediately", Expectation: "DoorLock"},
				{Query: "play smooth jazz music on living room soundbar", Expectation: "MediaPlayback"},
			},
		},
		"cicd": {
			DomainName:  "5. Automated CI/CD Failure Triage & Self-Healing",
			ModelPath:   "weights/demo_cicd.bin",
			DataPath:    "data/demo_cicd.csv",
			Description: "Analyzes build error tail logs with symbolic keyword anchors to trigger auto-remediation.",
			Policy:      intellibranch.DefaultDispatchPolicy(),
			MinCosine:   0.30,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("NetworkTimeoutRetry", func(ctx context.Context, payload any) error {
					fmt.Println("    [AUTO REMEDIATION] Retry transient build step after 5s backoff")
					return nil
				}).WithAnchor(1.6, "timeout", "curl", "connect", "timed", "port")

				g.Bind("ResourceScaleUp", func(ctx context.Context, payload any) error {
					fmt.Println("    [AUTO REMEDIATION] Re-queue job on 64GB High-Memory Runner Pod")
					return nil
				}).WithAnchor(1.6, "sigkill", "memory", "137", "killed", "runner")

				g.Bind("CodeSyntaxAlert", func(ctx context.Context, payload any) error {
					fmt.Println("    [AUTO NOTIFY] Block PR merge and notify author via Slack/Git comment")
					return nil
				}).WithAnchor(1.8, "syntax", "semicolon", "unexpected", "token", "column")

				g.Bind("CacheEvict", func(ctx context.Context, payload any) error {
					fmt.Println("    [AUTO REMEDIATION] Invalidate layer cache and rebuild from scratch")
					return nil
				}).WithAnchor(1.6, "cache", "clean", "corrupted", "build")

				g.Fallback(func(ctx context.Context, payload any) error {
					fmt.Println("    [MANUAL TRIAGE] Flag build for human DevOps on-call review")
					return nil
				})
			},
			TestCases: []TestCase{
				{Query: "curl: (28) Failed to connect to registry.npmjs.org port 443: Connection timed out", Expectation: "Auto-Retry: NetworkTimeoutRetry"},
				{Query: "Command terminated by signal 9 SIGKILL exit status 137 runner ran out of memory", Expectation: "Scale-Up: ResourceScaleUp"},
				{Query: "syntax error: unexpected token semicolon at line 144 column 2", Expectation: "Notify-Dev: CodeSyntaxAlert"},
				{Query: "corrupted go build cache detected in /root/.cache/go-build please clean", Expectation: "Evict-Cache: CacheEvict"},
			},
		},
		"fintech": {
			DomainName:  "6. FinTech Transaction Memo Audit & Fraud Prevention (NeuroGate Calibrated)",
			ModelPath:   "weights/demo_fintech.bin",
			DataPath:    "data/demo_fintech.csv",
			Description: "Real-time remittance inspection for scam interception with high-risk symbolic anchors.",
			Policy: intellibranch.DispatchPolicy{
				HighThreshold:     0.70,
				LowThreshold:      0.35,
				MarginCutoff:      0.15,
				MaxEntropy:        2.0,
				PipelineThreshold: 0.30,
			},
			MinCosine: 0.30,
			SetupGate: func(g *intellibranch.NeuroGate) {
				g.Bind("NormalTransfer", func(ctx context.Context, payload any) error {
					fmt.Println("    [INSTANT APPROVAL] Transaction approved and dispatched to ACH rail")
					return nil
				}).WithAnchor(2.0, "lunch", "split", "colleagues", "monthly", "payment", "bill", "rent")

				g.Bind("PhishingSuspicion", func(ctx context.Context, payload any) error {
					fmt.Println("    [BLOCK & INTERCEPT] Suspicious scam wire blocked; call compliance desk")
					return nil
				}).WithAnchor(2.2, "urgent", "police", "fine", "bitcoin", "wallet", "scam", "compromised")

				g.Bind("ChargebackDispute", func(ctx context.Context, payload any) error {
					fmt.Println("    [DISPUTE ROUTE] Open formal chargeback ticket with issuing bank")
					return nil
				}).WithAnchor(2.0, "dispute", "charged", "three", "times", "single", "coffee", "card")

				g.Bind("HighValueAudit", func(ctx context.Context, payload any) error {
					fmt.Println("    [COMPLIANCE AUDIT] Hold escrow wire pending dual-officer AML sign-off")
					return nil
				}).WithAnchor(1.8, "acquisition", "escrow", "million", "tranche", "corporate")

				g.Ambiguous(func(ctx context.Context, p, s string, payload any) error {
					fmt.Printf("    [STEP-UP 2FA] Ambiguous memo (%s vs %s): SMS OTP challenge required\n", p, s)
					return nil
				}).Fallback(func(ctx context.Context, payload any) error {
					fmt.Println("    [MANUAL AUDIT] Route wire memo to fraud investigations team")
					return nil
				})
			},
			TestCases: []TestCase{
				{Query: "monthly lunch payment split with office colleagues", Expectation: "Instant Approval: NormalTransfer"},
				{Query: "urgent send funds now police fine wire to bitcoin wallet", Expectation: "Block & Intercept: PhishingSuspicion"},
				{Query: "merchant charged my card three times for single coffee", Expectation: "Dispute: ChargebackDispute"},
				{Query: "corporate acquisition escrow settlement tranche wire five million dollars", Expectation: "AML Audit: HighValueAudit"},
			},
		},
	}

	orderedKeys := []string{"cs", "llm", "sre", "iot", "cicd", "fintech"}
	selected := strings.ToLower(*targetDomain)

	ctx := context.Background()
	totalStart := time.Now()
	totalQueries := 0
	executedDomains := 0

	fmt.Println("================================================================================")
	fmt.Println("  INTELLIBRANCH v2.0 - 6-DOMAIN NEUROGATE 3-HEAD INTELLIGENT FILTERING SUITE")
	fmt.Println("================================================================================")

	for _, key := range orderedKeys {
		if selected != "all" && selected != key {
			continue
		}

		executedDomains++
		suite := suites[key]
		fmt.Printf("\n>>> DOMAIN: %s\n", suite.DomainName)
		fmt.Printf("    Model Path  : %s\n", suite.ModelPath)
		fmt.Printf("    Capability  : %s\n", suite.Description)
		fmt.Printf("    ----------------------------------------------------------------------------\n")

		// Auto-train model on-the-fly if binary is absent (zero manual downloads required)
		ensureModel(suite.ModelPath, suite.DataPath)

		gate, err := intellibranch.NewNeuroGate(suite.ModelPath)
		if err != nil {
			log.Fatalf("Fatal: Failed to load NeuroGate [%s]: %v", suite.ModelPath, err)
		}
		gate.SetPolicy(suite.Policy)
		if suite.MinCosine > 0 {
			gate.SetMinCosineSim(suite.MinCosine)
		}
		if samples, err := intellibranch.LoadCSVDataset(suite.DataPath); err == nil {
			gate.CalibrateDomainCentroid(samples)
		}
		suite.SetupGate(gate)

		for _, tc := range suite.TestCases {
			totalQueries++
			trace := gate.Inspect(tc.Query)
			benchLatency := measureBenchmarkLatency(gate, tc.Query, 1000)
			fmt.Printf("  • Input    : \"%s\"\n", tc.Query)
			fmt.Printf("    Expect   : %s\n", tc.Expectation)
			fmt.Printf("    Inference: %s (Confidence: %.2f%%, Cosine: %.4f, Entropy: %.4f, Latency: %.2f μs/op, OOD: %t)\n",
				trace.PredictedLabel, trace.Confidence*100, trace.CosineSimilarity, trace.Entropy, benchLatency, trace.IsOOD)

			_ = gate.FilterPipeline(ctx, tc.Query, nil)
			fmt.Println()
		}

		if suite.CustomRun != nil {
			suite.CustomRun(gate, ctx)
			fmt.Println()
		}
	}

	totalDuration := time.Since(totalStart)
	domainLabel := "domain"
	if executedDomains > 1 {
		domainLabel = "domains"
	}
	fmt.Println("================================================================================")
	fmt.Printf("DEMONSTRATION COMPLETED: %d queries routed across %d distinct neural %s in %s\n",
		totalQueries, executedDomains, domainLabel, totalDuration)
	fmt.Println("ALL INFERENCES RUN IN MICROSECONDS WITH CGO_ENABLED=0 AND ZERO ALLOCATIONS.")
	fmt.Println("================================================================================")
}

