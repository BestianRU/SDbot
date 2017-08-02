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

	//Init map of authorized users
	au, err := user.NewAuthUser(c)
	if err != nil {
		log.Println("Error load authorized users:", c.AuthUser)
		panic(err)
	}

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

	in, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	out := make(chan tgbotapi.Chattable, 512)
	go ReadMessages(in, out, au, c)
	//go SendMessages(bot, out)
	//send message
	for msg := range out {
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		//	log.Println(m)

	}
}

//auth authorise user, return true if user is valid
// func auth(phone string) bool {
// 	return true
// }

//ReadMessages from telegram
func ReadMessages(in tgbotapi.UpdatesChannel, out chan tgbotapi.Chattable, au *user.AuthUser, c *cfg.Cfg) {
	for update := range in {
		if update.Message == nil {
			continue
		}
		//read command
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "auth":
				//	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "/auth")
				out <- RequestPhone(int64(update.Message.From.ID), c)
			}
			continue
		}

		//user not authorized
		if _, err := au.GetByTId(uint64(update.Message.From.ID)); err != nil {
			if update.Message.Contact == nil {
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), c.Msg.MsgNotAuth)
				out <- msg
				continue
			}
			if update.Message.Contact.PhoneNumber == "" {
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), c.Msg.MsgNotAuth)
				out <- msg
				continue
			}
			//if client send your phone number
			u, err := user.GetUserFromSQLByPhone(update.Message.Contact.PhoneNumber, c)
			//phone number not found
			if err != nil {
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), c.Msg.PhoneNotFound)
				out <- msg
				continue
			}
			//add new user
			u.TId = uint64(update.Message.From.ID)
			err = au.Add(u, c)
			if err != nil {
				log.Println(err)
				continue
			}
			msg := tgbotapi.NewMessage(int64(update.Message.From.ID), c.Msg.AuthMsg)
			out <- msg

		}

	}
}

//RequestPhone create message with button getting phone number
func RequestPhone(id int64, c *cfg.Cfg) tgbotapi.MessageConfig {

	var msg tgbotapi.MessageConfig
	phoneButton := tgbotapi.NewKeyboardButtonContact(c.Msg.TextPhoneButton)
	row := tgbotapi.NewKeyboardButtonRow(phoneButton)
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
	}
	keyboard.Keyboard = append(keyboard.Keyboard, row)
	msg.ReplyMarkup = keyboard
	msg.Text = c.Msg.RequestPhone
	msg.ChatID = id
	return msg

}

//SendMessages send message to telegram
// func SendMessages(bot *tgbotapi.BotAPI, out chan tgbotapi.Chattable) {
// 	for msg := range out {
// 		bot.Send(msg)
// 	}
// }
