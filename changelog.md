2021-05-13:
Implement PPU triggering NMI.

2021-05-11:
Basics of reading and writing from CPU to the PPU implemented. 
This allows some roms to load palette information correctly.

2021-04-25: 
Started implementing the PPU.
CHR ROM is being rendered, although palette is for now hardcoded.

2021-04-24
Refactored and return back to use a single Go package for all code (circular dependencies are killing me)
Dropped Bus.
Created Memory interface with a CPUMemory implementation. This handles reading from the proper device
based on address ranges.
Loading roms is more encapsulated inside `GamePak`.

2021-04-22
First graphics! Start a window for the emulator and render some debug data, like program address and disassembled code being executed.

2021-04-19
Pass nestest up to first illegal opcode!