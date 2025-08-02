package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

/////// main.go //////////////////////////////////////////////////////////////
// Only the code for the rocket game, using abstractions from engine.go for //
// readability.                                                             //
// TODO: color for the actual game - not just the                           //
// title screen, powerups, make the game fair for defferent terminal sizes  //
//////////////////////////////////////////////////////////////////////////////

type gameState struct {
	rocketPos     uint
	playerScore   uint
	currentRocket int
	obstacleMap   [][]bool // The first list is the y the second is the x, if the bool is true then there is an obstacle
	terminal      *terminalManager
}

const (
	rocketWidth         = 13
	rocketLeftPointing  = -1
	rocketNormal        = 0
	rocketRightPointing = 1
)

func (game *gameState) doCollisionDetection() {
	var rocketBodyStartAndEnd []uint
	switch game.currentRocket {
	case rocketNormal:
		rocketBodyStartAndEnd = []uint{
			4, 10,
			5, 9,
			2, 12,
			2, 12,
			3, 11,
			4, 10,
			5, 9,
			5, 9,
			6, 8,
			7, 7,
		}
	case rocketLeftPointing:
		rocketBodyStartAndEnd = []uint{
			7, 13,
			6, 12,
			2, 11,
			2, 12,
			2, 12,
			2, 10,
			2, 8,
			1, 6,
			0, 5,
			0, 2,
		}
	case rocketRightPointing:
		rocketBodyStartAndEnd = []uint{
			1, 7,
			2, 8,
			3, 12,
			2, 12,
			2, 12,
			4, 12,
			6, 12,
			8, 13,
			9, 14,
			12, 14,
		}
	}
	obstacleMapY := int(game.playerScore)
	for i := 0; i < 20; i += 2 {
		obstacleMapY %= len(game.obstacleMap)
		if game.obstacleMap[obstacleMapY][game.rocketPos+rocketBodyStartAndEnd[i]] || game.obstacleMap[obstacleMapY][game.rocketPos+rocketBodyStartAndEnd[i+1]] {
			exitGame("Well done commander, you drove the rocket into an XXL rock!")
		}
		obstacleMapY++
	}
}

func (game *gameState) drawRocket() {
	maxRocketPos := uint(game.terminal.noOfColumns - rocketWidth)
	if game.rocketPos > maxRocketPos {
		game.rocketPos = maxRocketPos
	} else if game.rocketPos < 1 {
		game.rocketPos = 1
	}
	game.doCollisionDetection()
	switch game.currentRocket {
	case rocketNormal:
		drawMultilineText(
			[]string{
				`!!!!\=====/!!!!`,
				`!!   | S |   !!`,
				`!!___| C |___!!`,
				`!!\  | O |  /!!`,
				`!!!\ | R | /!!!`,
				`!!!!\| E |/!!!!`,
				fmt.Sprintf(`!!!!!|%3d|!!!!!`, game.playerScore),
				`!!!!!\   /!!!!!`,
				`!!!!!!\ /!!!!!!`,
				`!!!!!!!*!!!!!!!`,
			},
			game.rocketPos, 1)
	case rocketLeftPointing:
		drawMultilineText(
			[]string{
				`!!!!!!!|=====/ `,
				`!!     / S / !!`,
				`!!|-__/ C /  !!`,
				`!!|  / O /-__ !`,
				`!!| / R /  _/ !`,
				`!!|/ E / _/ !!!`,
				fmt.Sprintf(`! /%3d/_/ !!!!!`, game.playerScore),
				` /   / !!!!!!!!`,
				`/  _/ !!!!!!!!!`,
				`|_/ !!!!!!!!!!!`,
			},
			game.rocketPos, 1)
	case rocketRightPointing:
		drawMultilineText(
			[]string{
				` \=====|!!!!!!!`,
				`!! \ S \     !!`,
				`!!  \ C \__-|!!`,
				`! __-\ O \  |!!`,
				`! \_  \ R \ |!!`,
				`!!! \_ \ E \|!!`,
				fmt.Sprintf(`!!!!! \_\%3d\ !`, game.playerScore),
				`!!!!!!!! \   \ `,
				`!!!!!!!!! \_  \`,
				`!!!!!!!!!!! \_|`,
			},
			game.rocketPos, 1)
	}
}

func (game *gameState) drawRock(text []string, x uint) {
	y := uint(game.terminal.noOfLines - 3)
	for _, line := range text {
		pos := x
		fmt.Printf("\033[%d;%dH", y, x)
		for _, letter := range line {
			if letter != '!' {
				fmt.Print(string(letter))
				game.obstacleMap[game.convertYToObstacleMapIndex(y-1)][pos] = true
			} else {
				fmt.Print("\033[C")
			}
			pos++
		}
		y++
	}
}

func (game *gameState) convertYToObstacleMapIndex(y uint) uint {
	return (y + game.playerScore) % uint(len(game.obstacleMap))
}

func (game *gameState) rocketFly() {
	game.playerScore++
	fmt.Print("\033[9999;9999H\n")
	game.obstacleMap[game.convertYToObstacleMapIndex(uint(game.terminal.noOfLines-1))] = make([]bool, game.terminal.noOfColumns)
	if game.playerScore >= 999 {
		exitGame("Congratulations, you appear to have a player score longer then the rocket's counter can comprehend, this means that you've won the game.")
	}
	xPos := uint(rand.Intn(int(game.terminal.noOfColumns) - 12))
	switch rand.Intn(2) {
	case 0:
		game.drawRock([]string{
			"!_____",
			`/     \`,
			"|    _/",
			`!\__/`,
		}, xPos)
	case 1:
		game.drawRock([]string{
			"!!__!!!!___",
			`!/  \__/   \`,
			"|       ___/",
			`!\_____/`,
		}, xPos)
	case 2:
		game.drawRock([]string{
			`!___`,
			`/   \`,
			`\___/`,
		}, xPos)
	}
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
		game.obstacleMap = make([][]bool, game.terminal.noOfLines)
		for i := range game.obstacleMap {
			game.obstacleMap[i] = make([]bool, game.terminal.noOfColumns)
		}
		game.drawRocket()
	}
	game.terminal.resizeFunction()

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
				exitGame("You quit the game by pressing Q, bye bye.")
			default:
				continue
			}
			game.drawRocket()
		}
	}()

	// Main game loop //
	for {
		for i := 0; i < 10; i++ {
			if game.currentRocket != rocketNormal {
				game.rocketPos += uint(game.currentRocket) // Update the rocket position dependent on which way the rocket is pointing
				game.drawRocket()
			}
			time.Sleep(100 * time.Millisecond)
		}
		game.rocketFly()
		game.drawRocket()
	}
}
