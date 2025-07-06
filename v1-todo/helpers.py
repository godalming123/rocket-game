import tty
import os
import time
import sys
import select
import copy
import termios

### GENERIC HELPERS ###

def changeMinAndMax(value, oldMin, oldMax, newMin, newMax):
    """Converts `value` from a number between `oldMin` and `oldMax` to a number between `newMin` and `newMax`"""
    between0and1 = (value - oldMin) / (oldMax - oldMin)
    return (between0and1 * (newMax - newMin)) + newMin

def fold(data, initialState, func):
    state = initialState
    for i, elem in enumerate(data):
        state = func(state, elem)
    return state

def replaceSubstringAt(string: str, pos: int, newText: str) -> str:
	return string[:pos] + newText + string[pos + len(newText):]

def bitFieldSetBit(bit_field: int, bit: int, value: bool):
    if value == True:
        bit_field |= bit
    else:
        bit_field &= ~bit

### CONSOLE OUTPUT HELPERS ###

def clear():
    print("\033[2J", end="")

def cursorToTopLeft():
    print("\033[H", end="")

def typeLineOfText(text, minWait, maxWait, end="\n"):
    waitDelta = maxWait - minWait
    for i, letter in enumerate(text):
        extraWait = changeMinAndMax(i, 0, len(text), -1, 1)**2.0
        time.sleep((extraWait * waitDelta) + minWait)
        sys.stdout.write(letter)
        sys.stdout.flush()
    print(end, end="")

def typeText(text, start="\n", end="\n", minWait=0.03, maxWait=0.08):
    print(start, end="")
    def foldFunc(state, data):
        if len(state) == 0 or len(state[len(state)-1]) + len(data) > 80:
             return state + [data]
        else:
            return state[:len(state)-1] + [state[len(state)-1] + " " + data]
    lines = fold(text.split(), [], foldFunc)
    for line in lines[:len(lines)-1]:
        typeLineOfText(line, minWait, maxWait)
    typeLineOfText(lines[len(lines)-1], minWait, maxWait, end)
    sys.stdout.flush()

### USER INPUT HELPERS ###

def setInputSettings(fd, settingsBase, *, immediatelyReadableKeypresses: bool, showUserInput: bool, allowCtrlC: bool):
    # See https://www.man7.org/linux/man-pages/man3/termios.3.html
    newSettings = copy.deepcopy(settingsBase)
    # bitFieldSetBit(newSettings[3], termios.ICANON, not immediatelyReadableKeypresses)
    # bitFieldSetBit(newSettings[3], termios.ECHO, showUserInput)
    # bitFieldSetBit(newSettings[3], termios.IGNBRK, not allowCtrlC)
    termios.tcsetattr(fd, termios.TCSAFLUSH, newSettings)
    if immediatelyReadableKeypresses:
        tty.setcbreak(fd)

def pressToContinue(prompt="\n↩ "):
    input(prompt)
    print()

def getString(prompt, start="\n"):
    typeText(prompt, start)
    return input(" => ")

def tryReadStdinChar() -> str:
    """Checks if something exists in stdin. If so, read the first charecter that is in stdin, and return it. Otherwise, return a blank string."""
    if select.select([sys.stdin], [], [], 0.0)[0]:
        return sys.stdin.read(1)
    return ""

### GAME SCREEN ABSTRACTION ###

class gameScreen:
    def __init__(self):
        self.width = os.get_terminal_size().columns
        self.height = os.get_terminal_size().lines - 1
        self.contents = " " * (self.width * self.height)
        self.collisionString = self.contents

    def drawAscii(self, x: int, y: int, text: list[str], collisionText: list[str]):
        startingIndex = (self.width * max(y, 0)) + x
        for index, (line, collisionLine) in enumerate(zip(text, collisionText)):
            if y + index < 0:
                continue
            if y + index >= self.height:
                break
            self.contents = replaceSubstringAt(self.contents, startingIndex, line)
            self.collisionString = replaceSubstringAt(self.collisionString, startingIndex, collisionLine)
            startingIndex += self.width

    def moveAscii(self, x1, y1, x2, y2, width, height):
        # cut the text
        fillerChars = " " * width
        startingIndex = (self.width * (y1)) + x1
        text = []
        collisionText = []
        for lineOn in range(height):
            if 0 <= (y1 + lineOn) < self.height:  # if text fits in screen vertical height

                # Copy old text
                text += self.contents[startingIndex:startingIndex + width] + "\n"
                collisionText += self.collisionString[startingIndex:startingIndex + width]

                # Clear old text
                self.contents = replaceSubstringAt(self.contents, startingIndex, fillerChars)
                self.collisionString = replaceSubstringAt(self.collisionString, startingIndex, fillerChars)
            startingIndex += self.width  # increment the starting index

        # paste the text
        self.drawAscii(x2, y2, text, collisionText)

    def drawRect(self, x1, y1, cols, lines):
        startingIndex = (self.width * y1) + x1
        fillerChars = " " * cols
        for line in range(lines):
            self.contents = replaceSubstringAt(self.contents, startingIndex, fillerChars)
            self.collisionString = replaceSubstringAt(self.collisionString, startingIndex, fillerChars)
            startingIndex += self.width

    def moveContentsDown(self, deltaY, newRows):
        self.contents = newRows + self.contents[0:-(self.width * deltaY)]
        self.collisionString = newRows + self.collisionString[:-(self.width * deltaY)]

    def print(self):
        print(self.contents, end="")
