package secure

import (
  "github.com/wscherphof/expeertise/router"
)

func init() {
  router.GET    ("/account", IfSecureHandle(UpdateAccountForm, SignUpForm))
  router.POST   ("/account", SignUp)
  router.PUT    ("/account", SecureHandle(UpdateAccount))
  // TODO: router.DELETE ("/account", Authenticate(TerminateAccount))

  router.GET    ("/session", LogInForm)
  router.POST   ("/session", LogIn)
  router.DELETE ("/session", LogOut)

  router.GET    ("/account/activation",          ActivateForm)
  router.GET    ("/account/activation/",         ActivateForm)
  router.GET    ("/account/activation/:uid",     ActivateForm)
  router.PUT    ("/account/activation",          Activate)
  router.GET    ("/account/activationcode",      ActivationCodeForm)
  router.GET    ("/account/activationcode/",     ActivationCodeForm)
  router.GET    ("/account/activationcode/:uid", ActivationCodeForm)
  router.POST   ("/account/activationcode",      ActivationCode)
  
  router.GET    ("/account/passwordcode",        PasswordCodeForm)
  router.GET    ("/account/passwordcode/",       PasswordCodeForm)
  router.GET    ("/account/passwordcode/:uid",   PasswordCodeForm)
  router.POST   ("/account/passwordcode",        PasswordCode)
  router.GET    ("/account/password/:uid",       PasswordForm)
  router.PUT    ("/account/password",            ChangePassword)
  
  router.GET    ("/account/emailaddresscode",       SecureHandle(EmailAddressCodeForm))
  router.GET    ("/account/emailaddresscode/",      SecureHandle(EmailAddressCodeForm))
  router.POST   ("/account/emailaddresscode",       SecureHandle(EmailAddressCode))
  router.GET    ("/account/emailaddress/*filepath", SecureHandle(EmailAddressForm))
  router.PUT    ("/account/emailaddress",           SecureHandle(ChangeEmailAddress))
}
