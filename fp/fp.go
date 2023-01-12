package fp

import (
	"fmt"
	"reflect"
	"strings"
)

type Fp struct {
	ConsoleOut bool
	SkipZero   bool
	Prefix     string
	ReplacePkg string
	Out        []string
}

func (f *Fp) printf(format string, args ...any) {
	if f.ConsoleOut {
		fmt.Printf(format, args...)
	}
}

func (f *Fp) FormatPrint(data interface{}, out []string, level int, lastType reflect.Kind) []string {
	valueOf := reflect.ValueOf(data)
	typeOf := reflect.TypeOf(data)

	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
		typeOf = typeOf.Elem()
	}

	switch valueOf.Kind() {
	case reflect.Struct:
		if valueOf.IsZero() && f.SkipZero {
			break
		}

		structName := valueOf.Type().String()
		if f.ReplacePkg != "" && valueOf.Type().PkgPath() == f.ReplacePkg {
			structName = strings.Replace(structName, f.ReplacePkg+".", "", 1)
		}

		if lastType == reflect.Slice {
			structName = ""
		}

		f.printf("%s%s{\n", strings.Repeat(f.Prefix, level), structName)

		var ok bool
		var field reflect.StructField
		out = append(out, structName+"{")
		for i, length := 0, valueOf.NumField(); i < length; i++ {
			fieldName := typeOf.Field(i).Name
			if field, ok = typeOf.FieldByName(fieldName); !ok {
				continue
			}

			value := valueOf.FieldByName(fieldName)
			if value.IsZero() && f.SkipZero {
				continue
			}

			switch value.Kind() {
			case reflect.Struct:
				out = append(out, f.FormatPrint(value.Interface(), out, level+1, reflect.Struct)...)
			case reflect.Slice:
				structName = value.Type().String()
				if f.ReplacePkg != "" && strings.Contains(structName, f.ReplacePkg+".") {
					structName = strings.Replace(structName, f.ReplacePkg+".", "", 1)
				}

				str := fmt.Sprintf("%s: %s{", fieldName, structName)
				f.printf("%s%s\n", strings.Repeat(f.Prefix, level+1), str)

				var temp []string
				out = append(out, str)
				for i := 0; i < value.Len(); i++ {
					if value.Index(i).IsZero() && f.SkipZero {
						continue
					}

					temp = f.FormatPrint(value.Index(i).Interface(), temp, level+2, reflect.Slice)
				}

				out = append(out, strings.TrimRight(strings.Join(temp, ""), ","))
				f.printf("%s},\n", strings.Repeat(f.Prefix, level+1))
				out = append(out, "},")
			case reflect.String:
				str := fmt.Sprintf("%s: \"%+v\"", field.Name, value.Interface())
				f.printf("%s%s,\n", strings.Repeat(f.Prefix, level+1), str)
				out = append(out, str+",")
			default:
				str := fmt.Sprintf("%s: %+v", field.Name, value.Interface())
				f.printf("%s%s,\n", strings.Repeat(f.Prefix, level+1), str)
				out = append(out, str+",")
			}
		}
		out[len(out)-1] = strings.TrimRight(out[len(out)-1], ",")

		if lastType != reflect.Slice {
			out = append(out, "}")
			f.printf("%s}\n", strings.Repeat(f.Prefix, level))
			break
		}

		f.printf("%s},\n", strings.Repeat(f.Prefix, level))
		out = append(out, "},")
	case reflect.Slice:
		structName := valueOf.Type().String()
		if f.ReplacePkg != "" && strings.Contains(structName, f.ReplacePkg+".") {
			structName = strings.Replace(structName, f.ReplacePkg+".", "", 1)
		}

		str := fmt.Sprintf("%s{", structName)
		f.printf("%s%s\n", strings.Repeat(f.Prefix, level+1), str)

		var temp []string
		out = append(out, str)
		for i := 0; i < valueOf.Len(); i++ {
			if valueOf.Index(i).IsZero() && f.SkipZero {
				continue
			}

			temp = f.FormatPrint(valueOf.Index(i).Interface(), temp, level+2, reflect.Slice)
		}

		out = append(out, strings.TrimRight(strings.Join(temp, ""), ","))
		f.printf("%s}\n", strings.Repeat(f.Prefix, level+1))
		out = append(out, "}")
	}

	return out
}
