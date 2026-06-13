package ports

import "context"

// LLMClassifier classifies text via cloud LLM (L3).
type LLMClassifier interface {
	Classify(ctx context.Context, text string) (category string, err error)
}
