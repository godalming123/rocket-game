package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

/////// main.go //////////////////////////////////////////////////////////////
// Only the code for the rocket game, using abstractions from engine.go for //
// readability.                                                             //
//////////////////////////////////////////////////////////////////////////////

type gameState struct {
	rocketPos     uint
	playerScore   uint
	currentRocket int
	terminal      *terminalManager
}

const (
	rocketWidth        = 13
	rocketLeftPointing = iota - 2
	rocketNormal
	rocketRightPointing
)

func (gameData *gameState) drawRocket() {
	maxRocketPos := uint(gameData.terminal.noOfColumns - rocketWidth)
	if gameData.rocketPos >= maxRocketPos {
		gameData.rocketPos = maxRocketPos
	} else if gameData.rocketPos < 1 {
		gameData.rocketPos = 1
	}
	switch gameData.currentRocket {
	case rocketNormal:
		drawMultilineText(
			[]string{
				`    \=====/    `,
				`     | S |     `,
				`  ___| C |___  `,
				`  \  | O |  /  `,
				`   \ | R | /   `,
				`    \| E |/    `,
				fmt.Sprintf(`     |%3d|     `, gameData.playerScore),
				`     \   /     `,
				`      \ /      `,
				`       *       `,
			},
			gameData.rocketPos, 1) // Draw the rocket
	case rocketLeftPointing:
		drawMultilineText(
			[]string{
				`       |=====/ `,
				`       / S /   `,
				`  |-__/ C /    `,
				`  |  / O /-__  `,
				`  | / R /  _/  `,
				`  |/ E / _/    `,
				fmt.Sprintf(`  /%3d/_/      `, gameData.playerScore),
				` /   /         `,
				`/  _/          `,
				`|_/            `,
			},
			gameData.rocketPos, 1) // Draw the rocket
	case rocketRightPointing:
		drawMultilineText(
			[]string{
				` \=====|       `,
				`   \ S \       `,
				`    \ C \__-|  `,
				`  __-\ O \  |  `,
				`  \_  \ R \ |  `,
				`    \_ \ E \|  `,
				fmt.Sprintf(`      \_\%3d\  `, gameData.playerScore),
				`         \   \ `,
				`          \_  \`,
				`            \_|`,
			},
			gameData.rocketPos, 1) // Draw the rocket
	}
}

func (game *gameState) rocketFly() {
	game.playerScore++
	if game.playerScore >= 999 {
		fmt.Println("\033[2J\033[HCongratulations, you appear to have a player score longer then the rocket's counter can comprehend. Good job!")
		os.Exit(0)
	}
	if game.playerScore % 2 == 0 {
		switch rand.Intn(2) {
			case 0:
				drawMultilineText([]string{
					" _____",
					`/     \`,
					"|    _/",
					` \__/`,
				}, uint(rand.Intn(int(game.terminal.noOfColumns))), uint(game.terminal.noOfLines)-3)
			case 1:
				drawMultilineText([]string{
					"  __    ___",
					` /  \__/   \`,
					"|       ___/",
					` \_____/`,
				}, uint(rand.Intn(int(game.terminal.noOfColumns))), uint(game.terminal.noOfLines)-3)
			case 2:
				drawMultilineText([]string{
					` ___`,
					`/   \`,
					`\___/`,
				}, uint(rand.Intn(int(game.terminal.noOfColumns))), uint(game.terminal.noOfLines)-2)
			}
	}
	fmt.Print("\033[9999;9999H\n")
	game.drawRocket()
}

func main() {
	// Init //
	var game gameState                     // Create the game object
	fmt.Print("\033[?25l")                 // Hide the cursor
	game.terminal = setupTerminal(func() { // Setup stuff for handling the terminal
		// This function is ran now and subsequently on every terminal resize
		drawMultilineText(
			[]string{
				" Language: \033[34mgo\033[0m | License: \033[34mMIT\033[0m | Author: \033[34mgodalming123\033[0m",
				"\033[93m _______ _           \033[91m _____            _        _",
				"\033[93m|__   __| |          \033[91m|  __ \\          | |      | |",
				"\033[93m   | |  | |__   ___  \033[91m| |__) |___   ___| | _____| |_",
				"\033[93m   | |  | '_ \\ / _ \\ \033[91m|  _  // _ \\ / __| |/ / _ \\ __|",
				"\033[93m   | |  | | | |  __/ \033[91m| | \\ \\ (_) | (__|   <  __/ |_",
				"\033[93m   |_|  |_| |_|\\___| \033[91m|_|  \\_\\___/ \\___|_|\\_\\___|\\__|",
				"\033[92m           / ____|",
				"           | |  __  __ _ _ __ ___   ___",
				"           | | |_ |/ _` | '_ ` _ \\ / _ \\",
				"           | |__| | (_| | | | | | |  __/",
				`            \_____|\__,_|_| |_| |_|\___|`,
				"              \033[0mPress any key to start"},
			uint(game.terminal.noOfColumns/2)-26, uint(game.terminal.noOfLines/2)-6)
	})
	_, _, err := game.terminal.reader.ReadRune() // Wait until the user presses a key
	if err != nil {                              // Handle errors if necarsarry
		log.Fatal(err)
	}

	// Create the actual game //
	fmt.Print("\033[2J") // Clear the screen
	game.terminal.resizeFunction = func() {
		game.drawRocket()
	}

	// Loop to handle keypresses //
	go func() {
		for {
			rune, _, err := game.terminal.reader.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			switch rune {
			case 's':
				if game.currentRocket == rocketNormal {
					game.rocketFly()
				} else {
					game.currentRocket = rocketNormal
				}
			case 'a':
				if game.currentRocket == rocketLeftPointing {
					game.rocketPos--
				} else {
					game.currentRocket = rocketLeftPointing
				}
			case 'd':
				if game.currentRocket == rocketRightPointing {
					game.rocketPos++
				} else {
					game.currentRocket = rocketRightPointing
				}
			case 'q':
				os.Exit(0)
			default:
				continue
			}
			game.drawRocket()
		}
	}()

	// Loop for moving the rocket every 150ms //
	go func() {
		for {
			if game.currentRocket != rocketNormal {
				game.rocketPos = uint(int(game.rocketPos) + game.currentRocket) // Update the rocket position dependent on which way the rocket is pointing
				game.drawRocket()
			}
			time.Sleep(150 * time.Millisecond)
		}
	}()

	// Main game loop //
	for {
		game.rocketFly()
		time.Sleep(1 * time.Second)
	}
}
