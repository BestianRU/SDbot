package main

import (
	"SDbot/cfg"
	"SDbot/user"
	"log"

	tgbotapi "github.com/DmitryBugrov/telegram-bot-api"
)

func main() {
	log.Println("Bot is starting!")
	log.Println("Load config from: sdbotcfg.json")
	c := new(cfg.Cfg)
	err := c.Load()
	if err != nil {
		log.Println("Error reading config:")
		panic(err)
	}

	// user, err := user.GetUserFromSQLByPhone("79622754090", c)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(user)

	//Init map of authorized users
	au := user.NewAuthUser(c)

	//Init bot
	bot, err := tgbotapi.NewBotAPI(c.T.Token)
	if err != nil {
		log.Println("Error connecting to telegram:")
		panic(err)
	}
	bot.Debug = c.T.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = c.T.Timeout

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		//user not authorized
		if _, err = au.GetByTId(uint64(update.Message.From.ID)); err != nil {
			msg := tgbotapi.NewMessage(int64(update.Message.From.ID), c.MsgNotAuth)
			bot.Send(msg)
		}
		//	log.Printf("[%d] %s", update.Message.From.ID, update.Message.Text)
		//	phoneButton := tgbotapi.NewKeyboardButtonContact("Send my your phone number")
		//	var msg tgbotapi.KeyboardMsg

		// row := make([]tgbotapi.KeyboardButton, 1)
		// row = append(row, phoneButton)
		// msg.Keyboard = append(msg.Keyboard, row)

		// msg.ResizeKeyboard = true

		// msg.ChatID = update.Message.Chat.ID
		// msg.Text = "Send my your phone number"

	}
}

//auth authorise user, return true if user is valid
func auth(phone string) bool {
	return true
}
