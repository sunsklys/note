package template

import (
	"math"
	"testing"
)

func TestUnionSet(t *testing.T) {
	testes := []struct {
		name string
		set  [][]int
		want int
	}{
		{"一个集合", [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}, 1},
		{"多个集合", [][]int{{1, 2}, {1, 3}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}, 4},
	}

	for _, test := range testes {
		t.Run(test.name, func(t *testing.T) {
			u := NewUnionSet(math.MaxInt8)
			for _, s := range test.set {
				u.Union(s[0], s[1])
				if u.Find(s[0]) != u.Find(s[1]) {
					t.Fatal("结果: ", u.Father, "期望值: ", test.set)
				}
			}
			t.Log(u.Find(1))
		})
	}
}

type UnionSet struct {
	Father []int
}

func NewUnionSet(n int) *UnionSet {
	u := &UnionSet{
		Father: make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		u.Father[i] = i
	}

	return u
}

func (u *UnionSet) Union(i int, j int) {
	iFather := u.Father[i]
	jFather := u.Father[j]
	u.Father[iFather] = jFather
}

func (u *UnionSet) Find(i int) int {
	if i == u.Father[i] {
		return i
	} else {
		u.Father[i] = u.Find(u.Father[i])
		return u.Father[i]
	}
}
