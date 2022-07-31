package result

import "fmt"

type Result[Ok any] struct {
	ok  Ok
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

func (r *Result[Ok]) IsError() bool {
	return !r.IsOk()
}

// Unwrap Can panic if this Result is an error. Prefer Match over this
func (r *Result[Ok]) Unwrap() Ok {
	if r.IsError() {
		panic("cannot get the value of an error")
	}

	return r.ok
}

func Map[From any, To any](res *Result[From], mapping func(From) To) *Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return Ok[To](mapping(res.ok))
}

func Bind[From any, To any](res *Result[From], binding func(From) *Result[To]) *Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return binding(res.ok)
}

func Match[Ok any, To any](res *Result[Ok], ok func(Ok) To, failed func(error) To) To {
	if res.IsError() {
		return failed(res.err)
	}

	return ok(res.ok)
}

func Ok[Ok any](it Ok) *Result[Ok] {
	return &Result[Ok]{
		ok:  it,
		err: nil,
	}
}

func Error[Ok any](err error) *Result[Ok] {
	return &Result[Ok]{
		ok:  *new(Ok),
		err: err,
	}
}

func FromTupleOf[Of any](it Of, err error) *Result[Of] {
	if err != nil {
		return Error[Of](err)
	}

	return Ok[Of](it)
}
