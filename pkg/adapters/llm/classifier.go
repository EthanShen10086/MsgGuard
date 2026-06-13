package llmadapter

import (
	"context"
	"strings"

	vllm "github.com/EthanShen10086/voxera-kit/llm"
	"github.com/EthanShen10086/voxera-kit/llm/prompt"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type Classifier struct {
	router *vllm.Router
}

func NewClassifier(router *vllm.Router) *Classifier {
	return &Classifier{router: router}
}

func (c *Classifier) Classify(ctx context.Context, text string) (string, error) {
	if c.router == nil {
		return "ham", nil
	}
	tmpl := prompt.Classify
	_, userPrompt := tmpl.Render(map[string]any{
		"Categories": "spam, promotion, ham, phishing, transaction",
		"Text":       text,
	})
	resp, err := c.router.Route(ctx, vllm.Request{
		Messages: []vllm.Message{
			{Role: vllm.RoleSystem, Content: tmpl.System},
			{Role: vllm.RoleUser, Content: userPrompt},
		},
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.ToLower(resp.Content)), nil
}

var _ ports.LLMClassifier = (*Classifier)(nil)
