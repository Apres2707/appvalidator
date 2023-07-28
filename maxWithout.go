package appvalidator

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// maxWithout validates the property value to be less than or equal to the specified threshold if the listed variables
// are empty.
// Paths of dependent properties are described from the current property, the nesting of structure properties is
// separated by a dot.
// Paths of different dependent properties are separated by a space.
// example: `validate:"max_without=Filters.Name Filters.Age Query.Text 1000"`
var maxWithout validator.Func = func(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), " ")

	dependencyFieldList := params[:len(params)-1]
	for _, fieldPath := range dependencyFieldList {
		fieldPathList := strings.Split(fieldPath, ".")
		if isPropertyNotEmpty(fl.Parent(), fieldPathList) {
			return true
		}
	}

	maxValueString := params[len(params)-1]

	return isLessOrEqual(fl.Field(), maxValueString)
}

func isPropertyNotEmpty(val reflect.Value, pathList []string) bool {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if len(pathList) > 0 {
		val = val.FieldByName(pathList[0])
		if val.IsZero() {
			return false
		}
		return isPropertyNotEmpty(val, pathList[1:])
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Map:
		return !val.IsNil() && val.Len() > 0
	case reflect.Pointer, reflect.Interface, reflect.Chan, reflect.Func:
		return !val.IsNil()
	default:
		return !val.IsZero()
	}
}

func isLessOrEqual(value reflect.Value, maxValueString string) bool {
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		maxValue, err := strconv.ParseUint(maxValueString, 10, 64)
		if err != nil {
			return false
		}

		return value.Uint() <= maxValue
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var maxValue int64
		if value.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(maxValueString)
			if err != nil {
				// attempt parsing as an integer assuming nanosecond precision
				maxValue, err = strconv.ParseInt(maxValueString, 10, 64)
				if err != nil {
					return false
				}
			} else {
				maxValue = int64(d)
			}
		} else {
			var err error
			maxValue, err = strconv.ParseInt(maxValueString, 10, 64)
			if err != nil {
				return false
			}
		}

		return value.Int() <= maxValue
	case reflect.String:
		maxValue, err := strconv.ParseInt(maxValueString, 10, 64)
		if err != nil {
			return false
		}

		return int64(utf8.RuneCountInString(value.String())) <= maxValue
	case reflect.Slice, reflect.Map, reflect.Array:
		maxValue, err := strconv.ParseInt(maxValueString, 10, 64)
		if err != nil {
			return false
		}

		return int64(value.Len()) <= maxValue
	case reflect.Float32, reflect.Float64:
		maxValue, err := strconv.ParseFloat(maxValueString, 64)
		if err != nil {
			return false
		}

		return value.Float() <= maxValue
	case reflect.Struct:
		if value.Type() == reflect.TypeOf(time.Time{}) {
			maxValue, err := time.Parse("2006-01-02T15:04:05-07:00", maxValueString)
			if err != nil {
				return false
			}
			t := value.Interface().(time.Time)

			return t.Before(maxValue) || t.Equal(maxValue)
		}
	}

	return false
}
