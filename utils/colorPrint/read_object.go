package colorPrint

import (
	"fmt"
	"reflect"
	"strings"
)

func readObject(object any, index int) {
    indent := "  "
    useIndent := func(index int) string {
        if index == 0 {
            return ""
        }
        return strings.Repeat(indent, index)
    }

    val := reflect.ValueOf(object)
    kind := val.Kind()

    switch kind {
    case reflect.Map:
        fmt.Println(RedP("Type"), ": map")
        for _, key := range val.MapKeys() {
            value := val.MapIndex(key)
            fmt.Println(useIndent(index), BlueP(fmt.Sprintf("%v", key)), ": ")
            readObject(value.Interface(), index+1)
        }
    case reflect.Slice:
        fmt.Println(useIndent(index), BlueP("[]"+val.Type().Elem().Name()), GreyP("{}"))
        for i := 0; i < val.Len(); i++ {
            value := val.Index(i)
            fmt.Println(useIndent(index), BlueP(fmt.Sprintf("[%d]", i)), ": ")
            readObject(value.Interface(), index+1)
            fmt.Println(useIndent(index), GreyP("}"))
        }
    case reflect.Struct:

        for i := 0; i < val.NumField(); i++ {
            field := val.Type().Field(i)
            value := val.Field(i)
            
            
            //Si la valeur n'est pas un objet quelque soit le type, on affiche la valeur
            if value.Kind() != reflect.Struct && (value.Kind() == reflect.String ||
                value.Kind() == reflect.Int ||
                value.Kind() == reflect.Uint ||
                value.Kind() == reflect.Bool) {
                fmt.Println(useIndent(index), BlueP(field.Name), ": ", determineType(value.Interface()))
                continue
            }

            fmt.Println(useIndent(index), BlueP(field.Name), GreenP(value.Type().Name()), ": ")
            readObject(value.Interface(), index+1)
        }
    default:
        fmt.Println(RedP("Type in Else "), ": ", kind)
        fmt.Println(RedP("Type"), ": ", kind)
        fmt.Println(useIndent(index), BlueP("Value"), ": ", object)
    }
}

func determineType(object any) any {
    val := reflect.ValueOf(object)
    kind := val.Kind()
    switch kind {
    case reflect.String:
        return GreyP(`string `) + val.String()
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return GreyP(kind.String()+" ") + YellowP(fmt.Sprintf("%d", val.Int()))
    default:
        println(RedP("Type in determineType Else of  determineType"), kind, kind == reflect.Struct)
        return object
    }
}


func ObjectLog(content ...any) {
    readObject(content, 0)
}