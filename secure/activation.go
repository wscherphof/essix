package secure

import (
	"github.com/dchest/captcha"
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/expeertise/model/account"
	"github.com/wscherphof/expeertise/router"
	"net/http"
)

func activationEmail(r *http.Request, acc *account.Account) (error, string) {
	return sendEmail(r, acc.UID, acc.Name(), "activation", acc.ActivationCode, "")
}

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	return router.Template("secure", "activation", "", map[string]interface{}{
		"UID":  ps.ByName("uid"),
		"Code": r.FormValue("code"),
	})(w, r, ps)
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	if acc, e, conflict := account.Activate(r.FormValue("uid"), r.FormValue("code")); e != nil {
		err = router.NewError(e, "secure", "activation")
		err.Conflict = conflict
	} else {
		router.Template("secure", "activation_success", "", map[string]interface{}{
			"Name": acc.Name(),
		})(w, r, ps)
	}
	return
}

func ActivationCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	return router.Template("secure", "activation_resend", "", map[string]interface{}{
		"UID":       ps.ByName("uid"),
		"CaptchaId": captcha.New(),
	})(w, r, ps)
}

func ActivationCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (err *router.Error) {
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		err = router.NewError(captcha.ErrNotFound)
		err.Conflict = true
	} else if acc, e, conflict := account.GetInsecure(r.FormValue("uid")); e != nil {
		err = router.NewError(e)
		err.Conflict = conflict
	} else if acc.IsActive() {
		err = router.NewError(account.ErrAlreadyActivated)
		err.Conflict = true
	} else if e, remark := activationEmail(r, acc); e != nil {
		err = router.NewError(e)
	} else {
		router.Template("secure", "activation_resend_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"UID":    acc.UID,
			"Remark": remark,
		})(w, r, ps)
	}
	return
}
