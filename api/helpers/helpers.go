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
				vv.Add(tagval, fmt.Sprint(val.Field(i).Interface()))
			}
		}
	}
	return vv
}
