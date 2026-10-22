package types

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

// check is a small set of test assertions on the standard library: assert
// reports a failure and carries on, require stops the test.
type check struct{ fatal bool }

var (
	assert  = check{}
	require = check{fatal: true}
)

func (c check) fail(t *testing.T, msg string, extra []any) {
	t.Helper()
	if len(extra) > 0 {
		msg += ": " + fmt.Sprint(extra...)
	}
	if c.fatal {
		t.Fatal(msg)
	}
	t.Error(msg)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	}
	return false
}

func (c check) Equal(t *testing.T, want, got any, extra ...any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		c.fail(t, fmt.Sprintf("got %#v, want %#v", got, want), extra)
	}
}

func (c check) True(t *testing.T, v bool, extra ...any) {
	t.Helper()
	if !v {
		c.fail(t, "expected true", extra)
	}
}

func (c check) False(t *testing.T, v bool, extra ...any) {
	t.Helper()
	if v {
		c.fail(t, "expected false", extra)
	}
}

func (c check) NoError(t *testing.T, err error, extra ...any) {
	t.Helper()
	if err != nil {
		c.fail(t, fmt.Sprintf("unexpected error: %v", err), extra)
	}
}

func (c check) Error(t *testing.T, err error, extra ...any) {
	t.Helper()
	if err == nil {
		c.fail(t, "expected an error", extra)
	}
}

func (c check) NotNil(t *testing.T, v any, extra ...any) {
	t.Helper()
	if isNil(v) {
		c.fail(t, "expected a non-nil value", extra)
	}
}

func (c check) InDelta(t *testing.T, want, got, delta float64, extra ...any) {
	t.Helper()
	if math.IsNaN(got) || math.Abs(want-got) > delta {
		c.fail(t, fmt.Sprintf("got %v, want %v ± %v", got, want, delta), extra)
	}
}

func (c check) Len(t *testing.T, v any, n int, extra ...any) {
	t.Helper()
	if l := reflect.ValueOf(v).Len(); l != n {
		c.fail(t, fmt.Sprintf("length %d, want %d", l, n), extra)
	}
}

func (c check) NotEmpty(t *testing.T, v any, extra ...any) {
	t.Helper()
	if isNil(v) || reflect.ValueOf(v).IsZero() || lenOf(v) == 0 {
		c.fail(t, "expected a non-empty value", extra)
	}
}

func lenOf(v any) int {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len()
	}
	return 1
}

// contains reports whether s (a string, slice, array or map) holds elem:
// a substring, an element or a key.
func contains(s, elem any) bool {
	if str, ok := s.(string); ok {
		sub, ok := elem.(string)
		return ok && strings.Contains(str, sub)
	}
	rv := reflect.ValueOf(s)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if reflect.DeepEqual(rv.Index(i).Interface(), elem) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			if reflect.DeepEqual(k.Interface(), elem) {
				return true
			}
		}
	}
	return false
}

func (c check) Contains(t *testing.T, s, elem any, extra ...any) {
	t.Helper()
	if !contains(s, elem) {
		c.fail(t, fmt.Sprintf("%#v does not contain %#v", s, elem), extra)
	}
}

func (c check) NotContains(t *testing.T, s, elem any, extra ...any) {
	t.Helper()
	if contains(s, elem) {
		c.fail(t, fmt.Sprintf("%#v should not contain %#v", s, elem), extra)
	}
}
