package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/essix/util"
	"github.com/wscherphof/msg"
	"github.com/wscherphof/secure"
	"net/http"
)

func PasswordTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "password", "PasswordTokenForm")
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		t.Set("ratelimit", token)
		t.Run()
	}
}

func PasswordToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "password", "PasswordToken"); t == nil {
		return
	} else if account, err, conflict := model.GetAccount("", r.FormValue("email")); err != nil {
		template.Error(w, r, err, conflict)
	} else if !account.IsActive() {
		template.Error(w, r, model.ErrNotActivated, conflict)
	} else if err := account.CreatePasswordToken(); err != nil {
		template.Error(w, r, err, false)
	} else {
		format := msg.Msg(r)("Time format")
		expires := account.PasswordToken.Expires.Format(format)
		path := "/account/password"
		path += "?token=" + account.PasswordToken.Value
		path += "&id=" + account.ID
		path += "&expires=" + util.URLEncodeString(expires)
		if err, remark := sendEmail(r, account.Email, "PasswordToken", path, expires); err != nil {
			template.Error(w, r, err, false)
		} else {
			t.Set("remark", remark)
			t.Run()
		}
	}
}

func ChangePasswordForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := template.GET(w, r, "password", "ChangePasswordForm")
	id, token, expires, cancel := r.FormValue("id"), r.FormValue("token"), r.FormValue("expires"), r.FormValue("cancel")
	expires, _ = util.URLDecodeString(expires)
	if cancel == "true" {
		model.ClearPasswordToken(id, token)
		template.Run(w, r, "password", "ChangePassword-cancel", "", nil)
	} else {
		t.Set("id", id)
		t.Set("token", token)
		t.Set("expires", expires)
		t.Run()
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t := template.PRG(w, r, "password", "ChangePassword"); t == nil {
		return
	} else if account, err, conflict := model.GetAccount(
		r.FormValue("id"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, conflict := account.ChangePassword(
		r.FormValue("token"),
		r.FormValue("pwd1"),
		r.FormValue("pwd2"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else {
		secure.LogOut(w, r, false)
		t.Run()
	}
}
