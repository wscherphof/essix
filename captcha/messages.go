package captcha

import (
  "github.com/wscherphof/msg"
  "github.com/dchest/captcha"
)

func DefineMessages () {
  var m, a = msg.Init()

  m("Captcha")
  a("nl", "Typ de code zoals hieronder afgebeeld")
  a("en", "Enter the code as depicted below")

  m("Captcha image")
  a("nl", "Afbeelding van captchacode")
  a("en", "Image of captcha code")

  m(captcha.ErrNotFound.Error())
  a("nl", "De 'captcha' code ontbreekt of is onjuist")
  a("en", "The 'captcha' code is missing or incorrect")
}
