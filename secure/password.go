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
		email := template.Email(r, "password", "PasswordToken-email", "lang")
		format := msg.Msg(r)("Time format")
		expires := account.PasswordToken.Expires.Format(format)
		email.Set("expires", expires)
		link := "https://" + r.Host + "/account/password"
		link += "?token=" + account.PasswordToken.Value
		link += "&id=" + account.ID
		link += "&expires=" + util.URLEncodeString(expires)
		email.Set("link", link)
		if err, message := email.Run(account.Email, "Reset password"); err != nil {
			template.Error(w, r, err, false)
		} else {
			t.Set("message", message)
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
