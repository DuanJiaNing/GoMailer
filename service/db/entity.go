package db

import "time"

type User struct {
	ID         int64
	InsertTime time.Time

	Username string
	password string
}

type UserApp struct {
	ID         int64
	InsertTime time.Time

	UserID int64

	AppName string
	Host    string
}

type EndPoint struct {
	ID         int64
	InsertTime time.Time

	UserAppID       int64
	DialerID        int64
	TemplateID      int64
	UserID          int64
	DeliverStrategy int32
	EnableReCaptcha bool

	Name string
}

type Receiver struct {
	ID         int64
	InsertTime time.Time

	EndPointID int64

	Address      string
	ReceiverType int32
}

type Dialer struct {
	ID         int64
	InsertTime time.Time

	UserID int64

	Host         string
	Port         int32
	AuthUsername string
	AuthPassword string

	Name string
}

type Template struct {
	ID         int64
	InsertTime time.Time

	UserID int64

	Template    string
	ContentType string
}
