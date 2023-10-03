package main

import (
	"fmt"
	"os"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// variables constantes
const (
	//4 KB de RAM
	memorySize = 4096
	//Nb de Commandes = 16
	commandCount = 16
	//Taille de Stack
	stackSize = 16
	//Dimension de l'écran
	screenSize = 64 * 32
	OpcodeFX29 = 0xF029
	OpcodeFX33 = 0xF033
	OpcodeDXYN = 0xD000
)

// structure du Chip8
type Chip8 struct {
	input       [16]bool // État des touches du clavier (0-F)     // Temporisation du délai
	soundTimer  byte     // Temporisation du son
	keyPressed  bool     // Un drapeau pour indiquer si une touche est pressée
	keyRegister uint16
	memory      [memorySize]byte   //mémoire du Chip8
	command     [commandCount]byte // cmd V0 à VF
	//uint16 = ints entiers non negatifs de 16 bits
	indexRegister  uint16            // Opérations pour la mémoire
	programCounter uint16            // Pointeur du programme en execution
	stack          [stackSize]uint16 // prise en compte du flow de l'execution du programme
	stackPointer   byte              //Prise en compte du niveau de stack
	delayTimer     byte              // temps de délai entre chaque évènement // Taille de l'écran du Chip8
	fontset        [80]byte          // Tableau de police (pour stocker les chiffres 0-9 et certaines lettres)
	screen         [64 * 32]bool     // Écran représenté par un tableau de booléens
	pixelUpdate    bool              // Un drapeau pour indiquer si l'écran a été modifié

}

// fonction de chargement d'un programme Chip8 (ROM)
func (c *Chip8) loadProgram() error {
	filename := "game/Tetris.ch8" // Chemin de la ROM
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	for i, b := range data {
		c.memory[i+0x200] = b
	}
	c.fontset = [80]byte{
		// Les valeurs de police pour 0-9
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9

		// Les valeurs de police pour A-F (10-15)
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	return nil
}
func (c *Chip8) drawSprite(x, y, n uint16) {
	// Parcourez chaque ligne du sprite (n lignes)
	for row := uint16(0); row < n; row++ {
		spriteRow := c.memory[c.indexRegister+row]

		// Parcourez chaque pixel de la ligne (8 pixels par byte)
		for col := uint16(0); col < 8; col++ {
			// Vérifiez chaque bit de la ligne du sprite
			if (spriteRow & (0x80 >> col)) != 0 {
				// Calculez les coordonnées réelles (modulo pour gérer le débordement)
				realX := (x + col) % 64
				realY := (y + row) % 32

				// Mettez à jour l'état du pixel
				c.screen[realY*64+realX] = !c.screen[realY*64+realX]
			}
		}
	}
}

func (c *Chip8) emulateCycle() {
	// Recherche de l'opcode dans la mémoire
	opcode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

	// Décodage et execution de l'opcode (Si plus de codes, ajouter des cas)
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00E0:
			// Opcode 0x00E0 : Efface l'écran
			// Implémentez le code pour effacer l'écran ici
			// c.screen doit être effacé
			c.programCounter += 2
		case 0x00EE:
			// Opcode 0x00EE : Retour d'une sous-routine
			// Implémentez le code pour revenir d'une sous-routine ici
			// Vous devriez dépiler l'adresse de retour depuis la pile
			// c.programCounter doit être mis à cette adresse

			// Opérateur FX29 : Met le pointeur de police I à l'emplacement du caractère stocké dans VX.
		case OpcodeFX29:
			x := uint16((opcode & 0x0F00) >> 8)
			character := c.command[x]
			c.indexRegister = uint16(character) * 5 // Chaque caractère a une largeur de 5 pixels
			c.programCounter += 2

		// Opérateur FX33 : Stocke la représentation binaire-décimale de VX à l'adresse I, I+1 et I+2.
		case OpcodeFX33:
			x := uint16((opcode & 0x0F00) >> 8)
			value := c.command[x]
			c.memory[c.indexRegister] = value / 100         // Centaines
			c.memory[c.indexRegister+1] = (value / 10) % 10 // Dizaines
			c.memory[c.indexRegister+2] = value % 10        // Unités
			c.programCounter += 2

		// Opérateur DXYN : Dessine un sprite en utilisant les données stockées à l'adresse I à partir de (VX, VY), avec une hauteur de N pixels.
		case OpcodeDXYN:
			x := uint16((opcode & 0x0F00) >> 8)
			y := uint16((opcode & 0x00F0) >> 4)
			height := uint16(opcode & 0x000F)

			// Vous devrez implémenter la logique pour dessiner le sprite ici,
			// en utilisant les données à l'adresse I et en mettant à jour l'écran (c.screen).
			// N'oubliez pas de gérer les collisions et de définir le drapeau pixelUpdate en conséquence.

			c.pixelUpdate = false

			// Dessiner le sprite ligne par ligne
			for row := uint16(0); row < height; row++ {
				spriteData := c.memory[c.indexRegister+row]

				// Dessiner chaque pixel du sprite
				for col := uint16(0); col < 8; col++ {
					// Utilisez les opérations de bit pour extraire chaque pixel du sprite
					pixel := (spriteData >> (7 - col)) & 0x1
					posX := (x + col) % 64
					posY := (y + row) % 32

					// Si le pixel est allumé (1) et le pixel sur l'écran est déjà allumé (1),
					// définissez le drapeau de collision (c.pixelUpdate)
					if pixel == 1 && c.screen[posY*64+posX] {
						c.pixelUpdate = true
					}

					// Effectuer une opération XOR pour dessiner le pixel (gestion de collision)
					c.screen[posY*64+posX] = c.screen[posY*64+posX] != (pixel == 1)
				}
			}

		case 0xF00A:
			x := uint16((opcode & 0x0F00) >> 8)

			// Vérifiez si une touche est pressée
			keyPressed := false
			for i := 0; i < len(c.input); i++ {
				if c.input[i] {
					c.command[x] = byte(i)
					keyPressed = true
					break
				}
			}

			if !keyPressed {
				// Si aucune touche n'est pressée, réexécutez cet opcode
				c.programCounter -= 2
			}

		case 0xF015:
			x := uint16((opcode & 0x0F00) >> 8)
			c.delayTimer = c.command[x]
			c.programCounter += 2

		case 0xF018:
			x := uint16((opcode & 0x0F00) >> 8)
			c.soundTimer = c.command[x]
			c.programCounter += 2

			c.programCounter += 2

			c.programCounter += 2
		default:
			fmt.Printf("Unknown opcode: 0x%X\n", opcode)
		}
	}
}

type Game struct {
	chip8 Chip8
}

func (g *Game) Update() error {
	// Gérer les événements de touche (exemple avec les touches 0-9)
	if ebiten.IsKeyPressed(ebiten.Key0) {
		g.chip8.input[0x0] = true
	} else {
		g.chip8.input[0x0] = false
	}
	// Répétez ce processus pour les autres touches (0-9)...
	// Par exemple, pour la touche 1 : g.chip8.Input[0x1] = ebiten.IsKeyPressed(ebiten.Key1)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello World") //display on screen
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP8 EMULATOR") // titre de la fenetre

	// Créez une instance de Chip8
	chip8 := &Chip8{}

	// Chargez la ROM du jeu
	if err := chip8.loadProgram(); err != nil {
		log.Fatalf("Erreur lors du chargement de la ROM : %v", err)
	}

	// Lancez l'émulateur
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
