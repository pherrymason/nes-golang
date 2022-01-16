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
  <img src="var/visual%20evolution/Screenshot%20from%202021-04-22%2019-02-13.png" width="30%" alt="rendering decompilation"/>
  <img src="var/visual%20evolution/Screenshot%20from%202021-04-22%2019-32-13.png" width="30%" alt="proper decompilation"/>
  <img src="var/visual%20evolution/Screenshot%20from%202021-04-24%2020-18-57.png" width="30%" alt="first try rendering pattern table"/>
  <img src="var/visual%20evolution/Screenshot%20from%202021-04-25%2000-19-39.png" width="30%" alt="Renders pattern table"/>
  <img src="var/visual%20evolution/Screenshot%20from%202021-04-25%2000-36-54.png" width="30%" alt="testing palette selection"/>
  <img src="var/visual%20evolution/Screenshot%20from%202022-01-16%2017-52-50.png" width="30%" alt="Renders background nestest"/>
</p>
