package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
	Param string // Параметр для подстановки в текст ошибки
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	err := ""
	for _, val := range v {
		err = err + val.Field + ": " + fmt.Sprintf(val.Err.Error(), val.Param) + "\n"
	}

	return err
}

var allowedTypes = map[string]byte{
	reflect.Int.String():    1,
	reflect.String.String(): 1,
}

var (
	ErrUnsupportedType         = errors.New("type is not supported by validator")
	ErrUnsupportedRule         = errors.New("rule is unsupported")
	ErrIncorrectValidationRule = errors.New("unable to parse validation rule")
	ErrIncorrectRuleValue      = errors.New("unable to parse rule value")
)

var (
	ErrIncorrectMinValue = errors.New("value must be more than %s")
	ErrIncorrectMaxValue = errors.New("value must be less than %s")
)

var (
	ErrIncorrectValueByRegexp = errors.New("value must be satisfy a regexp %s")
	ErrIncorrectValueOneOf    = errors.New("value must be one of: %s")
)

var ErrIncorrectLength = errors.New("value length must be equal %s")

func Validate(v interface{}) error {
	var errors ValidationErrors

	t := reflect.ValueOf(v)
	// Проверяем, что входящий параметр - структура
	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Type().Field(i)

		tag, ok := field.Tag.Lookup("validate")
		// Пропускаем элементы без тэга validate
		if !ok {
			continue
		}

		var curType string
		// Если это слайс, то за текущий тип берём тип элемента слайса
		if field.Type.Kind() == reflect.Slice {
			curType = field.Type.Elem().String()
		} else {
			curType = field.Type.String()
		}

		// Проверяем, поддерживается ли тип поля
		if _, ok := allowedTypes[curType]; !ok {
			errors = append(
				errors,
				ValidationError{Field: field.Name, Err: ErrUnsupportedType},
			)
			continue
		}

		// Собираем все правила для поля
		rules := strings.Split(tag, "|")

		for _, rule := range rules {
			// Правило должно состоять из двух частей: тип и значение
			ruleData := strings.SplitN(rule, ":", 2)
			if len(ruleData) != 2 {
				errors = append(
					errors,
					ValidationError{Field: field.Name, Err: ErrIncorrectValidationRule},
				)
				continue
			}
			ruleType := ruleData[0]
			ruleValue := ruleData[1]

			// Если это слайс - валидируем содержимое этого слайса
			if field.Type.Kind() == reflect.Slice {
				for i := 0; i < t.FieldByName(field.Name).Len(); i++ {
					err := processValidation(ruleType, ruleValue, t.FieldByName(field.Name).Index(i))
					if err != nil {
						errors = append(
							errors,
							ValidationError{Field: field.Name, Err: err, Param: ruleValue},
						)
					}
				}
			} else {
				err := processValidation(ruleType, ruleValue, t.FieldByName(field.Name))
				if err != nil {
					errors = append(
						errors,
						ValidationError{Field: field.Name, Err: err, Param: ruleValue},
					)
				}
			}
		}
	}

	return errors
}

func processValidation(ruleType string, ruleValue string, value reflect.Value) error {
	switch ruleType {
	case "min":
		formattedRuleValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ErrIncorrectRuleValue
		}
		if value.Int() < int64(formattedRuleValue) {
			return ErrIncorrectMinValue
		}
	case "max":
		formattedRuleValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ErrIncorrectRuleValue
		}
		if value.Int() > int64(formattedRuleValue) {
			return ErrIncorrectMaxValue
		}
	case "len":
		formattedRuleValue, err := strconv.Atoi(ruleValue)
		if err != nil {
			return ErrIncorrectRuleValue
		}
		if value.Len() != formattedRuleValue {
			return ErrIncorrectLength
		}
	case "regexp":
		regexp, err := regexp.Compile(ruleValue)
		if err != nil {
			return err
		}
		if !regexp.MatchString(value.String()) {
			return ErrIncorrectValueByRegexp
		}
	case "in":
		ruleValues := strings.Split(ruleValue, ",")
		if !inArray(value.String(), ruleValues) {
			return ErrIncorrectValueOneOf
		}
	default:
		return ErrUnsupportedRule
	}

	return nil
}

func inArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}

	return false
}
