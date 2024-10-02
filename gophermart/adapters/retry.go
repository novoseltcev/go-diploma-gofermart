package adapters

import log "github.com/sirupsen/logrus"

func Retry[T any](f func() (*T, error), retries int) (result *T, err error) {
    for {
		result, err = f()
		if err == nil {
			return result, err
		}

		if retries == 0 {
			return nil, err
		}

		log.WithField("retries", retries).Warn("retrying...")
		retries--
	}
}
