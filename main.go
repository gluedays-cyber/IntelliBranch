package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"intellibranch/pkg/intellibranch"
)

// 1. Business Logic Handlers
func handleRefund(ctx context.Context, payload any) error {
	fmt.Printf("[ACTION: Refund]   Processing refund for: '%v'\n", payload)
	return nil
}

func handleDelivery(ctx context.Context, payload any) error {
	fmt.Printf("[ACTION: Delivery] Querying shipment tracking for: '%v'\n", payload)
	return nil
}

func handleAccount(ctx context.Context, payload any) error {
	fmt.Printf("[ACTION: Account]  Initiating account security for: '%v'\n", payload)
	return nil
}

func handleFallback(ctx context.Context, payload any) error {
	fmt.Printf("[FALLBACK: Safety] Isolated low-confidence request: '%v'\n", payload)
	return nil
}

func main() {
	modelPath := "weights/intent.bin"

	// Auto-compile model if missing (ensures instant zero-config clone & run)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Println("Model weights not found. Compiling from data/sample_dataset.csv...")
		samples, err := intellibranch.LoadCSVDataset("data/sample_dataset.csv")
		if err != nil {
			log.Fatalf("Failed to load dataset: %v", err)
		}
		cfg := intellibranch.DefaultTrainConfig()
		cfg.Epochs = 50
		cfg.LearningRate = 0.005
		cfg.TargetVocabSize = 250

		model, err := intellibranch.TrainModel(samples, cfg)
		if err != nil {
			log.Fatalf("Training failed: %v", err)
		}
		_ = os.MkdirAll("weights", 0755)
		if err := intellibranch.SaveBinaryModel(modelPath, model); err != nil {
			log.Fatalf("Failed to save model: %v", err)
		}
		log.Println("Model compilation completed.")
	}

	// 2. Load compiled binary weights into memory (0.60 calibrated threshold)
	router, err := intellibranch.NewRouter(modelPath, 0.60)
	if err != nil {
		log.Fatalf("Router initialization failed: %v", err)
	}

	// 3. Bind routes directly inside main.go
	router.
		Bind("Refund", handleRefund).
		Bind("Delivery", handleDelivery).
		Bind("Account", handleAccount).
		Fallback(handleFallback)

	// 4. Execute microsecond branch dispatch
	testQueries := []string{
		"I want to cancel my payment and request a refund",
		"When will my delivery package arrive",
		"Forgot my account password",
		"Please refund my purchase",
		"Track my shipment status",
		"Completely random gibberish noise 12345!@#$",
		"hey where is my stuff it was supposed to get here yesterday",
		"can u cancel order #49281? i bought it by mistake",
		"bruh the reset link is not sending to my email, fix this",
		"got charged twice on my card, refund the extra charge asap",
		"item arrived totally smashed, want my money back",
		"cant log into my acct keeps saying wrong password",
		"tracking says delivered but nothing is in my mailbox",
		"yo i typed the wrong apt number, can someone update the address before it ships",
		"sent the return box a week ago, when do i get my refund?",
		"locked out of my account after 3 tries... help pls",
		"ordered a large but you guys sent me a small",
		"any update on order #88412? hasnt moved in 4 days",
		"how do i just delete my account permanently? done with this site",
		"driver dumped the package in the rain, everything inside is ruined",
		"promo code didnt apply at checkout, can u refund the difference",
		"need a real person, this bot is completely useless",
		"can i change the delivery date? nobody will be home this friday",
		"my card was charged but never received any confirmation email or receipt",
		"lost access to my 2FA phone number, how do i get back in",
		"package has been stuck in transit for 10 days straight, is it lost or what",
	}

	fmt.Println("=== IntelliBranch Server Routing Started ===")
	ctx := context.Background()
	for _, query := range testQueries {
		if err := router.Dispatch(ctx, query, query); err != nil {
			log.Printf("Dispatch error: %v", err)
		}
	}
	fmt.Println("=== All queries dispatched in microseconds ===")
}
