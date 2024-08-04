// Package traps designates the traps that can be encountered
// by a user.
package traps

const (
	// GETC gets a character from the keyboard.
	GETC = 0x20

	// OUT writes a characters to the console display.
	OUT = 0x21

	// PUTS writes a string of ASCII characters to the console display.
	PUTS = 0x22

	// IN prints a prompt on the screen and reads a single character
	// from the keyboard.
	IN = 0x23

	// PUTSP writes a string of ASCII characters to the console.
	PUTSP = 0x24

	// HALT halts execution and prints a message to the console.
	HALT = 0x25
)
