package main

import (
	"log"
	"SDbot/cfg"
	"SDbot/user"
	"SDbot/telegram-bot-api.v4"
)



func main()  {
	log.Println(" Bot was starting!")
	log.Println("Load config from: ./sdbotcfg.json")
	c:=new(cfg.Cfg)
	err:=c.Load()
	if err!=nil {
		panic(err )
	}
	user.GetUserFromSQLByPhone("",c)

	//Init bot
	bot, err := tgbotapi.NewBotAPI(c.T.Token)
	if err != nil {
		panic(err )
	}
	bot.Debug = c.T.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = c.T.Timeout
 
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%d] %s", update.Message.From.ID , update.Message.Text)
		phoneButton:=tgbotapi.NewKeyboardButtonContact("Send my your phone number")
		var msg tgbotapi.KeyboardMsg
		
		row:=make([]tgbotapi.KeyboardButton, 1)
		row=append(row,phoneButton)
		msg.Keyboard=append(msg.Keyboard,row)
		
		msg.ResizeKeyboard=true
				
		msg.ChatID=update.Message.Chat.ID
		msg.Text="Send my your phone number"
		
		
		bot.Send(msg)
	
	}
}

//auth authorise user, return true if user is valid
func auth(phone string) bool {
	return true
}

