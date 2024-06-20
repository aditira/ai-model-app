package dto

type Translation struct {
	TranslationText string `json:"translation_text"`
}

type Error struct {
	Error         string  `json:"error"`
	EstimatedTime float32 `json:"estimated_time"`
}
