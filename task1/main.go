package main

import "fmt"

func Filter(s []int, isIncluded func(elem, _ int) bool) (res []int) {
	for i := 0; i < len(s); i++ {
		if isIncluded(s[i], i) {
			res = append(res, s[i])
		}
	}
	return
}

func main() {
	fmt.Println("Even", Filter([]int{1, 2, 3, 4, 5}, func(item, index int) bool { return item%2 == 0 }))
}
