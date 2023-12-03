package chain

import (
	"fmt"

	bedrock "github.com/megaproaktiv/go-rag-kendra-bedrock/bedrock"
	kendra "github.com/megaproaktiv/go-rag-kendra-bedrock/kendra"
	"golang.org/x/exp/slog"
)

func RagChain(question string) (string, error) {
	slog.Info("Lambda start")
	slog.Info("Kendra start")
	query, err := kendra.Retrieve(kendra.Client, question)
	if err != nil {
		slog.Error("Kendra retrieve error", "Error", err)
	}
	slog.Info("Kendra end")

	pre := ` This is a friendly conversation between a human and an AI.
		The AI is conversational and provides many specific details from its context.
		If the AI does not know the answer to a question, it truthfully says that it
		does not know.`

	post := `Instruction: You are a friendly service guy.
		Based on this text, give a detailed answer to the following question: \n` + question + ` answers with "I can't say anything about that",
			if the data in the document is not sufficient.
	`

	prompt := pre
	for _, doc := range query.ResultItems {
		prompt = prompt + " Document Title: " + *doc.DocumentTitle
		prompt = prompt + " Document Excerpt: " + *doc.Content
	}
	prompt = prompt + post

	answer := bedrock.Chat(prompt)

	slog.Info("OpenAI end")
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return "", err
	}

	fmt.Println(answer)
	slog.Info("Lambda end")
	return answer, nil
}
