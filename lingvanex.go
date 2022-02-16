package lingvanex

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/powerman/structlog"
)

const (
	defaultRetryMax    = 5
	defaultConnTimeout = 15 * time.Second
)

type Client interface {
	SetRetryMax(retryMax int) Client
	SetConnTimeout(connTimeout time.Duration) Client
	SetURL(apiURL string) Client
	SetKey(apiKey string) Client
	SetEnableTransliteration(enable bool) Client
	GetLanguages(ctx context.Context, code string) ([]Language, error)
	Translate(ctx context.Context, q, source, target string) (*TranslateResponse, error)
}

type client struct {
	transliteration bool
	retryMax        int
	apiURL          string
	apiKey          string

	languagesURL string
	translateURL string

	connTimeout time.Duration

	httpClient *retryablehttp.Client
	logger     *structlog.Logger
}

func NewTranslator(apiURL string) Client {
	c := &client{
		apiURL:       apiURL,
		retryMax:     defaultRetryMax,
		connTimeout:  defaultConnTimeout,
		languagesURL: fmt.Sprintf("%s/getLanguages", apiURL),
		translateURL: fmt.Sprintf("%s/translate", apiURL),
		logger:       structlog.New(),
	}
	c.initClient()
	return c
}

func (c *client) SetRetryMax(retryMax int) Client {
	c.retryMax = retryMax
	return c
}

func (c *client) SetConnTimeout(connTimeout time.Duration) Client {
	c.connTimeout = connTimeout
	return c
}

func (c *client) SetURL(apiURL string) Client {
	c.apiURL = apiURL
	return c
}

func (c *client) SetKey(apiKey string) Client {
	c.apiKey = apiKey
	return c
}

func (c *client) SetEnableTransliteration(enable bool) Client {
	c.transliteration = enable
	return c
}

func (c *client) initClient() {
	client := &http.Client{
		Timeout: c.connTimeout,
		Transport: &http.Transport{
			DialContext:         (&net.Dialer{}).DialContext,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}

	httpClient := &retryablehttp.Client{
		RetryMax:   c.retryMax,
		Logger:     c.logger,
		HTTPClient: client,
		Backoff:    retryablehttp.LinearJitterBackoff,
		CheckRetry: retryablehttp.DefaultRetryPolicy,
		RequestLogHook: func(logger retryablehttp.Logger, request *http.Request, i int) {
			if i > 0 {
				logger.Printf("retry url %s attempt %d", request.URL.Path, i)
			}
		},
	}
	c.httpClient = httpClient
}
