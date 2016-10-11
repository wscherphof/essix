package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func NewAccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "account", "NewAccountForm", "", map[string]interface{}{
			"ratelimit": token,
		})
	}
}

func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if prg := router.PRG(w, r, "account", "NewAccount", "", "id", "remark"); prg == nil {
		return
	} else if account, err, conflict := model.NewAccount(
		r.FormValue("email"),
		r.FormValue("pwd1"),
		r.FormValue("pwd2"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := activateEmail(r, account); err != nil {
		template.Error(w, r, err, false)
	} else {
		prg(account.ID, remark)
	}
}

func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	template.Run(w, r, "account", "Account", "", map[string]interface{}{
		"email": account.Email,
	})
}
