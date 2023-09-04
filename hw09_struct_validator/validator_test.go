package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
	Custom1 struct {
		Any int `validate:"incorrectRule:anything"`
	}
	Custom2 struct {
		Any int `validate:"regexp"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123412341234123412341234123412341234",
				Age:    123,
				Phones: []string{"123123", "12345678910"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrIncorrectMaxValue,
					Param: "50",
				},
				ValidationError{
					Field: "Email",
					Err:   ErrIncorrectValueByRegexp,
					Param: "^\\w+@\\w+\\.\\w+$",
				},
				ValidationError{
					Field: "Role",
					Err:   ErrUnsupportedType,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrIncorrectLength,
					Param: "11",
				},
			},
		},
		{
			in: App{
				Version: "versi",
			},
			expectedErr: ValidationErrors(nil),
		},
		{
			in:          Token{},
			expectedErr: ValidationErrors(nil),
		},
		{
			in: Response{
				Code: 444,
				Body: "anything",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrIncorrectValueOneOf,
					Param: "200,404,500",
				},
			},
		},
		{
			in: Custom1{},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Any",
					Err:   ErrUnsupportedRule,
					Param: "anything",
				},
			},
		},
		{
			in: Custom2{},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Any",
					Err:   ErrIncorrectValidationRule,
					Param: "",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			fmt.Println(err)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
