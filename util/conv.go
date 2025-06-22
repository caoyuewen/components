package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ToJsonE(v interface{}) (string, error) {
	if v == nil {
		return "null", nil
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ToJson(v interface{}) string {
	if v == nil {
		return "null"
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ToStringE(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	switch x := v.(type) {
	case string:
		return x, nil
	case []byte:
		return string(x), nil
	case fmt.Stringer:
		return x.String(), nil
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func ToInt64E(v interface{}) (int64, error) {
	if v == nil {
		return 0, nil
	}
	switch x := v.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		return toInt64ViaFloat64(x), nil
	case []byte:
		return ToInt64E(string(x))
	case string:
		x = strings.TrimSpace(x)
		if x == "" {
			return 0, nil
		}
		if strings.Contains(x, ".") {
			f, err := strconv.ParseFloat(x, 64)
			if err != nil {
				return 0, err
			}
			return int64(f), nil
		}
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func ToIntE(v interface{}) (int, error) {
	i64, err := ToInt64E(v)
	if err != nil {
		return 0, err
	}
	if i64 > math.MaxInt || i64 < math.MinInt {
		return 0, errors.New("int64 out of int range")
	}
	return int(i64), nil
}

func ToFloat64E(v interface{}) (float64, error) {
	if v == nil {
		return 0, nil
	}
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return toFloat64(x)
	case bool:
		if x {
			return 1.0, nil
		}
		return 0.0, nil
	case []byte:
		return strconv.ParseFloat(strings.TrimSpace(string(x)), 64)
	case string:
		return strconv.ParseFloat(strings.TrimSpace(x), 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func ToBoolE(v interface{}) (bool, error) {
	if v == nil {
		return false, nil
	}
	switch x := v.(type) {
	case bool:
		return x, nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		f, _ := ToFloat64E(x)
		return f != 0, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(x)) {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off", "":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean string: %s", x)
		}
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func ToString(v interface{}) string {
	s, _ := ToStringE(v)
	return s
}

func ToInt64(v interface{}) int64 {
	i, _ := ToInt64E(v)
	return i
}

func ToInt(v interface{}) int {
	i, _ := ToIntE(v)
	return i
}

func ToFloat64(v interface{}) float64 {
	f, _ := ToFloat64E(v)
	return f
}

func ToBool(v interface{}) bool {
	b, _ := ToBoolE(v)
	return b
}

func MustToString(v interface{}) string {
	s, err := ToStringE(v)
	if err != nil {
		panic(err)
	}
	return s
}

func MustToInt64(v interface{}) int64 {
	i, err := ToInt64E(v)
	if err != nil {
		panic(err)
	}
	return i
}

func toFloat64(v interface{}) (float64, error) {
	switch x := v.(type) {
	case int:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
		//// 注意：可能超过 float64 精度范围
		//if x > math.MaxFloat64 {
		//	return 0, fmt.Errorf("uint64 %d overflows float64", x)
		//}
		//return float64(x), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func toInt64ViaFloat64(v interface{}) int64 {
	f, err := toFloat64(v)
	if err != nil {
		return 0
	}
	if f > float64(math.MaxInt64) {
		return math.MaxInt64
	}
	if f < float64(math.MinInt64) {
		return math.MinInt64
	}
	return int64(f)
}
