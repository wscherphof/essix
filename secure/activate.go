package secure

import (
	"github.com/julienschmidt/httprouter"
	"github.com/wscherphof/essix/model"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/template"
	"net/http"
)

func ActivateForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, err := ratelimit.NewToken(r, "/account/activate/token"); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "activate", "ActivateForm", "", map[string]interface{}{
			"id":        r.FormValue("id"),
			"token":      r.FormValue("token"),
			"ratelimit": token,
		})
	}
}

func Activate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.GetAccount(r.FormValue("id")); err != nil {
		template.ErrorTail(w, r, err, conflict, "activate", "Activate")
	} else if err, conflict = account.Activate(r.FormValue("token")); err != nil {
		template.ErrorTail(w, r, err, conflict, "activate", "Activate")
	} else {
		template.Run(w, r, "activate", "Activate", "", nil)
	}
}

func ActivateTokenForm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if token, err := ratelimit.NewToken(r); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "activate", "ActivateTokenForm", "", map[string]interface{}{
			"ratelimit": token,
		})
	}
}

func ActivateToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if account, err, conflict := model.GetAccount(r.FormValue("id"), r.FormValue("email")); err != nil {
		template.Error(w, r, err, conflict)
	} else if account.IsActive() {
		template.Error(w, r, model.ErrAlreadyActivated, true)
	} else if err, remark := activateEmail(r, account); err != nil {
		template.Error(w, r, err, false)
	} else {
		template.Run(w, r, "activate", "ActivateToken", "", map[string]interface{}{
			"id":     account.ID,
			"remark": remark,
		})
	}
}
