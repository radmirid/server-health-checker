package checking

import (
	"net/http"
	"time"
)

type worker struct {
	client *http.Client
}

func newWorker(timeout time.Duration) *worker {
	return &worker{
		&http.Client{
			Timeout: timeout,
		},
	}
}

func (w worker) handle(j Job) Result {
	result := Result{Link: j.Link}

	now := time.Now()

	resp, err := w.client.Get(j.Link)
	if err != nil {
		result.Error = err

		return result
	}

	result.StatusCode = resp.StatusCode
	result.ResponseTime = time.Since(now)

	return result
}
