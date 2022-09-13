package textmsg

import (
	"fmt"
	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"social/internal/config"
	"social/internal/functions"
	"social/internal/menu"
	"social/pkg/logging"
	"strconv"
	"strings"
	"time"
)

var MesInfo tgb.Update

func WithTextQueryDo(update tgb.Update, bot *tgb.BotAPI, logger *logging.Logger, cfg *config.Config) {

	db, _ := functions.NewFuncList(cfg, logger)

	//defer db.DbClose()

	myGroup, _ := db.GetModeratorsGroup()

	command, err := functions.TrimSymbolsFromSlice(strings.Split(update.Message.Text, " "), cfg)
	if err != nil {
		logger.Info("error trim symbols from message")
	}

	if strings.Contains(strings.ToLower(command[0]), "chatinfo") {

		delMes := update.Message.MessageID
		chatId := update.Message.Chat.ID
		groupTitle := update.Message.Chat.Title
		groupName := update.Message.Chat.FirstName

		count, err := bot.GetChatMembersCount(tgb.ChatMemberCountConfig{
			ChatConfig: tgb.ChatConfig{
				ChatID:             chatId,
				SuperGroupUsername: groupName,
			},
		})
		logger.Info("count: ", count)
		if err != nil {
			logger.Error(err)
		}

		go func() {
			_, _ = bot.Send(tgb.NewDeleteMessage(chatId, delMes))
		}()

		msg := tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, fmt.Sprintf(
			"ID группы: `%d`\nИмя группы: *%s*, %d", chatId, groupTitle, count))
		msg.ParseMode = "markdown"

		_, err = bot.Send(msg)
		if err != nil {
			logger.Error(err)
		}
	}

	for _, jobGroup := range myGroup {

		if update.Message.Chat.ID == jobGroup.UserGroupID || update.Message.Chat.ID == jobGroup.ModerGroupID {

			// trim symbols
			if len(update.Message.Text) > 0 {

				command, err := functions.TrimSymbolsFromSlice(strings.Split(update.Message.Text, " "), cfg)
				if err != nil {
					logger.Info("error trim symbols from message")
				}

				// check bad words in chat messages +
				_, b, err := db.CheckBadWords(command)
				if err != nil {
					logger.Error("bad words error: ", err)
				}

				// message to chat where found bad word, copy to moderators groups +
				if b && strings.ToLower(command[0]) != "мат" {

					//msgID := update.Message.MessageID
					badText := update.Message.Text

					//go func() {
					//	_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, msgID))
					//}()

					badGuyName := update.Message.From.FirstName
					badGuyNick := update.Message.From.UserName
					badGuyID := update.Message.From.ID
					groupName := update.Message.Chat.Title

					moderatorsGroups, err := db.GetModeratorsGroup()
					if err != nil {
						logger.Error(err)
					}

					for _, v := range moderatorsGroups {

						if v.UserGroupID == update.Message.Chat.ID {

							modMess := tgb.NewMessage(v.ModerGroupID, fmt.Sprintf(
								"Найдены нецензурные выражения:\nГруппа: %s\nИмя пользователя: %s\nНик пользователя: "+
									"@%s\nID пользователя: %d\nТекст сообщения: %s\nВремя: %s\nОригинал сообщения удален из чата.",
								groupName, badGuyName, badGuyNick, badGuyID, badText, time.Now().Format(config.StructDateTimeFormat)))
							_, _ = bot.Send(modMess)

							//_, _ = bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, fmt.Sprintf(
							//	"Найдены нецензурные выражения:\nГруппа: %s\nИмя пользователя: %s\nНик пользователя: "+
							//		"@%s\nID пользователя: %d\nТекст сообщения: %s\nВремя: %s\nОригинал сообщения удален из чата.",
							//	groupName, badGuyName, badGuyNick, badGuyID, badText, time.Now().Format(config.StructDateTimeFormat))))

						}

					}

					cleanAnswer := tgb.NewMessage(update.Message.Chat.ID, fmt.Sprintf(
						cfg.MsgText.MsgOfBadWordToUserChat, groupName))
					del, _ := bot.Send(cleanAnswer)

					go func() {
						time.Sleep(30 * time.Second)
						_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, del.MessageID))
					}()
				}

				// add moderator group
				if strings.Contains(strings.ToLower(command[0]), "addmoderatorgroup") {

					if update.Message.Chat.ID == cfg.ModersGroupID.ModeratorsGroup {

						newModGroup, err := strconv.ParseInt(command[1], 10, 64)
						if err != nil {
							logger.Error(err)

						}

						info, _ := bot.Send(tgb.NewMessage(newModGroup, "test"))

						b, _, err := db.AddModeratorsGroup(newModGroup, info.Chat.Title)
						if err != nil {
							logger.Error(err)
						}
						_, _ = bot.Send(tgb.NewDeleteMessage(newModGroup, info.MessageID))

						if b && err != nil {

							_, _ = bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, fmt.Sprintf("Такая группа уже есть: %d", newModGroup)))
						}

						if b && err == nil {

							_, _ = bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, fmt.Sprintf("Успешно добавлена: %d", newModGroup)))

						}
					}
				}

				//  new users count +
				if strings.Contains(strings.ToLower(command[0]), "chatinfo") {

					delMes := update.Message.MessageID
					chatId := update.Message.Chat.ID
					groupTitle := update.Message.Chat.Title
					groupName := update.Message.Chat.FirstName

					count, err := bot.GetChatMembersCount(tgb.ChatMemberCountConfig{
						ChatConfig: tgb.ChatConfig{
							ChatID:             chatId,
							SuperGroupUsername: groupName,
						},
					})
					logger.Info("count: ", count)
					if err != nil {
						logger.Error(err)
					}

					go func() {
						_, _ = bot.Send(tgb.NewDeleteMessage(chatId, delMes))
					}()

					msg := tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, fmt.Sprintf(
						"ID группы: `%d`\nИмя группы: *%s*, %d", chatId, groupTitle, count))
					msg.ParseMode = "markdown"

					_, err = bot.Send(msg)
					if err != nil {
						logger.Error(err)
					}
				}

				// add bad words to the base +
				if strings.Contains(strings.ToLower(command[0]), "мат") && len(command) > 1 {

					moderators, err := db.GetModeratorsGroup()
					if err != nil {
						logger.Error(err)
					}

					for _, moder := range moderators {

						if update.Message.Chat.ID == moder.ModerGroupID {

							go func() {
								time.Sleep(5 * time.Second)
								_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))

							}()

							b, err := db.AddBadWord(command[1])
							if err != nil {
								log.Println(err)
							}

							if b == true && err == nil {

								del, _ := bot.Send(tgb.NewMessage(update.Message.Chat.ID, "Уже есть."))

								go func() {

									time.Sleep(5 * time.Second)
									_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, del.MessageID))

								}()

							} else if b == true && err != nil {

								del, _ := bot.Send(tgb.NewMessage(update.Message.Chat.ID, "Добавлено."))

								go func() {

									time.Sleep(5 * time.Second)
									_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, del.MessageID))
								}()
							}
						}
					}
				}

				// menu moderators groups +
				if strings.Contains(strings.ToLower(command[0]), "меню") {

					if err != nil {
						return
					}

					if jobGroup.ModerGroupID == update.Message.Chat.ID {

						_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))

						msg := tgb.NewMessage(update.Message.Chat.ID, "Меню закроется через 1 минуту")

						msg.ReplyMarkup = menu.NumericKeyboard

						delMes, err := bot.Send(msg)
						if err != nil {
							logger.Error(err)
						}

						go func() {
							time.Sleep(60 * time.Second)
							_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, delMes.MessageID))
						}()
						break

					} else if jobGroup.UserGroupID == update.Message.Chat.ID {

						_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))

						msg := tgb.NewMessage(update.Message.Chat.ID, "Меню закроется через 1 минуту")

						msg.ReplyMarkup = menu.UserKeyboard

						delMes, err := bot.Send(msg)
						if err != nil {
							logger.Error(err)
						}

						go func() {
							time.Sleep(60 * time.Second)
							_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, delMes.MessageID))
						}()
						break
					}

				}

				// add moder group
				if strings.Contains(strings.ToLower(command[0]), "add-moder-group") {

					MesInfo = update

					go func() {
						time.Sleep(60 * time.Second)
						MesInfo.Message.Chat.ID = 0
					}()

					chatId := update.Message.Chat.ID
					chatName := update.Message.Chat.Title
					userName := update.Message.From.FirstName
					userNick := update.Message.From.UserName

					_, _ = bot.Send(tgb.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))

					text := fmt.Sprintf("Новый запрос на добавление группы администраторов:\nНазвание группы: %s\n"+
						"Уникальный номер группы: %d\nИмя пользователя: %s\nНик пользователя: @%s\nВремя запроса: %s\n"+
						"Подтвердите в течении 60 секунд, или проигнорируйте сообщение.",
						chatName, chatId, userName, userNick, time.Now().Format(config.StructDateTimeFormat))

					msg := tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, text)
					msg.ReplyMarkup = tgb.NewInlineKeyboardMarkup(tgb.NewInlineKeyboardRow(menu.Button4, menu.Button11))

					msgDel, err := bot.Send(msg)
					if err != nil {
						logger.Error(err)
					}

					go func() {
						time.Sleep(60 * time.Second)
						_, _ = bot.Send(tgb.NewDeleteMessage(msgDel.Chat.ID, msgDel.MessageID))
					}()
				}

				// user groups link with moder group
				if strings.Contains(strings.ToLower(command[0]), "add-moder-user-link") {

					if update.Message.Chat.ID == cfg.ModersGroupID.ModeratorsGroup {

						moderGroup, _ := strconv.ParseInt(command[1], 10, 64)
						userGroup, _ := strconv.ParseInt(command[2], 10, 64)

						moderInfo, err1 := bot.Send(tgb.NewMessage(moderGroup, "test"))
						if err != nil {
							logger.Error(err1)
						}

						userInfo, err2 := bot.Send(tgb.NewMessage(userGroup, "test"))
						if err2 != nil {
							logger.Error(err2)
						}

						_, _ = bot.Send(tgb.NewDeleteMessage(moderInfo.Chat.ID, moderInfo.MessageID))
						_, _ = bot.Send(tgb.NewDeleteMessage(userInfo.Chat.ID, userInfo.MessageID))

						b, err = db.AddUserGroupList(moderGroup, userGroup, moderInfo.Chat.Title, userInfo.Chat.Title)
						if !b && err != nil {
							logger.Error(err)
						}

						if b && err == nil {
							_, _ = bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, "Такая группа пользователей уже есть в базе."))
						}
						if !b && err == nil {
							_, _ = bot.Send(tgb.NewMessage(cfg.ModersGroupID.ModeratorsGroup, "Успешно добавлено."))
						}
					}
				}

				if len(command) > 1 {

					if strings.ToLower(command[0]) == "спокойной" && strings.ToLower(command[1]) == "ночи" {

						if update.Message.From.FirstName == "Az" {

							msg := tgb.NewMessage(update.Message.Chat.ID, "Спокойной ночи, Создатель!")
							_, _ = bot.Send(msg)

						} else {

							name := update.Message.From.FirstName
							msg := tgb.NewMessage(update.Message.Chat.ID, fmt.Sprintf("И тебе спокойной ночи, %s", name))
							_, _ = bot.Send(msg)
						}
					}
				} //  Good night

				switch strings.ToLower(command[0]) {

				case "привет", "хаюшки", "здарова", "салют", "здравствуйте", "приветик", "здравствуй", "приветствую":

					name := update.Message.From.FirstName
					msg := "Здравствуй, " + name
					del, _ := bot.Send(tgb.NewMessage(update.Message.Chat.ID, msg))

					go func() {
						time.Sleep(60 * time.Second)
						msDel := tgb.NewDeleteMessage(update.Message.Chat.ID, del.MessageID)
						_, _ = bot.Send(msDel)
					}()

				case "доброе":

					if len(command) > 1 {

						if strings.ToLower(command[1]) == "утро" {

							if update.Message.From.FirstName == "Az" {

								msg := tgb.NewMessage(update.Message.Chat.ID, "С Добрым Утром! Счастлив видеть тебя, Создатель!")
								_, _ = bot.Send(msg)

							} else {

								name := update.Message.From.FirstName
								goodMorn := "Здравствуй, " + name + ". "
								text, _ := db.GoodMorningMessage()
								msg := goodMorn + text
								_, _ = bot.Send(tgb.NewMessage(update.Message.Chat.ID, msg))
							}
						}
					}

				case "утро":

					if len(command) > 1 {

						if strings.ToLower(command[1]) == "доброе" {

							if update.Message.From.FirstName == "Az" {

								msg := tgb.NewMessage(update.Message.Chat.ID, "С Добрым Утром! Счастлив видеть тебя, Создатель!")
								_, _ = bot.Send(msg)

							} else {

								name := update.Message.From.FirstName
								goodMorn := "И тебе доброго утра, " + name + ". "
								text, _ := db.GoodMorningMessage()
								msg := goodMorn + text
								_, _ = bot.Send(tgb.NewMessage(update.Message.Chat.ID, msg))
							}
						}
					}

				default:

				}
			}

		}
	}

}
