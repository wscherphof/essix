package config

import (
	"github.com/wscherphof/essix/entity"
	"strings"
)

type item struct {
	*entity.Base
}

func init() {
	entity.Register("", "config")
}

func initItem(key string) *item {
	return &item{Base: &entity.Base{ID: strings.ToLower(key)}}
}

func Get(key string, result interface{}) {
	i := initItem(key)

}
