package result

import (
	"errors"
	"fmt"
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
