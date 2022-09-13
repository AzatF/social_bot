package functions

import (
	"database/sql"
	"errors"
	"fmt"
	tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"os"
	"path"
	"social/internal/config"
	"social/internal/model"
	"social/pkg/logging"
	"strings"
	"time"
)

type list struct {
	cfg    *config.Config
	logger *logging.Logger
	db     *sql.DB
}

func NewFuncList(cfg *config.Config, logger *logging.Logger) (FuncList, error) {

	err := os.MkdirAll(cfg.DBFilePath, 0777)
	if err != nil {
		logger.Error(err)
	}

	liteDb, err := sql.Open("sqlite3", path.Join(cfg.DBFilePath, "tg_bot_db.db"))
	if err != nil {
		logger.Fatalf("error open database %v", err)
	}

	return &list{
		cfg:    cfg,
		logger: logger,
		db:     liteDb,
	}, nil

}

func (l *list) NewData() error {

	_, err := l.db.Exec(config.MakeModeratorsGroupTable)
	if err != nil {
		return err
	}

	_, err = l.db.Exec(config.MakeBadWordsTable)
	if err != nil {
		return err
	}

	_, err = l.db.Exec(config.MakeUserGroupTable)
	if err != nil {
		return err
	}

	_, err = l.db.Exec(config.MakeBotMessageGoodMorningTable)
	if err != nil {
		return err
	}

	_, err = l.db.Exec(config.MakeBotMessageComplimentTable)
	if err != nil {
		return err
	}

	_, err = l.db.Exec(config.MakeAnecdoteTable)
	if err != nil {
		return err
	}

	return nil
}

type FuncList interface {
	NewData() error
	GoodMorningMessage() (string, error)
	Compliment(sex string) (compliment string, err error)
	CheckSex(un string) (haveUser bool, s string, err error)
	AddUserGroupList(moderGroup, userGroup int64, moderTitle, userTitle string) (bool, error)
	CheckBadWords(badList []string) (clearList []string, haveBadWords bool, err error)
	AddBadWord(word string) (bool, error)
	AddModeratorsGroup(group int64, title string) (haveGroup bool, modGroups []model.ModeratorsGroup, err error)
	GetModeratorsGroup() (groups []model.ModeratorsGroup, err error)
	AddNewJubileeUser(newUser *tgb.User, serial int, update tgb.Update) error
	AddNewUser(newUser *tgb.User, sex string, groupTitle string, groupId int64) error
	GetJubileeUsers() (jubUsers []model.GroupUser, err error)
	GetAllJubileeUsers() (jubUsers []model.GroupUser, err error)
	GetAnecdote() (text string, err error)
	AddAnecdote(text string) error
	AddSex(name, sex string) error
	DbClose() error
}

func TrimSymbolsFromSlice(s []string, cfg *config.Config) (words []string, err error) {

	var messageUpd []string

	for _, k := range s {

		k = strings.Trim(k, cfg.MsgText.MsgTrimSymbol)
		messageUpd = append(messageUpd, k)
	}

	words = messageUpd

	return words, nil
}

func (l *list) GoodMorningMessage() (string, error) {

	var goodMorning model.GoodMorningMSG
	var goodMorningMsgs []model.GoodMorningMSG

	rows, err := l.db.Query("SELECT * FROM botMessageGoodMorning")
	if err != nil {
		return "", err
	}

	for rows.Next() {
		err = rows.Scan(&goodMorning.ID, &goodMorning.Text)
		if err != nil {
			return "", err
		}
		goodMorningMsgs = append(goodMorningMsgs, goodMorning)
	}

	if len(goodMorningMsgs) > 0 {

		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(len(goodMorningMsgs)) + 1

		return goodMorningMsgs[num].Text, nil
	}

	return "", nil
}

func (l *list) Compliment(sex string) (compliment string, err error) {

	var compText model.ComplimentMSG
	var compAllText []model.ComplimentMSG
	var compForSexText []string

	log.Println("запрос таблицы")
	rows, err := l.db.Query("SELECT * FROM botMessageCompliment")
	if err != nil {
		log.Println("ошибка запроса botMessageCompliment")
	}

	for rows.Next() {
		err = rows.Scan(&compText.ID, &compText.Text, &compText.Sex)
		compAllText = append(compAllText, compText)
	}

	for _, v := range compAllText {

		if v.Sex == sex {
			compForSexText = append(compForSexText, v.Text)
		}
	}

	var num int
	if len(compForSexText) != 0 {

		rand.Seed(time.Now().UnixNano())
		l.logger.Infof("len compliment map %d", len(compForSexText))
		num = rand.Intn(len(compForSexText))
		return compForSexText[num], nil
	}

	l.logger.Infof("len compliment map %d", len(compForSexText))

	return "", nil

}

