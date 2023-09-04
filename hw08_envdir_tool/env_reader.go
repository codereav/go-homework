package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment, len(dirEntries))
	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil || info.IsDir() {
			continue
		}
		envVal := EnvValue{}
		fileName := info.Name()
		if info.Size() == 0 {
			envVal.NeedRemove = true
		} else {
			file, err := os.Open(dir + "/" + fileName)
			if err != nil {
				return nil, fmt.Errorf("could not open file: %w", err)
			}
			// Читаем строку, закрываем файл
			reader := bufio.NewReader(file)
			row, err := reader.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("could not read string: %w", err)
			}
			err = file.Close()
			if err != nil {
				return nil, fmt.Errorf("can't close file: %w", err)
			}
			// Заменяем байт 0 на перенос строки
			val := strings.ReplaceAll(row, string([]byte{0}), string('\n'))
			// Убираем переносы строк, табы и пробелы справа
			val = strings.TrimRight(val, string('\n'))
			val = strings.TrimRight(val, string('\t'))
			val = strings.TrimRight(val, " ")

			envVal.Value = val
		}
		// Убираем пробелы из названия переменной
		varName := strings.ReplaceAll(fileName, "=", "")

		result[varName] = envVal
	}

	return result, nil
}
