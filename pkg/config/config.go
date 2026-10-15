package config

import (
	"errors"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func Get[V GenericType](key string, defaultValue ...V) V {
	v, err := GenericParseType[V](os.Getenv(key))
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return v
}

type GenericType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~bool | ~string | ~[]byte
}

func GenericParseType[V GenericType](str string) (V, error) {
	var v V
	switch any(v).(type) {
	case int:
		result, err := strconv.ParseInt(str, 10, 0)
		return any(int(result)).(V), err
	case int8:
		result, err := strconv.ParseInt(str, 10, 8)
		return any(int8(result)).(V), err
	case int16:
		result, err := strconv.ParseInt(str, 10, 16)
		return any(int16(result)).(V), err
	case int32:
		result, err := strconv.ParseInt(str, 10, 32)
		return any(int32(result)).(V), err
	case int64:
		result, err := strconv.ParseInt(str, 10, 64)
		return any(int64(result)).(V), err
	case uint:
		result, err := strconv.ParseUint(str, 10, 0)
		return any(uint(result)).(V), err
	case uint8:
		result, err := strconv.ParseUint(str, 10, 8)
		return any(uint8(result)).(V), err
	case uint16:
		result, err := strconv.ParseUint(str, 10, 16)
		return any(uint16(result)).(V), err
	case uint32:
		result, err := strconv.ParseUint(str, 10, 32)
		return any(uint32(result)).(V), err
	case uint64:
		result, err := strconv.ParseUint(str, 10, 64)
		return any(uint64(result)).(V), err
	case float32:
		result, err := strconv.ParseFloat(str, 32)
		return any(float32(result)).(V), err
	case float64:
		result, err := strconv.ParseFloat(str, 64)
		return any(float64(result)).(V), err
	case bool:
		result, err := strconv.ParseBool(str)
		return any(result).(V), err
	case string:
		if str == "" {
			return v, errors.New("parsed result is empty string")
		}
		return any(str).(V), nil
	case []byte:
		if str == "" {
			return v, errors.New("parsed result is empty bytes")
		}
		return any([]byte(str)).(V), nil
	default:
		return v, errors.New("unsupported generic type")
	}
}
