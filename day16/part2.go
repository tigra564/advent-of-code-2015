package main

import (
	"fmt"
)

func main() {
	result := findSue(func(traitsMFCSAM, traitsSue map[string]int) bool {
		for k, v := range traitsSue {
			switch k {
			case "cats", "trees":
				if v <= traitsMFCSAM[k] {
					return false
				}
			case "pomeranians", "goldfish":
				if v >= traitsMFCSAM[k] {
					return false
				}
			default:
				if v != traitsMFCSAM[k] {
					return false
				}
			}
		}
		return true
	})

	fmt.Println("Result:", result)
}
