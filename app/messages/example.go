package messages

import (
	msg "github.com/wscherphof/essix/messages"
)

func init() {
	msg.Key("Edit profile").
		Set("nl", "Mijn profiel").
		Set("en", "My profile details")

	msg.Key("Country").
		Set("nl", "Land").
		Set("en", "Country")

	msg.Key("Time zone").
		Set("nl", "Tijdzone").
		Set("en", "Time zone")

	msg.Key("First name").
		Set("nl", "Voornaam").
		Set("en", "First name")

	msg.Key("Last name").
		Set("nl", "Achternaam").
		Set("en", "Last Name")

	msg.Key("Time format").
		Set("nl", "2-1-2006, 15.04 u. (MST)").
		Set("en", "1/2/2006 at 3:04pm (MST)")

	msg.Key("Last modified").
		Set("nl", "Laatste wijziging").
		Set("en", "Last modified")

	msg.Key("Save changes").
		Set("nl", "Wijzigingen opslaan").
		Set("en", "Save changes")
}
