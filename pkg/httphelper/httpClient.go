package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTPOption func(c *HTTPCli)

type IHttpClient interface {
	Post(ctx context.Context, urll string, reqbody interface{}, resp interface{}, options ...HTTPOption) error
}

type HTTPCli struct {
	client http.Client
	req    *http.Request

	sseStream chan string
}

var _ IHttpClient = (*HTTPCli)(nil)

func NewHTTPClient() *HTTPCli {
	return &HTTPCli{
		client:    http.Client{},
		sseStream: nil,
	}
}

func (c *HTTPCli) Post(ctx context.Context, urll string, reqbody interface{}, respbody interface{}, options ...HTTPOption) error {
	// options = append(options, WithHeader(HeaderMap{"content-type": "application/json"}))

	// if reqbody == nil {
	// 	return
	// }

	resp, err := c.httpInner(ctx, "POST", urll, reqbody, options...)
	if err != nil {
		return err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Println("[post] result: ", string(result))

	if len(result) != 0 && respbody != nil {
		err = json.Unmarshal(result, &respbody)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *HTTPCli) EncodeJSONBody(body interface{}) (*bytes.Buffer, error) {
	if body != nil {
		var bodyJSON []byte
		switch body := body.(type) {
		case *bytes.Buffer:
			return body, nil
		case []byte:
			bodyJSON = body
		default:
			var err error
			bodyJSON, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
		return bytes.NewBuffer(bodyJSON), nil
	}
	return bytes.NewBuffer(nil), nil
}

func (c *HTTPCli) httpInner(ctx context.Context, method, url string, body interface{}, options ...HTTPOption) (*http.Response, error) {
	var err error

	bodyBuffer, err := c.EncodeJSONBody(body)
	if err != nil {
		return nil, err
	}
	log.Printf("[http-req] body: %+v\n", bodyBuffer.String())

	c.req, err = http.NewRequestWithContext(ctx, method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(c)
	}

	resp, err := c.client.Do(c.req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		_, errIo := io.ReadAll(resp.Body)
		if errIo != nil {
			err = errIo
			return resp, err
		}

		return resp, err
	}
	// Note: close rsp at outer function
	return resp, nil
}
