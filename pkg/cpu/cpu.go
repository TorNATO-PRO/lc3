// Package cpu implements the CPU for executing
// the LC3 virtual machine. I also plan to showcase
// a mix of functional and object oriented programming
// for your viewing pleasure.
package cpu

import (
	"bufio"
	"fmt"
	"lc3/internal/constants"
	"lc3/pkg/cflags"
	"lc3/pkg/opcodes"
	"lc3/pkg/registers"
	"lc3/pkg/traps"
	"os"
)

// opTable specifies a table of operations and corresponding functions.
var opTable = map[uint16]func(cpu *cpu, cancel func()) error{
	opcodes.OPADD:  handleAdd,
	opcodes.OPBR:   handleBr,
	opcodes.OPLD:   handleLoad,
	opcodes.OPST:   handleStore,
	opcodes.OPJSR:  handleJumpSubroutine,
	opcodes.OPAND:  handleAnd,
	opcodes.OPLDR:  handleLoadR,
	opcodes.OPSTR:  handleStr,
	opcodes.OPRTI:  unhandledOpcode,
	opcodes.OPNOT:  handleNot,
	opcodes.OPLDI:  handleLoadIndirect,
	opcodes.OPSTI:  handleStoreIndirect,
	opcodes.OPJMP:  handleJmp,
	opcodes.OPRES:  unhandledOpcode,
	opcodes.OPLEA:  handleLoadEffectiveAddress,
	opcodes.OPTRAP: handleTrap,
}

var trapTable = map[uint16]func(cpu *cpu, cancel func()) error{
	traps.GETC:  handleGetC,
	traps.OUT:   handleOut,
	traps.PUTS:  handlePuts,
	traps.IN:    handleIn,
	traps.PUTSP: handlePutsP,
	traps.HALT:  handleHalt,
}

// CPU defines an interface that we expect for a
// LC3 CPU implementation. Given an initial memory state,
// we should be able to run the program!.
type CPU interface {
	// Run runs the CPU given an initial memory state.
	Run(memory [constants.MemoryMax]uint16) error
}

// cpu defines our default CPU implementation.
type cpu struct {
	// memory is the current place in memory
	// that we are at.
	memory [constants.MemoryMax]uint16

	// registers denotes the current workbench state
	// of the CPU.
	registers [registers.RCOUNT]uint16

	// op represents the current opcode that the CPU
	// is executing.
	op uint16

	// instr represents the current instruction
	// executing on the CPU.
	instr uint16
}

// NewCPU defines a new CPU.
func NewCPU() *cpu {
	var regs [registers.RCOUNT]uint16

	cpu := cpu{
		registers: regs,
	}

	cpu.registers[registers.RCOND] = cflags.FLZRO

	// set the PC to starting position,
	// in the case of LC3, 0x3000 is the default starting
	// position for whatever reason.
	cpu.registers[registers.RPC] = 0x3000

	return &cpu
}

// Run runs the CPU over the memory.
func (c *cpu) Run(memory [constants.MemoryMax]uint16) error {
	c.memory = memory

	err := c.Loop(func(op uint16, cancel func()) error {
		fn, ok := opTable[op]

		if !ok {
			return fmt.Errorf("unrecognized operation %d", op)
		}

		err := fn(c, cancel)

		return err
	})

	return err
}

// Loop takes in a continuation for the function
// that could potentially return an error, and executes
// it, breaking on either the nil or a call to the cancel
// function.
func (c *cpu) Loop(loopCont func(op uint16, cancel func()) error) error {
	running := true

	cancel := func() {
		running = false
	}

	exec := 0

	for running {
		err := c.Step()

		// fmt.Println(c.instr)

		if err != nil {
			return err
		}

		err = loopCont(c.op, cancel)

		if err != nil {
			return err
		}

		exec++
	}

	return nil
}

// Step steps the CPU along.
func (c *cpu) Step() error {
	// read the memory location of the program counter.
	instr, err := c.memoryRead(c.registers[registers.RPC])
	if err != nil {
		return err
	}

	// increment the program counter.
	c.incrProgramCounter()

	c.op = instr >> 12

	c.instr = instr

	return nil
}

// incrProgramCounter increments the program counter.
func (c *cpu) incrProgramCounter() uint16 {
	pc := c.registers[registers.RPC]
	pc++
	c.registers[registers.RPC] = pc
	return pc
}

// memoryRead reads a value from the current memory address.
func (c *cpu) memoryRead(address uint16) (uint16, error) {
	if address == registers.MRKBSR {
		reader := bufio.NewReader(os.Stdin)

		key, err := reader.ReadByte()

		if err != nil {
			return 0, err
		}

		if uint16(key) != 0 {
			c.memory[registers.MRKBSR] = 1 << 15
			c.memory[registers.MRKBDR] = uint16(key)
		} else {
			c.memory[registers.MRKBSR] = 0
		}

	}

	return c.memory[address], nil
}

