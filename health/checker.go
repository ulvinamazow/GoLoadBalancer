package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Checker(method string, URL string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, URL, nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return 0, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	code := resp.StatusCode

	return code, nil
}
