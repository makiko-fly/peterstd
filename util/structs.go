package util

import (
	"reflect"
	"strings"
)

func StructFieldsMap(s interface{}, fields ...string) map[string]interface{} {
	var (
		res     = make(map[string]interface{})
		sv      = reflect.Indirect(reflect.ValueOf(s))
		st      = sv.Type()
		slen    = st.NumField()
		tagname = "json"
	)

	for i := 0; i < slen; i++ {
		if st.Field(i).Anonymous {
			submap := StructFieldsMap(sv.Field(i).Interface(), fields...)
			for k, v := range submap {
				res[k] = v
			}
		} else {
			tag := st.Field(i).Tag.Get(tagname)
			tag = strings.SplitN(tag, ",", 2)[0]
			res[tag] = sv.Field(i).Interface()
		}
	}

	if len(fields) == 0 {
		return res
	}

	fieldsMap := make(map[string]bool)
	for _, field := range fields {
		fieldsMap[field] = true
	}

	// Remove not need fields
	for k := range res {
		if ok := fieldsMap[k]; !ok {
			delete(res, k)
		}
	}

	// Add null fields
	for _, field := range fields {
		if _, ok := res[field]; !ok {
			res[field] = nil
		}
	}
	return res
}
