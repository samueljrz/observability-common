package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customError struct {
	msg string
}

func (err customError) Error() string {
	return err.msg
}

func TestGetErrorName(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no error",
			args: args{},
			want: "",
		},
		{
			name: "new std error",
			args: args{
				err: errors.New("an error occurred in the test"),
			},
			want: "error",
		},
		{
			name: "custom error",
			args: args{
				err: customError{msg: "a custom error occurred in the test"},
			},
			want: "customError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorName(tt.args.err)
			assert.Equalf(t, got, tt.want, "GetErrorName() = %v, want %v", got, tt.want)
		})
	}
}
