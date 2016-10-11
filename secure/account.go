package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/essix/router"
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
	switch r.Method {
	case "GET":
		template.Run(w, r, "account", "NewAccount", "", map[string]interface{}{
			"id":     r.FormValue("id"),
			"remark": r.FormValue("remark"),
		})
	case "POST":
		if account, err, conflict := model.NewAccount(
			r.FormValue("email"),
			r.FormValue("pwd1"),
			r.FormValue("pwd2"),
		); err != nil {
			template.Error(w, r, err, conflict)
		} else if err, remark := activateEmail(r, account); err != nil {
			template.Error(w, r, err, false)
		} else {
			router.Redirect(w, r, map[string]string{
				"id":     account.ID,
				"remark": remark,
			})
		}
	}
}

func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	account := Authentication(w, r)
	template.Run(w, r, "account", "Account", "", map[string]interface{}{
		"email": account.Email,
	})
}
