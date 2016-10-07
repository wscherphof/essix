package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func AccountForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, e := ratelimit.NewToken(r); e != nil {
		template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "AccountForm", "", map[string]interface{}{
			"RateLimitToken": token,
		})
	}
}

func NewAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if acc, e, conflict := model.NewAccount(r.FormValue("uid"), r.FormValue("pwd1"), r.FormValue("pwd2")); e != nil {
		template.Error(w, r, e, conflict)
		// } else if e, remark := activationEmail(r, acc); e != nil {
		// 	template.Error(w, r, e, false)
	} else {
		template.Run(w, r, "secure", "NewAccount", "", map[string]interface{}{
			"uid":    acc.ID,
			"name":   acc.Name(),
			"remark": "remark",
		})
	}
}
