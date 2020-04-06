package e

import (
	"errors"
	"fmt"
	"testing"
)

const (
	CodeUnexpected = "unexpected_error"
	CodeDatabase   = "database_error"
	CodeInternal   = "internal_error"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
		want string
	}{
		{
			name: "new constructs error",
			fn: func() error {
				const op = "Foo"

				return New(op, CodeDatabase, "cannot do something")
			},
			want: "Foo: cannot do something [database_error]",
		},
		{
			name: "wrap adds op",
			fn: func() error {
				const op = "Inner"
				err := New(op, CodeInternal, "cannot do something")

				const op2 = "Outer"
				return Wrap(op2, err)
			},
			want: "Outer: Inner: cannot do something [internal_error]",
		},
		{
			name: "wrap adds op and optionalInfo",
			fn: func() error {
				const op = "Inner"
				err := New(op, CodeInternal, "cannot do something")

				const op2 = "Outer"
				return Wrap(op2, err, "optional info here")
			},
			want: "Outer: (optional info here): Inner: cannot do something [internal_error]",
		},
		{
			name: "can wrap non-package errors",
			fn: func() error {
				const op = "Foo"
				err := errors.New("basic error")
				return Wrap(op, err)
			},
			want: "Foo: basic error",
		},
		{
			name: "can non-pkg wrap Error",
			fn: func() error {
				const op = "Inner"
				err := New(op, CodeInternal, "cannot do something")

				const op2 = "Outer"
				err = Wrap(op2, err)

				wrap := fmt.Errorf("not encouraged but compatible: %w", err)

				const op3 = "Outer2"
				return Wrap(op3, wrap)
			},
			want: "Outer2: not encouraged but compatible: Outer: Inner: cannot do something [internal_error]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil && err.Error() != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", err, tt.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		fn   func() string
		want string
	}{
		{
			name: "unset message returns blank",
			fn: func() string {
				op := "Foo"
				err := New(op, CodeUnexpected, "unexpected error occurred")
				return ErrorMessage(err)
			},
			want: "",
		},
		{
			name: "set message returns correctly",
			fn: func() string {
				op := "Foo"
				err := New(op, CodeUnexpected, "unexpected error occurred").SetClientMsg("oh no")
				return ErrorMessage(err)
			},
			want: "oh no",
		},
		{
			name: "multiple messages but outermost message is returned",
			fn: func() string {
				op := "Foo"
				err1 := New(op, CodeUnexpected, "bar").SetClientMsg("don't show this")

				op2 := "Foo2"
				err2 := Wrap(op2, err1).SetClientMsg("show this")

				return ErrorMessage(err2)
			},
			want: "show this",
		},
		{
			name: "works with non-pkg wrapping",
			fn: func() string {
				const op = "Inner"
				err := New(op, CodeInternal, "cannot do something")

				const op2 = "Outer"
				err = Wrap(op2, err).SetClientMsg("wrapped by fmt.Errorf")

				wrap := fmt.Errorf("not encouraged but compatible: %w", err)

				const op3 = "Outer2"
				wrap2 := Wrap(op3, wrap)

				return ErrorMessage(wrap2)
			},
			want: "wrapped by fmt.Errorf",
		},
		{
			name: "cleared message does not get returned",
			fn: func() string {
				const op = "Foo"
				err := New(op, CodeInternal, "fail fail fail").SetClientMsg("clear me!")

				const op2 = "Outer"
				err = Wrap(op2, err).SetClientMsg("clear me too!")

				err = err.ClearClientMsg()
				return ErrorMessage(err)
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestErrorCode(t *testing.T) {
	tests := []struct {
		name string
		fn   func() string
		want string
	}{
		{
			name: "unset code returns blank",
			fn: func() string {
				op := "Foo"
				err := New(op, "", "unexpected error occurred")
				return ErrorCode(err)
			},
			want: "",
		},
		{
			name: "set code with New() returns correctly",
			fn: func() string {
				op := "Foo"
				err := New(op, CodeUnexpected, "unexpected error occurred")
				return ErrorCode(err)
			},
			want: CodeUnexpected,
		},
		{
			name: "set code with Wrap() returns correctly",
			fn: func() string {
				op := "Foo"
				err := errors.New("db error occurred")
				wrap := Wrap(op, err).SetCode(CodeDatabase)
				return ErrorCode(wrap)
			},
			want: CodeDatabase,
		},
		{
			name: "multiple codes but outermost code is returned",
			fn: func() string {
				op := "Foo"
				err1 := New(op, CodeUnexpected, "bar")

				op2 := "Foo2"
				err2 := Wrap(op2, err1).SetCode(CodeInternal)

				return ErrorCode(err2)
			},
			want: CodeInternal,
		},
		{
			name: "works with non-pkg wrapping",
			fn: func() string {
				const op = "Inner"
				err := New(op, CodeInternal, "cannot do something")

				const op2 = "Outer"
				err = Wrap(op2, err)

				wrap := fmt.Errorf("not encouraged but compatible: %w", err)

				const op3 = "Outer2"
				wrap2 := Wrap(op3, wrap)

				return ErrorCode(wrap2)
			},
			want: CodeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(); got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}
