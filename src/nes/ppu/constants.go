package ppu

// Registers

// PPUCTRL NMI enable (V), PPU master/slave (P), sprite height (H),
// background tile select (B), sprite tile select (S), increment mode (I),
// name table select (NN)
const PPUCTRL = 0x2000

// PPUMASK color emphasis (BGR), sprite enable (s), background enable (b),
const PPUMASK = 0x2001

// PPSTATUS sprite left column enable (M), background left column enable (m), greyscale (G)
// vblank (V), sprite 0 hit (S), sprite overflow (O); read resets write pair for $2005/$2006
const PPUSTATUS = 0x2002
const OAMADDR = 0x2003
const OAMDATA = 0x2004
const PPUSCROLL = 0x2005
const PPUADDR = 0x2006
const PPUDATA = 0x2007
const OAMDMA = 0x4014

const NES_PALETTE_COLORS = 64

// Memory sizes
const OAMDATA_SIZE = 256
const NAMETABLE_SIZE = 1024
const PALETTE_SIZE = 32

// Screen constants

const PPU_SCREEN_SPACE_CYCLES_BY_SCANLINE = 256
const PPU_CYCLES_BY_SCANLINE = 341
const PPU_SCREEN_SPACE_SCANLINES = 240
const VBLANK_START_SCANLINE = 241
const VBLANK_END_SCNALINE = 261
const PPU_SCANLINES = 261
const PPU_VBLANK_START_CYCLE = (PPU_SCREEN_SPACE_SCANLINES + 1) * PPU_CYCLES_BY_SCANLINE
const PPU_VBLANK_END_CYCLE = PPU_SCANLINES * PPU_CYCLES_BY_SCANLINE

const TILE_WIDTH = 8
const TILE_HEIGHT = 8
const TILE_PIXELS = 8 * 8
