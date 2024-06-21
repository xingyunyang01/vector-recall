package aihelper

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/xingyunyang01/vector-recall/pkg/httphelper"
)

func NewOpenAiClient() *openai.Client {
	token := os.Getenv("DashScope")
	//dashscope_url := "https://dashscope.aliyuncs.com/compatible-mode/v1"
	dashscope_url := "http://8.136.125.54/compatible-mode/v1"
	//dashscope_url := "https://api.moonshot.cn/v1"
	config := openai.DefaultConfig(token)
	config.BaseURL = dashscope_url

	return openai.NewClientWithConfig(config)
}

// 文本转为向量
func SimpleGetVec(prompt string) ([]float64, error) {
	textInputs := []string{prompt}
	req := &Request{
		Model: TextEmbeddingV2,
		Params: Params{
			TextType: TypeDocument, // 默认值
		},
		Input: Input{
			Texts: textInputs,
		},
	}
	resp, err := createEmbedding(req, httphelper.NewHTTPClient())
	if err != nil {
		return nil, err
	}

	var embedding []float64

	if len(resp.Output.Embeddings) > 0 {
		embedding = resp.Output.Embeddings[0].Embedding
	}

	return embedding, nil
}

func GetVec(textInputs []string) ([]float64, error) {
	//textInputs := []string{question, anwser}
	req := &Request{
		Model: TextEmbeddingV2,
		Params: Params{
			TextType: TypeDocument, // 默认值
		},
		Input: Input{
			Texts: textInputs,
		},
	}
	resp, err := createEmbedding(req, httphelper.NewHTTPClient())
	if err != nil {
		return nil, err
	}

	var embedding []float64

	if len(resp.Output.Embeddings) > 0 {
		embedding = resp.Output.Embeddings[0].Embedding
	}
	return embedding, nil
}

func createEmbedding(req *Request, cli httphelper.IHttpClient) (*Response, error) {
	token := os.Getenv("DashScope")

	if req.Model == "" {
		req.Model = TextEmbeddingV2
	}
	if req.Params.TextType == "" {
		req.Params.TextType = TypeDocument
	}

	resp := Response{}
	tokenOption := httphelper.WithTokenHeaderOption(token)
	headerOption := httphelper.WithHeader(httphelper.HeaderMap{"content-type": "application/json"})
	err := cli.Post(context.Background(), embeddingURL, req, &resp, tokenOption, headerOption)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
