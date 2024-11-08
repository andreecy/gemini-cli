package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
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

type Prompt struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

func convertRespToPrompts(resp *genai.GenerateContentResponse) []Prompt {
	var prompts []Prompt
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				switch v := part.(type) {
				case genai.Text:
					// check if text start with "{"
					// if yes, try to unmarshal it as JSON object
					if strings.HasPrefix(string(v), "{") {
						// if yes, try to unmarshal it as JSON
						var prompt Prompt
						err := json.Unmarshal([]byte(v), &prompt)
						if err != nil {
							log.Fatal(err)
						}
						prompts = append(prompts, prompt)
					} else if strings.HasPrefix(string(v), "[") {
						// try to unmarshal it as JSON array
						err := json.Unmarshal([]byte(v), &prompts)
						if err != nil {
							log.Fatal(err)
						}
					}
				default:
					fmt.Println("Unknown part type:", v)
				}
			}
		}
	}
	return prompts
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

func howto(prompt string) {
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

	prompts := convertRespToPrompts(resp)

	fmt.Printf("──────────────── %s ────────────────\n", "Command")
	p := prompts[0]
	fmt.Println(p.Command)
	fmt.Printf("\n")

	selectAction := promptui.Select{
		Label: "Select Action",
		Items: []string{"✅ Run this command", "❌ Cancel"},
	}

	actionIdx, _, err := selectAction.Run()

	if err != nil {
		fmt.Printf("Prompt cancelled %v\n", err)
		return
	}

	switch actionIdx {
	case 0:
		// Run the command
		fmt.Println("Running command...")
		cmd := exec.Command("/bin/bash", "-c", p.Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error running command:", err)
		}
		break
	case 1:
		// Cancel
		fmt.Println("Cancelled")
		break
	}
}
