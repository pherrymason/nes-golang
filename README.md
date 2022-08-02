[![CircleCI](https://circleci.com/gh/raulferras/nes-golang.svg?style=shield)](https://circleci.com/gh/raulferras/nes-golang)

# Status

- CPU: 100% of "legal" opcodes implemented
- PPU: Implemented scanline rendering. 
  - Renders contents of CHR ROM
  - Renders background.
  - Preliminar sprite rendering.
- APU: 0%
- MMU: 0%

## Visual evolution
<p align="center">
  <img src="assets/visual%20evolution/01-cpu/nestest-first-load.png" width="30%" alt="rendering decompilation"/>
  <img src="assets/visual%20evolution/01-cpu/nestest-improve-disasm.png" width="30%" alt="proper decompilation"/>
  <img src="assets/visual%20evolution/02-ppu-chr/supermariobros-chr-noise.png" width="30%" alt="first try rendering pattern table"/>
  <img src="assets/visual%20evolution/02-ppu-chr/supermariobros-chr.png" width="30%" alt="first try rendering pattern table"/>
  <img src="assets/visual%20evolution/02-ppu-chr/supermariobros-chr-2.png" width="30%" alt="first try rendering pattern table"/>
  <img src="assets/visual%20evolution/02-ppu-chr/supermariobros-chr-3.png" width="30%" alt="first try rendering pattern table"/>
  <img src="assets/visual%20evolution/02-ppu-chr/nestest-chr.png" width="30%" alt="first try rendering pattern table"/>
  <img src="assets/visual%20evolution/03-ppu-background/nestest-background-1.png" width="30%" alt="Renders background nestest"/>
  <img src="assets/visual%20evolution/03-ppu-background/nestest-background-2.png" width="30%" alt="Renders background nestest"/>
  <img src="assets/visual%20evolution/03-ppu-background/pacman-title-1.png" width="30%" alt="Renders background nestest"/>
  <img src="assets/visual%20evolution/03-ppu-background/supermariobros-title-1.png" width="30%" alt="Renders background super mario bros"/>
<img src="assets/visual%20evolution/03-ppu-background/donkey-kong-title-1.png" width="30%" alt="Renders background donkey kong, title screen"/>
  <img src="assets/visual%20evolution/03-ppu-background/donkey-kong-1.png" width="30%" alt="Renders background donkey kong, optimizations allow to see demo mode"/>
  <img src="assets/visual%20evolution/03-ppu-background/donkey-kong-title-2.png" width="30%" alt="Renders background donkey kong, title screen, small fixes in colors"/>
  <img src="assets/visual%20evolution/03-ppu-background/donkey-kong-2.png" width="30%" alt="Renders background donkey kong, fixes in colors"/>
  <img src="assets/visual%20evolution/03-ppu-background/donkey-kong-title-3.png" width="30%" alt="Renders donkey kong title, colors finally fixed"/>
  <img src="assets/visual%20evolution/03-ppu-background/donkey-kong-3.png" width="30%" alt="Renders donkey kong, colors finally fixed"/>
  <img src="assets/visual%20evolution/03-ppu-background/supermariobros-title-4.png" width="30%" alt="Super Mario Bros title screen, colors finally fixed by implementing transparent background colors"/>
  <img src="assets/visual%20evolution/05-ppu-sprite-rendering/donkey-kong-demo-1.png" width="30%" alt="Donkey Kong demo, preliminar sprite rendering"/>
</p>
