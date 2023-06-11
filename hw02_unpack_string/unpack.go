package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Unpack(s string) (string, error) {
	var result strings.Builder
	runes := []rune(s)
	var prev rune
	isEscaped := false
	for i, current := range runes {
		if i == 0 && unicode.IsDigit(current) {
			return "", getErr("Строка не должна начинаться с цифр")
		}
		if i > 0 {
			prev = runes[i-1]
		}

		// Переходим к следующему символу, если текущий символ - неэкранированный бэкслэш
		if isBackslash(current) && !isEscaped {
			isEscaped = true
			continue
		}

		if i > 2 && unicode.IsDigit(current) && unicode.IsDigit(prev) && runes[i-2] != '\\' {
			return "", getErr("Не должно быть несколько неэкранированных цифр подряд")
		}

		if isEscaped {
			isEscaped = false
			if unicode.IsDigit(current) || isBackslash(current) {
				result.WriteString(string(current))
				continue
			}
			return "", getErr("Экранировать можно только цифры или бэкслэш")
		}

		// Убираем предыдущий символ, если текущий равен 0
		if string(current) == "0" {
			tmp := result.String()
			result.Reset()
			result.WriteString(tmp[:(len(tmp) - 1)])
			continue
		}

		// Дублируем предыдущий символ n-1 раз
		if unicode.IsDigit(current) {
			n, _ := strconv.Atoi(string(current))
			if n == 0 {
				continue
			}
			result.WriteString(strings.Repeat(string(prev), n-1))
			continue
		}
		result.WriteString(string(current))
	}

	return result.String(), nil
}

func isBackslash(r rune) bool {
	return r == '\\'
}

func getErr(s string) error {
	return errors.New(s)
}
