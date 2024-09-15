package test_test

import (
	"errors"
	"testing"

	"go-mod.ewintr.nl/go-kit/test"
)

func TestTest(t *testing.T) {
	t.Run("assert", func(t *testing.T) {
		condition := true
		test.Assert(t, condition, "expected condition to be true")
	})

	t.Run("ok", func(t *testing.T) {
		var condition error
		test.Assert(t, condition == nil, "expected condition to be true")
		test.OK(t, condition)
	})

	t.Run("not-nil", func(t *testing.T) {
		var condition error
		condition = errors.New("some error here")
		test.NotNil(t, condition)
	})

	t.Run("nil", func(t *testing.T) {
		var condition error
		test.Nil(t, condition)
	})

	t.Run("equals", func(t *testing.T) {
		for _, tc := range []struct {
			message  string
			expected interface{}
			result   interface{}
		}{
			{
				message: "when expected is zero value",
			},
			{
				message:  "when expected is nil",
				expected: nil,
			},
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{"testing"},
				result:   struct{ test string }{"testing"},
			},
			{
				message:  "when expected and result are strings",
				expected: "testing",
				result:   "testing",
			},
		} {
			t.Log(tc.message)
			{
				test.Equals(t, tc.expected, tc.result)
			}
		}
	})

	t.Run("not-zero", func(t *testing.T) {
		for _, tc := range []struct {
			message  string
			expected interface{}
		}{
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{"testing"},
			},
			{
				message:  "when expected and result are strings",
				expected: "testing",
			},
			{
				message:  "when expected and result are integers",
				expected: 1,
			},
		} {
			t.Log(tc.message)
			{
				test.NotZero(t, tc.expected)
			}
		}
	})

	t.Run("zero", func(t *testing.T) {
		for _, tc := range []struct {
			message  string
			expected interface{}
		}{
			{
				message:  "when expected and result are struct",
				expected: struct{ test string }{},
			},
			{
				message:  "when expected and result are strings",
				expected: "",
			},
			{
				message:  "when expected and result are integers",
				expected: 0,
			},
		} {
			t.Log(tc.message)
			{
				test.Zero(t, tc.expected)
			}
		}
	})

	t.Run("includes", func(t *testing.T) {
		result := "The quick brown fox jumps over the lazy dog"
		expected := "jumps"
		test.Includes(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		test.Includes(t, expected, resultList...)
	})

	t.Run("includes-i", func(t *testing.T) {
		result := "The quick brown fox jumps over the lazy dog"
		expected := "JUMPS"
		test.IncludesI(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		test.IncludesI(t, expected, resultList...)
	})

	t.Run("not-includes", func(t *testing.T) {
		result := "The quick brown fox jumps over the lazy dog"
		expected := "hippo"
		test.NotIncludes(t, expected, result)

		resultList := []string{"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
		test.NotIncludes(t, expected, resultList...)
	})

	t.Run("includes-slice", func(t *testing.T) {
		expected := []string{"B"}
		original := []string{"A", "B", "C"}
		test.IncludesSlice(t, expected, original)

		expectedI := []int{5}
		originalI := []int{1, 2, 3, 4, 5, 6, 7}
		test.IncludesSlice(t, expectedI, originalI)

		expectedE := []interface{}{5, "B"}
		originalE := []interface{}{1, 2, 3, 4, 5, 6, 7, "A", "B", "C"}
		test.IncludesSlice(t, expectedE, originalE)
	})

	t.Run("includes-map", func(t *testing.T) {
		expected := map[string]string{"B": "B"}
		original := map[string]string{
			"A": "A",
			"B": "B",
			"C": "C",
		}
		test.IncludesMap(t, expected, original)
	})
}
