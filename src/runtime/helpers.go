package runtime

import (
	"fmt"
	"strconv"
)

//
// Конвертация
//

func ToInt(val RuntimeVal) IntVal {
	switch v := val.(type) {
	case IntVal:
		return v
	case FloatVal:
		return IntVal{
			Value: int64(v.Value),
		}
	case StringVal:
		i, err := strconv.Atoi(v.Value)
		if err != nil {
			panic(fmt.Errorf("cannot convert string to int: %v", err))
		}
		return IntVal{
			Value: int64(i),
		}
	case BoolVal:
		if v.Value {
			return IntVal{
				Value: 1,
			}
		}
		return IntVal{
			Value: 0,
		}
	default:
		panic(fmt.Errorf("unsupported type for int conversion: %T", val))
	}
}

func ToFloat(val RuntimeVal) FloatVal {
	switch v := val.(type) {
	case IntVal:
		return FloatVal{
			Value: float64(v.Value),
		}
	case FloatVal:
		return v
	case StringVal:
		f, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			panic(fmt.Errorf("cannot convert string to float: %v", err))
		}
		return FloatVal{
			Value: f,
		}
	case BoolVal:
		if v.Value {
			return FloatVal{
				Value: 1.0,
			}
		}
		return FloatVal{
			Value: 0.0,
		}
	default:
		panic(fmt.Errorf("unsupported type for float conversion: %T", val))
	}
}

func ToString(val RuntimeVal) StringVal {
	switch v := val.(type) {
	case IntVal:
		return StringVal{
			Value: fmt.Sprintf("%d", v.Value),
		}
	case FloatVal:
		return StringVal{
			Value: fmt.Sprintf("%f", v.Value),
		}
	case StringVal:
		return v
	case BoolVal:
		return StringVal{
			Value: strconv.FormatBool(v.Value),
		}
	default:
		panic(fmt.Errorf("unsupported type for string conversion: %T", val))
	}
}

func ToBool(val RuntimeVal) BoolVal {
	switch v := val.(type) {
	case IntVal:
		return BoolVal{
			Value: v.Value > 0,
		}
	case FloatVal:
		return BoolVal{
			Value: v.Value > 0,
		}
	case StringVal:
		return BoolVal{
			Value: v.Value != "",
		}
	case BoolVal:
		return v
	default:
		panic(fmt.Errorf("unsupported type for bool conversion: %T", val))
	}
}

//
// Форматирование
//

func Format(val RuntimeVal) string {
	switch valType := val.(type) {
	case IntVal:
		return strconv.FormatInt(valType.Value, 10)
	case FloatVal:
		return strconv.FormatFloat(valType.Value, 'f', 1, 64)
	case StringVal:
		return valType.Value
	case BoolVal:
		if valType.Value {
			return "true"
		} else {
			return "false"
		}
	case NullVal:
		return "null"
	case ArrayVal:
		result := "["
		for i, elem := range valType.Elements {
			result += Format(elem)
			if i+1 < len(valType.Elements) {
				result += ", "
			}
		}
		result += "]"
		return result
	case ObjectVal:
		result := "{"
		i := 0
		for name, elem := range valType.Elements {
			i++
			result += name + ": " + Format(elem)
			if i+1 < len(valType.Elements) {
				result += ", "
			}
		}
		result += "}"
		return result
	case FunctionVal:
		result := valType.Name + "("
		for _, param := range valType.Params {
			result += param + ", "
		}
		result += ")"
		return result
	case NativeFnVal:
		return valType.Name + "()"
	default:
		return fmt.Sprintf("%#v", val)
	}
}

func sprint_format(vals []RuntimeVal) string {
	result := ""
	for _, val := range vals {
		result += Format(val)
	}
	return result
}
