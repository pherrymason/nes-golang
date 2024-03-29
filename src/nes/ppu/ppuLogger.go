package ppu

import (
	"bufio"
	"fmt"
	"github.com/FMNSSun/hexit"
	"os"
	"strconv"
	"strings"
)

const PPU_LOG_BUFFER_MAXSIZE = 120000

type PPUState struct {
	ppuControl Control
	ppuStatus  Status
	ppuMask    Mask // Controls the rendering of sprites and backgrounds
	vRam       LoopyRegister
	tRam       LoopyRegister
	fineX      uint8
	readBuffer byte

	nextTileId    byte
	nextAttribute byte
	nextLowTile   byte
	nextHighTile  byte

	shifterTileLow       uint16
	shifterTileHigh      uint16
	shifterAttributeLow  uint16
	shifterAttributeHigh uint16

	cycle           uint32
	renderCycle     uint16
	currentScanline Scanline
	nmi             bool
}

func (s *PPUState) String() string {
	var msg strings.Builder
	msg.Grow(150)

	//msg.WriteString(fmt.Sprintf("C: %d ", s.cycle))
	msg.WriteString("rc:")
	msg.WriteString(strconv.FormatUint(uint64(s.renderCycle), 10))
	msg.WriteString(" sl:")
	msg.WriteString(strconv.FormatUint(uint64(s.currentScanline), 10))
	msg.WriteString(" CTRL: ")
	msg.WriteString(hexit.HexUint8Str(s.ppuControl.Value()))
	msg.WriteString(" STATUS: ")
	msg.WriteString(hexit.HexUint8Str(s.ppuStatus.Value()))

	msg.WriteString("  vRam: ")
	msg.WriteString(hexit.HexUint16Str(uint16(s.vRam.address())))
	msg.WriteString(" buffer: ")
	msg.WriteString(hexit.HexUint8Str(s.readBuffer))

	msg.WriteString("\n")

	return msg.String()
}

type logger2c02 struct {
	enabled    bool
	file       *os.File
	fileBuffer *bufio.Writer
	outputPath string
	snapshots  []PPUState
}

func NewLogger2c02(enabled bool, outputPath string) *logger2c02 {
	f, err := os.Create(outputPath)

	if err != nil {
		panic(fmt.Sprintf("Could not create log file: %s", outputPath))
	}

	logger := logger2c02{
		enabled:    enabled,
		file:       f,
		fileBuffer: bufio.NewWriterSize(f, PPU_LOG_BUFFER_MAXSIZE*10),
		outputPath: outputPath,
	}

	return &logger
}

func (logger *logger2c02) log(ppu *P2c02) {

	if len(logger.snapshots) == PPU_LOG_BUFFER_MAXSIZE {
		logger.logToFile()
		logger.snapshots = logger.snapshots[:0]
	}

	state := PPUState{
		ppuControl:           ppu.PpuControl,
		ppuStatus:            ppu.PpuStatus,
		ppuMask:              ppu.PpuMask,
		vRam:                 ppu.vRam,
		tRam:                 ppu.tRam,
		fineX:                ppu.fineX,
		readBuffer:           ppu.readBuffer,
		nextTileId:           ppu.bgNextTileId,
		nextAttribute:        ppu.bgNextAttribute,
		nextLowTile:          ppu.bgNextLowTile,
		nextHighTile:         ppu.bgNextHighTile,
		shifterTileLow:       ppu.bgShifterTileLow,
		shifterTileHigh:      ppu.bgShifterTileHigh,
		shifterAttributeLow:  ppu.bgShifterAttributeLow,
		shifterAttributeHigh: ppu.bgShifterAttributeHigh,
		cycle:                ppu.cycle,
		renderCycle:          ppu.renderCycle,
		currentScanline:      ppu.currentScanline,
		nmi:                  ppu.nmi,
	}

	logger.snapshots = append(logger.snapshots, state)
}

func (logger *logger2c02) logToFile() {
	for _, state := range logger.snapshots {
		logger.fileBuffer.WriteString(state.String())
	}
	logger.fileBuffer.Flush()
	logger.file.Sync()
}

func (logger *logger2c02) Close() {
	defer logger.file.Close()
	logger.logToFile()
	logger.file.Sync()
}
