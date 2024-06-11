package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Payload struct {
	Inputs string `json:"inputs"`
}

func main() {
	fmt.Println("AI MODEL APP")
	godotenv.Load()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/mask", handlerMaskModel)
	r.POST("/translate", handlerTranslateModel)
	r.Run()
}

func handlerMaskModel(c *gin.Context) {
	// Read token in .env file
	token := os.Getenv("HUGGINGFACE_TOKEN")

	data := Payload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while encoding payload": err.Error()})
		return
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/bert-base-uncased", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while creating request:": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token+"") // Ganti "YOUR_HUGGING_FACE_TOKEN" dengan token Anda.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while making request:": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while reading response body:": err.Error()})
		return
	}

	fmt.Println("Logs: ", string(respBody))

	var responseJSON []map[string]interface{}
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while unmarshalling response body:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseJSON)
}

func handlerTranslateModel(c *gin.Context) {
	// Read token in .env file
	token := os.Getenv("HUGGINGFACE_TOKEN")

	data := Payload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while encoding payload": err.Error()})
		return
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/Helsinki-NLP/opus-mt-id-en", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while creating request:": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token+"") // Ganti "YOUR_HUGGING_FACE_TOKEN" dengan token Anda.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while making request:": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while reading response body:": err.Error()})
		return
	}

	fmt.Println("Logs: ", string(respBody))

	var responseJSON []map[string]interface{}
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while unmarshalling response body:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseJSON)
}
