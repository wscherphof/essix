package secure

import (
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
	"github.com/wscherphof/env"
)

func init() {
	var limit = env.GetInt("RATELIMIT", 60)

	router.GET("/account", SignUpForm)
	router.POST("/account", ratelimit.Handle(SignUp, limit))

	router.GET("/session", LogInForm)
	router.POST("/session", ratelimit.Handle(LogIn, limit))
	router.DELETE("/session", LogOut)

	router.GET("/account/activation/:uid", ActivateForm)
	router.GET("/account/activation", ActivateForm)
	router.PUT("/account/activation", Activate)
	router.GET("/account/activationcode/:uid", ActivationCodeForm)
	router.GET("/account/activationcode", ActivationCodeForm)
	router.POST("/account/activationcode", ratelimit.Handle(ActivationCode, limit))

	router.GET("/account/passwordcode/:uid", PasswordCodeForm)
	router.GET("/account/passwordcode", PasswordCodeForm)
	router.POST("/account/passwordcode", ratelimit.Handle(PasswordCode, limit))
	router.GET("/account/password/:uid", PasswordForm)
	router.PUT("/account/password", ChangePassword)

	router.GET("/account/emailaddresscode", SecureHandle(EmailAddressCodeForm))
	router.POST("/account/emailaddresscode", SecureHandle(EmailAddressCode))
	router.GET("/account/emailaddress/*filepath", SecureHandle(EmailAddressForm))
	router.PUT("/account/emailaddress", SecureHandle(ChangeEmailAddress))

	router.GET("/account/terminatecode", SecureHandle(TerminateCodeForm))
	router.POST("/account/terminatecode", SecureHandle(TerminateCode))
	router.GET("/account/terminate/*filepath", SecureHandle(TerminateForm))
	router.DELETE("/account", SecureHandle(Terminate))
}
