package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type SingletonService struct {
	NotEmptyStruct bool
}

type UserService struct {
	NotEmptyStruct bool
}

type MessageService struct {
	NotEmptyStruct bool
}

type Container struct {
	dependencies map[string]any
}

func NewContainer() *Container {
	return &Container{
		dependencies: make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor any) {
	if _, ok := c.dependencies[name]; ok {
		return
	}

	c.dependencies[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, constructor any) {
	if _, ok := c.dependencies[name]; ok {
		return
	}

	fn, ok := constructor.(func() any)
	if !ok {
		return
	}

	c.dependencies[name] = fn()
}

func (c *Container) Resolve(name string) (any, error) {
	constructor, ok := c.dependencies[name]
	if !ok {
		return nil, errors.New("unknown dependency")
	}

	switch fn := constructor.(type) {
	case (func() any):
		return fn(), nil
	case any:
		return fn, nil
	default:
		return nil, errors.New("bad constructor")
	}
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() any {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() any {
		return &MessageService{}
	})
	container.RegisterSingletonType("SingletonService", func() any {
		return &SingletonService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	singletonService1, err := container.Resolve("SingletonService")
	assert.NoError(t, err)
	singletonService2, err := container.Resolve("SingletonService")
	assert.NoError(t, err)

	s1 := singletonService1.(*SingletonService)
	s2 := singletonService2.(*SingletonService)
	assert.True(t, s1 == s2)
}
