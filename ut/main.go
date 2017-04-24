package main

import (
	"fmt"
	//"os"
	"strings"
)


func FormatOutput(field1 string, field2 string, total_len int) string {
	//中间至少保留 3个"-"的占位符
	output := make([]byte, total_len)
	len1 := len(field1)
	len2 := len(field2)

	j := 0
	for i := 0; i < total_len; i++ {

		if i < len1 {
			output[i] = field1[i]
			continue
		}
		//至少填充3个占位符
		if i < len1 + 3 {
			output[i] = '-'
			continue
		}

		if i < total_len - len2 {
			output[i] = '-'
			continue
		}

		output[i] = field2[j]
		j++
	}

	return  string(output)
}

func main() {
	relative_name := "img/img-test/test/test"
	//生成本地目录时，要把远程目录去掉
	fmt.Println(strings.TrimPrefix(relative_name, "img/"))
}