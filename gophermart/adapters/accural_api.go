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
	TooManyRequestsErr = errors.New("Too many requests")
	UnexpectedErr	   = errors.New("Unexpected error")
)

type AccuralApi struct {
	client *http.Client
	baseUrl string
	timeout time.Duration
	retries int
}

func NewAccuralApi(client *http.Client, baseUrl string, timeout time.Duration, retries int) *AccuralApi {
	return &AccuralApi{client: client, baseUrl: baseUrl, timeout: timeout, retries: retries}
}


type Order struct {
	Number string  `json:"number"`
	Status string  `json:"status"`
	Accural uint64 `json:"accural,omitempty"`
}

func (api *AccuralApi) GetOrderAccuralStatus(ctx context.Context, order string) (result *Order, err error) {
	var (
		res *http.Response
		retries int = api.retries
	)

    for retries > 0 {
		ctx, cancel := context.WithTimeout(ctx, api.timeout)
		defer cancel()

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/orders/%s", api.baseUrl, order), http.NoBody)
		if err != nil {
			return nil, err
		}

		res, err = api.client.Do(req)
		if err != nil {
			retries -= 1
			continue
		}
	}
	if err != nil {
		return nil, err
	}

	var buf *bytes.Buffer
	if res.Body != nil {
		defer res.Body.Close()
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
		return nil, TooManyRequestsErr
	case http.StatusInternalServerError | http.StatusBadGateway | http.StatusServiceUnavailable:
		return nil, UnexpectedErr
	default:
		log.WithField("code", res.StatusCode).Warn("Unhandled status code from accural api")
		return nil, UnexpectedErr
	}
}
