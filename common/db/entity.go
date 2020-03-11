package db

import (
	"errors"
	"time"
)

type deliverStrategy string
type receiverType string
type mailState string

var (
	deliverStrategyName = []string{
		"DELIVER_IMMEDIATELY",
		"STAGING",
	}

	receiverTypeName = []string{
		"TO",
		"CC",
		"BCC",
	}

	mailStateName = []string{
		"STAGING",
		"DELIVER_SUCCESS",
		"DELIVER_FAILED",
	}
)

type User struct {
	ID         int64
	InsertTime time.Time

	Username string
	Password string
}

type UserApp struct {
	ID         int64
	InsertTime time.Time

	UserID int64

	Name string
	Host string
}

type EndPoint struct {
	ID         int64
	InsertTime time.Time

	UserAppID  int64
	DialerID   int64
	TemplateID int64
	UserID     int64

	Name string
}

type EndPointPreference struct {
	ID         int64
	InsertTime time.Time

	EndPointID int64

	DeliverStrategy string
	EnableReCaptcha bool
}

type Receiver struct {
	ID         int64
	InsertTime time.Time

	EndPointID int64
	UserID     int64
	UserAppID  int64

	Address      string
	ReceiverType string
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

type Mail struct {
	ID         int64
	InsertTime time.Time

	EndPointID   int64
	State        string
	DeliveryTime time.Time
	Content      string
}

func ReceiverType(name string) (receiverType, error) {
	for _, n := range receiverTypeName {
		if name == n {
			return receiverType(name), nil
		}
	}

	return "", errors.New("receiver type illegal")
}

func DeliverStrategy(name string) (deliverStrategy, error) {
	for _, n := range deliverStrategyName {
		if name == n {
			return deliverStrategy(name), nil
		}
	}

	return "", errors.New("deliver strategy illegal")
}

func MailState(name string) (mailState, error) {
	for _, n := range mailStateName {
		if name == n {
			return mailState(name), nil
		}
	}

	return "", errors.New("mail state illegal")
}

func (r receiverType) Name() string {
	return string(r)
}

func (r deliverStrategy) Name() string {
	return string(r)
}

func (r mailState) Name() string {
	return string(r)
}
