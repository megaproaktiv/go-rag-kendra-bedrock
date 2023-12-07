package chain

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	bedrock "github.com/megaproaktiv/go-rag-kendra-bedrock/bedrock"
	kendra "github.com/megaproaktiv/go-rag-kendra-bedrock/kendra"
	"golang.org/x/exp/slog"
)

type TemplateData struct {
	Question string
	Document string
}

func RagChain(question kendra.Query) (string, *[]kendra.Document, error) {
	slog.Info("Lambda start")
	slog.Info("Kendra start")
	documents := make([]kendra.Document, 0)
	query, err := kendra.Retrieve(kendra.Client, question)
	if err != nil {
		slog.Error("Kendra retrieve error", "Error", err)
	}
	slog.Info("Kendra end")

	// Build the prompt
	// from Instructions, Question, and Documents
	promptTemplate, err := os.ReadFile("prompt.tmpl")
	var templateStr string
	if err != nil {
		log.Println("Error reading template, using standard:", err)
		templateStr = ` This is a friendly conversation between a human and an AI.
		The AI is conversational and provides many specific details from its context.
		If the AI does not know the answer to a question, it truthfully says that it
		does not know.
		Instruction: You are a friendly service guy.
	  Based on this text, give a detailed answer to the following question:
				{{.Question}}
		Answers with "I can't say anything about that",
		if the data in the document is not sufficient.
		<documents>
		{{.Document}}
		</documents>
	`
	}
	templateStr = string(promptTemplate)

	preExcerpt := "<document>\n"
	postExcerpt := "</document>\n"

	documentExcerpts := ""
	for _, doc := range query.ResultItems {
		document := kendra.Document{
			Excerpt: doc.Content,
			Title:   doc.DocumentTitle,
		}
		document.Page = kendraPage(doc)
		documentExcerpts += preExcerpt
		documentExcerpts += " Document Title: " + *doc.DocumentTitle + "\n"
		documentExcerpts += " Document Excerpt: " + *doc.Content
		documentExcerpts += postExcerpt
		documents = append(documents, document)
	}

	tmpl, err := template.New("Prompt").Parse(templateStr)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	data := TemplateData{
		Question: question.Question,
		Document: documentExcerpts,
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	// Extract the string from the buffer
	prompt := buffer.String()

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
