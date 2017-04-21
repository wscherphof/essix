package model

import (
	"github.com/wscherphof/entity"
)

type Profile struct {
	*entity.Base
	FirstName string
	LastName  string
	Country   string
	TimeZone  string
}

func init() {
	entity.Register(InitProfile(""))
}

func InitProfile(account string) (profile *Profile) {
	profile = &Profile{Base: &entity.Base{ID: account}}
	return
}
