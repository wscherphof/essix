package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("ErrDuplicatePrimaryKey").
		Set("nl", `Er bestaat al een record met deze gegevens`).
		Set("en", `A record with these data already exists.`)
}
