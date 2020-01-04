package controllers

import (
	"fmt"
	appv1alpha1 "github.com/easyteam-fr/emojis/operator/api/v1alpha1"
	"net/http"
	"os"
)

var (
	application = "http://emojis:8081"
)

func init() {
	if os.Getenv("EMOJIS_ENDPOINT") != "" {
		application = os.Getenv("EMOJIS_ENDPOINT")
	}
}

func emojiDelete(emoji *appv1alpha1.Emoji) error {
	client := &http.Client{}
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/emojis/%s", application, emoji.Name),
		nil,
	)
	if err == nil {
		req.Header.Add("Content-Type", "application/json")
		_, err = client.Do(req)
	}
	return err
}

func emojiCreateOrUpdate(emoji *appv1alpha1.Emoji) error {

	client := &http.Client{}
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/emojis/%s", application, emoji.Name),
		nil,
	)
	if err == nil {
		req.Header.Add("Content-Type", "application/json")
		_, err = client.Do(req)
	}
	return err
}
