package main

import (
	"fmt"
)

func applyPreliminaryDamage(player *fighter) {
	player.hp--
}

func main() {
	result := solve()
	fmt.Println("Result:", result)
}
