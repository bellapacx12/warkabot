package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	url := os.Getenv("BACKEND_URL") + "?telegram_id=" + fmt.Sprint(tgID)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("HTTP error:", err)
		return "", false
	}
	defer resp.Body.Close()

	// 🔥 IMPORTANT
	if resp.StatusCode == 404 {
		return "", false // user not found
	}

	if resp.StatusCode != 200 {
		log.Println("Bad status:", resp.StatusCode)
		return "", false
	}

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("Decode error:", err)
		return "", false
	}

	token := result["token"]
	if token == "" {
		return "", false
	}

	return token, true
}