package service

import (
	"bytes"
	"io"
	"net/http"
)

type Service interface {
	RequestHuggingface(payload []byte, modelName string) ([]byte, error)
}

type service struct {
	token string
}

// requestHuggingface implements Service.
func (s *service) RequestHuggingface(payload []byte, modelName string) ([]byte, error) {
	body := bytes.NewReader(payload)

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
