package httpclient

import (
	"net/http"
	"time"
)

var instance *http.Client

func Instance() *http.Client {
	if instance != nil {
		return instance
	}
	instance = &http.Client{
		Timeout: 5 * time.Second,
	}

	return instance
}
