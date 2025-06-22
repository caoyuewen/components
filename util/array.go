package util

import "reflect"

func Include(lst interface{}, val interface{}) bool {
	v := reflect.ValueOf(lst)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Interface() == val {
				return true
			}

		}
	}
	return false
}

func Exclude(lst interface{}, val interface{}) bool {
	v := reflect.ValueOf(lst)
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Interface() == val {
				return false
			}
		}
	}
	return true
}
