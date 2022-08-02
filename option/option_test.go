package option

import (
	"fmt"
	"github.com/MisterKaiou/go-functional/unit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringSome(t *testing.T) {
	opt := Some(42)

	assert.Equal(t, "42", opt.String())
}

func TestStringNone(t *testing.T) {
	opt := None[unit.Unit]()

	assert.Equal(t, "None", opt.String())
}

func TestNone(t *testing.T) {
	opt := None[unit.Unit]()

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

	boundRes := Bind(opt, func(val int) Option[string] { return Some(fmt.Sprint(val)) })

	assert.Equal(t, "42", boundRes.some)
}

func TestBindNone(t *testing.T) {
	opt := None[int]()

	boundRes := Bind(opt, func(val int) Option[string] { return Some(fmt.Sprint(val)) })

	assert.Nil(t, boundRes.some)
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

func TestContains(t *testing.T) {
	expected := 262
	opt := Some(expected)
	optNone := None[int]()

	assert.True(t, Contains(opt, expected))
	assert.False(t, Contains(optNone, expected))
}

func TestDefaultValue(t *testing.T) {
	expectedSome := 357
	expectedNone := 357
	some := Some(expectedSome)
	none := None[int]()

	assert.Equal(t, expectedSome, DefaultValue(some, 42))
	assert.Equal(t, expectedNone, DefaultValue(none, expectedNone))
}

func TestDefaultWith(t *testing.T) {
	expectedSome := 357
	expectedNone := 0
	some := Some(expectedSome)
	none := None[int]()

	assert.Equal(t, expectedSome, DefaultWith(some, func() int { return 0 }))
	assert.Equal(t, expectedNone, DefaultWith(none, func() int { return 0 }))
}

func TestExists(t *testing.T) {
	value := 42
	opt := Some(value)
	none := None[int]()

	assert.True(t, Exists(opt, func(i int) bool { return i > 41 }))
	assert.False(t, Exists(none, func(i int) bool { return i > 41 }))
}

func TestFilter(t *testing.T) {
	value := 42
	opt := Some(value)
	none := None[int]()
	predicate := func(i int) bool { return i == 42 }

	filterRetSome := Filter(opt, predicate)

	assert.Equal(t, filterRetSome, opt)

	filterRetNone := Filter(none, predicate)

	assert.Equal(t, None[int](), filterRetNone)

}

func TestFold(t *testing.T) {
	value := 667
	expected := 777
	expectedErr := 110
	res := Some(value)
	err := None[int]()

	assert.Equal(t, expected, Fold(res, 110, func(s int, i int) int { return s + i }))
	assert.Equal(t, expectedErr, Fold(err, expectedErr, func(s int, i int) int { return s + 1 }))
}

func TestIter(t *testing.T) {
	value := 0
	expected := 1
	opt := Some(&value)
	none := None[*int]()
	incrementPtr := func(i *int) unit.Unit { *i++; return unit.Unit{} }

	Iter(opt, incrementPtr)

	assert.Equal(t, expected, *(opt.some.(*int)))
	assert.Same(t, &value, opt.some)

	Iter(none, incrementPtr)

	assert.Nil(t, none.some)
}
