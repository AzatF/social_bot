package model

import "time"

type BadWords struct {
	ID   int
	Word string
}

type GoodMorningMSG struct {
	ID   int
	Text string
}

type ComplimentMSG struct {
	ID   int
	Text string
	Sex  string
}

type ModeratorsGroup struct {
	ID              int
	ModerGroupID    int64
	ModerGroupTitle string
	UserGroupID     int64
	UserGroupTitle  string
}

type GroupUser struct {
	ID        int
	Serial    int
	UserID    int
	UserName  string
	UserNick  string
	Time      time.Time
	GroupName string
	GroupID   int64
	//Marked    int
	Sex    string
	Rating int
}

type AnecdoteStruct struct {
	ID   int
	Cat  int
	Text string
}
