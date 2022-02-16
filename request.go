package lingvanex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	platform = "api"
)

func (c *client) GetLanguages(ctx context.Context, code string) ([]Language, error) {
	params := url.Values{}
	params.Add("platform", platform)
	if len(code) > 0 {
		params.Add("code", code)
	}

	reqURL := fmt.Sprintf("%s?%s", c.languagesURL, params.Encode())
	req, err := c.getRequest(ctx, reqURL)
	if err != nil {
		return nil, err
	}

	body, err := c.parseBody(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp LanguagesResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Err) > 0 {
		return nil, errors.New(resp.Err)
	}

	return resp.Result, nil
}

func (c *client) Translate(ctx context.Context, q, source, target string) (*TranslateResponse, error) {
	params := TranslateRequest{
		Data:            q,
		From:            source,
		To:              target,
		Platform:        platform,
		Transliteration: c.transliteration,
	}

	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	req, err := c.postRequest(ctx, reqBody, c.translateURL)
	if err != nil {
		return nil, err
	}

	body, err := c.parseBody(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp TranslateResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Err) > 0 {
		return nil, errors.New(resp.Err)
	}

	return &resp, nil
}

func (c *client) getRequest(ctx context.Context, reqURL string) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", c.apiKey)
	return r, nil
}

func (c *client) postRequest(ctx context.Context, body []byte, reqURL string) (*http.Request, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", c.apiKey)
	return r, nil
}

func (c *client) parseBody(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	var rReq = new(retryablehttp.Request)
	rReq.Request = req
	resp, err := c.httpClient.Do(rReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
