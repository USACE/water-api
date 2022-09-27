package helpers

import (
	"fmt"
	"net/url"
	"reflect"
)

func StructToQueryValues(data interface{}) url.Values {
	val, typ := reflect.ValueOf(data), reflect.TypeOf(data)
	vv := url.Values{}
	for i := 0; i < val.NumField(); i++ {

		if tagval, ok := typ.Field(i).Tag.Lookup("querystring"); ok {
			if tagval != "" {

				if typ.Field(i).Type.Kind() == reflect.Slice {

					field := reflect.ValueOf(val.Field(i).Interface())

					var level string

					for x := 1; x < field.Len(); x++ {
						level = ""
						level += field.Index(x - 1).Field(x - 1).String()
						level += ","
						level += fmt.Sprintf("%.2f", field.Index(x-1).Field(x).Float())
						//fmt.Println("level", level)
						vv.Add(tagval, level)
					}

				} else {
					vv.Add(tagval, fmt.Sprint(val.Field(i).Interface()))
				}
			}
		}
	}
	return vv
}
