package script_svc

import (
	"context"
	"io"
	"net/http"
	"time"
)

func requestSyncUrl(ctx context.Context, syncUrl string) (string, error) {
	c := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, syncUrl, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(b), nil
}
