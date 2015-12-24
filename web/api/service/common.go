package service

import (
	"fmt"
	"net/http"
)

func common(m string, r *http.Request) error {
	if r.Method != m {
		return fmt.Errorf("method not allowed")
	}
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("unable to process argument")
	}

	return nil
}
