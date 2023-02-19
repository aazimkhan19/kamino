package kamino_test

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/LastPossum/kamino"
)

type simpleStruct struct {
	Int     int
	Float64 float64
	String  string
}

type alltogether struct {
	Int                 int
	Float64             float64
	String              string
	ArrayOfInt          [10]int
	ArrayOfSimpleStruct [5]simpleStruct
	SliceOfInt          []int
	SliceOfSimpleStruct []simpleStruct
	Nested              simpleStruct

	IntPtrs    []*int
	StructPtrs []*simpleStruct
}

var (
	intInstance1 = 1114
	intInstance2 = 1387

	simpleStructIntance = simpleStruct{-1, -1, "-1"}

	alltogetherInstance = alltogether{
		Int:        10,
		Float64:    20.,
		String:     "30",
		ArrayOfInt: [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		ArrayOfSimpleStruct: [5]simpleStruct{
			{1, 1, "1"},
			{2, 2, "2"},
			{3, 3, "3"},
			{4, 4, "4"},
			{5, 5, "5"},
		},
		SliceOfInt: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		SliceOfSimpleStruct: []simpleStruct{
			{1, 1, "1"},
			{2, 2, "2"},
			{3, 3, "3"},
			{4, 4, "4"},
			{5, 5, "5"},
		},
		Nested: simpleStruct{
			1, 2, "3",
		},
		IntPtrs:    []*int{&intInstance1, &intInstance2, &intInstance1},
		StructPtrs: []*simpleStruct{&simpleStructIntance, &simpleStructIntance},
	}
)

func TestClone2(t *testing.T) {
	t.Run("primitive types", func(t *testing.T) {
		var (
			i64  int64      = 1
			f64  float64    = 1.0
			b               = true
			c128 complex128 = 4i + 1
			s               = "blah-blah"
		)

		goti64, _ := kamino.Clone(i64)
		assert.Equal(t, i64, goti64)

		gotf64, _ := kamino.Clone(f64)
		assert.Equal(t, f64, gotf64)

		gotb, _ := kamino.Clone(b)
		assert.Equal(t, b, gotb)

		gotc128, _ := kamino.Clone(c128)
		assert.Equal(t, c128, gotc128)

		gots, _ := kamino.Clone(s)
		assert.Equal(t, s, gots)
	})

	t.Run("arrays", func(t *testing.T) {
		arrInts := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		arrIntsClone, _ := kamino.Clone(arrInts)
		assert.Equal(t, arrInts, arrIntsClone)

		arrAny := [...]any{0, 1, 2, 3, nil, 5, 6, &simpleStructIntance, 8, 9, 10, 11, 12, 13, 14, 15}

		gotArrAny, _ := kamino.Clone(arrAny)
		assert.Equal(t, arrAny, gotArrAny)
	})

	t.Run("slices", func(t *testing.T) {
		sliceInts := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		sliceIntsClone, _ := kamino.Clone(sliceInts)
		assert.Equal(t, sliceInts, sliceIntsClone)

		sliceAny := []any{0, 1, 2, 3, nil, 5, 6, &simpleStructIntance, 8, 9, 10, 11, 12, 13, 14, 15}

		gotSliceAny, _ := kamino.Clone(sliceAny)
		assert.Equal(t, sliceAny, gotSliceAny)

		sliceOfSlice := [][]int{{1, 2, 4}, {4, 5}, {6}}
		gotSliceOfSlice, _ := kamino.Clone(sliceOfSlice)
		assert.Equal(t, sliceOfSlice, gotSliceOfSlice)
		sliceOfSlice[0][1] = 7
		assert.NotEqual(t, sliceOfSlice, gotSliceOfSlice)
	})

	t.Run("map", func(t *testing.T) {
		m := make(map[string]int)
		m["1"] = 1
		m["2"] = 2
		m["3"] = 3

		got, _ := kamino.Clone(m)
		assert.Equal(t, m, got)

		m2 := make(map[[3]int][]int)
		m2[[3]int{1, 2, 3}] = []int{1, 2, 3}
		m2[[3]int{2, 3}] = []int{2, 3}
		m2[[3]int{3}] = []int{3}

		got2, _ := kamino.Clone(m2)
		assert.Equal(t, m2, got2)

		i := 1
		m3 := map[string]*int{
			"1": &i,
		}
		got3, _ := kamino.Clone(m3)
		assert.Equal(t, m3, got3)

		*got3["1"] = 10
		assert.NotEqual(t, m3, got3)
	})

	t.Run("simple struct", func(t *testing.T) {
		type simpleStruct struct {
			Int     int
			Float64 float64
			String  string
		}
		param := simpleStruct{
			10, 20., "30",
		}

		got, _ := kamino.Clone(param)
		assert.Equal(t, param, got)
	})

	t.Run("pointer", func(t *testing.T) {
		var (
			i64     int64 = 1
			sstruct       = simpleStruct{2, 3, "4"}

			i64Ptr     = &i64
			sstructPtr = &sstruct
		)

		i64PtrClone, _ := kamino.Clone(i64Ptr)
		assert.Equal(t, *i64Ptr, *i64PtrClone)

		*i64PtrClone = -1
		assert.Equal(t, *i64Ptr, int64(1))

		sstructPtrClone, _ := kamino.Clone(sstructPtr)
		assert.Equal(t, sstructPtr, sstructPtrClone)

		sstructPtrClone.Int = -2
		sstructPtrClone.Float64 = -3
		sstructPtrClone.String = "-4"

		assert.Equal(t, sstructPtr.Int, 2)
		assert.Equal(t, sstructPtr.Float64, float64(3))
		assert.Equal(t, sstructPtr.String, "4")
	})

	t.Run("ptrCircles", func(t *testing.T) {
		type circled struct {
			PtrA *int
			PtrB *int

			This *circled
		}

		i := 10
		c := circled{PtrA: &i, PtrB: &i}
		c.This = &c

		cPtrClone, _ := kamino.Clone(&c)

		assert.Equal(t, cPtrClone, &c)
		assert.Equal(t, unsafe.Pointer(cPtrClone), unsafe.Pointer(cPtrClone.This))
		assert.Equal(t, unsafe.Pointer(cPtrClone.PtrA), unsafe.Pointer(cPtrClone.PtrB))
	})

	t.Run("alltogether", func(t *testing.T) {
		got, _ := kamino.Clone(alltogetherInstance)
		assert.Equal(t, alltogetherInstance, got)
	})

	t.Run("original is not affected by clone mutation", func(t *testing.T) {
		original := alltogether{
			Int:        1,
			Float64:    2,
			String:     "3",
			ArrayOfInt: [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			SliceOfInt: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			Nested:     simpleStruct{1, 2, "3"},
		}

		originalHandCopy := alltogether{
			Int:        1,
			Float64:    2,
			String:     "3",
			ArrayOfInt: [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			SliceOfInt: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			Nested:     simpleStruct{1, 2, "3"},
		}
		assert.Equal(t, original, originalHandCopy)

		clone, err := kamino.Clone(original)
		assert.NoError(t, err)

		assert.Equal(t, original, clone)
		assert.Equal(t, original, originalHandCopy)

		clone.Int = 2
		clone.Float64 = 3
		clone.String = "another string"
		clone.ArrayOfInt[5] = -1
		clone.SliceOfInt[5] = -5

		assert.NotEqual(t, original, clone)
		assert.Equal(t, original, originalHandCopy)
	})
}

func TestInterface(t *testing.T) {
	x := []interface{}{nil}
	y, _ := kamino.Clone(x)
	assert.Equal(t, x, y)

	var a interface{}
	b, _ := kamino.Clone(a)
	assert.True(t, a == b)
}

func TestTwoNils(t *testing.T) {
	type Foo struct {
		A int
	}
	type Bar struct {
		B int
	}
	type FooBar struct {
		Foo  *Foo
		Bar  *Bar
		Foo2 *Foo
		Bar2 *Bar
	}

	src := &FooBar{
		Foo2: &Foo{1},
		Bar2: &Bar{2},
	}

	dst, _ := kamino.Clone(src)

	assert.Equal(t, src, dst)
}

func ptrTo[T any](v T) *T {
	return &v
}

func TestCloneUnexported(t *testing.T) {
	type Foo struct {
		A    int
		a    int
		Ptr1 *int
		ptr2 *int
	}

	foo := Foo{
		1, 2, ptrTo(3), ptrTo(4),
	}

	bar, _ := kamino.Clone(foo)
	assert.Equal(t, foo, bar)

	assert.NotEqual(t, unsafe.Pointer(foo.Ptr1), unsafe.Pointer(bar.Ptr1))
	assert.Equal(t, unsafe.Pointer(foo.ptr2), unsafe.Pointer(bar.ptr2))
}

func TestForceCloneUnexported(t *testing.T) {
	type Foo struct {
		A    int
		a    int
		Ptr1 *int
		ptr2 *int
	}

	foo := Foo{
		1, 2, ptrTo(3), ptrTo(4),
	}

	bar, _ := kamino.Clone(foo, kamino.WithForceUnexported())
	assert.Equal(t, foo, bar)

	assert.NotEqual(t, unsafe.Pointer(foo.Ptr1), unsafe.Pointer(bar.Ptr1))
	assert.NotEqual(t, unsafe.Pointer(foo.ptr2), unsafe.Pointer(bar.ptr2))
}

type Fooer interface {
	foo() int
}

type fooer struct {
	i int
}

func (f *fooer) foo() int {
	return f.i
}

func TestCopyInterface(t *testing.T) {
	type fooerWrapper struct {
		F Fooer
	}

	fi := &fooer{i: 10}

	fw := fooerWrapper{
		F: fi,
	}

	got, _ := kamino.Clone(fw)

	assert.Equal(t, got, fw)
	fi.i = 20
	assert.NotEqual(t, got.F.foo(), fw.F.foo())
}

func TestCopyNestedTime(t *testing.T) {
	type nestedTime struct {
		T time.Time
	}

	nt := nestedTime{time.Now()}
	got, _ := kamino.Clone(nt)

	assert.Equal(t, got.T, nt.T)
}

func TestCopyNestedNil(t *testing.T) {
	type nestedNil struct {
		X any
	}

	nn := nestedNil{}
	got, _ := kamino.Clone(nn)

	assert.Equal(t, got, nn)
}