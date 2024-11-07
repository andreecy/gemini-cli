package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}

func generateContent(prompt string) {
	ctx := context.Background()

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)
}

func startChat(prompt string) {
	ctx := context.Background()

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	finalPrompt := fmt.Sprintf(`%s. Respond using this JSON schema:
  Prompt = {'command': string, description: string}
	           Return: Array<Prompt>`, prompt)

	resp, err := model.GenerateContent(ctx, genai.Text(finalPrompt))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)
}
