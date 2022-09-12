package helpers

import (
	"io"
	"net/http"
	"net/url"
)

func HTTPGet(u *url.URL) (string, error) {
	urlstr := u.String()
	resp, err := http.Get(urlstr)
	if err != nil {
		return "", err
	}
	// TODO; Check Response Status (200), Handle Non-200 status gracefully
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Convert to string
	b := string(body)
	return b, nil
}
