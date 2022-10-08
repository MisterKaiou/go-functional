package result

import (
	"fmt"

	"github.com/MisterKaiou/go-functional/option"
	"github.com/MisterKaiou/go-functional/unit"
)

type ok[Of any] any

// Of represents a result that can be either an input of given type, or error.
type Of[Ok any] struct {
	ok  ok[Ok]
	err error
}

// The value returned when calling this method depends on the state it represents. If ok return fmt.String applied to
// its internal value; if error, return error.Error()
func (r *Of[Ok]) String() string {
	if r.IsError() {
		return r.err.Error()
	}

	return fmt.Sprint(r.ok)
}

func (r *Of[Ok]) IsOk() bool {
	return r.err == nil
}

func IsOk[T any](res Of[T]) bool {
	return res.IsOk()
}

func (r *Of[Ok]) IsError() bool {
	return !r.IsOk()
}

func IsError[T any](res Of[T]) bool {
	return res.IsError()
}

// Unwrap can panic if this Of is an error. Prefer Match over this
func (r *Of[Ok]) Unwrap() Ok {
	if r.IsError() {
		panic("cannot get the value of an error")
	}

	return r.ok.(Ok)
}

// UnwrapError can panic if this Of is not an error.
func (r *Of[Ok]) UnwrapError() error {
	if r.IsOk() {
		panic("cannot get the error of an sucessful result")
	}

	return r.err
}

// Map applies the mapping function on the result's internal value if is not an error, and returns a new result.
func Map[From any, To any](res Of[From], mapping func(From) To) Of[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return Ok[To](mapping(res.ok.(From)))
}

// MapError applies the given mapping function on the result's internal error, if it is an error, and returns a new
// result, else returns the same instance provided.
func MapError[T any](res Of[T], mapping func(error) error) Of[T] {
	if res.IsError() {
		return Error[T](mapping(res.err))
	}

	return res
}

// Bind accepts a function that takes the Of internal value and returns another Of
func Bind[From any, To any](res Of[From], binding func(From) Of[To]) Of[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return binding(res.ok.(From))
}

// Match accepts two functions that return a value of the same type, but the first one receives the result
// contained in it and the second one receives the error.
func Match[Ok any, To any](res Of[Ok], ok func(Ok) To, failed func(error) To) To {
	if res.IsError() {
		return failed(res.err)
	}

	return ok(res.ok.(Ok))
}

// Ok creates a new Of representing an Ok state.
func Ok[Ok any](it Ok) Of[Ok] {
	return Of[Ok]{
		ok:  it,
		err: nil,
	}
}

// Error create a new Of that represents an Error state.
func Error[Ok any](err error) Of[Ok] {
	return Of[Ok]{
		ok:  nil,
		err: err,
	}
}

// FromTupleOf creates a new Of base on the values provided. If err is not nil, this result will represent an error.
func FromTupleOf[T any](it T, err error) Of[T] {
	if err != nil {
		return Error[T](err)
	}

	return Ok[T](it)
}

// Contains compare the content of the provided Of against the given expected value.
func Contains[T comparable](res Of[T], expected T) bool {
	if res.IsError() {
		return false
	}

	return res.ok == expected
}

// DefaultValue returns the inner value of this Of or the provided default value.
func DefaultValue[T any](res Of[T], or T) T {
	if res.IsError() {
		return or
	}

	return res.ok.(T)
}

// DefaultWith returns the inner value of this Of or executes the provided function with its inner error.
func DefaultWith[T any](res Of[T], def func(error) T) T {
	if res.IsError() {
		return def(res.err)
	}

	return res.ok.(T)
}

// Exists tests the Of inner value against the given predicate.
func Exists[T any](res Of[T], predicate func(T) bool) bool {
	if res.IsError() {
		return false
	}

	return predicate(res.ok.(T))
}

// Fold applies the folder function passing the provided state and the Of inner value to it and returns the updated State.
func Fold[T, State any](res Of[T], state State, folder func(State, T) State) State {
	if res.IsError() {
		return state
	}

	return folder(state, res.ok.(T))
}

// FoldTo applies the folder function passing the provided state and the Of inner value and returns a new value from it.
func FoldTo[T, State, To any](res Of[T], state State, folder func(State, T) To) Of[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	return Ok(folder(state, res.ok.(T)))
}

// FoldM applies the folder function, passing the provided state and the Of inner value to it and returns State
// wrapped in a Of.
func FoldM[T, State any](res Of[T], state State, folder func(State, T) State) Of[State] {
	if res.IsError() {
		return Error[State](res.err)
	}

	return Ok(folder(state, res.ok.(T)))
}

// CombineBy applies the combiner function on State and the current Of by unwrapping them.
func CombineBy[It, With, To any](res Of[It], state Of[With], combiner func(With, It) To) Of[To] {
	if res.IsError() {
		return Error[To](res.err)
	}

	if state.IsError() {
		return Error[To](state.err)
	}

	return Ok(combiner(state.ok.(With), res.ok.(It)))
}

// Iter applies the given action to the inner value of the Of provided.
func Iter[T any](res Of[T], action func(it T) unit.Unit) unit.Unit {
	if res.IsError() {
		return unit.Unit{}
	}

	return action(res.ok.(T))
}

// Flatten returns a Of from a Of of Of.
func Flatten[T any](res Of[Of[T]]) Of[T] {
	if res.IsError() {
		return Error[T](res.err)
	}

	return res.ok.(Of[T])
}

// ToOption creates an Of from the given Of. If error, the returned Of will be None, else Some with the inner
// value.
func ToOption[T any](res Of[T]) option.Of[T] {
	if res.IsError() {
		return option.None[T]()
	}

	return option.Some(res.ok.(T))
}
