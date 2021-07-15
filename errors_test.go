package errors

import (
	"testing"
)

func TestContain(t *testing.T) {
	e1 := New("the first error")
	e2 := Wrap(e1, "the second error")
	e3 := Wrap(e2, "the third error")

	if !Contain(e3, e1) {
		t.Errorf("failed to check")
	}

}

func TestStack(t *testing.T) {
	e1 := New("the first error")
	e2 := Wrap(e1, "the second error", Ref("level", "second"))
	e3 := Wrap(e2, "the third error")

	if Stack(e3)[1] != e2 {
		t.Errorf("failed to compare")
	}

}

func TestFileLine(t *testing.T) {
	e1 := New("the first error")
	e2 := Wrap(e1, "the second error", Ref("level", "second"))
	e3 := Wrap(e2, "the third error")

	file, line := FileLine(e3)

	if file != "errors_test.go" || line != 32 {
		t.Errorf("failed to call FileLine")
	}

}

func TestPrimary(t *testing.T) {
	e1 := New("the first error")
	e2 := Wrap(e1, "the second error", Ref("level", "second"))
	e3 := Wrap(e2, "the third error")

	if Primary(e3) != e1 {
		t.Errorf("failed to call Primary")
	}

}

func TestUnwrap(t *testing.T) {
	e1 := New("the first error")
	e2 := Wrap(e1, "the second error", Ref("level", "second"))
	e3 := Wrap(e2, "the third error")

	if Unwrap(e3) != e2 {
		t.Errorf("failed to call Unwrap")
	}

}
