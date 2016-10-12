package secure

import (
	"github.com/wscherphof/env"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
)

func init() {
	var limit = env.GetInt("RATELIMIT", 60)

	router.GET("/session", LogInForm)
	router.PUT("/session", ratelimit.Handle(LogIn, limit))
	router.DELETE("/session", LogOut)

	router.GET("/account", IfHandle(Account, NewAccountForm))
	router.POST("/account", ratelimit.Handle(NewAccount, limit))
	router.GET("/account/post", NewAccount)

	router.GET("/account/activate/token", ActivateTokenForm)
	router.PUT("/account/activate/token", ratelimit.Handle(ActivateToken, limit))
	router.GET("/account/activate/token/put", ActivateToken)

	router.GET("/account/activate", ActivateForm)
	router.PUT("/account/activate", Activate)
	router.GET("/account/activate/put", Activate)

	router.GET("/account/password/token", PasswordTokenForm)
	router.PUT("/account/password/token", ratelimit.Handle(PasswordToken, limit))
	router.GET("/account/password/token/put", PasswordToken)

	router.GET("/account/password", ChangePasswordForm)
	router.PUT("/account/password", ChangePassword)
	router.GET("/account/password/put", ChangePassword)

	router.GET("/account/email/token", Handle(EmailTokenForm))
	router.PUT("/account/email/token", Handle(EmailToken))
	router.GET("/account/email/token/put", Handle(EmailToken))

	router.GET("/account/email", Handle(ChangeEmailForm))
	router.PUT("/account/email", Handle(ChangeEmail))
	router.GET("/account/email/put", Handle(ChangeEmail))

	router.GET("/account/suspend/token", Handle(SuspendTokenForm))
	router.PUT("/account/suspend/token", Handle(SuspendToken))
	router.GET("/account/suspend/token/put", Handle(SuspendToken))

	router.GET("/account/suspend", Handle(SuspendForm))
	router.DELETE("/account", Handle(Suspend))
	router.GET("/account/delete", Handle(Suspend))
}
