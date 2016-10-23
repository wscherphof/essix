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
}
