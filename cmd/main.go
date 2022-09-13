package main

import (
	"flag"
	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"social/internal/callbackmsg"
	"social/internal/chatmembers"
	"social/internal/config"
	"social/internal/functions"
	"social/internal/textmsg"
	"social/pkg/client/telegram"
	"social/pkg/logging"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "./etc/tgbot/.env", "config file path")
}

func main() {

	log.Printf("config initializing from %s", cfgPath)
	cfg := config.GetConfig(cfgPath)

	log.Printf("logger initializing level: %s", cfg.AppConfig.LogLevel)
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	updChan, bot, err := telegram.StartBotByChan(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := functions.NewFuncList(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = db.NewData()
	if err != nil {
		logger.Error(err)
	}

	moderGroup, err := db.GetModeratorsGroup()
	if err != nil {
		logger.Error(err)
	}

	if len(moderGroup) == 0 {

		groupInfo, _ := bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, "test"))
		_, _, err = db.AddModeratorsGroup(cfg.ModersGroupID.ModeratorsGroup, groupInfo.Chat.Title)
		if err != nil {
			logger.Info(err)
		}
		_, _ = bot.Send(tgb.NewDeleteMessage(groupInfo.Chat.ID, groupInfo.MessageID))
	}

	for {

		update := <-updChan

		if update.Message != nil {

			if update.Message.Text != "" {
				// text messages operations
				textmsg.WithTextQueryDo(update, bot, logger, cfg)

			} else if update.Message.NewChatMembers != nil {

				chatmembers.WithChatMembersDo(update, bot, logger, cfg)
			}

		} else if update.CallbackQuery != nil {

			callbackmsg.WithCallBackDo(update, bot, logger, cfg)
			// TODO inline help

		} else if update.InlineQuery != nil {

			query := update.InlineQuery.Query
			logger.Printf("response from Inline query: %s", query)
		}
	}
}
