package bot

import (
	"log"
	"os"

	"bingo-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🤖 Bot started:", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		// 🟢 START
		if update.Message.Text == "/start" {

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"🎯 Welcome to Bingo!\n\n👉 Share your phone to continue",
			)

			btn := tgbotapi.NewKeyboardButtonContact("📱 Share Contact")

			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(btn),
			)

			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}

		// 🟢 CONTACT
		if update.Message.Contact != nil {

			tgID := update.Message.From.ID
			name := update.Message.From.UserName
			if name == "" {
				name = update.Message.From.FirstName
			}
			phone := update.Message.Contact.PhoneNumber

			token, err := services.RegisterUser(tgID, name, phone)
			if err != nil {
				log.Println("Register error:", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Try again"))
				continue
			}

			loginURL := os.Getenv("FRONTEND_URL") + "/login?token=" + token

			btn := tgbotapi.NewInlineKeyboardButtonURL("🎮 Play Bingo", loginURL)
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(btn),
			)

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"✅ You're in!\nTap below to play 👇",
			)
			msg.ReplyMarkup = keyboard

			bot.Send(msg)
		}
	}
}