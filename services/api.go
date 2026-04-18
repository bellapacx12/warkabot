package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func RegisterUser(tgID int64, name, phone string) (string, error) {

	body := map[string]interface{}{
		"telegram_id": tgID,
		"name":        name,
		"phone":       phone,
	}

	jsonData, _ := json.Marshal(body)

	resp, err := http.Post(
		os.Getenv("BACKEND_URL"),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["token"], nil
}
func GetUserToken(tgID int64) (string, bool) {

	url := fmt.Sprintf("%s/auth/telegram?telegram_id=%d", os.Getenv("BACKEND_URL"), tgID)

	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", false
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	return result["token"], true
}