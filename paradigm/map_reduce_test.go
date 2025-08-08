package paradigm

import (
	"strings"
	"testing"
)

func TestMapReduce(t *testing.T) {
	list := []string{"Abc", "Def", "Ghi"}
	newList := MapStrToStr(list, func(s string) string {
		return strings.ToUpper(s)
	})
	t.Log(newList)

	sum := Reduce(list, func(s string) int {
		return len(s)
	})
	t.Log(sum)
}

func MapStrToStr(arr []string, fn func(s string) string) []string {
	var newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		newArr = append(newArr, fn(arr[i]))
	}

	return newArr
}

func Reduce(arr []string, fn func(s string) int) int {
	s := 0
	for i := 0; i < len(arr); i++ {
		s += fn(arr[i])
	}

	return s
}
