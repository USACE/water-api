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

				// handle levels which are stored in an array of structs
				if typ.Field(i).Type.Kind() == reflect.Slice && tagval == "level" {

					recordSlice := reflect.ValueOf(val.Field(i).Interface())

					var level string

					for x := 0; x < recordSlice.Len(); x++ {

						entry := recordSlice.Index(x)
						// fmt.Println("Field entries: ", recordSlice.Len())
						// fmt.Println(entry.Field(0).String())

						fieldKey := entry.Field(0).String()
						fieldVal := entry.Field(1).Interface()

						// fmt.Println(fieldKey, " -> ", fieldVal)
						// fmt.Println(reflect.TypeOf(fieldVal))

						level = fieldKey + "," + fmt.Sprintf("%.2f", fieldVal)

						// fmt.Println("level", level)
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
