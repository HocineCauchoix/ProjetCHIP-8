package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	memorySize   = 4096
	commandCount = 16
	stackSize    = 16
	screenSize   = 64 * 32
	OpcodeFX29   = 0xF029
	OpcodeFX33   = 0xF033
	OpcodeDXYN   = 0xD000
)

type Chip8 struct {
	Register       [16]uint8
	input          [16]bool
	soundTimer     byte
	keyPressed     bool
	keyRegister    uint16
	memory         [memorySize]byte
	command        [commandCount]byte
	indexRegister  uint16
	programCounter uint16
	stack          [stackSize]uint16
	stackPointer   byte
	delayTimer     byte
	fontset        [80]byte
	screen         [screenSize]byte
	pixelUpdate    byte
	Opcode         uint16
}

func (c *Chip8) fetchOpcode() {
	c.Opcode = uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
}

func (c *Chip8) executeOpcode() bool {
	fmt.Println("OPCODE : ", c.Opcode&0xF000)
	switch c.Opcode & 0xF000 {
	case 0x0000:
		switch c.Opcode & 0x000FF {
		case 0x00E0:
			c.screen = [screenSize]byte{} // Clear the screen
			c.programCounter += 2
			println("CLS")
		case 0x000E:
			if c.stackPointer == 0 {
				log.Panic("Attempted to return from an empty stack")
			}
			c.stackPointer--
			c.programCounter = c.stack[c.stackPointer]
			c.programCounter += 2
		default:
			// panicUnknownOpcode(c.Opcode)
		}
	case 0x1000:
		c.programCounter = c.Opcode & 0x0FFF
	case 0x2000:
		c.stack[c.stackPointer] = c.programCounter
		c.stackPointer++
		c.programCounter = c.Opcode & 0x0FFF
	case 0x3000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		if c.Register[x] == uint8(c.Opcode&0x00FF) {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x4000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		if c.Register[x] != uint8(c.Opcode&0x00FF) {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x5000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		y := uint16((c.Opcode & 0x00F0) >> 4)
		if c.Register[x] == c.Register[y] {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x6000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		c.Register[x] = uint8(c.Opcode & 0x00FF)
		println("LD V")
		c.programCounter += 2
	// Ajoutez ici les autres opcodes manquants en suivant le modÃ¨le ci-dessus
	case 0x7000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		c.Register[x] += uint8(c.Opcode & 0x00FF)
		c.programCounter += 2
	case 0x8000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		y := uint16((c.Opcode & 0x00F0) >> 4)
		switch c.Opcode & 0x000F {
		case 0x0000:
			c.Register[x] = c.Register[y]
			c.programCounter += 2
		case 0x0001:
			c.Register[x] |= c.Register[y]
			c.programCounter += 2
		case 0x0002:
			c.Register[x] &= c.Register[y]
			c.programCounter += 2
		case 0x0003:
			c.Register[x] ^= c.Register[y]
			c.programCounter += 2
		case 0x0004:
			if c.Register[y] > (0xFF - c.Register[x]) {
				c.Register[0xF] = 1
			} else {
				c.Register[0xF] = 0
			}
			c.Register[x] += c.Register[y]
			c.programCounter += 2
		case 0x0005:
			if c.Register[y] > c.Register[x] {
				c.Register[0xF] = 0
			} else {
				c.Register[0xF] = 1
			}
			c.Register[x] -= c.Register[y]
			c.programCounter += 2
		case 0x0006:
			c.Register[0xF] = c.Register[x] & 0x1
			c.Register[x] >>= 1
			c.programCounter += 2
		case 0x0007:
			if c.Register[x] > c.Register[y] {
				c.Register[0xF] = 0
			} else {
				c.Register[0xF] = 1
			}
			c.Register[x] = c.Register[y] - c.Register[x]
			c.programCounter += 2
		case 0x000E:
			c.Register[0xF] = c.Register[x] >> 7
			c.Register[x] <<= 1
			c.programCounter += 2
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0x9000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		y := uint16((c.Opcode & 0x00F0) >> 4)
		if c.Register[x] != c.Register[y] {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0xA000:
		c.indexRegister = c.Opcode & 0x0FFF
		println("LD I")
		c.programCounter += 2
	case 0xB000:
		c.programCounter = (c.Opcode & 0x0FFF) + uint16(c.Register[0])
	case 0xC000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		randomValue := randomByte()
		c.Register[x] = randomValue & uint8(c.Opcode&0x00FF)
		c.programCounter += 2
	case 0xD000:
		x := uint16(c.Register[(c.Opcode&0x0F00)>>8])
		y := uint16(c.Register[(c.Opcode&0x00F0)>>4])
		height := uint16(c.Opcode & 0x000F)
		var pixel uint16

		c.Register[0xF] = 0
		for yline := uint16(0); yline < height; yline++ {
			pixel = uint16(c.memory[c.indexRegister+yline])
			for xline := uint16(0); xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.screen[x+xline+((y+yline)*64)] == 1 {
						c.Register[0xF] = 1
					}
					c.screen[x+xline+((y+yline)*64)] ^= 1
				}
			}
		}
		println("DRW")
		c.programCounter += 2
	case 0xE000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		switch c.Opcode & 0x00FF {
		case 0x009E:
			if c.input[c.Register[x]] {
				c.programCounter += 4
			} else {
				c.programCounter += 2
			}
		case 0x00A1:
			if !c.input[c.Register[x]] {
				c.programCounter += 4
			} else {
				c.programCounter += 2
			}
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0xF000:
		x := uint16((c.Opcode & 0x0F00) >> 8)
		switch c.Opcode & 0x00FF {
		case 0x0007:
			c.Register[x] = c.delayTimer
			c.programCounter += 2
		case 0x000A:
			keyPress := false
			for i := uint8(0); i < 16; i++ {
				if c.input[i] {
					c.Register[x] = i
					keyPress = true
				}
			}
			if !keyPress {
				return true
			}
			c.programCounter += 2
		case 0x0015:
			c.delayTimer = c.Register[x]
			c.programCounter += 2
		case 0x0018:
			c.soundTimer = c.Register[x]
			c.programCounter += 2
		case 0x001E:
			if c.indexRegister+uint16(c.Register[x]) > 0xFFF {
				c.Register[0xF] = 1
			} else {
				c.Register[0xF] = 0
			}
			c.indexRegister += uint16(c.Register[x])
			c.programCounter += 2
		case 0x0029:
			c.indexRegister = uint16(c.Register[x]) * 0x5
			c.programCounter += 2
		case 0x0033:
			c.memory[c.indexRegister] = c.Register[x] / 100
			c.memory[c.indexRegister+1] = (c.Register[x] / 10) % 10
			c.memory[c.indexRegister+2] = c.Register[x] % 10
			c.programCounter += 2
		case 0x0055:
			for i := uint16(0); i <= x; i++ {
				c.memory[c.indexRegister+i] = c.Register[i]
			}
			c.indexRegister += x + 1
			c.programCounter += 2
		case 0x0065:
			for i := uint16(0); i <= x; i++ {
				c.Register[i] = c.memory[c.indexRegister+i]
			}
			c.indexRegister += x + 1
			c.programCounter += 2
		default:
			panicUnknownOpcode(c.Opcode)
		}
	default:
		panicUnknownOpcode(c.Opcode)
	}
	return false
}

func panicUnknownOpcode(opcode uint16) {
	log.Panicf("Unknown opcode %v", opcode)
}

func (g *Game) Update() error {

	g.Chip8.fetchOpcode()
	g.Chip8.executeOpcode()
	for i := 0; i <= 9; i++ {
		if ebiten.IsKeyPressed(ebiten.Key(byte(i) + byte(ebiten.Key0))) {
			g.Chip8.input[byte(i)] = true
		} else {
			g.Chip8.input[byte(i)] = false
		}
	}

	return nil
}

func randomByte() uint8 {
	return 0
}

type Game struct {
	Chip8 Chip8
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			idx := y*64 + x
			if g.Chip8.screen[idx] == 1 {
				screen.Set(x, y, color.White)
			} else {
				screen.Set(x, y, color.Black)
			}
		}
	}

	ebitenutil.DebugPrint(screen, "")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32
}

func (c *Chip8) loadProgram(rom []byte) error {
	if len(rom) > 0 {
		fmt.Println("lecture...")
	}

	if len(rom) > len(c.memory)-512 {
		return fmt.Errorf("ROM is too large to fit in memory")
	}

	for i := uint16(0); i < uint16(len(rom)); i++ {
		c.memory[512+i] = rom[i]
	}

	if c.memory[512] != 0 {
		println("memory load...")
	}

	// copy(c.memory[0x200:], rom)
	c.programCounter = 0x200

	return nil
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP8 EMULATOR")

	chip8 := &Chip8{}

	romFileName := "game/invaders.ch8"
	romData, err := os.ReadFile(romFileName)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier ROM : %v", err)
	}

	if err := chip8.loadProgram(romData); err != nil {
		log.Fatalf("Erreur lors du chargement de la ROM : %v", err)
	}

	game := &Game{Chip8: *chip8}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
