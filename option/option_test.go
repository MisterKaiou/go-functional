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

	boundRes := Bind(opt, func(val int) Of[string] { return Some(fmt.Sprint(val)) })

	assert.Equal(t, "42", boundRes.some)
}

func TestBindNone(t *testing.T) {
	opt := None[int]()

	boundRes := Bind(opt, func(val int) Of[string] { return Some(fmt.Sprint(val)) })

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
	some := Some(value)
	none := None[int]()

	assert.Equal(t, expected, Fold(some, 110, func(s int, i int) int { return s + i }))
	assert.Equal(t, expectedErr, Fold(none, expectedErr, func(s int, i int) int { return s + 1 }))
}

func TestFoldTo(t *testing.T) {
	res := Some(3500)
	expected := "Number is: 3500"

	foldedResult := Fold(res, "Number is: ", func(s string, i int) string {
		return fmt.Sprint(s, i)
	})

	assert.Equal(t, expected, foldedResult)
}

func TestFoldM(t *testing.T) {
	value := 667
	expected := 777
	expectedErr := 110
	res := Some(value)
	err := None[int]()

	foldSome := FoldM(res, 110, func(s int, i int) int { return s + i })

	assert.Equal(t, expected, foldSome.some)

	foldNone := FoldM(err, expectedErr, func(s int, i int) int { return s + 1 })

	assert.Equal(t, err, foldNone)
}

func TestCombineByNoError(t *testing.T) {
	left := Some(42)
	right := Some("nice")

	combined := CombineBy(left, right, func(s string, i int) bool {
		return s == "nice" || i == 42
	})

	assert.True(t, combined.IsSome())
	assert.True(t, combined.some.(bool))
}

func TestCombineByWithError(t *testing.T) {
	combiningFunc := func(s string, i int) bool {
		return s == "nice" || i == 42
	}

	left := None[int]()
	right := Some("nice")

	combined := CombineBy(left, right, combiningFunc)

	assert.True(t, combined.IsNone())
	assert.Nil(t, combined.some)

	left = Some(69)
	right = None[string]()

	combined = CombineBy(left, right, combiningFunc)

	assert.True(t, combined.IsNone())
	assert.Nil(t, combined.some)

}

func TestFlatten(t *testing.T) {
	toFlatten := Some(Some(42))
	expectedSome := Some(42)

	flattenRes := Flatten(toFlatten)

	assert.Equal(t, expectedSome, flattenRes)
}

func TestFlattenError(t *testing.T) {
	okErr := Some(None[int]())

	inner := Flatten(okErr)

	assert.True(t, okErr.IsSome())
	assert.True(t, inner.IsNone())

	errErr := None[Of[int]]()

	inner = Flatten(errErr)

	assert.True(t, inner.IsNone())
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

func BenchmarkFoldVsMapVsFoldM(b *testing.B) {
	b.Run("Fold", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			some := Some(0)

			Fold(some, i, func(st int, it int) int { return st + it })
		}
	})

	b.Run("Map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			some := Some(0)

			Map(some, func(it int) int { return it + i })
		}
	})

	b.Run("FoldM", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			some := Some(0)

			FoldM(some, i, func(st int, it int) int { return st + it })
		}
	})
}
