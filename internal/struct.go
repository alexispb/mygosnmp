package internal

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/alexispb/mygosnmp/generics"
	"github.com/alexispb/mygosnmp/oid"
)

func StructString(s interface{}, depth int) string {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return "not struct"
	}

	var sb strings.Builder
	fprintStruct(&sb, v, depth)
	return strings.Trim(sb.String(), "\n")
}

func tryFprintOrString(sb *strings.Builder, v reflect.Value) bool {
	if method := v.MethodByName("Fprint"); method.IsValid() {
		t := method.Type()
		if t.NumIn() == 1 && t.In(0).String() == "*strings.Builder" && t.NumOut() == 0 {
			method.Call([]reflect.Value{reflect.ValueOf(sb)})
			return true
		}
	}
	if method := v.MethodByName("String"); method.IsValid() {
		t := method.Type()
		if t.NumIn() == 0 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.String {
			res := method.Call([]reflect.Value{})[0].String()
			sb.WriteString(res)
			return true
		}
	}
	return false
}

func fprintStruct(sb *strings.Builder, v reflect.Value, depth int) {
	prefix := Prefix(depth)
	for i, n, typ := 0, v.NumField(), v.Type(); i < n; i++ {
		sb.WriteByte('\n')
		sb.WriteString(prefix)
		sb.WriteString(typ.Field(i).Name)
		sb.WriteString(": ")
		fprintValue(sb, v.Field(i), depth+1)
	}
}

func fprintValue(sb *strings.Builder, v reflect.Value, depth int) {
	if tryFprintOrString(sb, v) {
		return
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(strconv.FormatInt(v.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(strconv.FormatUint(v.Uint(), 10))

	case reflect.Bool:
		sb.WriteString(strconv.FormatBool(v.Bool()))

	case reflect.String:
		sb.WriteString(v.String())

	case reflect.Struct:
		fprintStruct(sb, v, depth)

	case reflect.Slice:
		switch {
		case v.Len() == 0:
			sb.WriteString("<empty>")
		case v.Type().Elem().Kind() == reflect.Uint32:
			oid.Fprint(sb, v.Interface().([]uint32))
		default:
			prefix := Prefix(depth)
			for i := 0; i < v.Len(); i++ {
				sb.WriteByte('\n')
				sb.WriteString(prefix)
				fprintValue(sb, v.Index(i), depth+1)
			}
		}

	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			sb.WriteString("<nil>")
		} else {
			fprintValue(sb, v.Elem(), depth)
		}

	default:
		sb.WriteString("unsupported ")
		sb.WriteString(v.Kind().String())
	}
}

const (
	diffmark1 = '\u250C'
	diffmark2 = '\u2514'
)

func StringLinesDiff(s1, s2 string) string {
	if s1 == s2 {
		return ""
	}

	lines1 := strings.Split(s1, "\n")
	lines2 := strings.Split(s2, "\n")
	maxlen := generics.Max(len(lines1), len(lines2))

	var sb strings.Builder

	for i := 0; i < maxlen; i++ {
		line1 := ""
		if i < len(lines1) {
			line1 = lines1[i]
		}
		line2 := ""
		if i < len(lines2) {
			line2 = lines2[i]
		}
		if line1 == line2 {
			sb.WriteByte(' ')
			sb.WriteString(line1)
		} else {
			sb.WriteRune(diffmark1)
			sb.WriteString(line1)
			sb.WriteByte('\n')
			sb.WriteRune(diffmark2)
			sb.WriteString(line2)
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
