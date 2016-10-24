package messages

import (
	"github.com/wscherphof/msg"
)

func init() {
	msg.Key("").
		Set("nl", "").
		Set("en", "")

	msg.Key("Hi").
		Set("nl", "Hoi").
		Set("en", "Hi")

	msg.Key("Current").
		Set("nl", "Huidig").
		Set("en", "Current")

	msg.Key("New").
		Set("nl", "Nieuw").
		Set("en", "New")

	msg.Key("change").
		Set("nl", "wijzigen").
		Set("en", "change")

	msg.Key("Confirm").
		Set("nl", "Bevestigen").
		Set("en", "Confirm")
}
