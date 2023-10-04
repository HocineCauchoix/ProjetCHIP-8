package main

import (
	"fmt"
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
}

func (c *Chip8) loadProgram(rom []byte) error {
	if len(rom) > len(c.memory)-0x200 {
		return fmt.Errorf("ROM is too large to fit in memory")
	}

	copy(c.memory[0x200:], rom)
	c.programCounter = 0x200

	return nil
}

func (c *Chip8) fetchOpcode() {
	opcode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
	switch opcode & 0xF000 {

	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00EE:
			if c.stackPointer == 0 {
				log.Panic("Attempted to return from an empty stack")
			}
			c.stackPointer--
			c.programCounter = c.stack[c.stackPointer]
		case 0x00E0:
			for i := range c.screen {
				c.screen[i] = 0
			}
			c.programCounter += 2
		default:
			fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case OpcodeFX29:
			x := uint16((opcode & 0x0F00) >> 8)
			character := c.command[x]
			c.indexRegister = uint16(character) * 5
			c.programCounter += 2
		case OpcodeFX33:
			x := uint16((opcode & 0x0F00) >> 8)
			value := c.command[x]
			c.memory[c.indexRegister] = value / 100
			c.memory[c.indexRegister+1] = (value / 10) % 10
			c.memory[c.indexRegister+2] = value % 10
			c.programCounter += 2

		default:
			fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	case 0xD000:
		println("case draw")
		x := uint16((opcode & 0x0F00) >> 8)
		y := uint16((opcode & 0x00F0) >> 4)
		height := uint16(opcode & 0x000F)

		for row := uint16(0); row < height; row++ {
			spriteData := c.memory[c.indexRegister+row]
			for col := uint16(0); col < 8; col++ {
				pixel := (spriteData >> (7 - col)) & 0x1
				index := x + row + (y + col*64)

				if pixel == 1 {

					c.screen[index] = 1
				} else {
					c.screen[index] = 0
				}
				println(c.screen[index])

			}
		}

	}
}

type Game struct {
	chip8 Chip8
}

func (g *Game) Update() error {
	for i := 0; i <= 9; i++ {
		if ebiten.IsKeyPressed(ebiten.Key(byte(i) + byte(ebiten.Key0))) {
			g.chip8.input[byte(i)] = true
		} else {
			g.chip8.input[byte(i)] = false
		}
	}

	g.chip8.fetchOpcode()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP8 EMULATOR")

	chip8 := &Chip8{}

	romFileName := "game/1-chip8-logo.ch8"
	romData, err := os.ReadFile(romFileName)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier ROM : %v", err)
	}

	if err := chip8.loadProgram(romData); err != nil {
		log.Fatalf("Erreur lors du chargement de la ROM : %v", err)
	}

	if err := ebiten.RunGame(&Game{chip8: *chip8}); err != nil {
		log.Fatal(err)
	}
}
