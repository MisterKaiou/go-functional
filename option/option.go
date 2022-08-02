package option

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/unit"
)

type some[Of any] any

// Option represents a value that can be either something or nothing.
type Option[Of any] struct {
	some some[Of]
}

func (o *Option[Of]) IsSome() bool {
	return o.some != nil
}

func IsSome[Of any](o Option[Of]) bool {
	return o.IsSome()
}

func (o *Option[Of]) IsNone() bool {
	return !o.IsSome()
}

func IsNone[Of any](o Option[Of]) bool {
	return o.IsNone()
}

// The value returned when calling this method depends on the state it represents. If Some return fmt.String applied to
// its internal value; else returns "None"
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

// Map applies the mapping function on the option's internal value if it is Some, and returns a new Option.
func Map[Of any, To any](opt Option[Of], mapping func(Of) To) Option[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return Some(mapping(opt.some.(Of)))
}

// Bind accepts a function that takes the Option internal value and returns another Option
func Bind[Of any, To any](opt Option[Of], binding func(Of) Option[To]) Option[To] {
	if opt.IsNone() {
		return None[To]()
	}

	return binding(opt.some.(Of))
}

// Match accepts two functions that return a value of the same type, but the first one receives the value stored in
// the given Option and the second one receives no parameters.
func Match[Of any, To any](opt Option[Of], ok func(Of) To, failed func() To) To {
	if opt.IsNone() {
		return failed()
	}

	return ok(opt.some.(Of))
}

// Some creates a new Option representing Some state
func Some[Of any](it Of) Option[Of] {
	return Option[Of]{
		some: it,
	}
}

// None creates a new Option representing None state
func None[Of any]() Option[Of] {
	return Option[Of]{}
}

// Contains compare the content of the provided Option against the given expected value.
func Contains[Of comparable](opt Option[Of], expected Of) bool {
	if opt.IsNone() {
		return false
	}

	return opt.some.(Of) == expected
}

// DefaultValue returns the inner value of this Option or the provided default value.
func DefaultValue[Of any](opt Option[Of], or Of) Of {
	if opt.IsNone() {
		return or
	}

	return opt.some.(Of)
}

// DefaultWith returns the inner value of this Option or executes the provided function
func DefaultWith[Of any](opt Option[Of], def func() Of) Of {
	if opt.IsNone() {
		return def()
	}

	return opt.some.(Of)
}

// Exists tests the Option inner value against the given predicate.
func Exists[Of any](opt Option[Of], predicate func(Of) bool) bool {
	if opt.IsNone() {
		return false
	}

	return predicate(opt.some.(Of))
}

// Filter returns Some if the Option inner value satisfies the condition, else returns None
func Filter[Of any](opt Option[Of], predicate func(Of) bool) Option[Of] {
	if opt.IsNone() {
		return opt
	}

	if !predicate(opt.some.(Of)) {
		return None[Of]()
	}

	return opt
}

// Fold applies the folder function, passing the provided state and the Option inner value to it and returns State.
func Fold[Of, State any](opt Option[Of], state State, folder func(State, Of) State) State {
	if opt.IsNone() {
		return state
	}

	return folder(state, opt.some.(Of))
}

// FoldM applies the folder function, passing the provided state and the Option inner value to it and returns State
// wrapped in an Option.
func FoldM[Of, State any](opt Option[Of], state State, folder func(State, Of) State) Option[State] {
	if opt.IsNone() {
		return None[State]()
	}

	return Some(folder(state, opt.some.(Of)))
}

// Iter applies the given action to the inner value of the Option provided.
func Iter[Of any](opt Option[Of], action func(it Of) unit.Unit) unit.Unit {
	if opt.IsNone() {
		return unit.Unit{}
	}

	return action(opt.some.(Of))
}
