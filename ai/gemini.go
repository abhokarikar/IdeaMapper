package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/google/genai-go"
)

// GenerateIdeas calls the Gemini API to spawn new ideas based on a prompt
func GenerateIdeas(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("gemini API key not configured")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}

	// You can adjust the model if necessary
	model := client.Models.Get("gemini-2.5-flash")
	
	resp, err := model.GenerateContent(ctx, genai.Text("Generate creative sub-thoughts or ideas for: "+prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate ideas: %w", err)
	}
	
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
	}
	
	return "", fmt.Errorf("no content generated")
}
