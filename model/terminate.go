package model

func (a *Account) CreateTerminateCode(sure bool) (err error, conflict bool) {
	if !sure {
		err, conflict = ErrCodeIncorrect, true
	} else {
		a.TerminateCode = NewCode()
		err = a.Update(a)
	}
	return
}

func (a *Account) ClearTerminateCode(code string) (err error) {
	if a.TerminateCode == code {
		a.TerminateCode = ""
		err = a.Update(a)
	}
	return
}

func (a *Account) Terminate(code string, sure bool) (err error, conflict bool) {
	uid := a.ID
	if !sure {
		err, conflict = ErrCodeIncorrect, true
	} else if acc, e, c := getAccount(uid); e != nil {
		err, conflict = e, c
	} else if acc.TerminateCode == "" {
		err, conflict = ErrCodeUnset, true
	} else if code == "" || code != acc.TerminateCode {
		err, conflict = ErrCodeIncorrect, true
	} else {
		err = acc.Delete(acc)
	}
	return
}
