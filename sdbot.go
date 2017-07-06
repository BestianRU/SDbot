package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"

)



func main()  {
	log.Println(" Bot was starting!")
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
 
	updates, err := bot.GetUpdatesChan(u)

	// var keyboard tgbotapi.ReplyKeyboardMarkup
	// keyboard.ResizeKeyboard=true
	// var phoneButton tgbotapi.KeyboardButton
	// phoneButton.Text="Send my your phone number"	
	// phoneButton.RequestContact=true
	// keyboard.Keyboard=append(keyboard.Keyboard,phoneButton)
	
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.ID , update.Message.Text)
		phoneButton:=tgbotapi.NewKeyboardButtonContact("Send my your phone number")
		var msg tgbotapi.KeyboardMsg
		
		row:=make([]tgbotapi.KeyboardButton, 1)
		row=append(row,phoneButton)
		msg.Keyboard=append(msg.Keyboard,row)
		
		msg.ResizeKeyboard=true
				
		msg.ChatID=update.Message.Chat.ID
		msg.Text="Send my your phone number"
		
		
		bot.Send(msg)
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
		
	//	bot.Send(msg)
	}
}

//func auth()