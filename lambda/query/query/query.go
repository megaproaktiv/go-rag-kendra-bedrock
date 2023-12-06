package query

import (
	"golang.org/x/exp/slog"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/megaproaktiv/go-rag-kendra-bedrock/chain"
	"github.com/megaproaktiv/go-rag-kendra-bedrock/kendra"
)

type QueryRequest struct {
	Question string `json:"question"`
	Category string `json:"category"`
	Version  string `json:"version"`
}

type Response struct {
	Answer    string            `json:"answer"`
	Documents []kendra.Document `json:"documents"`
}

func Query(c *gin.Context) {
	var req QueryRequest

	err := c.BindJSON(&req)
	if err != nil {
		slog.Error("Error loading input parameter", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	question := req.Question
	log.Println("Question received", question)
	log.Println("Category", req.Category)
	log.Println("Version", req.Version)

	query := kendra.Query{
		Question: question,
	}

	if req.Category != "" {
		query.Category = &req.Category
	}
	if req.Version != "" {
		query.Version = &req.Version
	}
	answer, documents, err := chain.RagChain(query)
	if err != nil {
		slog.Error("Error chain rag", "error", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := Response{
		Answer:    answer,
		Documents: *documents,
	}
	c.JSON(200, response)
}
