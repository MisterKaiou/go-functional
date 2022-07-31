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

func (o *Option[Of]) IsNone() bool {
	return !o.IsSome()
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
