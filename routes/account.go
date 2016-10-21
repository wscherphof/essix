package routes

import (
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/routes/account"
	"github.com/wscherphof/essix/template"
	"github.com/wscherphof/secure"
)

func init() {
	// LogIn & LogOut
	router.GET("/session", template.Handle("session", "LogInForm", ""))
	router.PUT("/session", ratelimit.Handle(account.LogIn))
	router.DELETE("/session", account.LogOut)

	// Edit current, or create new account
	router.GET("/account", secure.IfHandle(
		account.EditAccount,
		template.Handle("account", "NewAccountForm", ""),
	))
	router.POST("/account", ratelimit.Handle(account.NewAccount))
	router.GET("/account/post", account.NewAccount)

	// Resend activate token
	router.GET("/account/activate/token", template.Handle("activate", "ActivateTokenForm", ""))
	router.PUT("/account/activate/token", ratelimit.Handle(account.ActivateToken))
	router.GET("/account/activate/token/put", account.ActivateToken)

	// Activate account w/ token
	router.GET("/account/activate", account.ActivateForm)
	router.PUT("/account/activate", account.Activate)
	router.GET("/account/activate/put", account.Activate)

	// Request password change token
	router.GET("/account/password/token", template.Handle("password", "PasswordTokenForm", ""))
	router.PUT("/account/password/token", ratelimit.Handle(account.PasswordToken))
	router.GET("/account/password/token/put", account.PasswordToken)

	// Change password w/ token
	router.GET("/account/password", account.ChangePasswordForm)
	router.PUT("/account/password", account.ChangePassword)
	router.GET("/account/password/put", account.ChangePassword)

	// Request email change token
	router.GET("/account/email/token", secure.Handle(account.EmailTokenForm))
	router.PUT("/account/email/token", secure.Handle(account.EmailToken))
	router.GET("/account/email/token/put", secure.Handle(account.EmailToken))

	// Change email w/ token
	router.GET("/account/email", secure.Handle(account.ChangeEmailForm))
	router.PUT("/account/email", secure.Handle(account.ChangeEmail))
	router.GET("/account/email/put", secure.Handle(account.ChangeEmail))

	// Request account suspension token
	router.GET("/account/suspend/token", secure.Handle(template.Handle("suspend", "SuspendTokenForm", "")))
	router.PUT("/account/suspend/token", secure.Handle(account.SuspendToken))
	router.GET("/account/suspend/token/put", secure.Handle(account.SuspendToken))

	// Suspend account w/ token
	router.GET("/account/suspend", secure.Handle(account.SuspendForm))
	router.DELETE("/account", secure.Handle(account.Suspend))
	router.GET("/account/delete", account.Suspend)
}
