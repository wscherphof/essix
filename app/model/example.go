package model

import (
	"github.com/wscherphof/entity"
)

type Profile struct {
	*entity.Base
	FirstName string
	LastName  string
	Country   string
	Postcode  string
}

func InitProfile(account string) (profile *Profile) {
	profile = &Profile{Base: &entity.Base{ID: account}}
	return
}

func (p *Profile) FullName() (name string) {
	name = ""
	if len(p.FirstName) > 0 {
		name = p.FirstName
	}
	if len(p.LastName) > 0 {
		if len(name) > 0 {
			name = name + " "
		}
		name = name + p.LastName
	}
	if len(name) == 0 {
		name = p.Account
	}
	return
}

func (p *Profile) Name() (name string) {
	if len(p.FirstName) > 0 {
		name = p.FirstName
	} else if len(p.LastName) > 0 {
		name = p.LastName
	} else {
		name = p.Account
	}
	return
}
