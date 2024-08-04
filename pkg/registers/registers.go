// Package registers contains information relevant to
// registers.
// A register is a slot for storing a single value on the
// CPU and are essentially the workbench of the CPU. LC3
// has 10 total registers.
package registers

const (
	// RR0 is the 0-index general purpose register.
	RR0 = iota

	// RR1 is the 1-index general purpose register.
	RR1

	// RR2 is the 2-index general purpose register.
	RR2

	// RR3 is the 3-index general purpose register.
	RR3

	// RR4 is the 4-index general purpose register.
	RR4

	// RR5 is the 5-index general purpose register.
	RR5

	// RR6 is the 6-index general purpose register.
	RR6

	// RR7 is the 7-index general purpose register.
	RR7

	// RPC is the program counter register.
	RPC

	// RCOND is the condition flags register.
	RCOND

	// RCOUNT counts how many registers there are..
	RCOUNT
)

var registers [RCOUNT]uint16

// Get returns the registers that are available for
// the CPU to use.
func GetRegisters() *[RCOUNT]uint16 {
	return &registers
}

const (
	// MRKBSR is a memory mapped register used to interact with the
	// keyboard status.
	MRKBSR = 0xFE00

	// MRKBDR is a memory mapped register used to interact with the
	// keyboard data.
	MRKBDR = 0xFE02
)
