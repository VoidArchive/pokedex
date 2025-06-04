package battle

import (
	"fmt"
	"math/rand"

	"github.com/voidarchive/pokedex/internal/pokeapi"
	"github.com/voidarchive/pokedex/internal/shared/constants"
)

func SimulateBattle(r *rand.Rand, playerPokemon pokeapi.UserPokemon, opponentPokemon pokeapi.PokemonData) (xpGained int) {
	fmt.Printf("\n%s--- Battle Start: %s%s%s %s(Lvl %s%d%s)%s %svs %s%s%s %s---%s\n",
		constants.ColorBrightCyan,
		constants.ColorGreen, playerPokemon.Name, constants.ColorReset,
		constants.ColorYellow,
		constants.ColorBrightCyan, playerPokemon.Level, constants.ColorReset,
		constants.ColorYellow,
		constants.ColorBrightCyan,
		constants.ColorRed, opponentPokemon.Name, constants.ColorReset,
		constants.ColorBrightCyan,
		constants.ColorReset)

	playerHP, _ := playerPokemon.GetStat("hp")
	opponentHP, _ := opponentPokemon.GetStat("hp")
	playerAttack, _ := playerPokemon.GetStat("attack")
	opponentAttack, _ := opponentPokemon.GetStat("attack")
	playerDefense, _ := playerPokemon.GetStat("defense")
	opponentDefense, _ := opponentPokemon.GetStat("defense")
	playerSpeed, _ := playerPokemon.GetStat("speed")
	opponentSpeed, _ := opponentPokemon.GetStat("speed")

	turn := 1
	for playerHP > 0 && opponentHP > 0 {
		fmt.Printf("\n%s--- Turn %s%d%s ---%s\n", constants.ColorCyan, constants.ColorYellow, turn, constants.ColorCyan, constants.ColorReset)
		fmt.Printf("  %s%s HP: %s%d%s | %s%s HP: %s%d%s\n",
			constants.ColorGreen, playerPokemon.Name, constants.ColorBrightGreen, playerHP, constants.ColorReset,
			constants.ColorRed, opponentPokemon.Name, constants.ColorBrightRed, opponentHP, constants.ColorReset)

		p1Name := playerPokemon.Name
		p2Name := opponentPokemon.Name

		var attackerName, defenderName string
		var attackerHP, defenderHP *int
		var attackerAtk, defenderAtk int
		var attackerDef, defenderDef int
		var attackerColor, defenderColor string // Colors for attacker and defender names

		if opponentSpeed > playerSpeed { // Opponent (p2) is faster
			attackerName, defenderName = p2Name, p1Name
			attackerHP, defenderHP = &opponentHP, &playerHP
			attackerAtk, defenderAtk = opponentAttack, playerAttack
			attackerDef, defenderDef = opponentDefense, playerDefense
			attackerColor, defenderColor = constants.ColorRed, constants.ColorGreen // Opponent is Red, Player is Green
		} else { // Player (p1) is faster or speeds are equal
			attackerName, defenderName = p1Name, p2Name
			attackerHP, defenderHP = &playerHP, &opponentHP
			attackerAtk, defenderAtk = playerAttack, opponentAttack
			attackerDef, defenderDef = playerDefense, opponentDefense
			attackerColor, defenderColor = constants.ColorGreen, constants.ColorRed // Player is Green, Opponent is Red
		}

		// First attacker's move
		fmt.Printf("  %s%s%s attacks %s%s%s!\n", attackerColor, attackerName, constants.ColorReset, defenderColor, defenderName, constants.ColorReset)
		damage := CalculateDamage(r, attackerAtk, defenderDef)
		*defenderHP -= damage
		fmt.Printf("  %s%s%s takes %s%d%s damage.\n", defenderColor, defenderName, constants.ColorReset, constants.ColorBrightRed, damage, constants.ColorReset)

		if *defenderHP <= 0 {
			fmt.Printf("  %s%s%s fainted!%s\n", defenderColor, defenderName, constants.ColorRed, constants.ColorReset)
			break
		}

		// Second attacker's move (swap roles)
		// Previous defender is now attacker, previous attacker is now defender
		fmt.Printf("  %s%s%s attacks %s%s%s!\n", defenderColor, defenderName, constants.ColorReset, attackerColor, attackerName, constants.ColorReset) // Original defender (now attacker) attacks original attacker (now defender)
		damage = CalculateDamage(r, defenderAtk, attackerDef)                                                                                          // Original defender's attack vs original attacker's defense
		*attackerHP -= damage
		fmt.Printf("  %s%s%s takes %s%d%s damage.\n", attackerColor, attackerName, constants.ColorReset, constants.ColorBrightRed, damage, constants.ColorReset)

		if *attackerHP <= 0 {
			fmt.Printf("  %s%s%s fainted!%s\n", attackerColor, attackerName, constants.ColorRed, constants.ColorReset)
			break
		}
		turn++
	}

	fmt.Printf("\n%s--- Battle End ---%s\n", constants.ColorBrightCyan, constants.ColorReset)
	if playerHP > 0 {
		fmt.Printf("%s%s%s wins!%s\n", constants.ColorBrightGreen, playerPokemon.Name, constants.ColorBrightGreen, constants.ColorReset)
		xpGained = opponentPokemon.BaseExperience
		if xpGained <= 0 {
			xpGained = 10
		}
	} else {
		fmt.Printf("%s%s%s wins!%s\n", constants.ColorBrightRed, opponentPokemon.Name, constants.ColorBrightRed, constants.ColorReset)
		xpGained = 0
	}
	fmt.Println(constants.ColorBrightCyan + "--------------------" + constants.ColorReset)
	return xpGained
}

func CalculateDamage(r *rand.Rand, attack, defense int) int {
	if defense <= 0 {
		defense = 1
	}
	damage := max((attack*20)/(defense+10), 1)

	randomFactor := 0.8 + r.Float64()*(1.2-0.8)
	damage = max(int(float64(damage)*randomFactor), 1)
	return damage
}
