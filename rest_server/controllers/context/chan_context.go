package context

import (
	"sync"
)

const (
	TokenChannel = "TokenChannel"
)

var context *Context

var once sync.Once

func GetInstance() *Context {
	once.Do(func() {
		context = &Context{}
		context.data = make(map[string]interface{})
	})

	return context
}

type Context struct {
	data map[string]interface{}
}

func (o *Context) Put(key string, value interface{}) {
	o.data[key] = value
}

func (o *Context) Get(key string) (interface{}, bool) {
	val, exists := o.data[key]
	return val, exists
}
