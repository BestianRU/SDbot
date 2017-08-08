# Telegram bot for ServiceDesk GLPI

The bot is waiting new notifications for users and send them telegram messages. For reading new notifications the bot must connect to mysql database GLPI. 

##Install:##
1. register new bot: https://core.telegram.org/bots and get uniqual token;
2. configure bot - file sdbotcfg.json
   Set telegram token and parameters for connection to mysql database GLPI
    ```
    {
        "telegram": {
            "token":"token for bot",
            "timeout":60,
            "debug":false
        }, 
        "mysql"  :{
            "host":"127.0.0.1",
            "port":"3306",
            "database":"glpi",
            "user":"glpi",
            "pass":"123456"
        },
        "authUser":"authuser.json",
        "notificationsPeriod":60,
        "messages":{
            "msgNotAuth": "You are not authorized. Please, send me command /auth and allow me reading your phone number.",
            "textPhoneButton":"Send my your phone number",
            "requestPhone":"Please click to button below",
            "phoneNotFound":"Sorry, but your phone number not found in Service Desk, please check your contact information in Service Desk",
            "authMsg":"You are authorized!",
            "idontknow":"Sorry, but I don't know that mean: "
        }
    }
    ```
3. Start sdbot

##Using##
After starting you can find bot by name in your telegramm client. For your identification you must allow reading your mobile number. For it, send command    "/auth". If your phone number exist in GLPI you authorized. Now you will receive all the messages from the GLPI, which the servicedesk sends you by e-mail.



