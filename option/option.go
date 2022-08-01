package option

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/unit"
)

type none unit.Unit
type some[Of any] any

type Option[Of any] struct {
	none *none
	some some[Of]
}

func (o *Option[Of]) IsSome() bool {
	return o.none == nil
}

func IsSome[Of any](o *Option[Of]) bool {
	return o.IsSome()
}

func (o *Option[Of]) IsNone() bool {
	return !o.IsSome()
}

func IsNone[Of any](o *Option[Of]) bool {
	return o.IsNone()
}

func (o *Option[Of]) String() string {
	if o.IsNone() {
		return "None"
	}

	return fmt.Sprint(o.some)
}

// Unwrap Can panic if this Result is an error. Prefer Match over this
func (o *Option[Of]) Unwrap() Of {
	if o.IsNone() {
		panic("cannot get the value of nothing (none)")
	}

	return o.some.(Of)
}

func Map[Of any, To any](opt *Option[Of], mapping func(Of) To) *Option[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return Some(mapping(opt.some.(Of)))
}

func Bind[Of any, To any](opt *Option[Of], binding func(Of) *Option[To]) *Option[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return binding(opt.some.(Of))
}

func Match[Of any, To any](opt *Option[Of], ok func(Of) To, failed func() To) To {
	if opt.IsNone() {
		return failed()
	}

	return ok(opt.some.(Of))
}

func Some[Of any](it Of) *Option[Of] {
	return &Option[Of]{
		some: it,
	}
}

func None[Of any]() *Option[Of] {
	return &Option[Of]{none: &none{}}
}

func Contains[Of comparable](opt *Option[Of], expected Of) bool {
	if opt.IsNone() {
		return false
	}

	return opt.some.(Of) == expected
}

func DefaultValue[Of any](opt *Option[Of], or Of) Of {
	if opt.IsNone() {
		return or
	}

	return opt.some.(Of)
}

func DefaultWith[Of any](opt *Option[Of], def func() Of) Of {
	if opt.IsNone() {
		return def()
	}

	return opt.some.(Of)
}

func Exists[Of any](opt *Option[Of], predicate func(Of) bool) bool {
	if opt.IsNone() {
		return false
	}

	return predicate(opt.some.(Of))
}

func Filter[Of any](opt *Option[Of], predicate func(Of) bool) *Option[Of] {
	if opt.IsNone() || !predicate(opt.some.(Of)) {
		return None[Of]()
	}

	return opt
}

func Fold[Of, State any](opt *Option[Of], state State, folder func(State, Of) State) State {
	if opt.IsNone() {
		return state
	}

	return folder(state, opt.some.(Of))
}

func Iter[Of any](opt *Option[Of], action func(it Of) unit.Unit) unit.Unit {
	if opt.IsNone() {
		return unit.Unit{}
	}

	return action(opt.some.(Of))
}
