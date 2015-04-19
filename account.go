package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/wscherphof/expeertise/data"
  "github.com/wscherphof/expeertise/model"
  "net/smtp"
  "github.com/wscherphof/expeertise/db"
  "log"
)

func SignUpForm (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  // TODO: captcha
  T("signup", "", map[string]interface{}{
    "Countries": data.Countries(),
  })(w, r, ps)
}

func SignUp (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  if account, err, conflict := model.NewAccount(r.FormValue); err != nil {
    if conflict {
      Error(w, r, ps, err, http.StatusConflict)
    } else {
      Error(w, r, ps, err)
    }
  } else {
    // TODO: confirmation email, formatted response, ...
    const CONFIG_TABLE string = "config"
    const EMAIL_PWD string = "EMAIL_PWD"
    type pwdType struct{
      Key string
      PWD string
    }
    opts := db.NewTableCreateOpts()
    opts.PrimaryKey = "Key"
    if cursor, _ := db.TableCreate(CONFIG_TABLE, opts); cursor != nil {
      log.Println("INFO: table created:", CONFIG_TABLE)
    }
    pwd := new(pwdType)
    if err, found := db.Get(CONFIG_TABLE, EMAIL_PWD, pwd); err != nil {
      Error(w, r, ps, err)
    } else if !(found) {
      db.Insert(CONFIG_TABLE, &pwdType{
        Key: EMAIL_PWD,
        PWD: "",
      })
      log.Println("INFO: record created:", EMAIL_PWD, "in table:", CONFIG_TABLE)
    }
    //
    auth := smtp.PlainAuth("", "expeertise@gmail.com", pwd.PWD, "smtp.gmail.com")
    msg := []byte("Subject: Test email from Go!\nThis is the email body.")
    if err := smtp.SendMail("smtp.gmail.com:587", auth, "expeertise@gmail.org", []string{account.UID}, msg); err != nil {
      Error(w, r, ps, err)
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("account created: " + account.UID))
  }
}
