package paradigm

import (
	"fmt"
	"reflect"
	"testing"
)

type data struct {
	num   int
	slice []int
}

func TestDeepEqual(t *testing.T) {
	v1 := data{
		slice: []int{1, 2, 3},
		num:   1,
	}
	v2 := data{
		num:   1,
		slice: []int{1, 2, 3},
	}
	fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))

	m1 := map[string]string{"one": "a", "two": "b"}
	m2 := map[string]string{"two": "b", "one": "a"}
	fmt.Println("m1 == m2:", reflect.DeepEqual(m1, m2))

	s1 := []int{1, 2, 3}
	s2 := []int{1, 2, 3}
	fmt.Println("s1 == s2:", reflect.DeepEqual(s1, s2))
}
