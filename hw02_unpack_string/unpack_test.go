package hw02unpackstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},

		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},

		{input: `О 2тус!!!\\\3`, expected: `О  тус!!!\3`},
		{input: `№;%!"№%:?*()_=а3\\\3`, expected: `№;%!"№%:?*()_=ааа\3`},
		{input: `\\№\\\\\\0`, expected: `\№\\`},
		{input: `\\Капсом Т\\ОЖ3Е Можно :)\0`, expected: `\Капсом Т\ОЖЖЖЕ Можно :)0`},
		{input: `\\♪2◕‿◕5`, expected: `\♪♪◕‿◕◕◕◕◕`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "zxcz123", "aaa10b", "zz43", "456", "\\dd"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Error(t, err, "actual error %q", err)
		})
	}
}
