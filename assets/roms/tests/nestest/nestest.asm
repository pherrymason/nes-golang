;-------------------------------------------------------------------------------
; nestest.nes disasembled by DISASM6 v1.5
;-------------------------------------------------------------------------------

;-------------------------------------------------------------------------------
; Registers
;-------------------------------------------------------------------------------
PPUCTRL              = $2000
PPUMASK              = $2001
PPUSTATUS            = $2002
OAMADDR              = $2003
OAMDATA              = $2004
PPUSCROLL            = $2005
PPUADDR              = $2006
PPUDATA              = $2007
SQ1_VOL              = $4000
SQ1_SWEEP            = $4001
SQ1_LO               = $4002
SQ1_HI               = $4003
SQ2_VOL              = $4004
SQ2_SWEEP            = $4005
SQ2_LO               = $4006
SQ2_HI               = $4007
TRI_LINEAR           = $4008
TRI_LO               = $400A
TRI_HI               = $400B
NOISE_VOL            = $400C
NOISE_LO             = $400E
NOISE_HI             = $400F
DMC_FREQ             = $4010
DMC_RAW              = $4011
DMC_START            = $4012
DMC_LEN              = $4013
OAM_DMA              = $4014
SND_CHN              = $4015
JOY1                 = $4016
JOY2                 = $4017

;-------------------------------------------------------------------------------
; iNES Header
;-------------------------------------------------------------------------------
            .db "NES", $1A     ; Header
            .db 1              ; 1 x 16k PRG banks
            .db 1              ; 1 x 8k CHR banks
            .db %00000000      ; Mirroring: Horizontal
                               ; SRAM: Not used
                               ; 512k Trainer: Not used
                               ; 4 Screen VRAM: Not used
                               ; Mapper: 0
            .db %00000000      ; RomType: NES
            .hex 00 00 00 00   ; iNES Tail 
            .hex 00 00 00 00    

;-------------------------------------------------------------------------------
; Program Origin
;-------------------------------------------------------------------------------
            .org $c000         ; Set program counter

;-------------------------------------------------------------------------------
; ROM Start
;-------------------------------------------------------------------------------
            jmp __start__c5f5         ; $c000: 4c f5 c5  

;-------------------------------------------------------------------------------
            rts                ; $c003: 60        

;-------------------------------------------------------------------------------
; reset vector
;-------------------------------------------------------------------------------
reset:      sei                ; $c004: 78              Set Interrupt Flag
            cld                ; $c005: d8              Clear Decimal Mode
            ldx #$ff           ; $c006: a2 ff           cpu.X = 0xFF
            txs                ; $c008: 9a              cpu.push(cpu.X)

wait1stVBlank:    lda PPUSTATUS               ; $c009: ad 02 20        cpu.A = ppu.Status
                  bpl wait1stVBlank           ; $c00c: 10 fb           If vertical blank started, go next

wait2ndVBlank:      lda PPUSTATUS             ; $c00e: ad 02 20        cpu.A = ppu.Status
                    bpl wait2ndVBlank         ; $c011: 10 fb           If vertical blank started, go next    

            lda #$00           ; $c013: a9 00           cpu.A = 0x00
            sta PPUCTRL        ; $c015: 8d 00 20        ppu.CTRL = cpu.A
            sta PPUMASK        ; $c018: 8d 01 20        ppu.MASK = cpu.A
            sta PPUSCROLL      ; $c01b: 8d 05 20        ppu.SCROLL = cpu.A
            sta PPUSCROLL      ; $c01e: 8d 05 20        ppu.SCROLL = cpu.A
            lda PPUSTATUS      ; $c021: ad 02 20        cpu.A = ppu.STATUS
            ldx #$20           ; $c024: a2 20           .
            stx PPUADDR        ; $c026: 8e 06 20        .
            ldx #$00           ; $c029: a2 00           .
            stx PPUADDR        ; $c02b: 8e 06 20        ppu.ADDR = 0x2000               // Let's write into NameTable
            ldx #$00           ; $c02e: a2 00           cpu.X = 0x00
            ldy #$0f           ; $c030: a0 0f           cpu.Y = 0x0F
            lda #$00           ; $c032: a9 00           cpu.A = 0x00

__c034:     sta PPUDATA        ; $c034: 8d 07 20        for y = 0; y < 4080; y++ {
            dex                ; $c037: ca                  ppu.DATA = cpu.A;           // Set all nametable to 0x00
            bne __c034         ; $c038: d0 fa           .        
            dey                ; $c03a: 88              .
            bne __c034         ; $c03b: d0 f7           }

            lda #$3f           ; $c03d: a9 3f           .
            sta PPUADDR        ; $c03f: 8d 06 20        .
            lda #$00           ; $c042: a9 00           .
            sta PPUADDR        ; $c044: 8d 06 20        ppu.ADDR = 0x3F00               // Set palette indexes
            ldx #$00           ; $c047: a2 00           cpu.x = 0x00        
                                                        // Copy palette indexes to PPU
__c049:     lda __ff78,x       ; $c049: bd 78 ff        do{
            sta PPUDATA        ; $c04c: 8d 07 20            cpu.A = $0xFF78 + cpu.X         
            inx                ; $c04f: e8                  ppu.DATA = cpu.A    
            cpx #$20           ; $c050: e0 20               x++
            bne __c049         ; $c052: d0 f5           } while (cpu.X != 0x20)

            lda #$c0           ; $c054: a9 c0           cpu.A = 0xC0
            sta JOY2           ; $c056: 8d 17 40        joypad2 = 0xC0                  // ??
            lda #$00           ; $c059: a9 00           cpu.A = 0x00
            sta SND_CHN        ; $c05b: 8d 15 40        soudnChannel = 0x00             // ??
            lda #$78           ; $c05e: a9 78           cpu.A = 0x78
            sta $d0            ; $c060: 85 d0           write(0xD0, cpu.A)              // ??
            lda #$fb           ; $c062: a9 fb           cpu.A = 0xFB
            sta $d1            ; $c064: 85 d1           write(0xD1, cpu.A)
            lda #$7f           ; $c066: a9 7f           cpu.A = 0x7F
            sta $d3            ; $c068: 85 d3           write(0xD3, cpu.A)
            ldy #$00           ; $c06a: a0 00           .
            sty PPUADDR        ; $c06c: 8c 06 20        .
            sty PPUADDR        ; $c06f: 8c 06 20        ppu.ADDR = 0x0000               // Set pattern table address
__c072:     lda #$00           ; $c072: a9 00           cpu.A = 0x00
            sta $d7            ; $c074: 85 d7           write(0xD7, cpu.A)
            lda #$07           ; $c076: a9 07           cpu.A = 0x07
            sta $d0            ; $c078: 85 d0           write(0xD0, cpu.A)
            lda #$c3           ; $c07a: a9 c3           cpu.A = 0xC3
            sta $d1            ; $c07c: 85 d1           write(0xD1, cpu.A)
            jsr __c2a7         ; $c07e: 20 a7 c2        
__c081:     jsr __c28d         ; $c081: 20 8d c2  
            ldx #$12           ; $c084: a2 12     
            jsr __c261         ; $c086: 20 61 c2  
            lda $d5            ; $c089: a5 d5     
            lsr                ; $c08b: 4a        
            lsr                ; $c08c: 4a        
            lsr                ; $c08d: 4a        
            bcs __c0ac         ; $c08e: b0 1c     
            lsr                ; $c090: 4a        
            bcs __c09f         ; $c091: b0 0c     
            lsr                ; $c093: 4a        
            bcs __c0bd         ; $c094: b0 27     
            lsr                ; $c096: 4a        
            bcs __c09c         ; $c097: b0 03     
            jmp __c081         ; $c099: 4c 81 c0  

;-------------------------------------------------------------------------------
__c09c:     jmp __c126         ; $c09c: 4c 26 c1  

;-------------------------------------------------------------------------------
__c09f:     jsr __c66f         ; $c09f: 20 6f c6  
            dec $d7            ; $c0a2: c6 d7     
            bpl __c081         ; $c0a4: 10 db     
            lda #$0d           ; $c0a6: a9 0d     
            sta $d7            ; $c0a8: 85 d7     
            bne __c081         ; $c0aa: d0 d5     
__c0ac:     jsr __c66f         ; $c0ac: 20 6f c6  
            inc $d7            ; $c0af: e6 d7     
            lda $d7            ; $c0b1: a5 d7     
            cmp #$0e           ; $c0b3: c9 0e     
            bcc __c081         ; $c0b5: 90 ca     
            lda #$00           ; $c0b7: a9 00     
            sta $d7            ; $c0b9: 85 d7     
            beq __c081         ; $c0bb: f0 c4     
__c0bd:     jsr __c689         ; $c0bd: 20 89 c6  
            lda $d7            ; $c0c0: a5 d7     
            beq __c0ca         ; $c0c2: f0 06     
            jsr __c0ed         ; $c0c4: 20 ed c0  
            jmp __c081         ; $c0c7: 4c 81 c0  

;-------------------------------------------------------------------------------
__c0ca:     lda #$00           ; $c0ca: a9 00     
            sta $d8            ; $c0cc: 85 d8     
            inc $d7            ; $c0ce: e6 d7     
__c0d0:     jsr __c0ed         ; $c0d0: 20 ed c0  
            inc $d7            ; $c0d3: e6 d7     
            lda $d7            ; $c0d5: a5 d7     
            cmp #$0e           ; $c0d7: c9 0e     
            bne __c0d0         ; $c0d9: d0 f5     
            lda #$00           ; $c0db: a9 00     
            sta $d7            ; $c0dd: 85 d7     
            lda $d8            ; $c0df: a5 d8     
            beq __c0e5         ; $c0e1: f0 02     
            lda #$ff           ; $c0e3: a9 ff     
__c0e5:     sta $00            ; $c0e5: 85 00     
            jsr __c1ed         ; $c0e7: 20 ed c1  
            jmp __c081         ; $c0ea: 4c 81 c0  

;-------------------------------------------------------------------------------
__c0ed:     lda $d7            ; $c0ed: a5 d7     
            asl                ; $c0ef: 0a        
            tax                ; $c0f0: aa        
            lda __c10a,x       ; $c0f1: bd 0a c1  
            sta $0200          ; $c0f4: 8d 00 02  
            lda __c10b,x       ; $c0f7: bd 0b c1  
            sta $0201          ; $c0fa: 8d 01 02  
            lda #$c1           ; $c0fd: a9 c1     
            pha                ; $c0ff: 48        
            lda #$de           ; $c100: a9 de     
            pha                ; $c102: 48        
            lda #$00           ; $c103: a9 00     
            sta $00            ; $c105: 85 00     
            jmp ($0200)        ; $c107: 6c 00 02  

;-------------------------------------------------------------------------------
__c10a:     .hex 2d            ; $c10a: 2d        Suspected data
__c10b:     .hex c7 2d         ; $c10b: c7 2d     Invalid Opcode - DCP $2d
            .hex c7 db         ; $c10d: c7 db     Invalid Opcode - DCP $db
            .hex c7 85         ; $c10f: c7 85     Invalid Opcode - DCP $85
            iny                ; $c111: c8        
            dec __f8cb,x       ; $c112: de cb f8  
            cmp __ceee         ; $c115: cd ee ce  
            ldx #$cf           ; $c118: a2 cf     
            .hex 74 d1         ; $c11a: 74 d1     Invalid Opcode - NOP $d1,x
            .hex fb d4 d4      ; $c11c: fb d4 d4  Invalid Opcode - ISC __d4d4,y
            cmp ($4a,x)        ; $c11f: c1 4a     
            .hex df b8 db      ; $c121: df b8 db  Invalid Opcode - DCP __dbb8,x
            tax                ; $c124: aa        
            .hex e1            ; $c125: e1        Suspected data
__c126:     lda #$00           ; $c126: a9 00     
            sta $d7            ; $c128: 85 d7     
            lda #$92           ; $c12a: a9 92     
            sta $d0            ; $c12c: 85 d0     
            lda #$c4           ; $c12e: a9 c4     
            sta $d1            ; $c130: 85 d1     
            jsr __c2a7         ; $c132: 20 a7 c2  
__c135:     jsr __c28d         ; $c135: 20 8d c2  
            ldx #$0f           ; $c138: a2 0f     
            jsr __c261         ; $c13a: 20 61 c2  
            lda $d5            ; $c13d: a5 d5     
            lsr                ; $c13f: 4a        
            lsr                ; $c140: 4a        
            lsr                ; $c141: 4a        
            bcs __c160         ; $c142: b0 1c     
            lsr                ; $c144: 4a        
            bcs __c153         ; $c145: b0 0c     
            lsr                ; $c147: 4a        
            bcs __c171         ; $c148: b0 27     
            lsr                ; $c14a: 4a        
            bcs __c150         ; $c14b: b0 03     
            jmp __c135         ; $c14d: 4c 35 c1  

;-------------------------------------------------------------------------------
__c150:     jmp __c072         ; $c150: 4c 72 c0  

;-------------------------------------------------------------------------------
__c153:     jsr __c66f         ; $c153: 20 6f c6  
            dec $d7            ; $c156: c6 d7     
            bpl __c135         ; $c158: 10 db     
            lda #$0a           ; $c15a: a9 0a     
            sta $d7            ; $c15c: 85 d7     
            bne __c135         ; $c15e: d0 d5     
__c160:     jsr __c66f         ; $c160: 20 6f c6  
            inc $d7            ; $c163: e6 d7     
            lda $d7            ; $c165: a5 d7     
            cmp #$0b           ; $c167: c9 0b     
            bcc __c135         ; $c169: 90 ca     
            lda #$00           ; $c16b: a9 00     
            sta $d7            ; $c16d: 85 d7     
            beq __c135         ; $c16f: f0 c4     
__c171:     jsr __c689         ; $c171: 20 89 c6  
            lda $d7            ; $c174: a5 d7     
            beq __c17e         ; $c176: f0 06     
            jsr __c1a1         ; $c178: 20 a1 c1  
            jmp __c135         ; $c17b: 4c 35 c1  

;-------------------------------------------------------------------------------
__c17e:     lda #$00           ; $c17e: a9 00     
            sta $d8            ; $c180: 85 d8     
            inc $d7            ; $c182: e6 d7     
__c184:     jsr __c1a1         ; $c184: 20 a1 c1  
            inc $d7            ; $c187: e6 d7     
            lda $d7            ; $c189: a5 d7     
            cmp #$0b           ; $c18b: c9 0b     
            bne __c184         ; $c18d: d0 f5     
            lda #$00           ; $c18f: a9 00     
            sta $d7            ; $c191: 85 d7     
            lda $d8            ; $c193: a5 d8     
            beq __c199         ; $c195: f0 02     
            lda #$ff           ; $c197: a9 ff     
__c199:     sta $00            ; $c199: 85 00     
            jsr __c1ed         ; $c19b: 20 ed c1  
            jmp __c135         ; $c19e: 4c 35 c1  

;-------------------------------------------------------------------------------
__c1a1:     lda $d7            ; $c1a1: a5 d7     
            asl                ; $c1a3: 0a        
            tax                ; $c1a4: aa        
            lda __c1be,x       ; $c1a5: bd be c1  
            sta $0200          ; $c1a8: 8d 00 02  
            lda __c1bf,x       ; $c1ab: bd bf c1  
            sta $0201          ; $c1ae: 8d 01 02  
            lda #$c1           ; $c1b1: a9 c1     
            pha                ; $c1b3: 48        
            lda #$de           ; $c1b4: a9 de     
            pha                ; $c1b6: 48        
            lda #$00           ; $c1b7: a9 00     
            sta $00            ; $c1b9: 85 00     
            jmp ($0200)        ; $c1bb: 6c 00 02  

;-------------------------------------------------------------------------------
__c1be:     .hex a3            ; $c1be: a3        Suspected data
__c1bf:     dec $a3            ; $c1bf: c6 a3     
            dec $1e            ; $c1c1: c6 1e     
            sbc $3d            ; $c1c3: e5 3d     
            .hex e7 d3         ; $c1c5: e7 d3     Invalid Opcode - ISC $d3
            inx                ; $c1c7: e8        
            asl $e9,x          ; $c1c8: 16 e9     
            stx $eb            ; $c1ca: 86 eb     
            inc $ed,x          ; $c1cc: f6 ed     
            ror $f0            ; $c1ce: 66 f0     
            dec $f2,x          ; $c1d0: d6 f2     
            lsr $f5            ; $c1d2: 46 f5     
            lda #$00           ; $c1d4: a9 00     
            sta $00            ; $c1d6: 85 00     
            jsr __d900         ; $c1d8: 20 00 d9  
            jsr __dae0         ; $c1db: 20 e0 da  
            nop                ; $c1de: ea        
            nop                ; $c1df: ea        
            nop                ; $c1e0: ea        
            lda $00            ; $c1e1: a5 00     
            beq __c1e7         ; $c1e3: f0 02     
            sta $d8            ; $c1e5: 85 d8     
__c1e7:     jmp __c1ed         ; $c1e7: 4c ed c1  

;-------------------------------------------------------------------------------
            jmp __c081         ; $c1ea: 4c 81 c0  

;-------------------------------------------------------------------------------
__c1ed:     jsr __c28d         ; $c1ed: 20 8d c2  
            lda #$00           ; $c1f0: a9 00     
            sta $d3            ; $c1f2: 85 d3     
            lda $d7            ; $c1f4: a5 d7     
            clc                ; $c1f6: 18        
            adc #$04           ; $c1f7: 69 04     
            asl                ; $c1f9: 0a        
            rol $d3            ; $c1fa: 26 d3     
            asl                ; $c1fc: 0a        
            rol $d3            ; $c1fd: 26 d3     
            asl                ; $c1ff: 0a        
            rol $d3            ; $c200: 26 d3     
            asl                ; $c202: 0a        
            rol $d3            ; $c203: 26 d3     
            asl                ; $c205: 0a        
            rol $d3            ; $c206: 26 d3     
            pha                ; $c208: 48        
            lda $d3            ; $c209: a5 d3     
            ora #$20           ; $c20b: 09 20     
            sta PPUADDR        ; $c20d: 8d 06 20  
            pla                ; $c210: 68        
            ora #$04           ; $c211: 09 04     
            sta PPUADDR        ; $c213: 8d 06 20  
            lda $00            ; $c216: a5 00     
            beq __c237         ; $c218: f0 1d     
            cmp #$ff           ; $c21a: c9 ff     
            beq __c244         ; $c21c: f0 26     
            lsr                ; $c21e: 4a        
            lsr                ; $c21f: 4a        
            lsr                ; $c220: 4a        
            lsr                ; $c221: 4a        
            tax                ; $c222: aa        
            lda __c251,x       ; $c223: bd 51 c2  
            sta PPUDATA        ; $c226: 8d 07 20  
            lda $00            ; $c229: a5 00     
            and #$0f           ; $c22b: 29 0f     
            tax                ; $c22d: aa        
            lda __c251,x       ; $c22e: bd 51 c2  
            sta PPUDATA        ; $c231: 8d 07 20  
            jmp __c294         ; $c234: 4c 94 c2  

;-------------------------------------------------------------------------------
__c237:     lda #$4f           ; $c237: a9 4f     
            sta PPUDATA        ; $c239: 8d 07 20  
            lda #$4b           ; $c23c: a9 4b     
            sta PPUDATA        ; $c23e: 8d 07 20  
            jmp __c294         ; $c241: 4c 94 c2  

;-------------------------------------------------------------------------------
__c244:     lda #$45           ; $c244: a9 45     
            sta PPUDATA        ; $c246: 8d 07 20  
            lda #$72           ; $c249: a9 72     
            sta PPUDATA        ; $c24b: 8d 07 20  
            jmp __c294         ; $c24e: 4c 94 c2  

;-------------------------------------------------------------------------------
__c251:     bmi __c284         ; $c251: 30 31     
            .hex 32            ; $c253: 32        Invalid Opcode - KIL 
            .hex 33 34         ; $c254: 33 34     Invalid Opcode - RLA ($34),y
            and $36,x          ; $c256: 35 36     
            .hex 37 38         ; $c258: 37 38     Invalid Opcode - RLA $38,x
            and $4241,y        ; $c25a: 39 41 42  
            .hex 43 44         ; $c25d: 43 44     Invalid Opcode - SRE ($44,x)
            eor $46            ; $c25f: 45 46     
__c261:     lda $d7            ; $c261: a5 d7     
            clc                ; $c263: 18        
            adc #$04           ; $c264: 69 04     
            tay                ; $c266: a8        
            lda #$84           ; $c267: a9 84     
            sta PPUCTRL        ; $c269: 8d 00 20  
            lda #$20           ; $c26c: a9 20     
            sta PPUADDR        ; $c26e: 8d 06 20  
            lda #$02           ; $c271: a9 02     
            sta PPUADDR        ; $c273: 8d 06 20  
            lda #$20           ; $c276: a9 20     
            dey                ; $c278: 88        
            iny                ; $c279: c8        
            bne __c27e         ; $c27a: d0 02     
            lda #$2a           ; $c27c: a9 2a     
__c27e:     sta PPUDATA        ; $c27e: 8d 07 20  
            dey                ; $c281: 88        
            dex                ; $c282: ca        
            .hex d0            ; $c283: d0        Suspected data
__c284:     sbc ($a9),y        ; $c284: f1 a9     
            .hex 80 8d         ; $c286: 80 8d     Invalid Opcode - NOP #$8d
            brk                ; $c288: 00        
            jsr $944c          ; $c289: 20 4c 94  
            .hex c2            ; $c28c: c2        Suspected data
__c28d:     lda $d2            ; $c28d: a5 d2     
__c28f:     cmp $d2            ; $c28f: c5 d2     
            beq __c28f         ; $c291: f0 fc     
            rts                ; $c293: 60        

;-------------------------------------------------------------------------------
__c294:     lda #$00           ; $c294: a9 00     
            sta PPUSCROLL      ; $c296: 8d 05 20  
            sta PPUSCROLL      ; $c299: 8d 05 20  
            lda #$00           ; $c29c: a9 00     
            sta PPUADDR        ; $c29e: 8d 06 20  
            lda #$00           ; $c2a1: a9 00     
            sta PPUADDR        ; $c2a3: 8d 06 20  
            rts                ; $c2a6: 60        

;-------------------------------------------------------------------------------
__c2a7:     lda #$00           ; $c2a7: a9 00           cpu.A = 0x00
            sta PPUCTRL        ; $c2a9: 8d 00 20        ppu.CTRL = cpu.A
            sta PPUMASK        ; $c2ac: 8d 01 20        ppu.MASK = cpu.A
            jsr clearScreen    ; $c2af: 20 ed c2        clearScreen()
            lda #$20           ; $c2b2: a9 20           .
            sta PPUADDR        ; $c2b4: 8d 06 20        .
            ldy #$00           ; $c2b7: a0 00           .
            sty PPUADDR        ; $c2b9: 8c 06 20            ppu.ADDR = 0x2000
__c2bc:     ldx #$20           ; $c2bc: a2 20               cpu.X = 0x20
__c2be:     lda ($d0),y                 ; $c2be: b1 d0      cpu.A = $0xD0 + cpu.Y
            beq ppuForBckgRendering     ; $c2c0: f0 20      if (cpu.Zero) {
                                                                ppuForBckgRendering()
                                                            }

            cmp #$ff                    ; $c2c2: c9 ff      if (cpu.A == 0xFF) {
                                                                cpu.Carry = true
                                                                cpu.Zero = true
                                                                cpu.NegativeFlag = (cpu.A - 0xFF) < 0
                                                            }
            beq __c2d3         ; $c2c4: f0 0d               if (cpu.Zero) {
                                                                __c2d3()
                                                            }
            sta PPUDATA        ; $c2c6: 8d 07 20  
            iny                ; $c2c9: c8        
            bne __c2ce         ; $c2ca: d0 02     
            inc $d1            ; $c2cc: e6 d1     
__c2ce:     dex                ; $c2ce: ca        
            bne __c2be         ; $c2cf: d0 ed     
            beq __c2bc         ; $c2d1: f0 e9     
__c2d3:     iny                ; $c2d3: c8              cpu.Y--
            bne __c2d8         ; $c2d4: d0 02           if (cpu.Y > 0) {
                                                            __c2d8()
                                                        }

            inc $d1            ; $c2d6: e6 d1           
__c2d8:     lda #$20           ; $c2d8: a9 20           cpu.A = 0x20
            sta PPUDATA        ; $c2da: 8d 07 20        ppu.DATA = cpu.A
            dex                ; $c2dd: ca              cpu.X--
            bne __c2d8         ; $c2de: d0 f8           if (cpu.X > 0) {
                                                            __c2d8()
                                                        }        
            beq __c2bc         ; $c2e0: f0 da     
ppuForBckgRendering:     lda #$80           ; $c2e2: a9 80           cpu.A = 0x80
            sta PPUCTRL        ; $c2e4: 8d 00 20        ppu.CTRL = 0x80     //  generateNMIAtVBlank
            lda #$0e           ; $c2e7: a9 0e           cpu.A = 0x0E
            sta PPUMASK        ; $c2e9: 8d 01 20        ppu.MASK = 0x0E     //  showBackgroundLeftEdge
                                                                            //  showSpritesLeftEdge
                                                                            //  showBackground

            rts                ; $c2ec: 60              

;-------------------------------------------------------------------------------
clearScreen:     lda #$20           ; $c2ed: a9 20           .
            sta PPUADDR        ; $c2ef: 8d 06 20        .
            lda #$00           ; $c2f2: a9 00           .
            sta PPUADDR        ; $c2f4: 8d 06 20        ppu.ADDR = $0x2000
            ldx #$1e           ; $c2f7: a2 1e           . 
            lda #$20           ; $c2f9: a9 20           . 
__c2fb:     ldy #$20           ; $c2fb: a0 20           . 
__c2fd:     sta PPUDATA        ; $c2fd: 8d 07 20        for (x = 0; x<30; x++)     
            dey                ; $c300: 88                  for (y = 0; y<32; y++) {     
            bne __c2fd         ; $c301: d0 fa                   ppu.DATA = 0x20;
                                                            }
            dex                ; $c303: ca              .
            bne __c2fb         ; $c304: d0 f5           }
            rts                ; $c306: 60        

;-------------------------------------------------------------------------------
            .hex ff ff ff      ; $c307: ff ff ff  Invalid Opcode - ISC $ffff,x
            .hex ff 20 20      ; $c30a: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2d20          ; $c30d: 20 20 2d  
            and $5220          ; $c310: 2d 20 52  
            adc $6e,x          ; $c313: 75 6e     
            jsr $6c61          ; $c315: 20 61 6c  
            jmp ($7420)        ; $c318: 6c 20 74  

;-------------------------------------------------------------------------------
            adc $73            ; $c31b: 65 73     
            .hex 74 73         ; $c31d: 74 73     Invalid Opcode - NOP $73,x
            .hex ff 20 20      ; $c31f: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2d20          ; $c322: 20 20 2d  
            and $4220          ; $c325: 2d 20 42  
            .hex 72            ; $c328: 72        Invalid Opcode - KIL 
            adc ($6e,x)        ; $c329: 61 6e     
            .hex 63 68         ; $c32b: 63 68     Invalid Opcode - RRA ($68,x)
            jsr $6574          ; $c32d: 20 74 65  
            .hex 73 74         ; $c330: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c332: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c334: 20 20 20  
            jsr $2d2d          ; $c337: 20 2d 2d  
            jsr $6c46          ; $c33a: 20 46 6c  
            adc ($67,x)        ; $c33d: 61 67     
            jsr $6574          ; $c33f: 20 74 65  
            .hex 73 74         ; $c342: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c344: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c346: 20 20 20  
            jsr $2d2d          ; $c349: 20 2d 2d  
            jsr $6d49          ; $c34c: 20 49 6d  
            adc $6465          ; $c34f: 6d 65 64  
            adc #$61           ; $c352: 69 61     
            .hex 74 65         ; $c354: 74 65     Invalid Opcode - NOP $65,x
            jsr $6574          ; $c356: 20 74 65  
            .hex 73 74         ; $c359: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c35b: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c35d: 20 20 20  
            jsr $2d2d          ; $c360: 20 2d 2d  
            jsr $6d49          ; $c363: 20 49 6d  
            bvs __c3d4         ; $c366: 70 6c     
            adc #$65           ; $c368: 69 65     
            .hex 64 20         ; $c36a: 64 20     Invalid Opcode - NOP $20
            .hex 74 65         ; $c36c: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c36e: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c370: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c372: 20 20 20  
            jsr $2d2d          ; $c375: 20 2d 2d  
            jsr $7453          ; $c378: 20 53 74  
            adc ($63,x)        ; $c37b: 61 63     
            .hex 6b 20         ; $c37d: 6b 20     Invalid Opcode - ARR #$20
            .hex 74 65         ; $c37f: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c381: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c383: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c385: 20 20 20  
            jsr $2d2d          ; $c388: 20 2d 2d  
            jsr $6341          ; $c38b: 20 41 63  
            .hex 63 75         ; $c38e: 63 75     Invalid Opcode - RRA ($75,x)
            adc $6c75          ; $c390: 6d 75 6c  
            adc ($74,x)        ; $c393: 61 74     
            .hex 6f 72 20      ; $c395: 6f 72 20  Invalid Opcode - RRA $2072
            .hex 74 65         ; $c398: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c39a: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c39c: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c39e: 20 20 20  
            jsr $2d2d          ; $c3a1: 20 2d 2d  
            jsr $4928          ; $c3a4: 20 28 49  
            ror $6964          ; $c3a7: 6e 64 69  
            .hex 72            ; $c3aa: 72        Invalid Opcode - KIL 
            adc $63            ; $c3ab: 65 63     
            .hex 74 2c         ; $c3ad: 74 2c     Invalid Opcode - NOP $2c,x
            cli                ; $c3af: 58        
            and #$20           ; $c3b0: 29 20     
            .hex 74 65         ; $c3b2: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c3b4: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c3b6: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c3b8: 20 20 20  
            jsr $2d2d          ; $c3bb: 20 2d 2d  
            jsr $655a          ; $c3be: 20 5a 65  
            .hex 72            ; $c3c1: 72        Invalid Opcode - KIL 
            .hex 6f 70 61      ; $c3c2: 6f 70 61  Invalid Opcode - RRA $6170
            .hex 67 65         ; $c3c5: 67 65     Invalid Opcode - RRA $65
            jsr $6574          ; $c3c7: 20 74 65  
            .hex 73 74         ; $c3ca: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c3cc: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c3ce: 20 20 20  
            jsr $2d2d          ; $c3d1: 20 2d 2d  
__c3d4:     jsr $6241          ; $c3d4: 20 41 62  
            .hex 73 6f         ; $c3d7: 73 6f     Invalid Opcode - RRA ($6f),y
            jmp ($7475)        ; $c3d9: 6c 75 74  

;-------------------------------------------------------------------------------
            adc $20            ; $c3dc: 65 20     
            .hex 74 65         ; $c3de: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c3e0: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c3e2: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c3e4: 20 20 20  
            jsr $2d2d          ; $c3e7: 20 2d 2d  
            jsr $4928          ; $c3ea: 20 28 49  
            ror $6964          ; $c3ed: 6e 64 69  
            .hex 72            ; $c3f0: 72        Invalid Opcode - KIL 
            adc $63            ; $c3f1: 65 63     
            .hex 74 29         ; $c3f3: 74 29     Invalid Opcode - NOP $29,x
            bit $2059          ; $c3f5: 2c 59 20  
            .hex 74 65         ; $c3f8: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c3fa: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c3fc: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c3fe: 20 20 20  
            jsr $2d2d          ; $c401: 20 2d 2d  
            jsr $6241          ; $c404: 20 41 62  
            .hex 73 6f         ; $c407: 73 6f     Invalid Opcode - RRA ($6f),y
            jmp ($7475)        ; $c409: 6c 75 74  

;-------------------------------------------------------------------------------
            adc $2c            ; $c40c: 65 2c     
            eor $7420,y        ; $c40e: 59 20 74  
            adc $73            ; $c411: 65 73     
            .hex 74 73         ; $c413: 74 73     Invalid Opcode - NOP $73,x
            .hex ff 20 20      ; $c415: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2d20          ; $c418: 20 20 2d  
            and $5a20          ; $c41b: 2d 20 5a  
            adc $72            ; $c41e: 65 72     
            .hex 6f 70 61      ; $c420: 6f 70 61  Invalid Opcode - RRA $6170
            .hex 67 65         ; $c423: 67 65     Invalid Opcode - RRA $65
            bit $2058          ; $c425: 2c 58 20  
            .hex 74 65         ; $c428: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c42a: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c42c: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c42e: 20 20 20  
            jsr $2d2d          ; $c431: 20 2d 2d  
            jsr $6241          ; $c434: 20 41 62  
            .hex 73 6f         ; $c437: 73 6f     Invalid Opcode - RRA ($6f),y
            jmp ($7475)        ; $c439: 6c 75 74  

;-------------------------------------------------------------------------------
            adc $2c            ; $c43c: 65 2c     
            cli                ; $c43e: 58        
            jsr $6574          ; $c43f: 20 74 65  
            .hex 73 74         ; $c442: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c444: 73 ff     Invalid Opcode - RRA ($ff),y
            .hex ff ff 20      ; $c446: ff ff 20  Invalid Opcode - ISC $20ff,x
            jsr $2020          ; $c449: 20 20 20  
            eor $70,x          ; $c44c: 55 70     
            .hex 2f 44 6f      ; $c44e: 2f 44 6f  Invalid Opcode - RLA $6f44
            .hex 77 6e         ; $c451: 77 6e     Invalid Opcode - RRA $6e,x
            .hex 3a            ; $c453: 3a        Invalid Opcode - NOP 
            jsr $6573          ; $c454: 20 73 65  
            jmp ($6365)        ; $c457: 6c 65 63  

;-------------------------------------------------------------------------------
            .hex 74 20         ; $c45a: 74 20     Invalid Opcode - NOP $20,x
            .hex 74 65         ; $c45c: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c45e: 73 74     Invalid Opcode - RRA ($74),y
            .hex ff 20 20      ; $c460: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2020          ; $c463: 20 20 20  
            jsr $7453          ; $c466: 20 53 74  
            adc ($72,x)        ; $c469: 61 72     
            .hex 74 3a         ; $c46b: 74 3a     Invalid Opcode - NOP $3a,x
            jsr $7572          ; $c46d: 20 72 75  
            ror $7420          ; $c470: 6e 20 74  
            adc $73            ; $c473: 65 73     
            .hex 74 ff         ; $c475: 74 ff     Invalid Opcode - NOP $ff,x
            jsr $2020          ; $c477: 20 20 20  
            jsr $5320          ; $c47a: 20 20 53  
            adc $6c            ; $c47d: 65 6c     
            adc $63            ; $c47f: 65 63     
            .hex 74 3a         ; $c481: 74 3a     Invalid Opcode - NOP $3a,x
            jsr $6e49          ; $c483: 20 49 6e  
            ror $61,x          ; $c486: 76 61     
            jmp ($6469)        ; $c488: 6c 69 64  

;-------------------------------------------------------------------------------
            jsr $706f          ; $c48b: 20 6f 70  
            .hex 73 21         ; $c48e: 73 21     Invalid Opcode - RRA ($21),y
            .hex ff 00 ff      ; $c490: ff 00 ff  Invalid Opcode - ISC __ff00,x
            .hex ff ff ff      ; $c493: ff ff ff  Invalid Opcode - ISC $ffff,x
            jsr $2020          ; $c496: 20 20 20  
            jsr $2d2d          ; $c499: 20 2d 2d  
            jsr $7552          ; $c49c: 20 52 75  
            ror $6120          ; $c49f: 6e 20 61  
            jmp ($206c)        ; $c4a2: 6c 6c 20  

;-------------------------------------------------------------------------------
            .hex 74 65         ; $c4a5: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c4a7: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c4a9: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c4ab: 20 20 20  
            jsr $2d2d          ; $c4ae: 20 2d 2d  
            jsr $4f4e          ; $c4b1: 20 4e 4f  
            bvc __c4d6         ; $c4b4: 50 20     
            .hex 74 65         ; $c4b6: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c4b8: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c4ba: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c4bc: 20 20 20  
            jsr $2d2d          ; $c4bf: 20 2d 2d  
            jsr $414c          ; $c4c2: 20 4c 41  
            cli                ; $c4c5: 58        
            jsr $6574          ; $c4c6: 20 74 65  
            .hex 73 74         ; $c4c9: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c4cb: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c4cd: 20 20 20  
            jsr $2d2d          ; $c4d0: 20 2d 2d  
            jsr $4153          ; $c4d3: 20 53 41  
__c4d6:     cli                ; $c4d6: 58        
            jsr $6574          ; $c4d7: 20 74 65  
            .hex 73 74         ; $c4da: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c4dc: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c4de: 20 20 20  
            jsr $2d2d          ; $c4e1: 20 2d 2d  
            jsr $4253          ; $c4e4: 20 53 42  
            .hex 43 20         ; $c4e7: 43 20     Invalid Opcode - SRE ($20,x)
            .hex 74 65         ; $c4e9: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c4eb: 73 74     Invalid Opcode - RRA ($74),y
            jsr $6f28          ; $c4ed: 20 28 6f  
            bvs __c555         ; $c4f0: 70 63     
            .hex 6f 64 65      ; $c4f2: 6f 64 65  Invalid Opcode - RRA $6564
            jsr $4530          ; $c4f5: 20 30 45  
            .hex 42            ; $c4f8: 42        Invalid Opcode - KIL 
            pla                ; $c4f9: 68        
            and #$ff           ; $c4fa: 29 ff     
            jsr $2020          ; $c4fc: 20 20 20  
            jsr $2d2d          ; $c4ff: 20 2d 2d  
            jsr $4344          ; $c502: 20 44 43  
            bvc __c527         ; $c505: 50 20     
            .hex 74 65         ; $c507: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c509: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c50b: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c50d: 20 20 20  
            jsr $2d2d          ; $c510: 20 2d 2d  
            jsr $5349          ; $c513: 20 49 53  
            .hex 42            ; $c516: 42        Invalid Opcode - KIL 
            jsr $6574          ; $c517: 20 74 65  
            .hex 73 74         ; $c51a: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c51c: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c51e: 20 20 20  
            jsr $2d2d          ; $c521: 20 2d 2d  
            jsr $4c53          ; $c524: 20 53 4c  
__c527:     .hex 4f 20 74      ; $c527: 4f 20 74  Invalid Opcode - SRE $7420
            adc $73            ; $c52a: 65 73     
            .hex 74 73         ; $c52c: 74 73     Invalid Opcode - NOP $73,x
            .hex ff 20 20      ; $c52e: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2d20          ; $c531: 20 20 2d  
            and $5220          ; $c534: 2d 20 52  
            jmp $2041          ; $c537: 4c 41 20  

;-------------------------------------------------------------------------------
            .hex 74 65         ; $c53a: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c53c: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c53e: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c540: 20 20 20  
            jsr $2d2d          ; $c543: 20 2d 2d  
            jsr $5253          ; $c546: 20 53 52  
            eor $20            ; $c549: 45 20     
            .hex 74 65         ; $c54b: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c54d: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c54f: 73 ff     Invalid Opcode - RRA ($ff),y
            jsr $2020          ; $c551: 20 20 20  
            .hex 20            ; $c554: 20        Suspected data
__c555:     and $202d          ; $c555: 2d 2d 20  
            .hex 52            ; $c558: 52        Invalid Opcode - KIL 
            .hex 52            ; $c559: 52        Invalid Opcode - KIL 
            eor ($20,x)        ; $c55a: 41 20     
            .hex 74 65         ; $c55c: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c55e: 73 74     Invalid Opcode - RRA ($74),y
            .hex 73 ff         ; $c560: 73 ff     Invalid Opcode - RRA ($ff),y
            .hex ff ff ff      ; $c562: ff ff ff  Invalid Opcode - ISC $ffff,x
            .hex ff ff 20      ; $c565: ff ff 20  Invalid Opcode - ISC $20ff,x
            jsr $2020          ; $c568: 20 20 20  
            eor $70,x          ; $c56b: 55 70     
            .hex 2f 44 6f      ; $c56d: 2f 44 6f  Invalid Opcode - RLA $6f44
            .hex 77 6e         ; $c570: 77 6e     Invalid Opcode - RRA $6e,x
            .hex 3a            ; $c572: 3a        Invalid Opcode - NOP 
            jsr $6573          ; $c573: 20 73 65  
            jmp ($6365)        ; $c576: 6c 65 63  

;-------------------------------------------------------------------------------
            .hex 74 20         ; $c579: 74 20     Invalid Opcode - NOP $20,x
            .hex 74 65         ; $c57b: 74 65     Invalid Opcode - NOP $65,x
            .hex 73 74         ; $c57d: 73 74     Invalid Opcode - RRA ($74),y
            .hex ff 20 20      ; $c57f: ff 20 20  Invalid Opcode - ISC $2020,x
            jsr $2020          ; $c582: 20 20 20  
            jsr $7453          ; $c585: 20 53 74  
            adc ($72,x)        ; $c588: 61 72     
            .hex 74 3a         ; $c58a: 74 3a     Invalid Opcode - NOP $3a,x
            jsr $7572          ; $c58c: 20 72 75  
            ror $7420          ; $c58f: 6e 20 74  
            adc $73            ; $c592: 65 73     
            .hex 74 ff         ; $c594: 74 ff     Invalid Opcode - NOP $ff,x
            jsr $2020          ; $c596: 20 20 20  
            jsr $5320          ; $c599: 20 20 53  
            adc $6c            ; $c59c: 65 6c     
            adc $63            ; $c59e: 65 63     
            .hex 74 3a         ; $c5a0: 74 3a     Invalid Opcode - NOP $3a,x
            jsr $6f4e          ; $c5a2: 20 4e 6f  
            .hex 72            ; $c5a5: 72        Invalid Opcode - KIL 
            adc $6c61          ; $c5a6: 6d 61 6c  
            jsr $706f          ; $c5a9: 20 6f 70  
            .hex 73 ff         ; $c5ac: 73 ff     Invalid Opcode - RRA ($ff),y
            brk                ; $c5ae: 00        

;-------------------------------------------------------------------------------
; nmi vector
;-------------------------------------------------------------------------------
nmi:        pha                ; $c5af: 48        
            txa                ; $c5b0: 8a        
            pha                ; $c5b1: 48        
            lda PPUSTATUS      ; $c5b2: ad 02 20  
            lda #$20           ; $c5b5: a9 20     
            sta PPUADDR        ; $c5b7: 8d 06 20  
            lda #$40           ; $c5ba: a9 40     
            sta PPUADDR        ; $c5bc: 8d 06 20  
            inc $d2            ; $c5bf: e6 d2     
            lda #$00           ; $c5c1: a9 00     
            sta PPUSCROLL      ; $c5c3: 8d 05 20  
            sta PPUSCROLL      ; $c5c6: 8d 05 20  
            lda #$00           ; $c5c9: a9 00     
            sta PPUADDR        ; $c5cb: 8d 06 20  
            lda #$00           ; $c5ce: a9 00     
            sta PPUADDR        ; $c5d0: 8d 06 20  
            ldx #$09           ; $c5d3: a2 09     
            stx JOY1           ; $c5d5: 8e 16 40  
            dex                ; $c5d8: ca        
            stx JOY1           ; $c5d9: 8e 16 40  
__c5dc:     lda JOY1           ; $c5dc: ad 16 40  
            lsr                ; $c5df: 4a        
            rol $d4            ; $c5e0: 26 d4     
            dex                ; $c5e2: ca        
            bne __c5dc         ; $c5e3: d0 f7     
            lda $d4            ; $c5e5: a5 d4     
            tax                ; $c5e7: aa        
            eor $d6            ; $c5e8: 45 d6     
            and $d4            ; $c5ea: 25 d4     
            sta $d5            ; $c5ec: 85 d5     
            stx $d6            ; $c5ee: 86 d6     
            pla                ; $c5f0: 68        
            tax                ; $c5f1: aa        
            pla                ; $c5f2: 68        
            rti                ; $c5f3: 40        

;-------------------------------------------------------------------------------
; irq/brk vector
;-------------------------------------------------------------------------------
irq:        rti                ; $c5f4: 40        

;-------------------------------------------------------------------------------
__start__c5f5:     ldx #$00    ; $c5f5: a2 00       Load 0x00 into X
            stx $00            ; $c5f7: 86 00       Save X into $0000  (RAM)       
            stx $10            ; $c5f9: 86 10       Save X into $0010  (RAM)
            stx $11            ; $c5fb: 86 11       Save X into $0011  (RAM)
            jsr testSetCarry__c72d         ; $c5fd: 20 2d c7  
            jsr __c7db         ; $c600: 20 db c7  
            jsr __c885         ; $c603: 20 85 c8  
            jsr __cbde         ; $c606: 20 de cb  
            jsr __cdf8         ; $c609: 20 f8 cd  
            jsr __ceee         ; $c60c: 20 ee ce  
            jsr __cfa2         ; $c60f: 20 a2 cf  
            jsr __d174         ; $c612: 20 74 d1  
            jsr __d4fb         ; $c615: 20 fb d4  
            jsr __d900         ; $c618: 20 00 d9  
            lda $00            ; $c61b: a5 00     
            sta $10            ; $c61d: 85 10     
            lda #$00           ; $c61f: a9 00     
            sta $00            ; $c621: 85 00     
            jsr __dae0         ; $c623: 20 e0 da  
            jsr __df4a         ; $c626: 20 4a df  
            jsr __dbb8         ; $c629: 20 b8 db  
            jsr __e1aa         ; $c62c: 20 aa e1  
            jsr __c6a3         ; $c62f: 20 a3 c6  
            jsr __e51e         ; $c632: 20 1e e5  
            jsr __e73d         ; $c635: 20 3d e7  
            jsr __e8d3         ; $c638: 20 d3 e8  
            jsr __e916         ; $c63b: 20 16 e9  
            jsr __eb86         ; $c63e: 20 86 eb  
            jsr __edf6         ; $c641: 20 f6 ed  
            jsr __f066         ; $c644: 20 66 f0  
            jsr __f2d6         ; $c647: 20 d6 f2  
            lda $00            ; $c64a: a5 00     
            sta $11            ; $c64c: 85 11     
            lda #$00           ; $c64e: a9 00     
            sta $00            ; $c650: 85 00     
            jsr __f546         ; $c652: 20 46 f5  
            lda $00            ; $c655: a5 00     
            ora $10            ; $c657: 05 10     
            ora $11            ; $c659: 05 11     
            beq __c66b         ; $c65b: f0 0e     
            jsr __c66f         ; $c65d: 20 6f c6  
            ldx $00            ; $c660: a6 00     
            stx $02            ; $c662: 86 02     
            ldx $10            ; $c664: a6 10     
            stx $03            ; $c666: 86 03     
            jmp __c66e         ; $c668: 4c 6e c6  

;-------------------------------------------------------------------------------
__c66b:     jsr __c689         ; $c66b: 20 89 c6  
__c66e:     rts                ; $c66e: 60        

;-------------------------------------------------------------------------------
__c66f:     lda #$03           ; $c66f: a9 03     
            sta SND_CHN        ; $c671: 8d 15 40  
            lda #$87           ; $c674: a9 87     
            sta SQ1_VOL        ; $c676: 8d 00 40  
            lda #$89           ; $c679: a9 89     
            sta SQ1_SWEEP      ; $c67b: 8d 01 40  
            lda #$f0           ; $c67e: a9 f0     
            sta SQ1_LO         ; $c680: 8d 02 40  
            lda #$00           ; $c683: a9 00     
            sta SQ1_HI         ; $c685: 8d 03 40  
            rts                ; $c688: 60        

;-------------------------------------------------------------------------------
__c689:     lda #$02           ; $c689: a9 02     
            sta SND_CHN        ; $c68b: 8d 15 40  
            lda #$3f           ; $c68e: a9 3f     
            sta SQ2_VOL        ; $c690: 8d 04 40  
            lda #$9a           ; $c693: a9 9a     
            sta SQ2_SWEEP      ; $c695: 8d 05 40  
            lda #$ff           ; $c698: a9 ff     
            sta SQ2_LO         ; $c69a: 8d 06 40  
            lda #$00           ; $c69d: a9 00     
            sta SQ2_HI         ; $c69f: 8d 07 40  
            rts                ; $c6a2: 60        

;-------------------------------------------------------------------------------
__c6a3:     ldy #$4e           ; $c6a3: a0 4e     
            lda #$ff           ; $c6a5: a9 ff     
            sta $01            ; $c6a7: 85 01     
            jsr __c6b0         ; $c6a9: 20 b0 c6  
            jsr __c6b7         ; $c6ac: 20 b7 c6  
            rts                ; $c6af: 60        

;-------------------------------------------------------------------------------
__c6b0:     lda #$ff           ; $c6b0: a9 ff     
            pha                ; $c6b2: 48        
            lda #$aa           ; $c6b3: a9 aa     
            bne __c6bc         ; $c6b5: d0 05     
__c6b7:     lda #$34           ; $c6b7: a9 34     
            pha                ; $c6b9: 48        
            lda #$55           ; $c6ba: a9 55     
__c6bc:     plp                ; $c6bc: 28        
            .hex 04 a9         ; $c6bd: 04 a9     Invalid Opcode - NOP $a9
            .hex 44 a9         ; $c6bf: 44 a9     Invalid Opcode - NOP $a9
            .hex 64 a9         ; $c6c1: 64 a9     Invalid Opcode - NOP $a9
            nop                ; $c6c3: ea        
            nop                ; $c6c4: ea        
            nop                ; $c6c5: ea        
__c6c6:     nop                ; $c6c6: ea        
            php                ; $c6c7: 08        
            pha                ; $c6c8: 48        
            .hex 0c a9 a9      ; $c6c9: 0c a9 a9  Invalid Opcode - NOP $a9a9
            nop                ; $c6cc: ea        
            nop                ; $c6cd: ea        
            nop                ; $c6ce: ea        
            nop                ; $c6cf: ea        
            php                ; $c6d0: 08        
            pha                ; $c6d1: 48        
            .hex 14 a9         ; $c6d2: 14 a9     Invalid Opcode - NOP $a9,x
            .hex 34 a9         ; $c6d4: 34 a9     Invalid Opcode - NOP $a9,x
            .hex 54 a9         ; $c6d6: 54 a9     Invalid Opcode - NOP $a9,x
            .hex 74 a9         ; $c6d8: 74 a9     Invalid Opcode - NOP $a9,x
            .hex d4 a9         ; $c6da: d4 a9     Invalid Opcode - NOP $a9,x
            .hex f4 a9         ; $c6dc: f4 a9     Invalid Opcode - NOP $a9,x
            nop                ; $c6de: ea        
            nop                ; $c6df: ea        
            nop                ; $c6e0: ea        
            nop                ; $c6e1: ea        
            php                ; $c6e2: 08        
            pha                ; $c6e3: 48        
            .hex 1a            ; $c6e4: 1a        Invalid Opcode - NOP 
            .hex 3a            ; $c6e5: 3a        Invalid Opcode - NOP 
            .hex 5a            ; $c6e6: 5a        Invalid Opcode - NOP 
            .hex 7a            ; $c6e7: 7a        Invalid Opcode - NOP 
            .hex da            ; $c6e8: da        Invalid Opcode - NOP 
            .hex fa            ; $c6e9: fa        Invalid Opcode - NOP 
            .hex 80 89         ; $c6ea: 80 89     Invalid Opcode - NOP #$89
            nop                ; $c6ec: ea        
            nop                ; $c6ed: ea        
            nop                ; $c6ee: ea        
            nop                ; $c6ef: ea        
            php                ; $c6f0: 08        
            pha                ; $c6f1: 48        
            .hex 1c a9 a9      ; $c6f2: 1c a9 a9  Invalid Opcode - NOP $a9a9,x
            .hex 3c a9 a9      ; $c6f5: 3c a9 a9  Invalid Opcode - NOP $a9a9,x
            .hex 5c a9 a9      ; $c6f8: 5c a9 a9  Invalid Opcode - NOP $a9a9,x
            .hex 7c a9 a9      ; $c6fb: 7c a9 a9  Invalid Opcode - NOP $a9a9,x
            .hex dc a9 a9      ; $c6fe: dc a9 a9  Invalid Opcode - NOP $a9a9,x
            .hex fc a9 a9      ; $c701: fc a9 a9  Invalid Opcode - NOP $a9a9,x
            nop                ; $c704: ea        
            nop                ; $c705: ea        
            nop                ; $c706: ea        
            nop                ; $c707: ea        
            php                ; $c708: 08        
            pha                ; $c709: 48        
            ldx #$05           ; $c70a: a2 05     
__c70c:     pla                ; $c70c: 68        
            cmp #$55           ; $c70d: c9 55     
            beq __c71b         ; $c70f: f0 0a     
            cmp #$aa           ; $c711: c9 aa     
            beq __c71b         ; $c713: f0 06     
            pla                ; $c715: 68        
            sty $00            ; $c716: 84 00     
            jmp __c728         ; $c718: 4c 28 c7  

;-------------------------------------------------------------------------------
__c71b:     pla                ; $c71b: 68        
            and #$cb           ; $c71c: 29 cb     
            cmp #$00           ; $c71e: c9 00     
            beq __c728         ; $c720: f0 06     
            cmp #$cb           ; $c722: c9 cb     
            beq __c728         ; $c724: f0 02     
            sty $00            ; $c726: 84 00     
__c728:     iny                ; $c728: c8        
            dex                ; $c729: ca        
            bne __c70c         ; $c72a: d0 e0     
            rts                ; $c72c: 60        

;-------------------------------------------------------------------------------
testSetCarry__c72d:     nop                              ; $c72d: ea          wait
                        sec                              ; $c72e: 38          set carry flag  
                        bcs testUnsetCarry__c735         ; $c72f: b0 04       if set correctly, test unset
                        ldx #$01                         ; $c731: a2 01       .
                        stx $00                          ; $c733: 86 00       Store "1" at $0x000 if carry was not set
testUnsetCarry__c735:   nop                              ; $c735: ea          wait
                        clc                              ; $c736: 18          clears carry flag
                        bcs __c73c                       ; $c737: b0 03       If carry flag still set, save failed result.
                        jmp __c740                       ; $c739: 4c 40 c7    Else go next test

;-------------------------------------------------------------------------------
; testUnsetCarry failed
__c73c:     ldx #$02           ; $c73c: a2 02       .
            stx $00            ; $c73e: 86 00       Store "2" at $0x000 if carry was not unset
__c740:     nop                ; $c740: ea        
            sec                ; $c741: 38          Set carry flag
            bcc __c747         ; $c742: 90 03       If carry flag is clear, save failed result: bccSawClear
            jmp __c74b         ; $c744: 4c 4b c7    Else go next test

;-------------------------------------------------------------------------------
; bccSawClear
__c747:     ldx #$03           ; $c747: a2 03       .
            stx $00            ; $c749: 86 00       Store "3" at $0x000 if BCC failed to see carry clear
__c74b:     nop                ; $c74b: ea        
            clc                ; $c74c: 18          clears carry flag
            bcc __c753         ; $c74d: 90 04       If carry flag is clear, go to
            ldx #$04           ; $c74f: a2 04       Else
            stx $00            ; $c751: 86 00           Store "4" at $0x000
__c753:     nop                ; $c753: ea        
            lda #$00           ; $c754: a9 00       Load test result into A
            beq __c75c         ; $c756: f0 04       If test result == 0 go to __c75c
            ldx #$05           ; $c758: a2 05     
            stx $00            ; $c75a: 86 00     
__c75c:     nop                ; $c75c: ea        
            lda #$40           ; $c75d: a9 40       Load "0x40" into A
            beq __c764         ; $c75f: f0 03       If ZeroFlag == 0 go to __c764
            jmp testBNE         ; $c761: 4c 68 c7   Else go testBNE

;-------------------------------------------------------------------------------
__c764:     ldx #$06           ; $c764: a2 06     
            stx $00            ; $c766: 86 00     
testBNE:    nop                ; $c768: ea        
            lda #$40           ; $c769: a9 40       Load "0x40" into A
            bne __c771         ; $c76b: d0 04       If ZeroFlag != 0 go to next test
            ldx #$07           ; $c76d: a2 07       Else
            stx $00            ; $c76f: 86 00           Store "0x07" at $0x0000
__c771:     nop                ; $c771: ea          
            lda #$00           ; $c772: a9 00       Load "0x00" into A
            bne __c779         ; $c774: d0 03       If ZeroFlag != 0 go save failed test
            jmp testBVS        ; $c776: 4c 7d c7    Else, test succeed. go next test

;-------------------------------------------------------------------------------
__c779:     ldx #$08           ; $c779: a2 08       .
            stx $00            ; $c77b: 86 00       Store "0x08" at $0x0000
testBVS:    nop                ; $c77d: ea          
            lda #$ff           ; $c77e: a9 ff       Load "0xFF" into A
            sta $01            ; $c780: 85 01       Write A at $0x0001
            bit $01            ; $c782: 24 01       Checks bits 7 and 6 between A and $0x0001
            bvs testBVC        ; $c784: 70 04       If Overflow is set, go to testBVC
            ldx #$09           ; $c786: a2 09       Else
            stx $00            ; $c788: 86 00           Store "0x09" at $0x0000 as failed test
testBVC:    nop                ; $c78a: ea        
            bit $01            ; $c78b: 24 01       Checks bit 7 and 6 between A and $0x0001
            bvc testBVCfailed  ; $c78d: 50 03       If Overflow is clear, go to testBVCfailed
            jmp __c796         ; $c78f: 4c 96 c7    Else, test succeed. go next test

;-------------------------------------------------------------------------------
testBVCfailed:     ldx #$0a           ; $c792: a2 0a     
                   stx $00            ; $c794: 86 00     
__c796:            nop                ; $c796: ea          
                   lda #$00           ; $c797: a9 00       Load "0x00" into A
                   sta $01            ; $c799: 85 01       Write A into $0x0001
                   bit $01            ; $c79b: 24 01       Checks bits 7 and 6 between A and $0x0001
                   bvc __c7a3         ; $c79d: 50 04       If Overflow is clear, go to 
                   ldx #$0b           ; $c79f: a2 0b     
                   stx $00            ; $c7a1: 86 00     
__c7a3:            nop                ; $c7a3: ea        
                   bit $01            ; $c7a4: 24 01       Checks bits 7 and 6 between A and $0x0001
                   bvs __c7ab         ; $c7a6: 70 03       If Overflow is set, go to __c7ab
                   jmp __c7af         ; $c7a8: 4c af c7    Else, test succeed. go next test

;-------------------------------------------------------------------------------
__c7ab:     ldx #$0c           ; $c7ab: a2 0c     
            stx $00            ; $c7ad: 86 00     
__c7af:     nop                ; $c7af: ea        
            lda #$00           ; $c7b0: a9 00     
            bpl __c7b8         ; $c7b2: 10 04     
            ldx #$0d           ; $c7b4: a2 0d     
            stx $00            ; $c7b6: 86 00     
__c7b8:     nop                ; $c7b8: ea        
            lda #$80           ; $c7b9: a9 80     
            bpl __c7c0         ; $c7bb: 10 03     
            jmp __c7d9         ; $c7bd: 4c d9 c7  

;-------------------------------------------------------------------------------
__c7c0:     ldx #$0e           ; $c7c0: a2 0e     
            stx $00            ; $c7c2: 86 00     
            nop                ; $c7c4: ea        
            lda #$80           ; $c7c5: a9 80     
            bmi __c7cd         ; $c7c7: 30 04     
            ldx #$0f           ; $c7c9: a2 0f     
            stx $00            ; $c7cb: 86 00     
__c7cd:     nop                ; $c7cd: ea        
            lda #$00           ; $c7ce: a9 00     
            bmi __c7d5         ; $c7d0: 30 03     
            jmp __c7d9         ; $c7d2: 4c d9 c7  

;-------------------------------------------------------------------------------
__c7d5:     ldx #$10           ; $c7d5: a2 10     
            stx $00            ; $c7d7: 86 00     
__c7d9:     nop                ; $c7d9: ea        
            rts                ; $c7da: 60        

;-------------------------------------------------------------------------------
__c7db:     nop                ; $c7db: ea        
            lda #$ff           ; $c7dc: a9 ff     
            sta $01            ; $c7de: 85 01     
            bit $01            ; $c7e0: 24 01     
            lda #$00           ; $c7e2: a9 00     
            sec                ; $c7e4: 38        
            sei                ; $c7e5: 78        
            sed                ; $c7e6: f8        
            php                ; $c7e7: 08        
            pla                ; $c7e8: 68        
            and #$ef           ; $c7e9: 29 ef     
            cmp #$6f           ; $c7eb: c9 6f     
            beq __c7f3         ; $c7ed: f0 04     
            ldx #$11           ; $c7ef: a2 11     
            stx $00            ; $c7f1: 86 00     
__c7f3:     nop                ; $c7f3: ea        
            lda #$40           ; $c7f4: a9 40     
            sta $01            ; $c7f6: 85 01     
            bit $01            ; $c7f8: 24 01     
            cld                ; $c7fa: d8        
            lda #$10           ; $c7fb: a9 10     
            clc                ; $c7fd: 18        
            php                ; $c7fe: 08        
            pla                ; $c7ff: 68        
            and #$ef           ; $c800: 29 ef     
            cmp #$64           ; $c802: c9 64     
            beq __c80a         ; $c804: f0 04     
            ldx #$12           ; $c806: a2 12     
            stx $00            ; $c808: 86 00     
__c80a:     nop                ; $c80a: ea        
            lda #$80           ; $c80b: a9 80     
            sta $01            ; $c80d: 85 01     
            bit $01            ; $c80f: 24 01     
            sed                ; $c811: f8        
            lda #$00           ; $c812: a9 00     
            sec                ; $c814: 38        
            php                ; $c815: 08        
            pla                ; $c816: 68        
            and #$ef           ; $c817: 29 ef     
            cmp #$2f           ; $c819: c9 2f     
            beq __c821         ; $c81b: f0 04     
            ldx #$13           ; $c81d: a2 13     
            stx $00            ; $c81f: 86 00     
__c821:     nop                ; $c821: ea        
            lda #$ff           ; $c822: a9 ff     
            pha                ; $c824: 48        
            plp                ; $c825: 28        
            bne __c831         ; $c826: d0 09     
            bpl __c831         ; $c828: 10 07     
            bvc __c831         ; $c82a: 50 05     
            bcc __c831         ; $c82c: 90 03     
            jmp __c835         ; $c82e: 4c 35 c8  

;-------------------------------------------------------------------------------
__c831:     ldx #$14           ; $c831: a2 14     
            stx $00            ; $c833: 86 00     
__c835:     nop                ; $c835: ea        
            lda #$04           ; $c836: a9 04     
            pha                ; $c838: 48        
            plp                ; $c839: 28        
            beq __c845         ; $c83a: f0 09     
            bmi __c845         ; $c83c: 30 07     
            bvs __c845         ; $c83e: 70 05     
            bcs __c845         ; $c840: b0 03     
            jmp __c849         ; $c842: 4c 49 c8  

;-------------------------------------------------------------------------------
__c845:     ldx #$15           ; $c845: a2 15     
            stx $00            ; $c847: 86 00     
__c849:     nop                ; $c849: ea        
            sed                ; $c84a: f8        
            lda #$ff           ; $c84b: a9 ff     
            sta $01            ; $c84d: 85 01     
            bit $01            ; $c84f: 24 01     
            clc                ; $c851: 18        
            lda #$00           ; $c852: a9 00     
            pha                ; $c854: 48        
            lda #$ff           ; $c855: a9 ff     
            pla                ; $c857: 68        
            bne __c863         ; $c858: d0 09     
            bmi __c863         ; $c85a: 30 07     
            bvc __c863         ; $c85c: 50 05     
            bcs __c863         ; $c85e: b0 03     
            jmp __c867         ; $c860: 4c 67 c8  

;-------------------------------------------------------------------------------
__c863:     ldx #$16           ; $c863: a2 16     
            stx $00            ; $c865: 86 00     
__c867:     nop                ; $c867: ea        
            lda #$00           ; $c868: a9 00     
            sta $01            ; $c86a: 85 01     
            bit $01            ; $c86c: 24 01     
            sec                ; $c86e: 38        
            lda #$ff           ; $c86f: a9 ff     
            pha                ; $c871: 48        
            lda #$00           ; $c872: a9 00     
            pla                ; $c874: 68        
            beq __c880         ; $c875: f0 09     
            bpl __c880         ; $c877: 10 07     
            bvs __c880         ; $c879: 70 05     
            bcc __c880         ; $c87b: 90 03     
            jmp __c884         ; $c87d: 4c 84 c8  

;-------------------------------------------------------------------------------
__c880:     ldx #$17           ; $c880: a2 17     
            stx $00            ; $c882: 86 00     
__c884:     rts                ; $c884: 60        

;-------------------------------------------------------------------------------
__c885:     nop                ; $c885: ea        
            clc                ; $c886: 18        
            lda #$ff           ; $c887: a9 ff     
            sta $01            ; $c889: 85 01     
            bit $01            ; $c88b: 24 01     
            lda #$55           ; $c88d: a9 55     
            ora #$aa           ; $c88f: 09 aa     
            bcs __c89e         ; $c891: b0 0b     
            bpl __c89e         ; $c893: 10 09     
            cmp #$ff           ; $c895: c9 ff     
            bne __c89e         ; $c897: d0 05     
            bvc __c89e         ; $c899: 50 03     
            jmp __c8a2         ; $c89b: 4c a2 c8  

;-------------------------------------------------------------------------------
__c89e:     ldx #$18           ; $c89e: a2 18     
            stx $00            ; $c8a0: 86 00     
__c8a2:     nop                ; $c8a2: ea        
            sec                ; $c8a3: 38        
            clv                ; $c8a4: b8        
            lda #$00           ; $c8a5: a9 00     
            ora #$00           ; $c8a7: 09 00     
            bne __c8b4         ; $c8a9: d0 09     
            bvs __c8b4         ; $c8ab: 70 07     
            bcc __c8b4         ; $c8ad: 90 05     
            bmi __c8b4         ; $c8af: 30 03     
            jmp __c8b8         ; $c8b1: 4c b8 c8  

;-------------------------------------------------------------------------------
__c8b4:     ldx #$19           ; $c8b4: a2 19     
            stx $00            ; $c8b6: 86 00     
__c8b8:     nop                ; $c8b8: ea        
            clc                ; $c8b9: 18        
            bit $01            ; $c8ba: 24 01     
            lda #$55           ; $c8bc: a9 55     
            and #$aa           ; $c8be: 29 aa     
            bne __c8cb         ; $c8c0: d0 09     
            bvc __c8cb         ; $c8c2: 50 07     
            bcs __c8cb         ; $c8c4: b0 05     
            bmi __c8cb         ; $c8c6: 30 03     
            jmp __c8cf         ; $c8c8: 4c cf c8  

;-------------------------------------------------------------------------------
__c8cb:     ldx #$1a           ; $c8cb: a2 1a     
            stx $00            ; $c8cd: 86 00     
__c8cf:     nop                ; $c8cf: ea        
            sec                ; $c8d0: 38        
            clv                ; $c8d1: b8        
            lda #$f8           ; $c8d2: a9 f8     
            and #$ef           ; $c8d4: 29 ef     
            bcc __c8e3         ; $c8d6: 90 0b     
            bpl __c8e3         ; $c8d8: 10 09     
            cmp #$e8           ; $c8da: c9 e8     
            bne __c8e3         ; $c8dc: d0 05     
            bvs __c8e3         ; $c8de: 70 03     
            jmp __c8e7         ; $c8e0: 4c e7 c8  

;-------------------------------------------------------------------------------
__c8e3:     ldx #$1b           ; $c8e3: a2 1b     
            stx $00            ; $c8e5: 86 00     
__c8e7:     nop                ; $c8e7: ea        
            clc                ; $c8e8: 18        
            bit $01            ; $c8e9: 24 01     
            lda #$5f           ; $c8eb: a9 5f     
            eor #$aa           ; $c8ed: 49 aa     
            bcs __c8fc         ; $c8ef: b0 0b     
            bpl __c8fc         ; $c8f1: 10 09     
            cmp #$f5           ; $c8f3: c9 f5     
            bne __c8fc         ; $c8f5: d0 05     
            bvc __c8fc         ; $c8f7: 50 03     
            jmp __c900         ; $c8f9: 4c 00 c9  

;-------------------------------------------------------------------------------
__c8fc:     ldx #$1c           ; $c8fc: a2 1c     
            stx $00            ; $c8fe: 86 00     
__c900:     nop                ; $c900: ea        
            sec                ; $c901: 38        
            clv                ; $c902: b8        
            lda #$70           ; $c903: a9 70     
            eor #$70           ; $c905: 49 70     
            bne __c912         ; $c907: d0 09     
            bvs __c912         ; $c909: 70 07     
            bcc __c912         ; $c90b: 90 05     
            bmi __c912         ; $c90d: 30 03     
            jmp __c916         ; $c90f: 4c 16 c9  

;-------------------------------------------------------------------------------
__c912:     ldx #$1d           ; $c912: a2 1d     
            stx $00            ; $c914: 86 00     
__c916:     nop                ; $c916: ea        
            clc                ; $c917: 18        
            bit $01            ; $c918: 24 01     
            lda #$00           ; $c91a: a9 00     
            adc #$69           ; $c91c: 69 69     
            bmi __c92b         ; $c91e: 30 0b     
            bcs __c92b         ; $c920: b0 09     
            cmp #$69           ; $c922: c9 69     
            bne __c92b         ; $c924: d0 05     
            bvs __c92b         ; $c926: 70 03     
            jmp __c92f         ; $c928: 4c 2f c9  

;-------------------------------------------------------------------------------
__c92b:     ldx #$1e           ; $c92b: a2 1e     
            stx $00            ; $c92d: 86 00     
__c92f:     nop                ; $c92f: ea        
            sec                ; $c930: 38        
            sed                ; $c931: f8        
            bit $01            ; $c932: 24 01     
            lda #$01           ; $c934: a9 01     
            adc #$69           ; $c936: 69 69     
            bmi __c945         ; $c938: 30 0b     
            bcs __c945         ; $c93a: b0 09     
            cmp #$6b           ; $c93c: c9 6b     
            bne __c945         ; $c93e: d0 05     
            bvs __c945         ; $c940: 70 03     
            jmp __c949         ; $c942: 4c 49 c9  

;-------------------------------------------------------------------------------
__c945:     ldx #$1f           ; $c945: a2 1f     
            stx $00            ; $c947: 86 00     
__c949:     nop                ; $c949: ea        
            cld                ; $c94a: d8        
            sec                ; $c94b: 38        
            clv                ; $c94c: b8        
            lda #$7f           ; $c94d: a9 7f     
            adc #$7f           ; $c94f: 69 7f     
            bpl __c95e         ; $c951: 10 0b     
            bcs __c95e         ; $c953: b0 09     
            cmp #$ff           ; $c955: c9 ff     
            bne __c95e         ; $c957: d0 05     
            bvc __c95e         ; $c959: 50 03     
            jmp __c962         ; $c95b: 4c 62 c9  

;-------------------------------------------------------------------------------
__c95e:     ldx #$20           ; $c95e: a2 20     
            stx $00            ; $c960: 86 00     
__c962:     nop                ; $c962: ea        
            clc                ; $c963: 18        
            bit $01            ; $c964: 24 01     
            lda #$7f           ; $c966: a9 7f     
            adc #$80           ; $c968: 69 80     
            bpl __c977         ; $c96a: 10 0b     
            bcs __c977         ; $c96c: b0 09     
            cmp #$ff           ; $c96e: c9 ff     
            bne __c977         ; $c970: d0 05     
            bvs __c977         ; $c972: 70 03     
            jmp __c97b         ; $c974: 4c 7b c9  

;-------------------------------------------------------------------------------
__c977:     ldx #$21           ; $c977: a2 21     
            stx $00            ; $c979: 86 00     
__c97b:     nop                ; $c97b: ea        
            sec                ; $c97c: 38        
            clv                ; $c97d: b8        
            lda #$7f           ; $c97e: a9 7f     
            adc #$80           ; $c980: 69 80     
            bne __c98d         ; $c982: d0 09     
            bmi __c98d         ; $c984: 30 07     
            bvs __c98d         ; $c986: 70 05     
            bcc __c98d         ; $c988: 90 03     
            jmp __c991         ; $c98a: 4c 91 c9  

;-------------------------------------------------------------------------------
__c98d:     ldx #$22           ; $c98d: a2 22     
            stx $00            ; $c98f: 86 00     
__c991:     nop                ; $c991: ea        
            sec                ; $c992: 38        
            clv                ; $c993: b8        
            lda #$9f           ; $c994: a9 9f     
            beq __c9a1         ; $c996: f0 09     
            bpl __c9a1         ; $c998: 10 07     
            bvs __c9a1         ; $c99a: 70 05     
            bcc __c9a1         ; $c99c: 90 03     
            jmp __c9a5         ; $c99e: 4c a5 c9  

;-------------------------------------------------------------------------------
__c9a1:     ldx #$23           ; $c9a1: a2 23     
            stx $00            ; $c9a3: 86 00     
__c9a5:     nop                ; $c9a5: ea        
            clc                ; $c9a6: 18        
            bit $01            ; $c9a7: 24 01     
            lda #$00           ; $c9a9: a9 00     
            bne __c9b6         ; $c9ab: d0 09     
            bmi __c9b6         ; $c9ad: 30 07     
            bvc __c9b6         ; $c9af: 50 05     
            bcs __c9b6         ; $c9b1: b0 03     
            jmp __c9ba         ; $c9b3: 4c ba c9  

;-------------------------------------------------------------------------------
__c9b6:     ldx #$23           ; $c9b6: a2 23     
            stx $00            ; $c9b8: 86 00     
__c9ba:     nop                ; $c9ba: ea        
            bit $01            ; $c9bb: 24 01     
            lda #$40           ; $c9bd: a9 40     
            cmp #$40           ; $c9bf: c9 40     
            bmi __c9cc         ; $c9c1: 30 09     
            bcc __c9cc         ; $c9c3: 90 07     
            bne __c9cc         ; $c9c5: d0 05     
            bvc __c9cc         ; $c9c7: 50 03     
            jmp __c9d0         ; $c9c9: 4c d0 c9  

;-------------------------------------------------------------------------------
__c9cc:     ldx #$24           ; $c9cc: a2 24     
            stx $00            ; $c9ce: 86 00     
__c9d0:     nop                ; $c9d0: ea        
            clv                ; $c9d1: b8        
            cmp #$3f           ; $c9d2: c9 3f     
            beq __c9df         ; $c9d4: f0 09     
            bmi __c9df         ; $c9d6: 30 07     
            bcc __c9df         ; $c9d8: 90 05     
            bvs __c9df         ; $c9da: 70 03     
            jmp __c9e3         ; $c9dc: 4c e3 c9  

;-------------------------------------------------------------------------------
__c9df:     ldx #$25           ; $c9df: a2 25     
            stx $00            ; $c9e1: 86 00     
__c9e3:     nop                ; $c9e3: ea        
            cmp #$41           ; $c9e4: c9 41     
            beq __c9ef         ; $c9e6: f0 07     
            bpl __c9ef         ; $c9e8: 10 05     
            bpl __c9ef         ; $c9ea: 10 03     
            jmp __c9f3         ; $c9ec: 4c f3 c9  

;-------------------------------------------------------------------------------
__c9ef:     ldx #$26           ; $c9ef: a2 26     
            stx $00            ; $c9f1: 86 00     
__c9f3:     nop                ; $c9f3: ea        
            lda #$80           ; $c9f4: a9 80     
            cmp #$00           ; $c9f6: c9 00     
            beq __ca01         ; $c9f8: f0 07     
            bpl __ca01         ; $c9fa: 10 05     
            bcc __ca01         ; $c9fc: 90 03     
            jmp __ca05         ; $c9fe: 4c 05 ca  

;-------------------------------------------------------------------------------
__ca01:     ldx #$27           ; $ca01: a2 27     
            stx $00            ; $ca03: 86 00     
__ca05:     nop                ; $ca05: ea        
            cmp #$80           ; $ca06: c9 80     
            bne __ca11         ; $ca08: d0 07     
            bmi __ca11         ; $ca0a: 30 05     
            bcc __ca11         ; $ca0c: 90 03     
            jmp __ca15         ; $ca0e: 4c 15 ca  

;-------------------------------------------------------------------------------
__ca11:     ldx #$28           ; $ca11: a2 28     
            stx $00            ; $ca13: 86 00     
__ca15:     nop                ; $ca15: ea        
            cmp #$81           ; $ca16: c9 81     
            bcs __ca21         ; $ca18: b0 07     
            beq __ca21         ; $ca1a: f0 05     
            bpl __ca21         ; $ca1c: 10 03     
            jmp __ca25         ; $ca1e: 4c 25 ca  

;-------------------------------------------------------------------------------
__ca21:     ldx #$29           ; $ca21: a2 29     
            stx $00            ; $ca23: 86 00     
__ca25:     nop                ; $ca25: ea        
            cmp #$7f           ; $ca26: c9 7f     
            bcc __ca31         ; $ca28: 90 07     
            beq __ca31         ; $ca2a: f0 05     
            bmi __ca31         ; $ca2c: 30 03     
            jmp __ca35         ; $ca2e: 4c 35 ca  

;-------------------------------------------------------------------------------
__ca31:     ldx #$2a           ; $ca31: a2 2a     
            stx $00            ; $ca33: 86 00     
__ca35:     nop                ; $ca35: ea        
            bit $01            ; $ca36: 24 01     
            ldy #$40           ; $ca38: a0 40     
            cpy #$40           ; $ca3a: c0 40     
            bne __ca47         ; $ca3c: d0 09     
            bmi __ca47         ; $ca3e: 30 07     
            bcc __ca47         ; $ca40: 90 05     
            bvc __ca47         ; $ca42: 50 03     
            jmp __ca4b         ; $ca44: 4c 4b ca  

;-------------------------------------------------------------------------------
__ca47:     ldx #$2b           ; $ca47: a2 2b     
            stx $00            ; $ca49: 86 00     
__ca4b:     nop                ; $ca4b: ea        
            clv                ; $ca4c: b8        
            cpy #$3f           ; $ca4d: c0 3f     
            beq __ca5a         ; $ca4f: f0 09     
            bmi __ca5a         ; $ca51: 30 07     
            bcc __ca5a         ; $ca53: 90 05     
            bvs __ca5a         ; $ca55: 70 03     
            jmp __ca5e         ; $ca57: 4c 5e ca  

;-------------------------------------------------------------------------------
__ca5a:     ldx #$2c           ; $ca5a: a2 2c     
            stx $00            ; $ca5c: 86 00     
__ca5e:     nop                ; $ca5e: ea        
            cpy #$41           ; $ca5f: c0 41     
            beq __ca6a         ; $ca61: f0 07     
            bpl __ca6a         ; $ca63: 10 05     
            bpl __ca6a         ; $ca65: 10 03     
            jmp __ca6e         ; $ca67: 4c 6e ca  

;-------------------------------------------------------------------------------
__ca6a:     ldx #$2d           ; $ca6a: a2 2d     
            stx $00            ; $ca6c: 86 00     
__ca6e:     nop                ; $ca6e: ea        
            ldy #$80           ; $ca6f: a0 80     
            cpy #$00           ; $ca71: c0 00     
            beq __ca7c         ; $ca73: f0 07     
            bpl __ca7c         ; $ca75: 10 05     
            bcc __ca7c         ; $ca77: 90 03     
            jmp __ca80         ; $ca79: 4c 80 ca  

;-------------------------------------------------------------------------------
__ca7c:     ldx #$2e           ; $ca7c: a2 2e     
            stx $00            ; $ca7e: 86 00     
__ca80:     nop                ; $ca80: ea        
            cpy #$80           ; $ca81: c0 80     
            bne __ca8c         ; $ca83: d0 07     
            bmi __ca8c         ; $ca85: 30 05     
            bcc __ca8c         ; $ca87: 90 03     
            jmp __ca90         ; $ca89: 4c 90 ca  

;-------------------------------------------------------------------------------
__ca8c:     ldx #$2f           ; $ca8c: a2 2f     
            stx $00            ; $ca8e: 86 00     
__ca90:     nop                ; $ca90: ea        
            cpy #$81           ; $ca91: c0 81     
            bcs __ca9c         ; $ca93: b0 07     
            beq __ca9c         ; $ca95: f0 05     
            bpl __ca9c         ; $ca97: 10 03     
            jmp __caa0         ; $ca99: 4c a0 ca  

;-------------------------------------------------------------------------------
__ca9c:     ldx #$30           ; $ca9c: a2 30     
            stx $00            ; $ca9e: 86 00     
__caa0:     nop                ; $caa0: ea        
            cpy #$7f           ; $caa1: c0 7f     
            bcc __caac         ; $caa3: 90 07     
            beq __caac         ; $caa5: f0 05     
            bmi __caac         ; $caa7: 30 03     
            jmp __cab0         ; $caa9: 4c b0 ca  

;-------------------------------------------------------------------------------
__caac:     ldx #$31           ; $caac: a2 31     
            stx $00            ; $caae: 86 00     
__cab0:     nop                ; $cab0: ea        
            bit $01            ; $cab1: 24 01     
            ldx #$40           ; $cab3: a2 40     
            cpx #$40           ; $cab5: e0 40     
            bne __cac2         ; $cab7: d0 09     
            bmi __cac2         ; $cab9: 30 07     
            bcc __cac2         ; $cabb: 90 05     
            bvc __cac2         ; $cabd: 50 03     
            jmp __cac6         ; $cabf: 4c c6 ca  

;-------------------------------------------------------------------------------
__cac2:     lda #$32           ; $cac2: a9 32     
            sta $00            ; $cac4: 85 00     
__cac6:     nop                ; $cac6: ea        
            clv                ; $cac7: b8        
            cpx #$3f           ; $cac8: e0 3f     
            beq __cad5         ; $caca: f0 09     
            bmi __cad5         ; $cacc: 30 07     
            bcc __cad5         ; $cace: 90 05     
            bvs __cad5         ; $cad0: 70 03     
            jmp __cad9         ; $cad2: 4c d9 ca  

;-------------------------------------------------------------------------------
__cad5:     lda #$33           ; $cad5: a9 33     
            sta $00            ; $cad7: 85 00     
__cad9:     nop                ; $cad9: ea        
            cpx #$41           ; $cada: e0 41     
            beq __cae5         ; $cadc: f0 07     
            bpl __cae5         ; $cade: 10 05     
            bpl __cae5         ; $cae0: 10 03     
            jmp __cae9         ; $cae2: 4c e9 ca  

;-------------------------------------------------------------------------------
__cae5:     lda #$34           ; $cae5: a9 34     
            sta $00            ; $cae7: 85 00     
__cae9:     nop                ; $cae9: ea        
            ldx #$80           ; $caea: a2 80     
            cpx #$00           ; $caec: e0 00     
            beq __caf7         ; $caee: f0 07     
            bpl __caf7         ; $caf0: 10 05     
            bcc __caf7         ; $caf2: 90 03     
            jmp __cafb         ; $caf4: 4c fb ca  

;-------------------------------------------------------------------------------
__caf7:     lda #$35           ; $caf7: a9 35     
            sta $00            ; $caf9: 85 00     
__cafb:     nop                ; $cafb: ea        
            cpx #$80           ; $cafc: e0 80     
            bne __cb07         ; $cafe: d0 07     
            bmi __cb07         ; $cb00: 30 05     
            bcc __cb07         ; $cb02: 90 03     
            jmp __cb0b         ; $cb04: 4c 0b cb  

;-------------------------------------------------------------------------------
__cb07:     lda #$36           ; $cb07: a9 36     
            sta $00            ; $cb09: 85 00     
__cb0b:     nop                ; $cb0b: ea        
            cpx #$81           ; $cb0c: e0 81     
            bcs __cb17         ; $cb0e: b0 07     
            beq __cb17         ; $cb10: f0 05     
            bpl __cb17         ; $cb12: 10 03     
            jmp __cb1b         ; $cb14: 4c 1b cb  

;-------------------------------------------------------------------------------
__cb17:     lda #$37           ; $cb17: a9 37     
            sta $00            ; $cb19: 85 00     
__cb1b:     nop                ; $cb1b: ea        
            cpx #$7f           ; $cb1c: e0 7f     
            bcc __cb27         ; $cb1e: 90 07     
            beq __cb27         ; $cb20: f0 05     
            bmi __cb27         ; $cb22: 30 03     
            jmp __cb2b         ; $cb24: 4c 2b cb  

;-------------------------------------------------------------------------------
__cb27:     lda #$38           ; $cb27: a9 38     
            sta $00            ; $cb29: 85 00     
__cb2b:     nop                ; $cb2b: ea        
            sec                ; $cb2c: 38        
            clv                ; $cb2d: b8        
            ldx #$9f           ; $cb2e: a2 9f     
            beq __cb3b         ; $cb30: f0 09     
            bpl __cb3b         ; $cb32: 10 07     
            bvs __cb3b         ; $cb34: 70 05     
            bcc __cb3b         ; $cb36: 90 03     
            jmp __cb3f         ; $cb38: 4c 3f cb  

;-------------------------------------------------------------------------------
__cb3b:     ldx #$39           ; $cb3b: a2 39     
            stx $00            ; $cb3d: 86 00     
__cb3f:     nop                ; $cb3f: ea        
            clc                ; $cb40: 18        
            bit $01            ; $cb41: 24 01     
            ldx #$00           ; $cb43: a2 00     
            bne __cb50         ; $cb45: d0 09     
            bmi __cb50         ; $cb47: 30 07     
            bvc __cb50         ; $cb49: 50 05     
            bcs __cb50         ; $cb4b: b0 03     
            jmp __cb54         ; $cb4d: 4c 54 cb  

;-------------------------------------------------------------------------------
__cb50:     ldx #$3a           ; $cb50: a2 3a     
            stx $00            ; $cb52: 86 00     
__cb54:     nop                ; $cb54: ea        
            sec                ; $cb55: 38        
            clv                ; $cb56: b8        
            ldy #$9f           ; $cb57: a0 9f     
            beq __cb64         ; $cb59: f0 09     
            bpl __cb64         ; $cb5b: 10 07     
            bvs __cb64         ; $cb5d: 70 05     
            bcc __cb64         ; $cb5f: 90 03     
            jmp __cb68         ; $cb61: 4c 68 cb  

;-------------------------------------------------------------------------------
__cb64:     ldx #$3b           ; $cb64: a2 3b     
            stx $00            ; $cb66: 86 00     
__cb68:     nop                ; $cb68: ea        
            clc                ; $cb69: 18        
            bit $01            ; $cb6a: 24 01     
            ldy #$00           ; $cb6c: a0 00     
            bne __cb79         ; $cb6e: d0 09     
            bmi __cb79         ; $cb70: 30 07     
            bvc __cb79         ; $cb72: 50 05     
            bcs __cb79         ; $cb74: b0 03     
            jmp __cb7d         ; $cb76: 4c 7d cb  

;-------------------------------------------------------------------------------
__cb79:     ldx #$3c           ; $cb79: a2 3c     
            stx $00            ; $cb7b: 86 00     
__cb7d:     nop                ; $cb7d: ea        
            lda #$55           ; $cb7e: a9 55     
            ldx #$aa           ; $cb80: a2 aa     
            ldy #$33           ; $cb82: a0 33     
            cmp #$55           ; $cb84: c9 55     
            bne __cbab         ; $cb86: d0 23     
            cpx #$aa           ; $cb88: e0 aa     
            bne __cbab         ; $cb8a: d0 1f     
            cpy #$33           ; $cb8c: c0 33     
            bne __cbab         ; $cb8e: d0 1b     
            cmp #$55           ; $cb90: c9 55     
            bne __cbab         ; $cb92: d0 17     
            cpx #$aa           ; $cb94: e0 aa     
            bne __cbab         ; $cb96: d0 13     
            cpy #$33           ; $cb98: c0 33     
            bne __cbab         ; $cb9a: d0 0f     
            cmp #$56           ; $cb9c: c9 56     
            beq __cbab         ; $cb9e: f0 0b     
            cpx #$ab           ; $cba0: e0 ab     
            beq __cbab         ; $cba2: f0 07     
            cpy #$34           ; $cba4: c0 34     
            beq __cbab         ; $cba6: f0 03     
            jmp __cbaf         ; $cba8: 4c af cb  

;-------------------------------------------------------------------------------
__cbab:     ldx #$3d           ; $cbab: a2 3d     
            stx $00            ; $cbad: 86 00     
__cbaf:     ldy #$71           ; $cbaf: a0 71     
            jsr __f931         ; $cbb1: 20 31 f9  
            sbc #$40           ; $cbb4: e9 40     
            jsr __f937         ; $cbb6: 20 37 f9  
            iny                ; $cbb9: c8        
            jsr __f947         ; $cbba: 20 47 f9  
            sbc #$3f           ; $cbbd: e9 3f     
            jsr __f94c         ; $cbbf: 20 4c f9  
            iny                ; $cbc2: c8        
            jsr __f95c         ; $cbc3: 20 5c f9  
            sbc #$41           ; $cbc6: e9 41     
            jsr __f962         ; $cbc8: 20 62 f9  
            iny                ; $cbcb: c8        
            jsr __f972         ; $cbcc: 20 72 f9  
            sbc #$00           ; $cbcf: e9 00     
            jsr __f976         ; $cbd1: 20 76 f9  
            iny                ; $cbd4: c8        
            jsr __f980         ; $cbd5: 20 80 f9  
            sbc #$7f           ; $cbd8: e9 7f     
            jsr __f984         ; $cbda: 20 84 f9  
            rts                ; $cbdd: 60        

;-------------------------------------------------------------------------------
__cbde:     nop                ; $cbde: ea        
            lda #$ff           ; $cbdf: a9 ff     
            sta $01            ; $cbe1: 85 01     
            lda #$44           ; $cbe3: a9 44     
            ldx #$55           ; $cbe5: a2 55     
            ldy #$66           ; $cbe7: a0 66     
            inx                ; $cbe9: e8        
            dey                ; $cbea: 88        
            cpx #$56           ; $cbeb: e0 56     
            bne __cc10         ; $cbed: d0 21     
            cpy #$65           ; $cbef: c0 65     
            bne __cc10         ; $cbf1: d0 1d     
            inx                ; $cbf3: e8        
            inx                ; $cbf4: e8        
            dey                ; $cbf5: 88        
            dey                ; $cbf6: 88        
            cpx #$58           ; $cbf7: e0 58     
            bne __cc10         ; $cbf9: d0 15     
            cpy #$63           ; $cbfb: c0 63     
            bne __cc10         ; $cbfd: d0 11     
            dex                ; $cbff: ca        
            iny                ; $cc00: c8        
            cpx #$57           ; $cc01: e0 57     
            bne __cc10         ; $cc03: d0 0b     
            cpy #$64           ; $cc05: c0 64     
            bne __cc10         ; $cc07: d0 07     
            cmp #$44           ; $cc09: c9 44     
            bne __cc10         ; $cc0b: d0 03     
            jmp __cc14         ; $cc0d: 4c 14 cc  

;-------------------------------------------------------------------------------
__cc10:     ldx #$3e           ; $cc10: a2 3e     
            stx $00            ; $cc12: 86 00     
__cc14:     nop                ; $cc14: ea        
            sec                ; $cc15: 38        
            ldx #$69           ; $cc16: a2 69     
            lda #$96           ; $cc18: a9 96     
            bit $01            ; $cc1a: 24 01     
            ldy #$ff           ; $cc1c: a0 ff     
            iny                ; $cc1e: c8        
            bne __cc5e         ; $cc1f: d0 3d     
            bmi __cc5e         ; $cc21: 30 3b     
            bcc __cc5e         ; $cc23: 90 39     
            bvc __cc5e         ; $cc25: 50 37     
            cpy #$00           ; $cc27: c0 00     
            bne __cc5e         ; $cc29: d0 33     
            iny                ; $cc2b: c8        
            beq __cc5e         ; $cc2c: f0 30     
            bmi __cc5e         ; $cc2e: 30 2e     
            bcc __cc5e         ; $cc30: 90 2c     
            bvc __cc5e         ; $cc32: 50 2a     
            clc                ; $cc34: 18        
            clv                ; $cc35: b8        
            ldy #$00           ; $cc36: a0 00     
            dey                ; $cc38: 88        
            beq __cc5e         ; $cc39: f0 23     
            bpl __cc5e         ; $cc3b: 10 21     
            bcs __cc5e         ; $cc3d: b0 1f     
            bvs __cc5e         ; $cc3f: 70 1d     
            cpy #$ff           ; $cc41: c0 ff     
            bne __cc5e         ; $cc43: d0 19     
            clc                ; $cc45: 18        
            dey                ; $cc46: 88        
            beq __cc5e         ; $cc47: f0 15     
            bpl __cc5e         ; $cc49: 10 13     
            bcs __cc5e         ; $cc4b: b0 11     
            bvs __cc5e         ; $cc4d: 70 0f     
            cpy #$fe           ; $cc4f: c0 fe     
            bne __cc5e         ; $cc51: d0 0b     
            cmp #$96           ; $cc53: c9 96     
            bne __cc5e         ; $cc55: d0 07     
            cpx #$69           ; $cc57: e0 69     
            bne __cc5e         ; $cc59: d0 03     
            jmp __cc62         ; $cc5b: 4c 62 cc  

;-------------------------------------------------------------------------------
__cc5e:     ldx #$3f           ; $cc5e: a2 3f     
            stx $00            ; $cc60: 86 00     
__cc62:     nop                ; $cc62: ea        
            sec                ; $cc63: 38        
            ldy #$69           ; $cc64: a0 69     
            lda #$96           ; $cc66: a9 96     
            bit $01            ; $cc68: 24 01     
            ldx #$ff           ; $cc6a: a2 ff     
            inx                ; $cc6c: e8        
            bne __ccac         ; $cc6d: d0 3d     
            bmi __ccac         ; $cc6f: 30 3b     
            bcc __ccac         ; $cc71: 90 39     
            bvc __ccac         ; $cc73: 50 37     
            cpx #$00           ; $cc75: e0 00     
            bne __ccac         ; $cc77: d0 33     
            inx                ; $cc79: e8        
            beq __ccac         ; $cc7a: f0 30     
            bmi __ccac         ; $cc7c: 30 2e     
            bcc __ccac         ; $cc7e: 90 2c     
            bvc __ccac         ; $cc80: 50 2a     
            clc                ; $cc82: 18        
            clv                ; $cc83: b8        
            ldx #$00           ; $cc84: a2 00     
            dex                ; $cc86: ca        
            beq __ccac         ; $cc87: f0 23     
            bpl __ccac         ; $cc89: 10 21     
            bcs __ccac         ; $cc8b: b0 1f     
            bvs __ccac         ; $cc8d: 70 1d     
            cpx #$ff           ; $cc8f: e0 ff     
            bne __ccac         ; $cc91: d0 19     
            clc                ; $cc93: 18        
            dex                ; $cc94: ca        
            beq __ccac         ; $cc95: f0 15     
            bpl __ccac         ; $cc97: 10 13     
            bcs __ccac         ; $cc99: b0 11     
            bvs __ccac         ; $cc9b: 70 0f     
            cpx #$fe           ; $cc9d: e0 fe     
            bne __ccac         ; $cc9f: d0 0b     
            cmp #$96           ; $cca1: c9 96     
            bne __ccac         ; $cca3: d0 07     
            cpy #$69           ; $cca5: c0 69     
            bne __ccac         ; $cca7: d0 03     
            jmp __ccb0         ; $cca9: 4c b0 cc  

;-------------------------------------------------------------------------------
__ccac:     ldx #$40           ; $ccac: a2 40     
            stx $00            ; $ccae: 86 00     
__ccb0:     nop                ; $ccb0: ea        
            lda #$85           ; $ccb1: a9 85     
            ldx #$34           ; $ccb3: a2 34     
            ldy #$99           ; $ccb5: a0 99     
            clc                ; $ccb7: 18        
            bit $01            ; $ccb8: 24 01     
            tay                ; $ccba: a8        
            beq __cceb         ; $ccbb: f0 2e     
            bcs __cceb         ; $ccbd: b0 2c     
            bvc __cceb         ; $ccbf: 50 2a     
            bpl __cceb         ; $ccc1: 10 28     
            cmp #$85           ; $ccc3: c9 85     
            bne __cceb         ; $ccc5: d0 24     
            cpx #$34           ; $ccc7: e0 34     
            bne __cceb         ; $ccc9: d0 20     
            cpy #$85           ; $cccb: c0 85     
            bne __cceb         ; $cccd: d0 1c     
            lda #$00           ; $cccf: a9 00     
            sec                ; $ccd1: 38        
            clv                ; $ccd2: b8        
            tay                ; $ccd3: a8        
            bne __cceb         ; $ccd4: d0 15     
            bcc __cceb         ; $ccd6: 90 13     
            bvs __cceb         ; $ccd8: 70 11     
            bmi __cceb         ; $ccda: 30 0f     
            cmp #$00           ; $ccdc: c9 00     
            bne __cceb         ; $ccde: d0 0b     
            cpx #$34           ; $cce0: e0 34     
            bne __cceb         ; $cce2: d0 07     
            cpy #$00           ; $cce4: c0 00     
            bne __cceb         ; $cce6: d0 03     
            jmp __ccef         ; $cce8: 4c ef cc  

;-------------------------------------------------------------------------------
__cceb:     ldx #$41           ; $cceb: a2 41     
            stx $00            ; $cced: 86 00     
__ccef:     nop                ; $ccef: ea        
            lda #$85           ; $ccf0: a9 85     
            ldx #$34           ; $ccf2: a2 34     
            ldy #$99           ; $ccf4: a0 99     
            clc                ; $ccf6: 18        
            bit $01            ; $ccf7: 24 01     
            tax                ; $ccf9: aa        
            beq __cd2a         ; $ccfa: f0 2e     
            bcs __cd2a         ; $ccfc: b0 2c     
            bvc __cd2a         ; $ccfe: 50 2a     
            bpl __cd2a         ; $cd00: 10 28     
            cmp #$85           ; $cd02: c9 85     
            bne __cd2a         ; $cd04: d0 24     
            cpx #$85           ; $cd06: e0 85     
            bne __cd2a         ; $cd08: d0 20     
            cpy #$99           ; $cd0a: c0 99     
            bne __cd2a         ; $cd0c: d0 1c     
            lda #$00           ; $cd0e: a9 00     
            sec                ; $cd10: 38        
            clv                ; $cd11: b8        
            tax                ; $cd12: aa        
            bne __cd2a         ; $cd13: d0 15     
            bcc __cd2a         ; $cd15: 90 13     
            bvs __cd2a         ; $cd17: 70 11     
            bmi __cd2a         ; $cd19: 30 0f     
            cmp #$00           ; $cd1b: c9 00     
            bne __cd2a         ; $cd1d: d0 0b     
            cpx #$00           ; $cd1f: e0 00     
            bne __cd2a         ; $cd21: d0 07     
            cpy #$99           ; $cd23: c0 99     
            bne __cd2a         ; $cd25: d0 03     
            jmp __cd2e         ; $cd27: 4c 2e cd  

;-------------------------------------------------------------------------------
__cd2a:     ldx #$42           ; $cd2a: a2 42     
            stx $00            ; $cd2c: 86 00     
__cd2e:     nop                ; $cd2e: ea        
            lda #$85           ; $cd2f: a9 85     
            ldx #$34           ; $cd31: a2 34     
            ldy #$99           ; $cd33: a0 99     
            clc                ; $cd35: 18        
            bit $01            ; $cd36: 24 01     
            tya                ; $cd38: 98        
            beq __cd69         ; $cd39: f0 2e     
            bcs __cd69         ; $cd3b: b0 2c     
            bvc __cd69         ; $cd3d: 50 2a     
            bpl __cd69         ; $cd3f: 10 28     
            cmp #$99           ; $cd41: c9 99     
            bne __cd69         ; $cd43: d0 24     
            cpx #$34           ; $cd45: e0 34     
            bne __cd69         ; $cd47: d0 20     
            cpy #$99           ; $cd49: c0 99     
            bne __cd69         ; $cd4b: d0 1c     
            ldy #$00           ; $cd4d: a0 00     
            sec                ; $cd4f: 38        
            clv                ; $cd50: b8        
            tya                ; $cd51: 98        
            bne __cd69         ; $cd52: d0 15     
            bcc __cd69         ; $cd54: 90 13     
            bvs __cd69         ; $cd56: 70 11     
            bmi __cd69         ; $cd58: 30 0f     
            cmp #$00           ; $cd5a: c9 00     
            bne __cd69         ; $cd5c: d0 0b     
            cpx #$34           ; $cd5e: e0 34     
            bne __cd69         ; $cd60: d0 07     
            cpy #$00           ; $cd62: c0 00     
            bne __cd69         ; $cd64: d0 03     
            jmp __cd6d         ; $cd66: 4c 6d cd  

;-------------------------------------------------------------------------------
__cd69:     ldx #$43           ; $cd69: a2 43     
            stx $00            ; $cd6b: 86 00     
__cd6d:     nop                ; $cd6d: ea        
            lda #$85           ; $cd6e: a9 85     
            ldx #$34           ; $cd70: a2 34     
            ldy #$99           ; $cd72: a0 99     
            clc                ; $cd74: 18        
            bit $01            ; $cd75: 24 01     
            txa                ; $cd77: 8a        
            beq __cda8         ; $cd78: f0 2e     
            bcs __cda8         ; $cd7a: b0 2c     
            bvc __cda8         ; $cd7c: 50 2a     
            bmi __cda8         ; $cd7e: 30 28     
            cmp #$34           ; $cd80: c9 34     
            bne __cda8         ; $cd82: d0 24     
            cpx #$34           ; $cd84: e0 34     
            bne __cda8         ; $cd86: d0 20     
            cpy #$99           ; $cd88: c0 99     
            bne __cda8         ; $cd8a: d0 1c     
            ldx #$00           ; $cd8c: a2 00     
            sec                ; $cd8e: 38        
            clv                ; $cd8f: b8        
            txa                ; $cd90: 8a        
            bne __cda8         ; $cd91: d0 15     
            bcc __cda8         ; $cd93: 90 13     
            bvs __cda8         ; $cd95: 70 11     
            bmi __cda8         ; $cd97: 30 0f     
            cmp #$00           ; $cd99: c9 00     
            bne __cda8         ; $cd9b: d0 0b     
            cpx #$00           ; $cd9d: e0 00     
            bne __cda8         ; $cd9f: d0 07     
            cpy #$99           ; $cda1: c0 99     
            bne __cda8         ; $cda3: d0 03     
            jmp __cdac         ; $cda5: 4c ac cd  

;-------------------------------------------------------------------------------
__cda8:     ldx #$44           ; $cda8: a2 44     
            stx $00            ; $cdaa: 86 00     
__cdac:     nop                ; $cdac: ea        
            tsx                ; $cdad: ba        
            stx $07ff          ; $cdae: 8e ff 07  
            ldy #$33           ; $cdb1: a0 33     
            ldx #$69           ; $cdb3: a2 69     
            lda #$84           ; $cdb5: a9 84     
            clc                ; $cdb7: 18        
            bit $01            ; $cdb8: 24 01     
            txs                ; $cdba: 9a        
            beq __cdef         ; $cdbb: f0 32     
            bpl __cdef         ; $cdbd: 10 30     
            bcs __cdef         ; $cdbf: b0 2e     
            bvc __cdef         ; $cdc1: 50 2c     
            cmp #$84           ; $cdc3: c9 84     
            bne __cdef         ; $cdc5: d0 28     
            cpx #$69           ; $cdc7: e0 69     
            bne __cdef         ; $cdc9: d0 24     
            cpy #$33           ; $cdcb: c0 33     
            bne __cdef         ; $cdcd: d0 20     
            ldy #$01           ; $cdcf: a0 01     
            lda #$04           ; $cdd1: a9 04     
            sec                ; $cdd3: 38        
            clv                ; $cdd4: b8        
            ldx #$00           ; $cdd5: a2 00     
            tsx                ; $cdd7: ba        
            beq __cdef         ; $cdd8: f0 15     
            bmi __cdef         ; $cdda: 30 13     
            bcc __cdef         ; $cddc: 90 11     
            bvs __cdef         ; $cdde: 70 0f     
            cpx #$69           ; $cde0: e0 69     
            bne __cdef         ; $cde2: d0 0b     
            cmp #$04           ; $cde4: c9 04     
            bne __cdef         ; $cde6: d0 07     
            cpy #$01           ; $cde8: c0 01     
            bne __cdef         ; $cdea: d0 03     
            jmp __cdf3         ; $cdec: 4c f3 cd  

;-------------------------------------------------------------------------------
__cdef:     ldx #$45           ; $cdef: a2 45     
            stx $00            ; $cdf1: 86 00     
__cdf3:     ldx $07ff          ; $cdf3: ae ff 07  
            txs                ; $cdf6: 9a        
            rts                ; $cdf7: 60        

;-------------------------------------------------------------------------------
__cdf8:     lda #$ff           ; $cdf8: a9 ff     
            sta $01            ; $cdfa: 85 01     
            tsx                ; $cdfc: ba        
            stx $07ff          ; $cdfd: 8e ff 07  
            nop                ; $ce00: ea        
            ldx #$80           ; $ce01: a2 80     
            txs                ; $ce03: 9a        
            lda #$33           ; $ce04: a9 33     
            pha                ; $ce06: 48        
            lda #$69           ; $ce07: a9 69     
            pha                ; $ce09: 48        
            tsx                ; $ce0a: ba        
            cpx #$7e           ; $ce0b: e0 7e     
            bne __ce2f         ; $ce0d: d0 20     
            pla                ; $ce0f: 68        
            cmp #$69           ; $ce10: c9 69     
            bne __ce2f         ; $ce12: d0 1b     
            pla                ; $ce14: 68        
            cmp #$33           ; $ce15: c9 33     
            bne __ce2f         ; $ce17: d0 16     
            tsx                ; $ce19: ba        
            cpx #$80           ; $ce1a: e0 80     
            bne __ce2f         ; $ce1c: d0 11     
            lda $0180          ; $ce1e: ad 80 01  
            cmp #$33           ; $ce21: c9 33     
            bne __ce2f         ; $ce23: d0 0a     
            lda $017f          ; $ce25: ad 7f 01  
            cmp #$69           ; $ce28: c9 69     
            bne __ce2f         ; $ce2a: d0 03     
            jmp __ce33         ; $ce2c: 4c 33 ce  

;-------------------------------------------------------------------------------
__ce2f:     ldx #$46           ; $ce2f: a2 46     
            stx $00            ; $ce31: 86 00     
__ce33:     nop                ; $ce33: ea        
            ldx #$80           ; $ce34: a2 80     
            txs                ; $ce36: 9a        
            jsr __ce3d         ; $ce37: 20 3d ce  
            jmp __ce5b         ; $ce3a: 4c 5b ce  

;-------------------------------------------------------------------------------
__ce3d:     tsx                ; $ce3d: ba        
            cpx #$7e           ; $ce3e: e0 7e     
            bne __ce5b         ; $ce40: d0 19     
            pla                ; $ce42: 68        
            pla                ; $ce43: 68        
            tsx                ; $ce44: ba        
            cpx #$80           ; $ce45: e0 80     
            bne __ce5b         ; $ce47: d0 12     
            lda #$00           ; $ce49: a9 00     
            jsr __ce4e         ; $ce4b: 20 4e ce  
__ce4e:     pla                ; $ce4e: 68        
            cmp #$4d           ; $ce4f: c9 4d     
            bne __ce5b         ; $ce51: d0 08     
            pla                ; $ce53: 68        
            cmp #$ce           ; $ce54: c9 ce     
            bne __ce5b         ; $ce56: d0 03     
            jmp __ce5f         ; $ce58: 4c 5f ce  

;-------------------------------------------------------------------------------
__ce5b:     ldx #$47           ; $ce5b: a2 47     
            stx $00            ; $ce5d: 86 00     
__ce5f:     nop                ; $ce5f: ea        
            lda #$ce           ; $ce60: a9 ce     
            pha                ; $ce62: 48        
            lda #$66           ; $ce63: a9 66     
            pha                ; $ce65: 48        
__ce66:     rts                ; $ce66: 60        

;-------------------------------------------------------------------------------
            ldx #$77           ; $ce67: a2 77     
            ldy #$69           ; $ce69: a0 69     
            clc                ; $ce6b: 18        
            bit $01            ; $ce6c: 24 01     
            lda #$83           ; $ce6e: a9 83     
            jsr __ce66         ; $ce70: 20 66 ce  
            beq __ce99         ; $ce73: f0 24     
            bpl __ce99         ; $ce75: 10 22     
            bcs __ce99         ; $ce77: b0 20     
            bvc __ce99         ; $ce79: 50 1e     
            cmp #$83           ; $ce7b: c9 83     
            bne __ce99         ; $ce7d: d0 1a     
            cpy #$69           ; $ce7f: c0 69     
            bne __ce99         ; $ce81: d0 16     
            cpx #$77           ; $ce83: e0 77     
            bne __ce99         ; $ce85: d0 12     
            sec                ; $ce87: 38        
            clv                ; $ce88: b8        
            lda #$00           ; $ce89: a9 00     
            jsr __ce66         ; $ce8b: 20 66 ce  
            bne __ce99         ; $ce8e: d0 09     
            bmi __ce99         ; $ce90: 30 07     
            bcc __ce99         ; $ce92: 90 05     
            bvs __ce99         ; $ce94: 70 03     
            jmp __ce9d         ; $ce96: 4c 9d ce  

;-------------------------------------------------------------------------------
__ce99:     ldx #$48           ; $ce99: a2 48     
            stx $00            ; $ce9b: 86 00     
__ce9d:     nop                ; $ce9d: ea        
            lda #$ce           ; $ce9e: a9 ce     
            pha                ; $cea0: 48        
            lda #$ae           ; $cea1: a9 ae     
            pha                ; $cea3: 48        
            lda #$65           ; $cea4: a9 65     
            pha                ; $cea6: 48        
            lda #$55           ; $cea7: a9 55     
            ldy #$88           ; $cea9: a0 88     
            ldx #$99           ; $ceab: a2 99     
            rti                ; $cead: 40        

;-------------------------------------------------------------------------------
            bmi __cee5         ; $ceae: 30 35     
            bvc __cee5         ; $ceb0: 50 33     
            beq __cee5         ; $ceb2: f0 31     
            bcc __cee5         ; $ceb4: 90 2f     
            cmp #$55           ; $ceb6: c9 55     
            bne __cee5         ; $ceb8: d0 2b     
            cpy #$88           ; $ceba: c0 88     
            bne __cee5         ; $cebc: d0 27     
            cpx #$99           ; $cebe: e0 99     
            bne __cee5         ; $cec0: d0 23     
            lda #$ce           ; $cec2: a9 ce     
            pha                ; $cec4: 48        
            lda #$ce           ; $cec5: a9 ce     
            pha                ; $cec7: 48        
            lda #$87           ; $cec8: a9 87     
            pha                ; $ceca: 48        
            lda #$55           ; $cecb: a9 55     
            rti                ; $cecd: 40        

;-------------------------------------------------------------------------------
            bpl __cee5         ; $cece: 10 15     
            bvs __cee5         ; $ced0: 70 13     
            bne __cee5         ; $ced2: d0 11     
            bcc __cee5         ; $ced4: 90 0f     
            cmp #$55           ; $ced6: c9 55     
            bne __cee5         ; $ced8: d0 0b     
            cpy #$88           ; $ceda: c0 88     
            bne __cee5         ; $cedc: d0 07     
            cpx #$99           ; $cede: e0 99     
            bne __cee5         ; $cee0: d0 03     
            jmp __cee9         ; $cee2: 4c e9 ce  

;-------------------------------------------------------------------------------
__cee5:     ldx #$49           ; $cee5: a2 49     
            stx $00            ; $cee7: 86 00     
__cee9:     ldx $07ff          ; $cee9: ae ff 07  
            txs                ; $ceec: 9a        
            rts                ; $ceed: 60        

;-------------------------------------------------------------------------------
__ceee:     ldx #$55           ; $ceee: a2 55     
            ldy #$69           ; $cef0: a0 69     
            lda #$ff           ; $cef2: a9 ff     
            sta $01            ; $cef4: 85 01     
            nop                ; $cef6: ea        
            bit $01            ; $cef7: 24 01     
            sec                ; $cef9: 38        
            lda #$01           ; $cefa: a9 01     
            lsr                ; $cefc: 4a        
            bcc __cf1c         ; $cefd: 90 1d     
            bne __cf1c         ; $ceff: d0 1b     
            bmi __cf1c         ; $cf01: 30 19     
            bvc __cf1c         ; $cf03: 50 17     
            cmp #$00           ; $cf05: c9 00     
            bne __cf1c         ; $cf07: d0 13     
            clv                ; $cf09: b8        
            lda #$aa           ; $cf0a: a9 aa     
            lsr                ; $cf0c: 4a        
            bcs __cf1c         ; $cf0d: b0 0d     
            beq __cf1c         ; $cf0f: f0 0b     
            bmi __cf1c         ; $cf11: 30 09     
            bvs __cf1c         ; $cf13: 70 07     
            cmp #$55           ; $cf15: c9 55     
            bne __cf1c         ; $cf17: d0 03     
            jmp __cf20         ; $cf19: 4c 20 cf  

;-------------------------------------------------------------------------------
__cf1c:     ldx #$4a           ; $cf1c: a2 4a     
            stx $00            ; $cf1e: 86 00     
__cf20:     nop                ; $cf20: ea        
            bit $01            ; $cf21: 24 01     
            sec                ; $cf23: 38        
            lda #$80           ; $cf24: a9 80     
            asl                ; $cf26: 0a        
            bcc __cf47         ; $cf27: 90 1e     
            bne __cf47         ; $cf29: d0 1c     
            bmi __cf47         ; $cf2b: 30 1a     
            bvc __cf47         ; $cf2d: 50 18     
            cmp #$00           ; $cf2f: c9 00     
            bne __cf47         ; $cf31: d0 14     
            clv                ; $cf33: b8        
            sec                ; $cf34: 38        
            lda #$55           ; $cf35: a9 55     
            asl                ; $cf37: 0a        
            bcs __cf47         ; $cf38: b0 0d     
            beq __cf47         ; $cf3a: f0 0b     
            bpl __cf47         ; $cf3c: 10 09     
            bvs __cf47         ; $cf3e: 70 07     
            cmp #$aa           ; $cf40: c9 aa     
            bne __cf47         ; $cf42: d0 03     
            jmp __cf4b         ; $cf44: 4c 4b cf  

;-------------------------------------------------------------------------------
__cf47:     ldx #$4b           ; $cf47: a2 4b     
            stx $00            ; $cf49: 86 00     
__cf4b:     nop                ; $cf4b: ea        
            bit $01            ; $cf4c: 24 01     
            sec                ; $cf4e: 38        
            lda #$01           ; $cf4f: a9 01     
            ror                ; $cf51: 6a        
            bcc __cf72         ; $cf52: 90 1e     
            beq __cf72         ; $cf54: f0 1c     
            bpl __cf72         ; $cf56: 10 1a     
            bvc __cf72         ; $cf58: 50 18     
            cmp #$80           ; $cf5a: c9 80     
            bne __cf72         ; $cf5c: d0 14     
            clv                ; $cf5e: b8        
            clc                ; $cf5f: 18        
            lda #$55           ; $cf60: a9 55     
            ror                ; $cf62: 6a        
            bcc __cf72         ; $cf63: 90 0d     
            beq __cf72         ; $cf65: f0 0b     
            bmi __cf72         ; $cf67: 30 09     
            bvs __cf72         ; $cf69: 70 07     
            cmp #$2a           ; $cf6b: c9 2a     
            bne __cf72         ; $cf6d: d0 03     
            jmp __cf76         ; $cf6f: 4c 76 cf  

;-------------------------------------------------------------------------------
__cf72:     ldx #$4c           ; $cf72: a2 4c     
            stx $00            ; $cf74: 86 00     
__cf76:     nop                ; $cf76: ea        
            bit $01            ; $cf77: 24 01     
            sec                ; $cf79: 38        
            lda #$80           ; $cf7a: a9 80     
            rol                ; $cf7c: 2a        
            bcc __cf9d         ; $cf7d: 90 1e     
            beq __cf9d         ; $cf7f: f0 1c     
            bmi __cf9d         ; $cf81: 30 1a     
            bvc __cf9d         ; $cf83: 50 18     
            cmp #$01           ; $cf85: c9 01     
            bne __cf9d         ; $cf87: d0 14     
            clv                ; $cf89: b8        
            clc                ; $cf8a: 18        
            lda #$55           ; $cf8b: a9 55     
            rol                ; $cf8d: 2a        
            bcs __cf9d         ; $cf8e: b0 0d     
            beq __cf9d         ; $cf90: f0 0b     
            bpl __cf9d         ; $cf92: 10 09     
            bvs __cf9d         ; $cf94: 70 07     
            cmp #$aa           ; $cf96: c9 aa     
            bne __cf9d         ; $cf98: d0 03     
            jmp __cfa1         ; $cf9a: 4c a1 cf  

;-------------------------------------------------------------------------------
__cf9d:     ldx #$4d           ; $cf9d: a2 4d     
            stx $00            ; $cf9f: 86 00     
__cfa1:     rts                ; $cfa1: 60        

;-------------------------------------------------------------------------------
__cfa2:     lda $00            ; $cfa2: a5 00     
            sta $07ff          ; $cfa4: 8d ff 07  
            lda #$00           ; $cfa7: a9 00     
            sta $80            ; $cfa9: 85 80     
            lda #$02           ; $cfab: a9 02     
            sta $81            ; $cfad: 85 81     
            lda #$ff           ; $cfaf: a9 ff     
            sta $01            ; $cfb1: 85 01     
            lda #$00           ; $cfb3: a9 00     
            sta $82            ; $cfb5: 85 82     
            lda #$03           ; $cfb7: a9 03     
            sta $83            ; $cfb9: 85 83     
            sta $84            ; $cfbb: 85 84     
            lda #$00           ; $cfbd: a9 00     
            sta $ff            ; $cfbf: 85 ff     
            lda #$04           ; $cfc1: a9 04     
            sta $00            ; $cfc3: 85 00     
            lda #$5a           ; $cfc5: a9 5a     
            sta $0200          ; $cfc7: 8d 00 02  
            lda #$5b           ; $cfca: a9 5b     
            sta $0300          ; $cfcc: 8d 00 03  
            lda #$5c           ; $cfcf: a9 5c     
            sta $0303          ; $cfd1: 8d 03 03  
            lda #$5d           ; $cfd4: a9 5d     
            sta $0400          ; $cfd6: 8d 00 04  
            ldx #$00           ; $cfd9: a2 00     
            lda ($80,x)        ; $cfdb: a1 80     
            cmp #$5a           ; $cfdd: c9 5a     
            bne __d000         ; $cfdf: d0 1f     
            inx                ; $cfe1: e8        
            inx                ; $cfe2: e8        
            lda ($80,x)        ; $cfe3: a1 80     
            cmp #$5b           ; $cfe5: c9 5b     
            bne __d000         ; $cfe7: d0 17     
            inx                ; $cfe9: e8        
            lda ($80,x)        ; $cfea: a1 80     
            cmp #$5c           ; $cfec: c9 5c     
            bne __d000         ; $cfee: d0 10     
            ldx #$00           ; $cff0: a2 00     
            lda ($ff,x)        ; $cff2: a1 ff     
            cmp #$5d           ; $cff4: c9 5d     
            bne __d000         ; $cff6: d0 08     
            ldx #$81           ; $cff8: a2 81     
            lda ($ff,x)        ; $cffa: a1 ff     
            cmp #$5a           ; $cffc: c9 5a     
            beq __d005         ; $cffe: f0 05     
__d000:     lda #$58           ; $d000: a9 58     
            sta $07ff          ; $d002: 8d ff 07  
__d005:     lda #$aa           ; $d005: a9 aa     
            ldx #$00           ; $d007: a2 00     
            sta ($80,x)        ; $d009: 81 80     
            inx                ; $d00b: e8        
            inx                ; $d00c: e8        
            lda #$ab           ; $d00d: a9 ab     
            sta ($80,x)        ; $d00f: 81 80     
            inx                ; $d011: e8        
            lda #$ac           ; $d012: a9 ac     
            sta ($80,x)        ; $d014: 81 80     
            ldx #$00           ; $d016: a2 00     
            lda #$ad           ; $d018: a9 ad     
            sta ($ff,x)        ; $d01a: 81 ff     
            lda $0200          ; $d01c: ad 00 02  
            cmp #$aa           ; $d01f: c9 aa     
            bne __d038         ; $d021: d0 15     
            lda $0300          ; $d023: ad 00 03  
            cmp #$ab           ; $d026: c9 ab     
            bne __d038         ; $d028: d0 0e     
            lda $0303          ; $d02a: ad 03 03  
            cmp #$ac           ; $d02d: c9 ac     
            bne __d038         ; $d02f: d0 07     
            lda $0400          ; $d031: ad 00 04  
            cmp #$ad           ; $d034: c9 ad     
            beq __d03d         ; $d036: f0 05     
__d038:     lda #$59           ; $d038: a9 59     
            sta $07ff          ; $d03a: 8d ff 07  
__d03d:     lda $07ff          ; $d03d: ad ff 07  
            sta $00            ; $d040: 85 00     
            lda #$00           ; $d042: a9 00     
            sta $0300          ; $d044: 8d 00 03  
            lda #$aa           ; $d047: a9 aa     
            sta $0200          ; $d049: 8d 00 02  
            ldx #$00           ; $d04c: a2 00     
            ldy #$5a           ; $d04e: a0 5a     
            jsr __f7b6         ; $d050: 20 b6 f7  
            ora ($80,x)        ; $d053: 01 80     
            jsr __f7c0         ; $d055: 20 c0 f7  
            iny                ; $d058: c8        
            jsr __f7ce         ; $d059: 20 ce f7  
            ora ($82,x)        ; $d05c: 01 82     
            jsr __f7d3         ; $d05e: 20 d3 f7  
            iny                ; $d061: c8        
            jsr __f7df         ; $d062: 20 df f7  
            and ($80,x)        ; $d065: 21 80     
            jsr __f7e5         ; $d067: 20 e5 f7  
            iny                ; $d06a: c8        
            lda #$ef           ; $d06b: a9 ef     
            sta $0300          ; $d06d: 8d 00 03  
            jsr __f7f1         ; $d070: 20 f1 f7  
            and ($82,x)        ; $d073: 21 82     
            jsr __f7f6         ; $d075: 20 f6 f7  
            iny                ; $d078: c8        
            jsr __f804         ; $d079: 20 04 f8  
            eor ($80,x)        ; $d07c: 41 80     
            jsr __f80a         ; $d07e: 20 0a f8  
            iny                ; $d081: c8        
            lda #$70           ; $d082: a9 70     
            sta $0300          ; $d084: 8d 00 03  
            jsr __f818         ; $d087: 20 18 f8  
            eor ($82,x)        ; $d08a: 41 82     
            jsr __f81d         ; $d08c: 20 1d f8  
            iny                ; $d08f: c8        
            lda #$69           ; $d090: a9 69     
            sta $0200          ; $d092: 8d 00 02  
            jsr __f829         ; $d095: 20 29 f8  
            adc ($80,x)        ; $d098: 61 80     
            jsr __f82f         ; $d09a: 20 2f f8  
            iny                ; $d09d: c8        
            jsr __f83d         ; $d09e: 20 3d f8  
            adc ($80,x)        ; $d0a1: 61 80     
            jsr __f843         ; $d0a3: 20 43 f8  
            iny                ; $d0a6: c8        
            lda #$7f           ; $d0a7: a9 7f     
            sta $0200          ; $d0a9: 8d 00 02  
            jsr __f851         ; $d0ac: 20 51 f8  
            adc ($80,x)        ; $d0af: 61 80     
            jsr __f856         ; $d0b1: 20 56 f8  
            iny                ; $d0b4: c8        
            lda #$80           ; $d0b5: a9 80     
            sta $0200          ; $d0b7: 8d 00 02  
            jsr __f864         ; $d0ba: 20 64 f8  
            adc ($80,x)        ; $d0bd: 61 80     
            jsr __f86a         ; $d0bf: 20 6a f8  
            iny                ; $d0c2: c8        
            jsr __f878         ; $d0c3: 20 78 f8  
            adc ($80,x)        ; $d0c6: 61 80     
            jsr __f87d         ; $d0c8: 20 7d f8  
            iny                ; $d0cb: c8        
            lda #$40           ; $d0cc: a9 40     
            sta $0200          ; $d0ce: 8d 00 02  
            jsr __f889         ; $d0d1: 20 89 f8  
            cmp ($80,x)        ; $d0d4: c1 80     
            jsr __f88e         ; $d0d6: 20 8e f8  
            iny                ; $d0d9: c8        
            pha                ; $d0da: 48        
            lda #$3f           ; $d0db: a9 3f     
            sta $0200          ; $d0dd: 8d 00 02  
            pla                ; $d0e0: 68        
            jsr __f89a         ; $d0e1: 20 9a f8  
            cmp ($80,x)        ; $d0e4: c1 80     
            jsr __f89c         ; $d0e6: 20 9c f8  
            iny                ; $d0e9: c8        
            pha                ; $d0ea: 48        
            lda #$41           ; $d0eb: a9 41     
            sta $0200          ; $d0ed: 8d 00 02  
            pla                ; $d0f0: 68        
            cmp ($80,x)        ; $d0f1: c1 80     
            jsr __f8a8         ; $d0f3: 20 a8 f8  
            iny                ; $d0f6: c8        
            pha                ; $d0f7: 48        
            lda #$00           ; $d0f8: a9 00     
            sta $0200          ; $d0fa: 8d 00 02  
            pla                ; $d0fd: 68        
            jsr __f8b2         ; $d0fe: 20 b2 f8  
            cmp ($80,x)        ; $d101: c1 80     
            jsr __f8b5         ; $d103: 20 b5 f8  
            iny                ; $d106: c8        
            pha                ; $d107: 48        
            lda #$80           ; $d108: a9 80     
            sta $0200          ; $d10a: 8d 00 02  
            pla                ; $d10d: 68        
            cmp ($80,x)        ; $d10e: c1 80     
            jsr __f8bf         ; $d110: 20 bf f8  
            iny                ; $d113: c8        
            pha                ; $d114: 48        
            lda #$81           ; $d115: a9 81     
            sta $0200          ; $d117: 8d 00 02  
            pla                ; $d11a: 68        
            cmp ($80,x)        ; $d11b: c1 80     
            jsr __f8c9         ; $d11d: 20 c9 f8  
            iny                ; $d120: c8        
            pha                ; $d121: 48        
            lda #$7f           ; $d122: a9 7f     
            sta $0200          ; $d124: 8d 00 02  
            pla                ; $d127: 68        
            cmp ($80,x)        ; $d128: c1 80     
            jsr __f8d3         ; $d12a: 20 d3 f8  
            iny                ; $d12d: c8        
            lda #$40           ; $d12e: a9 40     
            sta $0200          ; $d130: 8d 00 02  
            jsr __f931         ; $d133: 20 31 f9  
            sbc ($80,x)        ; $d136: e1 80     
            jsr __f937         ; $d138: 20 37 f9  
            iny                ; $d13b: c8        
            lda #$3f           ; $d13c: a9 3f     
            sta $0200          ; $d13e: 8d 00 02  
            jsr __f947         ; $d141: 20 47 f9  
            sbc ($80,x)        ; $d144: e1 80     
            jsr __f94c         ; $d146: 20 4c f9  
            iny                ; $d149: c8        
            lda #$41           ; $d14a: a9 41     
            sta $0200          ; $d14c: 8d 00 02  
            jsr __f95c         ; $d14f: 20 5c f9  
            sbc ($80,x)        ; $d152: e1 80     
            jsr __f962         ; $d154: 20 62 f9  
            iny                ; $d157: c8        
            lda #$00           ; $d158: a9 00     
            sta $0200          ; $d15a: 8d 00 02  
            jsr __f972         ; $d15d: 20 72 f9  
            sbc ($80,x)        ; $d160: e1 80     
            jsr __f976         ; $d162: 20 76 f9  
            iny                ; $d165: c8        
            lda #$7f           ; $d166: a9 7f     
            sta $0200          ; $d168: 8d 00 02  
            jsr __f980         ; $d16b: 20 80 f9  
            sbc ($80,x)        ; $d16e: e1 80     
            jsr __f984         ; $d170: 20 84 f9  
            rts                ; $d173: 60        

;-------------------------------------------------------------------------------
__d174:     lda #$55           ; $d174: a9 55     
            sta $78            ; $d176: 85 78     
            lda #$ff           ; $d178: a9 ff     
            sta $01            ; $d17a: 85 01     
            bit $01            ; $d17c: 24 01     
            ldy #$11           ; $d17e: a0 11     
            ldx #$23           ; $d180: a2 23     
            lda #$00           ; $d182: a9 00     
            lda $78            ; $d184: a5 78     
            beq __d198         ; $d186: f0 10     
            bmi __d198         ; $d188: 30 0e     
            cmp #$55           ; $d18a: c9 55     
            bne __d198         ; $d18c: d0 0a     
            cpy #$11           ; $d18e: c0 11     
            bne __d198         ; $d190: d0 06     
            cpx #$23           ; $d192: e0 23     
            bvc __d198         ; $d194: 50 02     
            beq __d19c         ; $d196: f0 04     
__d198:     lda #$76           ; $d198: a9 76     
            sta $00            ; $d19a: 85 00     
__d19c:     lda #$46           ; $d19c: a9 46     
            bit $01            ; $d19e: 24 01     
            sta $78            ; $d1a0: 85 78     
            beq __d1ae         ; $d1a2: f0 0a     
            bpl __d1ae         ; $d1a4: 10 08     
            bvc __d1ae         ; $d1a6: 50 06     
            lda $78            ; $d1a8: a5 78     
            cmp #$46           ; $d1aa: c9 46     
            beq __d1b2         ; $d1ac: f0 04     
__d1ae:     lda #$77           ; $d1ae: a9 77     
            sta $00            ; $d1b0: 85 00     
__d1b2:     lda #$55           ; $d1b2: a9 55     
            sta $78            ; $d1b4: 85 78     
            bit $01            ; $d1b6: 24 01     
            lda #$11           ; $d1b8: a9 11     
            ldx #$23           ; $d1ba: a2 23     
            ldy #$00           ; $d1bc: a0 00     
            ldy $78            ; $d1be: a4 78     
            beq __d1d2         ; $d1c0: f0 10     
            bmi __d1d2         ; $d1c2: 30 0e     
            cpy #$55           ; $d1c4: c0 55     
            bne __d1d2         ; $d1c6: d0 0a     
            cmp #$11           ; $d1c8: c9 11     
            bne __d1d2         ; $d1ca: d0 06     
            cpx #$23           ; $d1cc: e0 23     
            bvc __d1d2         ; $d1ce: 50 02     
            beq __d1d6         ; $d1d0: f0 04     
__d1d2:     lda #$78           ; $d1d2: a9 78     
            sta $00            ; $d1d4: 85 00     
__d1d6:     ldy #$46           ; $d1d6: a0 46     
            bit $01            ; $d1d8: 24 01     
            sty $78            ; $d1da: 84 78     
            beq __d1e8         ; $d1dc: f0 0a     
            bpl __d1e8         ; $d1de: 10 08     
            bvc __d1e8         ; $d1e0: 50 06     
            ldy $78            ; $d1e2: a4 78     
            cpy #$46           ; $d1e4: c0 46     
            beq __d1ec         ; $d1e6: f0 04     
__d1e8:     lda #$79           ; $d1e8: a9 79     
            sta $00            ; $d1ea: 85 00     
__d1ec:     bit $01            ; $d1ec: 24 01     
            lda #$55           ; $d1ee: a9 55     
            sta $78            ; $d1f0: 85 78     
            ldy #$11           ; $d1f2: a0 11     
            lda #$23           ; $d1f4: a9 23     
            ldx #$00           ; $d1f6: a2 00     
            ldx $78            ; $d1f8: a6 78     
            beq __d20c         ; $d1fa: f0 10     
            bmi __d20c         ; $d1fc: 30 0e     
            cpx #$55           ; $d1fe: e0 55     
            bne __d20c         ; $d200: d0 0a     
            cpy #$11           ; $d202: c0 11     
            bne __d20c         ; $d204: d0 06     
            cmp #$23           ; $d206: c9 23     
            bvc __d20c         ; $d208: 50 02     
            beq __d210         ; $d20a: f0 04     
__d20c:     lda #$7a           ; $d20c: a9 7a     
            sta $00            ; $d20e: 85 00     
__d210:     ldx #$46           ; $d210: a2 46     
            bit $01            ; $d212: 24 01     
            stx $78            ; $d214: 86 78     
            beq __d222         ; $d216: f0 0a     
            bpl __d222         ; $d218: 10 08     
            bvc __d222         ; $d21a: 50 06     
            ldx $78            ; $d21c: a6 78     
            cpx #$46           ; $d21e: e0 46     
            beq __d226         ; $d220: f0 04     
__d222:     lda #$7b           ; $d222: a9 7b     
            sta $00            ; $d224: 85 00     
__d226:     lda #$c0           ; $d226: a9 c0     
            sta $78            ; $d228: 85 78     
            ldx #$33           ; $d22a: a2 33     
            ldy #$88           ; $d22c: a0 88     
            lda #$05           ; $d22e: a9 05     
            bit $78            ; $d230: 24 78     
            bpl __d244         ; $d232: 10 10     
            bvc __d244         ; $d234: 50 0e     
            bne __d244         ; $d236: d0 0c     
            cmp #$05           ; $d238: c9 05     
            bne __d244         ; $d23a: d0 08     
            cpx #$33           ; $d23c: e0 33     
            bne __d244         ; $d23e: d0 04     
            cpy #$88           ; $d240: c0 88     
            beq __d248         ; $d242: f0 04     
__d244:     lda #$7c           ; $d244: a9 7c     
            sta $00            ; $d246: 85 00     
__d248:     lda #$03           ; $d248: a9 03     
            sta $78            ; $d24a: 85 78     
            lda #$01           ; $d24c: a9 01     
            bit $78            ; $d24e: 24 78     
            bmi __d25a         ; $d250: 30 08     
            bvs __d25a         ; $d252: 70 06     
            beq __d25a         ; $d254: f0 04     
            cmp #$01           ; $d256: c9 01     
            beq __d25e         ; $d258: f0 04     
__d25a:     lda #$7d           ; $d25a: a9 7d     
            sta $00            ; $d25c: 85 00     
__d25e:     ldy #$7e           ; $d25e: a0 7e     
            lda #$aa           ; $d260: a9 aa     
            sta $78            ; $d262: 85 78     
            jsr __f7b6         ; $d264: 20 b6 f7  
            ora $78            ; $d267: 05 78     
            jsr __f7c0         ; $d269: 20 c0 f7  
            iny                ; $d26c: c8        
            lda #$00           ; $d26d: a9 00     
            sta $78            ; $d26f: 85 78     
            jsr __f7ce         ; $d271: 20 ce f7  
            ora $78            ; $d274: 05 78     
            jsr __f7d3         ; $d276: 20 d3 f7  
            iny                ; $d279: c8        
            lda #$aa           ; $d27a: a9 aa     
            sta $78            ; $d27c: 85 78     
            jsr __f7df         ; $d27e: 20 df f7  
            and $78            ; $d281: 25 78     
            jsr __f7e5         ; $d283: 20 e5 f7  
            iny                ; $d286: c8        
            lda #$ef           ; $d287: a9 ef     
            sta $78            ; $d289: 85 78     
            jsr __f7f1         ; $d28b: 20 f1 f7  
            and $78            ; $d28e: 25 78     
            jsr __f7f6         ; $d290: 20 f6 f7  
            iny                ; $d293: c8        
            lda #$aa           ; $d294: a9 aa     
            sta $78            ; $d296: 85 78     
            jsr __f804         ; $d298: 20 04 f8  
            eor $78            ; $d29b: 45 78     
            jsr __f80a         ; $d29d: 20 0a f8  
            iny                ; $d2a0: c8        
            lda #$70           ; $d2a1: a9 70     
            sta $78            ; $d2a3: 85 78     
            jsr __f818         ; $d2a5: 20 18 f8  
            eor $78            ; $d2a8: 45 78     
            jsr __f81d         ; $d2aa: 20 1d f8  
            iny                ; $d2ad: c8        
            lda #$69           ; $d2ae: a9 69     
            sta $78            ; $d2b0: 85 78     
            jsr __f829         ; $d2b2: 20 29 f8  
            adc $78            ; $d2b5: 65 78     
            jsr __f82f         ; $d2b7: 20 2f f8  
            iny                ; $d2ba: c8        
            jsr __f83d         ; $d2bb: 20 3d f8  
            adc $78            ; $d2be: 65 78     
            jsr __f843         ; $d2c0: 20 43 f8  
            iny                ; $d2c3: c8        
            lda #$7f           ; $d2c4: a9 7f     
            sta $78            ; $d2c6: 85 78     
            jsr __f851         ; $d2c8: 20 51 f8  
            adc $78            ; $d2cb: 65 78     
            jsr __f856         ; $d2cd: 20 56 f8  
            iny                ; $d2d0: c8        
            lda #$80           ; $d2d1: a9 80     
            sta $78            ; $d2d3: 85 78     
            jsr __f864         ; $d2d5: 20 64 f8  
            adc $78            ; $d2d8: 65 78     
            jsr __f86a         ; $d2da: 20 6a f8  
            iny                ; $d2dd: c8        
            jsr __f878         ; $d2de: 20 78 f8  
            adc $78            ; $d2e1: 65 78     
            jsr __f87d         ; $d2e3: 20 7d f8  
            iny                ; $d2e6: c8        
            lda #$40           ; $d2e7: a9 40     
            sta $78            ; $d2e9: 85 78     
            jsr __f889         ; $d2eb: 20 89 f8  
            cmp $78            ; $d2ee: c5 78     
            jsr __f88e         ; $d2f0: 20 8e f8  
            iny                ; $d2f3: c8        
            pha                ; $d2f4: 48        
            lda #$3f           ; $d2f5: a9 3f     
            sta $78            ; $d2f7: 85 78     
            pla                ; $d2f9: 68        
            jsr __f89a         ; $d2fa: 20 9a f8  
            cmp $78            ; $d2fd: c5 78     
            jsr __f89c         ; $d2ff: 20 9c f8  
            iny                ; $d302: c8        
            pha                ; $d303: 48        
            lda #$41           ; $d304: a9 41     
            sta $78            ; $d306: 85 78     
            pla                ; $d308: 68        
            cmp $78            ; $d309: c5 78     
            jsr __f8a8         ; $d30b: 20 a8 f8  
            iny                ; $d30e: c8        
            pha                ; $d30f: 48        
            lda #$00           ; $d310: a9 00     
            sta $78            ; $d312: 85 78     
            pla                ; $d314: 68        
            jsr __f8b2         ; $d315: 20 b2 f8  
            cmp $78            ; $d318: c5 78     
            jsr __f8b5         ; $d31a: 20 b5 f8  
            iny                ; $d31d: c8        
            pha                ; $d31e: 48        
            lda #$80           ; $d31f: a9 80     
            sta $78            ; $d321: 85 78     
            pla                ; $d323: 68        
            cmp $78            ; $d324: c5 78     
            jsr __f8bf         ; $d326: 20 bf f8  
            iny                ; $d329: c8        
            pha                ; $d32a: 48        
            lda #$81           ; $d32b: a9 81     
            sta $78            ; $d32d: 85 78     
            pla                ; $d32f: 68        
            cmp $78            ; $d330: c5 78     
            jsr __f8c9         ; $d332: 20 c9 f8  
            iny                ; $d335: c8        
            pha                ; $d336: 48        
            lda #$7f           ; $d337: a9 7f     
            sta $78            ; $d339: 85 78     
            pla                ; $d33b: 68        
            cmp $78            ; $d33c: c5 78     
            jsr __f8d3         ; $d33e: 20 d3 f8  
            iny                ; $d341: c8        
            lda #$40           ; $d342: a9 40     
            sta $78            ; $d344: 85 78     
            jsr __f931         ; $d346: 20 31 f9  
            sbc $78            ; $d349: e5 78     
            jsr __f937         ; $d34b: 20 37 f9  
            iny                ; $d34e: c8        
            lda #$3f           ; $d34f: a9 3f     
            sta $78            ; $d351: 85 78     
            jsr __f947         ; $d353: 20 47 f9  
            sbc $78            ; $d356: e5 78     
            jsr __f94c         ; $d358: 20 4c f9  
            iny                ; $d35b: c8        
            lda #$41           ; $d35c: a9 41     
            sta $78            ; $d35e: 85 78     
            jsr __f95c         ; $d360: 20 5c f9  
            sbc $78            ; $d363: e5 78     
            jsr __f962         ; $d365: 20 62 f9  
            iny                ; $d368: c8        
            lda #$00           ; $d369: a9 00     
            sta $78            ; $d36b: 85 78     
            jsr __f972         ; $d36d: 20 72 f9  
            sbc $78            ; $d370: e5 78     
            jsr __f976         ; $d372: 20 76 f9  
            iny                ; $d375: c8        
            lda #$7f           ; $d376: a9 7f     
            sta $78            ; $d378: 85 78     
            jsr __f980         ; $d37a: 20 80 f9  
            sbc $78            ; $d37d: e5 78     
            jsr __f984         ; $d37f: 20 84 f9  
            iny                ; $d382: c8        
            lda #$40           ; $d383: a9 40     
            sta $78            ; $d385: 85 78     
            jsr __f889         ; $d387: 20 89 f8  
            tax                ; $d38a: aa        
            cpx $78            ; $d38b: e4 78     
            jsr __f88e         ; $d38d: 20 8e f8  
            iny                ; $d390: c8        
            lda #$3f           ; $d391: a9 3f     
            sta $78            ; $d393: 85 78     
            jsr __f89a         ; $d395: 20 9a f8  
            cpx $78            ; $d398: e4 78     
            jsr __f89c         ; $d39a: 20 9c f8  
            iny                ; $d39d: c8        
            lda #$41           ; $d39e: a9 41     
            sta $78            ; $d3a0: 85 78     
            cpx $78            ; $d3a2: e4 78     
            jsr __f8a8         ; $d3a4: 20 a8 f8  
            iny                ; $d3a7: c8        
            lda #$00           ; $d3a8: a9 00     
            sta $78            ; $d3aa: 85 78     
            jsr __f8b2         ; $d3ac: 20 b2 f8  
            tax                ; $d3af: aa        
            cpx $78            ; $d3b0: e4 78     
            jsr __f8b5         ; $d3b2: 20 b5 f8  
            iny                ; $d3b5: c8        
            lda #$80           ; $d3b6: a9 80     
            sta $78            ; $d3b8: 85 78     
            cpx $78            ; $d3ba: e4 78     
            jsr __f8bf         ; $d3bc: 20 bf f8  
            iny                ; $d3bf: c8        
            lda #$81           ; $d3c0: a9 81     
            sta $78            ; $d3c2: 85 78     
            cpx $78            ; $d3c4: e4 78     
            jsr __f8c9         ; $d3c6: 20 c9 f8  
            iny                ; $d3c9: c8        
            lda #$7f           ; $d3ca: a9 7f     
            sta $78            ; $d3cc: 85 78     
            cpx $78            ; $d3ce: e4 78     
            jsr __f8d3         ; $d3d0: 20 d3 f8  
            iny                ; $d3d3: c8        
            tya                ; $d3d4: 98        
            tax                ; $d3d5: aa        
            lda #$40           ; $d3d6: a9 40     
            sta $78            ; $d3d8: 85 78     
            jsr __f8dd         ; $d3da: 20 dd f8  
            cpy $78            ; $d3dd: c4 78     
            jsr __f8e2         ; $d3df: 20 e2 f8  
            inx                ; $d3e2: e8        
            lda #$3f           ; $d3e3: a9 3f     
            sta $78            ; $d3e5: 85 78     
            jsr __f8ee         ; $d3e7: 20 ee f8  
            cpy $78            ; $d3ea: c4 78     
            jsr __f8f0         ; $d3ec: 20 f0 f8  
            inx                ; $d3ef: e8        
            lda #$41           ; $d3f0: a9 41     
            sta $78            ; $d3f2: 85 78     
            cpy $78            ; $d3f4: c4 78     
            jsr __f8fc         ; $d3f6: 20 fc f8  
            inx                ; $d3f9: e8        
            lda #$00           ; $d3fa: a9 00     
            sta $78            ; $d3fc: 85 78     
            jsr __f906         ; $d3fe: 20 06 f9  
            cpy $78            ; $d401: c4 78     
            jsr __f909         ; $d403: 20 09 f9  
            inx                ; $d406: e8        
            lda #$80           ; $d407: a9 80     
            sta $78            ; $d409: 85 78     
            cpy $78            ; $d40b: c4 78     
            jsr __f913         ; $d40d: 20 13 f9  
            inx                ; $d410: e8        
            lda #$81           ; $d411: a9 81     
            sta $78            ; $d413: 85 78     
            cpy $78            ; $d415: c4 78     
            jsr __f91d         ; $d417: 20 1d f9  
            inx                ; $d41a: e8        
            lda #$7f           ; $d41b: a9 7f     
            sta $78            ; $d41d: 85 78     
            cpy $78            ; $d41f: c4 78     
            jsr __f927         ; $d421: 20 27 f9  
            inx                ; $d424: e8        
            txa                ; $d425: 8a        
            tay                ; $d426: a8        
            jsr __f990         ; $d427: 20 90 f9  
            sta $78            ; $d42a: 85 78     
            lsr $78            ; $d42c: 46 78     
            lda $78            ; $d42e: a5 78     
            jsr __f99d         ; $d430: 20 9d f9  
            iny                ; $d433: c8        
            sta $78            ; $d434: 85 78     
            lsr $78            ; $d436: 46 78     
            lda $78            ; $d438: a5 78     
            jsr __f9ad         ; $d43a: 20 ad f9  
            iny                ; $d43d: c8        
            jsr __f9bd         ; $d43e: 20 bd f9  
            sta $78            ; $d441: 85 78     
            asl $78            ; $d443: 06 78     
            lda $78            ; $d445: a5 78     
            jsr __f9c3         ; $d447: 20 c3 f9  
            iny                ; $d44a: c8        
            sta $78            ; $d44b: 85 78     
            asl $78            ; $d44d: 06 78     
            lda $78            ; $d44f: a5 78     
            jsr __f9d4         ; $d451: 20 d4 f9  
            iny                ; $d454: c8        
            jsr __f9e4         ; $d455: 20 e4 f9  
            sta $78            ; $d458: 85 78     
            ror $78            ; $d45a: 66 78     
            lda $78            ; $d45c: a5 78     
            jsr __f9ea         ; $d45e: 20 ea f9  
            iny                ; $d461: c8        
            sta $78            ; $d462: 85 78     
            ror $78            ; $d464: 66 78     
            lda $78            ; $d466: a5 78     
            jsr __f9fb         ; $d468: 20 fb f9  
            iny                ; $d46b: c8        
            jsr __fa0a         ; $d46c: 20 0a fa  
            sta $78            ; $d46f: 85 78     
            rol $78            ; $d471: 26 78     
            lda $78            ; $d473: a5 78     
            jsr __fa10         ; $d475: 20 10 fa  
            iny                ; $d478: c8        
            sta $78            ; $d479: 85 78     
            rol $78            ; $d47b: 26 78     
            lda $78            ; $d47d: a5 78     
            jsr __fa21         ; $d47f: 20 21 fa  
            lda #$ff           ; $d482: a9 ff     
            sta $78            ; $d484: 85 78     
            sta $01            ; $d486: 85 01     
            bit $01            ; $d488: 24 01     
            sec                ; $d48a: 38        
            inc $78            ; $d48b: e6 78     
            bne __d49b         ; $d48d: d0 0c     
            bmi __d49b         ; $d48f: 30 0a     
            bvc __d49b         ; $d491: 50 08     
            bcc __d49b         ; $d493: 90 06     
            lda $78            ; $d495: a5 78     
            cmp #$00           ; $d497: c9 00     
            beq __d49f         ; $d499: f0 04     
__d49b:     lda #$ab           ; $d49b: a9 ab     
            sta $00            ; $d49d: 85 00     
__d49f:     lda #$7f           ; $d49f: a9 7f     
            sta $78            ; $d4a1: 85 78     
            clv                ; $d4a3: b8        
            clc                ; $d4a4: 18        
            inc $78            ; $d4a5: e6 78     
            beq __d4b5         ; $d4a7: f0 0c     
            bpl __d4b5         ; $d4a9: 10 0a     
            bvs __d4b5         ; $d4ab: 70 08     
            bcs __d4b5         ; $d4ad: b0 06     
            lda $78            ; $d4af: a5 78     
            cmp #$80           ; $d4b1: c9 80     
            beq __d4b9         ; $d4b3: f0 04     
__d4b5:     lda #$ac           ; $d4b5: a9 ac     
            sta $00            ; $d4b7: 85 00     
__d4b9:     lda #$00           ; $d4b9: a9 00     
            sta $78            ; $d4bb: 85 78     
            bit $01            ; $d4bd: 24 01     
            sec                ; $d4bf: 38        
            dec $78            ; $d4c0: c6 78     
            beq __d4d0         ; $d4c2: f0 0c     
            bpl __d4d0         ; $d4c4: 10 0a     
            bvc __d4d0         ; $d4c6: 50 08     
            bcc __d4d0         ; $d4c8: 90 06     
            lda $78            ; $d4ca: a5 78     
            cmp #$ff           ; $d4cc: c9 ff     
            beq __d4d4         ; $d4ce: f0 04     
__d4d0:     lda #$ad           ; $d4d0: a9 ad     
            sta $00            ; $d4d2: 85 00     
__d4d4:     lda #$80           ; $d4d4: a9 80     
            sta $78            ; $d4d6: 85 78     
            clv                ; $d4d8: b8        
            clc                ; $d4d9: 18        
            dec $78            ; $d4da: c6 78     
            beq __d4ea         ; $d4dc: f0 0c     
            bmi __d4ea         ; $d4de: 30 0a     
            bvs __d4ea         ; $d4e0: 70 08     
            bcs __d4ea         ; $d4e2: b0 06     
            lda $78            ; $d4e4: a5 78     
            cmp #$7f           ; $d4e6: c9 7f     
            beq __d4ee         ; $d4e8: f0 04     
__d4ea:     lda #$ae           ; $d4ea: a9 ae     
            sta $00            ; $d4ec: 85 00     
__d4ee:     lda #$01           ; $d4ee: a9 01     
            sta $78            ; $d4f0: 85 78     
            dec $78            ; $d4f2: c6 78     
            beq __d4fa         ; $d4f4: f0 04     
            lda #$af           ; $d4f6: a9 af     
            sta $00            ; $d4f8: 85 00     
__d4fa:     rts                ; $d4fa: 60        

;-------------------------------------------------------------------------------
__d4fb:     lda #$55           ; $d4fb: a9 55     
            sta $0678          ; $d4fd: 8d 78 06  
            lda #$ff           ; $d500: a9 ff     
            sta $01            ; $d502: 85 01     
            bit $01            ; $d504: 24 01     
            ldy #$11           ; $d506: a0 11     
            ldx #$23           ; $d508: a2 23     
            lda #$00           ; $d50a: a9 00     
            lda $0678          ; $d50c: ad 78 06  
            beq __d521         ; $d50f: f0 10     
            bmi __d521         ; $d511: 30 0e     
            cmp #$55           ; $d513: c9 55     
            bne __d521         ; $d515: d0 0a     
            cpy #$11           ; $d517: c0 11     
            bne __d521         ; $d519: d0 06     
            cpx #$23           ; $d51b: e0 23     
            bvc __d521         ; $d51d: 50 02     
            beq __d525         ; $d51f: f0 04     
__d521:     lda #$b0           ; $d521: a9 b0     
            sta $00            ; $d523: 85 00     
__d525:     lda #$46           ; $d525: a9 46     
            bit $01            ; $d527: 24 01     
            sta $0678          ; $d529: 8d 78 06  
            beq __d539         ; $d52c: f0 0b     
            bpl __d539         ; $d52e: 10 09     
            bvc __d539         ; $d530: 50 07     
            lda $0678          ; $d532: ad 78 06  
            cmp #$46           ; $d535: c9 46     
            beq __d53d         ; $d537: f0 04     
__d539:     lda #$b1           ; $d539: a9 b1     
            sta $00            ; $d53b: 85 00     
__d53d:     lda #$55           ; $d53d: a9 55     
            sta $0678          ; $d53f: 8d 78 06  
            bit $01            ; $d542: 24 01     
            lda #$11           ; $d544: a9 11     
            ldx #$23           ; $d546: a2 23     
            ldy #$00           ; $d548: a0 00     
            ldy $0678          ; $d54a: ac 78 06  
            beq __d55f         ; $d54d: f0 10     
            bmi __d55f         ; $d54f: 30 0e     
            cpy #$55           ; $d551: c0 55     
            bne __d55f         ; $d553: d0 0a     
            cmp #$11           ; $d555: c9 11     
            bne __d55f         ; $d557: d0 06     
            cpx #$23           ; $d559: e0 23     
            bvc __d55f         ; $d55b: 50 02     
            beq __d563         ; $d55d: f0 04     
__d55f:     lda #$b2           ; $d55f: a9 b2     
            sta $00            ; $d561: 85 00     
__d563:     ldy #$46           ; $d563: a0 46     
            bit $01            ; $d565: 24 01     
            sty $0678          ; $d567: 8c 78 06  
            beq __d577         ; $d56a: f0 0b     
            bpl __d577         ; $d56c: 10 09     
            bvc __d577         ; $d56e: 50 07     
            ldy $0678          ; $d570: ac 78 06  
            cpy #$46           ; $d573: c0 46     
            beq __d57b         ; $d575: f0 04     
__d577:     lda #$b3           ; $d577: a9 b3     
            sta $00            ; $d579: 85 00     
__d57b:     bit $01            ; $d57b: 24 01     
            lda #$55           ; $d57d: a9 55     
            sta $0678          ; $d57f: 8d 78 06  
            ldy #$11           ; $d582: a0 11     
            lda #$23           ; $d584: a9 23     
            ldx #$00           ; $d586: a2 00     
            ldx $0678          ; $d588: ae 78 06  
            beq __d59d         ; $d58b: f0 10     
            bmi __d59d         ; $d58d: 30 0e     
            cpx #$55           ; $d58f: e0 55     
            bne __d59d         ; $d591: d0 0a     
            cpy #$11           ; $d593: c0 11     
            bne __d59d         ; $d595: d0 06     
            cmp #$23           ; $d597: c9 23     
            bvc __d59d         ; $d599: 50 02     
            beq __d5a1         ; $d59b: f0 04     
__d59d:     lda #$b4           ; $d59d: a9 b4     
            sta $00            ; $d59f: 85 00     
__d5a1:     ldx #$46           ; $d5a1: a2 46     
            bit $01            ; $d5a3: 24 01     
            stx $0678          ; $d5a5: 8e 78 06  
            beq __d5b5         ; $d5a8: f0 0b     
            bpl __d5b5         ; $d5aa: 10 09     
            bvc __d5b5         ; $d5ac: 50 07     
            ldx $0678          ; $d5ae: ae 78 06  
            cpx #$46           ; $d5b1: e0 46     
            beq __d5b9         ; $d5b3: f0 04     
__d5b5:     lda #$b5           ; $d5b5: a9 b5     
            sta $00            ; $d5b7: 85 00     
__d5b9:     lda #$c0           ; $d5b9: a9 c0     
            sta $0678          ; $d5bb: 8d 78 06  
            ldx #$33           ; $d5be: a2 33     
            ldy #$88           ; $d5c0: a0 88     
            lda #$05           ; $d5c2: a9 05     
            bit $0678          ; $d5c4: 2c 78 06  
            bpl __d5d9         ; $d5c7: 10 10     
            bvc __d5d9         ; $d5c9: 50 0e     
            bne __d5d9         ; $d5cb: d0 0c     
            cmp #$05           ; $d5cd: c9 05     
            bne __d5d9         ; $d5cf: d0 08     
            cpx #$33           ; $d5d1: e0 33     
            bne __d5d9         ; $d5d3: d0 04     
            cpy #$88           ; $d5d5: c0 88     
            beq __d5dd         ; $d5d7: f0 04     
__d5d9:     lda #$b6           ; $d5d9: a9 b6     
            sta $00            ; $d5db: 85 00     
__d5dd:     lda #$03           ; $d5dd: a9 03     
            sta $0678          ; $d5df: 8d 78 06  
            lda #$01           ; $d5e2: a9 01     
            bit $0678          ; $d5e4: 2c 78 06  
            bmi __d5f1         ; $d5e7: 30 08     
            bvs __d5f1         ; $d5e9: 70 06     
            beq __d5f1         ; $d5eb: f0 04     
            cmp #$01           ; $d5ed: c9 01     
            beq __d5f5         ; $d5ef: f0 04     
__d5f1:     lda #$b7           ; $d5f1: a9 b7     
            sta $00            ; $d5f3: 85 00     
__d5f5:     ldy #$b8           ; $d5f5: a0 b8     
            lda #$aa           ; $d5f7: a9 aa     
            sta $0678          ; $d5f9: 8d 78 06  
            jsr __f7b6         ; $d5fc: 20 b6 f7  
            ora $0678          ; $d5ff: 0d 78 06  
            jsr __f7c0         ; $d602: 20 c0 f7  
            iny                ; $d605: c8        
            lda #$00           ; $d606: a9 00     
            sta $0678          ; $d608: 8d 78 06  
            jsr __f7ce         ; $d60b: 20 ce f7  
            ora $0678          ; $d60e: 0d 78 06  
            jsr __f7d3         ; $d611: 20 d3 f7  
            iny                ; $d614: c8        
            lda #$aa           ; $d615: a9 aa     
            sta $0678          ; $d617: 8d 78 06  
            jsr __f7df         ; $d61a: 20 df f7  
            and $0678          ; $d61d: 2d 78 06  
            jsr __f7e5         ; $d620: 20 e5 f7  
            iny                ; $d623: c8        
            lda #$ef           ; $d624: a9 ef     
            sta $0678          ; $d626: 8d 78 06  
            jsr __f7f1         ; $d629: 20 f1 f7  
            and $0678          ; $d62c: 2d 78 06  
            jsr __f7f6         ; $d62f: 20 f6 f7  
            iny                ; $d632: c8        
            lda #$aa           ; $d633: a9 aa     
            sta $0678          ; $d635: 8d 78 06  
            jsr __f804         ; $d638: 20 04 f8  
            eor $0678          ; $d63b: 4d 78 06  
            jsr __f80a         ; $d63e: 20 0a f8  
            iny                ; $d641: c8        
            lda #$70           ; $d642: a9 70     
            sta $0678          ; $d644: 8d 78 06  
            jsr __f818         ; $d647: 20 18 f8  
            eor $0678          ; $d64a: 4d 78 06  
            jsr __f81d         ; $d64d: 20 1d f8  
            iny                ; $d650: c8        
            lda #$69           ; $d651: a9 69     
            sta $0678          ; $d653: 8d 78 06  
            jsr __f829         ; $d656: 20 29 f8  
            adc $0678          ; $d659: 6d 78 06  
            jsr __f82f         ; $d65c: 20 2f f8  
            iny                ; $d65f: c8        
            jsr __f83d         ; $d660: 20 3d f8  
            adc $0678          ; $d663: 6d 78 06  
            jsr __f843         ; $d666: 20 43 f8  
            iny                ; $d669: c8        
            lda #$7f           ; $d66a: a9 7f     
            sta $0678          ; $d66c: 8d 78 06  
            jsr __f851         ; $d66f: 20 51 f8  
            adc $0678          ; $d672: 6d 78 06  
            jsr __f856         ; $d675: 20 56 f8  
            iny                ; $d678: c8        
            lda #$80           ; $d679: a9 80     
            sta $0678          ; $d67b: 8d 78 06  
            jsr __f864         ; $d67e: 20 64 f8  
            adc $0678          ; $d681: 6d 78 06  
            jsr __f86a         ; $d684: 20 6a f8  
            iny                ; $d687: c8        
            jsr __f878         ; $d688: 20 78 f8  
            adc $0678          ; $d68b: 6d 78 06  
            jsr __f87d         ; $d68e: 20 7d f8  
            iny                ; $d691: c8        
            lda #$40           ; $d692: a9 40     
            sta $0678          ; $d694: 8d 78 06  
            jsr __f889         ; $d697: 20 89 f8  
            cmp $0678          ; $d69a: cd 78 06  
            jsr __f88e         ; $d69d: 20 8e f8  
            iny                ; $d6a0: c8        
            pha                ; $d6a1: 48        
            lda #$3f           ; $d6a2: a9 3f     
            sta $0678          ; $d6a4: 8d 78 06  
            pla                ; $d6a7: 68        
            jsr __f89a         ; $d6a8: 20 9a f8  
            cmp $0678          ; $d6ab: cd 78 06  
            jsr __f89c         ; $d6ae: 20 9c f8  
            iny                ; $d6b1: c8        
            pha                ; $d6b2: 48        
            lda #$41           ; $d6b3: a9 41     
            sta $0678          ; $d6b5: 8d 78 06  
            pla                ; $d6b8: 68        
            cmp $0678          ; $d6b9: cd 78 06  
            jsr __f8a8         ; $d6bc: 20 a8 f8  
            iny                ; $d6bf: c8        
            pha                ; $d6c0: 48        
            lda #$00           ; $d6c1: a9 00     
            sta $0678          ; $d6c3: 8d 78 06  
            pla                ; $d6c6: 68        
            jsr __f8b2         ; $d6c7: 20 b2 f8  
            cmp $0678          ; $d6ca: cd 78 06  
            jsr __f8b5         ; $d6cd: 20 b5 f8  
            iny                ; $d6d0: c8        
            pha                ; $d6d1: 48        
            lda #$80           ; $d6d2: a9 80     
            sta $0678          ; $d6d4: 8d 78 06  
            pla                ; $d6d7: 68        
            cmp $0678          ; $d6d8: cd 78 06  
            jsr __f8bf         ; $d6db: 20 bf f8  
            iny                ; $d6de: c8        
            pha                ; $d6df: 48        
            lda #$81           ; $d6e0: a9 81     
            sta $0678          ; $d6e2: 8d 78 06  
            pla                ; $d6e5: 68        
            cmp $0678          ; $d6e6: cd 78 06  
            jsr __f8c9         ; $d6e9: 20 c9 f8  
            iny                ; $d6ec: c8        
            pha                ; $d6ed: 48        
            lda #$7f           ; $d6ee: a9 7f     
            sta $0678          ; $d6f0: 8d 78 06  
            pla                ; $d6f3: 68        
            cmp $0678          ; $d6f4: cd 78 06  
            jsr __f8d3         ; $d6f7: 20 d3 f8  
            iny                ; $d6fa: c8        
            lda #$40           ; $d6fb: a9 40     
            sta $0678          ; $d6fd: 8d 78 06  
            jsr __f931         ; $d700: 20 31 f9  
            sbc $0678          ; $d703: ed 78 06  
            jsr __f937         ; $d706: 20 37 f9  
            iny                ; $d709: c8        
            lda #$3f           ; $d70a: a9 3f     
            sta $0678          ; $d70c: 8d 78 06  
            jsr __f947         ; $d70f: 20 47 f9  
            sbc $0678          ; $d712: ed 78 06  
            jsr __f94c         ; $d715: 20 4c f9  
            iny                ; $d718: c8        
            lda #$41           ; $d719: a9 41     
            sta $0678          ; $d71b: 8d 78 06  
            jsr __f95c         ; $d71e: 20 5c f9  
            sbc $0678          ; $d721: ed 78 06  
            jsr __f962         ; $d724: 20 62 f9  
            iny                ; $d727: c8        
            lda #$00           ; $d728: a9 00     
            sta $0678          ; $d72a: 8d 78 06  
            jsr __f972         ; $d72d: 20 72 f9  
            sbc $0678          ; $d730: ed 78 06  
            jsr __f976         ; $d733: 20 76 f9  
            iny                ; $d736: c8        
            lda #$7f           ; $d737: a9 7f     
            sta $0678          ; $d739: 8d 78 06  
            jsr __f980         ; $d73c: 20 80 f9  
            sbc $0678          ; $d73f: ed 78 06  
            jsr __f984         ; $d742: 20 84 f9  
            iny                ; $d745: c8        
            lda #$40           ; $d746: a9 40     
            sta $0678          ; $d748: 8d 78 06  
            jsr __f889         ; $d74b: 20 89 f8  
            tax                ; $d74e: aa        
            cpx $0678          ; $d74f: ec 78 06  
            jsr __f88e         ; $d752: 20 8e f8  
            iny                ; $d755: c8        
            lda #$3f           ; $d756: a9 3f     
            sta $0678          ; $d758: 8d 78 06  
            jsr __f89a         ; $d75b: 20 9a f8  
            cpx $0678          ; $d75e: ec 78 06  
            jsr __f89c         ; $d761: 20 9c f8  
            iny                ; $d764: c8        
            lda #$41           ; $d765: a9 41     
            sta $0678          ; $d767: 8d 78 06  
            cpx $0678          ; $d76a: ec 78 06  
            jsr __f8a8         ; $d76d: 20 a8 f8  
            iny                ; $d770: c8        
            lda #$00           ; $d771: a9 00     
            sta $0678          ; $d773: 8d 78 06  
            jsr __f8b2         ; $d776: 20 b2 f8  
            tax                ; $d779: aa        
            cpx $0678          ; $d77a: ec 78 06  
            jsr __f8b5         ; $d77d: 20 b5 f8  
            iny                ; $d780: c8        
            lda #$80           ; $d781: a9 80     
            sta $0678          ; $d783: 8d 78 06  
            cpx $0678          ; $d786: ec 78 06  
            jsr __f8bf         ; $d789: 20 bf f8  
            iny                ; $d78c: c8        
            lda #$81           ; $d78d: a9 81     
            sta $0678          ; $d78f: 8d 78 06  
            cpx $0678          ; $d792: ec 78 06  
            jsr __f8c9         ; $d795: 20 c9 f8  
            iny                ; $d798: c8        
            lda #$7f           ; $d799: a9 7f     
            sta $0678          ; $d79b: 8d 78 06  
            cpx $0678          ; $d79e: ec 78 06  
            jsr __f8d3         ; $d7a1: 20 d3 f8  
            iny                ; $d7a4: c8        
            tya                ; $d7a5: 98        
            tax                ; $d7a6: aa        
            lda #$40           ; $d7a7: a9 40     
            sta $0678          ; $d7a9: 8d 78 06  
            jsr __f8dd         ; $d7ac: 20 dd f8  
            cpy $0678          ; $d7af: cc 78 06  
            jsr __f8e2         ; $d7b2: 20 e2 f8  
            inx                ; $d7b5: e8        
            lda #$3f           ; $d7b6: a9 3f     
            sta $0678          ; $d7b8: 8d 78 06  
            jsr __f8ee         ; $d7bb: 20 ee f8  
            cpy $0678          ; $d7be: cc 78 06  
            jsr __f8f0         ; $d7c1: 20 f0 f8  
            inx                ; $d7c4: e8        
            lda #$41           ; $d7c5: a9 41     
            sta $0678          ; $d7c7: 8d 78 06  
            cpy $0678          ; $d7ca: cc 78 06  
            jsr __f8fc         ; $d7cd: 20 fc f8  
            inx                ; $d7d0: e8        
            lda #$00           ; $d7d1: a9 00     
            sta $0678          ; $d7d3: 8d 78 06  
            jsr __f906         ; $d7d6: 20 06 f9  
            cpy $0678          ; $d7d9: cc 78 06  
            jsr __f909         ; $d7dc: 20 09 f9  
            inx                ; $d7df: e8        
            lda #$80           ; $d7e0: a9 80     
            sta $0678          ; $d7e2: 8d 78 06  
            cpy $0678          ; $d7e5: cc 78 06  
            jsr __f913         ; $d7e8: 20 13 f9  
            inx                ; $d7eb: e8        
            lda #$81           ; $d7ec: a9 81     
            sta $0678          ; $d7ee: 8d 78 06  
            cpy $0678          ; $d7f1: cc 78 06  
            jsr __f91d         ; $d7f4: 20 1d f9  
            inx                ; $d7f7: e8        
            lda #$7f           ; $d7f8: a9 7f     
            sta $0678          ; $d7fa: 8d 78 06  
            cpy $0678          ; $d7fd: cc 78 06  
            jsr __f927         ; $d800: 20 27 f9  
            inx                ; $d803: e8        
            txa                ; $d804: 8a        
            tay                ; $d805: a8        
            jsr __f990         ; $d806: 20 90 f9  
            sta $0678          ; $d809: 8d 78 06  
            lsr $0678          ; $d80c: 4e 78 06  
            lda $0678          ; $d80f: ad 78 06  
            jsr __f99d         ; $d812: 20 9d f9  
            iny                ; $d815: c8        
            sta $0678          ; $d816: 8d 78 06  
            lsr $0678          ; $d819: 4e 78 06  
            lda $0678          ; $d81c: ad 78 06  
            jsr __f9ad         ; $d81f: 20 ad f9  
            iny                ; $d822: c8        
            jsr __f9bd         ; $d823: 20 bd f9  
            sta $0678          ; $d826: 8d 78 06  
            asl $0678          ; $d829: 0e 78 06  
            lda $0678          ; $d82c: ad 78 06  
            jsr __f9c3         ; $d82f: 20 c3 f9  
            iny                ; $d832: c8        
            sta $0678          ; $d833: 8d 78 06  
            asl $0678          ; $d836: 0e 78 06  
            lda $0678          ; $d839: ad 78 06  
            jsr __f9d4         ; $d83c: 20 d4 f9  
            iny                ; $d83f: c8        
            jsr __f9e4         ; $d840: 20 e4 f9  
            sta $0678          ; $d843: 8d 78 06  
            ror $0678          ; $d846: 6e 78 06  
            lda $0678          ; $d849: ad 78 06  
            jsr __f9ea         ; $d84c: 20 ea f9  
            iny                ; $d84f: c8        
            sta $0678          ; $d850: 8d 78 06  
            ror $0678          ; $d853: 6e 78 06  
            lda $0678          ; $d856: ad 78 06  
            jsr __f9fb         ; $d859: 20 fb f9  
            iny                ; $d85c: c8        
            jsr __fa0a         ; $d85d: 20 0a fa  
            sta $0678          ; $d860: 8d 78 06  
            rol $0678          ; $d863: 2e 78 06  
            lda $0678          ; $d866: ad 78 06  
            jsr __fa10         ; $d869: 20 10 fa  
            iny                ; $d86c: c8        
            sta $0678          ; $d86d: 8d 78 06  
            rol $0678          ; $d870: 2e 78 06  
            lda $0678          ; $d873: ad 78 06  
            jsr __fa21         ; $d876: 20 21 fa  
            lda #$ff           ; $d879: a9 ff     
            sta $0678          ; $d87b: 8d 78 06  
            sta $01            ; $d87e: 85 01     
            bit $01            ; $d880: 24 01     
            sec                ; $d882: 38        
            inc $0678          ; $d883: ee 78 06  
            bne __d895         ; $d886: d0 0d     
            bmi __d895         ; $d888: 30 0b     
            bvc __d895         ; $d88a: 50 09     
            bcc __d895         ; $d88c: 90 07     
            lda $0678          ; $d88e: ad 78 06  
            cmp #$00           ; $d891: c9 00     
            beq __d899         ; $d893: f0 04     
__d895:     lda #$e5           ; $d895: a9 e5     
            sta $00            ; $d897: 85 00     
__d899:     lda #$7f           ; $d899: a9 7f     
            sta $0678          ; $d89b: 8d 78 06  
            clv                ; $d89e: b8        
            clc                ; $d89f: 18        
            inc $0678          ; $d8a0: ee 78 06  
            beq __d8b2         ; $d8a3: f0 0d     
            bpl __d8b2         ; $d8a5: 10 0b     
            bvs __d8b2         ; $d8a7: 70 09     
            bcs __d8b2         ; $d8a9: b0 07     
            lda $0678          ; $d8ab: ad 78 06  
            cmp #$80           ; $d8ae: c9 80     
            beq __d8b6         ; $d8b0: f0 04     
__d8b2:     lda #$e6           ; $d8b2: a9 e6     
            sta $00            ; $d8b4: 85 00     
__d8b6:     lda #$00           ; $d8b6: a9 00     
            sta $0678          ; $d8b8: 8d 78 06  
            bit $01            ; $d8bb: 24 01     
            sec                ; $d8bd: 38        
            dec $0678          ; $d8be: ce 78 06  
            beq __d8d0         ; $d8c1: f0 0d     
            bpl __d8d0         ; $d8c3: 10 0b     
            bvc __d8d0         ; $d8c5: 50 09     
            bcc __d8d0         ; $d8c7: 90 07     
            lda $0678          ; $d8c9: ad 78 06  
            cmp #$ff           ; $d8cc: c9 ff     
            beq __d8d4         ; $d8ce: f0 04     
__d8d0:     lda #$e7           ; $d8d0: a9 e7     
            sta $00            ; $d8d2: 85 00     
__d8d4:     lda #$80           ; $d8d4: a9 80     
            sta $0678          ; $d8d6: 8d 78 06  
            clv                ; $d8d9: b8        
            clc                ; $d8da: 18        
            dec $0678          ; $d8db: ce 78 06  
            beq __d8ed         ; $d8de: f0 0d     
            bmi __d8ed         ; $d8e0: 30 0b     
            bvs __d8ed         ; $d8e2: 70 09     
            bcs __d8ed         ; $d8e4: b0 07     
            lda $0678          ; $d8e6: ad 78 06  
            cmp #$7f           ; $d8e9: c9 7f     
            beq __d8f1         ; $d8eb: f0 04     
__d8ed:     lda #$e8           ; $d8ed: a9 e8     
            sta $00            ; $d8ef: 85 00     
__d8f1:     lda #$01           ; $d8f1: a9 01     
            sta $0678          ; $d8f3: 8d 78 06  
            dec $0678          ; $d8f6: ce 78 06  
            beq __d8ff         ; $d8f9: f0 04     
            lda #$e9           ; $d8fb: a9 e9     
            sta $00            ; $d8fd: 85 00     
__d8ff:     rts                ; $d8ff: 60        

;-------------------------------------------------------------------------------
__d900:     lda #$a3           ; $d900: a9 a3     
            sta $33            ; $d902: 85 33     
            lda #$89           ; $d904: a9 89     
            sta $0300          ; $d906: 8d 00 03  
            lda #$12           ; $d909: a9 12     
            sta $0245          ; $d90b: 8d 45 02  
            lda #$ff           ; $d90e: a9 ff     
            sta $01            ; $d910: 85 01     
            ldx #$65           ; $d912: a2 65     
            lda #$00           ; $d914: a9 00     
            sta $89            ; $d916: 85 89     
            lda #$03           ; $d918: a9 03     
            sta $8a            ; $d91a: 85 8a     
            ldy #$00           ; $d91c: a0 00     
            sec                ; $d91e: 38        
            lda #$00           ; $d91f: a9 00     
            clv                ; $d921: b8        
            lda ($89),y        ; $d922: b1 89     
            beq __d932         ; $d924: f0 0c     
            bcc __d932         ; $d926: 90 0a     
            bvs __d932         ; $d928: 70 08     
            cmp #$89           ; $d92a: c9 89     
            bne __d932         ; $d92c: d0 04     
            cpx #$65           ; $d92e: e0 65     
            beq __d936         ; $d930: f0 04     
__d932:     lda #$ea           ; $d932: a9 ea     
            sta $00            ; $d934: 85 00     
__d936:     lda #$ff           ; $d936: a9 ff     
            sta $97            ; $d938: 85 97     
            sta $98            ; $d93a: 85 98     
            bit $98            ; $d93c: 24 98     
            ldy #$34           ; $d93e: a0 34     
            lda ($97),y        ; $d940: b1 97     
            cmp #$a3           ; $d942: c9 a3     
            bne __d948         ; $d944: d0 02     
            bcs __d94c         ; $d946: b0 04     
__d948:     lda #$eb           ; $d948: a9 eb     
            sta $00            ; $d94a: 85 00     
__d94c:     lda $00            ; $d94c: a5 00     
            pha                ; $d94e: 48        
            lda #$46           ; $d94f: a9 46     
            sta $ff            ; $d951: 85 ff     
            lda #$01           ; $d953: a9 01     
            sta $00            ; $d955: 85 00     
            ldy #$ff           ; $d957: a0 ff     
            lda ($ff),y        ; $d959: b1 ff     
            cmp #$12           ; $d95b: c9 12     
            beq __d963         ; $d95d: f0 04     
            pla                ; $d95f: 68        
            lda #$ec           ; $d960: a9 ec     
            pha                ; $d962: 48        
__d963:     pla                ; $d963: 68        
            sta $00            ; $d964: 85 00     
            ldx #$ed           ; $d966: a2 ed     
            lda #$00           ; $d968: a9 00     
            sta $33            ; $d96a: 85 33     
            lda #$04           ; $d96c: a9 04     
            sta $34            ; $d96e: 85 34     
            ldy #$00           ; $d970: a0 00     
            clc                ; $d972: 18        
            lda #$ff           ; $d973: a9 ff     
            sta $01            ; $d975: 85 01     
            bit $01            ; $d977: 24 01     
            lda #$aa           ; $d979: a9 aa     
            sta $0400          ; $d97b: 8d 00 04  
            lda #$55           ; $d97e: a9 55     
            ora ($33),y        ; $d980: 11 33     
            bcs __d98c         ; $d982: b0 08     
            bpl __d98c         ; $d984: 10 06     
            cmp #$ff           ; $d986: c9 ff     
            bne __d98c         ; $d988: d0 02     
            bvs __d98e         ; $d98a: 70 02     
__d98c:     stx $00            ; $d98c: 86 00     
__d98e:     inx                ; $d98e: e8        
            sec                ; $d98f: 38        
            clv                ; $d990: b8        
            lda #$00           ; $d991: a9 00     
            ora ($33),y        ; $d993: 11 33     
            beq __d99d         ; $d995: f0 06     
            bvs __d99d         ; $d997: 70 04     
            bcc __d99d         ; $d999: 90 02     
            bmi __d99f         ; $d99b: 30 02     
__d99d:     stx $00            ; $d99d: 86 00     
__d99f:     inx                ; $d99f: e8        
            clc                ; $d9a0: 18        
            bit $01            ; $d9a1: 24 01     
            lda #$55           ; $d9a3: a9 55     
            and ($33),y        ; $d9a5: 31 33     
            bne __d9af         ; $d9a7: d0 06     
            bvc __d9af         ; $d9a9: 50 04     
            bcs __d9af         ; $d9ab: b0 02     
            bpl __d9b1         ; $d9ad: 10 02     
__d9af:     stx $00            ; $d9af: 86 00     
__d9b1:     inx                ; $d9b1: e8        
            sec                ; $d9b2: 38        
            clv                ; $d9b3: b8        
            lda #$ef           ; $d9b4: a9 ef     
            sta $0400          ; $d9b6: 8d 00 04  
            lda #$f8           ; $d9b9: a9 f8     
            and ($33),y        ; $d9bb: 31 33     
            bcc __d9c7         ; $d9bd: 90 08     
            bpl __d9c7         ; $d9bf: 10 06     
            cmp #$e8           ; $d9c1: c9 e8     
            bne __d9c7         ; $d9c3: d0 02     
            bvc __d9c9         ; $d9c5: 50 02     
__d9c7:     stx $00            ; $d9c7: 86 00     
__d9c9:     inx                ; $d9c9: e8        
            clc                ; $d9ca: 18        
            bit $01            ; $d9cb: 24 01     
            lda #$aa           ; $d9cd: a9 aa     
            sta $0400          ; $d9cf: 8d 00 04  
            lda #$5f           ; $d9d2: a9 5f     
            eor ($33),y        ; $d9d4: 51 33     
            bcs __d9e0         ; $d9d6: b0 08     
            bpl __d9e0         ; $d9d8: 10 06     
            cmp #$f5           ; $d9da: c9 f5     
            bne __d9e0         ; $d9dc: d0 02     
            bvs __d9e2         ; $d9de: 70 02     
__d9e0:     stx $00            ; $d9e0: 86 00     
__d9e2:     inx                ; $d9e2: e8        
            sec                ; $d9e3: 38        
            clv                ; $d9e4: b8        
            lda #$70           ; $d9e5: a9 70     
            sta $0400          ; $d9e7: 8d 00 04  
            eor ($33),y        ; $d9ea: 51 33     
            bne __d9f4         ; $d9ec: d0 06     
            bvs __d9f4         ; $d9ee: 70 04     
            bcc __d9f4         ; $d9f0: 90 02     
            bpl __d9f6         ; $d9f2: 10 02     
__d9f4:     stx $00            ; $d9f4: 86 00     
__d9f6:     inx                ; $d9f6: e8        
            clc                ; $d9f7: 18        
            bit $01            ; $d9f8: 24 01     
            lda #$69           ; $d9fa: a9 69     
            sta $0400          ; $d9fc: 8d 00 04  
            lda #$00           ; $d9ff: a9 00     
            adc ($33),y        ; $da01: 71 33     
            bmi __da0d         ; $da03: 30 08     
            bcs __da0d         ; $da05: b0 06     
            cmp #$69           ; $da07: c9 69     
            bne __da0d         ; $da09: d0 02     
            bvc __da0f         ; $da0b: 50 02     
__da0d:     stx $00            ; $da0d: 86 00     
__da0f:     inx                ; $da0f: e8        
            sec                ; $da10: 38        
            bit $01            ; $da11: 24 01     
            lda #$00           ; $da13: a9 00     
            adc ($33),y        ; $da15: 71 33     
            bmi __da21         ; $da17: 30 08     
            bcs __da21         ; $da19: b0 06     
            cmp #$6a           ; $da1b: c9 6a     
            bne __da21         ; $da1d: d0 02     
            bvc __da23         ; $da1f: 50 02     
__da21:     stx $00            ; $da21: 86 00     
__da23:     inx                ; $da23: e8        
            sec                ; $da24: 38        
            clv                ; $da25: b8        
            lda #$7f           ; $da26: a9 7f     
            sta $0400          ; $da28: 8d 00 04  
            adc ($33),y        ; $da2b: 71 33     
            bpl __da37         ; $da2d: 10 08     
            bcs __da37         ; $da2f: b0 06     
            cmp #$ff           ; $da31: c9 ff     
            bne __da37         ; $da33: d0 02     
            bvs __da39         ; $da35: 70 02     
__da37:     stx $00            ; $da37: 86 00     
__da39:     inx                ; $da39: e8        
            clc                ; $da3a: 18        
            bit $01            ; $da3b: 24 01     
            lda #$80           ; $da3d: a9 80     
            sta $0400          ; $da3f: 8d 00 04  
            lda #$7f           ; $da42: a9 7f     
            adc ($33),y        ; $da44: 71 33     
            bpl __da50         ; $da46: 10 08     
            bcs __da50         ; $da48: b0 06     
            cmp #$ff           ; $da4a: c9 ff     
            bne __da50         ; $da4c: d0 02     
            bvc __da52         ; $da4e: 50 02     
__da50:     stx $00            ; $da50: 86 00     
__da52:     inx                ; $da52: e8        
            sec                ; $da53: 38        
            clv                ; $da54: b8        
            lda #$80           ; $da55: a9 80     
            sta $0400          ; $da57: 8d 00 04  
            lda #$7f           ; $da5a: a9 7f     
            adc ($33),y        ; $da5c: 71 33     
            bne __da66         ; $da5e: d0 06     
            bmi __da66         ; $da60: 30 04     
            bvs __da66         ; $da62: 70 02     
            bcs __da68         ; $da64: b0 02     
__da66:     stx $00            ; $da66: 86 00     
__da68:     inx                ; $da68: e8        
            bit $01            ; $da69: 24 01     
            lda #$40           ; $da6b: a9 40     
            sta $0400          ; $da6d: 8d 00 04  
            cmp ($33),y        ; $da70: d1 33     
            bmi __da7a         ; $da72: 30 06     
            bcc __da7a         ; $da74: 90 04     
            bne __da7a         ; $da76: d0 02     
            bvs __da7c         ; $da78: 70 02     
__da7a:     stx $00            ; $da7a: 86 00     
__da7c:     inx                ; $da7c: e8        
            clv                ; $da7d: b8        
            dec $0400          ; $da7e: ce 00 04  
            cmp ($33),y        ; $da81: d1 33     
            beq __da8b         ; $da83: f0 06     
            bmi __da8b         ; $da85: 30 04     
            bcc __da8b         ; $da87: 90 02     
            bvc __da8d         ; $da89: 50 02     
__da8b:     stx $00            ; $da8b: 86 00     
__da8d:     inx                ; $da8d: e8        
            inc $0400          ; $da8e: ee 00 04  
            inc $0400          ; $da91: ee 00 04  
            cmp ($33),y        ; $da94: d1 33     
            beq __da9a         ; $da96: f0 02     
            bmi __da9c         ; $da98: 30 02     
__da9a:     stx $00            ; $da9a: 86 00     
__da9c:     inx                ; $da9c: e8        
            lda #$00           ; $da9d: a9 00     
            sta $0400          ; $da9f: 8d 00 04  
            lda #$80           ; $daa2: a9 80     
            cmp ($33),y        ; $daa4: d1 33     
            beq __daac         ; $daa6: f0 04     
            bpl __daac         ; $daa8: 10 02     
            bcs __daae         ; $daaa: b0 02     
__daac:     stx $00            ; $daac: 86 00     
__daae:     inx                ; $daae: e8        
            ldy #$80           ; $daaf: a0 80     
            sty $0400          ; $dab1: 8c 00 04  
            ldy #$00           ; $dab4: a0 00     
            cmp ($33),y        ; $dab6: d1 33     
            bne __dabe         ; $dab8: d0 04     
            bmi __dabe         ; $daba: 30 02     
            bcs __dac0         ; $dabc: b0 02     
__dabe:     stx $00            ; $dabe: 86 00     
__dac0:     inx                ; $dac0: e8        
            inc $0400          ; $dac1: ee 00 04  
            cmp ($33),y        ; $dac4: d1 33     
            bcs __dacc         ; $dac6: b0 04     
            beq __dacc         ; $dac8: f0 02     
            bmi __dace         ; $daca: 30 02     
__dacc:     stx $00            ; $dacc: 86 00     
__dace:     inx                ; $dace: e8        
            dec $0400          ; $dacf: ce 00 04  
            dec $0400          ; $dad2: ce 00 04  
            cmp ($33),y        ; $dad5: d1 33     
            bcc __dadd         ; $dad7: 90 04     
            beq __dadd         ; $dad9: f0 02     
            bpl __dadf         ; $dadb: 10 02     
__dadd:     stx $00            ; $dadd: 86 00     
__dadf:     rts                ; $dadf: 60        

;-------------------------------------------------------------------------------
__dae0:     lda #$00           ; $dae0: a9 00     
            sta $33            ; $dae2: 85 33     
            lda #$04           ; $dae4: a9 04     
            sta $34            ; $dae6: 85 34     
            ldy #$00           ; $dae8: a0 00     
            ldx #$01           ; $daea: a2 01     
            bit $01            ; $daec: 24 01     
            lda #$40           ; $daee: a9 40     
            sta $0400          ; $daf0: 8d 00 04  
            sec                ; $daf3: 38        
            sbc ($33),y        ; $daf4: f1 33     
            bmi __db02         ; $daf6: 30 0a     
            bcc __db02         ; $daf8: 90 08     
            bne __db02         ; $dafa: d0 06     
            bvs __db02         ; $dafc: 70 04     
            cmp #$00           ; $dafe: c9 00     
            beq __db04         ; $db00: f0 02     
__db02:     stx $00            ; $db02: 86 00     
__db04:     inx                ; $db04: e8        
            clv                ; $db05: b8        
            sec                ; $db06: 38        
            lda #$40           ; $db07: a9 40     
            dec $0400          ; $db09: ce 00 04  
            sbc ($33),y        ; $db0c: f1 33     
            beq __db1a         ; $db0e: f0 0a     
            bmi __db1a         ; $db10: 30 08     
            bcc __db1a         ; $db12: 90 06     
            bvs __db1a         ; $db14: 70 04     
            cmp #$01           ; $db16: c9 01     
            beq __db1c         ; $db18: f0 02     
__db1a:     stx $00            ; $db1a: 86 00     
__db1c:     inx                ; $db1c: e8        
            lda #$40           ; $db1d: a9 40     
            sec                ; $db1f: 38        
            bit $01            ; $db20: 24 01     
            inc $0400          ; $db22: ee 00 04  
            inc $0400          ; $db25: ee 00 04  
            sbc ($33),y        ; $db28: f1 33     
            bcs __db36         ; $db2a: b0 0a     
            beq __db36         ; $db2c: f0 08     
            bpl __db36         ; $db2e: 10 06     
            bvs __db36         ; $db30: 70 04     
            cmp #$ff           ; $db32: c9 ff     
            beq __db38         ; $db34: f0 02     
__db36:     stx $00            ; $db36: 86 00     
__db38:     inx                ; $db38: e8        
            clc                ; $db39: 18        
            lda #$00           ; $db3a: a9 00     
            sta $0400          ; $db3c: 8d 00 04  
            lda #$80           ; $db3f: a9 80     
            sbc ($33),y        ; $db41: f1 33     
            bcc __db49         ; $db43: 90 04     
            cmp #$7f           ; $db45: c9 7f     
            beq __db4b         ; $db47: f0 02     
__db49:     stx $00            ; $db49: 86 00     
__db4b:     inx                ; $db4b: e8        
            sec                ; $db4c: 38        
            lda #$7f           ; $db4d: a9 7f     
            sta $0400          ; $db4f: 8d 00 04  
            lda #$81           ; $db52: a9 81     
            sbc ($33),y        ; $db54: f1 33     
            bvc __db5e         ; $db56: 50 06     
            bcc __db5e         ; $db58: 90 04     
            cmp #$02           ; $db5a: c9 02     
            beq __db60         ; $db5c: f0 02     
__db5e:     stx $00            ; $db5e: 86 00     
__db60:     inx                ; $db60: e8        
            lda #$00           ; $db61: a9 00     
            lda #$87           ; $db63: a9 87     
            sta ($33),y        ; $db65: 91 33     
            lda $0400          ; $db67: ad 00 04  
            cmp #$87           ; $db6a: c9 87     
            beq __db70         ; $db6c: f0 02     
            stx $00            ; $db6e: 86 00     
__db70:     inx                ; $db70: e8        
            lda #$7e           ; $db71: a9 7e     
            sta $0200          ; $db73: 8d 00 02  
            lda #$db           ; $db76: a9 db     
            sta $0201          ; $db78: 8d 01 02  
            jmp ($0200)        ; $db7b: 6c 00 02  

;-------------------------------------------------------------------------------
            lda #$00           ; $db7e: a9 00     
            sta $02ff          ; $db80: 8d ff 02  
            lda #$01           ; $db83: a9 01     
            sta $0300          ; $db85: 8d 00 03  
            lda #$03           ; $db88: a9 03     
            sta $0200          ; $db8a: 8d 00 02  
            lda #$a9           ; $db8d: a9 a9     
            sta $0100          ; $db8f: 8d 00 01  
            lda #$55           ; $db92: a9 55     
            sta $0101          ; $db94: 8d 01 01  
            lda #$60           ; $db97: a9 60     
            sta $0102          ; $db99: 8d 02 01  
            lda #$a9           ; $db9c: a9 a9     
            sta $0300          ; $db9e: 8d 00 03  
            lda #$aa           ; $dba1: a9 aa     
            sta $0301          ; $dba3: 8d 01 03  
            lda #$60           ; $dba6: a9 60     
            sta $0302          ; $dba8: 8d 02 03  
            jsr __dbb5         ; $dbab: 20 b5 db  
            cmp #$aa           ; $dbae: c9 aa     
            beq __dbb4         ; $dbb0: f0 02     
            stx $00            ; $dbb2: 86 00     
__dbb4:     rts                ; $dbb4: 60        

;-------------------------------------------------------------------------------
__dbb5:     jmp ($02ff)        ; $dbb5: 6c ff 02  

;-------------------------------------------------------------------------------
__dbb8:     lda #$ff           ; $dbb8: a9 ff     
            sta $01            ; $dbba: 85 01     
            lda #$aa           ; $dbbc: a9 aa     
            sta $33            ; $dbbe: 85 33     
            lda #$bb           ; $dbc0: a9 bb     
            sta $89            ; $dbc2: 85 89     
            ldx #$00           ; $dbc4: a2 00     
            lda #$66           ; $dbc6: a9 66     
            bit $01            ; $dbc8: 24 01     
            sec                ; $dbca: 38        
            ldy #$00           ; $dbcb: a0 00     
            ldy $33,x          ; $dbcd: b4 33     
            bpl __dbe3         ; $dbcf: 10 12     
            beq __dbe3         ; $dbd1: f0 10     
            bvc __dbe3         ; $dbd3: 50 0e     
            bcc __dbe3         ; $dbd5: 90 0c     
            cmp #$66           ; $dbd7: c9 66     
            bne __dbe3         ; $dbd9: d0 08     
            cpx #$00           ; $dbdb: e0 00     
            bne __dbe3         ; $dbdd: d0 04     
            cpy #$aa           ; $dbdf: c0 aa     
            beq __dbe7         ; $dbe1: f0 04     
__dbe3:     lda #$08           ; $dbe3: a9 08     
            sta $00            ; $dbe5: 85 00     
__dbe7:     ldx #$8a           ; $dbe7: a2 8a     
            lda #$66           ; $dbe9: a9 66     
            clv                ; $dbeb: b8        
            clc                ; $dbec: 18        
            ldy #$00           ; $dbed: a0 00     
            ldy $ff,x          ; $dbef: b4 ff     
            bpl __dc05         ; $dbf1: 10 12     
            beq __dc05         ; $dbf3: f0 10     
            bvs __dc05         ; $dbf5: 70 0e     
            bcs __dc05         ; $dbf7: b0 0c     
            cpy #$bb           ; $dbf9: c0 bb     
            bne __dc05         ; $dbfb: d0 08     
            cmp #$66           ; $dbfd: c9 66     
            bne __dc05         ; $dbff: d0 04     
            cpx #$8a           ; $dc01: e0 8a     
            beq __dc09         ; $dc03: f0 04     
__dc05:     lda #$09           ; $dc05: a9 09     
            sta $00            ; $dc07: 85 00     
__dc09:     bit $01            ; $dc09: 24 01     
            sec                ; $dc0b: 38        
            ldy #$44           ; $dc0c: a0 44     
            ldx #$00           ; $dc0e: a2 00     
            sty $33,x          ; $dc10: 94 33     
            lda $33            ; $dc12: a5 33     
            bcc __dc2e         ; $dc14: 90 18     
            cmp #$44           ; $dc16: c9 44     
            bne __dc2e         ; $dc18: d0 14     
            bvc __dc2e         ; $dc1a: 50 12     
            clc                ; $dc1c: 18        
            clv                ; $dc1d: b8        
            ldy #$99           ; $dc1e: a0 99     
            ldx #$80           ; $dc20: a2 80     
            sty $85,x          ; $dc22: 94 85     
            lda $05            ; $dc24: a5 05     
            bcs __dc2e         ; $dc26: b0 06     
            cmp #$99           ; $dc28: c9 99     
            bne __dc2e         ; $dc2a: d0 02     
            bvc __dc32         ; $dc2c: 50 04     
__dc2e:     lda #$0a           ; $dc2e: a9 0a     
            sta $00            ; $dc30: 85 00     
__dc32:     ldy #$0b           ; $dc32: a0 0b     
            lda #$aa           ; $dc34: a9 aa     
            ldx #$78           ; $dc36: a2 78     
            sta $78            ; $dc38: 85 78     
            jsr __f7b6         ; $dc3a: 20 b6 f7  
            ora $00,x          ; $dc3d: 15 00     
            jsr __f7c0         ; $dc3f: 20 c0 f7  
            iny                ; $dc42: c8        
            lda #$00           ; $dc43: a9 00     
            sta $78            ; $dc45: 85 78     
            jsr __f7ce         ; $dc47: 20 ce f7  
            ora $00,x          ; $dc4a: 15 00     
            jsr __f7d3         ; $dc4c: 20 d3 f7  
            iny                ; $dc4f: c8        
            lda #$aa           ; $dc50: a9 aa     
            sta $78            ; $dc52: 85 78     
            jsr __f7df         ; $dc54: 20 df f7  
            and $00,x          ; $dc57: 35 00     
            jsr __f7e5         ; $dc59: 20 e5 f7  
            iny                ; $dc5c: c8        
            lda #$ef           ; $dc5d: a9 ef     
            sta $78            ; $dc5f: 85 78     
            jsr __f7f1         ; $dc61: 20 f1 f7  
            and $00,x          ; $dc64: 35 00     
            jsr __f7f6         ; $dc66: 20 f6 f7  
            iny                ; $dc69: c8        
            lda #$aa           ; $dc6a: a9 aa     
            sta $78            ; $dc6c: 85 78     
            jsr __f804         ; $dc6e: 20 04 f8  
            eor $00,x          ; $dc71: 55 00     
            jsr __f80a         ; $dc73: 20 0a f8  
            iny                ; $dc76: c8        
            lda #$70           ; $dc77: a9 70     
            sta $78            ; $dc79: 85 78     
            jsr __f818         ; $dc7b: 20 18 f8  
            eor $00,x          ; $dc7e: 55 00     
            jsr __f81d         ; $dc80: 20 1d f8  
            iny                ; $dc83: c8        
            lda #$69           ; $dc84: a9 69     
            sta $78            ; $dc86: 85 78     
            jsr __f829         ; $dc88: 20 29 f8  
            adc $00,x          ; $dc8b: 75 00     
            jsr __f82f         ; $dc8d: 20 2f f8  
            iny                ; $dc90: c8        
            jsr __f83d         ; $dc91: 20 3d f8  
            adc $00,x          ; $dc94: 75 00     
            jsr __f843         ; $dc96: 20 43 f8  
            iny                ; $dc99: c8        
            lda #$7f           ; $dc9a: a9 7f     
            sta $78            ; $dc9c: 85 78     
            jsr __f851         ; $dc9e: 20 51 f8  
            adc $00,x          ; $dca1: 75 00     
            jsr __f856         ; $dca3: 20 56 f8  
            iny                ; $dca6: c8        
            lda #$80           ; $dca7: a9 80     
            sta $78            ; $dca9: 85 78     
            jsr __f864         ; $dcab: 20 64 f8  
            adc $00,x          ; $dcae: 75 00     
            jsr __f86a         ; $dcb0: 20 6a f8  
            iny                ; $dcb3: c8        
            jsr __f878         ; $dcb4: 20 78 f8  
            adc $00,x          ; $dcb7: 75 00     
            jsr __f87d         ; $dcb9: 20 7d f8  
            iny                ; $dcbc: c8        
            lda #$40           ; $dcbd: a9 40     
            sta $78            ; $dcbf: 85 78     
            jsr __f889         ; $dcc1: 20 89 f8  
            cmp $00,x          ; $dcc4: d5 00     
            jsr __f88e         ; $dcc6: 20 8e f8  
            iny                ; $dcc9: c8        
            pha                ; $dcca: 48        
            lda #$3f           ; $dccb: a9 3f     
            sta $78            ; $dccd: 85 78     
            pla                ; $dccf: 68        
            jsr __f89a         ; $dcd0: 20 9a f8  
            cmp $00,x          ; $dcd3: d5 00     
            jsr __f89c         ; $dcd5: 20 9c f8  
            iny                ; $dcd8: c8        
            pha                ; $dcd9: 48        
            lda #$41           ; $dcda: a9 41     
            sta $78            ; $dcdc: 85 78     
            pla                ; $dcde: 68        
            cmp $00,x          ; $dcdf: d5 00     
            jsr __f8a8         ; $dce1: 20 a8 f8  
            iny                ; $dce4: c8        
            pha                ; $dce5: 48        
            lda #$00           ; $dce6: a9 00     
            sta $78            ; $dce8: 85 78     
            pla                ; $dcea: 68        
            jsr __f8b2         ; $dceb: 20 b2 f8  
            cmp $00,x          ; $dcee: d5 00     
            jsr __f8b5         ; $dcf0: 20 b5 f8  
            iny                ; $dcf3: c8        
            pha                ; $dcf4: 48        
            lda #$80           ; $dcf5: a9 80     
            sta $78            ; $dcf7: 85 78     
            pla                ; $dcf9: 68        
            cmp $00,x          ; $dcfa: d5 00     
            jsr __f8bf         ; $dcfc: 20 bf f8  
            iny                ; $dcff: c8        
            pha                ; $dd00: 48        
            lda #$81           ; $dd01: a9 81     
            sta $78            ; $dd03: 85 78     
            pla                ; $dd05: 68        
            cmp $00,x          ; $dd06: d5 00     
            jsr __f8c9         ; $dd08: 20 c9 f8  
            iny                ; $dd0b: c8        
            pha                ; $dd0c: 48        
            lda #$7f           ; $dd0d: a9 7f     
            sta $78            ; $dd0f: 85 78     
            pla                ; $dd11: 68        
            cmp $00,x          ; $dd12: d5 00     
            jsr __f8d3         ; $dd14: 20 d3 f8  
            iny                ; $dd17: c8        
            lda #$40           ; $dd18: a9 40     
            sta $78            ; $dd1a: 85 78     
            jsr __f931         ; $dd1c: 20 31 f9  
            sbc $00,x          ; $dd1f: f5 00     
            jsr __f937         ; $dd21: 20 37 f9  
            iny                ; $dd24: c8        
            lda #$3f           ; $dd25: a9 3f     
            sta $78            ; $dd27: 85 78     
            jsr __f947         ; $dd29: 20 47 f9  
            sbc $00,x          ; $dd2c: f5 00     
            jsr __f94c         ; $dd2e: 20 4c f9  
            iny                ; $dd31: c8        
            lda #$41           ; $dd32: a9 41     
            sta $78            ; $dd34: 85 78     
            jsr __f95c         ; $dd36: 20 5c f9  
            sbc $00,x          ; $dd39: f5 00     
            jsr __f962         ; $dd3b: 20 62 f9  
            iny                ; $dd3e: c8        
            lda #$00           ; $dd3f: a9 00     
            sta $78            ; $dd41: 85 78     
            jsr __f972         ; $dd43: 20 72 f9  
            sbc $00,x          ; $dd46: f5 00     
            jsr __f976         ; $dd48: 20 76 f9  
            iny                ; $dd4b: c8        
            lda #$7f           ; $dd4c: a9 7f     
            sta $78            ; $dd4e: 85 78     
            jsr __f980         ; $dd50: 20 80 f9  
            sbc $00,x          ; $dd53: f5 00     
            jsr __f984         ; $dd55: 20 84 f9  
            lda #$aa           ; $dd58: a9 aa     
            sta $33            ; $dd5a: 85 33     
            lda #$bb           ; $dd5c: a9 bb     
            sta $89            ; $dd5e: 85 89     
            ldx #$00           ; $dd60: a2 00     
            ldy #$66           ; $dd62: a0 66     
            bit $01            ; $dd64: 24 01     
            sec                ; $dd66: 38        
            lda #$00           ; $dd67: a9 00     
            lda $33,x          ; $dd69: b5 33     
            bpl __dd7f         ; $dd6b: 10 12     
            beq __dd7f         ; $dd6d: f0 10     
            bvc __dd7f         ; $dd6f: 50 0e     
            bcc __dd7f         ; $dd71: 90 0c     
            cpy #$66           ; $dd73: c0 66     
            bne __dd7f         ; $dd75: d0 08     
            cpx #$00           ; $dd77: e0 00     
            bne __dd7f         ; $dd79: d0 04     
            cmp #$aa           ; $dd7b: c9 aa     
            beq __dd83         ; $dd7d: f0 04     
__dd7f:     lda #$22           ; $dd7f: a9 22     
            sta $00            ; $dd81: 85 00     
__dd83:     ldx #$8a           ; $dd83: a2 8a     
            ldy #$66           ; $dd85: a0 66     
            clv                ; $dd87: b8        
            clc                ; $dd88: 18        
            lda #$00           ; $dd89: a9 00     
            lda $ff,x          ; $dd8b: b5 ff     
            bpl __dda1         ; $dd8d: 10 12     
            beq __dda1         ; $dd8f: f0 10     
            bvs __dda1         ; $dd91: 70 0e     
            bcs __dda1         ; $dd93: b0 0c     
            cmp #$bb           ; $dd95: c9 bb     
            bne __dda1         ; $dd97: d0 08     
            cpy #$66           ; $dd99: c0 66     
            bne __dda1         ; $dd9b: d0 04     
            cpx #$8a           ; $dd9d: e0 8a     
            beq __dda5         ; $dd9f: f0 04     
__dda1:     lda #$23           ; $dda1: a9 23     
            sta $00            ; $dda3: 85 00     
__dda5:     bit $01            ; $dda5: 24 01     
            sec                ; $dda7: 38        
            lda #$44           ; $dda8: a9 44     
            ldx #$00           ; $ddaa: a2 00     
            sta $33,x          ; $ddac: 95 33     
            lda $33            ; $ddae: a5 33     
            bcc __ddca         ; $ddb0: 90 18     
            cmp #$44           ; $ddb2: c9 44     
            bne __ddca         ; $ddb4: d0 14     
            bvc __ddca         ; $ddb6: 50 12     
            clc                ; $ddb8: 18        
            clv                ; $ddb9: b8        
            lda #$99           ; $ddba: a9 99     
            ldx #$80           ; $ddbc: a2 80     
            sta $85,x          ; $ddbe: 95 85     
            lda $05            ; $ddc0: a5 05     
            bcs __ddca         ; $ddc2: b0 06     
            cmp #$99           ; $ddc4: c9 99     
            bne __ddca         ; $ddc6: d0 02     
            bvc __ddce         ; $ddc8: 50 04     
__ddca:     lda #$24           ; $ddca: a9 24     
            sta $00            ; $ddcc: 85 00     
__ddce:     ldy #$25           ; $ddce: a0 25     
            ldx #$78           ; $ddd0: a2 78     
            jsr __f990         ; $ddd2: 20 90 f9  
            sta $00,x          ; $ddd5: 95 00     
            lsr $00,x          ; $ddd7: 56 00     
            lda $00,x          ; $ddd9: b5 00     
            jsr __f99d         ; $dddb: 20 9d f9  
            iny                ; $ddde: c8        
            sta $00,x          ; $dddf: 95 00     
            lsr $00,x          ; $dde1: 56 00     
            lda $00,x          ; $dde3: b5 00     
            jsr __f9ad         ; $dde5: 20 ad f9  
            iny                ; $dde8: c8        
            jsr __f9bd         ; $dde9: 20 bd f9  
            sta $00,x          ; $ddec: 95 00     
            asl $00,x          ; $ddee: 16 00     
            lda $00,x          ; $ddf0: b5 00     
            jsr __f9c3         ; $ddf2: 20 c3 f9  
            iny                ; $ddf5: c8        
            sta $00,x          ; $ddf6: 95 00     
            asl $00,x          ; $ddf8: 16 00     
            lda $00,x          ; $ddfa: b5 00     
            jsr __f9d4         ; $ddfc: 20 d4 f9  
            iny                ; $ddff: c8        
            jsr __f9e4         ; $de00: 20 e4 f9  
            sta $00,x          ; $de03: 95 00     
            ror $00,x          ; $de05: 76 00     
            lda $00,x          ; $de07: b5 00     
            jsr __f9ea         ; $de09: 20 ea f9  
            iny                ; $de0c: c8        
            sta $00,x          ; $de0d: 95 00     
            ror $00,x          ; $de0f: 76 00     
            lda $00,x          ; $de11: b5 00     
            jsr __f9fb         ; $de13: 20 fb f9  
            iny                ; $de16: c8        
            jsr __fa0a         ; $de17: 20 0a fa  
            sta $00,x          ; $de1a: 95 00     
            rol $00,x          ; $de1c: 36 00     
            lda $00,x          ; $de1e: b5 00     
            jsr __fa10         ; $de20: 20 10 fa  
            iny                ; $de23: c8        
            sta $00,x          ; $de24: 95 00     
            rol $00,x          ; $de26: 36 00     
            lda $00,x          ; $de28: b5 00     
            jsr __fa21         ; $de2a: 20 21 fa  
            lda #$ff           ; $de2d: a9 ff     
            sta $00,x          ; $de2f: 95 00     
            sta $01            ; $de31: 85 01     
            bit $01            ; $de33: 24 01     
            sec                ; $de35: 38        
            inc $00,x          ; $de36: f6 00     
            bne __de46         ; $de38: d0 0c     
            bmi __de46         ; $de3a: 30 0a     
            bvc __de46         ; $de3c: 50 08     
            bcc __de46         ; $de3e: 90 06     
            lda $00,x          ; $de40: b5 00     
            cmp #$00           ; $de42: c9 00     
            beq __de4a         ; $de44: f0 04     
__de46:     lda #$2d           ; $de46: a9 2d     
            sta $00            ; $de48: 85 00     
__de4a:     lda #$7f           ; $de4a: a9 7f     
            sta $00,x          ; $de4c: 95 00     
            clv                ; $de4e: b8        
            clc                ; $de4f: 18        
            inc $00,x          ; $de50: f6 00     
            beq __de60         ; $de52: f0 0c     
            bpl __de60         ; $de54: 10 0a     
            bvs __de60         ; $de56: 70 08     
            bcs __de60         ; $de58: b0 06     
            lda $00,x          ; $de5a: b5 00     
            cmp #$80           ; $de5c: c9 80     
            beq __de64         ; $de5e: f0 04     
__de60:     lda #$2e           ; $de60: a9 2e     
            sta $00            ; $de62: 85 00     
__de64:     lda #$00           ; $de64: a9 00     
            sta $00,x          ; $de66: 95 00     
            bit $01            ; $de68: 24 01     
            sec                ; $de6a: 38        
            dec $00,x          ; $de6b: d6 00     
            beq __de7b         ; $de6d: f0 0c     
            bpl __de7b         ; $de6f: 10 0a     
            bvc __de7b         ; $de71: 50 08     
            bcc __de7b         ; $de73: 90 06     
            lda $00,x          ; $de75: b5 00     
            cmp #$ff           ; $de77: c9 ff     
            beq __de7f         ; $de79: f0 04     
__de7b:     lda #$2f           ; $de7b: a9 2f     
            sta $00            ; $de7d: 85 00     
__de7f:     lda #$80           ; $de7f: a9 80     
            sta $00,x          ; $de81: 95 00     
            clv                ; $de83: b8        
            clc                ; $de84: 18        
            dec $00,x          ; $de85: d6 00     
            beq __de95         ; $de87: f0 0c     
            bmi __de95         ; $de89: 30 0a     
            bvs __de95         ; $de8b: 70 08     
            bcs __de95         ; $de8d: b0 06     
            lda $00,x          ; $de8f: b5 00     
            cmp #$7f           ; $de91: c9 7f     
            beq __de99         ; $de93: f0 04     
__de95:     lda #$30           ; $de95: a9 30     
            sta $00            ; $de97: 85 00     
__de99:     lda #$01           ; $de99: a9 01     
            sta $00,x          ; $de9b: 95 00     
            dec $00,x          ; $de9d: d6 00     
            beq __dea5         ; $de9f: f0 04     
            lda #$31           ; $dea1: a9 31     
            sta $00            ; $dea3: 85 00     
__dea5:     lda #$33           ; $dea5: a9 33     
            sta $78            ; $dea7: 85 78     
            lda #$44           ; $dea9: a9 44     
            ldy #$78           ; $deab: a0 78     
            ldx #$00           ; $dead: a2 00     
            sec                ; $deaf: 38        
            bit $01            ; $deb0: 24 01     
            ldx $00,y          ; $deb2: b6 00     
            bcc __dec8         ; $deb4: 90 12     
            bvc __dec8         ; $deb6: 50 10     
            bmi __dec8         ; $deb8: 30 0e     
            beq __dec8         ; $deba: f0 0c     
            cpx #$33           ; $debc: e0 33     
            bne __dec8         ; $debe: d0 08     
            cpy #$78           ; $dec0: c0 78     
            bne __dec8         ; $dec2: d0 04     
            cmp #$44           ; $dec4: c9 44     
            beq __decc         ; $dec6: f0 04     
__dec8:     lda #$32           ; $dec8: a9 32     
            sta $00            ; $deca: 85 00     
__decc:     lda #$97           ; $decc: a9 97     
            sta $7f            ; $dece: 85 7f     
            lda #$47           ; $ded0: a9 47     
            ldy #$ff           ; $ded2: a0 ff     
            ldx #$00           ; $ded4: a2 00     
            clc                ; $ded6: 18        
            clv                ; $ded7: b8        
            ldx $80,y          ; $ded8: b6 80     
            bcs __deee         ; $deda: b0 12     
            bvs __deee         ; $dedc: 70 10     
            bpl __deee         ; $dede: 10 0e     
            beq __deee         ; $dee0: f0 0c     
            cpx #$97           ; $dee2: e0 97     
            bne __deee         ; $dee4: d0 08     
            cpy #$ff           ; $dee6: c0 ff     
            bne __deee         ; $dee8: d0 04     
            cmp #$47           ; $deea: c9 47     
            beq __def2         ; $deec: f0 04     
__deee:     lda #$33           ; $deee: a9 33     
            sta $00            ; $def0: 85 00     
__def2:     lda #$00           ; $def2: a9 00     
            sta $7f            ; $def4: 85 7f     
            lda #$47           ; $def6: a9 47     
            ldy #$ff           ; $def8: a0 ff     
            ldx #$69           ; $defa: a2 69     
            clc                ; $defc: 18        
            clv                ; $defd: b8        
            stx $80,y          ; $defe: 96 80     
            bcs __df1a         ; $df00: b0 18     
            bvs __df1a         ; $df02: 70 16     
            bmi __df1a         ; $df04: 30 14     
            beq __df1a         ; $df06: f0 12     
            cpx #$69           ; $df08: e0 69     
            bne __df1a         ; $df0a: d0 0e     
            cpy #$ff           ; $df0c: c0 ff     
            bne __df1a         ; $df0e: d0 0a     
            cmp #$47           ; $df10: c9 47     
            bne __df1a         ; $df12: d0 06     
            lda $7f            ; $df14: a5 7f     
            cmp #$69           ; $df16: c9 69     
            beq __df1e         ; $df18: f0 04     
__df1a:     lda #$34           ; $df1a: a9 34     
            sta $00            ; $df1c: 85 00     
__df1e:     lda #$f5           ; $df1e: a9 f5     
            sta $4f            ; $df20: 85 4f     
            lda #$47           ; $df22: a9 47     
            ldy #$4f           ; $df24: a0 4f     
            bit $01            ; $df26: 24 01     
            ldx #$00           ; $df28: a2 00     
            sec                ; $df2a: 38        
            stx $00,y          ; $df2b: 96 00     
            bcc __df45         ; $df2d: 90 16     
            bvc __df45         ; $df2f: 50 14     
            bmi __df45         ; $df31: 30 12     
            bne __df45         ; $df33: d0 10     
            cpx #$00           ; $df35: e0 00     
            bne __df45         ; $df37: d0 0c     
            cpy #$4f           ; $df39: c0 4f     
            bne __df45         ; $df3b: d0 08     
            cmp #$47           ; $df3d: c9 47     
            bne __df45         ; $df3f: d0 04     
            lda $4f            ; $df41: a5 4f     
            beq __df49         ; $df43: f0 04     
__df45:     lda #$35           ; $df45: a9 35     
            sta $00            ; $df47: 85 00     
__df49:     rts                ; $df49: 60        

;-------------------------------------------------------------------------------
__df4a:     lda #$89           ; $df4a: a9 89     
            sta $0300          ; $df4c: 8d 00 03  
            lda #$a3           ; $df4f: a9 a3     
            sta $33            ; $df51: 85 33     
            lda #$12           ; $df53: a9 12     
            sta $0245          ; $df55: 8d 45 02  
            ldx #$65           ; $df58: a2 65     
            ldy #$00           ; $df5a: a0 00     
            sec                ; $df5c: 38        
            lda #$00           ; $df5d: a9 00     
            clv                ; $df5f: b8        
            lda $0300,y        ; $df60: b9 00 03  
            beq __df71         ; $df63: f0 0c     
            bcc __df71         ; $df65: 90 0a     
            bvs __df71         ; $df67: 70 08     
            cmp #$89           ; $df69: c9 89     
            bne __df71         ; $df6b: d0 04     
            cpx #$65           ; $df6d: e0 65     
            beq __df75         ; $df6f: f0 04     
__df71:     lda #$36           ; $df71: a9 36     
            sta $00            ; $df73: 85 00     
__df75:     lda #$ff           ; $df75: a9 ff     
            sta $01            ; $df77: 85 01     
            bit $01            ; $df79: 24 01     
            ldy #$34           ; $df7b: a0 34     
            lda $ffff,y        ; $df7d: b9 ff ff  
            cmp #$a3           ; $df80: c9 a3     
            bne __df86         ; $df82: d0 02     
            bcs __df8a         ; $df84: b0 04     
__df86:     lda #$37           ; $df86: a9 37     
            sta $00            ; $df88: 85 00     
__df8a:     lda #$46           ; $df8a: a9 46     
            sta $ff            ; $df8c: 85 ff     
            ldy #$ff           ; $df8e: a0 ff     
            lda $0146,y        ; $df90: b9 46 01  
            cmp #$12           ; $df93: c9 12     
            beq __df9b         ; $df95: f0 04     
            lda #$38           ; $df97: a9 38     
            sta $00            ; $df99: 85 00     
__df9b:     ldx #$39           ; $df9b: a2 39     
            clc                ; $df9d: 18        
            lda #$ff           ; $df9e: a9 ff     
            sta $01            ; $dfa0: 85 01     
            bit $01            ; $dfa2: 24 01     
            lda #$aa           ; $dfa4: a9 aa     
            sta $0400          ; $dfa6: 8d 00 04  
            lda #$55           ; $dfa9: a9 55     
            ldy #$00           ; $dfab: a0 00     
            ora $0400,y        ; $dfad: 19 00 04  
            bcs __dfba         ; $dfb0: b0 08     
            bpl __dfba         ; $dfb2: 10 06     
            cmp #$ff           ; $dfb4: c9 ff     
            bne __dfba         ; $dfb6: d0 02     
            bvs __dfbc         ; $dfb8: 70 02     
__dfba:     stx $00            ; $dfba: 86 00     
__dfbc:     inx                ; $dfbc: e8        
            sec                ; $dfbd: 38        
            clv                ; $dfbe: b8        
            lda #$00           ; $dfbf: a9 00     
            ora $0400,y        ; $dfc1: 19 00 04  
            beq __dfcc         ; $dfc4: f0 06     
            bvs __dfcc         ; $dfc6: 70 04     
            bcc __dfcc         ; $dfc8: 90 02     
            bmi __dfce         ; $dfca: 30 02     
__dfcc:     stx $00            ; $dfcc: 86 00     
__dfce:     inx                ; $dfce: e8        
            clc                ; $dfcf: 18        
            bit $01            ; $dfd0: 24 01     
            lda #$55           ; $dfd2: a9 55     
            and $0400,y        ; $dfd4: 39 00 04  
            bne __dfdf         ; $dfd7: d0 06     
            bvc __dfdf         ; $dfd9: 50 04     
            bcs __dfdf         ; $dfdb: b0 02     
            bpl __dfe1         ; $dfdd: 10 02     
__dfdf:     stx $00            ; $dfdf: 86 00     
__dfe1:     inx                ; $dfe1: e8        
            sec                ; $dfe2: 38        
            clv                ; $dfe3: b8        
            lda #$ef           ; $dfe4: a9 ef     
            sta $0400          ; $dfe6: 8d 00 04  
            lda #$f8           ; $dfe9: a9 f8     
            and $0400,y        ; $dfeb: 39 00 04  
            bcc __dff8         ; $dfee: 90 08     
            bpl __dff8         ; $dff0: 10 06     
            cmp #$e8           ; $dff2: c9 e8     
            bne __dff8         ; $dff4: d0 02     
            bvc __dffa         ; $dff6: 50 02     
__dff8:     stx $00            ; $dff8: 86 00     
__dffa:     inx                ; $dffa: e8        
            clc                ; $dffb: 18        
            bit $01            ; $dffc: 24 01     
            lda #$aa           ; $dffe: a9 aa     
            sta $0400          ; $e000: 8d 00 04  
            lda #$5f           ; $e003: a9 5f     
            eor $0400,y        ; $e005: 59 00 04  
            bcs __e012         ; $e008: b0 08     
            bpl __e012         ; $e00a: 10 06     
            cmp #$f5           ; $e00c: c9 f5     
            bne __e012         ; $e00e: d0 02     
            bvs __e014         ; $e010: 70 02     
__e012:     stx $00            ; $e012: 86 00     
__e014:     inx                ; $e014: e8        
            sec                ; $e015: 38        
            clv                ; $e016: b8        
            lda #$70           ; $e017: a9 70     
            sta $0400          ; $e019: 8d 00 04  
            eor $0400,y        ; $e01c: 59 00 04  
            bne __e027         ; $e01f: d0 06     
            bvs __e027         ; $e021: 70 04     
            bcc __e027         ; $e023: 90 02     
            bpl __e029         ; $e025: 10 02     
__e027:     stx $00            ; $e027: 86 00     
__e029:     inx                ; $e029: e8        
            clc                ; $e02a: 18        
            bit $01            ; $e02b: 24 01     
            lda #$69           ; $e02d: a9 69     
            sta $0400          ; $e02f: 8d 00 04  
            lda #$00           ; $e032: a9 00     
            adc $0400,y        ; $e034: 79 00 04  
            bmi __e041         ; $e037: 30 08     
            bcs __e041         ; $e039: b0 06     
            cmp #$69           ; $e03b: c9 69     
            bne __e041         ; $e03d: d0 02     
            bvc __e043         ; $e03f: 50 02     
__e041:     stx $00            ; $e041: 86 00     
__e043:     inx                ; $e043: e8        
            sec                ; $e044: 38        
            bit $01            ; $e045: 24 01     
            lda #$00           ; $e047: a9 00     
            adc $0400,y        ; $e049: 79 00 04  
            bmi __e056         ; $e04c: 30 08     
            bcs __e056         ; $e04e: b0 06     
            cmp #$6a           ; $e050: c9 6a     
            bne __e056         ; $e052: d0 02     
            bvc __e058         ; $e054: 50 02     
__e056:     stx $00            ; $e056: 86 00     
__e058:     inx                ; $e058: e8        
            sec                ; $e059: 38        
            clv                ; $e05a: b8        
            lda #$7f           ; $e05b: a9 7f     
            sta $0400          ; $e05d: 8d 00 04  
            adc $0400,y        ; $e060: 79 00 04  
            bpl __e06d         ; $e063: 10 08     
            bcs __e06d         ; $e065: b0 06     
            cmp #$ff           ; $e067: c9 ff     
            bne __e06d         ; $e069: d0 02     
            bvs __e06f         ; $e06b: 70 02     
__e06d:     stx $00            ; $e06d: 86 00     
__e06f:     inx                ; $e06f: e8        
            clc                ; $e070: 18        
            bit $01            ; $e071: 24 01     
            lda #$80           ; $e073: a9 80     
            sta $0400          ; $e075: 8d 00 04  
            lda #$7f           ; $e078: a9 7f     
            adc $0400,y        ; $e07a: 79 00 04  
            bpl __e087         ; $e07d: 10 08     
            bcs __e087         ; $e07f: b0 06     
            cmp #$ff           ; $e081: c9 ff     
            bne __e087         ; $e083: d0 02     
            bvc __e089         ; $e085: 50 02     
__e087:     stx $00            ; $e087: 86 00     
__e089:     inx                ; $e089: e8        
            sec                ; $e08a: 38        
            clv                ; $e08b: b8        
            lda #$80           ; $e08c: a9 80     
            sta $0400          ; $e08e: 8d 00 04  
            lda #$7f           ; $e091: a9 7f     
            adc $0400,y        ; $e093: 79 00 04  
            bne __e09e         ; $e096: d0 06     
            bmi __e09e         ; $e098: 30 04     
            bvs __e09e         ; $e09a: 70 02     
            bcs __e0a0         ; $e09c: b0 02     
__e09e:     stx $00            ; $e09e: 86 00     
__e0a0:     inx                ; $e0a0: e8        
            bit $01            ; $e0a1: 24 01     
            lda #$40           ; $e0a3: a9 40     
            sta $0400          ; $e0a5: 8d 00 04  
            cmp $0400,y        ; $e0a8: d9 00 04  
            bmi __e0b3         ; $e0ab: 30 06     
            bcc __e0b3         ; $e0ad: 90 04     
            bne __e0b3         ; $e0af: d0 02     
            bvs __e0b5         ; $e0b1: 70 02     
__e0b3:     stx $00            ; $e0b3: 86 00     
__e0b5:     inx                ; $e0b5: e8        
            clv                ; $e0b6: b8        
            dec $0400          ; $e0b7: ce 00 04  
            cmp $0400,y        ; $e0ba: d9 00 04  
            beq __e0c5         ; $e0bd: f0 06     
            bmi __e0c5         ; $e0bf: 30 04     
            bcc __e0c5         ; $e0c1: 90 02     
            bvc __e0c7         ; $e0c3: 50 02     
__e0c5:     stx $00            ; $e0c5: 86 00     
__e0c7:     inx                ; $e0c7: e8        
            inc $0400          ; $e0c8: ee 00 04  
            inc $0400          ; $e0cb: ee 00 04  
            cmp $0400,y        ; $e0ce: d9 00 04  
            beq __e0d5         ; $e0d1: f0 02     
            bmi __e0d7         ; $e0d3: 30 02     
__e0d5:     stx $00            ; $e0d5: 86 00     
__e0d7:     inx                ; $e0d7: e8        
            lda #$00           ; $e0d8: a9 00     
            sta $0400          ; $e0da: 8d 00 04  
            lda #$80           ; $e0dd: a9 80     
            cmp $0400,y        ; $e0df: d9 00 04  
            beq __e0e8         ; $e0e2: f0 04     
            bpl __e0e8         ; $e0e4: 10 02     
            bcs __e0ea         ; $e0e6: b0 02     
__e0e8:     stx $00            ; $e0e8: 86 00     
__e0ea:     inx                ; $e0ea: e8        
            ldy #$80           ; $e0eb: a0 80     
            sty $0400          ; $e0ed: 8c 00 04  
            ldy #$00           ; $e0f0: a0 00     
            cmp $0400,y        ; $e0f2: d9 00 04  
            bne __e0fb         ; $e0f5: d0 04     
            bmi __e0fb         ; $e0f7: 30 02     
            bcs __e0fd         ; $e0f9: b0 02     
__e0fb:     stx $00            ; $e0fb: 86 00     
__e0fd:     inx                ; $e0fd: e8        
            inc $0400          ; $e0fe: ee 00 04  
            cmp $0400,y        ; $e101: d9 00 04  
            bcs __e10a         ; $e104: b0 04     
            beq __e10a         ; $e106: f0 02     
            bmi __e10c         ; $e108: 30 02     
__e10a:     stx $00            ; $e10a: 86 00     
__e10c:     inx                ; $e10c: e8        
            dec $0400          ; $e10d: ce 00 04  
            dec $0400          ; $e110: ce 00 04  
            cmp $0400,y        ; $e113: d9 00 04  
            bcc __e11c         ; $e116: 90 04     
            beq __e11c         ; $e118: f0 02     
            bpl __e11e         ; $e11a: 10 02     
__e11c:     stx $00            ; $e11c: 86 00     
__e11e:     inx                ; $e11e: e8        
            bit $01            ; $e11f: 24 01     
            lda #$40           ; $e121: a9 40     
            sta $0400          ; $e123: 8d 00 04  
            sec                ; $e126: 38        
            sbc $0400,y        ; $e127: f9 00 04  
            bmi __e136         ; $e12a: 30 0a     
            bcc __e136         ; $e12c: 90 08     
            bne __e136         ; $e12e: d0 06     
            bvs __e136         ; $e130: 70 04     
            cmp #$00           ; $e132: c9 00     
            beq __e138         ; $e134: f0 02     
__e136:     stx $00            ; $e136: 86 00     
__e138:     inx                ; $e138: e8        
            clv                ; $e139: b8        
            sec                ; $e13a: 38        
            lda #$40           ; $e13b: a9 40     
            dec $0400          ; $e13d: ce 00 04  
            sbc $0400,y        ; $e140: f9 00 04  
            beq __e14f         ; $e143: f0 0a     
            bmi __e14f         ; $e145: 30 08     
            bcc __e14f         ; $e147: 90 06     
            bvs __e14f         ; $e149: 70 04     
            cmp #$01           ; $e14b: c9 01     
            beq __e151         ; $e14d: f0 02     
__e14f:     stx $00            ; $e14f: 86 00     
__e151:     inx                ; $e151: e8        
            lda #$40           ; $e152: a9 40     
            sec                ; $e154: 38        
            bit $01            ; $e155: 24 01     
            inc $0400          ; $e157: ee 00 04  
            inc $0400          ; $e15a: ee 00 04  
            sbc $0400,y        ; $e15d: f9 00 04  
            bcs __e16c         ; $e160: b0 0a     
            beq __e16c         ; $e162: f0 08     
            bpl __e16c         ; $e164: 10 06     
            bvs __e16c         ; $e166: 70 04     
            cmp #$ff           ; $e168: c9 ff     
            beq __e16e         ; $e16a: f0 02     
__e16c:     stx $00            ; $e16c: 86 00     
__e16e:     inx                ; $e16e: e8        
            clc                ; $e16f: 18        
            lda #$00           ; $e170: a9 00     
            sta $0400          ; $e172: 8d 00 04  
            lda #$80           ; $e175: a9 80     
            sbc $0400,y        ; $e177: f9 00 04  
            bcc __e180         ; $e17a: 90 04     
            cmp #$7f           ; $e17c: c9 7f     
            beq __e182         ; $e17e: f0 02     
__e180:     stx $00            ; $e180: 86 00     
__e182:     inx                ; $e182: e8        
            sec                ; $e183: 38        
            lda #$7f           ; $e184: a9 7f     
            sta $0400          ; $e186: 8d 00 04  
            lda #$81           ; $e189: a9 81     
            sbc $0400,y        ; $e18b: f9 00 04  
            bvc __e196         ; $e18e: 50 06     
            bcc __e196         ; $e190: 90 04     
            cmp #$02           ; $e192: c9 02     
            beq __e198         ; $e194: f0 02     
__e196:     stx $00            ; $e196: 86 00     
__e198:     inx                ; $e198: e8        
            lda #$00           ; $e199: a9 00     
            lda #$87           ; $e19b: a9 87     
            sta $0400,y        ; $e19d: 99 00 04  
            lda $0400          ; $e1a0: ad 00 04  
            cmp #$87           ; $e1a3: c9 87     
            beq __e1a9         ; $e1a5: f0 02     
            stx $00            ; $e1a7: 86 00     
__e1a9:     rts                ; $e1a9: 60        

;-------------------------------------------------------------------------------
__e1aa:     lda #$ff           ; $e1aa: a9 ff     
            sta $01            ; $e1ac: 85 01     
            lda #$aa           ; $e1ae: a9 aa     
            sta $0633          ; $e1b0: 8d 33 06  
            lda #$bb           ; $e1b3: a9 bb     
            sta $0689          ; $e1b5: 8d 89 06  
            ldx #$00           ; $e1b8: a2 00     
            lda #$66           ; $e1ba: a9 66     
            bit $01            ; $e1bc: 24 01     
            sec                ; $e1be: 38        
            ldy #$00           ; $e1bf: a0 00     
            ldy $0633,x        ; $e1c1: bc 33 06  
            bpl __e1d8         ; $e1c4: 10 12     
            beq __e1d8         ; $e1c6: f0 10     
            bvc __e1d8         ; $e1c8: 50 0e     
            bcc __e1d8         ; $e1ca: 90 0c     
            cmp #$66           ; $e1cc: c9 66     
            bne __e1d8         ; $e1ce: d0 08     
            cpx #$00           ; $e1d0: e0 00     
            bne __e1d8         ; $e1d2: d0 04     
            cpy #$aa           ; $e1d4: c0 aa     
            beq __e1dc         ; $e1d6: f0 04     
__e1d8:     lda #$51           ; $e1d8: a9 51     
            sta $00            ; $e1da: 85 00     
__e1dc:     ldx #$8a           ; $e1dc: a2 8a     
            lda #$66           ; $e1de: a9 66     
            clv                ; $e1e0: b8        
            clc                ; $e1e1: 18        
            ldy #$00           ; $e1e2: a0 00     
            ldy $05ff,x        ; $e1e4: bc ff 05  
            bpl __e1fb         ; $e1e7: 10 12     
            beq __e1fb         ; $e1e9: f0 10     
            bvs __e1fb         ; $e1eb: 70 0e     
            bcs __e1fb         ; $e1ed: b0 0c     
            cpy #$bb           ; $e1ef: c0 bb     
            bne __e1fb         ; $e1f1: d0 08     
            cmp #$66           ; $e1f3: c9 66     
            bne __e1fb         ; $e1f5: d0 04     
            cpx #$8a           ; $e1f7: e0 8a     
            beq __e1ff         ; $e1f9: f0 04     
__e1fb:     lda #$52           ; $e1fb: a9 52     
            sta $00            ; $e1fd: 85 00     
__e1ff:     ldy #$53           ; $e1ff: a0 53     
            lda #$aa           ; $e201: a9 aa     
            ldx #$78           ; $e203: a2 78     
            sta $0678          ; $e205: 8d 78 06  
            jsr __f7b6         ; $e208: 20 b6 f7  
            ora $0600,x        ; $e20b: 1d 00 06  
            jsr __f7c0         ; $e20e: 20 c0 f7  
            iny                ; $e211: c8        
            lda #$00           ; $e212: a9 00     
            sta $0678          ; $e214: 8d 78 06  
            jsr __f7ce         ; $e217: 20 ce f7  
            ora $0600,x        ; $e21a: 1d 00 06  
            jsr __f7d3         ; $e21d: 20 d3 f7  
            iny                ; $e220: c8        
            lda #$aa           ; $e221: a9 aa     
            sta $0678          ; $e223: 8d 78 06  
            jsr __f7df         ; $e226: 20 df f7  
            and $0600,x        ; $e229: 3d 00 06  
            jsr __f7e5         ; $e22c: 20 e5 f7  
            iny                ; $e22f: c8        
            lda #$ef           ; $e230: a9 ef     
            sta $0678          ; $e232: 8d 78 06  
            jsr __f7f1         ; $e235: 20 f1 f7  
            and $0600,x        ; $e238: 3d 00 06  
            jsr __f7f6         ; $e23b: 20 f6 f7  
            iny                ; $e23e: c8        
            lda #$aa           ; $e23f: a9 aa     
            sta $0678          ; $e241: 8d 78 06  
            jsr __f804         ; $e244: 20 04 f8  
            eor $0600,x        ; $e247: 5d 00 06  
            jsr __f80a         ; $e24a: 20 0a f8  
            iny                ; $e24d: c8        
            lda #$70           ; $e24e: a9 70     
            sta $0678          ; $e250: 8d 78 06  
            jsr __f818         ; $e253: 20 18 f8  
            eor $0600,x        ; $e256: 5d 00 06  
            jsr __f81d         ; $e259: 20 1d f8  
            iny                ; $e25c: c8        
            lda #$69           ; $e25d: a9 69     
            sta $0678          ; $e25f: 8d 78 06  
            jsr __f829         ; $e262: 20 29 f8  
            adc $0600,x        ; $e265: 7d 00 06  
            jsr __f82f         ; $e268: 20 2f f8  
            iny                ; $e26b: c8        
            jsr __f83d         ; $e26c: 20 3d f8  
            adc $0600,x        ; $e26f: 7d 00 06  
            jsr __f843         ; $e272: 20 43 f8  
            iny                ; $e275: c8        
            lda #$7f           ; $e276: a9 7f     
            sta $0678          ; $e278: 8d 78 06  
            jsr __f851         ; $e27b: 20 51 f8  
            adc $0600,x        ; $e27e: 7d 00 06  
            jsr __f856         ; $e281: 20 56 f8  
            iny                ; $e284: c8        
            lda #$80           ; $e285: a9 80     
            sta $0678          ; $e287: 8d 78 06  
            jsr __f864         ; $e28a: 20 64 f8  
            adc $0600,x        ; $e28d: 7d 00 06  
            jsr __f86a         ; $e290: 20 6a f8  
            iny                ; $e293: c8        
            jsr __f878         ; $e294: 20 78 f8  
            adc $0600,x        ; $e297: 7d 00 06  
            jsr __f87d         ; $e29a: 20 7d f8  
            iny                ; $e29d: c8        
            lda #$40           ; $e29e: a9 40     
            sta $0678          ; $e2a0: 8d 78 06  
            jsr __f889         ; $e2a3: 20 89 f8  
            cmp $0600,x        ; $e2a6: dd 00 06  
            jsr __f88e         ; $e2a9: 20 8e f8  
            iny                ; $e2ac: c8        
            pha                ; $e2ad: 48        
            lda #$3f           ; $e2ae: a9 3f     
            sta $0678          ; $e2b0: 8d 78 06  
            pla                ; $e2b3: 68        
            jsr __f89a         ; $e2b4: 20 9a f8  
            cmp $0600,x        ; $e2b7: dd 00 06  
            jsr __f89c         ; $e2ba: 20 9c f8  
            iny                ; $e2bd: c8        
            pha                ; $e2be: 48        
            lda #$41           ; $e2bf: a9 41     
            sta $0678          ; $e2c1: 8d 78 06  
            pla                ; $e2c4: 68        
            cmp $0600,x        ; $e2c5: dd 00 06  
            jsr __f8a8         ; $e2c8: 20 a8 f8  
            iny                ; $e2cb: c8        
            pha                ; $e2cc: 48        
            lda #$00           ; $e2cd: a9 00     
            sta $0678          ; $e2cf: 8d 78 06  
            pla                ; $e2d2: 68        
            jsr __f8b2         ; $e2d3: 20 b2 f8  
            cmp $0600,x        ; $e2d6: dd 00 06  
            jsr __f8b5         ; $e2d9: 20 b5 f8  
            iny                ; $e2dc: c8        
            pha                ; $e2dd: 48        
            lda #$80           ; $e2de: a9 80     
            sta $0678          ; $e2e0: 8d 78 06  
            pla                ; $e2e3: 68        
            cmp $0600,x        ; $e2e4: dd 00 06  
            jsr __f8bf         ; $e2e7: 20 bf f8  
            iny                ; $e2ea: c8        
            pha                ; $e2eb: 48        
            lda #$81           ; $e2ec: a9 81     
            sta $0678          ; $e2ee: 8d 78 06  
            pla                ; $e2f1: 68        
            cmp $0600,x        ; $e2f2: dd 00 06  
            jsr __f8c9         ; $e2f5: 20 c9 f8  
            iny                ; $e2f8: c8        
            pha                ; $e2f9: 48        
            lda #$7f           ; $e2fa: a9 7f     
            sta $0678          ; $e2fc: 8d 78 06  
            pla                ; $e2ff: 68        
            cmp $0600,x        ; $e300: dd 00 06  
            jsr __f8d3         ; $e303: 20 d3 f8  
            iny                ; $e306: c8        
            lda #$40           ; $e307: a9 40     
            sta $0678          ; $e309: 8d 78 06  
            jsr __f931         ; $e30c: 20 31 f9  
            sbc $0600,x        ; $e30f: fd 00 06  
            jsr __f937         ; $e312: 20 37 f9  
            iny                ; $e315: c8        
            lda #$3f           ; $e316: a9 3f     
            sta $0678          ; $e318: 8d 78 06  
            jsr __f947         ; $e31b: 20 47 f9  
            sbc $0600,x        ; $e31e: fd 00 06  
            jsr __f94c         ; $e321: 20 4c f9  
            iny                ; $e324: c8        
            lda #$41           ; $e325: a9 41     
            sta $0678          ; $e327: 8d 78 06  
            jsr __f95c         ; $e32a: 20 5c f9  
            sbc $0600,x        ; $e32d: fd 00 06  
            jsr __f962         ; $e330: 20 62 f9  
            iny                ; $e333: c8        
            lda #$00           ; $e334: a9 00     
            sta $0678          ; $e336: 8d 78 06  
            jsr __f972         ; $e339: 20 72 f9  
            sbc $0600,x        ; $e33c: fd 00 06  
            jsr __f976         ; $e33f: 20 76 f9  
            iny                ; $e342: c8        
            lda #$7f           ; $e343: a9 7f     
            sta $0678          ; $e345: 8d 78 06  
            jsr __f980         ; $e348: 20 80 f9  
            sbc $0600,x        ; $e34b: fd 00 06  
            jsr __f984         ; $e34e: 20 84 f9  
            lda #$aa           ; $e351: a9 aa     
            sta $0633          ; $e353: 8d 33 06  
            lda #$bb           ; $e356: a9 bb     
            sta $0689          ; $e358: 8d 89 06  
            ldx #$00           ; $e35b: a2 00     
            ldy #$66           ; $e35d: a0 66     
            bit $01            ; $e35f: 24 01     
            sec                ; $e361: 38        
            lda #$00           ; $e362: a9 00     
            lda $0633,x        ; $e364: bd 33 06  
            bpl __e37b         ; $e367: 10 12     
            beq __e37b         ; $e369: f0 10     
            bvc __e37b         ; $e36b: 50 0e     
            bcc __e37b         ; $e36d: 90 0c     
            cpy #$66           ; $e36f: c0 66     
            bne __e37b         ; $e371: d0 08     
            cpx #$00           ; $e373: e0 00     
            bne __e37b         ; $e375: d0 04     
            cmp #$aa           ; $e377: c9 aa     
            beq __e37f         ; $e379: f0 04     
__e37b:     lda #$6a           ; $e37b: a9 6a     
            sta $00            ; $e37d: 85 00     
__e37f:     ldx #$8a           ; $e37f: a2 8a     
            ldy #$66           ; $e381: a0 66     
            clv                ; $e383: b8        
            clc                ; $e384: 18        
            lda #$00           ; $e385: a9 00     
            lda $05ff,x        ; $e387: bd ff 05  
            bpl __e39e         ; $e38a: 10 12     
            beq __e39e         ; $e38c: f0 10     
            bvs __e39e         ; $e38e: 70 0e     
            bcs __e39e         ; $e390: b0 0c     
            cmp #$bb           ; $e392: c9 bb     
            bne __e39e         ; $e394: d0 08     
            cpy #$66           ; $e396: c0 66     
            bne __e39e         ; $e398: d0 04     
            cpx #$8a           ; $e39a: e0 8a     
            beq __e3a2         ; $e39c: f0 04     
__e39e:     lda #$6b           ; $e39e: a9 6b     
            sta $00            ; $e3a0: 85 00     
__e3a2:     bit $01            ; $e3a2: 24 01     
            sec                ; $e3a4: 38        
            lda #$44           ; $e3a5: a9 44     
            ldx #$00           ; $e3a7: a2 00     
            sta $0633,x        ; $e3a9: 9d 33 06  
            lda $0633          ; $e3ac: ad 33 06  
            bcc __e3cb         ; $e3af: 90 1a     
            cmp #$44           ; $e3b1: c9 44     
            bne __e3cb         ; $e3b3: d0 16     
            bvc __e3cb         ; $e3b5: 50 14     
            clc                ; $e3b7: 18        
            clv                ; $e3b8: b8        
            lda #$99           ; $e3b9: a9 99     
            ldx #$80           ; $e3bb: a2 80     
            sta $0585,x        ; $e3bd: 9d 85 05  
            lda $0605          ; $e3c0: ad 05 06  
            bcs __e3cb         ; $e3c3: b0 06     
            cmp #$99           ; $e3c5: c9 99     
            bne __e3cb         ; $e3c7: d0 02     
            bvc __e3cf         ; $e3c9: 50 04     
__e3cb:     lda #$6c           ; $e3cb: a9 6c     
            sta $00            ; $e3cd: 85 00     
__e3cf:     ldy #$6d           ; $e3cf: a0 6d     
            ldx #$6d           ; $e3d1: a2 6d     
            jsr __f990         ; $e3d3: 20 90 f9  
            sta $0600,x        ; $e3d6: 9d 00 06  
            lsr $0600,x        ; $e3d9: 5e 00 06  
            lda $0600,x        ; $e3dc: bd 00 06  
            jsr __f99d         ; $e3df: 20 9d f9  
            iny                ; $e3e2: c8        
            sta $0600,x        ; $e3e3: 9d 00 06  
            lsr $0600,x        ; $e3e6: 5e 00 06  
            lda $0600,x        ; $e3e9: bd 00 06  
            jsr __f9ad         ; $e3ec: 20 ad f9  
            iny                ; $e3ef: c8        
            jsr __f9bd         ; $e3f0: 20 bd f9  
            sta $0600,x        ; $e3f3: 9d 00 06  
            asl $0600,x        ; $e3f6: 1e 00 06  
            lda $0600,x        ; $e3f9: bd 00 06  
            jsr __f9c3         ; $e3fc: 20 c3 f9  
            iny                ; $e3ff: c8        
            sta $0600,x        ; $e400: 9d 00 06  
            asl $0600,x        ; $e403: 1e 00 06  
            lda $0600,x        ; $e406: bd 00 06  
            jsr __f9d4         ; $e409: 20 d4 f9  
            iny                ; $e40c: c8        
            jsr __f9e4         ; $e40d: 20 e4 f9  
            sta $0600,x        ; $e410: 9d 00 06  
            ror $0600,x        ; $e413: 7e 00 06  
            lda $0600,x        ; $e416: bd 00 06  
            jsr __f9ea         ; $e419: 20 ea f9  
            iny                ; $e41c: c8        
            sta $0600,x        ; $e41d: 9d 00 06  
            ror $0600,x        ; $e420: 7e 00 06  
            lda $0600,x        ; $e423: bd 00 06  
            jsr __f9fb         ; $e426: 20 fb f9  
            iny                ; $e429: c8        
            jsr __fa0a         ; $e42a: 20 0a fa  
            sta $0600,x        ; $e42d: 9d 00 06  
            rol $0600,x        ; $e430: 3e 00 06  
            lda $0600,x        ; $e433: bd 00 06  
            jsr __fa10         ; $e436: 20 10 fa  
            iny                ; $e439: c8        
            sta $0600,x        ; $e43a: 9d 00 06  
            rol $0600,x        ; $e43d: 3e 00 06  
            lda $0600,x        ; $e440: bd 00 06  
            jsr __fa21         ; $e443: 20 21 fa  
            lda #$ff           ; $e446: a9 ff     
            sta $0600,x        ; $e448: 9d 00 06  
            sta $01            ; $e44b: 85 01     
            bit $01            ; $e44d: 24 01     
            sec                ; $e44f: 38        
            inc $0600,x        ; $e450: fe 00 06  
            bne __e462         ; $e453: d0 0d     
            bmi __e462         ; $e455: 30 0b     
            bvc __e462         ; $e457: 50 09     
            bcc __e462         ; $e459: 90 07     
            lda $0600,x        ; $e45b: bd 00 06  
            cmp #$00           ; $e45e: c9 00     
            beq __e466         ; $e460: f0 04     
__e462:     lda #$75           ; $e462: a9 75     
            sta $00            ; $e464: 85 00     
__e466:     lda #$7f           ; $e466: a9 7f     
            sta $0600,x        ; $e468: 9d 00 06  
            clv                ; $e46b: b8        
            clc                ; $e46c: 18        
            inc $0600,x        ; $e46d: fe 00 06  
            beq __e47f         ; $e470: f0 0d     
            bpl __e47f         ; $e472: 10 0b     
            bvs __e47f         ; $e474: 70 09     
            bcs __e47f         ; $e476: b0 07     
            lda $0600,x        ; $e478: bd 00 06  
            cmp #$80           ; $e47b: c9 80     
            beq __e483         ; $e47d: f0 04     
__e47f:     lda #$76           ; $e47f: a9 76     
            sta $00            ; $e481: 85 00     
__e483:     lda #$00           ; $e483: a9 00     
            sta $0600,x        ; $e485: 9d 00 06  
            bit $01            ; $e488: 24 01     
            sec                ; $e48a: 38        
            dec $0600,x        ; $e48b: de 00 06  
            beq __e49d         ; $e48e: f0 0d     
            bpl __e49d         ; $e490: 10 0b     
            bvc __e49d         ; $e492: 50 09     
            bcc __e49d         ; $e494: 90 07     
            lda $0600,x        ; $e496: bd 00 06  
            cmp #$ff           ; $e499: c9 ff     
            beq __e4a1         ; $e49b: f0 04     
__e49d:     lda #$77           ; $e49d: a9 77     
            sta $00            ; $e49f: 85 00     
__e4a1:     lda #$80           ; $e4a1: a9 80     
            sta $0600,x        ; $e4a3: 9d 00 06  
            clv                ; $e4a6: b8        
            clc                ; $e4a7: 18        
            dec $0600,x        ; $e4a8: de 00 06  
            beq __e4ba         ; $e4ab: f0 0d     
            bmi __e4ba         ; $e4ad: 30 0b     
            bvs __e4ba         ; $e4af: 70 09     
            bcs __e4ba         ; $e4b1: b0 07     
            lda $0600,x        ; $e4b3: bd 00 06  
            cmp #$7f           ; $e4b6: c9 7f     
            beq __e4be         ; $e4b8: f0 04     
__e4ba:     lda #$78           ; $e4ba: a9 78     
            sta $00            ; $e4bc: 85 00     
__e4be:     lda #$01           ; $e4be: a9 01     
            sta $0600,x        ; $e4c0: 9d 00 06  
            dec $0600,x        ; $e4c3: de 00 06  
            beq __e4cc         ; $e4c6: f0 04     
            lda #$79           ; $e4c8: a9 79     
            sta $00            ; $e4ca: 85 00     
__e4cc:     lda #$33           ; $e4cc: a9 33     
            sta $0678          ; $e4ce: 8d 78 06  
            lda #$44           ; $e4d1: a9 44     
            ldy #$78           ; $e4d3: a0 78     
            ldx #$00           ; $e4d5: a2 00     
            sec                ; $e4d7: 38        
            bit $01            ; $e4d8: 24 01     
            ldx $0600,y        ; $e4da: be 00 06  
            bcc __e4f1         ; $e4dd: 90 12     
            bvc __e4f1         ; $e4df: 50 10     
            bmi __e4f1         ; $e4e1: 30 0e     
            beq __e4f1         ; $e4e3: f0 0c     
            cpx #$33           ; $e4e5: e0 33     
            bne __e4f1         ; $e4e7: d0 08     
            cpy #$78           ; $e4e9: c0 78     
            bne __e4f1         ; $e4eb: d0 04     
            cmp #$44           ; $e4ed: c9 44     
            beq __e4f5         ; $e4ef: f0 04     
__e4f1:     lda #$7a           ; $e4f1: a9 7a     
            sta $00            ; $e4f3: 85 00     
__e4f5:     lda #$97           ; $e4f5: a9 97     
            sta $067f          ; $e4f7: 8d 7f 06  
            lda #$47           ; $e4fa: a9 47     
            ldy #$ff           ; $e4fc: a0 ff     
            ldx #$00           ; $e4fe: a2 00     
            clc                ; $e500: 18        
            clv                ; $e501: b8        
            ldx $0580,y        ; $e502: be 80 05  
            bcs __e519         ; $e505: b0 12     
            bvs __e519         ; $e507: 70 10     
            bpl __e519         ; $e509: 10 0e     
            beq __e519         ; $e50b: f0 0c     
            cpx #$97           ; $e50d: e0 97     
            bne __e519         ; $e50f: d0 08     
            cpy #$ff           ; $e511: c0 ff     
            bne __e519         ; $e513: d0 04     
            cmp #$47           ; $e515: c9 47     
            beq __e51d         ; $e517: f0 04     
__e519:     lda #$7b           ; $e519: a9 7b     
            sta $00            ; $e51b: 85 00     
__e51d:     rts                ; $e51d: 60        

;-------------------------------------------------------------------------------
__e51e:     lda #$55           ; $e51e: a9 55     
            sta $0580          ; $e520: 8d 80 05  
            lda #$aa           ; $e523: a9 aa     
            sta $0432          ; $e525: 8d 32 04  
            lda #$80           ; $e528: a9 80     
            sta $43            ; $e52a: 85 43     
            lda #$05           ; $e52c: a9 05     
            sta $44            ; $e52e: 85 44     
            lda #$32           ; $e530: a9 32     
            sta $45            ; $e532: 85 45     
            lda #$04           ; $e534: a9 04     
            sta $46            ; $e536: 85 46     
            ldx #$03           ; $e538: a2 03     
            ldy #$77           ; $e53a: a0 77     
            lda #$ff           ; $e53c: a9 ff     
            sta $01            ; $e53e: 85 01     
            bit $01            ; $e540: 24 01     
            sec                ; $e542: 38        
            lda #$00           ; $e543: a9 00     
            .hex a3 40         ; $e545: a3 40     Invalid Opcode - LAX ($40,x)
            nop                ; $e547: ea        
            nop                ; $e548: ea        
            nop                ; $e549: ea        
            nop                ; $e54a: ea        
            beq __e55f         ; $e54b: f0 12     
            bmi __e55f         ; $e54d: 30 10     
            bvc __e55f         ; $e54f: 50 0e     
            bcc __e55f         ; $e551: 90 0c     
            cmp #$55           ; $e553: c9 55     
            bne __e55f         ; $e555: d0 08     
            cpx #$55           ; $e557: e0 55     
            bne __e55f         ; $e559: d0 04     
            cpy #$77           ; $e55b: c0 77     
            beq __e563         ; $e55d: f0 04     
__e55f:     lda #$7c           ; $e55f: a9 7c     
            sta $00            ; $e561: 85 00     
__e563:     ldx #$05           ; $e563: a2 05     
            ldy #$33           ; $e565: a0 33     
            clv                ; $e567: b8        
            clc                ; $e568: 18        
            lda #$00           ; $e569: a9 00     
            .hex a3 40         ; $e56b: a3 40     Invalid Opcode - LAX ($40,x)
            nop                ; $e56d: ea        
            nop                ; $e56e: ea        
            nop                ; $e56f: ea        
            nop                ; $e570: ea        
            beq __e585         ; $e571: f0 12     
            bpl __e585         ; $e573: 10 10     
            bvs __e585         ; $e575: 70 0e     
            bcs __e585         ; $e577: b0 0c     
            cmp #$aa           ; $e579: c9 aa     
            bne __e585         ; $e57b: d0 08     
            cpx #$aa           ; $e57d: e0 aa     
            bne __e585         ; $e57f: d0 04     
            cpy #$33           ; $e581: c0 33     
            beq __e589         ; $e583: f0 04     
__e585:     lda #$7d           ; $e585: a9 7d     
            sta $00            ; $e587: 85 00     
__e589:     lda #$87           ; $e589: a9 87     
            sta $67            ; $e58b: 85 67     
            lda #$32           ; $e58d: a9 32     
            sta $68            ; $e58f: 85 68     
            ldy #$57           ; $e591: a0 57     
            bit $01            ; $e593: 24 01     
            sec                ; $e595: 38        
            lda #$00           ; $e596: a9 00     
            .hex a7 67         ; $e598: a7 67     Invalid Opcode - LAX $67
            nop                ; $e59a: ea        
            nop                ; $e59b: ea        
            nop                ; $e59c: ea        
            nop                ; $e59d: ea        
            beq __e5b2         ; $e59e: f0 12     
            bpl __e5b2         ; $e5a0: 10 10     
            bvc __e5b2         ; $e5a2: 50 0e     
            bcc __e5b2         ; $e5a4: 90 0c     
            cmp #$87           ; $e5a6: c9 87     
            bne __e5b2         ; $e5a8: d0 08     
            cpx #$87           ; $e5aa: e0 87     
            bne __e5b2         ; $e5ac: d0 04     
            cpy #$57           ; $e5ae: c0 57     
            beq __e5b6         ; $e5b0: f0 04     
__e5b2:     lda #$7e           ; $e5b2: a9 7e     
            sta $00            ; $e5b4: 85 00     
__e5b6:     ldy #$53           ; $e5b6: a0 53     
            clv                ; $e5b8: b8        
            clc                ; $e5b9: 18        
            lda #$00           ; $e5ba: a9 00     
            .hex a7 68         ; $e5bc: a7 68     Invalid Opcode - LAX $68
            nop                ; $e5be: ea        
            nop                ; $e5bf: ea        
            nop                ; $e5c0: ea        
            nop                ; $e5c1: ea        
            beq __e5d6         ; $e5c2: f0 12     
            bmi __e5d6         ; $e5c4: 30 10     
            bvs __e5d6         ; $e5c6: 70 0e     
            bcs __e5d6         ; $e5c8: b0 0c     
            cmp #$32           ; $e5ca: c9 32     
            bne __e5d6         ; $e5cc: d0 08     
            cpx #$32           ; $e5ce: e0 32     
            bne __e5d6         ; $e5d0: d0 04     
            cpy #$53           ; $e5d2: c0 53     
            beq __e5da         ; $e5d4: f0 04     
__e5d6:     lda #$7f           ; $e5d6: a9 7f     
            sta $00            ; $e5d8: 85 00     
__e5da:     lda #$87           ; $e5da: a9 87     
            sta $0577          ; $e5dc: 8d 77 05  
            lda #$32           ; $e5df: a9 32     
            sta $0578          ; $e5e1: 8d 78 05  
            ldy #$57           ; $e5e4: a0 57     
            bit $01            ; $e5e6: 24 01     
            sec                ; $e5e8: 38        
            lda #$00           ; $e5e9: a9 00     
            .hex af 77 05      ; $e5eb: af 77 05  Invalid Opcode - LAX $0577
            nop                ; $e5ee: ea        
            nop                ; $e5ef: ea        
            nop                ; $e5f0: ea        
            nop                ; $e5f1: ea        
            beq __e606         ; $e5f2: f0 12     
            bpl __e606         ; $e5f4: 10 10     
            bvc __e606         ; $e5f6: 50 0e     
            bcc __e606         ; $e5f8: 90 0c     
            cmp #$87           ; $e5fa: c9 87     
            bne __e606         ; $e5fc: d0 08     
            cpx #$87           ; $e5fe: e0 87     
            bne __e606         ; $e600: d0 04     
            cpy #$57           ; $e602: c0 57     
            beq __e60a         ; $e604: f0 04     
__e606:     lda #$80           ; $e606: a9 80     
            sta $00            ; $e608: 85 00     
__e60a:     ldy #$53           ; $e60a: a0 53     
            clv                ; $e60c: b8        
            clc                ; $e60d: 18        
            lda #$00           ; $e60e: a9 00     
            .hex af 78 05      ; $e610: af 78 05  Invalid Opcode - LAX $0578
            nop                ; $e613: ea        
            nop                ; $e614: ea        
            nop                ; $e615: ea        
            nop                ; $e616: ea        
            beq __e62b         ; $e617: f0 12     
            bmi __e62b         ; $e619: 30 10     
            bvs __e62b         ; $e61b: 70 0e     
            bcs __e62b         ; $e61d: b0 0c     
            cmp #$32           ; $e61f: c9 32     
            bne __e62b         ; $e621: d0 08     
            cpx #$32           ; $e623: e0 32     
            bne __e62b         ; $e625: d0 04     
            cpy #$53           ; $e627: c0 53     
            beq __e62f         ; $e629: f0 04     
__e62b:     lda #$81           ; $e62b: a9 81     
            sta $00            ; $e62d: 85 00     
__e62f:     lda #$ff           ; $e62f: a9 ff     
            sta $43            ; $e631: 85 43     
            lda #$04           ; $e633: a9 04     
            sta $44            ; $e635: 85 44     
            lda #$32           ; $e637: a9 32     
            sta $45            ; $e639: 85 45     
            lda #$04           ; $e63b: a9 04     
            sta $46            ; $e63d: 85 46     
            lda #$55           ; $e63f: a9 55     
            sta $0580          ; $e641: 8d 80 05  
            lda #$aa           ; $e644: a9 aa     
            sta $0432          ; $e646: 8d 32 04  
            ldx #$03           ; $e649: a2 03     
            ldy #$81           ; $e64b: a0 81     
            bit $01            ; $e64d: 24 01     
            sec                ; $e64f: 38        
            lda #$00           ; $e650: a9 00     
            .hex b3 43         ; $e652: b3 43     Invalid Opcode - LAX ($43),y
            nop                ; $e654: ea        
            nop                ; $e655: ea        
            nop                ; $e656: ea        
            nop                ; $e657: ea        
            beq __e66c         ; $e658: f0 12     
            bmi __e66c         ; $e65a: 30 10     
            bvc __e66c         ; $e65c: 50 0e     
            bcc __e66c         ; $e65e: 90 0c     
            cmp #$55           ; $e660: c9 55     
            bne __e66c         ; $e662: d0 08     
            cpx #$55           ; $e664: e0 55     
            bne __e66c         ; $e666: d0 04     
            cpy #$81           ; $e668: c0 81     
            beq __e670         ; $e66a: f0 04     
__e66c:     lda #$82           ; $e66c: a9 82     
            sta $00            ; $e66e: 85 00     
__e670:     ldx #$05           ; $e670: a2 05     
            ldy #$00           ; $e672: a0 00     
            clv                ; $e674: b8        
            clc                ; $e675: 18        
            lda #$00           ; $e676: a9 00     
            .hex b3 45         ; $e678: b3 45     Invalid Opcode - LAX ($45),y
            nop                ; $e67a: ea        
            nop                ; $e67b: ea        
            nop                ; $e67c: ea        
            nop                ; $e67d: ea        
            beq __e692         ; $e67e: f0 12     
            bpl __e692         ; $e680: 10 10     
            bvs __e692         ; $e682: 70 0e     
            bcs __e692         ; $e684: b0 0c     
            cmp #$aa           ; $e686: c9 aa     
            bne __e692         ; $e688: d0 08     
            cpx #$aa           ; $e68a: e0 aa     
            bne __e692         ; $e68c: d0 04     
            cpy #$00           ; $e68e: c0 00     
            beq __e696         ; $e690: f0 04     
__e692:     lda #$83           ; $e692: a9 83     
            sta $00            ; $e694: 85 00     
__e696:     lda #$87           ; $e696: a9 87     
            sta $67            ; $e698: 85 67     
            lda #$32           ; $e69a: a9 32     
            sta $68            ; $e69c: 85 68     
            ldy #$57           ; $e69e: a0 57     
            bit $01            ; $e6a0: 24 01     
            sec                ; $e6a2: 38        
            lda #$00           ; $e6a3: a9 00     
            .hex b7 10         ; $e6a5: b7 10     Invalid Opcode - LAX $10,y
            nop                ; $e6a7: ea        
            nop                ; $e6a8: ea        
            nop                ; $e6a9: ea        
            nop                ; $e6aa: ea        
            beq __e6bf         ; $e6ab: f0 12     
            bpl __e6bf         ; $e6ad: 10 10     
            bvc __e6bf         ; $e6af: 50 0e     
            bcc __e6bf         ; $e6b1: 90 0c     
            cmp #$87           ; $e6b3: c9 87     
            bne __e6bf         ; $e6b5: d0 08     
            cpx #$87           ; $e6b7: e0 87     
            bne __e6bf         ; $e6b9: d0 04     
            cpy #$57           ; $e6bb: c0 57     
            beq __e6c3         ; $e6bd: f0 04     
__e6bf:     lda #$84           ; $e6bf: a9 84     
            sta $00            ; $e6c1: 85 00     
__e6c3:     ldy #$ff           ; $e6c3: a0 ff     
            clv                ; $e6c5: b8        
            clc                ; $e6c6: 18        
            lda #$00           ; $e6c7: a9 00     
            .hex b7 69         ; $e6c9: b7 69     Invalid Opcode - LAX $69,y
            nop                ; $e6cb: ea        
            nop                ; $e6cc: ea        
            nop                ; $e6cd: ea        
            nop                ; $e6ce: ea        
            beq __e6e3         ; $e6cf: f0 12     
            bmi __e6e3         ; $e6d1: 30 10     
            bvs __e6e3         ; $e6d3: 70 0e     
            bcs __e6e3         ; $e6d5: b0 0c     
            cmp #$32           ; $e6d7: c9 32     
            bne __e6e3         ; $e6d9: d0 08     
            cpx #$32           ; $e6db: e0 32     
            bne __e6e3         ; $e6dd: d0 04     
            .hex c0            ; $e6df: c0        Suspected data
__e6e0:     .hex ff f0 04      ; $e6e0: ff f0 04  Invalid Opcode - ISC $04f0,x
__e6e3:     lda #$85           ; $e6e3: a9 85     
            sta $00            ; $e6e5: 85 00     
            lda #$87           ; $e6e7: a9 87     
            sta $0587          ; $e6e9: 8d 87 05  
            lda #$32           ; $e6ec: a9 32     
            sta $0588          ; $e6ee: 8d 88 05  
            ldy #$30           ; $e6f1: a0 30     
            bit $01            ; $e6f3: 24 01     
            sec                ; $e6f5: 38        
            lda #$00           ; $e6f6: a9 00     
            .hex bf 57 05      ; $e6f8: bf 57 05  Invalid Opcode - LAX $0557,y
            nop                ; $e6fb: ea        
            nop                ; $e6fc: ea        
            nop                ; $e6fd: ea        
            nop                ; $e6fe: ea        
            beq __e713         ; $e6ff: f0 12     
            bpl __e713         ; $e701: 10 10     
            bvc __e713         ; $e703: 50 0e     
            bcc __e713         ; $e705: 90 0c     
            cmp #$87           ; $e707: c9 87     
            bne __e713         ; $e709: d0 08     
            cpx #$87           ; $e70b: e0 87     
            bne __e713         ; $e70d: d0 04     
            cpy #$30           ; $e70f: c0 30     
            beq __e717         ; $e711: f0 04     
__e713:     lda #$86           ; $e713: a9 86     
            sta $00            ; $e715: 85 00     
__e717:     ldy #$40           ; $e717: a0 40     
            clv                ; $e719: b8        
            clc                ; $e71a: 18        
            lda #$00           ; $e71b: a9 00     
            .hex bf 48 05      ; $e71d: bf 48 05  Invalid Opcode - LAX $0548,y
            nop                ; $e720: ea        
            nop                ; $e721: ea        
            nop                ; $e722: ea        
            nop                ; $e723: ea        
            beq __e738         ; $e724: f0 12     
            bmi __e738         ; $e726: 30 10     
            bvs __e738         ; $e728: 70 0e     
            bcs __e738         ; $e72a: b0 0c     
            cmp #$32           ; $e72c: c9 32     
            bne __e738         ; $e72e: d0 08     
            cpx #$32           ; $e730: e0 32     
            bne __e738         ; $e732: d0 04     
            cpy #$40           ; $e734: c0 40     
            beq __e73c         ; $e736: f0 04     
__e738:     lda #$87           ; $e738: a9 87     
            sta $00            ; $e73a: 85 00     
__e73c:     rts                ; $e73c: 60        

;-------------------------------------------------------------------------------
__e73d:     lda #$c0           ; $e73d: a9 c0     
            sta $01            ; $e73f: 85 01     
            lda #$00           ; $e741: a9 00     
            sta $0489          ; $e743: 8d 89 04  
            lda #$89           ; $e746: a9 89     
            sta $60            ; $e748: 85 60     
            lda #$04           ; $e74a: a9 04     
            sta $61            ; $e74c: 85 61     
            ldy #$44           ; $e74e: a0 44     
            ldx #$17           ; $e750: a2 17     
            lda #$3e           ; $e752: a9 3e     
            bit $01            ; $e754: 24 01     
            clc                ; $e756: 18        
            .hex 83 49         ; $e757: 83 49     Invalid Opcode - SAX ($49,x)
            nop                ; $e759: ea        
            nop                ; $e75a: ea        
            nop                ; $e75b: ea        
            nop                ; $e75c: ea        
            bne __e778         ; $e75d: d0 19     
            bcs __e778         ; $e75f: b0 17     
            bvc __e778         ; $e761: 50 15     
            bpl __e778         ; $e763: 10 13     
            cmp #$3e           ; $e765: c9 3e     
            bne __e778         ; $e767: d0 0f     
            cpy #$44           ; $e769: c0 44     
            bne __e778         ; $e76b: d0 0b     
            cpx #$17           ; $e76d: e0 17     
            bne __e778         ; $e76f: d0 07     
            lda $0489          ; $e771: ad 89 04  
            cmp #$16           ; $e774: c9 16     
            beq __e77c         ; $e776: f0 04     
__e778:     lda #$88           ; $e778: a9 88     
            sta $00            ; $e77a: 85 00     
__e77c:     ldy #$44           ; $e77c: a0 44     
            ldx #$7a           ; $e77e: a2 7a     
            lda #$66           ; $e780: a9 66     
            sec                ; $e782: 38        
            clv                ; $e783: b8        
            .hex 83 e6         ; $e784: 83 e6     Invalid Opcode - SAX ($e6,x)
            nop                ; $e786: ea        
            nop                ; $e787: ea        
            nop                ; $e788: ea        
            nop                ; $e789: ea        
            beq __e7a5         ; $e78a: f0 19     
            bcc __e7a5         ; $e78c: 90 17     
            bvs __e7a5         ; $e78e: 70 15     
            bmi __e7a5         ; $e790: 30 13     
            cmp #$66           ; $e792: c9 66     
            bne __e7a5         ; $e794: d0 0f     
            cpy #$44           ; $e796: c0 44     
            bne __e7a5         ; $e798: d0 0b     
            cpx #$7a           ; $e79a: e0 7a     
            bne __e7a5         ; $e79c: d0 07     
            lda $0489          ; $e79e: ad 89 04  
            cmp #$62           ; $e7a1: c9 62     
            beq __e7a9         ; $e7a3: f0 04     
__e7a5:     lda #$89           ; $e7a5: a9 89     
            sta $00            ; $e7a7: 85 00     
__e7a9:     lda #$ff           ; $e7a9: a9 ff     
            sta $49            ; $e7ab: 85 49     
            ldy #$44           ; $e7ad: a0 44     
            ldx #$aa           ; $e7af: a2 aa     
            lda #$55           ; $e7b1: a9 55     
            bit $01            ; $e7b3: 24 01     
            clc                ; $e7b5: 18        
            .hex 87 49         ; $e7b6: 87 49     Invalid Opcode - SAX $49
            nop                ; $e7b8: ea        
            nop                ; $e7b9: ea        
            nop                ; $e7ba: ea        
            nop                ; $e7bb: ea        
            beq __e7d6         ; $e7bc: f0 18     
            bcs __e7d6         ; $e7be: b0 16     
            bvc __e7d6         ; $e7c0: 50 14     
            bpl __e7d6         ; $e7c2: 10 12     
            cmp #$55           ; $e7c4: c9 55     
            bne __e7d6         ; $e7c6: d0 0e     
            cpy #$44           ; $e7c8: c0 44     
            bne __e7d6         ; $e7ca: d0 0a     
            cpx #$aa           ; $e7cc: e0 aa     
            bne __e7d6         ; $e7ce: d0 06     
            lda $49            ; $e7d0: a5 49     
            cmp #$00           ; $e7d2: c9 00     
            beq __e7da         ; $e7d4: f0 04     
__e7d6:     lda #$8a           ; $e7d6: a9 8a     
            sta $00            ; $e7d8: 85 00     
__e7da:     lda #$00           ; $e7da: a9 00     
            sta $56            ; $e7dc: 85 56     
            ldy #$58           ; $e7de: a0 58     
            ldx #$ef           ; $e7e0: a2 ef     
            lda #$66           ; $e7e2: a9 66     
            sec                ; $e7e4: 38        
            clv                ; $e7e5: b8        
            .hex 87 56         ; $e7e6: 87 56     Invalid Opcode - SAX $56
            nop                ; $e7e8: ea        
            nop                ; $e7e9: ea        
            nop                ; $e7ea: ea        
            nop                ; $e7eb: ea        
            beq __e806         ; $e7ec: f0 18     
            bcc __e806         ; $e7ee: 90 16     
            bvs __e806         ; $e7f0: 70 14     
            bmi __e806         ; $e7f2: 30 12     
            cmp #$66           ; $e7f4: c9 66     
            bne __e806         ; $e7f6: d0 0e     
            cpy #$58           ; $e7f8: c0 58     
            bne __e806         ; $e7fa: d0 0a     
            cpx #$ef           ; $e7fc: e0 ef     
            bne __e806         ; $e7fe: d0 06     
            lda $56            ; $e800: a5 56     
            cmp #$66           ; $e802: c9 66     
            beq __e80a         ; $e804: f0 04     
__e806:     lda #$8b           ; $e806: a9 8b     
            sta $00            ; $e808: 85 00     
__e80a:     lda #$ff           ; $e80a: a9 ff     
            sta $0549          ; $e80c: 8d 49 05  
            ldy #$e5           ; $e80f: a0 e5     
            ldx #$af           ; $e811: a2 af     
            lda #$f5           ; $e813: a9 f5     
            bit $01            ; $e815: 24 01     
            clc                ; $e817: 18        
            .hex 8f 49 05      ; $e818: 8f 49 05  Invalid Opcode - SAX $0549
            nop                ; $e81b: ea        
            nop                ; $e81c: ea        
            nop                ; $e81d: ea        
            nop                ; $e81e: ea        
            beq __e83a         ; $e81f: f0 19     
            bcs __e83a         ; $e821: b0 17     
            bvc __e83a         ; $e823: 50 15     
            bpl __e83a         ; $e825: 10 13     
            cmp #$f5           ; $e827: c9 f5     
            bne __e83a         ; $e829: d0 0f     
            cpy #$e5           ; $e82b: c0 e5     
            bne __e83a         ; $e82d: d0 0b     
            cpx #$af           ; $e82f: e0 af     
            bne __e83a         ; $e831: d0 07     
            lda $0549          ; $e833: ad 49 05  
            cmp #$a5           ; $e836: c9 a5     
            beq __e83e         ; $e838: f0 04     
__e83a:     lda #$8c           ; $e83a: a9 8c     
            sta $00            ; $e83c: 85 00     
__e83e:     lda #$00           ; $e83e: a9 00     
            sta $0556          ; $e840: 8d 56 05  
            ldy #$58           ; $e843: a0 58     
            ldx #$b3           ; $e845: a2 b3     
            lda #$97           ; $e847: a9 97     
            sec                ; $e849: 38        
            clv                ; $e84a: b8        
            .hex 8f 56 05      ; $e84b: 8f 56 05  Invalid Opcode - SAX $0556
            nop                ; $e84e: ea        
            nop                ; $e84f: ea        
            nop                ; $e850: ea        
            nop                ; $e851: ea        
            beq __e86d         ; $e852: f0 19     
            bcc __e86d         ; $e854: 90 17     
            bvs __e86d         ; $e856: 70 15     
            bpl __e86d         ; $e858: 10 13     
            cmp #$97           ; $e85a: c9 97     
            bne __e86d         ; $e85c: d0 0f     
            cpy #$58           ; $e85e: c0 58     
            bne __e86d         ; $e860: d0 0b     
            cpx #$b3           ; $e862: e0 b3     
            bne __e86d         ; $e864: d0 07     
            lda $0556          ; $e866: ad 56 05  
            cmp #$93           ; $e869: c9 93     
            beq __e871         ; $e86b: f0 04     
__e86d:     lda #$8d           ; $e86d: a9 8d     
            sta $00            ; $e86f: 85 00     
__e871:     lda #$ff           ; $e871: a9 ff     
            sta $49            ; $e873: 85 49     
            ldy #$ff           ; $e875: a0 ff     
            ldx #$aa           ; $e877: a2 aa     
            lda #$55           ; $e879: a9 55     
            bit $01            ; $e87b: 24 01     
            clc                ; $e87d: 18        
            .hex 97 4a         ; $e87e: 97 4a     Invalid Opcode - SAX $4a,y
            nop                ; $e880: ea        
            nop                ; $e881: ea        
            nop                ; $e882: ea        
            nop                ; $e883: ea        
            beq __e89e         ; $e884: f0 18     
            bcs __e89e         ; $e886: b0 16     
            bvc __e89e         ; $e888: 50 14     
            bpl __e89e         ; $e88a: 10 12     
            cmp #$55           ; $e88c: c9 55     
            bne __e89e         ; $e88e: d0 0e     
            cpy #$ff           ; $e890: c0 ff     
            bne __e89e         ; $e892: d0 0a     
            cpx #$aa           ; $e894: e0 aa     
            bne __e89e         ; $e896: d0 06     
            lda $49            ; $e898: a5 49     
            cmp #$00           ; $e89a: c9 00     
            beq __e8a2         ; $e89c: f0 04     
__e89e:     lda #$8e           ; $e89e: a9 8e     
            sta $00            ; $e8a0: 85 00     
__e8a2:     lda #$00           ; $e8a2: a9 00     
            sta $56            ; $e8a4: 85 56     
            ldy #$06           ; $e8a6: a0 06     
            ldx #$ef           ; $e8a8: a2 ef     
            lda #$66           ; $e8aa: a9 66     
            sec                ; $e8ac: 38        
            clv                ; $e8ad: b8        
            .hex 97 50         ; $e8ae: 97 50     Invalid Opcode - SAX $50,y
            nop                ; $e8b0: ea        
            nop                ; $e8b1: ea        
            nop                ; $e8b2: ea        
            nop                ; $e8b3: ea        
            beq __e8ce         ; $e8b4: f0 18     
            bcc __e8ce         ; $e8b6: 90 16     
            bvs __e8ce         ; $e8b8: 70 14     
            bmi __e8ce         ; $e8ba: 30 12     
            cmp #$66           ; $e8bc: c9 66     
            bne __e8ce         ; $e8be: d0 0e     
            cpy #$06           ; $e8c0: c0 06     
            bne __e8ce         ; $e8c2: d0 0a     
            cpx #$ef           ; $e8c4: e0 ef     
            bne __e8ce         ; $e8c6: d0 06     
            lda $56            ; $e8c8: a5 56     
            cmp #$66           ; $e8ca: c9 66     
            beq __e8d2         ; $e8cc: f0 04     
__e8ce:     lda #$8f           ; $e8ce: a9 8f     
            sta $00            ; $e8d0: 85 00     
__e8d2:     rts                ; $e8d2: 60        

;-------------------------------------------------------------------------------
__e8d3:     ldy #$90           ; $e8d3: a0 90     
            jsr __f931         ; $e8d5: 20 31 f9  
            .hex eb 40         ; $e8d8: eb 40     Invalid Opcode - SBC #$40
            nop                ; $e8da: ea        
            nop                ; $e8db: ea        
            nop                ; $e8dc: ea        
            nop                ; $e8dd: ea        
            jsr __f937         ; $e8de: 20 37 f9  
            iny                ; $e8e1: c8        
            jsr __f947         ; $e8e2: 20 47 f9  
            .hex eb 3f         ; $e8e5: eb 3f     Invalid Opcode - SBC #$3f
            nop                ; $e8e7: ea        
            nop                ; $e8e8: ea        
            nop                ; $e8e9: ea        
            nop                ; $e8ea: ea        
            jsr __f94c         ; $e8eb: 20 4c f9  
            iny                ; $e8ee: c8        
            jsr __f95c         ; $e8ef: 20 5c f9  
            .hex eb 41         ; $e8f2: eb 41     Invalid Opcode - SBC #$41
            nop                ; $e8f4: ea        
            nop                ; $e8f5: ea        
            nop                ; $e8f6: ea        
            nop                ; $e8f7: ea        
            jsr __f962         ; $e8f8: 20 62 f9  
            iny                ; $e8fb: c8        
            jsr __f972         ; $e8fc: 20 72 f9  
            .hex eb 00         ; $e8ff: eb 00     Invalid Opcode - SBC #$00
            nop                ; $e901: ea        
            nop                ; $e902: ea        
            nop                ; $e903: ea        
            nop                ; $e904: ea        
            jsr __f976         ; $e905: 20 76 f9  
            iny                ; $e908: c8        
            jsr __f980         ; $e909: 20 80 f9  
            .hex eb 7f         ; $e90c: eb 7f     Invalid Opcode - SBC #$7f
            nop                ; $e90e: ea        
            nop                ; $e90f: ea        
            nop                ; $e910: ea        
            nop                ; $e911: ea        
            jsr __f984         ; $e912: 20 84 f9  
            rts                ; $e915: 60        

;-------------------------------------------------------------------------------
__e916:     lda #$ff           ; $e916: a9 ff     
            sta $01            ; $e918: 85 01     
            ldy #$95           ; $e91a: a0 95     
            ldx #$02           ; $e91c: a2 02     
            lda #$47           ; $e91e: a9 47     
            sta $47            ; $e920: 85 47     
            lda #$06           ; $e922: a9 06     
            sta $48            ; $e924: 85 48     
            lda #$eb           ; $e926: a9 eb     
            sta $0647          ; $e928: 8d 47 06  
            jsr __fa31         ; $e92b: 20 31 fa  
            .hex c3 45         ; $e92e: c3 45     Invalid Opcode - DCP ($45,x)
            nop                ; $e930: ea        
            nop                ; $e931: ea        
            nop                ; $e932: ea        
            nop                ; $e933: ea        
            jsr __fa37         ; $e934: 20 37 fa  
            lda $0647          ; $e937: ad 47 06  
            cmp #$ea           ; $e93a: c9 ea     
            beq __e940         ; $e93c: f0 02     
            sty $00            ; $e93e: 84 00     
__e940:     iny                ; $e940: c8        
            lda #$00           ; $e941: a9 00     
            sta $0647          ; $e943: 8d 47 06  
            jsr __fa42         ; $e946: 20 42 fa  
            .hex c3 45         ; $e949: c3 45     Invalid Opcode - DCP ($45,x)
            nop                ; $e94b: ea        
            nop                ; $e94c: ea        
            nop                ; $e94d: ea        
            nop                ; $e94e: ea        
            jsr __fa47         ; $e94f: 20 47 fa  
            lda $0647          ; $e952: ad 47 06  
            cmp #$ff           ; $e955: c9 ff     
            beq __e95b         ; $e957: f0 02     
            sty $00            ; $e959: 84 00     
__e95b:     iny                ; $e95b: c8        
            lda #$37           ; $e95c: a9 37     
            sta $0647          ; $e95e: 8d 47 06  
            jsr __fa54         ; $e961: 20 54 fa  
            .hex c3 45         ; $e964: c3 45     Invalid Opcode - DCP ($45,x)
            nop                ; $e966: ea        
            nop                ; $e967: ea        
            nop                ; $e968: ea        
            nop                ; $e969: ea        
            jsr __fa59         ; $e96a: 20 59 fa  
            lda $0647          ; $e96d: ad 47 06  
            cmp #$36           ; $e970: c9 36     
            beq __e976         ; $e972: f0 02     
            sty $00            ; $e974: 84 00     
__e976:     iny                ; $e976: c8        
            lda #$eb           ; $e977: a9 eb     
            sta $47            ; $e979: 85 47     
            jsr __fa31         ; $e97b: 20 31 fa  
            .hex c7 47         ; $e97e: c7 47     Invalid Opcode - DCP $47
            nop                ; $e980: ea        
            nop                ; $e981: ea        
            nop                ; $e982: ea        
            nop                ; $e983: ea        
            jsr __fa37         ; $e984: 20 37 fa  
            lda $47            ; $e987: a5 47     
            cmp #$ea           ; $e989: c9 ea     
            beq __e98f         ; $e98b: f0 02     
            sty $00            ; $e98d: 84 00     
__e98f:     iny                ; $e98f: c8        
            lda #$00           ; $e990: a9 00     
            sta $47            ; $e992: 85 47     
            jsr __fa42         ; $e994: 20 42 fa  
            .hex c7 47         ; $e997: c7 47     Invalid Opcode - DCP $47
            nop                ; $e999: ea        
            nop                ; $e99a: ea        
            nop                ; $e99b: ea        
            nop                ; $e99c: ea        
            jsr __fa47         ; $e99d: 20 47 fa  
            lda $47            ; $e9a0: a5 47     
            cmp #$ff           ; $e9a2: c9 ff     
            beq __e9a8         ; $e9a4: f0 02     
            sty $00            ; $e9a6: 84 00     
__e9a8:     iny                ; $e9a8: c8        
            lda #$37           ; $e9a9: a9 37     
            sta $47            ; $e9ab: 85 47     
            jsr __fa54         ; $e9ad: 20 54 fa  
            .hex c7 47         ; $e9b0: c7 47     Invalid Opcode - DCP $47
            nop                ; $e9b2: ea        
            nop                ; $e9b3: ea        
            nop                ; $e9b4: ea        
            nop                ; $e9b5: ea        
            jsr __fa59         ; $e9b6: 20 59 fa  
            lda $47            ; $e9b9: a5 47     
            cmp #$36           ; $e9bb: c9 36     
            beq __e9c1         ; $e9bd: f0 02     
            sty $00            ; $e9bf: 84 00     
__e9c1:     iny                ; $e9c1: c8        
            lda #$eb           ; $e9c2: a9 eb     
            sta $0647          ; $e9c4: 8d 47 06  
            jsr __fa31         ; $e9c7: 20 31 fa  
            .hex cf 47 06      ; $e9ca: cf 47 06  Invalid Opcode - DCP $0647
            nop                ; $e9cd: ea        
            nop                ; $e9ce: ea        
            nop                ; $e9cf: ea        
            nop                ; $e9d0: ea        
            jsr __fa37         ; $e9d1: 20 37 fa  
            lda $0647          ; $e9d4: ad 47 06  
            cmp #$ea           ; $e9d7: c9 ea     
            beq __e9dd         ; $e9d9: f0 02     
            sty $00            ; $e9db: 84 00     
__e9dd:     iny                ; $e9dd: c8        
            lda #$00           ; $e9de: a9 00     
            sta $0647          ; $e9e0: 8d 47 06  
            jsr __fa42         ; $e9e3: 20 42 fa  
            .hex cf 47 06      ; $e9e6: cf 47 06  Invalid Opcode - DCP $0647
            nop                ; $e9e9: ea        
            nop                ; $e9ea: ea        
            nop                ; $e9eb: ea        
            nop                ; $e9ec: ea        
            jsr __fa47         ; $e9ed: 20 47 fa  
            lda $0647          ; $e9f0: ad 47 06  
            cmp #$ff           ; $e9f3: c9 ff     
            beq __e9f9         ; $e9f5: f0 02     
            sty $00            ; $e9f7: 84 00     
__e9f9:     iny                ; $e9f9: c8        
            lda #$37           ; $e9fa: a9 37     
            sta $0647          ; $e9fc: 8d 47 06  
            jsr __fa54         ; $e9ff: 20 54 fa  
            .hex cf 47 06      ; $ea02: cf 47 06  Invalid Opcode - DCP $0647
            nop                ; $ea05: ea        
            nop                ; $ea06: ea        
            nop                ; $ea07: ea        
            nop                ; $ea08: ea        
            jsr __fa59         ; $ea09: 20 59 fa  
            lda $0647          ; $ea0c: ad 47 06  
            cmp #$36           ; $ea0f: c9 36     
            beq __ea15         ; $ea11: f0 02     
            sty $00            ; $ea13: 84 00     
__ea15:     lda #$eb           ; $ea15: a9 eb     
            sta $0647          ; $ea17: 8d 47 06  
            lda #$48           ; $ea1a: a9 48     
            sta $45            ; $ea1c: 85 45     
            lda #$05           ; $ea1e: a9 05     
            sta $46            ; $ea20: 85 46     
            ldy #$ff           ; $ea22: a0 ff     
            jsr __fa31         ; $ea24: 20 31 fa  
            .hex d3 45         ; $ea27: d3 45     Invalid Opcode - DCP ($45),y
            nop                ; $ea29: ea        
            nop                ; $ea2a: ea        
            php                ; $ea2b: 08        
            pha                ; $ea2c: 48        
            ldy #$9e           ; $ea2d: a0 9e     
            pla                ; $ea2f: 68        
            plp                ; $ea30: 28        
            jsr __fa37         ; $ea31: 20 37 fa  
            lda $0647          ; $ea34: ad 47 06  
            cmp #$ea           ; $ea37: c9 ea     
            beq __ea3d         ; $ea39: f0 02     
            sty $00            ; $ea3b: 84 00     
__ea3d:     ldy #$ff           ; $ea3d: a0 ff     
            lda #$00           ; $ea3f: a9 00     
            sta $0647          ; $ea41: 8d 47 06  
            jsr __fa42         ; $ea44: 20 42 fa  
            .hex d3 45         ; $ea47: d3 45     Invalid Opcode - DCP ($45),y
            nop                ; $ea49: ea        
            nop                ; $ea4a: ea        
            php                ; $ea4b: 08        
            pha                ; $ea4c: 48        
            ldy #$9f           ; $ea4d: a0 9f     
            pla                ; $ea4f: 68        
            plp                ; $ea50: 28        
            jsr __fa47         ; $ea51: 20 47 fa  
            lda $0647          ; $ea54: ad 47 06  
            cmp #$ff           ; $ea57: c9 ff     
            beq __ea5d         ; $ea59: f0 02     
            sty $00            ; $ea5b: 84 00     
__ea5d:     ldy #$ff           ; $ea5d: a0 ff     
            lda #$37           ; $ea5f: a9 37     
            sta $0647          ; $ea61: 8d 47 06  
            jsr __fa54         ; $ea64: 20 54 fa  
            .hex d3 45         ; $ea67: d3 45     Invalid Opcode - DCP ($45),y
            nop                ; $ea69: ea        
            nop                ; $ea6a: ea        
            php                ; $ea6b: 08        
            pha                ; $ea6c: 48        
            ldy #$a0           ; $ea6d: a0 a0     
            pla                ; $ea6f: 68        
            plp                ; $ea70: 28        
            jsr __fa59         ; $ea71: 20 59 fa  
            lda $0647          ; $ea74: ad 47 06  
            cmp #$36           ; $ea77: c9 36     
            beq __ea7d         ; $ea79: f0 02     
            sty $00            ; $ea7b: 84 00     
__ea7d:     ldy #$a1           ; $ea7d: a0 a1     
            ldx #$ff           ; $ea7f: a2 ff     
            lda #$eb           ; $ea81: a9 eb     
            sta $47            ; $ea83: 85 47     
            jsr __fa31         ; $ea85: 20 31 fa  
            .hex d7 48         ; $ea88: d7 48     Invalid Opcode - DCP $48,x
            nop                ; $ea8a: ea        
            nop                ; $ea8b: ea        
            nop                ; $ea8c: ea        
            nop                ; $ea8d: ea        
            jsr __fa37         ; $ea8e: 20 37 fa  
            lda $47            ; $ea91: a5 47     
            cmp #$ea           ; $ea93: c9 ea     
            beq __ea99         ; $ea95: f0 02     
            sty $00            ; $ea97: 84 00     
__ea99:     iny                ; $ea99: c8        
            lda #$00           ; $ea9a: a9 00     
            sta $47            ; $ea9c: 85 47     
            jsr __fa42         ; $ea9e: 20 42 fa  
            .hex d7 48         ; $eaa1: d7 48     Invalid Opcode - DCP $48,x
            nop                ; $eaa3: ea        
            nop                ; $eaa4: ea        
            nop                ; $eaa5: ea        
            nop                ; $eaa6: ea        
            jsr __fa47         ; $eaa7: 20 47 fa  
            lda $47            ; $eaaa: a5 47     
            cmp #$ff           ; $eaac: c9 ff     
            beq __eab2         ; $eaae: f0 02     
            sty $00            ; $eab0: 84 00     
__eab2:     iny                ; $eab2: c8        
            lda #$37           ; $eab3: a9 37     
            sta $47            ; $eab5: 85 47     
            jsr __fa54         ; $eab7: 20 54 fa  
            .hex d7 48         ; $eaba: d7 48     Invalid Opcode - DCP $48,x
            nop                ; $eabc: ea        
            nop                ; $eabd: ea        
            nop                ; $eabe: ea        
            nop                ; $eabf: ea        
            jsr __fa59         ; $eac0: 20 59 fa  
            lda $47            ; $eac3: a5 47     
            cmp #$36           ; $eac5: c9 36     
            beq __eacb         ; $eac7: f0 02     
            sty $00            ; $eac9: 84 00     
__eacb:     lda #$eb           ; $eacb: a9 eb     
            sta $0647          ; $eacd: 8d 47 06  
            ldy #$ff           ; $ead0: a0 ff     
            jsr __fa31         ; $ead2: 20 31 fa  
            .hex db 48 05      ; $ead5: db 48 05  Invalid Opcode - DCP $0548,y
            nop                ; $ead8: ea        
            nop                ; $ead9: ea        
            php                ; $eada: 08        
            pha                ; $eadb: 48        
            ldy #$a4           ; $eadc: a0 a4     
            pla                ; $eade: 68        
            plp                ; $eadf: 28        
            jsr __fa37         ; $eae0: 20 37 fa  
            lda $0647          ; $eae3: ad 47 06  
            cmp #$ea           ; $eae6: c9 ea     
            beq __eaec         ; $eae8: f0 02     
            sty $00            ; $eaea: 84 00     
__eaec:     ldy #$ff           ; $eaec: a0 ff     
            lda #$00           ; $eaee: a9 00     
            sta $0647          ; $eaf0: 8d 47 06  
            jsr __fa42         ; $eaf3: 20 42 fa  
            .hex db 48 05      ; $eaf6: db 48 05  Invalid Opcode - DCP $0548,y
            nop                ; $eaf9: ea        
            nop                ; $eafa: ea        
            php                ; $eafb: 08        
            pha                ; $eafc: 48        
            ldy #$a5           ; $eafd: a0 a5     
            pla                ; $eaff: 68        
            plp                ; $eb00: 28        
            jsr __fa47         ; $eb01: 20 47 fa  
            lda $0647          ; $eb04: ad 47 06  
            cmp #$ff           ; $eb07: c9 ff     
            beq __eb0d         ; $eb09: f0 02     
            sty $00            ; $eb0b: 84 00     
__eb0d:     ldy #$ff           ; $eb0d: a0 ff     
            lda #$37           ; $eb0f: a9 37     
            sta $0647          ; $eb11: 8d 47 06  
            jsr __fa54         ; $eb14: 20 54 fa  
            .hex db 48 05      ; $eb17: db 48 05  Invalid Opcode - DCP $0548,y
            nop                ; $eb1a: ea        
            nop                ; $eb1b: ea        
            php                ; $eb1c: 08        
            pha                ; $eb1d: 48        
            ldy #$a6           ; $eb1e: a0 a6     
            pla                ; $eb20: 68        
            plp                ; $eb21: 28        
            jsr __fa59         ; $eb22: 20 59 fa  
            lda $0647          ; $eb25: ad 47 06  
            cmp #$36           ; $eb28: c9 36     
            beq __eb2e         ; $eb2a: f0 02     
            sty $00            ; $eb2c: 84 00     
__eb2e:     ldy #$a7           ; $eb2e: a0 a7     
            ldx #$ff           ; $eb30: a2 ff     
            lda #$eb           ; $eb32: a9 eb     
            sta $0647          ; $eb34: 8d 47 06  
            jsr __fa31         ; $eb37: 20 31 fa  
            .hex df 48 05      ; $eb3a: df 48 05  Invalid Opcode - DCP $0548,x
            nop                ; $eb3d: ea        
            nop                ; $eb3e: ea        
            nop                ; $eb3f: ea        
            nop                ; $eb40: ea        
            jsr __fa37         ; $eb41: 20 37 fa  
            lda $0647          ; $eb44: ad 47 06  
            cmp #$ea           ; $eb47: c9 ea     
            beq __eb4d         ; $eb49: f0 02     
            sty $00            ; $eb4b: 84 00     
__eb4d:     iny                ; $eb4d: c8        
            lda #$00           ; $eb4e: a9 00     
            sta $0647          ; $eb50: 8d 47 06  
            jsr __fa42         ; $eb53: 20 42 fa  
            .hex df 48 05      ; $eb56: df 48 05  Invalid Opcode - DCP $0548,x
            nop                ; $eb59: ea        
            nop                ; $eb5a: ea        
            nop                ; $eb5b: ea        
            nop                ; $eb5c: ea        
            jsr __fa47         ; $eb5d: 20 47 fa  
            lda $0647          ; $eb60: ad 47 06  
            cmp #$ff           ; $eb63: c9 ff     
            beq __eb69         ; $eb65: f0 02     
            sty $00            ; $eb67: 84 00     
__eb69:     iny                ; $eb69: c8        
            lda #$37           ; $eb6a: a9 37     
            sta $0647          ; $eb6c: 8d 47 06  
            jsr __fa54         ; $eb6f: 20 54 fa  
            .hex df 48 05      ; $eb72: df 48 05  Invalid Opcode - DCP $0548,x
            nop                ; $eb75: ea        
            nop                ; $eb76: ea        
            nop                ; $eb77: ea        
            nop                ; $eb78: ea        
            jsr __fa59         ; $eb79: 20 59 fa  
            lda $0647          ; $eb7c: ad 47 06  
            cmp #$36           ; $eb7f: c9 36     
            beq __eb85         ; $eb81: f0 02     
            sty $00            ; $eb83: 84 00     
__eb85:     rts                ; $eb85: 60        

;-------------------------------------------------------------------------------
__eb86:     lda #$ff           ; $eb86: a9 ff     
            sta $01            ; $eb88: 85 01     
            ldy #$aa           ; $eb8a: a0 aa     
            ldx #$02           ; $eb8c: a2 02     
            lda #$47           ; $eb8e: a9 47     
            sta $47            ; $eb90: 85 47     
            lda #$06           ; $eb92: a9 06     
            sta $48            ; $eb94: 85 48     
            lda #$eb           ; $eb96: a9 eb     
            sta $0647          ; $eb98: 8d 47 06  
            jsr __fab1         ; $eb9b: 20 b1 fa  
            .hex e3 45         ; $eb9e: e3 45     Invalid Opcode - ISC ($45,x)
            nop                ; $eba0: ea        
            nop                ; $eba1: ea        
            nop                ; $eba2: ea        
            nop                ; $eba3: ea        
            jsr __fab7         ; $eba4: 20 b7 fa  
            lda $0647          ; $eba7: ad 47 06  
            cmp #$ec           ; $ebaa: c9 ec     
            beq __ebb0         ; $ebac: f0 02     
            sty $00            ; $ebae: 84 00     
__ebb0:     iny                ; $ebb0: c8        
            lda #$ff           ; $ebb1: a9 ff     
            sta $0647          ; $ebb3: 8d 47 06  
            jsr __fac2         ; $ebb6: 20 c2 fa  
            .hex e3 45         ; $ebb9: e3 45     Invalid Opcode - ISC ($45,x)
            nop                ; $ebbb: ea        
            nop                ; $ebbc: ea        
            nop                ; $ebbd: ea        
            nop                ; $ebbe: ea        
            jsr __fac7         ; $ebbf: 20 c7 fa  
            lda $0647          ; $ebc2: ad 47 06  
            cmp #$00           ; $ebc5: c9 00     
            beq __ebcb         ; $ebc7: f0 02     
            sty $00            ; $ebc9: 84 00     
__ebcb:     iny                ; $ebcb: c8        
            lda #$37           ; $ebcc: a9 37     
            sta $0647          ; $ebce: 8d 47 06  
            jsr __fad4         ; $ebd1: 20 d4 fa  
            .hex e3 45         ; $ebd4: e3 45     Invalid Opcode - ISC ($45,x)
            nop                ; $ebd6: ea        
            nop                ; $ebd7: ea        
            nop                ; $ebd8: ea        
            nop                ; $ebd9: ea        
            jsr __fada         ; $ebda: 20 da fa  
            lda $0647          ; $ebdd: ad 47 06  
            cmp #$38           ; $ebe0: c9 38     
            beq __ebe6         ; $ebe2: f0 02     
            sty $00            ; $ebe4: 84 00     
__ebe6:     iny                ; $ebe6: c8        
            lda #$eb           ; $ebe7: a9 eb     
            sta $47            ; $ebe9: 85 47     
            jsr __fab1         ; $ebeb: 20 b1 fa  
            .hex e7 47         ; $ebee: e7 47     Invalid Opcode - ISC $47
            nop                ; $ebf0: ea        
            nop                ; $ebf1: ea        
            nop                ; $ebf2: ea        
            nop                ; $ebf3: ea        
            jsr __fab7         ; $ebf4: 20 b7 fa  
            lda $47            ; $ebf7: a5 47     
            cmp #$ec           ; $ebf9: c9 ec     
            beq __ebff         ; $ebfb: f0 02     
            sty $00            ; $ebfd: 84 00     
__ebff:     iny                ; $ebff: c8        
            lda #$ff           ; $ec00: a9 ff     
            sta $47            ; $ec02: 85 47     
            jsr __fac2         ; $ec04: 20 c2 fa  
            .hex e7 47         ; $ec07: e7 47     Invalid Opcode - ISC $47
            nop                ; $ec09: ea        
            nop                ; $ec0a: ea        
            nop                ; $ec0b: ea        
            nop                ; $ec0c: ea        
            jsr __fac7         ; $ec0d: 20 c7 fa  
            lda $47            ; $ec10: a5 47     
            cmp #$00           ; $ec12: c9 00     
            beq __ec18         ; $ec14: f0 02     
            sty $00            ; $ec16: 84 00     
__ec18:     iny                ; $ec18: c8        
            lda #$37           ; $ec19: a9 37     
            sta $47            ; $ec1b: 85 47     
            jsr __fad4         ; $ec1d: 20 d4 fa  
            .hex e7 47         ; $ec20: e7 47     Invalid Opcode - ISC $47
            nop                ; $ec22: ea        
            nop                ; $ec23: ea        
            nop                ; $ec24: ea        
            nop                ; $ec25: ea        
            jsr __fada         ; $ec26: 20 da fa  
            lda $47            ; $ec29: a5 47     
            cmp #$38           ; $ec2b: c9 38     
            beq __ec31         ; $ec2d: f0 02     
            sty $00            ; $ec2f: 84 00     
__ec31:     iny                ; $ec31: c8        
            lda #$eb           ; $ec32: a9 eb     
            sta $0647          ; $ec34: 8d 47 06  
            jsr __fab1         ; $ec37: 20 b1 fa  
            .hex ef 47 06      ; $ec3a: ef 47 06  Invalid Opcode - ISC $0647
            nop                ; $ec3d: ea        
            nop                ; $ec3e: ea        
            nop                ; $ec3f: ea        
            nop                ; $ec40: ea        
            jsr __fab7         ; $ec41: 20 b7 fa  
            lda $0647          ; $ec44: ad 47 06  
            cmp #$ec           ; $ec47: c9 ec     
            beq __ec4d         ; $ec49: f0 02     
            sty $00            ; $ec4b: 84 00     
__ec4d:     iny                ; $ec4d: c8        
            lda #$ff           ; $ec4e: a9 ff     
            sta $0647          ; $ec50: 8d 47 06  
            jsr __fac2         ; $ec53: 20 c2 fa  
            .hex ef 47 06      ; $ec56: ef 47 06  Invalid Opcode - ISC $0647
            nop                ; $ec59: ea        
            nop                ; $ec5a: ea        
            nop                ; $ec5b: ea        
            nop                ; $ec5c: ea        
            jsr __fac7         ; $ec5d: 20 c7 fa  
            lda $0647          ; $ec60: ad 47 06  
            cmp #$00           ; $ec63: c9 00     
            beq __ec69         ; $ec65: f0 02     
            sty $00            ; $ec67: 84 00     
__ec69:     iny                ; $ec69: c8        
            lda #$37           ; $ec6a: a9 37     
            sta $0647          ; $ec6c: 8d 47 06  
            jsr __fad4         ; $ec6f: 20 d4 fa  
            .hex ef 47 06      ; $ec72: ef 47 06  Invalid Opcode - ISC $0647
            nop                ; $ec75: ea        
            nop                ; $ec76: ea        
            nop                ; $ec77: ea        
            nop                ; $ec78: ea        
            jsr __fada         ; $ec79: 20 da fa  
            lda $0647          ; $ec7c: ad 47 06  
            cmp #$38           ; $ec7f: c9 38     
            beq __ec85         ; $ec81: f0 02     
            sty $00            ; $ec83: 84 00     
__ec85:     lda #$eb           ; $ec85: a9 eb     
            sta $0647          ; $ec87: 8d 47 06  
            lda #$48           ; $ec8a: a9 48     
            sta $45            ; $ec8c: 85 45     
            lda #$05           ; $ec8e: a9 05     
            sta $46            ; $ec90: 85 46     
            ldy #$ff           ; $ec92: a0 ff     
            jsr __fab1         ; $ec94: 20 b1 fa  
            .hex f3 45         ; $ec97: f3 45     Invalid Opcode - ISC ($45),y
            nop                ; $ec99: ea        
            nop                ; $ec9a: ea        
            php                ; $ec9b: 08        
            pha                ; $ec9c: 48        
            ldy #$b3           ; $ec9d: a0 b3     
            pla                ; $ec9f: 68        
            plp                ; $eca0: 28        
            jsr __fab7         ; $eca1: 20 b7 fa  
            lda $0647          ; $eca4: ad 47 06  
            cmp #$ec           ; $eca7: c9 ec     
            beq __ecad         ; $eca9: f0 02     
            sty $00            ; $ecab: 84 00     
__ecad:     ldy #$ff           ; $ecad: a0 ff     
            lda #$ff           ; $ecaf: a9 ff     
            sta $0647          ; $ecb1: 8d 47 06  
            jsr __fac2         ; $ecb4: 20 c2 fa  
            .hex f3 45         ; $ecb7: f3 45     Invalid Opcode - ISC ($45),y
            nop                ; $ecb9: ea        
            nop                ; $ecba: ea        
            php                ; $ecbb: 08        
            pha                ; $ecbc: 48        
            ldy #$b4           ; $ecbd: a0 b4     
            pla                ; $ecbf: 68        
            plp                ; $ecc0: 28        
            jsr __fac7         ; $ecc1: 20 c7 fa  
            lda $0647          ; $ecc4: ad 47 06  
            cmp #$00           ; $ecc7: c9 00     
            beq __eccd         ; $ecc9: f0 02     
            sty $00            ; $eccb: 84 00     
__eccd:     ldy #$ff           ; $eccd: a0 ff     
            lda #$37           ; $eccf: a9 37     
            sta $0647          ; $ecd1: 8d 47 06  
            jsr __fad4         ; $ecd4: 20 d4 fa  
            .hex f3 45         ; $ecd7: f3 45     Invalid Opcode - ISC ($45),y
            nop                ; $ecd9: ea        
            nop                ; $ecda: ea        
            php                ; $ecdb: 08        
            pha                ; $ecdc: 48        
            ldy #$b5           ; $ecdd: a0 b5     
            pla                ; $ecdf: 68        
            plp                ; $ece0: 28        
            jsr __fada         ; $ece1: 20 da fa  
            lda $0647          ; $ece4: ad 47 06  
            cmp #$38           ; $ece7: c9 38     
            beq __eced         ; $ece9: f0 02     
            sty $00            ; $eceb: 84 00     
__eced:     ldy #$b6           ; $eced: a0 b6     
            ldx #$ff           ; $ecef: a2 ff     
            lda #$eb           ; $ecf1: a9 eb     
            sta $47            ; $ecf3: 85 47     
            jsr __fab1         ; $ecf5: 20 b1 fa  
            .hex f7 48         ; $ecf8: f7 48     Invalid Opcode - ISC $48,x
            nop                ; $ecfa: ea        
            nop                ; $ecfb: ea        
            nop                ; $ecfc: ea        
            nop                ; $ecfd: ea        
            jsr __fab7         ; $ecfe: 20 b7 fa  
            lda $47            ; $ed01: a5 47     
            cmp #$ec           ; $ed03: c9 ec     
            beq __ed09         ; $ed05: f0 02     
            sty $00            ; $ed07: 84 00     
__ed09:     iny                ; $ed09: c8        
            lda #$ff           ; $ed0a: a9 ff     
            sta $47            ; $ed0c: 85 47     
            jsr __fac2         ; $ed0e: 20 c2 fa  
            .hex f7 48         ; $ed11: f7 48     Invalid Opcode - ISC $48,x
            nop                ; $ed13: ea        
            nop                ; $ed14: ea        
            nop                ; $ed15: ea        
            nop                ; $ed16: ea        
            jsr __fac7         ; $ed17: 20 c7 fa  
            lda $47            ; $ed1a: a5 47     
            cmp #$00           ; $ed1c: c9 00     
            beq __ed22         ; $ed1e: f0 02     
            sty $00            ; $ed20: 84 00     
__ed22:     iny                ; $ed22: c8        
            lda #$37           ; $ed23: a9 37     
            sta $47            ; $ed25: 85 47     
            jsr __fad4         ; $ed27: 20 d4 fa  
            .hex f7 48         ; $ed2a: f7 48     Invalid Opcode - ISC $48,x
            nop                ; $ed2c: ea        
            nop                ; $ed2d: ea        
            nop                ; $ed2e: ea        
            nop                ; $ed2f: ea        
            jsr __fada         ; $ed30: 20 da fa  
            lda $47            ; $ed33: a5 47     
            cmp #$38           ; $ed35: c9 38     
            beq __ed3b         ; $ed37: f0 02     
            sty $00            ; $ed39: 84 00     
__ed3b:     lda #$eb           ; $ed3b: a9 eb     
            sta $0647          ; $ed3d: 8d 47 06  
            ldy #$ff           ; $ed40: a0 ff     
            jsr __fab1         ; $ed42: 20 b1 fa  
            .hex fb 48 05      ; $ed45: fb 48 05  Invalid Opcode - ISC $0548,y
            nop                ; $ed48: ea        
            nop                ; $ed49: ea        
            php                ; $ed4a: 08        
            pha                ; $ed4b: 48        
            ldy #$b9           ; $ed4c: a0 b9     
            pla                ; $ed4e: 68        
            plp                ; $ed4f: 28        
            jsr __fab7         ; $ed50: 20 b7 fa  
            lda $0647          ; $ed53: ad 47 06  
            cmp #$ec           ; $ed56: c9 ec     
            beq __ed5c         ; $ed58: f0 02     
            sty $00            ; $ed5a: 84 00     
__ed5c:     ldy #$ff           ; $ed5c: a0 ff     
            lda #$ff           ; $ed5e: a9 ff     
            sta $0647          ; $ed60: 8d 47 06  
            jsr __fac2         ; $ed63: 20 c2 fa  
            .hex fb 48 05      ; $ed66: fb 48 05  Invalid Opcode - ISC $0548,y
            nop                ; $ed69: ea        
            nop                ; $ed6a: ea        
            php                ; $ed6b: 08        
            pha                ; $ed6c: 48        
            ldy #$ba           ; $ed6d: a0 ba     
            pla                ; $ed6f: 68        
            plp                ; $ed70: 28        
            jsr __fac7         ; $ed71: 20 c7 fa  
            lda $0647          ; $ed74: ad 47 06  
            cmp #$00           ; $ed77: c9 00     
            beq __ed7d         ; $ed79: f0 02     
            sty $00            ; $ed7b: 84 00     
__ed7d:     ldy #$ff           ; $ed7d: a0 ff     
            lda #$37           ; $ed7f: a9 37     
            sta $0647          ; $ed81: 8d 47 06  
            jsr __fad4         ; $ed84: 20 d4 fa  
            .hex fb 48 05      ; $ed87: fb 48 05  Invalid Opcode - ISC $0548,y
            nop                ; $ed8a: ea        
            nop                ; $ed8b: ea        
            php                ; $ed8c: 08        
            pha                ; $ed8d: 48        
            ldy #$bb           ; $ed8e: a0 bb     
            pla                ; $ed90: 68        
            plp                ; $ed91: 28        
            jsr __fada         ; $ed92: 20 da fa  
            lda $0647          ; $ed95: ad 47 06  
            cmp #$38           ; $ed98: c9 38     
            beq __ed9e         ; $ed9a: f0 02     
            sty $00            ; $ed9c: 84 00     
__ed9e:     ldy #$bc           ; $ed9e: a0 bc     
            ldx #$ff           ; $eda0: a2 ff     
            lda #$eb           ; $eda2: a9 eb     
            sta $0647          ; $eda4: 8d 47 06  
            jsr __fab1         ; $eda7: 20 b1 fa  
            .hex ff 48 05      ; $edaa: ff 48 05  Invalid Opcode - ISC $0548,x
            nop                ; $edad: ea        
            nop                ; $edae: ea        
            nop                ; $edaf: ea        
            nop                ; $edb0: ea        
            jsr __fab7         ; $edb1: 20 b7 fa  
            lda $0647          ; $edb4: ad 47 06  
            cmp #$ec           ; $edb7: c9 ec     
            beq __edbd         ; $edb9: f0 02     
            sty $00            ; $edbb: 84 00     
__edbd:     iny                ; $edbd: c8        
            lda #$ff           ; $edbe: a9 ff     
            sta $0647          ; $edc0: 8d 47 06  
            jsr __fac2         ; $edc3: 20 c2 fa  
            .hex ff 48 05      ; $edc6: ff 48 05  Invalid Opcode - ISC $0548,x
            nop                ; $edc9: ea        
            nop                ; $edca: ea        
            nop                ; $edcb: ea        
            nop                ; $edcc: ea        
            jsr __fac7         ; $edcd: 20 c7 fa  
            lda $0647          ; $edd0: ad 47 06  
            cmp #$00           ; $edd3: c9 00     
            beq __edd9         ; $edd5: f0 02     
            sty $00            ; $edd7: 84 00     
__edd9:     iny                ; $edd9: c8        
            lda #$37           ; $edda: a9 37     
            sta $0647          ; $eddc: 8d 47 06  
            jsr __fad4         ; $eddf: 20 d4 fa  
            .hex ff 48 05      ; $ede2: ff 48 05  Invalid Opcode - ISC $0548,x
            nop                ; $ede5: ea        
            nop                ; $ede6: ea        
            nop                ; $ede7: ea        
            nop                ; $ede8: ea        
            jsr __fada         ; $ede9: 20 da fa  
            lda $0647          ; $edec: ad 47 06  
            cmp #$38           ; $edef: c9 38     
            beq __edf5         ; $edf1: f0 02     
            sty $00            ; $edf3: 84 00     
__edf5:     rts                ; $edf5: 60        

;-------------------------------------------------------------------------------
__edf6:     lda #$ff           ; $edf6: a9 ff     
            sta $01            ; $edf8: 85 01     
            ldy #$bf           ; $edfa: a0 bf     
            ldx #$02           ; $edfc: a2 02     
            lda #$47           ; $edfe: a9 47     
            sta $47            ; $ee00: 85 47     
            lda #$06           ; $ee02: a9 06     
            sta $48            ; $ee04: 85 48     
            lda #$a5           ; $ee06: a9 a5     
            sta $0647          ; $ee08: 8d 47 06  
            jsr __fa7b         ; $ee0b: 20 7b fa  
            .hex 03 45         ; $ee0e: 03 45     Invalid Opcode - SLO ($45,x)
            nop                ; $ee10: ea        
            nop                ; $ee11: ea        
            nop                ; $ee12: ea        
            nop                ; $ee13: ea        
            jsr __fa81         ; $ee14: 20 81 fa  
            lda $0647          ; $ee17: ad 47 06  
            cmp #$4a           ; $ee1a: c9 4a     
            beq __ee20         ; $ee1c: f0 02     
            sty $00            ; $ee1e: 84 00     
__ee20:     iny                ; $ee20: c8        
            lda #$29           ; $ee21: a9 29     
            sta $0647          ; $ee23: 8d 47 06  
            jsr __fa8c         ; $ee26: 20 8c fa  
            .hex 03 45         ; $ee29: 03 45     Invalid Opcode - SLO ($45,x)
            nop                ; $ee2b: ea        
            nop                ; $ee2c: ea        
            nop                ; $ee2d: ea        
            nop                ; $ee2e: ea        
            jsr __fa91         ; $ee2f: 20 91 fa  
            lda $0647          ; $ee32: ad 47 06  
            cmp #$52           ; $ee35: c9 52     
            beq __ee3b         ; $ee37: f0 02     
            sty $00            ; $ee39: 84 00     
__ee3b:     iny                ; $ee3b: c8        
            lda #$37           ; $ee3c: a9 37     
            sta $0647          ; $ee3e: 8d 47 06  
            jsr __fa9e         ; $ee41: 20 9e fa  
            .hex 03 45         ; $ee44: 03 45     Invalid Opcode - SLO ($45,x)
            nop                ; $ee46: ea        
            nop                ; $ee47: ea        
            nop                ; $ee48: ea        
            nop                ; $ee49: ea        
            jsr __faa4         ; $ee4a: 20 a4 fa  
            lda $0647          ; $ee4d: ad 47 06  
            cmp #$6e           ; $ee50: c9 6e     
            beq __ee56         ; $ee52: f0 02     
            sty $00            ; $ee54: 84 00     
__ee56:     iny                ; $ee56: c8        
            lda #$a5           ; $ee57: a9 a5     
            sta $47            ; $ee59: 85 47     
            jsr __fa7b         ; $ee5b: 20 7b fa  
            .hex 07 47         ; $ee5e: 07 47     Invalid Opcode - SLO $47
            nop                ; $ee60: ea        
            nop                ; $ee61: ea        
            nop                ; $ee62: ea        
            nop                ; $ee63: ea        
            jsr __fa81         ; $ee64: 20 81 fa  
            lda $47            ; $ee67: a5 47     
            cmp #$4a           ; $ee69: c9 4a     
            beq __ee6f         ; $ee6b: f0 02     
            sty $00            ; $ee6d: 84 00     
__ee6f:     iny                ; $ee6f: c8        
            lda #$29           ; $ee70: a9 29     
            sta $47            ; $ee72: 85 47     
            jsr __fa8c         ; $ee74: 20 8c fa  
            .hex 07 47         ; $ee77: 07 47     Invalid Opcode - SLO $47
            nop                ; $ee79: ea        
            nop                ; $ee7a: ea        
            nop                ; $ee7b: ea        
            nop                ; $ee7c: ea        
            jsr __fa91         ; $ee7d: 20 91 fa  
            lda $47            ; $ee80: a5 47     
            cmp #$52           ; $ee82: c9 52     
            beq __ee88         ; $ee84: f0 02     
            sty $00            ; $ee86: 84 00     
__ee88:     iny                ; $ee88: c8        
            lda #$37           ; $ee89: a9 37     
            sta $47            ; $ee8b: 85 47     
            jsr __fa9e         ; $ee8d: 20 9e fa  
            .hex 07 47         ; $ee90: 07 47     Invalid Opcode - SLO $47
            nop                ; $ee92: ea        
            nop                ; $ee93: ea        
            nop                ; $ee94: ea        
            nop                ; $ee95: ea        
            jsr __faa4         ; $ee96: 20 a4 fa  
            lda $47            ; $ee99: a5 47     
            cmp #$6e           ; $ee9b: c9 6e     
            beq __eea1         ; $ee9d: f0 02     
            sty $00            ; $ee9f: 84 00     
__eea1:     iny                ; $eea1: c8        
            lda #$a5           ; $eea2: a9 a5     
            sta $0647          ; $eea4: 8d 47 06  
            jsr __fa7b         ; $eea7: 20 7b fa  
            .hex 0f 47 06      ; $eeaa: 0f 47 06  Invalid Opcode - SLO $0647
            nop                ; $eead: ea        
            nop                ; $eeae: ea        
            nop                ; $eeaf: ea        
            nop                ; $eeb0: ea        
            jsr __fa81         ; $eeb1: 20 81 fa  
            lda $0647          ; $eeb4: ad 47 06  
            cmp #$4a           ; $eeb7: c9 4a     
            beq __eebd         ; $eeb9: f0 02     
            sty $00            ; $eebb: 84 00     
__eebd:     iny                ; $eebd: c8        
            lda #$29           ; $eebe: a9 29     
            sta $0647          ; $eec0: 8d 47 06  
            jsr __fa8c         ; $eec3: 20 8c fa  
            .hex 0f 47 06      ; $eec6: 0f 47 06  Invalid Opcode - SLO $0647
            nop                ; $eec9: ea        
            nop                ; $eeca: ea        
            nop                ; $eecb: ea        
            nop                ; $eecc: ea        
            jsr __fa91         ; $eecd: 20 91 fa  
            lda $0647          ; $eed0: ad 47 06  
            cmp #$52           ; $eed3: c9 52     
            beq __eed9         ; $eed5: f0 02     
            sty $00            ; $eed7: 84 00     
__eed9:     iny                ; $eed9: c8        
            lda #$37           ; $eeda: a9 37     
            sta $0647          ; $eedc: 8d 47 06  
            jsr __fa9e         ; $eedf: 20 9e fa  
            .hex 0f 47 06      ; $eee2: 0f 47 06  Invalid Opcode - SLO $0647
            nop                ; $eee5: ea        
            nop                ; $eee6: ea        
            nop                ; $eee7: ea        
            nop                ; $eee8: ea        
            jsr __faa4         ; $eee9: 20 a4 fa  
            lda $0647          ; $eeec: ad 47 06  
            cmp #$6e           ; $eeef: c9 6e     
            beq __eef5         ; $eef1: f0 02     
            sty $00            ; $eef3: 84 00     
__eef5:     lda #$a5           ; $eef5: a9 a5     
            sta $0647          ; $eef7: 8d 47 06  
            lda #$48           ; $eefa: a9 48     
            sta $45            ; $eefc: 85 45     
            lda #$05           ; $eefe: a9 05     
            sta $46            ; $ef00: 85 46     
            ldy #$ff           ; $ef02: a0 ff     
            jsr __fa7b         ; $ef04: 20 7b fa  
            .hex 13 45         ; $ef07: 13 45     Invalid Opcode - SLO ($45),y
            nop                ; $ef09: ea        
            nop                ; $ef0a: ea        
            php                ; $ef0b: 08        
            pha                ; $ef0c: 48        
            ldy #$c8           ; $ef0d: a0 c8     
            pla                ; $ef0f: 68        
            plp                ; $ef10: 28        
            jsr __fa81         ; $ef11: 20 81 fa  
            lda $0647          ; $ef14: ad 47 06  
            cmp #$4a           ; $ef17: c9 4a     
            beq __ef1d         ; $ef19: f0 02     
            sty $00            ; $ef1b: 84 00     
__ef1d:     ldy #$ff           ; $ef1d: a0 ff     
            lda #$29           ; $ef1f: a9 29     
            sta $0647          ; $ef21: 8d 47 06  
            jsr __fa8c         ; $ef24: 20 8c fa  
            .hex 13 45         ; $ef27: 13 45     Invalid Opcode - SLO ($45),y
            nop                ; $ef29: ea        
            nop                ; $ef2a: ea        
            php                ; $ef2b: 08        
            pha                ; $ef2c: 48        
            ldy #$c9           ; $ef2d: a0 c9     
            pla                ; $ef2f: 68        
            plp                ; $ef30: 28        
            jsr __fa91         ; $ef31: 20 91 fa  
            lda $0647          ; $ef34: ad 47 06  
            cmp #$52           ; $ef37: c9 52     
            beq __ef3d         ; $ef39: f0 02     
            sty $00            ; $ef3b: 84 00     
__ef3d:     ldy #$ff           ; $ef3d: a0 ff     
            lda #$37           ; $ef3f: a9 37     
            sta $0647          ; $ef41: 8d 47 06  
            jsr __fa9e         ; $ef44: 20 9e fa  
            .hex 13 45         ; $ef47: 13 45     Invalid Opcode - SLO ($45),y
            nop                ; $ef49: ea        
            nop                ; $ef4a: ea        
            php                ; $ef4b: 08        
            pha                ; $ef4c: 48        
            ldy #$ca           ; $ef4d: a0 ca     
            pla                ; $ef4f: 68        
            plp                ; $ef50: 28        
            jsr __faa4         ; $ef51: 20 a4 fa  
            lda $0647          ; $ef54: ad 47 06  
            cmp #$6e           ; $ef57: c9 6e     
            beq __ef5d         ; $ef59: f0 02     
            sty $00            ; $ef5b: 84 00     
__ef5d:     ldy #$cb           ; $ef5d: a0 cb     
            ldx #$ff           ; $ef5f: a2 ff     
            lda #$a5           ; $ef61: a9 a5     
            sta $47            ; $ef63: 85 47     
            jsr __fa7b         ; $ef65: 20 7b fa  
            .hex 17 48         ; $ef68: 17 48     Invalid Opcode - SLO $48,x
            nop                ; $ef6a: ea        
            nop                ; $ef6b: ea        
            nop                ; $ef6c: ea        
            nop                ; $ef6d: ea        
            jsr __fa81         ; $ef6e: 20 81 fa  
            lda $47            ; $ef71: a5 47     
            cmp #$4a           ; $ef73: c9 4a     
            beq __ef79         ; $ef75: f0 02     
            sty $00            ; $ef77: 84 00     
__ef79:     iny                ; $ef79: c8        
            lda #$29           ; $ef7a: a9 29     
            sta $47            ; $ef7c: 85 47     
            jsr __fa8c         ; $ef7e: 20 8c fa  
            .hex 17 48         ; $ef81: 17 48     Invalid Opcode - SLO $48,x
            nop                ; $ef83: ea        
            nop                ; $ef84: ea        
            nop                ; $ef85: ea        
            nop                ; $ef86: ea        
            jsr __fa91         ; $ef87: 20 91 fa  
            lda $47            ; $ef8a: a5 47     
            cmp #$52           ; $ef8c: c9 52     
            beq __ef92         ; $ef8e: f0 02     
            sty $00            ; $ef90: 84 00     
__ef92:     iny                ; $ef92: c8        
            lda #$37           ; $ef93: a9 37     
            sta $47            ; $ef95: 85 47     
            jsr __fa9e         ; $ef97: 20 9e fa  
            .hex 17 48         ; $ef9a: 17 48     Invalid Opcode - SLO $48,x
            nop                ; $ef9c: ea        
            nop                ; $ef9d: ea        
            nop                ; $ef9e: ea        
            nop                ; $ef9f: ea        
            jsr __faa4         ; $efa0: 20 a4 fa  
            lda $47            ; $efa3: a5 47     
            cmp #$6e           ; $efa5: c9 6e     
            beq __efab         ; $efa7: f0 02     
            sty $00            ; $efa9: 84 00     
__efab:     lda #$a5           ; $efab: a9 a5     
            sta $0647          ; $efad: 8d 47 06  
            ldy #$ff           ; $efb0: a0 ff     
            jsr __fa7b         ; $efb2: 20 7b fa  
            .hex 1b 48 05      ; $efb5: 1b 48 05  Invalid Opcode - SLO $0548,y
            nop                ; $efb8: ea        
            nop                ; $efb9: ea        
            php                ; $efba: 08        
            pha                ; $efbb: 48        
            ldy #$ce           ; $efbc: a0 ce     
            pla                ; $efbe: 68        
            plp                ; $efbf: 28        
            jsr __fa81         ; $efc0: 20 81 fa  
            lda $0647          ; $efc3: ad 47 06  
            cmp #$4a           ; $efc6: c9 4a     
            beq __efcc         ; $efc8: f0 02     
            sty $00            ; $efca: 84 00     
__efcc:     ldy #$ff           ; $efcc: a0 ff     
            lda #$29           ; $efce: a9 29     
            sta $0647          ; $efd0: 8d 47 06  
            jsr __fa8c         ; $efd3: 20 8c fa  
            .hex 1b 48 05      ; $efd6: 1b 48 05  Invalid Opcode - SLO $0548,y
            nop                ; $efd9: ea        
            nop                ; $efda: ea        
            php                ; $efdb: 08        
            pha                ; $efdc: 48        
            ldy #$cf           ; $efdd: a0 cf     
            pla                ; $efdf: 68        
            plp                ; $efe0: 28        
            jsr __fa91         ; $efe1: 20 91 fa  
            lda $0647          ; $efe4: ad 47 06  
            cmp #$52           ; $efe7: c9 52     
            beq __efed         ; $efe9: f0 02     
            sty $00            ; $efeb: 84 00     
__efed:     ldy #$ff           ; $efed: a0 ff     
            lda #$37           ; $efef: a9 37     
            sta $0647          ; $eff1: 8d 47 06  
            jsr __fa9e         ; $eff4: 20 9e fa  
            .hex 1b 48 05      ; $eff7: 1b 48 05  Invalid Opcode - SLO $0548,y
            nop                ; $effa: ea        
            nop                ; $effb: ea        
            php                ; $effc: 08        
            pha                ; $effd: 48        
            ldy #$d0           ; $effe: a0 d0     
            pla                ; $f000: 68        
            plp                ; $f001: 28        
            jsr __faa4         ; $f002: 20 a4 fa  
            lda $0647          ; $f005: ad 47 06  
            cmp #$6e           ; $f008: c9 6e     
            beq __f00e         ; $f00a: f0 02     
            sty $00            ; $f00c: 84 00     
__f00e:     ldy #$d1           ; $f00e: a0 d1     
            ldx #$ff           ; $f010: a2 ff     
            lda #$a5           ; $f012: a9 a5     
            sta $0647          ; $f014: 8d 47 06  
            jsr __fa7b         ; $f017: 20 7b fa  
            .hex 1f 48 05      ; $f01a: 1f 48 05  Invalid Opcode - SLO $0548,x
            nop                ; $f01d: ea        
            nop                ; $f01e: ea        
            nop                ; $f01f: ea        
            nop                ; $f020: ea        
            jsr __fa81         ; $f021: 20 81 fa  
            lda $0647          ; $f024: ad 47 06  
            cmp #$4a           ; $f027: c9 4a     
            beq __f02d         ; $f029: f0 02     
            sty $00            ; $f02b: 84 00     
__f02d:     iny                ; $f02d: c8        
            lda #$29           ; $f02e: a9 29     
            sta $0647          ; $f030: 8d 47 06  
            jsr __fa8c         ; $f033: 20 8c fa  
            .hex 1f 48 05      ; $f036: 1f 48 05  Invalid Opcode - SLO $0548,x
            nop                ; $f039: ea        
            nop                ; $f03a: ea        
            nop                ; $f03b: ea        
            nop                ; $f03c: ea        
            jsr __fa91         ; $f03d: 20 91 fa  
            lda $0647          ; $f040: ad 47 06  
            cmp #$52           ; $f043: c9 52     
            beq __f049         ; $f045: f0 02     
            sty $00            ; $f047: 84 00     
__f049:     iny                ; $f049: c8        
            lda #$37           ; $f04a: a9 37     
            sta $0647          ; $f04c: 8d 47 06  
            jsr __fa9e         ; $f04f: 20 9e fa  
            .hex 1f 48 05      ; $f052: 1f 48 05  Invalid Opcode - SLO $0548,x
            nop                ; $f055: ea        
            nop                ; $f056: ea        
            nop                ; $f057: ea        
            nop                ; $f058: ea        
            jsr __faa4         ; $f059: 20 a4 fa  
            lda $0647          ; $f05c: ad 47 06  
            cmp #$6e           ; $f05f: c9 6e     
            beq __f065         ; $f061: f0 02     
            sty $00            ; $f063: 84 00     
__f065:     rts                ; $f065: 60        

;-------------------------------------------------------------------------------
__f066:     lda #$ff           ; $f066: a9 ff     
            sta $01            ; $f068: 85 01     
            ldy #$d4           ; $f06a: a0 d4     
            ldx #$02           ; $f06c: a2 02     
            lda #$47           ; $f06e: a9 47     
            sta $47            ; $f070: 85 47     
            lda #$06           ; $f072: a9 06     
            sta $48            ; $f074: 85 48     
            lda #$a5           ; $f076: a9 a5     
            sta $0647          ; $f078: 8d 47 06  
            jsr __fb53         ; $f07b: 20 53 fb  
            .hex 23 45         ; $f07e: 23 45     Invalid Opcode - RLA ($45,x)
            nop                ; $f080: ea        
            nop                ; $f081: ea        
            nop                ; $f082: ea        
            nop                ; $f083: ea        
            jsr __fb59         ; $f084: 20 59 fb  
            lda $0647          ; $f087: ad 47 06  
            cmp #$4a           ; $f08a: c9 4a     
            beq __f090         ; $f08c: f0 02     
            sty $00            ; $f08e: 84 00     
__f090:     iny                ; $f090: c8        
            lda #$29           ; $f091: a9 29     
            sta $0647          ; $f093: 8d 47 06  
            jsr __fb64         ; $f096: 20 64 fb  
            .hex 23 45         ; $f099: 23 45     Invalid Opcode - RLA ($45,x)
            nop                ; $f09b: ea        
            nop                ; $f09c: ea        
            nop                ; $f09d: ea        
            nop                ; $f09e: ea        
            jsr __fb69         ; $f09f: 20 69 fb  
            lda $0647          ; $f0a2: ad 47 06  
            cmp #$52           ; $f0a5: c9 52     
            beq __f0ab         ; $f0a7: f0 02     
            sty $00            ; $f0a9: 84 00     
__f0ab:     iny                ; $f0ab: c8        
            lda #$37           ; $f0ac: a9 37     
            sta $0647          ; $f0ae: 8d 47 06  
            jsr __fa68         ; $f0b1: 20 68 fa  
            .hex 23 45         ; $f0b4: 23 45     Invalid Opcode - RLA ($45,x)
            nop                ; $f0b6: ea        
            nop                ; $f0b7: ea        
            nop                ; $f0b8: ea        
            nop                ; $f0b9: ea        
            jsr __fa6e         ; $f0ba: 20 6e fa  
            lda $0647          ; $f0bd: ad 47 06  
            cmp #$6f           ; $f0c0: c9 6f     
            beq __f0c6         ; $f0c2: f0 02     
            sty $00            ; $f0c4: 84 00     
__f0c6:     iny                ; $f0c6: c8        
            lda #$a5           ; $f0c7: a9 a5     
            sta $47            ; $f0c9: 85 47     
            jsr __fb53         ; $f0cb: 20 53 fb  
            .hex 27 47         ; $f0ce: 27 47     Invalid Opcode - RLA $47
            nop                ; $f0d0: ea        
            nop                ; $f0d1: ea        
            nop                ; $f0d2: ea        
            nop                ; $f0d3: ea        
            jsr __fb59         ; $f0d4: 20 59 fb  
            lda $47            ; $f0d7: a5 47     
            cmp #$4a           ; $f0d9: c9 4a     
            beq __f0df         ; $f0db: f0 02     
            sty $00            ; $f0dd: 84 00     
__f0df:     iny                ; $f0df: c8        
            lda #$29           ; $f0e0: a9 29     
            sta $47            ; $f0e2: 85 47     
            jsr __fb64         ; $f0e4: 20 64 fb  
            .hex 27 47         ; $f0e7: 27 47     Invalid Opcode - RLA $47
            nop                ; $f0e9: ea        
            nop                ; $f0ea: ea        
            nop                ; $f0eb: ea        
            nop                ; $f0ec: ea        
            jsr __fb69         ; $f0ed: 20 69 fb  
            lda $47            ; $f0f0: a5 47     
            cmp #$52           ; $f0f2: c9 52     
            beq __f0f8         ; $f0f4: f0 02     
            sty $00            ; $f0f6: 84 00     
__f0f8:     iny                ; $f0f8: c8        
            lda #$37           ; $f0f9: a9 37     
            sta $47            ; $f0fb: 85 47     
            jsr __fa68         ; $f0fd: 20 68 fa  
            .hex 27 47         ; $f100: 27 47     Invalid Opcode - RLA $47
            nop                ; $f102: ea        
            nop                ; $f103: ea        
            nop                ; $f104: ea        
            nop                ; $f105: ea        
            jsr __fa6e         ; $f106: 20 6e fa  
            lda $47            ; $f109: a5 47     
            cmp #$6f           ; $f10b: c9 6f     
            beq __f111         ; $f10d: f0 02     
            sty $00            ; $f10f: 84 00     
__f111:     iny                ; $f111: c8        
            lda #$a5           ; $f112: a9 a5     
            sta $0647          ; $f114: 8d 47 06  
            jsr __fb53         ; $f117: 20 53 fb  
            .hex 2f 47 06      ; $f11a: 2f 47 06  Invalid Opcode - RLA $0647
            nop                ; $f11d: ea        
            nop                ; $f11e: ea        
            nop                ; $f11f: ea        
            nop                ; $f120: ea        
            jsr __fb59         ; $f121: 20 59 fb  
            lda $0647          ; $f124: ad 47 06  
            cmp #$4a           ; $f127: c9 4a     
            beq __f12d         ; $f129: f0 02     
            sty $00            ; $f12b: 84 00     
__f12d:     iny                ; $f12d: c8        
            lda #$29           ; $f12e: a9 29     
            sta $0647          ; $f130: 8d 47 06  
            jsr __fb64         ; $f133: 20 64 fb  
            .hex 2f 47 06      ; $f136: 2f 47 06  Invalid Opcode - RLA $0647
            nop                ; $f139: ea        
            nop                ; $f13a: ea        
            nop                ; $f13b: ea        
            nop                ; $f13c: ea        
            jsr __fb69         ; $f13d: 20 69 fb  
            lda $0647          ; $f140: ad 47 06  
            cmp #$52           ; $f143: c9 52     
            beq __f149         ; $f145: f0 02     
            sty $00            ; $f147: 84 00     
__f149:     iny                ; $f149: c8        
            lda #$37           ; $f14a: a9 37     
            sta $0647          ; $f14c: 8d 47 06  
            jsr __fa68         ; $f14f: 20 68 fa  
            .hex 2f 47 06      ; $f152: 2f 47 06  Invalid Opcode - RLA $0647
            nop                ; $f155: ea        
            nop                ; $f156: ea        
            nop                ; $f157: ea        
            nop                ; $f158: ea        
            jsr __fa6e         ; $f159: 20 6e fa  
            lda $0647          ; $f15c: ad 47 06  
            cmp #$6f           ; $f15f: c9 6f     
            beq __f165         ; $f161: f0 02     
            sty $00            ; $f163: 84 00     
__f165:     lda #$a5           ; $f165: a9 a5     
            sta $0647          ; $f167: 8d 47 06  
            lda #$48           ; $f16a: a9 48     
            sta $45            ; $f16c: 85 45     
            lda #$05           ; $f16e: a9 05     
            sta $46            ; $f170: 85 46     
            ldy #$ff           ; $f172: a0 ff     
            jsr __fb53         ; $f174: 20 53 fb  
            .hex 33 45         ; $f177: 33 45     Invalid Opcode - RLA ($45),y
            nop                ; $f179: ea        
            nop                ; $f17a: ea        
            php                ; $f17b: 08        
            pha                ; $f17c: 48        
            ldy #$dd           ; $f17d: a0 dd     
            pla                ; $f17f: 68        
            plp                ; $f180: 28        
            jsr __fb59         ; $f181: 20 59 fb  
            lda $0647          ; $f184: ad 47 06  
            cmp #$4a           ; $f187: c9 4a     
            beq __f18d         ; $f189: f0 02     
            sty $00            ; $f18b: 84 00     
__f18d:     ldy #$ff           ; $f18d: a0 ff     
            lda #$29           ; $f18f: a9 29     
            sta $0647          ; $f191: 8d 47 06  
            jsr __fb64         ; $f194: 20 64 fb  
            .hex 33 45         ; $f197: 33 45     Invalid Opcode - RLA ($45),y
            nop                ; $f199: ea        
            nop                ; $f19a: ea        
            php                ; $f19b: 08        
            pha                ; $f19c: 48        
            ldy #$de           ; $f19d: a0 de     
            pla                ; $f19f: 68        
            plp                ; $f1a0: 28        
            jsr __fb69         ; $f1a1: 20 69 fb  
            lda $0647          ; $f1a4: ad 47 06  
            cmp #$52           ; $f1a7: c9 52     
            beq __f1ad         ; $f1a9: f0 02     
            sty $00            ; $f1ab: 84 00     
__f1ad:     ldy #$ff           ; $f1ad: a0 ff     
            lda #$37           ; $f1af: a9 37     
            sta $0647          ; $f1b1: 8d 47 06  
            jsr __fa68         ; $f1b4: 20 68 fa  
            .hex 33 45         ; $f1b7: 33 45     Invalid Opcode - RLA ($45),y
            nop                ; $f1b9: ea        
            nop                ; $f1ba: ea        
            php                ; $f1bb: 08        
            pha                ; $f1bc: 48        
            ldy #$df           ; $f1bd: a0 df     
            pla                ; $f1bf: 68        
            plp                ; $f1c0: 28        
            jsr __fa6e         ; $f1c1: 20 6e fa  
            lda $0647          ; $f1c4: ad 47 06  
            cmp #$6f           ; $f1c7: c9 6f     
            beq __f1cd         ; $f1c9: f0 02     
            sty $00            ; $f1cb: 84 00     
__f1cd:     ldy #$e0           ; $f1cd: a0 e0     
            ldx #$ff           ; $f1cf: a2 ff     
            lda #$a5           ; $f1d1: a9 a5     
            sta $47            ; $f1d3: 85 47     
            jsr __fb53         ; $f1d5: 20 53 fb  
            .hex 37 48         ; $f1d8: 37 48     Invalid Opcode - RLA $48,x
            nop                ; $f1da: ea        
            nop                ; $f1db: ea        
            nop                ; $f1dc: ea        
            nop                ; $f1dd: ea        
            jsr __fb59         ; $f1de: 20 59 fb  
            lda $47            ; $f1e1: a5 47     
            cmp #$4a           ; $f1e3: c9 4a     
            beq __f1e9         ; $f1e5: f0 02     
            sty $00            ; $f1e7: 84 00     
__f1e9:     iny                ; $f1e9: c8        
            lda #$29           ; $f1ea: a9 29     
            sta $47            ; $f1ec: 85 47     
            jsr __fb64         ; $f1ee: 20 64 fb  
            .hex 37 48         ; $f1f1: 37 48     Invalid Opcode - RLA $48,x
            nop                ; $f1f3: ea        
            nop                ; $f1f4: ea        
            nop                ; $f1f5: ea        
            nop                ; $f1f6: ea        
            jsr __fb69         ; $f1f7: 20 69 fb  
            lda $47            ; $f1fa: a5 47     
            cmp #$52           ; $f1fc: c9 52     
            beq __f202         ; $f1fe: f0 02     
            sty $00            ; $f200: 84 00     
__f202:     iny                ; $f202: c8        
            lda #$37           ; $f203: a9 37     
            sta $47            ; $f205: 85 47     
            jsr __fa68         ; $f207: 20 68 fa  
            .hex 37 48         ; $f20a: 37 48     Invalid Opcode - RLA $48,x
            nop                ; $f20c: ea        
            nop                ; $f20d: ea        
            nop                ; $f20e: ea        
            nop                ; $f20f: ea        
            jsr __fa6e         ; $f210: 20 6e fa  
            lda $47            ; $f213: a5 47     
            cmp #$6f           ; $f215: c9 6f     
            beq __f21b         ; $f217: f0 02     
            sty $00            ; $f219: 84 00     
__f21b:     lda #$a5           ; $f21b: a9 a5     
            sta $0647          ; $f21d: 8d 47 06  
            ldy #$ff           ; $f220: a0 ff     
            jsr __fb53         ; $f222: 20 53 fb  
            .hex 3b 48 05      ; $f225: 3b 48 05  Invalid Opcode - RLA $0548,y
            nop                ; $f228: ea        
            nop                ; $f229: ea        
            php                ; $f22a: 08        
            pha                ; $f22b: 48        
            ldy #$e3           ; $f22c: a0 e3     
            pla                ; $f22e: 68        
            plp                ; $f22f: 28        
            jsr __fb59         ; $f230: 20 59 fb  
            lda $0647          ; $f233: ad 47 06  
            cmp #$4a           ; $f236: c9 4a     
            beq __f23c         ; $f238: f0 02     
            sty $00            ; $f23a: 84 00     
__f23c:     ldy #$ff           ; $f23c: a0 ff     
            lda #$29           ; $f23e: a9 29     
            sta $0647          ; $f240: 8d 47 06  
            jsr __fb64         ; $f243: 20 64 fb  
            .hex 3b 48 05      ; $f246: 3b 48 05  Invalid Opcode - RLA $0548,y
            nop                ; $f249: ea        
            nop                ; $f24a: ea        
            php                ; $f24b: 08        
            pha                ; $f24c: 48        
            ldy #$e4           ; $f24d: a0 e4     
            pla                ; $f24f: 68        
            plp                ; $f250: 28        
            jsr __fb69         ; $f251: 20 69 fb  
            lda $0647          ; $f254: ad 47 06  
            cmp #$52           ; $f257: c9 52     
            beq __f25d         ; $f259: f0 02     
            sty $00            ; $f25b: 84 00     
__f25d:     ldy #$ff           ; $f25d: a0 ff     
            lda #$37           ; $f25f: a9 37     
            sta $0647          ; $f261: 8d 47 06  
            jsr __fa68         ; $f264: 20 68 fa  
            .hex 3b 48 05      ; $f267: 3b 48 05  Invalid Opcode - RLA $0548,y
            nop                ; $f26a: ea        
            nop                ; $f26b: ea        
            php                ; $f26c: 08        
            pha                ; $f26d: 48        
            ldy #$e5           ; $f26e: a0 e5     
            pla                ; $f270: 68        
            plp                ; $f271: 28        
            jsr __fa6e         ; $f272: 20 6e fa  
            lda $0647          ; $f275: ad 47 06  
            cmp #$6f           ; $f278: c9 6f     
            beq __f27e         ; $f27a: f0 02     
            sty $00            ; $f27c: 84 00     
__f27e:     ldy #$e6           ; $f27e: a0 e6     
            ldx #$ff           ; $f280: a2 ff     
            lda #$a5           ; $f282: a9 a5     
            sta $0647          ; $f284: 8d 47 06  
            jsr __fb53         ; $f287: 20 53 fb  
            .hex 3f 48 05      ; $f28a: 3f 48 05  Invalid Opcode - RLA $0548,x
            nop                ; $f28d: ea        
            nop                ; $f28e: ea        
            nop                ; $f28f: ea        
            nop                ; $f290: ea        
            jsr __fb59         ; $f291: 20 59 fb  
            lda $0647          ; $f294: ad 47 06  
            cmp #$4a           ; $f297: c9 4a     
            beq __f29d         ; $f299: f0 02     
            sty $00            ; $f29b: 84 00     
__f29d:     iny                ; $f29d: c8        
            lda #$29           ; $f29e: a9 29     
            sta $0647          ; $f2a0: 8d 47 06  
            jsr __fb64         ; $f2a3: 20 64 fb  
            .hex 3f 48 05      ; $f2a6: 3f 48 05  Invalid Opcode - RLA $0548,x
            nop                ; $f2a9: ea        
            nop                ; $f2aa: ea        
            nop                ; $f2ab: ea        
            nop                ; $f2ac: ea        
            jsr __fb69         ; $f2ad: 20 69 fb  
            lda $0647          ; $f2b0: ad 47 06  
            cmp #$52           ; $f2b3: c9 52     
            beq __f2b9         ; $f2b5: f0 02     
            sty $00            ; $f2b7: 84 00     
__f2b9:     iny                ; $f2b9: c8        
            lda #$37           ; $f2ba: a9 37     
            sta $0647          ; $f2bc: 8d 47 06  
            jsr __fa68         ; $f2bf: 20 68 fa  
            .hex 3f 48 05      ; $f2c2: 3f 48 05  Invalid Opcode - RLA $0548,x
            nop                ; $f2c5: ea        
            nop                ; $f2c6: ea        
            nop                ; $f2c7: ea        
            nop                ; $f2c8: ea        
            jsr __fa6e         ; $f2c9: 20 6e fa  
            lda $0647          ; $f2cc: ad 47 06  
            cmp #$6f           ; $f2cf: c9 6f     
            beq __f2d5         ; $f2d1: f0 02     
            sty $00            ; $f2d3: 84 00     
__f2d5:     rts                ; $f2d5: 60        

;-------------------------------------------------------------------------------
__f2d6:     lda #$ff           ; $f2d6: a9 ff     
            sta $01            ; $f2d8: 85 01     
            ldy #$e9           ; $f2da: a0 e9     
            ldx #$02           ; $f2dc: a2 02     
            lda #$47           ; $f2de: a9 47     
            sta $47            ; $f2e0: 85 47     
            lda #$06           ; $f2e2: a9 06     
            sta $48            ; $f2e4: 85 48     
            lda #$a5           ; $f2e6: a9 a5     
            sta $0647          ; $f2e8: 8d 47 06  
            jsr __fb1d         ; $f2eb: 20 1d fb  
            .hex 43 45         ; $f2ee: 43 45     Invalid Opcode - SRE ($45,x)
            nop                ; $f2f0: ea        
            nop                ; $f2f1: ea        
            nop                ; $f2f2: ea        
            nop                ; $f2f3: ea        
            jsr __fb23         ; $f2f4: 20 23 fb  
            lda $0647          ; $f2f7: ad 47 06  
            cmp #$52           ; $f2fa: c9 52     
            beq __f300         ; $f2fc: f0 02     
            sty $00            ; $f2fe: 84 00     
__f300:     iny                ; $f300: c8        
            lda #$29           ; $f301: a9 29     
            sta $0647          ; $f303: 8d 47 06  
            jsr __fb2e         ; $f306: 20 2e fb  
            .hex 43 45         ; $f309: 43 45     Invalid Opcode - SRE ($45,x)
            nop                ; $f30b: ea        
            nop                ; $f30c: ea        
            nop                ; $f30d: ea        
            nop                ; $f30e: ea        
            jsr __fb33         ; $f30f: 20 33 fb  
            lda $0647          ; $f312: ad 47 06  
            cmp #$14           ; $f315: c9 14     
            beq __f31b         ; $f317: f0 02     
            sty $00            ; $f319: 84 00     
__f31b:     iny                ; $f31b: c8        
            lda #$37           ; $f31c: a9 37     
            sta $0647          ; $f31e: 8d 47 06  
            jsr __fb40         ; $f321: 20 40 fb  
            .hex 43 45         ; $f324: 43 45     Invalid Opcode - SRE ($45,x)
            nop                ; $f326: ea        
            nop                ; $f327: ea        
            nop                ; $f328: ea        
            nop                ; $f329: ea        
            jsr __fb46         ; $f32a: 20 46 fb  
            lda $0647          ; $f32d: ad 47 06  
            cmp #$1b           ; $f330: c9 1b     
            beq __f336         ; $f332: f0 02     
            sty $00            ; $f334: 84 00     
__f336:     iny                ; $f336: c8        
            lda #$a5           ; $f337: a9 a5     
            sta $47            ; $f339: 85 47     
            jsr __fb1d         ; $f33b: 20 1d fb  
            .hex 47 47         ; $f33e: 47 47     Invalid Opcode - SRE $47
            nop                ; $f340: ea        
            nop                ; $f341: ea        
            nop                ; $f342: ea        
            nop                ; $f343: ea        
            jsr __fb23         ; $f344: 20 23 fb  
            lda $47            ; $f347: a5 47     
            cmp #$52           ; $f349: c9 52     
            beq __f34f         ; $f34b: f0 02     
            sty $00            ; $f34d: 84 00     
__f34f:     iny                ; $f34f: c8        
            lda #$29           ; $f350: a9 29     
            sta $47            ; $f352: 85 47     
            jsr __fb2e         ; $f354: 20 2e fb  
            .hex 47 47         ; $f357: 47 47     Invalid Opcode - SRE $47
            nop                ; $f359: ea        
            nop                ; $f35a: ea        
            nop                ; $f35b: ea        
            nop                ; $f35c: ea        
            jsr __fb33         ; $f35d: 20 33 fb  
            lda $47            ; $f360: a5 47     
            cmp #$14           ; $f362: c9 14     
            beq __f368         ; $f364: f0 02     
            sty $00            ; $f366: 84 00     
__f368:     iny                ; $f368: c8        
            lda #$37           ; $f369: a9 37     
            sta $47            ; $f36b: 85 47     
            jsr __fb40         ; $f36d: 20 40 fb  
            .hex 47 47         ; $f370: 47 47     Invalid Opcode - SRE $47
            nop                ; $f372: ea        
            nop                ; $f373: ea        
            nop                ; $f374: ea        
            nop                ; $f375: ea        
            jsr __fb46         ; $f376: 20 46 fb  
            lda $47            ; $f379: a5 47     
            cmp #$1b           ; $f37b: c9 1b     
            beq __f381         ; $f37d: f0 02     
            sty $00            ; $f37f: 84 00     
__f381:     iny                ; $f381: c8        
            lda #$a5           ; $f382: a9 a5     
            sta $0647          ; $f384: 8d 47 06  
            jsr __fb1d         ; $f387: 20 1d fb  
            .hex 4f 47 06      ; $f38a: 4f 47 06  Invalid Opcode - SRE $0647
            nop                ; $f38d: ea        
            nop                ; $f38e: ea        
            nop                ; $f38f: ea        
            nop                ; $f390: ea        
            jsr __fb23         ; $f391: 20 23 fb  
            lda $0647          ; $f394: ad 47 06  
            cmp #$52           ; $f397: c9 52     
            beq __f39d         ; $f399: f0 02     
            sty $00            ; $f39b: 84 00     
__f39d:     iny                ; $f39d: c8        
            lda #$29           ; $f39e: a9 29     
            sta $0647          ; $f3a0: 8d 47 06  
            jsr __fb2e         ; $f3a3: 20 2e fb  
            .hex 4f 47 06      ; $f3a6: 4f 47 06  Invalid Opcode - SRE $0647
            nop                ; $f3a9: ea        
            nop                ; $f3aa: ea        
            nop                ; $f3ab: ea        
            nop                ; $f3ac: ea        
            jsr __fb33         ; $f3ad: 20 33 fb  
            lda $0647          ; $f3b0: ad 47 06  
            cmp #$14           ; $f3b3: c9 14     
            beq __f3b9         ; $f3b5: f0 02     
            sty $00            ; $f3b7: 84 00     
__f3b9:     iny                ; $f3b9: c8        
            lda #$37           ; $f3ba: a9 37     
            sta $0647          ; $f3bc: 8d 47 06  
            jsr __fb40         ; $f3bf: 20 40 fb  
            .hex 4f 47 06      ; $f3c2: 4f 47 06  Invalid Opcode - SRE $0647
            nop                ; $f3c5: ea        
            nop                ; $f3c6: ea        
            nop                ; $f3c7: ea        
            nop                ; $f3c8: ea        
            jsr __fb46         ; $f3c9: 20 46 fb  
            lda $0647          ; $f3cc: ad 47 06  
            cmp #$1b           ; $f3cf: c9 1b     
            beq __f3d5         ; $f3d1: f0 02     
            sty $00            ; $f3d3: 84 00     
__f3d5:     lda #$a5           ; $f3d5: a9 a5     
            sta $0647          ; $f3d7: 8d 47 06  
            lda #$48           ; $f3da: a9 48     
            sta $45            ; $f3dc: 85 45     
            lda #$05           ; $f3de: a9 05     
            sta $46            ; $f3e0: 85 46     
            ldy #$ff           ; $f3e2: a0 ff     
            jsr __fb1d         ; $f3e4: 20 1d fb  
            .hex 53 45         ; $f3e7: 53 45     Invalid Opcode - SRE ($45),y
            nop                ; $f3e9: ea        
            nop                ; $f3ea: ea        
            php                ; $f3eb: 08        
            pha                ; $f3ec: 48        
            ldy #$f2           ; $f3ed: a0 f2     
            pla                ; $f3ef: 68        
            plp                ; $f3f0: 28        
            jsr __fb23         ; $f3f1: 20 23 fb  
            lda $0647          ; $f3f4: ad 47 06  
            cmp #$52           ; $f3f7: c9 52     
            beq __f3fd         ; $f3f9: f0 02     
            sty $00            ; $f3fb: 84 00     
__f3fd:     ldy #$ff           ; $f3fd: a0 ff     
            lda #$29           ; $f3ff: a9 29     
            sta $0647          ; $f401: 8d 47 06  
            jsr __fb2e         ; $f404: 20 2e fb  
            .hex 53 45         ; $f407: 53 45     Invalid Opcode - SRE ($45),y
            nop                ; $f409: ea        
            nop                ; $f40a: ea        
            php                ; $f40b: 08        
            pha                ; $f40c: 48        
            ldy #$f3           ; $f40d: a0 f3     
            pla                ; $f40f: 68        
            plp                ; $f410: 28        
            jsr __fb33         ; $f411: 20 33 fb  
            lda $0647          ; $f414: ad 47 06  
            cmp #$14           ; $f417: c9 14     
            beq __f41d         ; $f419: f0 02     
            sty $00            ; $f41b: 84 00     
__f41d:     ldy #$ff           ; $f41d: a0 ff     
            lda #$37           ; $f41f: a9 37     
            sta $0647          ; $f421: 8d 47 06  
            jsr __fb40         ; $f424: 20 40 fb  
            .hex 53 45         ; $f427: 53 45     Invalid Opcode - SRE ($45),y
            nop                ; $f429: ea        
            nop                ; $f42a: ea        
            php                ; $f42b: 08        
            pha                ; $f42c: 48        
            ldy #$f4           ; $f42d: a0 f4     
            pla                ; $f42f: 68        
            plp                ; $f430: 28        
            jsr __fb46         ; $f431: 20 46 fb  
            lda $0647          ; $f434: ad 47 06  
            cmp #$1b           ; $f437: c9 1b     
            beq __f43d         ; $f439: f0 02     
            sty $00            ; $f43b: 84 00     
__f43d:     ldy #$f5           ; $f43d: a0 f5     
            ldx #$ff           ; $f43f: a2 ff     
            lda #$a5           ; $f441: a9 a5     
            sta $47            ; $f443: 85 47     
            jsr __fb1d         ; $f445: 20 1d fb  
            .hex 57 48         ; $f448: 57 48     Invalid Opcode - SRE $48,x
            nop                ; $f44a: ea        
            nop                ; $f44b: ea        
            nop                ; $f44c: ea        
            nop                ; $f44d: ea        
            jsr __fb23         ; $f44e: 20 23 fb  
            lda $47            ; $f451: a5 47     
            cmp #$52           ; $f453: c9 52     
            beq __f459         ; $f455: f0 02     
            sty $00            ; $f457: 84 00     
__f459:     iny                ; $f459: c8        
            lda #$29           ; $f45a: a9 29     
            sta $47            ; $f45c: 85 47     
            jsr __fb2e         ; $f45e: 20 2e fb  
            .hex 57 48         ; $f461: 57 48     Invalid Opcode - SRE $48,x
            nop                ; $f463: ea        
            nop                ; $f464: ea        
            nop                ; $f465: ea        
            nop                ; $f466: ea        
            jsr __fb33         ; $f467: 20 33 fb  
            lda $47            ; $f46a: a5 47     
            cmp #$14           ; $f46c: c9 14     
            beq __f472         ; $f46e: f0 02     
            sty $00            ; $f470: 84 00     
__f472:     iny                ; $f472: c8        
            lda #$37           ; $f473: a9 37     
            sta $47            ; $f475: 85 47     
            jsr __fb40         ; $f477: 20 40 fb  
            .hex 57 48         ; $f47a: 57 48     Invalid Opcode - SRE $48,x
            nop                ; $f47c: ea        
            nop                ; $f47d: ea        
            nop                ; $f47e: ea        
            nop                ; $f47f: ea        
            jsr __fb46         ; $f480: 20 46 fb  
            lda $47            ; $f483: a5 47     
            cmp #$1b           ; $f485: c9 1b     
            beq __f48b         ; $f487: f0 02     
            sty $00            ; $f489: 84 00     
__f48b:     lda #$a5           ; $f48b: a9 a5     
            sta $0647          ; $f48d: 8d 47 06  
            ldy #$ff           ; $f490: a0 ff     
            jsr __fb1d         ; $f492: 20 1d fb  
            .hex 5b 48 05      ; $f495: 5b 48 05  Invalid Opcode - SRE $0548,y
            nop                ; $f498: ea        
            nop                ; $f499: ea        
            php                ; $f49a: 08        
            pha                ; $f49b: 48        
            ldy #$f8           ; $f49c: a0 f8     
            pla                ; $f49e: 68        
            plp                ; $f49f: 28        
            jsr __fb23         ; $f4a0: 20 23 fb  
            lda $0647          ; $f4a3: ad 47 06  
            cmp #$52           ; $f4a6: c9 52     
            beq __f4ac         ; $f4a8: f0 02     
            sty $00            ; $f4aa: 84 00     
__f4ac:     ldy #$ff           ; $f4ac: a0 ff     
            lda #$29           ; $f4ae: a9 29     
            sta $0647          ; $f4b0: 8d 47 06  
            jsr __fb2e         ; $f4b3: 20 2e fb  
            .hex 5b 48 05      ; $f4b6: 5b 48 05  Invalid Opcode - SRE $0548,y
            nop                ; $f4b9: ea        
            nop                ; $f4ba: ea        
            php                ; $f4bb: 08        
            pha                ; $f4bc: 48        
            ldy #$f9           ; $f4bd: a0 f9     
            pla                ; $f4bf: 68        
            plp                ; $f4c0: 28        
            jsr __fb33         ; $f4c1: 20 33 fb  
            lda $0647          ; $f4c4: ad 47 06  
            cmp #$14           ; $f4c7: c9 14     
            beq __f4cd         ; $f4c9: f0 02     
            sty $00            ; $f4cb: 84 00     
__f4cd:     ldy #$ff           ; $f4cd: a0 ff     
            lda #$37           ; $f4cf: a9 37     
            sta $0647          ; $f4d1: 8d 47 06  
            jsr __fb40         ; $f4d4: 20 40 fb  
            .hex 5b 48 05      ; $f4d7: 5b 48 05  Invalid Opcode - SRE $0548,y
            nop                ; $f4da: ea        
            nop                ; $f4db: ea        
            php                ; $f4dc: 08        
            pha                ; $f4dd: 48        
            ldy #$fa           ; $f4de: a0 fa     
            pla                ; $f4e0: 68        
            plp                ; $f4e1: 28        
            jsr __fb46         ; $f4e2: 20 46 fb  
            lda $0647          ; $f4e5: ad 47 06  
            cmp #$1b           ; $f4e8: c9 1b     
            beq __f4ee         ; $f4ea: f0 02     
            sty $00            ; $f4ec: 84 00     
__f4ee:     ldy #$fb           ; $f4ee: a0 fb     
            ldx #$ff           ; $f4f0: a2 ff     
            lda #$a5           ; $f4f2: a9 a5     
            sta $0647          ; $f4f4: 8d 47 06  
            jsr __fb1d         ; $f4f7: 20 1d fb  
            .hex 5f 48 05      ; $f4fa: 5f 48 05  Invalid Opcode - SRE $0548,x
            nop                ; $f4fd: ea        
            nop                ; $f4fe: ea        
            nop                ; $f4ff: ea        
            nop                ; $f500: ea        
            jsr __fb23         ; $f501: 20 23 fb  
            lda $0647          ; $f504: ad 47 06  
            cmp #$52           ; $f507: c9 52     
            beq __f50d         ; $f509: f0 02     
            sty $00            ; $f50b: 84 00     
__f50d:     iny                ; $f50d: c8        
            lda #$29           ; $f50e: a9 29     
            sta $0647          ; $f510: 8d 47 06  
            jsr __fb2e         ; $f513: 20 2e fb  
            .hex 5f 48 05      ; $f516: 5f 48 05  Invalid Opcode - SRE $0548,x
            nop                ; $f519: ea        
            nop                ; $f51a: ea        
            nop                ; $f51b: ea        
            nop                ; $f51c: ea        
            jsr __fb33         ; $f51d: 20 33 fb  
            lda $0647          ; $f520: ad 47 06  
            cmp #$14           ; $f523: c9 14     
            beq __f529         ; $f525: f0 02     
            sty $00            ; $f527: 84 00     
__f529:     iny                ; $f529: c8        
            lda #$37           ; $f52a: a9 37     
            sta $0647          ; $f52c: 8d 47 06  
            jsr __fb40         ; $f52f: 20 40 fb  
            .hex 5f 48 05      ; $f532: 5f 48 05  Invalid Opcode - SRE $0548,x
            nop                ; $f535: ea        
            nop                ; $f536: ea        
            nop                ; $f537: ea        
            nop                ; $f538: ea        
            jsr __fb46         ; $f539: 20 46 fb  
            lda $0647          ; $f53c: ad 47 06  
            cmp #$1b           ; $f53f: c9 1b     
            beq __f545         ; $f541: f0 02     
            sty $00            ; $f543: 84 00     
__f545:     rts                ; $f545: 60        

;-------------------------------------------------------------------------------
__f546:     lda #$ff           ; $f546: a9 ff     
            sta $01            ; $f548: 85 01     
            ldy #$01           ; $f54a: a0 01     
            ldx #$02           ; $f54c: a2 02     
            lda #$47           ; $f54e: a9 47     
            sta $47            ; $f550: 85 47     
            lda #$06           ; $f552: a9 06     
            sta $48            ; $f554: 85 48     
            lda #$a5           ; $f556: a9 a5     
            sta $0647          ; $f558: 8d 47 06  
            jsr __fae9         ; $f55b: 20 e9 fa  
            .hex 63 45         ; $f55e: 63 45     Invalid Opcode - RRA ($45,x)
            nop                ; $f560: ea        
            nop                ; $f561: ea        
            nop                ; $f562: ea        
            nop                ; $f563: ea        
            jsr __faef         ; $f564: 20 ef fa  
            lda $0647          ; $f567: ad 47 06  
            cmp #$52           ; $f56a: c9 52     
            beq __f570         ; $f56c: f0 02     
            sty $00            ; $f56e: 84 00     
__f570:     iny                ; $f570: c8        
            lda #$29           ; $f571: a9 29     
            sta $0647          ; $f573: 8d 47 06  
            jsr __fafa         ; $f576: 20 fa fa  
            .hex 63 45         ; $f579: 63 45     Invalid Opcode - RRA ($45,x)
            nop                ; $f57b: ea        
            nop                ; $f57c: ea        
            nop                ; $f57d: ea        
            nop                ; $f57e: ea        
            jsr __faff         ; $f57f: 20 ff fa  
            lda $0647          ; $f582: ad 47 06  
            cmp #$14           ; $f585: c9 14     
            beq __f58b         ; $f587: f0 02     
            sty $00            ; $f589: 84 00     
__f58b:     iny                ; $f58b: c8        
            lda #$37           ; $f58c: a9 37     
            sta $0647          ; $f58e: 8d 47 06  
            jsr __fb0a         ; $f591: 20 0a fb  
            .hex 63 45         ; $f594: 63 45     Invalid Opcode - RRA ($45,x)
            nop                ; $f596: ea        
            nop                ; $f597: ea        
            nop                ; $f598: ea        
            nop                ; $f599: ea        
            jsr __fb10         ; $f59a: 20 10 fb  
            lda $0647          ; $f59d: ad 47 06  
            cmp #$9b           ; $f5a0: c9 9b     
            beq __f5a6         ; $f5a2: f0 02     
            sty $00            ; $f5a4: 84 00     
__f5a6:     iny                ; $f5a6: c8        
            lda #$a5           ; $f5a7: a9 a5     
            sta $47            ; $f5a9: 85 47     
            jsr __fae9         ; $f5ab: 20 e9 fa  
            .hex 67 47         ; $f5ae: 67 47     Invalid Opcode - RRA $47
            nop                ; $f5b0: ea        
            nop                ; $f5b1: ea        
            nop                ; $f5b2: ea        
            nop                ; $f5b3: ea        
            jsr __faef         ; $f5b4: 20 ef fa  
            lda $47            ; $f5b7: a5 47     
            cmp #$52           ; $f5b9: c9 52     
            beq __f5bf         ; $f5bb: f0 02     
            sty $00            ; $f5bd: 84 00     
__f5bf:     iny                ; $f5bf: c8        
            lda #$29           ; $f5c0: a9 29     
            sta $47            ; $f5c2: 85 47     
            jsr __fafa         ; $f5c4: 20 fa fa  
            .hex 67 47         ; $f5c7: 67 47     Invalid Opcode - RRA $47
            nop                ; $f5c9: ea        
            nop                ; $f5ca: ea        
            nop                ; $f5cb: ea        
            nop                ; $f5cc: ea        
            jsr __faff         ; $f5cd: 20 ff fa  
            lda $47            ; $f5d0: a5 47     
            cmp #$14           ; $f5d2: c9 14     
            beq __f5d8         ; $f5d4: f0 02     
            sty $00            ; $f5d6: 84 00     
__f5d8:     iny                ; $f5d8: c8        
            lda #$37           ; $f5d9: a9 37     
            sta $47            ; $f5db: 85 47     
            jsr __fb0a         ; $f5dd: 20 0a fb  
            .hex 67 47         ; $f5e0: 67 47     Invalid Opcode - RRA $47
            nop                ; $f5e2: ea        
            nop                ; $f5e3: ea        
            nop                ; $f5e4: ea        
            nop                ; $f5e5: ea        
            jsr __fb10         ; $f5e6: 20 10 fb  
            lda $47            ; $f5e9: a5 47     
            cmp #$9b           ; $f5eb: c9 9b     
            beq __f5f1         ; $f5ed: f0 02     
            sty $00            ; $f5ef: 84 00     
__f5f1:     iny                ; $f5f1: c8        
            lda #$a5           ; $f5f2: a9 a5     
            sta $0647          ; $f5f4: 8d 47 06  
            jsr __fae9         ; $f5f7: 20 e9 fa  
            .hex 6f 47 06      ; $f5fa: 6f 47 06  Invalid Opcode - RRA $0647
            nop                ; $f5fd: ea        
            nop                ; $f5fe: ea        
            nop                ; $f5ff: ea        
            nop                ; $f600: ea        
            jsr __faef         ; $f601: 20 ef fa  
            lda $0647          ; $f604: ad 47 06  
            cmp #$52           ; $f607: c9 52     
            beq __f60d         ; $f609: f0 02     
            sty $00            ; $f60b: 84 00     
__f60d:     iny                ; $f60d: c8        
            lda #$29           ; $f60e: a9 29     
            sta $0647          ; $f610: 8d 47 06  
            jsr __fafa         ; $f613: 20 fa fa  
            .hex 6f 47 06      ; $f616: 6f 47 06  Invalid Opcode - RRA $0647
            nop                ; $f619: ea        
            nop                ; $f61a: ea        
            nop                ; $f61b: ea        
            nop                ; $f61c: ea        
            jsr __faff         ; $f61d: 20 ff fa  
            lda $0647          ; $f620: ad 47 06  
            cmp #$14           ; $f623: c9 14     
            beq __f629         ; $f625: f0 02     
            sty $00            ; $f627: 84 00     
__f629:     iny                ; $f629: c8        
            lda #$37           ; $f62a: a9 37     
            sta $0647          ; $f62c: 8d 47 06  
            jsr __fb0a         ; $f62f: 20 0a fb  
            .hex 6f 47 06      ; $f632: 6f 47 06  Invalid Opcode - RRA $0647
            nop                ; $f635: ea        
            nop                ; $f636: ea        
            nop                ; $f637: ea        
            nop                ; $f638: ea        
            jsr __fb10         ; $f639: 20 10 fb  
            lda $0647          ; $f63c: ad 47 06  
            cmp #$9b           ; $f63f: c9 9b     
            beq __f645         ; $f641: f0 02     
            sty $00            ; $f643: 84 00     
__f645:     lda #$a5           ; $f645: a9 a5     
            sta $0647          ; $f647: 8d 47 06  
            lda #$48           ; $f64a: a9 48     
            sta $45            ; $f64c: 85 45     
            lda #$05           ; $f64e: a9 05     
            sta $46            ; $f650: 85 46     
            ldy #$ff           ; $f652: a0 ff     
            jsr __fae9         ; $f654: 20 e9 fa  
            .hex 73 45         ; $f657: 73 45     Invalid Opcode - RRA ($45),y
            nop                ; $f659: ea        
            nop                ; $f65a: ea        
            php                ; $f65b: 08        
            pha                ; $f65c: 48        
            ldy #$0a           ; $f65d: a0 0a     
            pla                ; $f65f: 68        
            plp                ; $f660: 28        
            jsr __faef         ; $f661: 20 ef fa  
            lda $0647          ; $f664: ad 47 06  
            cmp #$52           ; $f667: c9 52     
            beq __f66d         ; $f669: f0 02     
            sty $00            ; $f66b: 84 00     
__f66d:     ldy #$ff           ; $f66d: a0 ff     
            lda #$29           ; $f66f: a9 29     
            sta $0647          ; $f671: 8d 47 06  
            jsr __fafa         ; $f674: 20 fa fa  
            .hex 73 45         ; $f677: 73 45     Invalid Opcode - RRA ($45),y
            nop                ; $f679: ea        
            nop                ; $f67a: ea        
            php                ; $f67b: 08        
            pha                ; $f67c: 48        
            ldy #$0b           ; $f67d: a0 0b     
            pla                ; $f67f: 68        
            plp                ; $f680: 28        
            jsr __faff         ; $f681: 20 ff fa  
            lda $0647          ; $f684: ad 47 06  
            cmp #$14           ; $f687: c9 14     
            beq __f68d         ; $f689: f0 02     
            sty $00            ; $f68b: 84 00     
__f68d:     ldy #$ff           ; $f68d: a0 ff     
            lda #$37           ; $f68f: a9 37     
            sta $0647          ; $f691: 8d 47 06  
            jsr __fb0a         ; $f694: 20 0a fb  
            .hex 73 45         ; $f697: 73 45     Invalid Opcode - RRA ($45),y
            nop                ; $f699: ea        
            nop                ; $f69a: ea        
            php                ; $f69b: 08        
            pha                ; $f69c: 48        
            ldy #$0c           ; $f69d: a0 0c     
            pla                ; $f69f: 68        
            plp                ; $f6a0: 28        
            jsr __fb10         ; $f6a1: 20 10 fb  
            lda $0647          ; $f6a4: ad 47 06  
            cmp #$9b           ; $f6a7: c9 9b     
            beq __f6ad         ; $f6a9: f0 02     
            sty $00            ; $f6ab: 84 00     
__f6ad:     ldy #$0d           ; $f6ad: a0 0d     
            ldx #$ff           ; $f6af: a2 ff     
            lda #$a5           ; $f6b1: a9 a5     
            sta $47            ; $f6b3: 85 47     
            jsr __fae9         ; $f6b5: 20 e9 fa  
            .hex 77 48         ; $f6b8: 77 48     Invalid Opcode - RRA $48,x
            nop                ; $f6ba: ea        
            nop                ; $f6bb: ea        
            nop                ; $f6bc: ea        
            nop                ; $f6bd: ea        
            jsr __faef         ; $f6be: 20 ef fa  
            lda $47            ; $f6c1: a5 47     
            cmp #$52           ; $f6c3: c9 52     
            beq __f6c9         ; $f6c5: f0 02     
            sty $00            ; $f6c7: 84 00     
__f6c9:     iny                ; $f6c9: c8        
            lda #$29           ; $f6ca: a9 29     
            sta $47            ; $f6cc: 85 47     
            jsr __fafa         ; $f6ce: 20 fa fa  
            .hex 77 48         ; $f6d1: 77 48     Invalid Opcode - RRA $48,x
            nop                ; $f6d3: ea        
            nop                ; $f6d4: ea        
            nop                ; $f6d5: ea        
            nop                ; $f6d6: ea        
            jsr __faff         ; $f6d7: 20 ff fa  
            lda $47            ; $f6da: a5 47     
            cmp #$14           ; $f6dc: c9 14     
            beq __f6e2         ; $f6de: f0 02     
            sty $00            ; $f6e0: 84 00     
__f6e2:     iny                ; $f6e2: c8        
            lda #$37           ; $f6e3: a9 37     
            sta $47            ; $f6e5: 85 47     
            jsr __fb0a         ; $f6e7: 20 0a fb  
            .hex 77 48         ; $f6ea: 77 48     Invalid Opcode - RRA $48,x
            nop                ; $f6ec: ea        
            nop                ; $f6ed: ea        
            nop                ; $f6ee: ea        
            nop                ; $f6ef: ea        
            jsr __fb10         ; $f6f0: 20 10 fb  
            lda $47            ; $f6f3: a5 47     
            cmp #$9b           ; $f6f5: c9 9b     
            beq __f6fb         ; $f6f7: f0 02     
            sty $00            ; $f6f9: 84 00     
__f6fb:     lda #$a5           ; $f6fb: a9 a5     
            sta $0647          ; $f6fd: 8d 47 06  
            ldy #$ff           ; $f700: a0 ff     
            jsr __fae9         ; $f702: 20 e9 fa  
            .hex 7b 48 05      ; $f705: 7b 48 05  Invalid Opcode - RRA $0548,y
            nop                ; $f708: ea        
            nop                ; $f709: ea        
            php                ; $f70a: 08        
            pha                ; $f70b: 48        
            ldy #$10           ; $f70c: a0 10     
            pla                ; $f70e: 68        
            plp                ; $f70f: 28        
            jsr __faef         ; $f710: 20 ef fa  
            lda $0647          ; $f713: ad 47 06  
            cmp #$52           ; $f716: c9 52     
            beq __f71c         ; $f718: f0 02     
            sty $00            ; $f71a: 84 00     
__f71c:     ldy #$ff           ; $f71c: a0 ff     
            lda #$29           ; $f71e: a9 29     
            sta $0647          ; $f720: 8d 47 06  
            jsr __fafa         ; $f723: 20 fa fa  
            .hex 7b 48 05      ; $f726: 7b 48 05  Invalid Opcode - RRA $0548,y
            nop                ; $f729: ea        
            nop                ; $f72a: ea        
            php                ; $f72b: 08        
            pha                ; $f72c: 48        
            ldy #$11           ; $f72d: a0 11     
            pla                ; $f72f: 68        
            plp                ; $f730: 28        
            jsr __faff         ; $f731: 20 ff fa  
            lda $0647          ; $f734: ad 47 06  
            cmp #$14           ; $f737: c9 14     
            beq __f73d         ; $f739: f0 02     
            sty $00            ; $f73b: 84 00     
__f73d:     ldy #$ff           ; $f73d: a0 ff     
            lda #$37           ; $f73f: a9 37     
            sta $0647          ; $f741: 8d 47 06  
            jsr __fb0a         ; $f744: 20 0a fb  
            .hex 7b 48 05      ; $f747: 7b 48 05  Invalid Opcode - RRA $0548,y
            nop                ; $f74a: ea        
            nop                ; $f74b: ea        
            php                ; $f74c: 08        
            pha                ; $f74d: 48        
            ldy #$12           ; $f74e: a0 12     
            pla                ; $f750: 68        
            plp                ; $f751: 28        
            jsr __fb10         ; $f752: 20 10 fb  
            lda $0647          ; $f755: ad 47 06  
            cmp #$9b           ; $f758: c9 9b     
            beq __f75e         ; $f75a: f0 02     
            sty $00            ; $f75c: 84 00     
__f75e:     ldy #$13           ; $f75e: a0 13     
            ldx #$ff           ; $f760: a2 ff     
            lda #$a5           ; $f762: a9 a5     
            sta $0647          ; $f764: 8d 47 06  
            jsr __fae9         ; $f767: 20 e9 fa  
            .hex 7f 48 05      ; $f76a: 7f 48 05  Invalid Opcode - RRA $0548,x
            nop                ; $f76d: ea        
            nop                ; $f76e: ea        
            nop                ; $f76f: ea        
            nop                ; $f770: ea        
            jsr __faef         ; $f771: 20 ef fa  
            lda $0647          ; $f774: ad 47 06  
            cmp #$52           ; $f777: c9 52     
            beq __f77d         ; $f779: f0 02     
            sty $00            ; $f77b: 84 00     
__f77d:     iny                ; $f77d: c8        
            lda #$29           ; $f77e: a9 29     
            sta $0647          ; $f780: 8d 47 06  
            jsr __fafa         ; $f783: 20 fa fa  
            .hex 7f 48 05      ; $f786: 7f 48 05  Invalid Opcode - RRA $0548,x
            nop                ; $f789: ea        
            nop                ; $f78a: ea        
            nop                ; $f78b: ea        
            nop                ; $f78c: ea        
            jsr __faff         ; $f78d: 20 ff fa  
            lda $0647          ; $f790: ad 47 06  
            cmp #$14           ; $f793: c9 14     
            beq __f799         ; $f795: f0 02     
            sty $00            ; $f797: 84 00     
__f799:     iny                ; $f799: c8        
            lda #$37           ; $f79a: a9 37     
            sta $0647          ; $f79c: 8d 47 06  
            jsr __fb0a         ; $f79f: 20 0a fb  
            .hex 7f 48 05      ; $f7a2: 7f 48 05  Invalid Opcode - RRA $0548,x
            nop                ; $f7a5: ea        
            nop                ; $f7a6: ea        
            nop                ; $f7a7: ea        
            nop                ; $f7a8: ea        
            jsr __fb10         ; $f7a9: 20 10 fb  
            lda $0647          ; $f7ac: ad 47 06  
            cmp #$9b           ; $f7af: c9 9b     
            beq __f7b5         ; $f7b1: f0 02     
            sty $00            ; $f7b3: 84 00     
__f7b5:     rts                ; $f7b5: 60        

;-------------------------------------------------------------------------------
__f7b6:     clc                ; $f7b6: 18        
            lda #$ff           ; $f7b7: a9 ff     
            sta $01            ; $f7b9: 85 01     
            bit $01            ; $f7bb: 24 01     
            lda #$55           ; $f7bd: a9 55     
            rts                ; $f7bf: 60        

;-------------------------------------------------------------------------------
__f7c0:     bcs __f7cb         ; $f7c0: b0 09     
            bpl __f7cb         ; $f7c2: 10 07     
            cmp #$ff           ; $f7c4: c9 ff     
            bne __f7cb         ; $f7c6: d0 03     
            bvc __f7cb         ; $f7c8: 50 01     
            rts                ; $f7ca: 60        

;-------------------------------------------------------------------------------
__f7cb:     sty $00            ; $f7cb: 84 00     
            rts                ; $f7cd: 60        

;-------------------------------------------------------------------------------
__f7ce:     sec                ; $f7ce: 38        
            clv                ; $f7cf: b8        
            lda #$00           ; $f7d0: a9 00     
            rts                ; $f7d2: 60        

;-------------------------------------------------------------------------------
__f7d3:     bne __f7dc         ; $f7d3: d0 07     
            bvs __f7dc         ; $f7d5: 70 05     
            bcc __f7dc         ; $f7d7: 90 03     
            bmi __f7dc         ; $f7d9: 30 01     
            rts                ; $f7db: 60        

;-------------------------------------------------------------------------------
__f7dc:     sty $00            ; $f7dc: 84 00     
            rts                ; $f7de: 60        

;-------------------------------------------------------------------------------
__f7df:     clc                ; $f7df: 18        
            bit $01            ; $f7e0: 24 01     
            lda #$55           ; $f7e2: a9 55     
            rts                ; $f7e4: 60        

;-------------------------------------------------------------------------------
__f7e5:     bne __f7ee         ; $f7e5: d0 07     
            bvc __f7ee         ; $f7e7: 50 05     
            bcs __f7ee         ; $f7e9: b0 03     
            bmi __f7ee         ; $f7eb: 30 01     
            rts                ; $f7ed: 60        

;-------------------------------------------------------------------------------
__f7ee:     sty $00            ; $f7ee: 84 00     
            rts                ; $f7f0: 60        

;-------------------------------------------------------------------------------
__f7f1:     sec                ; $f7f1: 38        
            clv                ; $f7f2: b8        
            lda #$f8           ; $f7f3: a9 f8     
            rts                ; $f7f5: 60        

;-------------------------------------------------------------------------------
__f7f6:     bcc __f801         ; $f7f6: 90 09     
            bpl __f801         ; $f7f8: 10 07     
            cmp #$e8           ; $f7fa: c9 e8     
            bne __f801         ; $f7fc: d0 03     
            bvs __f801         ; $f7fe: 70 01     
            rts                ; $f800: 60        

;-------------------------------------------------------------------------------
__f801:     sty $00            ; $f801: 84 00     
            rts                ; $f803: 60        

;-------------------------------------------------------------------------------
__f804:     clc                ; $f804: 18        
            bit $01            ; $f805: 24 01     
            lda #$5f           ; $f807: a9 5f     
            rts                ; $f809: 60        

;-------------------------------------------------------------------------------
__f80a:     bcs __f815         ; $f80a: b0 09     
            bpl __f815         ; $f80c: 10 07     
            cmp #$f5           ; $f80e: c9 f5     
            bne __f815         ; $f810: d0 03     
            bvc __f815         ; $f812: 50 01     
            rts                ; $f814: 60        

;-------------------------------------------------------------------------------
__f815:     sty $00            ; $f815: 84 00     
            rts                ; $f817: 60        

;-------------------------------------------------------------------------------
__f818:     sec                ; $f818: 38        
            clv                ; $f819: b8        
            lda #$70           ; $f81a: a9 70     
            rts                ; $f81c: 60        

;-------------------------------------------------------------------------------
__f81d:     bne __f826         ; $f81d: d0 07     
            bvs __f826         ; $f81f: 70 05     
            bcc __f826         ; $f821: 90 03     
            bmi __f826         ; $f823: 30 01     
            rts                ; $f825: 60        

;-------------------------------------------------------------------------------
__f826:     sty $00            ; $f826: 84 00     
            rts                ; $f828: 60        

;-------------------------------------------------------------------------------
__f829:     clc                ; $f829: 18        
            bit $01            ; $f82a: 24 01     
            lda #$00           ; $f82c: a9 00     
            rts                ; $f82e: 60        

;-------------------------------------------------------------------------------
__f82f:     bmi __f83a         ; $f82f: 30 09     
            bcs __f83a         ; $f831: b0 07     
            cmp #$69           ; $f833: c9 69     
            bne __f83a         ; $f835: d0 03     
            bvs __f83a         ; $f837: 70 01     
            rts                ; $f839: 60        

;-------------------------------------------------------------------------------
__f83a:     sty $00            ; $f83a: 84 00     
            rts                ; $f83c: 60        

;-------------------------------------------------------------------------------
__f83d:     sec                ; $f83d: 38        
            bit $01            ; $f83e: 24 01     
            lda #$00           ; $f840: a9 00     
            rts                ; $f842: 60        

;-------------------------------------------------------------------------------
__f843:     bmi __f84e         ; $f843: 30 09     
            bcs __f84e         ; $f845: b0 07     
            cmp #$6a           ; $f847: c9 6a     
            bne __f84e         ; $f849: d0 03     
            bvs __f84e         ; $f84b: 70 01     
            rts                ; $f84d: 60        

;-------------------------------------------------------------------------------
__f84e:     sty $00            ; $f84e: 84 00     
            rts                ; $f850: 60        

;-------------------------------------------------------------------------------
__f851:     sec                ; $f851: 38        
            clv                ; $f852: b8        
            lda #$7f           ; $f853: a9 7f     
            rts                ; $f855: 60        

;-------------------------------------------------------------------------------
__f856:     bpl __f861         ; $f856: 10 09     
            bcs __f861         ; $f858: b0 07     
            cmp #$ff           ; $f85a: c9 ff     
            bne __f861         ; $f85c: d0 03     
            bvc __f861         ; $f85e: 50 01     
            rts                ; $f860: 60        

;-------------------------------------------------------------------------------
__f861:     sty $00            ; $f861: 84 00     
            rts                ; $f863: 60        

;-------------------------------------------------------------------------------
__f864:     clc                ; $f864: 18        
            bit $01            ; $f865: 24 01     
            lda #$7f           ; $f867: a9 7f     
            rts                ; $f869: 60        

;-------------------------------------------------------------------------------
__f86a:     bpl __f875         ; $f86a: 10 09     
            bcs __f875         ; $f86c: b0 07     
            cmp #$ff           ; $f86e: c9 ff     
            bne __f875         ; $f870: d0 03     
            bvs __f875         ; $f872: 70 01     
            rts                ; $f874: 60        

;-------------------------------------------------------------------------------
__f875:     sty $00            ; $f875: 84 00     
            rts                ; $f877: 60        

;-------------------------------------------------------------------------------
__f878:     sec                ; $f878: 38        
            clv                ; $f879: b8        
            lda #$7f           ; $f87a: a9 7f     
            rts                ; $f87c: 60        

;-------------------------------------------------------------------------------
__f87d:     bne __f886         ; $f87d: d0 07     
            bmi __f886         ; $f87f: 30 05     
            bvs __f886         ; $f881: 70 03     
            bcc __f886         ; $f883: 90 01     
            rts                ; $f885: 60        

;-------------------------------------------------------------------------------
__f886:     sty $00            ; $f886: 84 00     
            rts                ; $f888: 60        

;-------------------------------------------------------------------------------
__f889:     bit $01            ; $f889: 24 01     
            lda #$40           ; $f88b: a9 40     
            rts                ; $f88d: 60        

;-------------------------------------------------------------------------------
__f88e:     bmi __f897         ; $f88e: 30 07     
            bcc __f897         ; $f890: 90 05     
            bne __f897         ; $f892: d0 03     
            bvc __f897         ; $f894: 50 01     
            rts                ; $f896: 60        

;-------------------------------------------------------------------------------
__f897:     sty $00            ; $f897: 84 00     
            rts                ; $f899: 60        

;-------------------------------------------------------------------------------
__f89a:     clv                ; $f89a: b8        
            rts                ; $f89b: 60        

;-------------------------------------------------------------------------------
__f89c:     beq __f8a5         ; $f89c: f0 07     
            bmi __f8a5         ; $f89e: 30 05     
            bcc __f8a5         ; $f8a0: 90 03     
            bvs __f8a5         ; $f8a2: 70 01     
            rts                ; $f8a4: 60        

;-------------------------------------------------------------------------------
__f8a5:     sty $00            ; $f8a5: 84 00     
            rts                ; $f8a7: 60        

;-------------------------------------------------------------------------------
__f8a8:     beq __f8af         ; $f8a8: f0 05     
            bpl __f8af         ; $f8aa: 10 03     
            bpl __f8af         ; $f8ac: 10 01     
            rts                ; $f8ae: 60        

;-------------------------------------------------------------------------------
__f8af:     sty $00            ; $f8af: 84 00     
            rts                ; $f8b1: 60        

;-------------------------------------------------------------------------------
__f8b2:     lda #$80           ; $f8b2: a9 80     
            rts                ; $f8b4: 60        

;-------------------------------------------------------------------------------
__f8b5:     beq __f8bc         ; $f8b5: f0 05     
            bpl __f8bc         ; $f8b7: 10 03     
            bcc __f8bc         ; $f8b9: 90 01     
            rts                ; $f8bb: 60        

;-------------------------------------------------------------------------------
__f8bc:     sty $00            ; $f8bc: 84 00     
            rts                ; $f8be: 60        

;-------------------------------------------------------------------------------
__f8bf:     bne __f8c6         ; $f8bf: d0 05     
            bmi __f8c6         ; $f8c1: 30 03     
            bcc __f8c6         ; $f8c3: 90 01     
            rts                ; $f8c5: 60        

;-------------------------------------------------------------------------------
__f8c6:     sty $00            ; $f8c6: 84 00     
            rts                ; $f8c8: 60        

;-------------------------------------------------------------------------------
__f8c9:     bcs __f8d0         ; $f8c9: b0 05     
__f8cb:     beq __f8d0         ; $f8cb: f0 03     
            bpl __f8d0         ; $f8cd: 10 01     
            rts                ; $f8cf: 60        

;-------------------------------------------------------------------------------
__f8d0:     sty $00            ; $f8d0: 84 00     
            rts                ; $f8d2: 60        

;-------------------------------------------------------------------------------
__f8d3:     bcc __f8da         ; $f8d3: 90 05     
            beq __f8da         ; $f8d5: f0 03     
            bmi __f8da         ; $f8d7: 30 01     
            rts                ; $f8d9: 60        

;-------------------------------------------------------------------------------
__f8da:     sty $00            ; $f8da: 84 00     
            rts                ; $f8dc: 60        

;-------------------------------------------------------------------------------
__f8dd:     bit $01            ; $f8dd: 24 01     
            ldy #$40           ; $f8df: a0 40     
            rts                ; $f8e1: 60        

;-------------------------------------------------------------------------------
__f8e2:     bmi __f8eb         ; $f8e2: 30 07     
            bcc __f8eb         ; $f8e4: 90 05     
            bne __f8eb         ; $f8e6: d0 03     
            bvc __f8eb         ; $f8e8: 50 01     
            rts                ; $f8ea: 60        

;-------------------------------------------------------------------------------
__f8eb:     stx $00            ; $f8eb: 86 00     
            rts                ; $f8ed: 60        

;-------------------------------------------------------------------------------
__f8ee:     clv                ; $f8ee: b8        
            rts                ; $f8ef: 60        

;-------------------------------------------------------------------------------
__f8f0:     beq __f8f9         ; $f8f0: f0 07     
            bmi __f8f9         ; $f8f2: 30 05     
            bcc __f8f9         ; $f8f4: 90 03     
            bvs __f8f9         ; $f8f6: 70 01     
            rts                ; $f8f8: 60        

;-------------------------------------------------------------------------------
__f8f9:     stx $00            ; $f8f9: 86 00     
            rts                ; $f8fb: 60        

;-------------------------------------------------------------------------------
__f8fc:     beq __f903         ; $f8fc: f0 05     
            bpl __f903         ; $f8fe: 10 03     
            bpl __f903         ; $f900: 10 01     
            rts                ; $f902: 60        

;-------------------------------------------------------------------------------
__f903:     stx $00            ; $f903: 86 00     
            rts                ; $f905: 60        

;-------------------------------------------------------------------------------
__f906:     ldy #$80           ; $f906: a0 80     
            rts                ; $f908: 60        

;-------------------------------------------------------------------------------
__f909:     beq __f910         ; $f909: f0 05     
            bpl __f910         ; $f90b: 10 03     
            bcc __f910         ; $f90d: 90 01     
            rts                ; $f90f: 60        

;-------------------------------------------------------------------------------
__f910:     stx $00            ; $f910: 86 00     
            rts                ; $f912: 60        

;-------------------------------------------------------------------------------
__f913:     bne __f91a         ; $f913: d0 05     
            bmi __f91a         ; $f915: 30 03     
            bcc __f91a         ; $f917: 90 01     
            rts                ; $f919: 60        

;-------------------------------------------------------------------------------
__f91a:     stx $00            ; $f91a: 86 00     
            rts                ; $f91c: 60        

;-------------------------------------------------------------------------------
__f91d:     bcs __f924         ; $f91d: b0 05     
            beq __f924         ; $f91f: f0 03     
            bpl __f924         ; $f921: 10 01     
            rts                ; $f923: 60        

;-------------------------------------------------------------------------------
__f924:     stx $00            ; $f924: 86 00     
            rts                ; $f926: 60        

;-------------------------------------------------------------------------------
__f927:     bcc __f92e         ; $f927: 90 05     
            beq __f92e         ; $f929: f0 03     
            bmi __f92e         ; $f92b: 30 01     
            rts                ; $f92d: 60        

;-------------------------------------------------------------------------------
__f92e:     stx $00            ; $f92e: 86 00     
            rts                ; $f930: 60        

;-------------------------------------------------------------------------------
__f931:     bit $01            ; $f931: 24 01     
            lda #$40           ; $f933: a9 40     
            sec                ; $f935: 38        
            rts                ; $f936: 60        

;-------------------------------------------------------------------------------
__f937:     bmi __f944         ; $f937: 30 0b     
            bcc __f944         ; $f939: 90 09     
            bne __f944         ; $f93b: d0 07     
            bvs __f944         ; $f93d: 70 05     
            cmp #$00           ; $f93f: c9 00     
            bne __f944         ; $f941: d0 01     
            rts                ; $f943: 60        

;-------------------------------------------------------------------------------
__f944:     sty $00            ; $f944: 84 00     
            rts                ; $f946: 60        

;-------------------------------------------------------------------------------
__f947:     clv                ; $f947: b8        
            sec                ; $f948: 38        
            lda #$40           ; $f949: a9 40     
            rts                ; $f94b: 60        

;-------------------------------------------------------------------------------
__f94c:     beq __f959         ; $f94c: f0 0b     
            bmi __f959         ; $f94e: 30 09     
            bcc __f959         ; $f950: 90 07     
            bvs __f959         ; $f952: 70 05     
            cmp #$01           ; $f954: c9 01     
            bne __f959         ; $f956: d0 01     
            rts                ; $f958: 60        

;-------------------------------------------------------------------------------
__f959:     sty $00            ; $f959: 84 00     
            rts                ; $f95b: 60        

;-------------------------------------------------------------------------------
__f95c:     lda #$40           ; $f95c: a9 40     
            sec                ; $f95e: 38        
            bit $01            ; $f95f: 24 01     
            rts                ; $f961: 60        

;-------------------------------------------------------------------------------
__f962:     bcs __f96f         ; $f962: b0 0b     
            beq __f96f         ; $f964: f0 09     
            bpl __f96f         ; $f966: 10 07     
            bvs __f96f         ; $f968: 70 05     
            cmp #$ff           ; $f96a: c9 ff     
            bne __f96f         ; $f96c: d0 01     
            rts                ; $f96e: 60        

;-------------------------------------------------------------------------------
__f96f:     sty $00            ; $f96f: 84 00     
            rts                ; $f971: 60        

;-------------------------------------------------------------------------------
__f972:     clc                ; $f972: 18        
            lda #$80           ; $f973: a9 80     
            rts                ; $f975: 60        

;-------------------------------------------------------------------------------
__f976:     bcc __f97d         ; $f976: 90 05     
            cmp #$7f           ; $f978: c9 7f     
            bne __f97d         ; $f97a: d0 01     
            rts                ; $f97c: 60        

;-------------------------------------------------------------------------------
__f97d:     sty $00            ; $f97d: 84 00     
            rts                ; $f97f: 60        

;-------------------------------------------------------------------------------
__f980:     sec                ; $f980: 38        
            lda #$81           ; $f981: a9 81     
            rts                ; $f983: 60        

;-------------------------------------------------------------------------------
__f984:     bvc __f98d         ; $f984: 50 07     
            bcc __f98d         ; $f986: 90 05     
            cmp #$02           ; $f988: c9 02     
            bne __f98d         ; $f98a: d0 01     
            rts                ; $f98c: 60        

;-------------------------------------------------------------------------------
__f98d:     sty $00            ; $f98d: 84 00     
            rts                ; $f98f: 60        

;-------------------------------------------------------------------------------
__f990:     ldx #$55           ; $f990: a2 55     
            lda #$ff           ; $f992: a9 ff     
            sta $01            ; $f994: 85 01     
            nop                ; $f996: ea        
            bit $01            ; $f997: 24 01     
            sec                ; $f999: 38        
            lda #$01           ; $f99a: a9 01     
            rts                ; $f99c: 60        

;-------------------------------------------------------------------------------
__f99d:     bcc __f9ba         ; $f99d: 90 1b     
            bne __f9ba         ; $f99f: d0 19     
            bmi __f9ba         ; $f9a1: 30 17     
            bvc __f9ba         ; $f9a3: 50 15     
            cmp #$00           ; $f9a5: c9 00     
            bne __f9ba         ; $f9a7: d0 11     
            clv                ; $f9a9: b8        
            lda #$aa           ; $f9aa: a9 aa     
            rts                ; $f9ac: 60        

;-------------------------------------------------------------------------------
__f9ad:     bcs __f9ba         ; $f9ad: b0 0b     
            beq __f9ba         ; $f9af: f0 09     
            bmi __f9ba         ; $f9b1: 30 07     
            bvs __f9ba         ; $f9b3: 70 05     
            cmp #$55           ; $f9b5: c9 55     
            bne __f9ba         ; $f9b7: d0 01     
            rts                ; $f9b9: 60        

;-------------------------------------------------------------------------------
__f9ba:     sty $00            ; $f9ba: 84 00     
            rts                ; $f9bc: 60        

;-------------------------------------------------------------------------------
__f9bd:     bit $01            ; $f9bd: 24 01     
            sec                ; $f9bf: 38        
            lda #$80           ; $f9c0: a9 80     
            rts                ; $f9c2: 60        

;-------------------------------------------------------------------------------
__f9c3:     bcc __f9e1         ; $f9c3: 90 1c     
            bne __f9e1         ; $f9c5: d0 1a     
            bmi __f9e1         ; $f9c7: 30 18     
            bvc __f9e1         ; $f9c9: 50 16     
            cmp #$00           ; $f9cb: c9 00     
            bne __f9e1         ; $f9cd: d0 12     
            clv                ; $f9cf: b8        
            lda #$55           ; $f9d0: a9 55     
            sec                ; $f9d2: 38        
            rts                ; $f9d3: 60        

;-------------------------------------------------------------------------------
__f9d4:     bcs __f9e1         ; $f9d4: b0 0b     
            beq __f9e1         ; $f9d6: f0 09     
            bpl __f9e1         ; $f9d8: 10 07     
            bvs __f9e1         ; $f9da: 70 05     
            cmp #$aa           ; $f9dc: c9 aa     
            bne __f9e1         ; $f9de: d0 01     
            rts                ; $f9e0: 60        

;-------------------------------------------------------------------------------
__f9e1:     sty $00            ; $f9e1: 84 00     
            rts                ; $f9e3: 60        

;-------------------------------------------------------------------------------
__f9e4:     bit $01            ; $f9e4: 24 01     
            sec                ; $f9e6: 38        
            lda #$01           ; $f9e7: a9 01     
            rts                ; $f9e9: 60        

;-------------------------------------------------------------------------------
__f9ea:     bcc __fa08         ; $f9ea: 90 1c     
            beq __fa08         ; $f9ec: f0 1a     
            bpl __fa08         ; $f9ee: 10 18     
            bvc __fa08         ; $f9f0: 50 16     
            cmp #$80           ; $f9f2: c9 80     
            bne __fa08         ; $f9f4: d0 12     
            clv                ; $f9f6: b8        
            clc                ; $f9f7: 18        
            lda #$55           ; $f9f8: a9 55     
            rts                ; $f9fa: 60        

;-------------------------------------------------------------------------------
__f9fb:     bcc __fa08         ; $f9fb: 90 0b     
            beq __fa08         ; $f9fd: f0 09     
            bmi __fa08         ; $f9ff: 30 07     
            bvs __fa08         ; $fa01: 70 05     
            cmp #$2a           ; $fa03: c9 2a     
            bne __fa08         ; $fa05: d0 01     
            rts                ; $fa07: 60        

;-------------------------------------------------------------------------------
__fa08:     sty $00            ; $fa08: 84 00     
__fa0a:     bit $01            ; $fa0a: 24 01     
            sec                ; $fa0c: 38        
            lda #$80           ; $fa0d: a9 80     
            rts                ; $fa0f: 60        

;-------------------------------------------------------------------------------
__fa10:     bcc __fa2e         ; $fa10: 90 1c     
            beq __fa2e         ; $fa12: f0 1a     
            bmi __fa2e         ; $fa14: 30 18     
            bvc __fa2e         ; $fa16: 50 16     
            cmp #$01           ; $fa18: c9 01     
            bne __fa2e         ; $fa1a: d0 12     
            clv                ; $fa1c: b8        
            clc                ; $fa1d: 18        
            lda #$55           ; $fa1e: a9 55     
            rts                ; $fa20: 60        

;-------------------------------------------------------------------------------
__fa21:     bcs __fa2e         ; $fa21: b0 0b     
            beq __fa2e         ; $fa23: f0 09     
            bpl __fa2e         ; $fa25: 10 07     
            bvs __fa2e         ; $fa27: 70 05     
            cmp #$aa           ; $fa29: c9 aa     
            bne __fa2e         ; $fa2b: d0 01     
            rts                ; $fa2d: 60        

;-------------------------------------------------------------------------------
__fa2e:     sty $00            ; $fa2e: 84 00     
            rts                ; $fa30: 60        

;-------------------------------------------------------------------------------
__fa31:     bit $01            ; $fa31: 24 01     
            clc                ; $fa33: 18        
            lda #$40           ; $fa34: a9 40     
            rts                ; $fa36: 60        

;-------------------------------------------------------------------------------
__fa37:     bvc __fa65         ; $fa37: 50 2c     
            bcs __fa65         ; $fa39: b0 2a     
            bmi __fa65         ; $fa3b: 30 28     
            cmp #$40           ; $fa3d: c9 40     
            bne __fa65         ; $fa3f: d0 24     
            rts                ; $fa41: 60        

;-------------------------------------------------------------------------------
__fa42:     clv                ; $fa42: b8        
            sec                ; $fa43: 38        
            lda #$ff           ; $fa44: a9 ff     
            rts                ; $fa46: 60        

;-------------------------------------------------------------------------------
__fa47:     bvs __fa65         ; $fa47: 70 1c     
            bne __fa65         ; $fa49: d0 1a     
            bmi __fa65         ; $fa4b: 30 18     
            bcc __fa65         ; $fa4d: 90 16     
            cmp #$ff           ; $fa4f: c9 ff     
            bne __fa65         ; $fa51: d0 12     
            rts                ; $fa53: 60        

;-------------------------------------------------------------------------------
__fa54:     bit $01            ; $fa54: 24 01     
            lda #$f0           ; $fa56: a9 f0     
            rts                ; $fa58: 60        

;-------------------------------------------------------------------------------
__fa59:     bvc __fa65         ; $fa59: 50 0a     
            beq __fa65         ; $fa5b: f0 08     
            bpl __fa65         ; $fa5d: 10 06     
            bcc __fa65         ; $fa5f: 90 04     
            cmp #$f0           ; $fa61: c9 f0     
            beq __fa67         ; $fa63: f0 02     
__fa65:     sty $00            ; $fa65: 84 00     
__fa67:     rts                ; $fa67: 60        

;-------------------------------------------------------------------------------
__fa68:     bit $01            ; $fa68: 24 01     
            sec                ; $fa6a: 38        
            lda #$75           ; $fa6b: a9 75     
            rts                ; $fa6d: 60        

;-------------------------------------------------------------------------------
__fa6e:     bvc __fae6         ; $fa6e: 50 76     
            beq __fae6         ; $fa70: f0 74     
            bmi __fae6         ; $fa72: 30 72     
            bcs __fae6         ; $fa74: b0 70     
            cmp #$65           ; $fa76: c9 65     
            bne __fae6         ; $fa78: d0 6c     
            rts                ; $fa7a: 60        

;-------------------------------------------------------------------------------
__fa7b:     bit $01            ; $fa7b: 24 01     
            clc                ; $fa7d: 18        
            lda #$b3           ; $fa7e: a9 b3     
            rts                ; $fa80: 60        

;-------------------------------------------------------------------------------
__fa81:     bvc __fae6         ; $fa81: 50 63     
            bcc __fae6         ; $fa83: 90 61     
            bpl __fae6         ; $fa85: 10 5f     
            cmp #$fb           ; $fa87: c9 fb     
            bne __fae6         ; $fa89: d0 5b     
            rts                ; $fa8b: 60        

;-------------------------------------------------------------------------------
__fa8c:     clv                ; $fa8c: b8        
            clc                ; $fa8d: 18        
            lda #$c3           ; $fa8e: a9 c3     
            rts                ; $fa90: 60        

;-------------------------------------------------------------------------------
__fa91:     bvs __fae6         ; $fa91: 70 53     
            beq __fae6         ; $fa93: f0 51     
            bpl __fae6         ; $fa95: 10 4f     
            bcs __fae6         ; $fa97: b0 4d     
            cmp #$d3           ; $fa99: c9 d3     
            bne __fae6         ; $fa9b: d0 49     
            rts                ; $fa9d: 60        

;-------------------------------------------------------------------------------
__fa9e:     bit $01            ; $fa9e: 24 01     
            sec                ; $faa0: 38        
            lda #$10           ; $faa1: a9 10     
            rts                ; $faa3: 60        

;-------------------------------------------------------------------------------
__faa4:     bvc __fae6         ; $faa4: 50 40     
            beq __fae6         ; $faa6: f0 3e     
            bmi __fae6         ; $faa8: 30 3c     
            bcs __fae6         ; $faaa: b0 3a     
            cmp #$7e           ; $faac: c9 7e     
            bne __fae6         ; $faae: d0 36     
            rts                ; $fab0: 60        

;-------------------------------------------------------------------------------
__fab1:     bit $01            ; $fab1: 24 01     
            clc                ; $fab3: 18        
            lda #$40           ; $fab4: a9 40     
            rts                ; $fab6: 60        

;-------------------------------------------------------------------------------
__fab7:     bvs __fae6         ; $fab7: 70 2d     
            bcs __fae6         ; $fab9: b0 2b     
            bmi __fae6         ; $fabb: 30 29     
            cmp #$53           ; $fabd: c9 53     
            bne __fae6         ; $fabf: d0 25     
            rts                ; $fac1: 60        

;-------------------------------------------------------------------------------
__fac2:     clv                ; $fac2: b8        
            sec                ; $fac3: 38        
            lda #$ff           ; $fac4: a9 ff     
            rts                ; $fac6: 60        

;-------------------------------------------------------------------------------
__fac7:     bvs __fae6         ; $fac7: 70 1d     
            beq __fae6         ; $fac9: f0 1b     
            bpl __fae6         ; $facb: 10 19     
            bcc __fae6         ; $facd: 90 17     
            cmp #$ff           ; $facf: c9 ff     
            bne __fae6         ; $fad1: d0 13     
            rts                ; $fad3: 60        

;-------------------------------------------------------------------------------
__fad4:     bit $01            ; $fad4: 24 01     
            sec                ; $fad6: 38        
            lda #$f0           ; $fad7: a9 f0     
            rts                ; $fad9: 60        

;-------------------------------------------------------------------------------
__fada:     bvs __fae6         ; $fada: 70 0a     
            beq __fae6         ; $fadc: f0 08     
            bpl __fae6         ; $fade: 10 06     
            bcc __fae6         ; $fae0: 90 04     
            cmp #$b8           ; $fae2: c9 b8     
            beq __fae8         ; $fae4: f0 02     
__fae6:     sty $00            ; $fae6: 84 00     
__fae8:     rts                ; $fae8: 60        

;-------------------------------------------------------------------------------
__fae9:     bit $01            ; $fae9: 24 01     
            clc                ; $faeb: 18        
            lda #$b2           ; $faec: a9 b2     
            rts                ; $faee: 60        

;-------------------------------------------------------------------------------
__faef:     bvs __fb1b         ; $faef: 70 2a     
            bcc __fb1b         ; $faf1: 90 28     
            bmi __fb1b         ; $faf3: 30 26     
            cmp #$05           ; $faf5: c9 05     
            bne __fb1b         ; $faf7: d0 22     
            rts                ; $faf9: 60        

;-------------------------------------------------------------------------------
__fafa:     clv                ; $fafa: b8        
            clc                ; $fafb: 18        
            lda #$42           ; $fafc: a9 42     
            rts                ; $fafe: 60        

;-------------------------------------------------------------------------------
__faff:     bvs __fb1b         ; $faff: 70 1a     
            bmi __fb1b         ; $fb01: 30 18     
            bcs __fb1b         ; $fb03: b0 16     
            cmp #$57           ; $fb05: c9 57     
            bne __fb1b         ; $fb07: d0 12     
            rts                ; $fb09: 60        

;-------------------------------------------------------------------------------
__fb0a:     bit $01            ; $fb0a: 24 01     
            sec                ; $fb0c: 38        
            lda #$75           ; $fb0d: a9 75     
            rts                ; $fb0f: 60        

;-------------------------------------------------------------------------------
__fb10:     bvs __fb1b         ; $fb10: 70 09     
            bmi __fb1b         ; $fb12: 30 07     
            bcc __fb1b         ; $fb14: 90 05     
            cmp #$11           ; $fb16: c9 11     
            bne __fb1b         ; $fb18: d0 01     
            rts                ; $fb1a: 60        

;-------------------------------------------------------------------------------
__fb1b:     sta $00            ; $fb1b: 85 00     
__fb1d:     bit $01            ; $fb1d: 24 01     
            clc                ; $fb1f: 18        
            lda #$b3           ; $fb20: a9 b3     
            rts                ; $fb22: 60        

;-------------------------------------------------------------------------------
__fb23:     bvc __fb75         ; $fb23: 50 50     
            bcc __fb75         ; $fb25: 90 4e     
            bpl __fb75         ; $fb27: 10 4c     
            cmp #$e1           ; $fb29: c9 e1     
            bne __fb75         ; $fb2b: d0 48     
            rts                ; $fb2d: 60        

;-------------------------------------------------------------------------------
__fb2e:     clv                ; $fb2e: b8        
            clc                ; $fb2f: 18        
            lda #$42           ; $fb30: a9 42     
            rts                ; $fb32: 60        

;-------------------------------------------------------------------------------
__fb33:     bvs __fb75         ; $fb33: 70 40     
            beq __fb75         ; $fb35: f0 3e     
            bmi __fb75         ; $fb37: 30 3c     
            bcc __fb75         ; $fb39: 90 3a     
            cmp #$56           ; $fb3b: c9 56     
            bne __fb75         ; $fb3d: d0 36     
            rts                ; $fb3f: 60        

;-------------------------------------------------------------------------------
__fb40:     bit $01            ; $fb40: 24 01     
            sec                ; $fb42: 38        
            lda #$75           ; $fb43: a9 75     
            rts                ; $fb45: 60        

;-------------------------------------------------------------------------------
__fb46:     bvc __fb75         ; $fb46: 50 2d     
            beq __fb75         ; $fb48: f0 2b     
            bmi __fb75         ; $fb4a: 30 29     
            bcc __fb75         ; $fb4c: 90 27     
            cmp #$6e           ; $fb4e: c9 6e     
            bne __fb75         ; $fb50: d0 23     
            rts                ; $fb52: 60        

;-------------------------------------------------------------------------------
__fb53:     bit $01            ; $fb53: 24 01     
            clc                ; $fb55: 18        
            lda #$b3           ; $fb56: a9 b3     
            rts                ; $fb58: 60        

;-------------------------------------------------------------------------------
__fb59:     bvc __fb75         ; $fb59: 50 1a     
            bcc __fb75         ; $fb5b: 90 18     
            bmi __fb75         ; $fb5d: 30 16     
            cmp #$02           ; $fb5f: c9 02     
            bne __fb75         ; $fb61: d0 12     
            rts                ; $fb63: 60        

;-------------------------------------------------------------------------------
__fb64:     clv                ; $fb64: b8        
            clc                ; $fb65: 18        
            lda #$42           ; $fb66: a9 42     
            rts                ; $fb68: 60        

;-------------------------------------------------------------------------------
__fb69:     bvs __fb75         ; $fb69: 70 0a     
            beq __fb75         ; $fb6b: f0 08     
            bmi __fb75         ; $fb6d: 30 06     
            bcs __fb75         ; $fb6f: b0 04     
            cmp #$42           ; $fb71: c9 42     
            beq __fb77         ; $fb73: f0 02     
__fb75:     sty $00            ; $fb75: 84 00     
__fb77:     rts                ; $fb77: 60        

;-------------------------------------------------------------------------------
            brk                ; $fb78: 00        
            brk                ; $fb79: 00        
            brk                ; $fb7a: 00        
            brk                ; $fb7b: 00        
            brk                ; $fb7c: 00        
            brk                ; $fb7d: 00        
            brk                ; $fb7e: 00        
            brk                ; $fb7f: 00        
            brk                ; $fb80: 00        
            brk                ; $fb81: 00        
            brk                ; $fb82: 00        
            brk                ; $fb83: 00        
            brk                ; $fb84: 00        
            brk                ; $fb85: 00        
            brk                ; $fb86: 00        
            brk                ; $fb87: 00        
            .hex 80 80         ; $fb88: 80 80     Invalid Opcode - NOP #$80
            .hex ff 80 80      ; $fb8a: ff 80 80  Invalid Opcode - ISC $8080,x
            brk                ; $fb8d: 00        
            brk                ; $fb8e: 00        
            brk                ; $fb8f: 00        
            brk                ; $fb90: 00        
            brk                ; $fb91: 00        
            .hex ff 00 00      ; $fb92: ff 00 00  Bad Addr Mode - ISC $0000,x
            brk                ; $fb95: 00        
            brk                ; $fb96: 00        
            brk                ; $fb97: 00        
            ora ($01,x)        ; $fb98: 01 01     
            .hex ff 01 01      ; $fb9a: ff 01 01  Invalid Opcode - ISC $0101,x
__fb9d:     brk                ; $fb9d: 00        
            brk                ; $fb9e: 00        
            brk                ; $fb9f: 00        
            brk                ; $fba0: 00        
            brk                ; $fba1: 00        
            brk                ; $fba2: 00        
            brk                ; $fba3: 00        
            brk                ; $fba4: 00        
            brk                ; $fba5: 00        
            brk                ; $fba6: 00        
            brk                ; $fba7: 00        
            .hex 7c fe 00      ; $fba8: 7c fe 00  Bad Addr Mode - NOP $00fe,x
            cpy #$c0           ; $fbab: c0 c0     
            .hex fe 7c 00      ; $fbad: fe 7c 00  Bad Addr Mode - INC $007c,x
            .hex fe fe 00      ; $fbb0: fe fe 00  Bad Addr Mode - INC $00fe,x
            beq __fb75         ; $fbb3: f0 c0     
            .hex fe fe 00      ; $fbb5: fe fe 00  Bad Addr Mode - INC $00fe,x
            dec $c6            ; $fbb8: c6 c6     
            .hex 02            ; $fbba: 02        Invalid Opcode - KIL 
            inc __c6c6,x       ; $fbbb: fe c6 c6  
            dec $00            ; $fbbe: c6 00     
            .hex cc d8 00      ; $fbc0: cc d8 00  Bad Addr Mode - CPY $00d8
            beq __fb9d         ; $fbc3: f0 d8     
            .hex cc c6 00      ; $fbc5: cc c6 00  Bad Addr Mode - CPY $00c6
            dec $ee            ; $fbc8: c6 ee     
            .hex 02            ; $fbca: 02        Invalid Opcode - KIL 
            dec $c6,x          ; $fbcb: d6 c6     
            dec $c6            ; $fbcd: c6 c6     
            brk                ; $fbcf: 00        
            dec $c6            ; $fbd0: c6 c6     
            .hex 02            ; $fbd2: 02        Invalid Opcode - KIL 
            dec $ce,x          ; $fbd3: d6 ce     
            dec $c6            ; $fbd5: c6 c6     
            brk                ; $fbd7: 00        
            .hex 7c fe 02      ; $fbd8: 7c fe 02  Invalid Opcode - NOP $02fe,x
            dec $c6            ; $fbdb: c6 c6     
            .hex fe 7c 00      ; $fbdd: fe 7c 00  Bad Addr Mode - INC $007c,x
            .hex fc fe 02      ; $fbe0: fc fe 02  Invalid Opcode - NOP $02fe,x
            .hex fc c0 c0      ; $fbe3: fc c0 c0  Invalid Opcode - NOP __c0c0,x
            cpy #$00           ; $fbe6: c0 00     
            .hex cc cc 00      ; $fbe8: cc cc 00  Bad Addr Mode - CPY $00cc
            sei                ; $fbeb: 78        
            bmi __fc1e         ; $fbec: 30 30     
            bmi __fbf0         ; $fbee: 30 00     
__fbf0:     clc                ; $fbf0: 18        
            clc                ; $fbf1: 18        
            clc                ; $fbf2: 18        
            clc                ; $fbf3: 18        
            clc                ; $fbf4: 18        
            clc                ; $fbf5: 18        
            clc                ; $fbf6: 18        
            brk                ; $fbf7: 00        
            .hex fc fe 02      ; $fbf8: fc fe 02  Invalid Opcode - NOP $02fe,x
            asl $1c            ; $fbfb: 06 1c     
__fbfd:     bvs __fbfd         ; $fbfd: 70 fe     
            brk                ; $fbff: 00        
            .hex fc fe 02      ; $fc00: fc fe 02  Invalid Opcode - NOP $02fe,x
            .hex 3c 3c 02      ; $fc03: 3c 3c 02  Invalid Opcode - NOP $023c,x
            inc $1800,x        ; $fc06: fe 00 18  
            clc                ; $fc09: 18        
            cld                ; $fc0a: d8        
            cld                ; $fc0b: d8        
            inc $1818,x        ; $fc0c: fe 18 18  
            brk                ; $fc0f: 00        
            .hex fe fe 00      ; $fc10: fe fe 00  Bad Addr Mode - INC $00fe,x
            .hex 80 fc         ; $fc13: 80 fc     Invalid Opcode - NOP #$fc
            asl $fe            ; $fc15: 06 fe     
            brk                ; $fc17: 00        
            .hex 7c fe 00      ; $fc18: 7c fe 00  Bad Addr Mode - NOP $00fe,x
            cpy #$fc           ; $fc1b: c0 fc     
            .hex c6            ; $fc1d: c6        Suspected data
__fc1e:     inc __fe00,x       ; $fc1e: fe 00 fe  
            inc $0c06,x        ; $fc21: fe 06 0c  
            clc                ; $fc24: 18        
            bpl __fc57         ; $fc25: 10 30     
            brk                ; $fc27: 00        
            brk                ; $fc28: 00        
            brk                ; $fc29: 00        
            brk                ; $fc2a: 00        
            brk                ; $fc2b: 00        
            brk                ; $fc2c: 00        
            brk                ; $fc2d: 00        
            brk                ; $fc2e: 00        
            brk                ; $fc2f: 00        
            brk                ; $fc30: 00        
            brk                ; $fc31: 00        
            brk                ; $fc32: 00        
            brk                ; $fc33: 00        
            brk                ; $fc34: 00        
            brk                ; $fc35: 00        
            brk                ; $fc36: 00        
            brk                ; $fc37: 00        
            brk                ; $fc38: 00        
            brk                ; $fc39: 00        
            brk                ; $fc3a: 00        
            brk                ; $fc3b: 00        
            brk                ; $fc3c: 00        
            brk                ; $fc3d: 00        
            brk                ; $fc3e: 00        
            brk                ; $fc3f: 00        
            brk                ; $fc40: 00        
            brk                ; $fc41: 00        
            brk                ; $fc42: 00        
            brk                ; $fc43: 00        
            brk                ; $fc44: 00        
            brk                ; $fc45: 00        
            brk                ; $fc46: 00        
            brk                ; $fc47: 00        
            brk                ; $fc48: 00        
            brk                ; $fc49: 00        
            brk                ; $fc4a: 00        
            brk                ; $fc4b: 00        
            brk                ; $fc4c: 00        
            brk                ; $fc4d: 00        
            brk                ; $fc4e: 00        
            brk                ; $fc4f: 00        
            brk                ; $fc50: 00        
            brk                ; $fc51: 00        
            brk                ; $fc52: 00        
            brk                ; $fc53: 00        
            brk                ; $fc54: 00        
            brk                ; $fc55: 00        
            brk                ; $fc56: 00        
__fc57:     brk                ; $fc57: 00        
            brk                ; $fc58: 00        
            brk                ; $fc59: 00        
            brk                ; $fc5a: 00        
            brk                ; $fc5b: 00        
            brk                ; $fc5c: 00        
            brk                ; $fc5d: 00        
            brk                ; $fc5e: 00        
            brk                ; $fc5f: 00        
            brk                ; $fc60: 00        
            brk                ; $fc61: 00        
            brk                ; $fc62: 00        
            brk                ; $fc63: 00        
            brk                ; $fc64: 00        
            brk                ; $fc65: 00        
            brk                ; $fc66: 00        
            brk                ; $fc67: 00        
            clc                ; $fc68: 18        
            clc                ; $fc69: 18        
            clc                ; $fc6a: 18        
            .hex ff ff 18      ; $fc6b: ff ff 18  Invalid Opcode - ISC $18ff,x
            clc                ; $fc6e: 18        
            clc                ; $fc6f: 18        
            clc                ; $fc70: 18        
            clc                ; $fc71: 18        
            clc                ; $fc72: 18        
            .hex ff ff 00      ; $fc73: ff ff 00  Bad Addr Mode - ISC $00ff,x
            brk                ; $fc76: 00        
            brk                ; $fc77: 00        
            brk                ; $fc78: 00        
            brk                ; $fc79: 00        
            brk                ; $fc7a: 00        
            brk                ; $fc7b: 00        
            brk                ; $fc7c: 00        
            brk                ; $fc7d: 00        
            brk                ; $fc7e: 00        
            brk                ; $fc7f: 00        
            clc                ; $fc80: 18        
            clc                ; $fc81: 18        
            clc                ; $fc82: 18        
            clc                ; $fc83: 18        
            brk                ; $fc84: 00        
            clc                ; $fc85: 18        
            clc                ; $fc86: 18        
            brk                ; $fc87: 00        
            .hex 33 33         ; $fc88: 33 33     Invalid Opcode - RLA ($33),y
            ror $00            ; $fc8a: 66 00     
            brk                ; $fc8c: 00        
            brk                ; $fc8d: 00        
            brk                ; $fc8e: 00        
            brk                ; $fc8f: 00        
            ror $66            ; $fc90: 66 66     
            .hex ff 66 ff      ; $fc92: ff 66 ff  Invalid Opcode - ISC __ff66,x
            ror $66            ; $fc95: 66 66     
            brk                ; $fc97: 00        
            clc                ; $fc98: 18        
            rol $3c60,x        ; $fc99: 3e 60 3c  
            asl $7c            ; $fc9c: 06 7c     
            clc                ; $fc9e: 18        
            brk                ; $fc9f: 00        
            .hex 62            ; $fca0: 62        Invalid Opcode - KIL 
            ror $0c            ; $fca1: 66 0c     
            clc                ; $fca3: 18        
            bmi __fd0c         ; $fca4: 30 66     
            lsr $00            ; $fca6: 46 00     
            .hex 3c 66 3c      ; $fca8: 3c 66 3c  Invalid Opcode - NOP $3c66,x
            sec                ; $fcab: 38        
            .hex 67 66         ; $fcac: 67 66     Invalid Opcode - RRA $66
            .hex 3f 00 0c      ; $fcae: 3f 00 0c  Invalid Opcode - RLA $0c00,x
            .hex 0c 18 00      ; $fcb1: 0c 18 00  Bad Addr Mode - NOP $0018
            brk                ; $fcb4: 00        
            brk                ; $fcb5: 00        
            brk                ; $fcb6: 00        
            brk                ; $fcb7: 00        
            .hex 0c 18 30      ; $fcb8: 0c 18 30  Invalid Opcode - NOP $3018
            bmi __fced         ; $fcbb: 30 30     
            clc                ; $fcbd: 18        
            .hex 0c 00 30      ; $fcbe: 0c 00 30  Invalid Opcode - NOP $3000
            clc                ; $fcc1: 18        
            .hex 0c 0c 0c      ; $fcc2: 0c 0c 0c  Invalid Opcode - NOP $0c0c
            clc                ; $fcc5: 18        
            bmi __fcc8         ; $fcc6: 30 00     
__fcc8:     brk                ; $fcc8: 00        
            ror $3c            ; $fcc9: 66 3c     
            .hex ff 3c 66      ; $fccb: ff 3c 66  Invalid Opcode - ISC $663c,x
            brk                ; $fcce: 00        
            brk                ; $fccf: 00        
            brk                ; $fcd0: 00        
            clc                ; $fcd1: 18        
            clc                ; $fcd2: 18        
            ror $1818,x        ; $fcd3: 7e 18 18  
            brk                ; $fcd6: 00        
            brk                ; $fcd7: 00        
            brk                ; $fcd8: 00        
            brk                ; $fcd9: 00        
            brk                ; $fcda: 00        
            brk                ; $fcdb: 00        
            brk                ; $fcdc: 00        
            clc                ; $fcdd: 18        
            clc                ; $fcde: 18        
            bmi __fce1         ; $fcdf: 30 00     
__fce1:     brk                ; $fce1: 00        
            brk                ; $fce2: 00        
            .hex 6e 3b 00      ; $fce3: 6e 3b 00  Bad Addr Mode - ROR $003b
            brk                ; $fce6: 00        
            brk                ; $fce7: 00        
            brk                ; $fce8: 00        
            brk                ; $fce9: 00        
            brk                ; $fcea: 00        
            brk                ; $fceb: 00        
            brk                ; $fcec: 00        
__fced:     clc                ; $fced: 18        
            clc                ; $fcee: 18        
            brk                ; $fcef: 00        
            brk                ; $fcf0: 00        
            .hex 03 06         ; $fcf1: 03 06     Invalid Opcode - SLO ($06,x)
            .hex 0c 18 30      ; $fcf3: 0c 18 30  Invalid Opcode - NOP $3018
            rts                ; $fcf6: 60        

;-------------------------------------------------------------------------------
            brk                ; $fcf7: 00        
            rol $6763,x        ; $fcf8: 3e 63 67  
            .hex 6b 73         ; $fcfb: 6b 73     Invalid Opcode - ARR #$73
            .hex 63 3e         ; $fcfd: 63 3e     Invalid Opcode - RRA ($3e,x)
            brk                ; $fcff: 00        
            .hex 0c 1c 0c      ; $fd00: 0c 1c 0c  Invalid Opcode - NOP $0c1c
            .hex 0c 0c 0c      ; $fd03: 0c 0c 0c  Invalid Opcode - NOP $0c0c
            .hex 3f 00 3e      ; $fd06: 3f 00 3e  Invalid Opcode - RLA $3e00,x
            .hex 63 63         ; $fd09: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 0e            ; $fd0b: 0e        Suspected data
__fd0c:     sec                ; $fd0c: 38        
            .hex 63 7f         ; $fd0d: 63 7f     Invalid Opcode - RRA ($7f,x)
            brk                ; $fd0f: 00        
            rol $6363,x        ; $fd10: 3e 63 63  
            asl $6363          ; $fd13: 0e 63 63  
            rol $0600,x        ; $fd16: 3e 00 06  
            asl $261e          ; $fd19: 0e 1e 26  
            .hex 7f 06 06      ; $fd1c: 7f 06 06  Invalid Opcode - RRA $0606,x
            brk                ; $fd1f: 00        
            .hex 7f 63 60      ; $fd20: 7f 63 60  Invalid Opcode - RRA $6063,x
            ror $6303,x        ; $fd23: 7e 03 63  
            rol $3e00,x        ; $fd26: 3e 00 3e  
            .hex 63 60         ; $fd29: 63 60     Invalid Opcode - RRA ($60,x)
            ror $6363,x        ; $fd2b: 7e 63 63  
            rol $7f00,x        ; $fd2e: 3e 00 7f  
            .hex 63 06         ; $fd31: 63 06     Invalid Opcode - RRA ($06,x)
            .hex 0c 18 18      ; $fd33: 0c 18 18  Invalid Opcode - NOP $1818
            .hex 3c 00 3e      ; $fd36: 3c 00 3e  Invalid Opcode - NOP $3e00,x
            .hex 63 63         ; $fd39: 63 63     Invalid Opcode - RRA ($63,x)
            rol $6363,x        ; $fd3b: 3e 63 63  
            rol $3e00,x        ; $fd3e: 3e 00 3e  
            .hex 63 63         ; $fd41: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 3f 03 63      ; $fd43: 3f 03 63  Invalid Opcode - RLA $6303,x
            .hex 3e 00 00      ; $fd46: 3e 00 00  Bad Addr Mode - ROL $0000,x
            brk                ; $fd49: 00        
            clc                ; $fd4a: 18        
            clc                ; $fd4b: 18        
            brk                ; $fd4c: 00        
            clc                ; $fd4d: 18        
            clc                ; $fd4e: 18        
            brk                ; $fd4f: 00        
            brk                ; $fd50: 00        
            brk                ; $fd51: 00        
            clc                ; $fd52: 18        
            clc                ; $fd53: 18        
            brk                ; $fd54: 00        
            clc                ; $fd55: 18        
            clc                ; $fd56: 18        
            bmi __fd67         ; $fd57: 30 0e     
            clc                ; $fd59: 18        
            bmi __fdbc         ; $fd5a: 30 60     
            bmi __fd76         ; $fd5c: 30 18     
            .hex 0e 00 00      ; $fd5e: 0e 00 00  Bad Addr Mode - ASL $0000
            brk                ; $fd61: 00        
            ror $7e00,x        ; $fd62: 7e 00 7e  
            brk                ; $fd65: 00        
            brk                ; $fd66: 00        
__fd67:     brk                ; $fd67: 00        
            bvs __fd82         ; $fd68: 70 18     
            .hex 0c 06 0c      ; $fd6a: 0c 06 0c  Invalid Opcode - NOP $0c06
            clc                ; $fd6d: 18        
            bvs __fd70         ; $fd6e: 70 00     
__fd70:     ror $0363,x        ; $fd70: 7e 63 03  
            asl $1c            ; $fd73: 06 1c     
            brk                ; $fd75: 00        
__fd76:     clc                ; $fd76: 18        
            clc                ; $fd77: 18        
            .hex 7c c6 ce      ; $fd78: 7c c6 ce  Invalid Opcode - NOP __cec6,x
            inc __e6e0         ; $fd7b: ee e0 e6  
            .hex 7c 00 1c      ; $fd7e: 7c 00 1c  Invalid Opcode - NOP $1c00,x
            .hex 36            ; $fd81: 36        Suspected data
__fd82:     .hex 63 7f         ; $fd82: 63 7f     Invalid Opcode - RRA ($7f,x)
            .hex 63 63         ; $fd84: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 63 00         ; $fd86: 63 00     Invalid Opcode - RRA ($00,x)
            ror $6373          ; $fd88: 6e 73 63  
            ror $6363,x        ; $fd8b: 7e 63 63  
            ror $1e00,x        ; $fd8e: 7e 00 1e  
            .hex 33 60         ; $fd91: 33 60     Invalid Opcode - RLA ($60),y
            rts                ; $fd93: 60        

;-------------------------------------------------------------------------------
            rts                ; $fd94: 60        

;-------------------------------------------------------------------------------
            .hex 33 1e         ; $fd95: 33 1e     Invalid Opcode - RLA ($1e),y
            brk                ; $fd97: 00        
            jmp ($6376)        ; $fd98: 6c 76 63  

;-------------------------------------------------------------------------------
            .hex 63 63         ; $fd9b: 63 63     Invalid Opcode - RRA ($63,x)
            ror $7c            ; $fd9d: 66 7c     
            brk                ; $fd9f: 00        
            .hex 7f 31 30      ; $fda0: 7f 31 30  Invalid Opcode - RRA $3031,x
            .hex 3c 30 31      ; $fda3: 3c 30 31  Invalid Opcode - NOP $3130,x
            .hex 7f 00 7f      ; $fda6: 7f 00 7f  Invalid Opcode - RRA $7f00,x
            and ($30),y        ; $fda9: 31 30     
            .hex 3c 30 30      ; $fdab: 3c 30 30  Invalid Opcode - NOP $3030,x
            sei                ; $fdae: 78        
            brk                ; $fdaf: 00        
            asl $6033,x        ; $fdb0: 1e 33 60  
            .hex 67 63         ; $fdb3: 67 63     Invalid Opcode - RRA $63
            .hex 37 1d         ; $fdb5: 37 1d     Invalid Opcode - RLA $1d,x
            brk                ; $fdb7: 00        
            .hex 63 63         ; $fdb8: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 63 7f         ; $fdba: 63 7f     Invalid Opcode - RRA ($7f,x)
__fdbc:     .hex 63 63         ; $fdbc: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 63 00         ; $fdbe: 63 00     Invalid Opcode - RRA ($00,x)
            .hex 3c 18 18      ; $fdc0: 3c 18 18  Invalid Opcode - NOP $1818,x
            clc                ; $fdc3: 18        
            clc                ; $fdc4: 18        
            clc                ; $fdc5: 18        
            .hex 3c 00 1f      ; $fdc6: 3c 00 1f  Invalid Opcode - NOP $1f00,x
            asl $06            ; $fdc9: 06 06     
            asl $06            ; $fdcb: 06 06     
            ror $3c            ; $fdcd: 66 3c     
            brk                ; $fdcf: 00        
            ror $66            ; $fdd0: 66 66     
            jmp ($6c78)        ; $fdd2: 6c 78 6c  

;-------------------------------------------------------------------------------
            .hex 67 63         ; $fdd5: 67 63     Invalid Opcode - RRA $63
            brk                ; $fdd7: 00        
            sei                ; $fdd8: 78        
            bmi __fe3b         ; $fdd9: 30 60     
            rts                ; $fddb: 60        

;-------------------------------------------------------------------------------
            .hex 63 63         ; $fddc: 63 63     Invalid Opcode - RRA ($63,x)
            ror $6300,x        ; $fdde: 7e 00 63  
            .hex 77 7f         ; $fde1: 77 7f     Invalid Opcode - RRA $7f,x
            .hex 6b 63         ; $fde3: 6b 63     Invalid Opcode - ARR #$63
            .hex 63 63         ; $fde5: 63 63     Invalid Opcode - RRA ($63,x)
            brk                ; $fde7: 00        
            .hex 63 73         ; $fde8: 63 73     Invalid Opcode - RRA ($73,x)
            .hex 7b 6f 67      ; $fdea: 7b 6f 67  Invalid Opcode - RRA $676f,y
            .hex 63 63         ; $fded: 63 63     Invalid Opcode - RRA ($63,x)
            brk                ; $fdef: 00        
            .hex 1c 36 63      ; $fdf0: 1c 36 63  Invalid Opcode - NOP $6336,x
            .hex 63 63         ; $fdf3: 63 63     Invalid Opcode - RRA ($63,x)
            rol $1c,x          ; $fdf5: 36 1c     
            brk                ; $fdf7: 00        
            ror $6373          ; $fdf8: 6e 73 63  
            ror $6060,x        ; $fdfb: 7e 60 60  
            rts                ; $fdfe: 60        

;-------------------------------------------------------------------------------
            brk                ; $fdff: 00        
__fe00:     .hex 1c 36 63      ; $fe00: 1c 36 63  Invalid Opcode - NOP $6336,x
            .hex 6b 67         ; $fe03: 6b 67     Invalid Opcode - ARR #$67
            rol $1d,x          ; $fe05: 36 1d     
            brk                ; $fe07: 00        
            ror $6373          ; $fe08: 6e 73 63  
            ror $676c,x        ; $fe0b: 7e 6c 67  
            .hex 63 00         ; $fe0e: 63 00     Invalid Opcode - RRA ($00,x)
            rol $6063,x        ; $fe10: 3e 63 60  
            rol $6303,x        ; $fe13: 3e 03 63  
            rol $7e00,x        ; $fe16: 3e 00 7e  
            .hex 5a            ; $fe19: 5a        Invalid Opcode - NOP 
            clc                ; $fe1a: 18        
            clc                ; $fe1b: 18        
            clc                ; $fe1c: 18        
            clc                ; $fe1d: 18        
            .hex 3c 00 73      ; $fe1e: 3c 00 73  Invalid Opcode - NOP $7300,x
            .hex 33 63         ; $fe21: 33 63     Invalid Opcode - RLA ($63),y
            .hex 63 63         ; $fe23: 63 63     Invalid Opcode - RRA ($63,x)
            ror $3c,x          ; $fe25: 76 3c     
            brk                ; $fe27: 00        
            .hex 73 33         ; $fe28: 73 33     Invalid Opcode - RRA ($33),y
            .hex 63 63         ; $fe2a: 63 63     Invalid Opcode - RRA ($63,x)
            ror $3c            ; $fe2c: 66 3c     
            clc                ; $fe2e: 18        
            brk                ; $fe2f: 00        
            .hex 73 33         ; $fe30: 73 33     Invalid Opcode - RRA ($33),y
            .hex 63 6b         ; $fe32: 63 6b     Invalid Opcode - RRA ($6b,x)
            .hex 7f 77 63      ; $fe34: 7f 77 63  Invalid Opcode - RRA $6377,x
            brk                ; $fe37: 00        
            .hex 63 63         ; $fe38: 63 63     Invalid Opcode - RRA ($63,x)
            .hex 36            ; $fe3a: 36        Suspected data
__fe3b:     .hex 1c 36 63      ; $fe3b: 1c 36 63  Invalid Opcode - NOP $6336,x
            .hex 63 00         ; $fe3e: 63 00     Invalid Opcode - RRA ($00,x)
            .hex 33 63         ; $fe40: 33 63     Invalid Opcode - RLA ($63),y
            .hex 63 36         ; $fe42: 63 36     Invalid Opcode - RRA ($36,x)
            .hex 1c 78 70      ; $fe44: 1c 78 70  Invalid Opcode - NOP $7078,x
            brk                ; $fe47: 00        
            .hex 7f 63 06      ; $fe48: 7f 63 06  Invalid Opcode - RRA $0663,x
            .hex 1c 33 63      ; $fe4b: 1c 33 63  Invalid Opcode - NOP $6333,x
            ror $3c00,x        ; $fe4e: 7e 00 3c  
            bmi __fe83         ; $fe51: 30 30     
            bmi __fe85         ; $fe53: 30 30     
            bmi __fe93         ; $fe55: 30 3c     
            brk                ; $fe57: 00        
            rti                ; $fe58: 40        

;-------------------------------------------------------------------------------
            rts                ; $fe59: 60        

;-------------------------------------------------------------------------------
            bmi __fe74         ; $fe5a: 30 18     
            .hex 0c 06 02      ; $fe5c: 0c 06 02  Invalid Opcode - NOP $0206
            brk                ; $fe5f: 00        
            .hex 3c 0c 0c      ; $fe60: 3c 0c 0c  Invalid Opcode - NOP $0c0c,x
            .hex 0c 0c 0c      ; $fe63: 0c 0c 0c  Invalid Opcode - NOP $0c0c
            .hex 3c 00 00      ; $fe66: 3c 00 00  Bad Addr Mode - NOP $0000,x
            clc                ; $fe69: 18        
            .hex 3c 7e 18      ; $fe6a: 3c 7e 18  Invalid Opcode - NOP $187e,x
            clc                ; $fe6d: 18        
            clc                ; $fe6e: 18        
            clc                ; $fe6f: 18        
            brk                ; $fe70: 00        
            brk                ; $fe71: 00        
            brk                ; $fe72: 00        
            brk                ; $fe73: 00        
__fe74:     brk                ; $fe74: 00        
            brk                ; $fe75: 00        
            .hex ff ff 30      ; $fe76: ff ff 30  Invalid Opcode - ISC $30ff,x
            bmi __fe93         ; $fe79: 30 18     
            brk                ; $fe7b: 00        
            brk                ; $fe7c: 00        
            brk                ; $fe7d: 00        
            brk                ; $fe7e: 00        
            brk                ; $fe7f: 00        
            brk                ; $fe80: 00        
            brk                ; $fe81: 00        
            .hex 3f            ; $fe82: 3f        Suspected data
__fe83:     .hex 63 63         ; $fe83: 63 63     Invalid Opcode - RRA ($63,x)
__fe85:     .hex 67 3b         ; $fe85: 67 3b     Invalid Opcode - RRA $3b
            brk                ; $fe87: 00        
            rts                ; $fe88: 60        

;-------------------------------------------------------------------------------
            rts                ; $fe89: 60        

;-------------------------------------------------------------------------------
            ror $6373          ; $fe8a: 6e 73 63  
            .hex 63 3e         ; $fe8d: 63 3e     Invalid Opcode - RRA ($3e,x)
            brk                ; $fe8f: 00        
            brk                ; $fe90: 00        
            brk                ; $fe91: 00        
            .hex 3e            ; $fe92: 3e        Suspected data
__fe93:     .hex 63 60         ; $fe93: 63 60     Invalid Opcode - RRA ($60,x)
            .hex 63 3e         ; $fe95: 63 3e     Invalid Opcode - RRA ($3e,x)
            brk                ; $fe97: 00        
            .hex 03 03         ; $fe98: 03 03     Invalid Opcode - SLO ($03,x)
            .hex 3b 67 63      ; $fe9a: 3b 67 63  Invalid Opcode - RLA $6367,y
            .hex 63 3e         ; $fe9d: 63 3e     Invalid Opcode - RRA ($3e,x)
            brk                ; $fe9f: 00        
            brk                ; $fea0: 00        
            brk                ; $fea1: 00        
            rol $7f61,x        ; $fea2: 3e 61 7f  
            rts                ; $fea5: 60        

;-------------------------------------------------------------------------------
            rol $0e00,x        ; $fea6: 3e 00 0e  
            clc                ; $fea9: 18        
            clc                ; $feaa: 18        
            .hex 3c 18 18      ; $feab: 3c 18 18  Invalid Opcode - NOP $1818,x
            .hex 3c 00 00      ; $feae: 3c 00 00  Bad Addr Mode - NOP $0000,x
            brk                ; $feb1: 00        
            rol $6360,x        ; $feb2: 3e 60 63  
            .hex 63 3d         ; $feb5: 63 3d     Invalid Opcode - RRA ($3d,x)
            brk                ; $feb7: 00        
            rts                ; $feb8: 60        

;-------------------------------------------------------------------------------
            rts                ; $feb9: 60        

;-------------------------------------------------------------------------------
            ror $6373          ; $feba: 6e 73 63  
            ror $67            ; $febd: 66 67     
            brk                ; $febf: 00        
            brk                ; $fec0: 00        
            brk                ; $fec1: 00        
            asl $0c0c,x        ; $fec2: 1e 0c 0c  
            .hex 0c 1e 00      ; $fec5: 0c 1e 00  Bad Addr Mode - NOP $001e
            brk                ; $fec8: 00        
            brk                ; $fec9: 00        
            .hex 3f 06 06      ; $feca: 3f 06 06  Invalid Opcode - RLA $0606,x
            asl $66            ; $fecd: 06 66     
            .hex 3c 60 60      ; $fecf: 3c 60 60  Invalid Opcode - NOP $6060,x
            ror $6e            ; $fed2: 66 6e     
            .hex 7c 67 63      ; $fed4: 7c 67 63  Invalid Opcode - NOP $6367,x
            brk                ; $fed7: 00        
            .hex 1c 0c 0c      ; $fed8: 1c 0c 0c  Invalid Opcode - NOP $0c0c,x
            .hex 0c 0c 0c      ; $fedb: 0c 0c 0c  Invalid Opcode - NOP $0c0c
            .hex 1e 00 00      ; $fede: 1e 00 00  Bad Addr Mode - ASL $0000,x
            brk                ; $fee1: 00        
            ror $6b7f          ; $fee2: 6e 7f 6b  
            .hex 62            ; $fee5: 62        Invalid Opcode - KIL 
            .hex 67 00         ; $fee6: 67 00     Invalid Opcode - RRA $00
            brk                ; $fee8: 00        
            brk                ; $fee9: 00        
            ror $6373          ; $feea: 6e 73 63  
            ror $67            ; $feed: 66 67     
            brk                ; $feef: 00        
            brk                ; $fef0: 00        
            brk                ; $fef1: 00        
            rol $6363,x        ; $fef2: 3e 63 63  
            .hex 63 3e         ; $fef5: 63 3e     Invalid Opcode - RRA ($3e,x)
            brk                ; $fef7: 00        
            brk                ; $fef8: 00        
            brk                ; $fef9: 00        
            rol $7363,x        ; $fefa: 3e 63 73  
            ror $6060          ; $fefd: 6e 60 60  
            brk                ; $ff00: 00        
            brk                ; $ff01: 00        
            rol $6763,x        ; $ff02: 3e 63 67  
            .hex 3b 03 03      ; $ff05: 3b 03 03  Invalid Opcode - RLA $0303,y
            brk                ; $ff08: 00        
            brk                ; $ff09: 00        
            ror $6373          ; $ff0a: 6e 73 63  
            .hex 7e 63 00      ; $ff0d: 7e 63 00  Bad Addr Mode - ROR $0063,x
            brk                ; $ff10: 00        
            brk                ; $ff11: 00        
            rol $1c71,x        ; $ff12: 3e 71 1c  
            .hex 47 3e         ; $ff15: 47 3e     Invalid Opcode - SRE $3e
            brk                ; $ff17: 00        
            asl $0c            ; $ff18: 06 0c     
            .hex 3f 18 18      ; $ff1a: 3f 18 18  Invalid Opcode - RLA $1818,x
            .hex 1b 0e 00      ; $ff1d: 1b 0e 00  Invalid Opcode - SLO $000e,y
            brk                ; $ff20: 00        
            brk                ; $ff21: 00        
            .hex 73 33         ; $ff22: 73 33     Invalid Opcode - RRA ($33),y
            .hex 63 67         ; $ff24: 63 67     Invalid Opcode - RRA ($67,x)
            .hex 3b 00 00      ; $ff26: 3b 00 00  Invalid Opcode - RLA $0000,y
            brk                ; $ff29: 00        
            .hex 73 33         ; $ff2a: 73 33     Invalid Opcode - RRA ($33),y
            .hex 63 66         ; $ff2c: 63 66     Invalid Opcode - RRA ($66,x)
            .hex 3c 00 00      ; $ff2e: 3c 00 00  Bad Addr Mode - NOP $0000,x
            brk                ; $ff31: 00        
            .hex 63 6b         ; $ff32: 63 6b     Invalid Opcode - RRA ($6b,x)
            .hex 7f 77 63      ; $ff34: 7f 77 63  Invalid Opcode - RRA $6377,x
            brk                ; $ff37: 00        
            brk                ; $ff38: 00        
            brk                ; $ff39: 00        
            .hex 63 36         ; $ff3a: 63 36     Invalid Opcode - RRA ($36,x)
            .hex 1c 36 63      ; $ff3c: 1c 36 63  Invalid Opcode - NOP $6336,x
            brk                ; $ff3f: 00        
            brk                ; $ff40: 00        
            brk                ; $ff41: 00        
            .hex 33 63         ; $ff42: 33 63     Invalid Opcode - RLA ($63),y
            .hex 63 3f         ; $ff44: 63 3f     Invalid Opcode - RRA ($3f,x)
            .hex 03 3e         ; $ff46: 03 3e     Invalid Opcode - SLO ($3e,x)
            brk                ; $ff48: 00        
            brk                ; $ff49: 00        
            .hex 7f 0e 1c      ; $ff4a: 7f 0e 1c  Invalid Opcode - RRA $1c0e,x
            sec                ; $ff4d: 38        
            .hex 7f 00 3c      ; $ff4e: 7f 00 3c  Invalid Opcode - RRA $3c00,x
            .hex 42            ; $ff51: 42        Invalid Opcode - KIL 
            sta $a1a1,y        ; $ff52: 99 a1 a1  
            sta $3c42,y        ; $ff55: 99 42 3c  
            brk                ; $ff58: 00        
            brk                ; $ff59: 00        
            brk                ; $ff5a: 00        
            brk                ; $ff5b: 00        
            brk                ; $ff5c: 00        
            brk                ; $ff5d: 00        
            brk                ; $ff5e: 00        
            brk                ; $ff5f: 00        
            brk                ; $ff60: 00        
            brk                ; $ff61: 00        
            brk                ; $ff62: 00        
            brk                ; $ff63: 00        
            brk                ; $ff64: 00        
            brk                ; $ff65: 00        
            brk                ; $ff66: 00        
            brk                ; $ff67: 00        
            brk                ; $ff68: 00        
            brk                ; $ff69: 00        
            brk                ; $ff6a: 00        
            brk                ; $ff6b: 00        
            brk                ; $ff6c: 00        
            brk                ; $ff6d: 00        
            brk                ; $ff6e: 00        
            brk                ; $ff6f: 00        
            brk                ; $ff70: 00        
            brk                ; $ff71: 00        
            brk                ; $ff72: 00        
            brk                ; $ff73: 00        
            brk                ; $ff74: 00        
            brk                ; $ff75: 00        
            brk                ; $ff76: 00        
            brk                ; $ff77: 00        
__ff78:     .hex 0f 06 12      ; $ff78: 0f 06 12  Invalid Opcode - SLO $1206
            .hex 33 33         ; $ff7b: 33 33     Invalid Opcode - RLA ($33),y
            asl $12            ; $ff7d: 06 12     
            .hex 33 38         ; $ff7f: 33 38     Invalid Opcode - RLA ($38),y
            asl $12            ; $ff81: 06 12     
            .hex 33 3a         ; $ff83: 33 3a     Invalid Opcode - RLA ($3a),y
            asl $12            ; $ff85: 06 12     
            .hex 33 0f         ; $ff87: 33 0f     Invalid Opcode - RLA ($0f),y
            asl $12            ; $ff89: 06 12     
            .hex 33 33         ; $ff8b: 33 33     Invalid Opcode - RLA ($33),y
            asl $12            ; $ff8d: 06 12     
            .hex 33 38         ; $ff8f: 33 38     Invalid Opcode - RLA ($38),y
            asl $12            ; $ff91: 06 12     
            .hex 33 3a         ; $ff93: 33 3a     Invalid Opcode - RLA ($3a),y
            asl $12            ; $ff95: 06 12     
            .hex 33 00         ; $ff97: 33 00     Invalid Opcode - RLA ($00),y
            brk                ; $ff99: 00        
            brk                ; $ff9a: 00        
            brk                ; $ff9b: 00        
            brk                ; $ff9c: 00        
            brk                ; $ff9d: 00        
            brk                ; $ff9e: 00        
            brk                ; $ff9f: 00        
            brk                ; $ffa0: 00        
            brk                ; $ffa1: 00        
            brk                ; $ffa2: 00        
            brk                ; $ffa3: 00        
            brk                ; $ffa4: 00        
            brk                ; $ffa5: 00        
            brk                ; $ffa6: 00        
            brk                ; $ffa7: 00        
            brk                ; $ffa8: 00        
            brk                ; $ffa9: 00        
            brk                ; $ffaa: 00        
            brk                ; $ffab: 00        
            brk                ; $ffac: 00        
            brk                ; $ffad: 00        
            brk                ; $ffae: 00        
            brk                ; $ffaf: 00        
            brk                ; $ffb0: 00        
            brk                ; $ffb1: 00        
            brk                ; $ffb2: 00        
            brk                ; $ffb3: 00        
            brk                ; $ffb4: 00        
            brk                ; $ffb5: 00        
            brk                ; $ffb6: 00        
            brk                ; $ffb7: 00        
            brk                ; $ffb8: 00        
            brk                ; $ffb9: 00        
            brk                ; $ffba: 00        
            brk                ; $ffbb: 00        
            brk                ; $ffbc: 00        
            brk                ; $ffbd: 00        
            brk                ; $ffbe: 00        
            brk                ; $ffbf: 00        
            brk                ; $ffc0: 00        
            brk                ; $ffc1: 00        
            brk                ; $ffc2: 00        
            brk                ; $ffc3: 00        
            brk                ; $ffc4: 00        
            brk                ; $ffc5: 00        
            brk                ; $ffc6: 00        
            brk                ; $ffc7: 00        
            brk                ; $ffc8: 00        
            brk                ; $ffc9: 00        
            brk                ; $ffca: 00        
            brk                ; $ffcb: 00        
            brk                ; $ffcc: 00        
            brk                ; $ffcd: 00        
            brk                ; $ffce: 00        
            brk                ; $ffcf: 00        
            brk                ; $ffd0: 00        
            brk                ; $ffd1: 00        
            brk                ; $ffd2: 00        
            brk                ; $ffd3: 00        
            brk                ; $ffd4: 00        
            brk                ; $ffd5: 00        
            brk                ; $ffd6: 00        
            brk                ; $ffd7: 00        
            brk                ; $ffd8: 00        
            brk                ; $ffd9: 00        
            brk                ; $ffda: 00        
            brk                ; $ffdb: 00        
            brk                ; $ffdc: 00        
            brk                ; $ffdd: 00        
            brk                ; $ffde: 00        
            brk                ; $ffdf: 00        
            brk                ; $ffe0: 00        
            brk                ; $ffe1: 00        
            brk                ; $ffe2: 00        
            brk                ; $ffe3: 00        
            brk                ; $ffe4: 00        
            brk                ; $ffe5: 00        
            brk                ; $ffe6: 00        
            brk                ; $ffe7: 00        
            brk                ; $ffe8: 00        
            brk                ; $ffe9: 00        
            brk                ; $ffea: 00        
            brk                ; $ffeb: 00        
            brk                ; $ffec: 00        
            brk                ; $ffed: 00        
            brk                ; $ffee: 00        
            brk                ; $ffef: 00        
            brk                ; $fff0: 00        
            brk                ; $fff1: 00        
            brk                ; $fff2: 00        
            brk                ; $fff3: 00        
            brk                ; $fff4: 00        
            brk                ; $fff5: 00        
            brk                ; $fff6: 00        
            brk                ; $fff7: 00        
            brk                ; $fff8: 00        
            brk                ; $fff9: 00        

;-------------------------------------------------------------------------------
; Vector Table
;-------------------------------------------------------------------------------
vectors:    .dw nmi                        ; $fffa: af c5     Vector table
            .dw reset                      ; $fffc: 04 c0     Vector table
            .dw irq                        ; $fffe: f4 c5     Vector table
