package secure

import (
	"github.com/wscherphof/env"
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
)

func init() {
	var limit = env.GetInt("RATELIMIT", 60)

	router.GET("/account", NewAccountForm)
	router.POST("/account", ratelimit.Handle(NewAccount, limit))

	router.GET("/account/activate", ActivateForm)
	router.PUT("/account/activate", Activate)
	router.GET("/account/activate/token", ActivateTokenForm)
	router.PUT("/account/activate/token", ratelimit.Handle(ActivateToken, limit))

	router.GET("/session", LogInForm)
	router.PUT("/session", ratelimit.Handle(LogIn, limit))
	router.DELETE("/session", LogOut)

	router.GET("/account/password/token", PasswordTokenForm)
	router.PUT("/account/password/token", ratelimit.Handle(PasswordToken, limit))
	router.GET("/account/password", ChangePasswordForm)
	router.PUT("/account/password", ChangePassword)

	// router.GET("/account/emailaddresscode", SecureHandle(EmailAddressCodeForm))
	// router.POST("/account/emailaddresscode", SecureHandle(EmailAddressCode))
	// router.GET("/account/emailaddress/*filepath", SecureHandle(EmailAddressForm))
	// router.PUT("/account/emailaddress", SecureHandle(ChangeEmailAddress))

	// router.GET("/account/terminatecode", SecureHandle(TerminateCodeForm))
	// router.POST("/account/terminatecode", SecureHandle(TerminateCode))
	// router.GET("/account/terminate/*filepath", SecureHandle(TerminateForm))
	// router.DELETE("/account", SecureHandle(Terminate))
}
