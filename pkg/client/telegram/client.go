package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"social/internal/config"
	"social/pkg/logging"
	"time"
)

func StartBotByChan(cfg *config.Config, logger *logging.Logger) (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI, error) {

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	logger.Infof("Bot %s started at: %v", bot.Self.UserName, time.Now().Format("2 January 2006 15:04"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.Limit = 0
	u.Offset = 0

	updates := bot.GetUpdatesChan(u)

	return updates, bot, nil

}

//func StartBotByHook(cfg config.Config, logger *logging.Logger) (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI, error) {
//
//	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
//	if err != nil {
//		log.Fatal("bot didn't started, Token error.")
//	}
//
//	wh, _ := tgbotapi.NewWebhookWithCert("https://www.example.com:8443/" + cfg.Telegram.Token, cfg.Telegram.Sert)
//
//	_, err = bot.Request(wh)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	info, err := bot.GetWebhookInfo()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if info.LastErrorDate != 0 {
//		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
//	}
//
//	updates := bot.ListenForWebhook("/" + bot.Token)
//	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
//
//	fmt.Printf("Bot %s started at: %v", bot.Self.UserName, time.Now().Format("2 January 2006 15:04"))
//
//	for update := range updates {
//		log.Printf("%+v\n", update)
//	}
//
//	// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 3560 -subj "//O=Org\CN=Test" -nodes
//	// https://github.com/go-telegram-bot-api/telegram-bot-api
//
//	return nil, nil, nil
//}
