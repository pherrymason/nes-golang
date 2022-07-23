[![CircleCI](https://circleci.com/gh/raulferras/nes-golang.svg?style=shield)](https://circleci.com/gh/raulferras/nes-golang)

# Status

- CPU: 100% of "legal" opcodes implemented
- PPU: Implementing direct rendering. 
  - Renders contents of CHR ROM
  - Renders background.
- APU: 0%
- MMU: 0%

## Visual evolution
<p align="center">
  <img src="var/visual%20evolution/01-cpu/nestest-first-load.png" width="30%" alt="rendering decompilation"/>
  <img src="var/visual%20evolution/01-cpu/nestest-improve-disasm.png" width="30%" alt="proper decompilation"/>
  <img src="var/visual%20evolution/02-ppu-chr/supermariobros-chr-noise.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/02-ppu-chr/supermariobros-chr.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/02-ppu-chr/supermariobros-chr-2.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/02-ppu-chr/supermariobros-chr-3.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/02-ppu-chr/nestest-chr.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/03-ppu-background/nestest-background-1.png" width="30%" alt="Renders background nestest"/>
  <img src="var/visual%20evolution/03-ppu-background/nestest-background-2.png" width="30%" alt="Renders background nestest"/>
  <img src="var/visual%20evolution/03-ppu-background/pacman.png" width="30%" alt="Renders background nestest"/>
  <img src="var/visual%20evolution/03-ppu-background/donkey-kong.png" width="30%" alt="Renders background nestest"/>
  <img src="var/visual%20evolution/03-ppu-background/supermariobros.png" width="30%" alt="Renders background super mario bros"/>
  <img src="var/visual%20evolution/03-ppu-background/donkey-kong-2.png" width="30%" alt="Renders background donkey kong, optimizations allow to see demo mode"/>
  <img src="var/visual%20evolution/03-ppu-background/donkey-kong-3.png" width="30%" alt="Renders background donkey kong, small fixes in colors"/>
  <img src="var/visual%20evolution/03-ppu-background/donkey-kong-4.png" width="30%" alt="Renders background donkey kong, small fixes in colors"/>
</p>
