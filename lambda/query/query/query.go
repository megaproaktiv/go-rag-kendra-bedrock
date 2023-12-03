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

	// answer := bedrock.Chat(question)
	answer, documents, err := chain.RagChain(question)
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
