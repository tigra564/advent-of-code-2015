package main

import (
	"fmt"
)

func applyPreliminaryDamage(player *fighter) {
	// do nothing
}

func main() {
	// Sanity checks
	// fmt.Println(runDemo1())  // await ocWon and 226
	// fmt.Println(runDemo2())  // await ocWon and 641

	result := solve()
	fmt.Println("Result:", result)
}
