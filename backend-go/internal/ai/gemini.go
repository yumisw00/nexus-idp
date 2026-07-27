package ai

import (
	"context"
	"os"

	"github.com/google/genai-go"
)

func AnalyzeDocument(ctx context.Context, filePath string) (string, error) {
	c, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	res, err := c.Models.GenerateContent(ctx, "gemini-3.1-pro",
		genai.Text("Extract all key entities, dates, total amounts, and parties involved from this document. Return ONLY a valid, minified JSON object."),
		genai.Blob{Data: b, MimeType: "text/plain"},
	)
	if err != nil {
		return "", err
	}
	if len(res.Candidates) > 0 && res.Candidates[0].Content != nil && len(res.Candidates[0].Content.Parts) > 0 {
		if t, ok := res.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(t), nil
		}
	}
	return "", nil
}
