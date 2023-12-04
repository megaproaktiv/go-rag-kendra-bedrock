package chain

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	bedrock "github.com/megaproaktiv/go-rag-kendra-bedrock/bedrock"
	kendra "github.com/megaproaktiv/go-rag-kendra-bedrock/kendra"
	"golang.org/x/exp/slog"
)

func RagChain(question kendra.Query) (string, *[]kendra.Document, error) {
	slog.Info("Lambda start")
	slog.Info("Kendra start")
	documents := make([]kendra.Document, 0)
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
	 Based on this text, give a detailed answer to the following question: \n` + question.Question + ` answers with "I can't say anything about that",
			if the data in the document is not sufficient.
	`

	prompt := pre
	for _, doc := range query.ResultItems {
		document := kendra.Document{
			Excerpt: doc.Content,
			Title:   doc.DocumentTitle,
		}
		document.Page = kendraPage(doc)
		prompt = prompt + " Document Title: " + *doc.DocumentTitle
		prompt = prompt + " Document Excerpt: " + *doc.Content
		documents = append(documents, document)
	}
	prompt = prompt + post

	answer := bedrock.Chat(prompt)

	slog.Info("OpenAI end")
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return "", nil, err
	}

	fmt.Println(answer)
	slog.Info("Lambda end")
	return answer, &documents, nil
}

func kendraPage(item types.RetrieveResultItem) *int {
	page := 0
	for _, attribute := range item.DocumentAttributes {
		if *attribute.Key == "_excerpt_page_number" {
			page = int(*attribute.Value.LongValue)
		}
	}
	return &page
}

func kendraCategory(item types.RetrieveResultItem) *string {
	var category *string
	for _, attribute := range item.DocumentAttributes {
		if *attribute.Key == "_category" {
			category = attribute.Value.StringValue
		}
	}
	return category
}

func kendraVersion(item types.RetrieveResultItem) *string {
	var version *string
	for _, attribute := range item.DocumentAttributes {
		if *attribute.Key == "_version" {
			version = attribute.Value.StringValue
		}
	}
	return version
}
