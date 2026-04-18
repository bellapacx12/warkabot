package main

import (
	"log"
	"net/http"
	"os"

	"bingo-bot/bot"
)

func main() {

	if os.Getenv("RUN_BOT") == "true" {
		go bot.Start()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bot is running"))
	})

	log.Println("🌐 Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}