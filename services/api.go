package services

import (
	"bytes"
	"encoding/json"
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