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
Passing the http.Request pointer to Msg() renders a function to do the
key-to-translation lookup:
	translation := Msg(r)("Hi")

You could include the function returned by Msg() to the FuncMap of your
template:
	template.FuncMap{
		"Msg": msg.Msg(r),
	},
And then use the mapped Msg function inside the template:
	{{Msg "Hi"}} {{.name}}

If no translation is found, the message Key is used as a fallback.

Environment variables:
MSG_DEFAULT: determines the default language to use, if no translation is found
matching the Accept-Language header. The default value for MSG_DEFAULT is "en".
GO_ENV: if not set to "production", then for testing purpoeses, translations
that used the default language get prepended with "D-", and failed translations,
that used the Key fallback, get prepended with "X-".
*/
package msg

import (
	"github.com/wscherphof/env"
	"net/http"
	"strings"
)

var (
	production      = (env.Get("GO_ENV", "") == "production")
	defaultLanguage = &LanguageType{}
)

func init() {
	defaultLanguage.Parse(env.Get("MSG_DEFAULT", "en"))
}

// MessageType holds the translations for a message key.
type MessageType map[string]string

// Set stores the translation of the message for the given language. Any old
// value is overwritten.
func (m MessageType) Set(language, translation string) MessageType {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]MessageType, 500)

// NumLang sets the initial capacity for translations in a new message.
var NumLang = 10

// Key returns the message stored under the given key, if it doesn't exist yet,
// it gets created.
func Key(key string) (message MessageType) {
	if m, ok := messageStore[key]; ok {
		message = m
	} else {
		message = make(MessageType, NumLang)
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

func (l *LanguageType) Parse(s string) {
	parts := strings.Split(s, "-")
	l.Full = s
	l.Main = parts[0]
	if len(parts) > 1 {
		l.Sub = parts[1]
	}
	return
}

var languageCache = make(map[string][]*LanguageType, 100)

func headerLangs(r *http.Request) []*LanguageType {
	acceptLanguage := strings.ToLower(r.Header.Get("Accept-Language"))
	if cached, ok := languageCache[acceptLanguage]; ok {
		return cached
	}
	langStrings := strings.Split(acceptLanguage, ",")
	languageCache[acceptLanguage] = make([]*LanguageType, len(langStrings))
	for i, v := range langStrings {
		langString := strings.Split(v, ";")[0] // cut the q parameter
		lang := &LanguageType{}
		lang.Parse(langString)
		languageCache[acceptLanguage][i] = lang
	}
	return languageCache[acceptLanguage]
}

// Language provides the first language in the "Accept-Language" header in the
// given http request.
func Language(r *http.Request) (language *LanguageType) {
	languages := headerLangs(r)
	if len(languages) > 0 {
		language = languages[0]
	} else {
		language = defaultLanguage
	}
	return
}

func translate(key string, language *LanguageType) (translation string) {
	if val, ok := messageStore[key][language.Full]; ok {
		translation = val
	} else if val, ok := messageStore[key][language.Sub]; ok {
		translation = val
	} else if val, ok := messageStore[key][language.Main]; ok {
		translation = val
	}
	return
}

// Msg looks up the translation for a message, using the
// language matching the Accept-Language header in the request.
func Msg(r *http.Request, key string) (translation string) {
	languages := headerLangs(r)
	if key == "" {
		return ""
	}
	for _, language := range languages {
		if translation = translate(key, language); translation != "" {
			return
		}
	}
	if translation = translate(key, defaultLanguage); translation != "" {
		if !production {
			translation = "D-" + translation
		}
	} else {
		translation = key
		if !production {
			translation = "X-" + translation
		}
	}
	return
}
