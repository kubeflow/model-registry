package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		// Test int
		intVal := 42
		intPtr := Of(intVal)
		assert.NotNil(t, intPtr)
		assert.Equal(t, intVal, *intPtr)

		// Test string
		strVal := "hello world"
		strPtr := Of(strVal)
		assert.NotNil(t, strPtr)
		assert.Equal(t, strVal, *strPtr)

		// Test bool
		boolVal := true
		boolPtr := Of(boolVal)
		assert.NotNil(t, boolPtr)
		assert.Equal(t, boolVal, *boolPtr)

		// Test float64
		floatVal := 3.14159
		floatPtr := Of(floatVal)
		assert.NotNil(t, floatPtr)
		assert.Equal(t, floatVal, *floatPtr)
	})

	t.Run("zero values", func(t *testing.T) {
		// Test zero int
		zeroInt := 0
		zeroIntPtr := Of(zeroInt)
		assert.NotNil(t, zeroIntPtr)
		assert.Equal(t, 0, *zeroIntPtr)

		// Test empty string
		emptyStr := ""
		emptyStrPtr := Of(emptyStr)
		assert.NotNil(t, emptyStrPtr)
		assert.Equal(t, "", *emptyStrPtr)

		// Test false bool
		falseBool := false
		falseBoolPtr := Of(falseBool)
		assert.NotNil(t, falseBoolPtr)
		assert.Equal(t, false, *falseBoolPtr)
	})

	t.Run("complex types", func(t *testing.T) {
		// Test slice
		sliceVal := []int{1, 2, 3}
		slicePtr := Of(sliceVal)
		assert.NotNil(t, slicePtr)
		assert.Equal(t, sliceVal, *slicePtr)

		// Test map
		mapVal := map[string]int{"a": 1, "b": 2}
		mapPtr := Of(mapVal)
		assert.NotNil(t, mapPtr)
		assert.Equal(t, mapVal, *mapPtr)

		// Test struct
		type testStruct struct {
			Name string
			Age  int
		}
		structVal := testStruct{Name: "John", Age: 30}
		structPtr := Of(structVal)
		assert.NotNil(t, structPtr)
		assert.Equal(t, structVal, *structPtr)
	})

	t.Run("pointer independence", func(t *testing.T) {
		// Verify that Of() creates a copy, not a reference to the original variable
		original := 10
		ptr1 := Of(original)

		// Modify original after creating pointer - pointer should still have original value
		originalValueBeforeChange := original
		original = 20

		// The pointer should contain the value at the time Of() was called
		assert.Equal(t, originalValueBeforeChange, *ptr1, "pointer should contain value from when Of() was called")
		assert.Equal(t, 10, *ptr1, "pointer should be independent of original variable")
		assert.Equal(t, 20, original, "original variable should have new value")

		// Verify that multiple calls create different pointers
		val := 42
		ptr2 := Of(val)
		ptr3 := Of(val)
		assert.Equal(t, *ptr2, *ptr3, "values should be equal")
		assert.NotSame(t, ptr2, ptr3, "pointers should be different instances")
	})
}

func TestIn(t *testing.T) {
	t.Run("non-nil pointers", func(t *testing.T) {
		// Test int pointer
		intVal := 42
		intPtr := &intVal
		result := In(intPtr)
		assert.Equal(t, intVal, result)

		// Test string pointer
		strVal := "hello"
		strPtr := &strVal
		result2 := In(strPtr)
		assert.Equal(t, strVal, result2)

		// Test bool pointer
		boolVal := true
		boolPtr := &boolVal
		result3 := In(boolPtr)
		assert.Equal(t, boolVal, result3)

		// Test float pointer
		floatVal := 2.718
		floatPtr := &floatVal
		result4 := In(floatPtr)
		assert.Equal(t, floatVal, result4)
	})

	t.Run("nil pointers return zero values", func(t *testing.T) {
		// Test nil int pointer
		var nilIntPtr *int
		result := In(nilIntPtr)
		assert.Equal(t, 0, result)

		// Test nil string pointer
		var nilStrPtr *string
		result2 := In(nilStrPtr)
		assert.Equal(t, "", result2)

		// Test nil bool pointer
		var nilBoolPtr *bool
		result3 := In(nilBoolPtr)
		assert.Equal(t, false, result3)

		// Test nil float pointer
		var nilFloatPtr *float64
		result4 := In(nilFloatPtr)
		assert.Equal(t, 0.0, result4)
	})

	t.Run("complex types", func(t *testing.T) {
		// Test slice pointer
		sliceVal := []string{"a", "b", "c"}
		slicePtr := &sliceVal
		result := In(slicePtr)
		assert.Equal(t, sliceVal, result)

		// Test nil slice pointer
		var nilSlicePtr *[]string
		result2 := In(nilSlicePtr)
		assert.Nil(t, result2)

		// Test map pointer
		mapVal := map[string]int{"key": 123}
		mapPtr := &mapVal
		result3 := In(mapPtr)
		assert.Equal(t, mapVal, result3)

		// Test nil map pointer
		var nilMapPtr *map[string]int
		result4 := In(nilMapPtr)
		assert.Nil(t, result4)
	})

	t.Run("struct types", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		// Test struct pointer
		person := Person{Name: "Alice", Age: 25}
		personPtr := &person
		result := In(personPtr)
		assert.Equal(t, person, result)

		// Test nil struct pointer returns zero value
		var nilPersonPtr *Person
		result2 := In(nilPersonPtr)
		expected := Person{} // zero value
		assert.Equal(t, expected, result2)
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("Of then In should preserve values", func(t *testing.T) {
		// Test int
		original := 123
		result := In(Of(original))
		assert.Equal(t, original, result)

		// Test string
		originalStr := "test string"
		resultStr := In(Of(originalStr))
		assert.Equal(t, originalStr, resultStr)

		// Test bool
		originalBool := true
		resultBool := In(Of(originalBool))
		assert.Equal(t, originalBool, resultBool)
	})

	t.Run("In then Of for non-nil pointers", func(t *testing.T) {
		original := 456
		ptr := &original

		// In(ptr) gets the value, Of() creates a new pointer
		newPtr := Of(In(ptr))

		// Values should be equal but pointers should be different
		assert.Equal(t, *ptr, *newPtr)
		assert.NotSame(t, ptr, newPtr)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("pointer to pointer", func(t *testing.T) {
		val := 42
		ptr := &val
		ptrToPtr := &ptr

		// In should dereference once
		result := In(ptrToPtr)
		assert.Equal(t, ptr, result)
		assert.Equal(t, val, *result)
	})

	t.Run("zero values are handled correctly", func(t *testing.T) {
		// Ensure zero values work correctly with both functions
		zeroInt := 0
		ptrToZero := Of(zeroInt)
		backToZero := In(ptrToZero)
		assert.Equal(t, 0, backToZero)

		// Nil pointer should return zero value
		var nilPtr *int
		zeroFromNil := In(nilPtr)
		assert.Equal(t, 0, zeroFromNil)
	})
}
