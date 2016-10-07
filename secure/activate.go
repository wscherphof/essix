package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	// "github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	template.Run(w, r, "secure", "ActivateForm", "", map[string]interface{}{
		"id":   r.FormValue("id"),
		"code": r.FormValue("code"),
	})
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.GetAccount(r.FormValue("id")); err != nil {
		template.Error(w, r, err, conflict, "secure", "Activate")
	} else if err, conflict = account.Activate(r.FormValue("code")); err != nil {
		template.Error(w, r, err, conflict, "secure", "Activate")
	} else {
		template.Run(w, r, "secure", "Activate", "", nil)
	}
}

// func ActivationCodeForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	if token, e := ratelimit.NewToken(r); e != nil {
// 		template.Error(w, r, e, false)
// 	} else {
// 		template.Run(w, r, "secure", "activation_resend", "", map[string]interface{}{
// 			"UID":            ps.ByName("uid"),
// 			"RateLimitToken": token,
// 		})
// 	}
// }

// func ActivationCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	if acc, e, conflict := model.GetAccountInsecure(r.FormValue("uid")); e != nil {
// 		template.Error(w, r, e, conflict)
// 	} else if acc.IsActive() {
// 		template.Error(w, r, model.ErrAlreadyActivated, true)
// 	} else if e, remark := activationEmail(r, acc); e != nil {
// 		template.Error(w, r, e, false)
// 	} else {
// 		template.Run(w, r, "secure", "activation_resend_success", "", map[string]interface{}{
// 			"Name":   acc.Name(),
// 			"UID":    acc.ID,
// 			"Remark": remark,
// 		})
// 	}
// }
