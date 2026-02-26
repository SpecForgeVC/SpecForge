package llm

import (
	"context"
	"errors"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) StreamGenerate(ctx context.Context, prompt string, chunks chan<- string) error {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		if len(response.Choices) > 0 {
			chunks <- response.Choices[0].Delta.Content
		}
	}
}

func (c *OpenAIClient) TestConnection(ctx context.Context) error {
	_, err := c.Generate(ctx, "ping")
	return err
}

func (c *OpenAIClient) ListModels(ctx context.Context) ([]string, error) {
	modelsList, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, m := range modelsList.Models {
		models = append(models, m.ID)
	}
	return models, nil
}
