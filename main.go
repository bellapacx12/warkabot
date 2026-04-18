package main

import (
	"log"
	"net/http"
	"os"

	"bingo-bot/bot"
)

func main() {
	// 🔥 Start bot in background
	go bot.Start()

	// 🌐 Fake HTTP server for Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🌐 Fake server running on port", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bot is running"))
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}