package captcha

import (
	"github.com/dchest/captcha"
	"github.com/wscherphof/msg"
)

func init() {
	msg.New("Captcha").
		Add("nl", "Typ de code zoals hieronder afgebeeld").
		Add("en", "Enter the code as depicted below")

	msg.New("Captcha image").
		Add("nl", "Afbeelding van captchacode").
		Add("en", "Image of captcha code")

	msg.New(captcha.ErrNotFound.Error()).
		Add("nl", "De 'captcha' code ontbreekt of is onjuist of verlopen").
		Add("en", "The 'captcha' code is missing or incorrect or expired")
}
