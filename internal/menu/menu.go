package menu

import tgb "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const UserMenu = "         Список слов на которые бот реагирует:  🛠  \n \n " +
	"1. Привет, Здравствуйте и т.п. приветствия.\n " +
	"2. Доброе утро.\n " +
	"3. Спокойной ночи.\n " +
	"4. В работе\n " +
	"5. У бота очень тонкая душевная организация, он почти падает в обморок от нецензурных выражений.\n" +
	"Смотрит с укором, качает головой."

const ComMenu = "     Список доступных вам команд:  🛠  \n \n" +
	"✅ `addmoderatorgroup` + номер _(добавление группы модераторов)._\n\n" +
	"✅ `add-moder-group` _(отправляет запрос из группы где есть бот, на добавление этой группы в список групп модераторов)._\n\n" +
	"✅ `add-moder-user-link` _(связывает группу модераторов и пользователей. Введите команду, " +
	"затем через пробел номер группы модераторов, затем через пробел номер группы пользователей. " +
	"Внимательно проверьте правильность номеров групп)._\n\n" +
	"✅ `chatinfo` _(отправляет информацию в админку о имени и ID группы, откуда отправляется команда.\n" +
	"Номер можно скопировать нажатием, для удобства добавления связи групп модераторов и пользователей.\n" +
	"Сообщение будет удалено из группы отправителя, если бот админ)_\n\n" +
	"✅ *Мат + слово* _(Слово будет добавлено в базу)._"

var NumericKeyboard = tgb.NewInlineKeyboardMarkup(
	tgb.NewInlineKeyboardRow(button1),
	tgb.NewInlineKeyboardRow(button2),
	tgb.NewInlineKeyboardRow(button8),
	tgb.NewInlineKeyboardRow(button3),
	tgb.NewInlineKeyboardRow(button9),
	tgb.NewInlineKeyboardRow(Button01),
	tgb.NewInlineKeyboardRow(Button02),
	tgb.NewInlineKeyboardRow(Button03),
	tgb.NewInlineKeyboardRow(Button12),
)

var UserKeyboard = tgb.NewInlineKeyboardMarkup(
	tgb.NewInlineKeyboardRow(Button01),
	tgb.NewInlineKeyboardRow(Button03),
	tgb.NewInlineKeyboardRow(Button02),
	tgb.NewInlineKeyboardRow(button9),
	tgb.NewInlineKeyboardRow(Button12),
)

var Button01 = tgb.NewInlineKeyboardButtonData("Комплимент", "compliment")
var Button02 = tgb.NewInlineKeyboardButtonData("Список слов на которые бот отвечает", "bot_say")
var Button03 = tgb.NewInlineKeyboardButtonData("Анекдот", "anecdote")
var Button04 = tgb.NewInlineKeyboardButtonData("", "")
var Button05 = tgb.NewInlineKeyboardButtonData("", "")
var Button06 = tgb.NewInlineKeyboardButtonData("", "")

var SexMenu = tgb.NewInlineKeyboardMarkup(
	tgb.NewInlineKeyboardRow(MaleButton),
	tgb.NewInlineKeyboardRow(FemaleButton),
)

var button1 = tgb.NewInlineKeyboardButtonData("Список команд", "com_list")
var button2 = tgb.NewInlineKeyboardButtonData("Список новых пользователей (крайние трое)", "jubilee_list")
var button8 = tgb.NewInlineKeyboardButtonData("Весь список новых пользователей ", "all_jubilee_list")
var button3 = tgb.NewInlineKeyboardButtonData("Список групп модераторов и пользователей.", "moderator_group_list")
var button9 = tgb.NewInlineKeyboardButtonData("Памятка.", "moderator_member")

var Button4 = tgb.NewInlineKeyboardButtonData("Добавить группу", "add_new_mod")
var Button5 = tgb.NewInlineKeyboardButtonData("Да, я уверен!", "add_new_mod_true")

var Button11 = tgb.NewInlineKeyboardButtonData("Отклонить", "remove_button")
var Button12 = tgb.NewInlineKeyboardButtonData("Закрыть", "remove_button")
var MaleButton = tgb.NewInlineKeyboardButtonData("Мужской", "male")
var FemaleButton = tgb.NewInlineKeyboardButtonData("Женский", "female")
