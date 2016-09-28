package messages

import (
	"github.com/wscherphof/msg"
)

func Key(key string) msg.Message {
	return msg.Key(key)
}

func init() {
	msg.Key("").
		Set("nl", "").
		Set("en", "")

	msg.Key("Hi").
		Set("nl", "Hoi").
		Set("en", "Hi")

	msg.Key("Time format").
		Set("nl", "2-1-2006, 15.04 u. (MST)").
		Set("en", "1/2/2006 at 3:04pm (MST)")

	msg.Key("Last modified").
		Set("nl", "Laatste wijziging").
		Set("en", "Last modified")

	msg.Key("Save changes").
		Set("nl", "Wijzigingen opslaan").
		Set("en", "Save changes")

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
