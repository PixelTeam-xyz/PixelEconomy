package PrettyPrint

import (
	"fmt"
	"reflect"
	"strings"
)

func PrettyPrint(value interface{}) {
	fmt.Println(ToPrettyStr(value, 0))
}

func PrettyPrintX(value interface{}, PrintlnFn func(msgs ...interface{})) {
	PrintlnFn(ToPrettyStr(value, 0))
}

// Recursive formating
func ToPrettyStr(value interface{}, depth int) string {
	indent := strings.Repeat("  ", depth)

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf("\033[32m\"%s\"\033[0m", v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("\033[92m%d\033[0m", v.Int())
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("\033[92m%f\033[0m", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("\033[38;5;208m%t\033[0m", v.Bool())
	case reflect.Slice, reflect.Array:
		var elems []string
		for i := 0; i < v.Len(); i++ {
			elems = append(elems, ToPrettyStr(v.Index(i).Interface(), depth+1))
		}
		return fmt.Sprintf("[\n%s%s\n]", indent+"  "+strings.Join(elems, ",\n"+indent+"  "), indent)
	case reflect.Map:
		var elms []string
		for _, key := range v.MapKeys() {
			elms = append(elms, fmt.Sprintf("%s%s: %s",
				indent+"  ",
				ToPrettyStr(key.Interface(), depth+1),
				ToPrettyStr(v.MapIndex(key).Interface(), depth+1)))
		}
		return fmt.Sprintf("{\n%s\n%s}", strings.Join(elms, ",\n"), indent)
	case reflect.Struct:
		var elems []string
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			fieldValue := v.Field(i).Interface()
			elems = append(elems, fmt.Sprintf("%s%s: %s",
				indent+"  ",
				fieldName,
				ToPrettyStr(fieldValue, depth+1)))
		}
		return fmt.Sprintf("{\n%s\n%s}", strings.Join(elems, ",\n"), indent)
	case reflect.Ptr:
		if v.IsNil() {
			return "\033[38;5;208mnil\033[0m"
		}
		return ToPrettyStr(v.Elem().Interface(), depth)
	default:
		return fmt.Sprintf("%v", value)
	}
}
