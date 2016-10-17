package messages

import (
	msg "github.com/wscherphof/essix/messages"
)

func init() {
	msg.Key("Example").
		Set("nl", "Voorbeeld").
		Set("en", "Example")
}
