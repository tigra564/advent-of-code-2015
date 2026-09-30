package main

import (
	"advent-of-code-2015/common"
	"math"
)

type item struct {
	cost, damage, armor int
}

var weapons = []item{
	{8, 4, 0},
	{10, 5, 0},
	{25, 6, 0},
	{40, 7, 0},
	{74, 8, 0},
}

var armor = []item{
	/* armor is optional */
	{0, 0, 0},

	{13, 0, 1},
	{31, 0, 2},
	{53, 0, 3},
	{75, 0, 4},
	{102, 0, 5},
}

var rings = []item{
	/* both rings are optional */
	{0, 0, 0},
	{0, 0, 0},

	{25, 1, 0},
	{50, 2, 0},
	{100, 3, 0},
	{20, 0, 1},
	{40, 0, 2},
	{80, 0, 3},
}

type fighter struct {
	hp, damage, armor int
}

func readFighter(lines []string) fighter {
	var f fighter

	f.hp = common.ExtractNumbers(lines[0])[0]
	f.damage = common.ExtractNumbers(lines[1])[0]
	f.armor = common.ExtractNumbers(lines[2])[0]

	return f
}

func evalFight(boss, player fighter, playerInv []item) (int, bool) {
	gear := item{0, 0, 0}
	for _, elem := range playerInv {
		gear.cost += elem.cost
		gear.damage += elem.damage
		gear.armor += elem.armor
	}

	player.damage += gear.damage
	player.armor += gear.armor

	damageToPlayer := max(1, boss.damage - player.armor)
	damageToBoss := max(1, player.damage - boss.armor)

	// Zero HP means death, so a fighter must stay with at least 1 HP to win
	playerRoundsToDeath :=  (player.hp - 1) / damageToPlayer
	bossRoundsToDeath := (boss.hp - 1) / damageToBoss

	return gear.cost, playerRoundsToDeath >= bossRoundsToDeath
}

func solve(mustWin bool) int {
	lines := common.ReadInputLines()
	boss := readFighter(lines)
	player := fighter{100, 0, 0}

	result := math.MaxInt
	evalFn := func (a, b int) int { return min(a, b) }

	if !mustWin {
		result = math.MinInt
		evalFn = func (a, b int) int { return max(a, b) }
	}

	playerInv := make([]item, 4, 4)
	for _, wp := range weapons {
		playerInv[0] = wp
		for _, ar := range armor {
			playerInv[1] = ar
			for n, rn1 := range rings {
				for _, rn2 := range rings[n+1:] {
					playerInv[2] = rn1
					playerInv[3] = rn2
					cost, win := evalFight(boss, player, playerInv)
					if mustWin == win {
						result = evalFn(result, cost)
					}
				}
			}
		}
	}

	return result
}
