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
		template.Run(w, r, "secure", "AccountForm", "", map[string]interface{}{
			"RateLimitToken": token,
		})
	}
}

func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.NewAccount(r.FormValue("email"),
		r.FormValue("pwd1"), r.FormValue("pwd2"),
	); err != nil {
		template.Error(w, r, err, conflict)
	} else if err, remark := sendEmail(r, account.Email,
		"Activate",
		"/account/activate?code="+account.ActivateCode+"&id="+account.ID,
	); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "secure", "NewAccount", "", map[string]interface{}{
			"id":     account.ID,
			"remark": remark,
		})
	}
}
