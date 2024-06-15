package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ai-model-app/model"
	"github.com/ai-model-app/model/dto"
	"github.com/ai-model-app/service"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	MaskModel(c *gin.Context)
	TranslateModel(c *gin.Context)
}

type handler struct {
	service service.Service
}

// handlerTranslateModel implements Handler.
func (h *handler) TranslateModel(c *gin.Context) {
	data := model.Payload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while encoding payload": err.Error()})
		return
	}

	respBody, err := h.service.RequestHuggingface(payloadBytes, "Helsinki-NLP/opus-mt-id-en")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error internal server:": err.Error()})
		return
	}

	fmt.Println("Resp: ", string(respBody))

	var responseJSON []dto.Translation
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while unmarshalling response body:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseJSON)
}

func (h *handler) MaskModel(c *gin.Context) {
	data := model.Payload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while encoding payload": err.Error()})
		return
	}

	// Interaksi ke Huggingface via Service
	respBody, err := h.service.RequestHuggingface(payloadBytes, "bert-base-uncased")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error internal server:": err.Error()})
		return
	}

	var responseJSON []map[string]interface{}
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error occurred while unmarshalling response body:": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseJSON)
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}
