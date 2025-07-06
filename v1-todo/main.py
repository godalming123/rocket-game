#!/user/bin/env python3

import time
import os
import sys
import tty
import random
import math
import termios
import helpers
import art

### INITIALISATION ###

fd = sys.stdin.fileno()
original_settings = termios.tcgetattr(fd)
seed = random.randint(1, 10)
random.seed(seed)

def exit(message: str):
    helpers.clear()
    if message != "":
        helpers.typeText(message, start="")
    termios.tcsetattr(fd, termios.TCSAFLUSH, original_settings)
    sys.exit(0)

### TITLE SCREEN ###

helpers.setInputSettings(
    fd,
    original_settings,
    immediatelyReadableKeypresses = True,
    showUserInput = False,
    allowCtrlC = False,
)
playIsSelected = True
while True:
    helpers.clear()
    print("                     Seed is:", seed)
    print(art.title)
    print("                   Press q to quit")
    print(art.playSelectedText if playIsSelected else art.aboutSelectedText)
    key = sys.stdin.read(1)[0]
    if key == "q":
        exit("")
    elif (key == "A" and not playIsSelected) or (key == "B" and playIsSelected):
        playIsSelected = not playIsSelected
    elif key == "\n":
        if playIsSelected:
            helpers.typeText("The story begins in a peaceful village with only a few inhabita...")
            helpers.pressToContinue()
            helpers.typeText("Oh lets just get on with it! We know what you're here for, why ofcourse, the rocket game!")
            helpers.pressToContinue()
            helpers.typeText("Use W, A , S and D to move, pick up powerups, and complete the game when you get to a score of 500. If you want to give up, press G.")
            break
        else:
            print(art.aboutText)
            helpers.pressToContinue()
helpers.setInputSettings(
    fd,
    original_settings,
    immediatelyReadableKeypresses = False,
    showUserInput = True,
    allowCtrlC = False,
)
playerName = helpers.getString("Please enter your name")

### INITIALISE THE GAME ###

helpers.setInputSettings(
    fd,
    original_settings,
    immediatelyReadableKeypresses = True,
    showUserInput = False,
    allowCtrlC = False,
)

screen = helpers.gameScreen()
playerScore = 1
rocketXPos = round((screen.width / 2) - 7.5)
rocketMode = "normal"
health = 1
totalHealthPacks = 0
showObjs = False

def getSecondsSinceEpoch() -> int:
    return time.time()

def playerScoreIndicator(playerScore: int) -> list[str]:
    stringScore = str(playerScore)
    dashes = "-" * len(stringScore)
    return [
        "/-" + dashes + "-\ ",
        "| " + stringScore + " |",
        "\\-" + dashes + "-/",
    ]

def playerHealthIndicator(playerHealth: int) -> list[str]:
    healthIndicatorText = []
    for notHealthOn in range(3 - playerHealth):
        healthIndicatorText += 4 * ["               "]
    for healthOn in range(playerHealth):
        healthIndicatorText += [
            r"               ",
            r"/-------------\ ",        
            r"| + First Aid |",        
            r"\-------------/",
        ]
    return healthIndicatorText

def handleCollisionDetection():
    # TODO: implement collision detection
    pass

### START THE GAME LOOP ###

firstSSE = SSE = getSecondsSinceEpoch()
redraw = True
while True:
    if redraw:
        helpers.cursorToTopLeft()
        handleCollisionDetection()

        # draw the rocket, player score indicator, and player health indicator
        screen.drawAscii(rocketXPos, screen.height - 10, art.rocketAsciis[rocketMode][0], art.rocketAsciis[rocketMode][1])
        screen.drawAscii(2, screen.height - 3, playerScoreIndicator(playerScore), [""] * 3)
        screen.drawAscii(screen.width - 17, screen.height - 12, playerHealthIndicator(health), [""] * 12)

        # print the screen
        if showObjs:
            print(screen.collisionString)
        else:
            screen.print()
    redraw = True

    # get user keypress
    c = helpers.tryReadStdinChar()
    if c == "a" and rocketXPos > 0:
        rocketXPos -= 1
    elif c == "d" and rocketXPos < screen.width - 15:
        rocketXPos += 1
    elif c == "o":
        showObjs = not showObjs
    elif c == "g":
        exit("Oh well, you gave up.")
    else:
        redraw = False

    # sleep and calculate SSE
    SSE = getSecondsSinceEpoch()

    # move the contents of the screen down
    charsToMove = math.trunc(SSE - firstSSE) - playerScore + 1
    for _ in range(charsToMove):
        redraw = True
        # increment player score
        playerScore += 1

        # move screen down
        screen.moveContentsDown(1, " " * screen.width)

        # draw new rocks
        rockX = random.randint(5, screen.width)
        rockY = random.randint(0, 10)
        rock = random.randint(0, 2)
        screen.drawAscii(rockX, rockY, art.rocks[rock * 2], art.rocks[(rock * 2) + 1])

        # draw powerups
        if random.randint(0, 1) == 1:
            powerupX = random.randint(5, screen.width)
            powerupY = random.randint(0, 10)
            if random.randint(0, 4):
                screen.drawAscii(powerupX, powerupY, art.powerups["health"][0], art.powerups["health"][1])
            else:
                screen.drawAscii(powerupX, powerupY, art.powerups["invinsible"][0], art.powerups["invinsible"][1])

        # do collision detection
        handleCollisionDetection()

    if playerScore > 500:
        exit("Congratulations, on completing the game and getting to the moon. You did many hard things and hot to a total score of 500 you also managed to get " + str(totalHealthPacks) + " healthpacks!")
