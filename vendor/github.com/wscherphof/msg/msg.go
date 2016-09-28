/*
Package msg provides a means to manage translations of text labels ("messages")
in a web application.

New messages are defined like this:
	msg.Key("Hello").
	  Set("en", "Hello, world").
	  Set("nl", "Hallo wereld")
	msg.Key("Hi").
	  Set("en", "Hi").
	  Set("nl", "Hoi")

When you ask for the translation of a certain message key, the user's language
is determined from the "Accept-Language" request header.
Passing the httprequest pointer to Msg() renders a function to do the
key-to-translation lookup:
	translation := Msg(r)("Hi")

You could include the function returned by Msg() to the FuncMap of your
template:
	template.FuncMap{
		"Msg": msg.Msg(r),
	},
And then use the mapped Msg function inside the template:
	{{ Msg "Hi" }} {{ .Name }}
*/
package msg

import (
	"net/http"
	"strings"
)

// Message holds the translations for a message key.
type Message map[string]string

// Set stores the translation of the message for the given language. Any old
// value is overwritten.
func (m Message) Set(language, translation string) Message {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]Message, 500)

// NumLang sets the initial capacity for translations in a new message.
var NumLang = 2

// Key returns the message stored under the given key, if it doesn't exist yet,
// it gets created.
func Key(key string) (message Message) {
	if m, ok := messageStore[key]; ok {
		message = m
	} else {
		message = make(Message, NumLang)
		messageStore[key] = message
	}
	return
}

// LanguageType defines a language.
type LanguageType struct {
	// e.g. "en-gb"
	Full string
	// e.g. "en"
	Main string
	// e.g. "gb"
	Sub string
}

var languageCache = make(map[string]LanguageType, 100)

// Language provides the first language in the "Accept-Language" header in the
// given http request.
func Language(r *http.Request) (language LanguageType) {
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptLanguage = strings.ToLower(acceptLanguage)
	if lang, ok := languageCache[acceptLanguage]; ok {
		language = lang
	} else {
		firstLanguage := strings.Split(acceptLanguage, ",")[0] // cut other languages
		firstLanguage = strings.Split(firstLanguage, ";")[0]   // cut the q parameter
		parts := strings.Split(firstLanguage, "-")
		lang = LanguageType{
			Full: firstLanguage,
			Main: parts[0],
		}
		if len(parts) > 1 {
			lang.Sub = parts[1]
		}
		languageCache[acceptLanguage] = lang
		language = lang
	}
	return
}

// Msg returns a function that looks up the translation for a certain message
// key in the given language.
func Msg(lang LanguageType, key string) (value string) {
	if val, ok := messageStore[key][lang.Full]; ok {
		value = val
	} else if val, ok := messageStore[key][lang.Sub]; ok {
		value = val
	} else if val, ok := messageStore[key][lang.Main]; ok {
		value = val
	} else {
		value = "X-" + key
	}
	return
}