// unable to write to a memory address.
func (c *cpu) memoryWrite(address uint16, val uint16) error {
	c.memory[address] = val

	return nil
}

// updateFlags updates the flags of a given register.
func (c *cpu) updateFlags(r uint16) {
	if c.registers[r] == 0 {
		c.registers[registers.RCOND] = cflags.FLZRO
	} else if c.registers[r]>>15 != 0 {
		c.registers[registers.RCOND] = cflags.FLNEG
	} else {
		c.registers[registers.RCOND] = cflags.FLPOS
	}
}

// unhandledOpcode specifies that an opcode has yet to
// be handled.
func unhandledOpcode(cpu *cpu, cancel func()) error {
	return fmt.Errorf("failed to handle opcode %x", cpu.op)
}

// handleAdd handles the add opcode.
func handleAdd(cpu *cpu, cancel func()) error {
	r0 := (cpu.instr >> 9) & 0x7
	r1 := (cpu.instr >> 6) & 0x7
	immFlag := (cpu.instr >> 5) & 0x1

	if immFlag == 1 {
		imm5 := signExtend(cpu.instr&0x1F, 5)
		cpu.registers[r0] = cpu.registers[r1] + imm5
	} else {
		r2 := cpu.instr & 0x7
		cpu.registers[r0] = cpu.registers[r1] + cpu.registers[r2]
	}

	cpu.updateFlags(r0)

	return nil
}

// handleAnd handles the and opcode.
func handleAnd(cpu *cpu, cancel func()) error {
	// destination register
	r0 := (cpu.instr >> 9) & 0x7

	// first operand
	r1 := (cpu.instr >> 6) & 0x7

	// imm flag
	immFlag := (cpu.instr >> 5) & 0x1

	if immFlag == 1 {
		imm5 := signExtend(cpu.instr&0x1F, 5)
		cpu.registers[r0] = cpu.registers[r1] & imm5
	} else {
		r2 := cpu.instr & 0x7
		cpu.registers[r0] = cpu.registers[r1] & cpu.registers[r2]
	}

	cpu.updateFlags(r0)

	return nil
}

// handleBr handles the conditional branch opcode.
func handleBr(cpu *cpu, cancel func()) error {
	condFlag := (cpu.instr >> 9) & 0x7
	pcOffset := signExtend(cpu.instr&0x1FF, 9)

	if (condFlag & cpu.registers[registers.RCOND]) != 0 {
		cpu.registers[registers.RPC] += pcOffset
	}

	return nil
}

// handleJmp handles the jump and ret opcodes.
func handleJmp(cpu *cpu, cancel func()) error {
	r1 := (cpu.instr >> 6) & 0x7
	cpu.registers[registers.RPC] = cpu.registers[r1]

	return nil
}

// handleJsr handles the jump to subroutine opcode.
func handleJumpSubroutine(cpu *cpu, cancel func()) error {
	cpu.registers[registers.RR7] = cpu.registers[registers.RPC]

	bit11 := (cpu.instr >> 11) & 0x1

	if bit11 == 0 {
		baseR := (cpu.instr >> 6) & 0x7
		cpu.registers[registers.RPC] = cpu.registers[baseR]
	} else {
		pcOffset := signExtend(cpu.instr&0x7FF, 11)
		cpu.registers[registers.RPC] += pcOffset
	}

	return nil
}

// handleLoad handles the load opcode.
func handleLoad(cpu *cpu, cancel func()) error {
	dr := (cpu.instr >> 9) & 0x7
	pcOffset := signExtend(cpu.instr&0x1FF, 9)

	data, err := cpu.memoryRead(cpu.registers[registers.RPC] + pcOffset)

	if err != nil {
		return err
	}

	cpu.registers[dr] = data
	cpu.updateFlags(dr)

	return nil
}

// handleLoadR handles the load base + offset opcode.
func handleLoadR(cpu *cpu, cancel func()) error {
	dr := (cpu.instr >> 9) & 0x7
	br := (cpu.instr >> 6) & 0x7
	offset := signExtend(cpu.instr&0x3F, 6)
	k, err := cpu.memoryRead(cpu.registers[br] + offset)

	if err != nil {
		return err
	}

	cpu.registers[dr] = k
	cpu.updateFlags(dr)
	return nil
}

// handleStore handles the store operation.
func handleStore(cpu *cpu, cancel func()) error {
	sr := (cpu.instr >> 9) & 0x7
	pcOffset := signExtend(cpu.instr&0x1FF, 9)
	loc := cpu.registers[registers.RPC] + pcOffset

	return cpu.memoryWrite(loc, cpu.registers[sr])
}

// handleStoreIndirect handles store indirect.
func handleStoreIndirect(cpu *cpu, cancel func()) error {
	pc := cpu.registers[registers.RPC]
	pcOffset := signExtend(cpu.instr&0x1FF, 9)
	addr, err := cpu.memoryRead(pc + pcOffset)
	if err != nil {
		return err
	}

	sr := (cpu.instr >> 9) & 0x7
	return cpu.memoryWrite(addr, cpu.registers[sr])
}

