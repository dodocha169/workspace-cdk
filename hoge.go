package main

import "fmt"

func newURLSet(count int) []string {
	var set []string
	for i := 1; i <= count; i++ {
		set = append(set, ""+string(i))
	}
	return set
}

func main() {
	fmt.Printf("(%%#v) %#v\n", newURLSet(5))
}
