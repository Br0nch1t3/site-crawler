package utilshttp

import (
	"fmt"
	"net/http"
)

type Error struct {
	res *http.Response
}

func (e Error) Error() string {
	return fmt.Sprintf("request failed with status %d for %s", e.res.StatusCode, e.res.Request.URL)
}

func ExtractError(res *http.Response) error {
	if res.StatusCode >= 400 && res.StatusCode <= 599 {
		return &Error{res: res}
	}

	return nil
}
