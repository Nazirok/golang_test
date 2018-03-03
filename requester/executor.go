package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func RequestIssueExecutor(ctx context.Context) (resp *http.Response, err error) {
	result := ctx.Value("result")
	client := &http.Client{}
	req := &http.Request{}

	switch result.Method {
	case "GET":
		req, err = http.NewRequest("", result.Url, nil)
		if err != nil {
			return
		}

	case "POST":
		temp, err := json.Marshal(result.Body)
		if err != nil {
			return
		}
		req, err = http.NewRequest("POST", result.Url, bytes.NewBuffer(temp))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", result.ContentType)

	default:
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return resp, err
}
