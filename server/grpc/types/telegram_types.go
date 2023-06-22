package types

import (
	"time"
)

type TelegramLoginData struct {
	AuthDate  time.Time
	UserId    int64
	Firstname string
	Lastname  string
	Username  string
	PhotoURL  string
}
