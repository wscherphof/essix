/*
Package msg manages translations of text labels ("messages") in a web application.

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

MSG_DEFAULT: determines the default language to use, if no translation is found
matching the Accept-Language header. The default value for MSG_DEFAULT is "en".

GO_ENV: if not set to "production", then translations that resorted to the
default language get prepended with "D-", and failed translations, falling back
to the message key, get prepended with "X-".

Messages and Translators are stored in memory. Translators are cached on their
Accept-Language header value.

Messages can also be used as multi-language text fields in data records:
	type Entity struct {
		Label msg.MessageType
	}
	entity := &Entity{Label: msg.New()}
	entity.Label.Set("en", "entity")
	entity.Label.Set("nl", "entiteit")
	...
	t := msg.Translator(r)
	label := t.Select(entity.Label)
*/
package msg

import (
	"github.com/wscherphof/env"
	"net/http"
	"os"
	"strings"
)

var (
	production      = (env.Get("GO_ENV", "") == "production")
	defaultLanguage = &languageType{}
)

func init() {
	defaultLanguage.parse(env.Get("MSG_DEFAULT", "en"))
}

/*
MessageType hold the translations for a message.
*/
type MessageType map[string]string

/*
New initialises a new MessageType.
*/
func New() MessageType {
	return make(MessageType, NumLang)
}

/*
Set stores the translation of the message for the given language. Any old
value is overwritten.
*/
func (m MessageType) Set(language, translation string) MessageType {
	language = strings.ToLower(language)
	m[language] = translation
	return m
}

var messageStore = make(map[string]MessageType, 500)

/*
NumLang sets the initial capacity for translations in a new message.
*/
var NumLang = 10

/*
Key returns the message stored under the given key, if it doesn't exist yet,
it gets created.
*/
func Key(key string) (message MessageType) {
	if m, ok := messageStore[key]; ok {
		message = m
	} else {
		message = New()
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

func (l *languageType) test(f func(lang string) (ok bool)) (ok bool) {
	for _, lang := range [3]string{l.Full, l.Main, l.Sub} {
		if f(lang) {
			return true
		}
	}
	return false
}

/*
TranslatorType knows about translations for the user's accepted languages.
*/
type TranslatorType struct {
	languages []*languageType
	files     map[string]string
}

var translatorCache = make(map[string]*TranslatorType, 100)

/*
Translator returns a (cached) TranslatorType.
*/
func Translator(r *http.Request) (t *TranslatorType) {
	acceptLanguage := strings.ToLower(r.Header.Get("Accept-Language"))
	if cached, ok := translatorCache[acceptLanguage]; ok {
		return cached
	}
	if acceptLanguage == "" {
		t = &TranslatorType{
			languages: make([]*languageType, 1),
			files:     make(map[string]string, 20),
		}
		t.languages[0] = defaultLanguage
	} else {
		langStrings := strings.Split(acceptLanguage, ",")
		t = &TranslatorType{
			languages: make([]*languageType, len(langStrings) + 1),
			files:     make(map[string]string, 20),
		}
		for i, v := range langStrings {
			langString := strings.Split(v, ";")[0] // cut the q parameter
			lang := &languageType{}
			lang.parse(langString)
			t.languages[i] = lang
		}
		t.languages[len(langStrings)] = defaultLanguage
	}
	translatorCache[acceptLanguage] = t
	return
}

/*
Get returns the translation for the message with the given key.
*/
func (t *TranslatorType) Get(key string) (translation string) {
	if key == "" {
		return ""
	}
	message := messageStore[key]
	return t.Select(message, key)
}

/*
Select returns the translation for a message.
*/
func (t *TranslatorType) Select(message MessageType, opt_default ...string) (translation string) {
	translate := func (lang string) (ok bool) {
		translation, ok = message[lang]
		return
	}
	for _, language := range t.languages {
		if language.test(translate) {
			if !production && language == defaultLanguage {
				translation = "D-" + translation
			}
			return
		}
	}
	if len(opt_default) == 1 {
		translation = opt_default[0]
		if !production {
			translation = "X-" + translation
		}
	}
	return
}

/*
File searches for an "inner" template fitting the "base" template, matching
the user's accepted languages.

Template names are without file name extension. The default extension is ".ace".

Example: if MSG_DEFAULT is "en", and the Accept-Languages header is empty,
	msg.File("/resources/templates", "home", "HomePage", ".tpl")
returns
	"HomePage-en", nil
if the file "/resources/templates/home/HomePage-en.tpl" exists.
*/
func (t *TranslatorType) File(location, dir, base string, opt_extension ...string) (inner string, err error) {
	extension := ".ace"
	if len(opt_extension) == 1 {
		extension = opt_extension[0]
	}
	template := location + "/" + dir + "/" + base
	if cached, ok := t.files[template]; ok {
		return cached, nil
	}
	var lang string
	for _, language := range t.languages {
		if lang, err = exists(template, extension, language); err == nil {
			inner = base + "-" + lang
			t.files[template] = inner
			return
		}
	}
	return
}

func exists(template, extension string, language *languageType) (lang string, err error) {
	language.test(func(l string) (ok bool){
		if err = stat(template, extension, l); err == nil {
			lang, ok = l, true
		}
		return
	})
	return
}

func stat(template, extension, lang string) (err error) {
	path := template + "-" + lang + extension
	_, err = os.Stat(path)
	return
}
