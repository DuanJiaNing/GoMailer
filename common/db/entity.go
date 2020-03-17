package db

import (
	"time"
)

type deliverStrategy string
type receiverType string
type mailState string

const (
	DeliverStrategy_DELIVER_IMMEDIATELY = deliverStrategy("DELIVER_IMMEDIATELY")
	DeliverStrategy_STAGING             = deliverStrategy("STAGING")

	MailState_STAGING         = mailState("STAGING")
	MailState_DELIVER_SUCCESS = mailState("DELIVER_SUCCESS")
	MailState_DELIVER_FAILED  = mailState("DELIVER_FAILED")

	ReceiverType_To  = receiverType("To")
	ReceiverType_Cc  = receiverType("Cc")
	ReceiverType_Bcc = receiverType("Bcc")
)

var (
	deliverStrategyName = []string{
		"DELIVER_IMMEDIATELY",
		"STAGING",
	}

	receiverTypeName = []string{
		"To",
		"Cc",
		"Bcc",
	}

	mailStateName = []string{
		"STAGING",
		"DELIVER_SUCCESS",
		"DELIVER_FAILED",
	}
)

type User struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	Username string
	Password string
}

type UserApp struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	UserId int64

	Name string
	Host string
}

type Endpoint struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	UserAppId  int64
	DialerId   int64
	TemplateId int64
	UserId     int64

	Name string
}

type EndpointPreference struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	EndpointId int64

	DeliverStrategy string
	EnableReCaptcha int32 // 1 enable 2 disable
}

type Receiver struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	EndpointId int64
	UserId     int64
	UserAppId  int64

	Address      string
	ReceiverType string
}

type Dialer struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	UserId int64

	Host         string
	Port         int
	AuthUsername string
	AuthPassword string

	Name string
}

type Template struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	UserId int64

	Template    string
	ContentType string
}

type Mail struct {
	Id         int64
	InsertTime time.Time `xorm:"created"`

	EndpointId   int64
	State        string
	DeliveryTime time.Time
	Content      string
}

func ReceiverType(name string) receiverType {
	for _, n := range receiverTypeName {
		if name == n {
			return receiverType(name)
		}
	}

	return ""
}

func DeliverStrategy(name string) deliverStrategy {
	for _, n := range deliverStrategyName {
		if name == n {
			return deliverStrategy(name)
		}
	}

	return ""
}

func MailState(name string) mailState {
	for _, n := range mailStateName {
		if name == n {
			return mailState(name)
		}
	}

	return ""
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
