package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func AccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "account", "AccountForm", "", map[string]interface{}{
			"ratelimit": token,
		})
	}
}

func activateEmail(r *http.Request, account *model.Account) (error, string) {
	return sendEmail(r, account.Email,
		"Activate",
		"/account/activate?token="+account.ActivateToken+"&id="+account.ID,
	)
}

func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.NewAccount(r.FormValue("email"),
		r.FormValue("pwd1"), r.FormValue("pwd2"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := activateEmail(r, account); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "account", "NewAccount", "", map[string]interface{}{
			"id":     account.ID,
			"remark": remark,
		})
	}
}
