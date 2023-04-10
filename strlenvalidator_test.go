package strlenvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}

	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "valid min",
			args: args{
				v: struct {
					Foo string `validate:"len:0,10"`
				}{""},
			},
			wantErr: false,
		},
		{
			name: "valid non bound",
			args: args{
				v: struct {
					Foo string `validate:"len:0,10"`
				}{"123"},
			},
			wantErr: false,
		},
		{
			name: "invalid = max",
			args: args{
				v: struct {
					Foo string `validate:"len:0,10"`
				}{"1234567890"},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "invalid > max",
			args: args{
				v: struct {
					Foo string `validate:"len:0,10"`
				}{"1234567890000"},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "invalid not a num",
			args: args{
				v: struct {
					Foo string `validate:"len:asdsadasd"`
				}{"1234567890000"},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
		{
			name: "invalid args num",
			args: args{
				v: struct {
					Foo string `validate:"len:1"`
				}{"1234567890000"},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 1)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
