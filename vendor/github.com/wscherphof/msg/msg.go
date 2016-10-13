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

The user's language is determined from the "Accept-Language" request header.
Pass the http.Request pointer to Translator():
	t := msg.Translator(r)
Then get the translation:
	message := t.Get("Hi")

Environment variables:
- MSG_DEFAULT: determines the default language to use, if no translation is found
matching the Accept-Language header. The default value for MSG_DEFAULT is "en".
- GO_ENV: if not set to "production", then translations that resorted to the
default language get prepended with "D-", and failed translations, falling back
to the message key, get prepended with "X-".

Messages and Translators are stored in memory. Translators are cached on their
Accept-Language header value.
*/
package msg

import (
	"github.com/wscherphof/env"
	"net/http"
	"strings"
)

var (
	production      = (env.Get("GO_ENV", "") == "production")
	defaultLanguage = &languageType{}
)

func init() {
	defaultLanguage.parse(env.Get("MSG_DEFAULT", "en"))
}

type messageType map[string]string

// Set stores the translation of the message for the given language. Any old
// value is overwritten.
func (m messageType) Set(language, translation string) messageType {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]messageType, 500)

// NumLang sets the initial capacity for translations in a new message.
var NumLang = 10

// Key returns the message stored under the given key, if it doesn't exist yet,
// it gets created.
func Key(key string) (message messageType) {
	if m, ok := messageStore[key]; ok {
		message = m
	} else {
		message = make(messageType, NumLang)
		messageStore[key] = message
	}
	return
}

type languageType struct {
	// e.g. "en-gb"
	Full string
	// e.g. "en"
	Main string
	// e.g. "gb"
	Sub string
}

func (l *languageType) parse(s string) {
	parts := strings.Split(s, "-")
	l.Full = s
	l.Main = parts[0]
	if len(parts) > 1 {
		l.Sub = parts[1]
	}
	return
}

var translatorCache = make(map[string]*translatorType, 100)

type translatorType struct {
	languages []*languageType
}

// Translator returns an object that knows how to lookup the translation for a
// message.
func Translator(r *http.Request) *translatorType {
	acceptLanguage := strings.ToLower(r.Header.Get("Accept-Language"))
	if cached, ok := translatorCache[acceptLanguage]; ok {
		return cached
	}
	langStrings := strings.Split(acceptLanguage, ",")
	t := &translatorType{make([]*languageType, len(langStrings))}
	for i, v := range langStrings {
		langString := strings.Split(v, ";")[0] // cut the q parameter
		lang := &languageType{}
		lang.parse(langString)
		t.languages[i] = lang
	}
	translatorCache[acceptLanguage] = t
	return t
}

func translate(key string, language *languageType) (translation string) {
	if val, ok := messageStore[key][language.Full]; ok {
		translation = val
	} else if val, ok := messageStore[key][language.Sub]; ok {
		translation = val
	} else if val, ok := messageStore[key][language.Main]; ok {
		translation = val
	}
	return
}

// Get returns the translation for a message.
func (t *translatorType) Get(key string) (translation string) {
	if key == "" {
		return ""
	}
	for _, language := range t.languages {
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

// Language provides the first language in the "Accept-Language" header in the
// given http request.
func (t *translatorType) Language() (language *languageType) {
	if len(t.languages) > 0 {
		language = t.languages[0]
	} else {
		language = defaultLanguage
	}
	return
}
