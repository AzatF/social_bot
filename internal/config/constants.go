package config

const StructDateTimeFormat = "2006-01-02 15:04"
const StructDateFormat = "2006-01-02"
const StructTimeFormat = "15:04"

const MakeUserGroupTable = "CREATE TABLE IF NOT EXISTS group_users " +
	"(id INTEGER PRIMARY KEY, " +
	"serial INTEGER NOT NULL DEFAULT 0, " +
	"user_id INTEGER NOT NULL, " +
	"user_name VARCHAR (30) NOT NULL, " +
	"user_nick VARCHAR (50) DEFAULT ('нет ника'), " +
	"time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
	"group_name VARCHAR (50) NOT NULL, " +
	"group_id INTEGER NOT NULL, " +
	"sex VARCHAR (10) DEFAULT ('none')," +
	"rate INTEGER DEFAULT 0)"

const MakeBotMessageGoodMorningTable = "CREATE TABLE IF NOT EXISTS botMessageGoodMorning" +
	"(id INTEGER PRIMARY KEY, message TEXT NOT NULL)"

const MakeBotMessageComplimentTable = "CREATE TABLE IF NOT EXISTS botMessageCompliment" +
	"(id INTEGER PRIMARY KEY, message TEXT NOT NULL, sex VARCHAR (10) NOT NULL)"

const MakeModeratorsGroupTable = "CREATE TABLE IF NOT EXISTS moderators " +
	"(id INTEGER PRIMARY KEY, moder_group_id INTEGER NOT NULL, moder_group_title TEXT DEFAULT no_title, " +
	"user_group_id INTEGER DEFAULT 0, user_group_title TEXT DEFAULT no_title)"

const MakeBadWordsTable = "CREATE TABLE IF NOT EXISTS bad_words " +
	"(id INTEGER PRIMARY KEY, word VARCHAR (30) NOT NULL)"

const MakeAnecdoteTable = "CREATE TABLE IF NOT EXISTS anecdote " +
	"(id INTEGER PRIMARY KEY, cat INTEGER NOT NULL,  text TEXT NOT NULL)"
