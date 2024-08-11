package assert

import (
	"fmt"
	"reflect"
	"testing"
)

const noArguments = 0

func AreEqual(
	t *testing.T,
	got, expected any,
	msgAndArgs ...any,
) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		message := formatMessage(msgAndArgs...)
		t.Errorf(`assert.Equal failed:
			Got:      %v (type %T)
			Expected: %v (type %T)
			%s`, got, got, expected, expected, message)
	}
}

func formatMessage(msgAndArgs ...any) string {
	if len(msgAndArgs) == noArguments {
		return ""
	}
	return fmt.Sprintf("\nMessage: %s", fmt.Sprint(msgAndArgs...))
}

func AreEqualErrs(t *testing.T, got, expected error, msgAndArgs ...any) {
	t.Helper()

	if got == nil && expected == nil {
		return
	}

	if got == nil || expected == nil {
		t.Errorf(`assert.EqualErrors failed:
			Got:      %v
			Expected: %v
			%s`, got, expected, formatMessage(msgAndArgs...))
		return
	}

	if got.Error() != expected.Error() {
		t.Errorf(`assert.EqualErrors failed:
			Got:      %v
			Expected: %v
			%s`, got.Error(), expected.Error(), formatMessage(msgAndArgs...))
	}
}

func IsNotNil(t *testing.T, value any, msgAndArgs ...any) {
	t.Helper()
	if isNilObject(value) {
		message := formatMessage(msgAndArgs...)
		t.Errorf(`assert.NotNil failed:
			Got:      %v (type %T)
			%s`, value, value, message)
	}
}

func IsNil(t *testing.T, value any, msgAndArgs ...any) {
	t.Helper()
	if !isNilObject(value) {
		message := formatMessage(msgAndArgs...)
		t.Errorf(`assert.IsNil failed:
			Got:      %v (type %T)
			%s`, value, value, message)
	}
}

func isNilObject(object any) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}
