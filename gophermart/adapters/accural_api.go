package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrUnexpected	   = errors.New("unexpected error")
)

type AccuralAPI struct {
	client *http.Client
	baseURL string
	timeout time.Duration
	retries int
}

func NewAccuralAPI(client *http.Client, baseURL string, timeout time.Duration, retries int) *AccuralAPI {
	return &AccuralAPI{client: client, baseURL: baseURL, timeout: timeout, retries: retries}
}


type Order struct {
	Number string  `json:"number"`
	Status string  `json:"status"`
	Accural uint64 `json:"accural,omitempty"`
}

func (api *AccuralAPI) GetOrderAccuralStatus(ctx context.Context, order string) (result *Order, err error) {
	var (
		res *http.Response
		retries = api.retries
	)

    for retries > 0 {
		ctx, cancel := context.WithTimeout(ctx, api.timeout)
		defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/orders/%s", api.baseURL, order), http.NoBody)
		if err != nil {
			return nil, err
		}

		res, err = api.client.Do(req)
		if err != nil {
			retries -= 1
		} else {
			defer res.Body.Close()
			break
		}
	}
	if err != nil {
		return nil, err
	}

	var buf *bytes.Buffer
	if res.Body != nil {
		if _, err := io.Copy(buf, res.Body); err != nil {
			return nil, err
		}
	}

	switch res.StatusCode {
	case http.StatusOK:
		return result, json.Unmarshal(buf.Bytes(), result)
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusInternalServerError | http.StatusBadGateway | http.StatusServiceUnavailable:
		return nil, ErrUnexpected
	default:
		log.WithField("code", res.StatusCode).Warn("Unhandled status code from accural api")
		return nil, ErrUnexpected
	}
}