func (l *list) CheckSex(un string) (haveUser bool, s string, err error) {

	var newUser model.GroupUser
	var newUsers []model.GroupUser

	rows, err := l.db.Query("SELECT (id, user_name, sex) FROM group_users")
	if err != nil {
		return
	}

	for rows.Next() {

		err = rows.Scan(&newUser.ID, &newUser.UserName, &newUser.Sex)

		if newUser.UserName == un {

			log.Println("user found in db!")
			return true, newUser.Sex, nil
		}

		newUsers = append(newUsers, newUser)
	}

	return false, "", errors.New("user not found")
}

func (l *list) AddSex(name, sex string) error {

	_, err := l.db.Exec(fmt.Sprintf("UPDATE group_users SET ('sex') = ('%s') WHERE user_name = ('%s')", sex, name))
	if err != nil {
		return err
	}

	return nil
}

func (l *list) CheckBadWords(badList []string) (clearList []string, haveBadWords bool, err error) {

	var badWords []string
	var badWord string
	haveBadWords = false

	rows, err := l.db.Query("SELECT (word) FROM bad_words")
	if err != nil {
		return nil, false, err
	}

	for rows.Next() {
		err = rows.Scan(&badWord)
		badWords = append(badWords, badWord)
	}

	for _, word := range badList {

		for _, bad := range badWords {

			if word == bad {

				l.logger.Infof("найдено совпадение матерного слова в базе: %s", word)
				haveBadWords = true
			}
		}
	}

	return clearList, haveBadWords, err

}

func (l *list) AddBadWord(word string) (bool, error) {

	var badWord model.BadWords
	var haveWord = false

	rows, err := l.db.Query("SELECT * FROM bad_words")
	if err != nil {
	}

	for rows.Next() {
		err = rows.Scan(&badWord.ID, &badWord.Word)
		if badWord.Word == word {
			haveWord = true
			return true, nil
		}
	}

	if !haveWord {

		_, err = l.db.Exec(fmt.Sprintf("INSERT INTO bad_words (word) VALUES ('%s')", word))
		if err != nil {
			return false, errors.New("ошибка добавления матерного слова в базу")

		} else {
			return true, errors.New("новое матерное слово занесено в базу")
		}
	}

	return true, errors.New("added")

}

func (l *list) AddModeratorsGroup(group int64, title string) (haveGroup bool, modGroups []model.ModeratorsGroup, err error) {

	var modGroup model.ModeratorsGroup
	haveGroup = false

	rows, err := l.db.Query("SELECT * FROM moderators")
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {

		err = rows.Scan(&modGroup.ID, &modGroup.ModerGroupID, &modGroup.ModerGroupTitle, &modGroup.UserGroupID, &modGroup.UserGroupTitle)
		modGroups = append(modGroups, modGroup)
	}

	for _, grp := range modGroups {

		if grp.ModerGroupID == group {

			haveGroup = true
			return haveGroup, modGroups, errors.New("have group")
		}
	}

	if !haveGroup && group != 0 {

		_, err = l.db.Exec(fmt.Sprintf("INSERT INTO moderators (moder_group_id, moder_group_title, user_group_id , user_group_title) VALUES ('%d', '%s', '0', 'Без пользователей.')", group, title))
		if err != nil {
			log.Println(err)
		}

		haveGroup = true
	}

	return haveGroup, modGroups, nil
}

func (l *list) GetModeratorsGroup() (groups []model.ModeratorsGroup, err error) {

	rows, err := l.db.Query("SELECT * FROM moderators")
	if err != nil {
		log.Println(err)
	}

	var group model.ModeratorsGroup

	for rows.Next() {

		err = rows.Scan(&group.ID, &group.ModerGroupID, &group.ModerGroupTitle, &group.UserGroupID, &group.UserGroupTitle)
		groups = append(groups, group)
	}

	return groups, nil

}

