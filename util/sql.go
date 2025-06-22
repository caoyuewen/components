package util

import (
	"bytes"
	"fmt"
	"io"
)

func OnDuplicate(table string, obj map[string]interface{}) (string, []interface{}) {
	var keys []string
	var values []interface{}

	for k, v := range obj {
		keys = append(keys, k)
		values = append(values, v)
	}

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "INSERT INTO `%s`(", table)
	for i := 0; len(keys) > i; i++ {
		if i > 0 {
			io.WriteString(buf, ", ")
		}
		fmt.Fprintf(buf, "`%s`", keys[i])
	}
	io.WriteString(buf, ") VALUES(")
	for i := 0; len(keys) > i; i++ {
		if i > 0 {
			io.WriteString(buf, ", ")
		}
		io.WriteString(buf, "?")
	}
	io.WriteString(buf, ") ON DUPLICATE KEY UPDATE ")
	for i := 0; len(keys) > i; i++ {
		if i > 0 {
			io.WriteString(buf, ", ")
		}
		k := keys[i]
		fmt.Fprintf(buf, "`%s` = ?", k)
		values = append(values, obj[k])
	}

	return buf.String(), values
}
