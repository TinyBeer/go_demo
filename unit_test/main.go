package main

import "strings"

func AddUpper(n int) int {
	res := 0
	for i := 1; i <= n; i++ {
		res += i
	}
	return res
}

func Split(str string, sep string) []string {
	result := []string{}
	i := strings.Index(str, sep)
	for i > -1 {
		result = append(result, str[:i])
		str = str[i+len(sep):]
		i = strings.Index(str, sep)
	}
	result = append(result, str)
	return result
}
