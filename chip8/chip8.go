package chip8

import (
	"fmt"
	"log"
)

const ( //variables constantes
	memorySize   = 4096
	commandCount = 16
	stackSize    = 16
	screenSize   = 64 * 32
	OpcodeFX29   = 0xF029
	OpcodeFX33   = 0xF033
	OpcodeDXYN   = 0xD000
)

type Chip8 struct { //Structure du Chip8
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

func randomByte() uint8 {
	return 0
}

func (c *Chip8) executeOpcode() bool {
	fmt.Println("OPCODE : ", c.Opcode&0xF000) // Affiche le code de l'opcode en cours
	switch c.Opcode & 0xF000 {                //prend les 4 bits les plus importants de l'opcode
	case 0x0000:
		switch c.Opcode & 0x000FF { //extrait les 8 bits les moins importants de l'opcode
		case 0x00E0:
			// Opcode pour effacer l'écran (CLS)
			c.screen = [screenSize]byte{}
			c.programCounter += 2
			println("CLS")
		case 0x000E:
			// Opcode pour le retour d'appel à partir de la sous-routine (RET)
			if c.stackPointer == 0 {
				log.Panic("Attempted to return from an empty stack")
			}
			c.stackPointer--
			c.programCounter = c.stack[c.stackPointer]
			c.programCounter += 2
		default:

		}
	case 0x1000:
		// Opcode pour définir le compteur de programme à une adresse spécifique (JP addr)
		c.programCounter = c.Opcode & 0x0FFF
	case 0x2000:
		// Opcode pour appeler une sous-routine (CALL addr)
		c.stack[c.stackPointer] = c.programCounter
		c.stackPointer++
		c.programCounter = c.Opcode & 0x0FFF
	case 0x3000:
		// Opcode pour sauter l'instruction suivante si un registre est égal à une valeur (SE Vx, byte)
		x := uint16((c.Opcode & 0x0F00) >> 8)
		if c.Register[x] == uint8(c.Opcode&0x00FF) {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x4000:
		// Opcode pour sauter l'instruction suivante si un registre n'est pas égal à une valeur (SNE Vx, byte)
		x := uint16((c.Opcode & 0x0F00) >> 8)
		if c.Register[x] != uint8(c.Opcode&0x00FF) {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x5000:
		// Opcode de comparaison entre deux registres.
		x := uint16((c.Opcode & 0x0F00) >> 8)
		y := uint16((c.Opcode & 0x00F0) >> 4)
		if c.Register[x] == c.Register[y] {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0x6000:
		// Opcode de chargement d'une valeur immédiate dans un registre.
		x := uint16((c.Opcode & 0x0F00) >> 8)
		c.Register[x] = uint8(c.Opcode & 0x00FF)
		println("LD V")
		c.programCounter += 2
	case 0x7000:
		// Opcode d'addition avec une valeur immédiate.
		x := uint16((c.Opcode & 0x0F00) >> 8)
		c.Register[x] += uint8(c.Opcode & 0x00FF)
		c.programCounter += 2
	case 0x8000:
		// Opcode de manipulation de registres.
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
			// Opcode d'addition avec retenue.
			if c.Register[y] > (0xFF - c.Register[x]) {
				c.Register[0xF] = 1
			} else {
				c.Register[0xF] = 0
			}
			c.Register[x] += c.Register[y]
			c.programCounter += 2
		case 0x0005:
			// Opcode de soustraction avec emprunt.
			if c.Register[y] > c.Register[x] {
				c.Register[0xF] = 0
			} else {
				c.Register[0xF] = 1
			}
			c.Register[x] -= c.Register[y]
			c.programCounter += 2
		case 0x0006:
			// Opcode de décalage à droite.
			c.Register[0xF] = c.Register[x] & 0x1
			c.Register[x] >>= 1
			c.programCounter += 2
		case 0x0007:
			// Opcode de différence entre Vy et Vx
			if c.Register[x] > c.Register[y] {
				c.Register[0xF] = 0
			} else {
				c.Register[0xF] = 1
			}
			c.Register[x] = c.Register[y] - c.Register[x]
			c.programCounter += 2
		case 0x000E:
			// Opcode de décalage à gauche
			c.Register[0xF] = c.Register[x] >> 7
			c.Register[x] <<= 1
			c.programCounter += 2
		default:
			panicUnknownOpcode(c.Opcode)
		}
	case 0x9000:
		// Opcode de saut si les registres Vx et Vy ne sont pas égaux
		x := uint16((c.Opcode & 0x0F00) >> 8)
		y := uint16((c.Opcode & 0x00F0) >> 4)
		if c.Register[x] != c.Register[y] {
			c.programCounter += 4
		} else {
			c.programCounter += 2
		}
	case 0xA000:
		// Opcode de chargement de la valeur immédiate dans le registre I
		c.indexRegister = c.Opcode & 0x0FFF
		println("LD I")
		c.programCounter += 2
	case 0xB000:
		// Opcode de saut à une adresse en ajoutant la valeur de V0
		c.programCounter = (c.Opcode & 0x0FFF) + uint16(c.Register[0])
	case 0xC000:
		// Opcode de génération d'un nombre aléatoire et de stockage dans Vx
		x := uint16((c.Opcode & 0x0F00) >> 8)
		randomValue := randomByte()
		c.Register[x] = randomValue & uint8(c.Opcode&0x00FF)
		c.programCounter += 2
	case 0xD000:
		// Opcode de dessin à l'écran
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
			// Opcode de saut conditionnel si la touche pressée correspond à la valeur dans Vx
			if c.input[c.Register[x]] {
				c.programCounter += 4
			} else {
				c.programCounter += 2
			}
		case 0x00A1:
			// Opcode de saut conditionnel si la touche pressée ne correspond pas à la valeur dans Vx
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
			// Opcode de chargement de la valeur du retard dans Vx
			c.Register[x] = c.delayTimer
			c.programCounter += 2
		case 0x000A:
			// Opcode d'attente d'une touche et de stockage de sa valeur dans Vx
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
			// Opcode de chargement de la valeur de Vx dans le retard
			c.delayTimer = c.Register[x]
			c.programCounter += 2
		case 0x0018:
			// Opcode de chargement de la valeur de Vx dans le timer sonore
			c.soundTimer = c.Register[x]
			c.programCounter += 2
		case 0x001E:
			// Opcode d'ajout de la valeur de Vx à I
			if c.indexRegister+uint16(c.Register[x]) > 0xFFF {
				c.Register[0xF] = 1
			} else {
				c.Register[0xF] = 0
			}
			c.indexRegister += uint16(c.Register[x])
			c.programCounter += 2
		case 0x0029:
			// Opcode de chargement de l'emplacement du caractère dans I
			c.indexRegister = uint16(c.Register[x]) * 0x5
			c.programCounter += 2
		case 0x0033:
			// Opcode de chargement des chiffres décimaux dans la mémoire
			c.memory[c.indexRegister] = c.Register[x] / 100
			c.memory[c.indexRegister+1] = (c.Register[x] / 10) % 10
			c.memory[c.indexRegister+2] = c.Register[x] % 10
			c.programCounter += 2
		case 0x0055:
			// Opcode de sauvegarde des registres V0 à Vx dans la mémoire
			for i := uint16(0); i <= x; i++ {
				c.memory[c.indexRegister+i] = c.Register[i]
			}
			c.indexRegister += x + 1
			c.programCounter += 2
		case 0x0065:
			// Opcode de chargement des registres V0 à Vx à partir de la mémoire
			for i := uint16(0); i <= x; i++ {
				c.Register[i] = c.memory[c.indexRegister+i]
			}
			c.indexRegister += x + 1
			c.programCounter += 2
		default:
			panicUnknownOpcode(c.Opcode) // Si l'opcode n'est pas reconnu, génère une erreur
		}
	default:
		panicUnknownOpcode(c.Opcode) // Si l'opcode n'est pas reconnu, génère une erreur
	}
	return false // Indique que l'exécution normale du programme doit se poursuivre
}

func panicUnknownOpcode(opcode uint16) {
	log.Panicf("Unknown opcode %v", opcode)
}

func (c *Chip8) LoadProgram(rom []byte) error { //chargement programme de la mémoire
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
