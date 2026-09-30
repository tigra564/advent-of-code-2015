package main

import (
	"advent-of-code-2015/common"
	"math"
)

type spell struct {
	cost, timer int
	damage, armor, hp, mana int
}

const numSpells = 5

var spells = [numSpells]spell{
	// Here instant effects have timer 1 and actually are applied at the boss's turn.
	// However, for the Part 1 this changes nothing.
	{53, 1, 4, 0, 0, 0},
	{73, 1, 2, 0, 2, 0},
	{113, 6, 0, 7, 0, 0},
	{173, 6, 3, 0, 0, 0},
	{229, 5, 0, 0, 0, 101},
}

// Must be in sync with the spells array.
const (
	spMagicMissile = iota
	spDrain
	spShield
	spPoison
	spRecharge
)


type fighter struct {
	hp, damage, armor, mana int
	// Only timer values for each spell ever available
	spellTimers [numSpells]int
	manaSpent int
}

func readFighter(lines []string) fighter {
	var f fighter

	f.hp = common.ExtractNumbers(lines[0])[0]
	f.damage = common.ExtractNumbers(lines[1])[0]

	return f
}

func applyActiveSpells(boss, player *fighter) {
	// We'll re-apply it if Shield is active
	player.armor = 0

	for n := range player.spellTimers {
		timer := &player.spellTimers[n]
		if *timer > 0 {
			sp := spells[n]

			boss.hp -= sp.damage
			player.hp += sp.hp
			player.armor += sp.armor
			player.mana += sp.mana

			*timer--
		}
	}
}

// Possible round outcomes
const (
	ocLost = -1
	ocInProgress = 0
	ocWon = 1
)

type cacheKey struct {
	playerMana, playerManaSpent int
	playerHP byte
	playerTimers [numSpells]byte
	bossHP byte
}

func makeCacheKey(boss, player *fighter) cacheKey {
	var key cacheKey
	key.playerMana = player.mana
	key.playerManaSpent = player.manaSpent
	key.playerHP = byte(player.hp)
	for n, t := range player.spellTimers {
		key.playerTimers[n] = byte(t)
	}
	key.bossHP = byte(boss.hp)

	return key
}

type roundResult struct {
	outcome int
	manaSpent int
}

type cache struct {
	manaSpentLimit int
	results map[cacheKey]roundResult
}

// Evaluates a round: the player's turn, then the boss's turn
func evalRound(boss, player *fighter, spellNum int) roundResult {
	applyPreliminaryDamage(player)
	if player.hp <= 0 {
		return roundResult{ocLost, player.manaSpent}
	}

	// Apply active spells (instant spells won't apply here)
	applyActiveSpells(boss, player)

	// Player's turn (the caller must check that the spell is not already active)
	// But here we cotrol mana. If it is insufficient for the spell, then do not cast
	if player.mana >= spells[spellNum].cost {
		player.spellTimers[spellNum] = spells[spellNum].timer
		player.mana -= spells[spellNum].cost
		player.manaSpent += spells[spellNum].cost
	} else {
		return roundResult{ocLost, player.manaSpent}
	}

	// Apply active spells (here instant spells apply if casted by the player in the previous turn)
	applyActiveSpells(boss, player)

	if boss.hp <= 0 {
		return roundResult{ocWon, player.manaSpent}
	}

	// Boss's turn
	player.hp -= max(1, boss.damage - player.armor)

	if player.hp <= 0 {
		return roundResult{ocLost, player.manaSpent}
	}

	// Nobody has won yet
	return roundResult{ocInProgress, player.manaSpent}
}

func realFightImpl(boss, player *fighter, memo *cache) roundResult {
	key := makeCacheKey(boss, player)

	if cached, ok := memo.results[key]; ok {
		return cached
	}

	best := roundResult{ocLost, memo.manaSpentLimit}

	for n := range spells {
		// A spell with timer == 1 will wear out before the player will be allowed to cast
		if player.spellTimers[n] > 1 {
			continue
		}
		newBoss := *boss
		newPlayer := *player
		rndRes := evalRound(&newBoss, &newPlayer, n)

		// We are not interested in lost or too long fights
		if rndRes.outcome == ocLost || rndRes.manaSpent > best.manaSpent {
			continue
		}
		if rndRes.outcome == ocInProgress {
			rndRes = realFightImpl(&newBoss, &newPlayer, memo)
		}

		// Regardless whether we get rndRes from evalRound or realFightImpl.
		// Here we test rndRes.manaSpent <= best.manaSpent (which is
		// initially memo.manaSpentLimit), because when we go upwards
		// the recursion tree, we will have them equal.
		if rndRes.outcome == ocWon && rndRes.manaSpent <= best.manaSpent {
			best = rndRes
		}
	}

	memo.results[key] = best
	if best.outcome == ocWon {
		memo.manaSpentLimit = best.manaSpent
	}
	return best
}	

func realFight(boss, player fighter) roundResult {
	memo := cache{math.MaxInt, make(map[cacheKey]roundResult)}
	return realFightImpl(&boss, &player, &memo)
}

func demoFight(boss, player fighter, spellNums []int) roundResult {
	for _, spellNum := range spellNums {
		rndRes := evalRound(&boss, &player, spellNum)
		if rndRes.outcome != ocInProgress {
			return rndRes
		}
	}

	return roundResult{ocInProgress, player.manaSpent}
}

func runDemo1() roundResult {
	var boss fighter
	boss.hp = 13
	boss.damage = 8

	var player fighter
	player.hp = 10
	player.mana = 250

	spellNums := []int{spPoison, spMagicMissile}

	return demoFight(boss, player, spellNums)
}

func runDemo2() roundResult {
	var boss fighter
	boss.hp = 14
	boss.damage = 8

	var player fighter
	player.hp = 10
	player.mana = 250

	spellNums := []int{spRecharge, spShield, spDrain, spPoison, spMagicMissile}

	return demoFight(boss, player, spellNums)
}

func solve() int {
	lines := common.ReadInputLines()
	boss := readFighter(lines)

	var player fighter
	player.hp = 50
	player.mana = 500

	fightRes := realFight(boss, player)
	return fightRes.manaSpent
}
