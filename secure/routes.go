package secure

import (
	"github.com/wscherphof/essix/ratelimit"
	"github.com/wscherphof/essix/router"
)

func init() {
	router.Router.GET("/account", IfSecureHandle(UpdateAccountForm, SignUpForm))
	router.Router.POST("/account", ratelimit.Handle(3600, SignUp))
	router.Router.PUT("/account", SecureHandle(UpdateAccount))

	router.Router.GET("/session", LogInForm)
	router.Router.POST("/session", ratelimit.Handle(60, LogIn))
	router.Router.DELETE("/session", LogOut)

	router.Router.GET("/account/activation/:uid", ActivateForm)
	router.Router.GET("/account/activation", ActivateForm)
	router.Router.PUT("/account/activation", Activate)
	router.Router.GET("/account/activationcode/:uid", ActivationCodeForm)
	router.Router.GET("/account/activationcode", ActivationCodeForm)
	router.Router.POST("/account/activationcode", ratelimit.Handle(3600, ActivationCode))

	router.Router.GET("/account/passwordcode/:uid", PasswordCodeForm)
	router.Router.GET("/account/passwordcode", PasswordCodeForm)
	router.Router.POST("/account/passwordcode", ratelimit.Handle(3600, PasswordCode))
	router.Router.GET("/account/password/:uid", PasswordForm)
	router.Router.PUT("/account/password", ChangePassword)

	router.Router.GET("/account/emailaddresscode", SecureHandle(EmailAddressCodeForm))
	router.Router.POST("/account/emailaddresscode", SecureHandle(EmailAddressCode))
	router.Router.GET("/account/emailaddress/*filepath", SecureHandle(EmailAddressForm))
	router.Router.PUT("/account/emailaddress", SecureHandle(ChangeEmailAddress))

	router.Router.GET("/account/terminatecode", SecureHandle(TerminateCodeForm))
	router.Router.POST("/account/terminatecode", SecureHandle(TerminateCode))
	router.Router.GET("/account/terminate/*filepath", SecureHandle(TerminateForm))
	router.Router.DELETE("/account", SecureHandle(Terminate))
}
