package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (interface{}, error)
	GetOrSet(key string, get interface{}, val func() (interface{}, error), opts ...Option) error
	Set(key string, val interface{}, opts ...Option) error
	Del(key string) error
}

type mapCache struct {
	sync.RWMutex
	cache         map[string]*expireValue
	defaultExpire int64
}

func NewMapCache(defaultExpire int64) Cache {
	return &mapCache{
		cache:         make(map[string]*expireValue),
		defaultExpire: defaultExpire,
	}
}

func (m *mapCache) Get(key string) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.cache[key]
	if ok {
		if val.expire > 0 && time.Now().Unix()-val.expire > val.createtime {
			return nil, ErrNotExist
		}
		return val.value, nil
	}
	return nil, ErrNotExist
}

func (m *mapCache) GetOrSet(key string, get interface{}, f func() (interface{}, error), opts ...Option) error {
	val, err := m.Get(key)
	if err != nil && err != ErrNotExist {
		return err
	}
	if err == ErrNotExist {
		val, err = f()
		if err != nil {
			return err
		}
		if err := m.Set(key, val, opts...); err != nil {
			return err
		}
	}
	copyInterface(get, val)
	return nil
}

func (m *mapCache) Set(key string, val interface{}, opts ...Option) error {
	m.Lock()
	defer m.Unlock()
	opt := MakeOption(opts...)
	if opt.expire == 0 {
		opt.expire = m.defaultExpire
	}
	m.cache[key] = &expireValue{
		value:      val,
		expire:     opt.expire,
		createtime: time.Now().Unix(),
	}
	return nil
}

func (m *mapCache) Del(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.cache, key)
	return nil
}
