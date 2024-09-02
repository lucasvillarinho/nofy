package assert

import (
	"errors"
	"testing"
)

func TestAreEqual(t *testing.T) {
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
		message  string
	}{
		{
			name:     "should return true when integers are equal",
			got:      42,
			expected: 42,
			message:  "Expected integers to be equal",
		},
		{
			name:     "should return true when strings are equal",
			got:      "hello",
			expected: "hello",
			message:  "Expected strings to be equal",
		},
		{
			name:     "should return true when slices are equal",
			got:      []int{1, 2, 3},
			expected: []int{1, 2, 3},
			message:  "Expected slices to be equal",
		},
		{
			name:     "should return true when maps are equal",
			got:      map[string]int{"a": 1, "b": 2},
			expected: map[string]int{"a": 1, "b": 2},
			message:  "Expected maps to be equal",
		},
		{
			name: "should return true when structs are equal",
			got: struct {
				Field1 string
				Field2 int
			}{
				Field1: "hello",
				Field2: 42,
			},
			expected: struct {
				Field1 string
				Field2 int
			}{
				Field1: "hello",
				Field2: 42,
			},
			message: "Expected structs to be equal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AreEqual(t, tt.got, tt.expected, tt.message)
		})
	}
}

func TestAreEqualFails(t *testing.T) {
	tests := []struct {
		name        string
		got         interface{}
		expected    interface{}
		expectedErr string
	}{
		{
			name:        "should fail when integers are not equal",
			got:         42,
			expected:    43,
			expectedErr: "Equal failed:\n\t\t\tGot:      42 (type int)\n\t\t\tExpected: 43 (type int)\n\t\t\t",
		},
		{
			name:        "should fail when strings are not equal",
			got:         "hello",
			expected:    "world",
			expectedErr: "Equal failed:\n\t\t\tGot:      hello (type string)\n\t\t\tExpected: world (type string)\n\t\t\t",
		},
		{
			name:        "should fail when slices are not equal",
			got:         []int{1, 2, 3},
			expected:    []int{1, 2, 4},
			expectedErr: "Equal failed:\n\t\t\tGot:      [1 2 3] (type []int)\n\t\t\tExpected: [1 2 4] (type []int)\n\t\t\t",
		},
		{
			name:        "should fail when maps are not equal",
			got:         map[string]int{"a": 1, "b": 2},
			expected:    map[string]int{"a": 1, "b": 3},
			expectedErr: "Equal failed:\n\t\t\tGot:      map[a:1 b:2] (type map[string]int)\n\t\t\tExpected: map[a:1 b:3] (type map[string]int)\n\t\t\t",
		},
		{
			name: "should fail when structs are not equal",
			got: struct {
				Field1 string
				Field2 int
			}{
				Field1: "hello",
				Field2: 42,
			},
			expected: struct {
				Field1 string
				Field2 int
			}{
				Field1: "world",
				Field2: 42,
			},
			expectedErr: "Equal failed:\n\t\t\tGot:      {hello 42} (type struct { Field1 string; Field2 int })\n\t\t\tExpected: {world 42} (type struct { Field1 string; Field2 int })\n\t\t\t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			AreEqual(subT, tt.got, tt.expected)
			if !subT.Failed() {
				t.Fatalf("Expected test %s to fail, but it didn't", tt.name)
			}
		})
	}
}

func TestAreEqualErrsSuccess(t *testing.T) {
	tests := []struct {
		got         error
		expectedErr error
		name        string
	}{
		{
			name:        "should pass when both errors are nil",
			got:         nil,
			expectedErr: nil,
		},
		{
			name:        "should pass when errors are the same",
			got:         errors.New("same error"),
			expectedErr: errors.New("same error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			AreEqualErrs(subT, tt.got, tt.expectedErr)

			if subT.Failed() {
				t.Fatalf("Expected test %s to pass, but it failed", tt.name)
			}
		})
	}
}

func TestAreEqualErrsFails(t *testing.T) {
	tests := []struct {
		got         error
		expectedErr error
		name        string
	}{
		{
			name:        "should fail when got is nil and expected is not nil",
			got:         nil,
			expectedErr: errors.New("expected error"),
		},
		{
			name:        "should fail when expected is nil and got is not nil",
			got:         errors.New("got error"),
			expectedErr: nil,
		},
		{
			name:        "should fail when errors are different",
			got:         errors.New("got error"),
			expectedErr: errors.New("expected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			AreEqualErrs(subT, tt.got, tt.expectedErr)

			if !subT.Failed() {
				t.Fatalf("Expected test %s to fail, but it didn't", tt.name)
			}
		})
	}
}

func TestIsNotNilSuccess(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{
			name:  "Non-nil integer",
			value: 42,
		},
		{
			name:  "Non-nil string",
			value: "hello",
		},
		{
			name:  "Non-nil slice",
			value: []int{1, 2, 3},
		},
		{
			name:  "Non-nil map",
			value: map[string]int{"a": 1, "b": 2},
		},
		{
			name:  "Non-nil struct",
			value: struct{ Field1 string }{Field1: "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			IsNotNil(subT, tt.value)

			if subT.Failed() {
				t.Fatalf("Expected test %s to pass, but it failed", tt.name)
			}
		})
	}
}

func TestIsNotNilFails(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{
			name:  "should fail when value is a nil interface",
			value: nil,
		},
		{
			name:  "should fail when value is a nil pointer",
			value: (*int)(nil),
		},
		{
			name:  "should fail when value is a nil slice",
			value: ([]int)(nil),
		},
		{
			name:  "should fail when value is a nil map",
			value: (map[string]int)(nil),
		},
		{
			name:  "should fail when value is a nil channel",
			value: (chan int)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			IsNotNil(subT, tt.value)

			if !subT.Failed() {
				t.Fatalf("Expected test %s to fail, but it didn't", tt.name)
			}
		})
	}
}

func TestIsNilSuccess(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{
			name:  "Nil interface",
			value: nil,
		},
		{
			name:  "Nil pointer",
			value: (*int)(nil),
		},
		{
			name:  "Nil slice",
			value: ([]int)(nil),
		},
		{
			name:  "Nil map",
			value: (map[string]int)(nil),
		},
		{
			name:  "Nil channel",
			value: (chan int)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			IsNil(subT, tt.value)

			if subT.Failed() {
				t.Fatalf("Expected test %s to pass, but it failed", tt.name)
			}
		})
	}
}

func TestIsNilFails(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{
			name:  "should pass when value is a nil interface",
			value: 42,
		},
		{
			name:  "should pass when value is a nil pointer",
			value: "hello",
		},
		{
			name:  "should pass when value is a nil slice",
			value: []int{1, 2, 3},
		},
		{
			name:  "should pass when value is a nil map",
			value: map[string]int{"a": 1, "b": 2},
		},
		{
			name:  "should pass when value is a nil channel",
			value: struct{ Field1 string }{Field1: "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subT := &testing.T{}
			IsNil(subT, tt.value)

			if !subT.Failed() {
				t.Fatalf("Expected test %s to fail, but it didn't", tt.name)
			}
		})
	}
}
