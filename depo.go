package depo

import (
	"errors"
	"reflect"
)

type DependencyPool struct {
	avaliableDeps map[string]any
}

func New(
	deps ...any,
) *DependencyPool {
	dep := map[string]any{}
	for i := range deps {
		dep[reflect.TypeOf(deps[i]).Elem().Name()] = deps[i]
	}
	return &DependencyPool{
		avaliableDeps: dep,
	}
}

func (d *DependencyPool) Use(fun any) error {
	if fun == nil {
		return errors.New("fun is nil")
	}

	typ := reflect.TypeOf(fun)
	if typ.Kind() != reflect.Func {
		return errors.New("fun is not a function")
	}

	numParams := typ.NumIn()

	callDeps := make([]reflect.Value, numParams)

	for i := range numParams {
		param := typ.In(i)

		if param.Kind() == reflect.Ptr {
			param = param.Elem()
		}

		if _, ok := d.avaliableDeps[param.Name()]; !ok {
			return errors.New("dependency not found")
		}
		callDeps[i] = reflect.ValueOf(d.avaliableDeps[param.Name()])
	}

	call := reflect.ValueOf(fun)
	call.Call(callDeps)

	return nil
}
