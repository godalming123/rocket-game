package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"
)

/////// engine.go /////////////////////////////////////////////
// All of the functionality to:                              //
//   - draw multiline text to any location in the terminal   //
//   - get the size of the terminal (with a resize callback) //
//   - read keypresses from the termianl                     //
// (TODO: only unix based systems are supported)             //
///////////////////////////////////////////////////////////////

// - draws a list of single line strings as one multiline string at position x and y
// - if any of the strings have an ! mark, then the original text will be left
func drawMultilineText(text []string, x uint, y uint) {
	for _, line := range text {
		fmt.Printf("\033[%d;%dH", y, x)
		for _, letter := range line {
			if letter == '!' {
				fmt.Print("\033[C")
			} else {
				fmt.Print(string(letter))
			}
		}
		y++
	}
}

// This struct contains things for managing the terminal such as:
//   - variables for how large the terminal is in size
//   - a reader to handle user keypresses
//   - a resize function which is ran whenever the terminal resizes
type terminalManager struct {
	// These first 4 items MUST be first
	noOfLines        uint16
	noOfColumns      uint16
	horizontalPixels uint16
	verticalPixels   uint16

	reader         *bufio.Reader
	resizeFunction func()
}

// Sets up:
//   - reading charecters from the terminal using the bufio reader in the returned struct
//   - getting the terminal size with the items in the returned struct
//   - handling the terminal size changing with the function that is:
//     A) accepted as an argument
//     B) is ran when this function is called
//     C) can be modified afterwoods with the returned struct's resizeFunction property
//     D) automatically clears the screen before it triggers the resize callback
func setupTerminal(terminalResizeFunction func()) *terminalManager {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run() // Disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()              // Do not display entered characters on the screen
	terminal := terminalManager{
		reader:         bufio.NewReader(os.Stdin),
		resizeFunction: terminalResizeFunction,
	}
	resizeChannel := make(chan os.Signal)
	signal.Notify(resizeChannel, syscall.SIGWINCH)
	go func() {
		for {
			fmt.Print("\033[2J")                                    // Clear the screen
			retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL, // Update the terminalSize
				uintptr(syscall.Stdin),
				uintptr(syscall.TIOCGWINSZ),
				uintptr(unsafe.Pointer(&terminal)))
			if int(retCode) == -1 {
				panic(errno)
			}
			terminal.resizeFunction()
			<-resizeChannel // Wait until the terminal resizes
		}
	}()
	return &terminal
}

// Clears the screen and prints a message in the top left before exiting the program
func exitGame(message string) {
	fmt.Printf("\033[2J\033[H%s\n", message)
	os.Exit(0)
}
