package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "Ask a question, such as 'How do I deploy and AKS cluster with an API gateway?'",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]

		fmt.Printf("You asked: %s\n\n", prompt)
		fmt.Print("Here are some suggested documents:\n\n")

		suggestedDocuments := getSuggestedDocuments(prompt)
		for i, doc := range suggestedDocuments {
			fmt.Printf("\t%d. %s\n", i+1, doc)
		}
		fmt.Println()

		requestAction(suggestedDocuments)
	},
}

func getSuggestedDocuments(prompt string) []string {
	results := []string{}

	azureOpenAIKey := os.Getenv("OPENAI_API_KEY")
	if azureOpenAIKey == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY environment variable not set\n")
		return results
	}
	azureOpenAIEndpoint := os.Getenv("OPENAI_ENDPOINT")
	if azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_ENDPOINT environment variable not set\n")
		return results
	}
	modelDeploymentID := os.Getenv("OPENAI_MODEL_DEPLOYMENT_NAME")
	if modelDeploymentID == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_MODEL_DEPLOYMENT_NAME environment variable not set\n")
		return results
	}

	maxTokens := int32(400)
	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		log.Printf("ERROR creating OpenAI Client: %s", err)
		return results
	}

	// This is a conversation in progress.
	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: azopenai.NewChatRequestSystemMessageContent("You are a helpful assistant.")},

		// The user asks a question
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("Does Azure OpenAI support customer managed keys?")},

		// The reply would come back from the model. You'd add it to the conversation so we can maintain context.
		&azopenai.ChatRequestAssistantMessage{Content: azopenai.NewChatRequestAssistantMessageContent("Yes, customer managed keys are supported by Azure OpenAI")},

		// The user answers the question based on the latest reply.
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("What other Azure Services support customer managed keys?")},

		// from here you'd keep iterating, sending responses back from ChatGPT
	}

	gotReply := false

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// This is a conversation in progress.
		// NOTE: all messages count against token usage for this API.
		Messages:       messages,
		DeploymentName: &modelDeploymentID,
		MaxTokens:      &maxTokens,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return results
	}

	for _, choice := range resp.Choices {
		gotReply = true

		if choice.ContentFilterResults != nil {
			fmt.Fprintf(os.Stderr, "Content filter results\n")

			if choice.ContentFilterResults.Error != nil {
				fmt.Fprintf(os.Stderr, "  Error:%v\n", choice.ContentFilterResults.Error)
			}

			fmt.Fprintf(os.Stderr, "  Hate: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Hate.Severity, *choice.ContentFilterResults.Hate.Filtered)
			fmt.Fprintf(os.Stderr, "  SelfHarm: sev: %v, filtered: %v\n", *choice.ContentFilterResults.SelfHarm.Severity, *choice.ContentFilterResults.SelfHarm.Filtered)
			fmt.Fprintf(os.Stderr, "  Sexual: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Sexual.Severity, *choice.ContentFilterResults.Sexual.Filtered)
			fmt.Fprintf(os.Stderr, "  Violence: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Violence.Severity, *choice.ContentFilterResults.Violence.Filtered)
		}

		if choice.Message != nil && choice.Message.Content != nil {
			fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
		}

		if choice.FinishReason != nil {
			// this choice's conversation is complete.
			fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", *choice.Index, *choice.FinishReason)
		}
	}

	if gotReply {
		fmt.Fprintf(os.Stderr, "Received chat completions reply\n")
	}

	return results
}

func requestAction(suggestedDocuments []string) {
	fmt.Println("Enter the ID number to view more information or enter `Q[uit]` to exit.")

	var input string
	for {
		fmt.Print("Your choice: ")
		fmt.Scanln(&input)
		if input == "Q" || input == "Quit" || input == "q" || input == "quit" {
			fmt.Println("Exiting...")
			return
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(suggestedDocuments) {
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		fmt.Printf("Details about %s...\n", suggestedDocuments[choice-1])
	}
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
