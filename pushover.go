package main

import (
	"net/http"
	"net/url"
	"os"
)

// NotifyPushover sends a notification using the Pushover API
func NotifyPushover(message string) error {
	token := os.Getenv("PUSHOVER_TOKEN")
	user := os.Getenv("PUSHOVER_TARGET")
	_, err := http.PostForm("https://api.pushover.net/1/messages.json", url.Values{
		"token":   {token},
		"user":    {user},
		"title":   {"Bins out!"},
		"message": {message},
	})
	if err != nil {
		return err
	}
	return nil
}
