package httpclient

import (
	"context"
	"io"
	"net/http"

	"github.com/BleSSSeDDD/reviewer-assignment/internal/domain"
)

type vacancyParser struct {
	client *http.Client
}

func NewParser() domain.Parser {
	return &vacancyParser{client: &http.Client{}}
}

func (p *vacancyParser) GetRawContent(ctx context.Context, url string) (string, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

// func (p *vacancyParser) GetRawContent(ctx context.Context, url string) (string, error) {
// 	if url == "ivan sigma" {
// 		return string("qwe"), nil
// 	}
// 	return string("da eto tak)"), nil
// }
