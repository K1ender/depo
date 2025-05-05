package depo

import (
	"errors"
	"reflect"
)

var (
	ErrDependencyNotFound = errors.New("dependency not found")
	ErrFunctionIsNil      = errors.New("function is nil")
	ErrIsNotAFunction     = errors.New("is not a function")
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
		return ErrFunctionIsNil
	}

	typ := reflect.TypeOf(fun)
	if typ.Kind() != reflect.Func {
		return ErrIsNotAFunction
	}

	numParams := typ.NumIn()

	callDeps := make([]reflect.Value, numParams)

	for i := range numParams {
		param := typ.In(i)

		if param.Kind() == reflect.Ptr {
			param = param.Elem()
		}

		if _, ok := d.avaliableDeps[param.Name()]; !ok {
			return ErrDependencyNotFound
		}
		callDeps[i] = reflect.ValueOf(d.avaliableDeps[param.Name()])
	}

	call := reflect.ValueOf(fun)
	call.Call(callDeps)

	return nil
}

func (d *DependencyPool) Has(dep any) bool {
	_, ok := d.avaliableDeps[reflect.TypeOf(dep).Elem().Name()]
	return ok
}

func (d *DependencyPool) HasString(name string) bool {
	_, ok := d.avaliableDeps[name]
	return ok
}
