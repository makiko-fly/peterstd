package util

import (
	"fmt"
	"reflect"
	"strings"
)

func getFieldNames(t reflect.Type) []string {
	fmt.Println(t.Kind())
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fieldNum := t.NumField()
	var fieldsNames = make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		jsonDesc, ok := t.Field(i).Tag.Lookup("json")
		if ok {
			fieldsNames = append(fieldsNames, strings.Split(jsonDesc, ",")[0])
		} else {
			fieldsNames = append(fieldsNames, t.Field(i).Name)
		}
	}

	return fieldsNames
}

// CompressSlice 加强版返回压缩函数
func CompressSlice(slice interface{}) interface{} {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return slice
	}

	s := reflect.ValueOf(slice)

	var fieldNames []string
	var fieldValues [][]interface{}

	fieldNames = getFieldNames(reflect.TypeOf(slice).Elem())

	for i := 0; i < s.Len(); i++ {
		record := reflect.Indirect(s.Index(i))
		numField := record.NumField()

		var values []interface{}
		for n := 0; n < numField; n++ {
			values = append(values, record.Field(n).Interface())
		}

		fieldValues = append(fieldValues, values)
	}

	return &ListListResp{
		Fields: fieldNames,
		Items:  fieldValues,
	}
}

type ListListResp struct {
	Fields []string        `json:"fields"`
	Items  [][]interface{} `json:"items"`
}
