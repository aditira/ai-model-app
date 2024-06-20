package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ai-model-app/model"
)

type Service interface {
	RequestHuggingface(payload model.Payload, modelName string) ([]byte, error)
}

type service struct {
	token string
}

// requestHuggingface implements Service.
func (s *service) RequestHuggingface(payload model.Payload, modelName string) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/"+modelName, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token) // Ganti "YOUR_HUGGING_FACE_TOKEN" dengan token Anda.

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func NewService(token string) Service {
	return &service{
		token: token,
	}
}
