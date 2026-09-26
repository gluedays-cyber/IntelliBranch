package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"intellibranch/pkg/intellibranch"
)

func main() {
	datasetPath := flag.String("data", "data/sample_dataset.csv", "Path to CSV training dataset")
	outputPath := flag.String("out", "weights/model.bin", "Output binary model path")
	epochs := flag.Int("epochs", 150, "Maximum training epochs")
	targetVocab := flag.Int("vocab", 150, "Target BPE vocabulary size")
	lr := flag.Float64("lr", 0.005, "Learning rate")
	flag.Parse()

	log.Printf("Loading dataset from: %s", *datasetPath)
	samples, err := intellibranch.LoadCSVDataset(*datasetPath)
	if err != nil {
		log.Fatalf("Failed to load dataset: %v", err)
	}
	log.Printf("Loaded %d training samples", len(samples))

	cfg := intellibranch.DefaultTrainConfig()
	cfg.Epochs = *epochs
	cfg.TargetVocabSize = *targetVocab
	cfg.LearningRate = float32(*lr)
	cfg.BatchSize = 16
	cfg.Patience = 10

	log.Println("Starting offline BPE + AdamW training pipeline...")
	model, err := intellibranch.TrainModel(samples, cfg)
	if err != nil {
		log.Fatalf("Training failed: %v", err)
	}

	outDir := filepath.Dir(*outputPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	log.Printf("Serializing trained model to Little-Endian binary: %s", *outputPath)
	if err := intellibranch.SaveBinaryModel(*outputPath, model); err != nil {
		log.Fatalf("Failed to save model: %v", err)
	}

	log.Println("Training and binary export completed successfully.")
	fmt.Printf("Model saved at: %s (Vocab: %d, Classes: %d)\n", *outputPath, model.Header.VocabSize, model.Header.NumClasses)
}
