package cache

import (
	"errors"
	"reflect"
	"strings"
)

var ErrNotExist = errors.New("not exist")

func Join(keys ...string) string {
	return strings.Join(keys, ":")
}

func copyInterface(dst interface{}, src interface{}) {
	dstof := reflect.ValueOf(dst)
	if dstof.Kind() == reflect.Ptr {
		el := dstof.Elem()
		srcof := reflect.ValueOf(src)
		if srcof.Kind() == reflect.Ptr {
			el.Set(srcof.Elem())
		} else {
			el.Set(srcof)
		}
	}
}