func (l *list) AddUserGroupList(moderGroup, userGroup int64, moderTitle, userTitle string) (bool, error) {

	var moderatorGroup model.ModeratorsGroup
	var moderatorGroups []model.ModeratorsGroup
	var haveGroup = false

	rows, err := l.db.Query("SELECT * FROM moderators")
	if err != nil {
		return false, err
	}

	for rows.Next() {

		err := rows.Scan(&moderatorGroup.ID, &moderatorGroup.ModerGroupID, &moderatorGroup.ModerGroupTitle, &moderatorGroup.UserGroupID, &moderatorGroup.UserGroupTitle)
		if err != nil {
			return false, err
		}

		moderatorGroups = append(moderatorGroups, moderatorGroup)
	}

	for _, group := range moderatorGroups {

		if group.ModerGroupID == moderGroup && group.UserGroupID == userGroup {

			haveGroup = true
			return true, nil

		} else if group.ModerGroupID == moderGroup && group.UserGroupID == 0 {

			_, err = l.db.Exec(fmt.Sprintf("UPDATE moderators SET (user_group_id, user_group_title) = ('%d', '%s') WHERE moder_group_id = ('%d')", userGroup, userTitle, moderGroup))
			if err != nil {
				return false, err
			}

			haveGroup = true
			return false, nil
		}
	}

	if !haveGroup {

		_, err = l.db.Exec(fmt.Sprintf("INSERT INTO moderators (moder_group_id, moder_group_title, user_group_id, user_group_title) VALUES ('%d', '%s', '%d', '%s')", moderGroup, moderTitle, userGroup, userTitle))
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (l *list) AddNewJubileeUser(newUser *tgb.User, serial int, update tgb.Update) error {

	t := time.Now().Local().Format(config.StructDateTimeFormat)

	_, err := l.db.Exec(fmt.Sprintf("INSERT INTO group_users (serial, user_id, user_name, user_nick, time, group_name, group_id) VALUES ('%d', '%d', '%s', '%s', '%s', '%s', '%d')", serial, newUser.ID, newUser.FirstName,
		newUser.UserName, t, update.Message.Chat.Title, update.Message.Chat.ID))

	if err != nil {
		l.logger.Error(err)
	}

	return nil
}

func (l *list) AddNewUser(newUser *tgb.User, sex string, groupTitle string, groupId int64) error {

	t := time.Now().Local().Format(config.StructDateTimeFormat)

	_, err := l.db.Exec(fmt.Sprintf("INSERT INTO group_users (serial, user_id, user_name, user_nick, time, group_name, group_id, sex) VALUES ('0', '%d', '%s', '%s', '%s', '%s', '%d', '%s')", newUser.ID, newUser.FirstName,
		newUser.UserName, t, groupTitle, groupId, sex))

	if err != nil {
		l.logger.Error(err)
	}

	return nil

}

func (l *list) GetJubileeUsers() (users []model.GroupUser, err error) {

	var user model.GroupUser

	rows, err := l.db.Query("SELECT * FROM group_users ORDER BY id DESC LIMIT 3 ")
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		err = rows.Scan(&user.ID, &user.Serial, &user.UserID, &user.UserName, &user.UserNick,
			&user.Time, &user.GroupName, &user.GroupID, &user.Sex, &user.Rating)
		users = append(users, user)
	}

	return users, nil

}

func (l *list) GetAllJubileeUsers() (jubUsers []model.GroupUser, err error) {

	var user model.GroupUser
	var users []model.GroupUser

	rows, err := l.db.Query("SELECT * FROM group_users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		err = rows.Scan(&user.ID, &user.Serial, &user.UserID, &user.UserName, &user.UserNick,
			&user.Time, &user.GroupName, &user.GroupID, &user.Sex, &user.Rating)
		users = append(users, user)
	}

	return users, nil

}

func (l *list) GetAnecdote() (text string, err error) {

	var s model.AnecdoteStruct
	var ss []model.AnecdoteStruct

	rows, err := l.db.Query("SELECT * FROM anecdote")
	if err != nil {
		return text, err
	}

	for rows.Next() {
		err = rows.Scan(&s.ID, &s.Cat, &s.Text)
		ss = append(ss, s)
	}

	if len(ss) > 0 {

		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(len(ss))
		return ss[num].Text, nil
	}

	return text, nil

}

func (l *list) AddAnecdote(text string) error {

	_, err := l.db.Exec(fmt.Sprintf("INSERT INTO anecdote VALUES ('%d', '%s')", 21, text))
	if err != nil {
		return err
	}

	return nil
}

func (l *list) DbClose() error {
	err := l.db.Close()
	if err != nil {
		return err
	}

	return nil
}
