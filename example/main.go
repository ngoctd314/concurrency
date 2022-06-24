package example

import "fmt"

// Exec ...
func Exec() {
	// execBoring()
	// execFoo()
	// execUpdatePosition()
	// execFibonaciGenerator()
	multily := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}

		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}
	fmt.Println(multily(add([]int{1, 1, 1}, 1), 2))
}
