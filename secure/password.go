package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
	"net/http"
)

func passwordEmail(r *http.Request, acc *model.Account) (error, string) {
	format := msg.Msg(r)("Time format")
	return sendEmail(r, acc.ID, acc.Name(), "password", acc.PasswordCode.Value, acc.PasswordCode.Expires.Format(format))
}

func PasswordCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "passwordcode", "", map[string]interface{}{
			"UID":            ps.ByName("uid"),
			"RateLimitToken": token,
		})
	}
}

func PasswordCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := r.FormValue("uid")
	if acc, e, conflict := model.GetAccountInsecure(uid); e != nil {
		template.Error(w, r, e, conflict)
	} else if !acc.IsActive() {
		template.Error(w, r, model.ErrNotActivated, conflict, "secure", "activation_resend")
	} else if e := acc.CreatePasswordCode(); e != nil {
		template.Error(w, r, e, false)
	} else if e, remark := passwordEmail(r, acc); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "passwordcode_success", "", map[string]interface{}{
			"Name":   acc.Name(),
			"Remark": remark,
		})
	}
}

func PasswordForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, code, extra, cancel := ps.ByName("uid"), r.FormValue("code"), r.FormValue("extra"), r.FormValue("cancel")
	expires, _ := util.URLDecode([]byte(extra))
	if cancel == "true" {
		model.ClearPasswordCode(uid, code)
		template.Run(w, r, "secure", "passwordcode_cancelled", "", nil)
	} else {
		template.Run(w, r, "secure", "password", "", map[string]interface{}{
			"UID":     uid,
			"Code":    code,
			"Expires": string(expires),
		})
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, code, pwd1, pwd2 := r.FormValue("uid"), r.FormValue("code"), r.FormValue("pwd1"), r.FormValue("pwd2")
	if acc, e, conflict := model.GetAccountInsecure(uid); e != nil {
		template.Error(w, r, e, conflict)
	} else if e, conflict := acc.ChangePassword(code, pwd1, pwd2); e != nil {
		template.Error(w, r, e, conflict)
	} else {
		secure.LogOut(w, r, false)
		template.Run(w, r, "secure", "password_success", "", nil)
	}
}
