package chatmembers

import (
	"fmt"
	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"social/internal/config"
	"social/internal/functions"
	"social/pkg/logging"
	"time"
)

var NewUserID int64

func WithChatMembersDo(update tgb.Update, bot *tgb.BotAPI, logger *logging.Logger, cfg *config.Config) {

	db, err := functions.NewFuncList(cfg, logger)
	if err != nil {
		logger.Error(err)
	}

	newUser := update.Message.NewChatMembers[0]
	NewUserID = newUser.ID
	chatId := update.Message.Chat.ID
	groupName := update.Message.Chat.Title

	logger.Infof("from members NewUserID %d", NewUserID)

	if !newUser.IsBot {

		moderGroupList, err := db.GetModeratorsGroup()
		if err != nil {
			logger.Error(err)
		}

		for _, group := range moderGroupList {

			if update.Message.Chat.ID == group.UserGroupID {

				count, err := bot.GetChatMembersCount(tgb.ChatMemberCountConfig{
					ChatConfig: tgb.ChatConfig{
						ChatID:             chatId,
						SuperGroupUsername: groupName,
					},
				})

				msg := tgb.NewMessage(chatId, fmt.Sprintf(cfg.MsgText.MsgToNewUser, newUser.FirstName))
				ans, _ := bot.Send(msg)

				go func() {

					time.Sleep(60 * time.Second)
					_, _ = bot.Send(tgb.NewDeleteMessage(chatId, ans.MessageID))
				}()

				err = db.AddNewJubileeUser(&newUser, count, update)
				if err != nil {
					logger.Error(err)
				}

				for _, group := range moderGroupList {

					if group.ModerGroupID != 0 {

						text := fmt.Sprintf("🎉 В группу: %s вступил новый пользователь!\nИмя: %s "+
							"\nНик: @%s, \nНомер вступления: %d. \nВремя вступления %s",
							groupName, newUser.FirstName, newUser.UserName, count,
							time.Now().Format(config.StructDateTimeFormat))
						msg := tgb.NewMessage(group.ModerGroupID, text)

						_, _ = bot.Send(msg)

					}
				}
			}
		}
	}
}
