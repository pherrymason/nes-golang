xx August 2020
-----------------
Started the nes emulator project. I chose golang, after a previous try with kotlin, I expect golang to make it easier for me to work with lower level data structures.

After some initial NES documentation reading, looks like the CPU is a good candidate to be the first hardware component to emulate.
The NES uses a MOS 6502 CPU, its main characteristics are:
- 3 registers: A, X, Y.
- A Status Processor: 1 byte, each of its bits represent flags that affect instructions execution.
- A Program Counter: 2 bytes, points to the memory address the cpu must process.
- A stack of 255 bytes.
- 5x instructions.

It has other features that I choose to ignore for now. So I start developing the structure for emulating those registers and the instructions.

The very basic and simplified structure of a CPU emulator could be the following:

- Read memory.
- Interpret the value as an instruction.
- Repeat.

Before being able to implement the cpu instructions, I learn that an each instruction call is composed by: <Operation code> <Operand>
The `Operand` might be specified as a literal value, as an address pointing a memory location where the value is stored, or as an address plus a combination of some CPU registers.

These modes for obtaining the instruction `Operand` are categorized as `Addressing Modes` and our CPU implementation must implement them in order to execute the instructions with the proper input values.

As in order to get the final operand value might require to read multiple memory locations even before to execute the instruction, the `Program Counter` might be affected, and even take some cycles of CPU to complete. This has to be taken in mind in the implementation.
For my first approach, every address mode mutates the `Program Counter` of the cpu directly.

xx september 2020
-----------------
After a lot of research and multiple iterations of reading the documentation found, I finished a first implementation of every `Address Mode` together with some unit tests.
Those unit tests will help me avoid regressions.

This allows me to start coding the different instructions implementations.
Most of the instructions are simple and easy enough, but two of them highlight: ADC and SBC
They are "Add with Carry" and "Substract with Carry" and were harder to implement as they require you to understand how the CPU uses the flag Overflow. Thankfully there's a lot of articles explaining the maths of these operations and how the overflow flag is updated.


xx-september-2020
-----------------
Finished implementing all official 6502 cpu opcodes. 
There's still remaining opcodes (the so called unofficial opcodes), but I will delay their implementation until I find a rom that uses them.
As I have understand, they are usually combinations of official opcodes.


xx-september-2020
-----------------
I start preparing to be able to run a test rom and validate the cpu opcodes implementation. This test rom comes with an execution log so I can compare mine with it.
Turns out I still need to implement a bunch of architecture to be able to do so:

- Rom loading
- Memory mapping
- Initial CPU state


Rom loading:
Memory mapping:



After rom is properly loaded and mapped in the right memory address, all seemed ok to continue... but nope.

First AddressModes implementation did not incremented the program counter after each memory read, so this had to be fixed.
Next step is to generate an execution log so I can compare it with a known green log and validate my implementation.
In order to do this, I need to render the current cpu state (register values) and the next instruction to be executed plus the operand plus the evaluated operand (if required).

Because I wanted to reuse the addressMode implementations, and those mutated the state, I could not use them as is in order to render the log before running the next cpu step as they would modify the ProgramCounter. A different approach would be that AddressMode implementations don't mutate the state, but generate an effect that needs to be interpreted if wanted in order to mutate the state.
This way, my logging routine can read the current state without mutation.




21 september 2020
-----------------
Continued validating nestest.log and found some bugs:
- AddressMode Relative evaluation was wrong by using an old value of the ProgramCounter so final address was off by 1 position.
- The opcode table was wrong on some places so wrong implementations were being called in the wrong places.
Also refactored how addressModes are evaluated so its logic can be reused more easily from the logger.

Also, at address 0xC7DA there is a RTS that is jumping to a wrong address. I need to check if stack pop is correctly handled 
or if the wrong addresses were being pushed.

????
----
....


26 July 2022
------------
Some code cleanup on the PPU before attempting to render line-by-line timming. Made the logic for incrementing scanline and cycles counters more easy to follow and predictable. 
Implementing getting background color from palettes > 0 fixes some tiles rendering wrong colors (Super Mario Bros title screen now renders correctly).
Added some missing PPU tests.