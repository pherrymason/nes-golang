CPU 
===
The cpu is a MOS 6502 microprocesor.
Opcodes:
 - http://obelisk.me.uk/6502/reference.html
 - https://www.masswerk.at/6502/6502_instruction_set.html

There are some "ilegal"/unofficial opcodes that some game also use:
 -
 


Power-Up state
==============
On power up, the CPU sets the following state:
- Program Counter: 0x34
- IRQ Disabled
- A, X, Y = 0
- Stack Pointer: 0xFD
- Memory address 0x4017 = 0x00 (frame irq enabled)
- Memory address 0x4015 = 0x00 (all channels disabled)
- Memory addresses 0x4000 to 0x400F = 0x00
- Memory addresses 0x4010 to 0x4013 = 0x00
- All 15 bits of noise channel LFSR set to 0.
- APU Frame Counter reset
- Memory addresses 0x0000 to 0x07FF is set to an unreliable state.

Reference: http://wiki.nesdev.com/w/index.php/CPU_power_up_state

After reset state
=================
- A, X, Y were not affected
- S was decremented by 3 (but nothing was written to the stack)[3]
- The I (IRQ disable) flag was set to true (status ORed with $04)
- The internal memory was unchanged
- APU mode in $4017 was unchanged
- APU was silenced ($4015 = 0)
- APU triangle phase is reset to 0 (i.e. outputs a value of 15, the first step of its waveform)
- APU DPCM output ANDed with 1 (upper 6 bits cleared)
- 2A03G: APU Frame Counter reset. (but 2A03letterless: APU frame counter retains old value)