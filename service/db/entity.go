package db

import "time"

type User struct {
	ID         int64
	Username   string
	password   string
	InsertTime time.Time
}

type UserApp struct {
	ID         int64
	UserID     int64
	AppName    string
	Host       string
	InsertTime time.Time
}

type AppApi struct {
	ID         int64
	UserID     int64
	UserAppID  int64
	ApiID      int64
	InsertTime time.Time
}

type Api struct {
	ID         int64
	URL        string
	State      int32
	InsertTime time.Time
}
