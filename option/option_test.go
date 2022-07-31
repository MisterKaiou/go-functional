package option

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/unit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNone(t *testing.T) {
	opt := None[unit.Unit]()

	assert.NotNil(t, opt.none)
	assert.Nil(t, opt.some)
}

func TestIsNone(t *testing.T) {
	opt := None[unit.Unit]()

	assert.True(t, opt.IsNone())
}

func TestSome(t *testing.T) {
	value := 777
	opt := Some(value)

	assert.Equal(t, value, opt.some)
	assert.Nil(t, opt.none)
}

func TestIsSome(t *testing.T) {
	opt := Some(42)

	assert.True(t, opt.IsSome())
}

func TestUnwrapSome(t *testing.T) {
	expected := 420
	opt := Some(expected)

	assert.Equal(t, expected, opt.Unwrap())
}

func TestUnwrapError(t *testing.T) {
	opt := None[unit.Unit]()

	assert.Panics(t, func() { opt.Unwrap() })
}

func TestMapSome(t *testing.T) {
	value := 1653
	opt := Some(value)

	mappedOpt := Map(opt, func(it int) string { return fmt.Sprint(it) })

	assert.Equal(t, "1653", mappedOpt.some)
}

func TestMapNone(t *testing.T) {
	opt := None[bool]()

	mappedOpt := Map(opt, func(it bool) string { return fmt.Sprint(it) })

	assert.Nil(t, mappedOpt.some)
}

func TestBindSome(t *testing.T) {
	value := 42
	opt := Some(value)

	boundRes := Bind(opt, func(val int) *Option[string] { return Some(fmt.Sprint(val)) })

	assert.Equal(t, "42", boundRes.some)
	assert.Nil(t, boundRes.none)
}

func TestBindNone(t *testing.T) {
	opt := None[int]()

	boundRes := Bind(opt, func(val int) *Option[string] { return Some(fmt.Sprint(val)) })

	assert.Nil(t, boundRes.some)
	assert.NotNil(t, boundRes.none)
}

func TestMatchSome(t *testing.T) {
	value := 723
	opt := Some(value)

	matchedOpt := Match(opt,
		func(ok int) string { return fmt.Sprint(ok) },
		func() string { return "Should not call this" })

	assert.Equal(t, "723", matchedOpt)
}

func TestMatchNone(t *testing.T) {
	expected := "Should call this"
	opt := None[int]()

	matchedOpt := Match(opt,
		func(ok int) string { return fmt.Sprint(ok) },
		func() string { return expected })

	assert.Equal(t, expected, matchedOpt)
}
