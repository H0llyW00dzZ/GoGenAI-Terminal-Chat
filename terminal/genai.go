package terminal

import (
	"context"
	"fmt"
	"time"

	genai "github.com/google/generative-ai-go/genai"
)

func PrintTypingChat(message string, delay time.Duration) {
	for _, char := range message {
		fmt.Printf(AnimatedChars, char)
		time.Sleep(delay)
	}
	fmt.Println()
}

func SendMessage(ctx context.Context, client *genai.Client, chatContext string) (string, error) {
	// this subject to changed if there is lots of models
	model := client.GenerativeModel(ModelAi)
	cs := model.StartChat()

	resp, err := cs.SendMessage(ctx, genai.Text(chatContext))
	if err != nil {
		return "", err
	}

	return printResponse(resp), nil
}

func printResponse(resp *genai.GenerateContentResponse) string {
	aiResponse := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// Print "AI:" prefix directly without typing effect
				fmt.Print(AiNerd)

				// Assuming 'part' can be printed directly and is of type string or has a String() method
				// Use the typing banner effect only for the part content
				PrintTypingChat(fmt.Sprint(part), 100*time.Millisecond)
				aiResponse += fmt.Sprint(part) // Collect AI response
			}
		}
	}
	fmt.Println(StripChars)
	return aiResponse
}
