package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func AnalyzeDocument(ctx context.Context, filePath string) (string, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := "text/plain"
	if strings.HasSuffix(filePath, ".pdf") {
		mimeType = "application/pdf"
	}

	docPart := &genai.Part{
		InlineData: &genai.Blob{
			Data:     data,
			MIMEType: mimeType,
		},
	}

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: "Extract all key entities, dates, total amounts, and parties involved from this document. Return ONLY a valid JSON object."},
				docPart,
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash-lite", contents, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return resp.Text(), nil
}
