package aihelper

type Request struct {
	Model  string `json:"model"`
	Input  Input  `json:"input"`
	Params Params `json:"parameters"`
}

type Input struct {
	Texts []string `json:"texts"`
}

type Params struct {
	TextType string `json:"text_type"` // query or document
}

type Response struct {
	Output Output `json:"output"`
	Usgae  struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

type Embeddings struct {
	TextIndex int       `json:"text_index"`
	Embedding []float64 `json:"embedding"`
}

type Output struct {
	Embeddings []Embeddings `json:"embeddings"`
}

type OutPutFormat struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
