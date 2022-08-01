package result

import (
	"errors"
	"fmt"
	"github.com/MisterKaiou/go-functional/option"
	"github.com/MisterKaiou/go-functional/unit"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	s := "result"
	valRes := Ok(s)

	// Testing Pass by Value
	assert.Equal(t, valRes.ok, s)
	assert.NotSame(t, &valRes.ok, &s)
	assert.Nil(t, valRes.err)

	refRes := Ok(&s)

	// Testing Pass by Reference
	assert.Equal(t, s, *(refRes.ok.(*string)))
	assert.Equal(t, &s, refRes.ok)
	assert.Nil(t, refRes.err)
}

func TestErr(t *testing.T) {
	err := errors.New("some error")
	res := Error[int](err)

	assert.Equal(t, err, res.err)
	assert.Equal(t, &err, &res.err)
}

func TestMapNoError(t *testing.T) {
	value := 42
	res := Ok(value)

	mappedRes := Map(res, func(val int) string { return fmt.Sprint(val) })

	assert.Equal(t, "42", mappedRes.ok)
}

func TestMapWithError(t *testing.T) {
	err := errors.New("some error")
	res := Error[int](err)

	mappedRes := Map(res, func(val int) bool { return val == 0 })

	assert.Equal(t, err, mappedRes.err)
	assert.Equal(t, &err, &mappedRes.err)
	assert.Nil(t, mappedRes.ok)
}

func TestMapErrorNoError(t *testing.T) {
	value := "something"
	res := Ok(value)

	mapped := MapError(res, func(err error) error { return errors.New("some error") })

	assert.Nil(t, res.err)
	assert.Equal(t, value, mapped.ok)
	assert.Same(t, res, mapped)
}

func TestMapErrorWithError(t *testing.T) {
	value := errors.New("something failed")
	res := Error[int](value)

	mapped := MapError(res, func(err error) error { return errors.New(fmt.Sprint("oh no ", err)) })

	assert.NotNil(t, res.err)
	assert.Equal(t, "oh no something failed", mapped.err.Error())
	assert.NotSame(t, res, mapped)
}

func TestBindNoError(t *testing.T) {
	value := 42
	res := Ok(value)

	boundRes := Bind(res, func(val int) *Result[string] { return Ok(fmt.Sprint(val)) })

	assert.Equal(t, "42", boundRes.ok)
	assert.Nil(t, boundRes.err)
}

func TestBindWithError(t *testing.T) {
	err := errors.New("some error")
	res := Error[int](err)

	boundRes := Bind(res, func(val int) *Result[bool] { return Ok(val == 0) })

	assert.Equal(t, err, boundRes.err)
	assert.Equal(t, &err, &boundRes.err)
	assert.Nil(t, boundRes.ok)
}

func TestMatchNoError(t *testing.T) {
	value := 42
	res := Ok(value)

	matchedRes := Match(res,
		func(ok int) string { return fmt.Sprint(ok) },
		func(err error) string { return err.Error() })

	assert.Equal(t, "42", matchedRes)
}

func TestMatchWithError(t *testing.T) {
	errorMessage := "oh shit, here we go again"
	value := errors.New(errorMessage)
	res := Error[int](value)

	matchedRes := Match(res,
		func(ok int) string { return fmt.Sprint(ok) },
		func(err error) string { return err.Error() })

	assert.Equal(t, errorMessage, matchedRes)
}

func TestIsOk(t *testing.T) {
	value := "Are we done yet"
	res := Ok(value)

	assert.True(t, res.IsOk())
}

func TestIsError(t *testing.T) {
	err := errors.New("another error")
	res := Error[string](err)

	assert.True(t, res.IsError())
}

func TestUnwrapOk(t *testing.T) {
	value := 420
	res := Ok(value)

	assert.Equal(t, value, res.Unwrap())
}

func TestUnwrapError(t *testing.T) {
	err := errors.New("its about to get real bad")
	res := Error[int](err)

	assert.Panics(t, func() { res.Unwrap() })
}

func TestStringOk(t *testing.T) {
	expected := "Hi!"
	res := Ok(expected)

	assert.Equal(t, expected, res.String())
}

func TestStringError(t *testing.T) {
	expected := errors.New("hello from the other side")
	res := Error[bool](expected)

	assert.Equal(t, res.err.Error(), res.String())
}

func TestFromTupleOfNoError(t *testing.T) {
	expected := 586
	funcThatReturnsATuple := func() (int, error) { return expected, nil }

	res := FromTupleOf[int](funcThatReturnsATuple())

	assert.Equal(t, expected, res.ok)
	assert.Nil(t, res.err)
}

func TestFromTupleOfWithError(t *testing.T) {
	expected := errors.New("oops")
	funcThatReturnsATuple := func() (int, error) { return 0, expected }

	res := FromTupleOf[int](funcThatReturnsATuple())

	assert.Equal(t, expected, res.err)
	assert.Nil(t, res.ok)
}

func TestContains(t *testing.T) {
	ok := Ok("something")
	err := Error[string](errors.New("error"))

	assert.True(t, Contains(ok, "something"))
	assert.False(t, Contains(err, "other thing"))
}

func TestDefaultValue(t *testing.T) {
	expectedOk := 357
	expectedErr := 357
	ok := Ok(expectedOk)
	err := Error[int](errors.New("error"))

	assert.Equal(t, expectedOk, DefaultValue(ok, 42))
	assert.Equal(t, expectedErr, DefaultValue(err, expectedErr))
}

func TestDefaultWith(t *testing.T) {
	expectedOk := 357
	expectedErr := 0
	ok := Ok(expectedOk)
	err := Error[int](errors.New("error"))

	assert.Equal(t, expectedOk, DefaultWith(ok, func(_ error) int { return 0 }))
	assert.Equal(t, expectedErr, DefaultWith(err, func(_ error) int { return 0 }))
}

func TestExists(t *testing.T) {
	value := 42
	res := Ok(value)
	err := Error[int](errors.New("error"))

	assert.True(t, Exists(res, func(i int) bool { return i > 41 }))
	assert.False(t, Exists(err, func(i int) bool { return i > 41 }))
}

func TestFold(t *testing.T) {
	value := 667
	expected := 777
	expectedErr := 110
	res := Ok(value)
	err := Error[int](errors.New("error"))

	assert.Equal(t, expected, Fold(res, 110, func(s int, i int) int { return s + i }))
	assert.Equal(t, expectedErr, Fold(err, expectedErr, func(s int, i int) int { return s + 1 }))
}

func TestIter(t *testing.T) {
	value := 0
	expected := 1
	res := Ok(&value)
	err := Error[*int](errors.New("error"))
	incrementPtr := func(i *int) unit.Unit { *i++; return unit.Unit{} }

	Iter(res, incrementPtr)

	assert.Equal(t, expected, *(res.ok.(*int)))
	assert.Same(t, &value, res.ok)

	Iter(err, incrementPtr)

	assert.Nil(t, err.ok)
}

func TestToOption(t *testing.T) {
	value := 146
	res := Ok(value)
	err := Error[int](errors.New("error"))

	assert.Equal(t, option.Some[int](value), ToOption(res))
	assert.Equal(t, option.None[int](), ToOption(err))
}
