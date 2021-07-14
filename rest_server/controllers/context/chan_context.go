package context

import (
	"sync"
)

const (
	TokenChannel = "TokenChannel"
)

var context *ChanContext

var once sync.Once

func GetChanInstance() *ChanContext {
	once.Do(func() {
		context = &ChanContext{}
		context.data = make(map[string]interface{})
	})

	return context
}

type ChanContext struct {
	data map[string]interface{}
}

func (o *ChanContext) Put(key string, value interface{}) {
	o.data[key] = value
}

func (o *ChanContext) Get(key string) (interface{}, bool) {
	val, exists := o.data[key]
	return val, exists
}
