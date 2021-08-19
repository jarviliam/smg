package context

import (
	"errors"
	"sync"
)

type Context struct {
	contextMap map[string]string
	lock       *sync.Mutex
}

func NewContext() *Context {
	c := &Context{
		contextMap: make(map[string]string),
		lock:       &sync.Mutex{},
	}
	return c
}

func (c *Context) Get(key string) string {
	if v, ok := c.contextMap[key]; ok {
		return v
	}
	return ""
}

func (c *Context) Put(key, value string, ovrride bool) error {
	if _, ok := c.contextMap[key]; ok && ovrride {
		c.contextMap[key] = value
	} else if ok && !ovrride {
		return errors.New("Key already exists")
	} else {
		c.contextMap[key] = value
	}
	return nil
}
