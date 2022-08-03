package result

import (
	"fmt"

	"github.com/MisterKaiou/go-functional/option"
	"github.com/MisterKaiou/go-functional/unit"
)

type ok[Of any] any

// Result represents a result that can be either an input of given type, or error.
type Result[Ok any] struct {
	ok  ok[Ok]
	err error
}

// The value returned when calling this method depends on the state it represents. If ok return fmt.String applied to
// its internal value; if error, return error.Error()
func (r *Result[Ok]) String() string {
	if r.IsError() {
		return r.err.Error()
	}

	return fmt.Sprint(r.ok)
}

func (r *Result[Ok]) IsOk() bool {
	return r.err == nil
}

func IsOk[Of any](res Result[Of]) bool {
	return res.IsOk()
}

func (r *Result[Ok]) IsError() bool {
	return !r.IsOk()
}

func IsError[Of any](res Result[Of]) bool {
	return res.IsError()
}

// Unwrap can panic if this Result is an error. Prefer Match over this
func (r *Result[Ok]) Unwrap() Ok {
	if r.IsError() {
		panic("cannot get the value of an error")
	}

	return r.ok.(Ok)
}

// Map applies the mapping function on the result's internal value if is not an error, and returns a new result.
func Map[From any, To any](res Result[From], mapping func(From) To) Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return Ok[To](mapping(res.ok.(From)))
}

// MapError applies the given mapping function on the result's internal error, if it is an error, and returns a new
// result, else returns the same instance provided.
func MapError[Of any](res Result[Of], mapping func(error) error) Result[Of] {
	if res.IsError() {
		return Error[Of](mapping(res.err))
	}

	return res
}

// Bind accepts a function that takes the Result internal value and returns another Result
func Bind[From any, To any](res Result[From], binding func(From) Result[To]) Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return binding(res.ok.(From))
}

// Match accepts two functions that return a value of the same type, but the first one receives the result
// contained in it and the second one receives the error.
func Match[Ok any, To any](res Result[Ok], ok func(Ok) To, failed func(error) To) To {
	if res.IsError() {
		return failed(res.err)
	}

	return ok(res.ok.(Ok))
}

// Ok creates a new Result representing an Ok state.
func Ok[Ok any](it Ok) Result[Ok] {
	return Result[Ok]{
		ok:  it,
		err: nil,
	}
}

// Error create a new Result that represents an Error state.
func Error[Ok any](err error) Result[Ok] {
	return Result[Ok]{
		ok:  nil,
		err: err,
	}
}

// FromTupleOf creates a new Result base on the values provided. If err is not nil, this result will represent an error.
func FromTupleOf[Of any](it Of, err error) Result[Of] {
	if err != nil {
		return Error[Of](err)
	}

	return Ok[Of](it)
}

// Contains compare the content of the provided Result against the given expected value.
func Contains[Of comparable](res Result[Of], expected Of) bool {
	if res.IsError() {
		return false
	}

	return res.ok == expected
}

// DefaultValue returns the inner value of this Result or the provided default value.
func DefaultValue[Of any](res Result[Of], or Of) Of {
	if res.IsError() {
		return or
	}

	return res.ok.(Of)
}

// DefaultWith returns the inner value of this Result or executes the provided function with its inner error.
func DefaultWith[Of any](res Result[Of], def func(error) Of) Of {
	if res.IsError() {
		return def(res.err)
	}

	return res.ok.(Of)
}

// Exists tests the Result inner value against the given predicate.
func Exists[Of any](res Result[Of], predicate func(Of) bool) bool {
	if res.IsError() {
		return false
	}

	return predicate(res.ok.(Of))
}

// Fold applies the folder function, passing the provided state and the Result inner value to it and returns State.
func Fold[Of, State any](res Result[Of], state State, folder func(State, Of) State) State {
	if res.IsError() {
		return state
	}

	return folder(state, res.ok.(Of))
}

// FoldM applies the folder function, passing the provided state and the Result inner value to it and returns State
// wrapped in a Result.
func FoldM[Of, State any](res Result[Of], state State, folder func(State, Of) State) Result[State] {
	if res.IsError() {
		return Error[State](res.err)
	}

	return Ok(folder(state, res.ok.(Of)))
}

// CombineBy applies the combiner function on State and the current Result by unwrapping them.
func CombineBy[It, With, To any](res Result[It], state Result[With], combiner func(With, It) To) Result[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	if state.IsError() {
		return Error[To](state.err)
	}

	return Ok(combiner(state.ok.(With), res.ok.(It)))
}

// Iter applies the given action to the inner value of the Result provided.
func Iter[Of any](res Result[Of], action func(it Of) unit.Unit) unit.Unit {
	if res.IsError() {
		return unit.Unit{}
	}

	return action(res.ok.(Of))
}

// Flatten returns a Result from a Result of Result.
func Flatten[Of any](res Result[Result[Of]]) Result[Of] {
	if res.IsError() {
		return Error[Of](res.err)
	}

	return res.ok.(Result[Of])
}

// ToOption creates an Option from the given Result. If error, the returned Option will be None, else Some with the inner
// value.
func ToOption[Of any](res Result[Of]) option.Option[Of] {
	if res.IsError() {
		return option.None[Of]()
	}

	return option.Some(res.ok.(Of))
}
