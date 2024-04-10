package server

import (
	"encoding/base64"
	"go-source/pkg/parser"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChunkRequest struct {
	Filename string `json:"filename" binding:"required"`
	Input    string `json:"input" binding:"required"`
}

func Run() error {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.POST("/chunk", func(c *gin.Context) {
		var json ChunkRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input, err := base64.StdEncoding.DecodeString(json.Input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chunker, err := parser.DetectLanguage(json.Filename, input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := parser.Parse(c.Request.Context(), chunker, input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result.File = json.Filename
		c.JSON(http.StatusOK, result)

	})
	return r.Run()
}
