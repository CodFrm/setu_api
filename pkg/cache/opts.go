package cache

import (
	"context"
)

type expireValue struct {
	value      interface{}
	expire     int64
	createtime int64
}

type options struct {
	ctx    context.Context
	expire int64
}

type Option func(opts *options)

func MakeOption(opts ...Option) *options {
	ret := &options{}
	for _, v := range opts {
		v(ret)
	}
	return ret
}

func Expire(expire int64) Option {
	return func(opts *options) {
		opts.expire = expire
	}
}

func WithContext(ctx context.Context) Option {
	return func(opts *options) {
		opts.ctx = ctx
	}
}
