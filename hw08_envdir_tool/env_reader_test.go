package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	testDir := "./testdata/env/"

	// Создаём файл с некорректным названием, чтобы убедиться, что из названия переменной удалится знак =
	brokenVarFile, _ := os.OpenFile("./testdata/env/BR=OKEN", os.O_CREATE|os.O_WRONLY, 0o777)
	brokenVarFile.WriteString("anyValue")
	brokenVarFile.Close()
	defer os.Remove("./testdata/env/BR=OKEN")

	env, err := ReadDir(testDir)
	if err != nil {
		fmt.Println(fmt.Errorf("could not ge env from dir: %w", err))
		return
	}
	require.Equal(t, "   foo\nwith new line", env["FOO"].Value)
	require.Equal(t, "bar", env["BAR"].Value)
	require.Equal(t, `"hello"`, env["HELLO"].Value)
	require.Equal(t, ``, env["EMPTY"].Value)
	require.Equal(t, true, env["UNSET"].NeedRemove)
	require.Equal(t, "anyValue", env["BROKEN"].Value)
}
