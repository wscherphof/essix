/*
Package messages defines translations for user facing texts.

See github.com/wscherphof/msg
*/
package messages

import (
	"github.com/wscherphof/msg"
)

/*
Key defines a new message, and returns it to call Set() to set a translations
for it.
*/
var Key = msg.Key

// Init just triggers the init function of the message files.
func Init() {
	return
}
