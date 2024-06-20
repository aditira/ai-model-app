package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/ai-model-app/model"
	"github.com/ai-model-app/model/dto"
	"github.com/ai-model-app/service"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	MaskModel(c *gin.Context)
	TranslateModel(c *gin.Context)
	ChatModel(c *gin.Context)
}

type handler struct {
	t       *template.Template
	service service.Service
}

func (h *handler) ChatModel(c *gin.Context) {
	data := model.Payload{
		Inputs: c.Query("message"),
	}

	template := "chat.html"

	if data.Inputs == "" {
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  "message cannot be empty",
		})
		return
	}

	respBody, err := h.service.RequestHuggingface(data, "HuggingFaceH4/zephyr-7b-beta")
	if err != nil {
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  err.Error(),
		})
		return
	}

	var responseJSON []dto.Chat
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		var responseError dto.Error
		err = json.Unmarshal(respBody, &responseError)
		if err != nil {
			h.t.ExecuteTemplate(c.Writer, template, gin.H{
				"Data": nil,
				"Err":  err.Error(),
			})
			return
		}
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  responseError.Error,
		})
		return
	}

	h.t.ExecuteTemplate(c.Writer, template, gin.H{
		"Data": responseJSON,
		"Err":  nil,
	})
}

func (h *handler) TranslateModel(c *gin.Context) {
	data := model.Payload{
		Inputs: c.Query("message"),
	}

	template := "translate.html"

	if data.Inputs == "" {
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  "message cannot be empty",
		})
		return
	}

	respBody, err := h.service.RequestHuggingface(data, "Helsinki-NLP/opus-mt-id-en")
	if err != nil {
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  err.Error(),
		})
		return
	}

	var responseJSON []dto.Translation
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		var responseError dto.Error
		err = json.Unmarshal(respBody, &responseError)
		if err != nil {
			h.t.ExecuteTemplate(c.Writer, template, gin.H{
				"Data": nil,
				"Err":  err.Error(),
			})
			return
		}
		h.t.ExecuteTemplate(c.Writer, template, gin.H{
			"Data": nil,
			"Err":  responseError.Error,
		})
		return
	}

	h.t.ExecuteTemplate(c.Writer, template, gin.H{
		"Data": responseJSON,
		"Err":  nil,
	})
}

func (h *handler) MaskModel(c *gin.Context) {
	data := model.Payload{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	respBody, err := h.service.RequestHuggingface(data, "bert-base-uncased")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseJSON []map[string]interface{}
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseJSON)
}

func NewHandler(t *template.Template, service service.Service) Handler {
	return &handler{
		t:       t,
		service: service,
	}
}
