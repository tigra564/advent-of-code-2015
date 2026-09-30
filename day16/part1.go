package main

import (
	"fmt"
)

func main() {
	result := findSue(func(traitsMFCSAM, traitsSue map[string]int) bool {
		for k, v := range traitsSue {
			if v != traitsMFCSAM[k] {
				return false
			}
		}
		return true
	})

	fmt.Println("Result:", result)
}
