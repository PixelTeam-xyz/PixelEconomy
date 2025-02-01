package utils

import (
    "reflect"
    "strconv"
    "strings"
)

func Contains(element any, block any) bool {
    blockValue := reflect.ValueOf(block)

    if blockValue.Kind() == reflect.Map {
        return blockValue.MapIndex(reflect.ValueOf(element)).IsValid()
    }

    if blockValue.Kind() == reflect.Array || blockValue.Kind() == reflect.Slice {
        for i := 0; i < blockValue.Len(); i++ {
            if reflect.DeepEqual(blockValue.Index(i).Interface(), element) {
                return true
            }
        }
    }

    return false
}

func HasPrefix(str string, prefix any) bool {
    switch p := prefix.(type) {
    case string:
        return strings.HasPrefix(str, p)
    case rune:
        return strings.HasPrefix(str, string(p))
    case int:
        return strings.HasPrefix(str, strconv.Itoa(p))
    case float64:
        return strings.HasPrefix(str, strconv.FormatFloat(p, 'f', -1, 64))
    case bool:
        return strings.HasPrefix(str, strconv.FormatBool(p))
    case byte:
        return strings.HasPrefix(str, string(p))
    case []byte:
        return strings.HasPrefix(str, string(p))
    default:
        return false
    }
}

func TrimPrefix(str string, prefix any) string {
    switch p := prefix.(type) {
    case string:
        return strings.TrimPrefix(str, p)
    case rune:
        return strings.TrimPrefix(str, string(p))
    case int:
        return strings.TrimPrefix(str, strconv.Itoa(p))
    case float64:
        return strings.TrimPrefix(str, strconv.FormatFloat(p, 'f', -1, 64))
    default:
        return str
    }
}

func TrimSuffix(str string, suffix any) string {
    switch s := suffix.(type) {
    case string:
        return strings.TrimSuffix(str, s)
    case rune:
        return strings.TrimSuffix(str, string(s))
    case int:
        return strings.TrimSuffix(str, strconv.Itoa(s))
    case float64:
        return strings.TrimSuffix(str, strconv.FormatFloat(s, 'f', -1, 64))
    default:
        return str
    }
}
