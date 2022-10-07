package option

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/unit"
)

type some[Of any] any

// Of represents a value that can be either something or nothing.
type Of[T any] struct {
	some some[T]
}

func (o *Of[T]) IsSome() bool {
	return o.some != nil
}

func IsSome[T any](o Of[T]) bool {
	return o.IsSome()
}

func (o *Of[T]) IsNone() bool {
	return !o.IsSome()
}

func IsNone[T any](o Of[T]) bool {
	return o.IsNone()
}

// The value returned when calling this method depends on the state it represents. If Some return fmt.String applied to
// its internal value; else returns "None"
func (o *Of[T]) String() string {
	if o.IsNone() {
		return "None"
	}

	return fmt.Sprint(o.some)
}

// Unwrap Can panic if this Result is an error. Prefer Match over this
func (o *Of[T]) Unwrap() T {
	if o.IsNone() {
		panic("cannot get the value of nothing (none)")
	}

	return o.some.(T)
}

// Map applies the mapping function on the option's internal value if it is Some, and returns a new Of.
func Map[T any, To any](opt Of[T], mapping func(T) To) Of[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return Some(mapping(opt.some.(T)))
}

// Bind accepts a function that takes the Of internal value and returns another Of
func Bind[T any, To any](opt Of[T], binding func(T) Of[To]) Of[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return binding(opt.some.(T))
}

// Match accepts two functions that return a value of the same type, but the first one receives the value stored in
// the given Of and the second one receives no parameters.
func Match[T any, To any](opt Of[T], ok func(T) To, failed func() To) To {
	if opt.IsNone() {
		return failed()
	}

	return ok(opt.some.(T))
}

// Some creates a new Of representing Some state
func Some[T any](it T) Of[T] {
	return Of[T]{
		some: it,
	}
}

// None creates a new Of representing None state
func None[T any]() Of[T] {
	return Of[T]{}
}

// Contains compare the content of the provided Of against the given expected value.
func Contains[T comparable](opt Of[T], expected T) bool {
	if opt.IsNone() {
		return false
	}

	return opt.some.(T) == expected
}

// DefaultValue returns the inner value of this Of or the provided default value.
func DefaultValue[T any](opt Of[T], or T) T {
	if opt.IsNone() {
		return or
	}

	return opt.some.(T)
}

// DefaultWith returns the inner value of this Of or executes the provided function
func DefaultWith[T any](opt Of[T], def func() T) T {
	if opt.IsNone() {
		return def()
	}

	return opt.some.(T)
}

// Exists tests the Of inner value against the given predicate.
func Exists[T any](opt Of[T], predicate func(T) bool) bool {
	if opt.IsNone() {
		return false
	}

	return predicate(opt.some.(T))
}

// Filter returns Some if the Of inner value satisfies the condition, else returns None
func Filter[T any](opt Of[T], predicate func(T) bool) Of[T] {
	if opt.IsNone() {
		return opt
	}

	if !predicate(opt.some.(T)) {
		return None[T]()
	}

	return opt
}

// Fold applies the folder function, passing the provided state and the Of inner value to it and returns State.
func Fold[T, State any](opt Of[T], state State, folder func(State, T) State) State {
	if opt.IsNone() {
		return state
	}

	return folder(state, opt.some.(T))
}

// FoldTo applies the folder function passing the provided state and the Of inner value and returns a new value from it.
func FoldTo[T, State, To any](opt Of[T], state State, folder func(State, T) To) Of[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return Some(folder(state, opt.some.(T)))
}

// FoldM applies the folder function, passing the provided state and the Of inner value to it and returns State
// wrapped in an Of.
func FoldM[T, State any](opt Of[T], state State, folder func(State, T) State) Of[State] {
	if opt.IsNone() {
		return None[State]()
	}

	return Some(folder(state, opt.some.(T)))
}

// CombineBy applies the combiner function on State and the current Result by unwrapping them.
func CombineBy[It, With, To any](opt Of[It], state Of[With], combiner func(With, It) To) Of[To] {
	if opt.IsNone() {
		return None[To]()
	}

	if state.IsNone() {
		return None[To]()
	}

	return Some(combiner(state.some.(With), opt.some.(It)))
}

// Flatten returns a Result from a Result of Result.
func Flatten[T any](opt Of[Of[T]]) Of[T] {
	if opt.IsNone() {
		return None[T]()
	}

	return opt.some.(Of[T])
}

// Iter applies the given action to the inner value of the Of provided.
func Iter[T any](opt Of[T], action func(it T) unit.Unit) unit.Unit {
	if opt.IsNone() {
		return unit.Unit{}
	}

	return action(opt.some.(T))
}
