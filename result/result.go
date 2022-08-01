package result

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/option"
	"github.com/MisterKaiou/go-functional/unit"
)

type ok[Of any] any

type Result[Ok any] struct {
	ok  ok[Ok]
	err error
}

func (r *Result[Ok]) String() string {
	if r.IsError() {
		return r.err.Error()
	}

	return fmt.Sprint(r.ok)
}

func (r *Result[Ok]) IsOk() bool {
	return r.err == nil
}

func IsOk[Of any](res *Result[Of]) bool {
	return res.IsOk()
}

func (r *Result[Ok]) IsError() bool {
	return !r.IsOk()
}

func IsError[Of any](res *Result[Of]) bool {
	return res.IsError()
}

// Unwrap Can panic if this Result is an error. Prefer Match over this
func (r *Result[Ok]) Unwrap() Ok {
	if r.IsError() {
		panic("cannot get the value of an error")
	}

	return r.ok.(Ok)
}

func Map[From any, To any](res *Result[From], mapping func(From) To) *Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return Ok[To](mapping(res.ok.(From)))
}

func MapError[Of any](res *Result[Of], mapping func(error) error) *Result[Of] {
	if res.IsError() {
		return Error[Of](mapping(res.err))
	}

	return res
}

func Bind[From any, To any](res *Result[From], binding func(From) *Result[To]) *Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return binding(res.ok.(From))
}

func Match[Ok any, To any](res *Result[Ok], ok func(Ok) To, failed func(error) To) To {
	if res.IsError() {
		return failed(res.err)
	}

	return ok(res.ok.(Ok))
}

func Ok[Ok any](it Ok) *Result[Ok] {
	return &Result[Ok]{
		ok:  it,
		err: nil,
	}
}

func Error[Ok any](err error) *Result[Ok] {
	return &Result[Ok]{
		ok:  nil,
		err: err,
	}
}

func FromTupleOf[Of any](it Of, err error) *Result[Of] {
	if err != nil {
		return Error[Of](err)
	}

	return Ok[Of](it)
}

func Contains[Of comparable](res *Result[Of], expected Of) bool {
	if res.IsError() {
		return false
	}

	return res.ok == expected
}

func DefaultValue[Of any](res *Result[Of], or Of) Of {
	if res.IsError() {
		return or
	}

	return res.ok.(Of)
}

func DefaultWith[Of any](res *Result[Of], def func(error) Of) Of {
	if res.IsError() {
		return def(res.err)
	}

	return res.ok.(Of)
}

func Exists[Of any](res *Result[Of], predicate func(Of) bool) bool {
	if res.IsError() {
		return false
	}

	return predicate(res.ok.(Of))
}

func Fold[Of, State any](res *Result[Of], state State, folder func(State, Of) State) State {
	if res.IsError() {
		return state
	}

	return folder(state, res.ok.(Of))
}

func Iter[Of any](res *Result[Of], action func(it Of) unit.Unit) unit.Unit {
	if res.IsError() {
		return unit.Unit{}
	}

	return action(res.ok.(Of))
}

func ToOption[Of any](res *Result[Of]) *option.Option[Of] {
	if res.IsError() {
		return option.None[Of]()
	}

	return option.Some(res.ok.(Of))
}