// handleStr handles the store base + offset operation.
func handleStr(cpu *cpu, cancel func()) error {
	sr := (cpu.instr >> 9) & 0x7
	baseR := (cpu.instr >> 6) & 0x7
	offset := signExtend(cpu.instr&0x3F, 6)
	return cpu.memoryWrite(cpu.registers[baseR]+offset, cpu.registers[sr])
}

// handleLoadEffectiveAddress handles loading the effective address.
func handleLoadEffectiveAddress(cpu *cpu, cancel func()) error {
	dr := (cpu.instr >> 9) & 0x7
	pcOffset := signExtend(cpu.instr&0x1FF, 9)
	cpu.registers[dr] = cpu.registers[registers.RPC] + pcOffset
	cpu.updateFlags(dr)
	return nil
}

// handleNot handles the not address.
func handleNot(cpu *cpu, cancel func()) error {
	dr := (cpu.instr >> 9) & 0x7
	sr := (cpu.instr >> 6) & 0x7
	cpu.registers[dr] = ^cpu.registers[sr]
	cpu.updateFlags(dr)
	return nil
}

// handleLoadIndirect handles indirectly loading stuff
// from the CPU.
func handleLoadIndirect(cpu *cpu, cancel func()) error {
	r0 := (cpu.instr >> 9) & 0x7

	pcOffset := signExtend(cpu.instr&0x1FF, 9)

	addr, err := cpu.memoryRead(cpu.registers[registers.RPC] + pcOffset)

	if err != nil {
		return err
	}

	val, err := cpu.memoryRead(addr)

	if err != nil {
		return err
	}

	cpu.registers[r0] = val

	cpu.updateFlags(r0)

	return nil
}

// handleTrap handles the trap opcode.
func handleTrap(cpu *cpu, cancel func()) error {
	cpu.registers[registers.RR7] = cpu.registers[registers.RPC]

	trap := cpu.instr & 0xFF

	handler, ok := trapTable[trap]

	if !ok {
		return fmt.Errorf("unrecognized trap %x", trap)
	}

	return handler(cpu, cancel)
}

// handleGetC handles the GetC trap.
func handleGetC(cpu *cpu, cancel func()) error {
	reader := bufio.NewReader(os.Stdin)

	byt, err := reader.ReadByte()
	if err != nil {
		return err
	}

	cpu.registers[registers.RR0] = uint16(byt)
	cpu.updateFlags(registers.RR0)

	return nil
}

// handlePut handles the Puts trap.
func handlePuts(cpu *cpu, cancel func()) error {
	for addr := cpu.registers[registers.RR0]; ; addr++ {
		char, err := cpu.memoryRead(addr)

		if err != nil {
			return err
		}

		if char == 0 {
			break
		}

		fmt.Printf("%c", char)
	}

	return nil
}

// handleOut handles the Out trap.
func handleOut(cpu *cpu, cancel func()) error {
	writer := bufio.NewWriter(os.Stdout)

	elem := byte(cpu.registers[registers.RR0])

	if err := writer.WriteByte(elem); err != nil {
		return err
	}

	return writer.Flush()
}

// handleIn handles the In trap.
func handleIn(cpu *cpu, cancel func()) error {
	fmt.Print("Enter a character: ")

	reader := bufio.NewReader(os.Stdin)

	writer := bufio.NewWriter(os.Stdin)

	byt, err := reader.ReadByte()

	if err != nil {
		return err
	}

	err = writer.WriteByte(byt)

	if err != nil {
		return err
	}

	cpu.registers[registers.RR0] = uint16(byt)

	cpu.updateFlags(registers.RR0)

	return writer.Flush()
}

// handlePutsP handles the PutsP trap.
func handlePutsP(cpu *cpu, cancel func()) error {
	writer := bufio.NewWriter(os.Stdout)

	for addr := cpu.registers[registers.RR0]; ; addr++ {
		char, err := cpu.memoryRead(addr)

		if err != nil {
			return err
		}

		if char == 0 {
			break
		}

		err = writer.WriteByte(byte(char & 0xFF))

		if err != nil {
			return err
		}

		symb := char >> 8

		if symb != 0 {
			if err := writer.WriteByte(byte(symb)); err != nil {
				return err
			}
		}
	}

	return writer.Flush()
}

// handleHalt handles the Halt trap.
func handleHalt(cpu *cpu, cancel func()) error {
	// this is why the cancel function is getting
	// passed around everywhere. See what I did
	// there???
	cancel()

	return nil
}

// signExtend extends the sign of an unsigned int16
// by of bitCount bits.
func signExtend(x, bitCount uint16) uint16 {
	if (x>>(bitCount-1))&1 != 0 {
		x |= 0xFFFF << bitCount
	}
	return x
}
