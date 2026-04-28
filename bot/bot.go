package bot

import (
	"context"
	"log"
	"os"

	"bingo-bot/services"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start() {
	ctx := context.Background()

	b, err := bot.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🤖 Bot started")

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {

			if update.Message == nil {
				return
			}

			tgID := update.Message.From.ID

			// 🟢 /start
			if update.Message.Text == "/start" {

				// 🔥 CHECK USER VIA API
				token, exists := services.GetUserToken(tgID)

				if exists {
					// ✅ Returning user
					loginURL := os.Getenv("FRONTEND_URL") + "/login?token=" + token

					btn := models.InlineKeyboardButton{
						Text: "🎮 Play Bingo",
						WebApp: &models.WebAppInfo{
							URL: loginURL,
						},
					}
                    directBtn := models.InlineKeyboardButton{
	Text: "🌐 Open in Browser",
	URL:  loginURL,
}
					keyboard := models.InlineKeyboardMarkup{
						InlineKeyboard: [][]models.InlineKeyboardButton{
							{btn},
							{directBtn},
						},
					}

					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      update.Message.Chat.ID,
						Text:        "🎯 Welcome back!\nTap below to play 👇",
						ReplyMarkup: keyboard,
					})

					return
				}

				// 🆕 New user → request contact
				btn := models.KeyboardButton{
					Text:           "📱 Share Contact",
					RequestContact: true,
				}

				keyboard := models.ReplyKeyboardMarkup{
					Keyboard: [][]models.KeyboardButton{
						{btn},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        "🎯 Welcome to Bingo!\n\n👉 Share your phone to continue",
					ReplyMarkup: keyboard,
				})

				return
			}

			// 🟢 CONTACT → REGISTER VIA API
			if update.Message.Contact != nil {

				name := update.Message.From.Username
				if name == "" {
					name = update.Message.From.FirstName
				}

				phone := update.Message.Contact.PhoneNumber

				token, err := services.RegisterUser(tgID, name, phone)
				if err != nil {
					log.Println("Register error:", err)

					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "❌ Try again",
					})
					return
				}

				loginURL := os.Getenv("FRONTEND_URL") + "/login?token=" + token

				btn := models.InlineKeyboardButton{
					Text: "🎮 Play Bingo",
					WebApp: &models.WebAppInfo{
						URL: loginURL,
					},
				}

				keyboard := models.InlineKeyboardMarkup{
					InlineKeyboard: [][]models.InlineKeyboardButton{
						{btn},
					},
				}

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      update.Message.Chat.ID,
					Text:        "✅ You're in!\nTap below to play 👇",
					ReplyMarkup: keyboard,
				})
			}
		},
	)

	b.Start(ctx)
}