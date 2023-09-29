package main

import (
	"fmt"
	"io/ioutil"

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
)

// structure du Chip8
type Chip8 struct {
	memory  [memorySize]byte   //mémoire du Chip8
	command [commandCount]byte // cmd V0 à VF
	//uint16 = ints entiers non negatifs de 16 bits
	indexRegister  uint16            // Opérations pour la mémoire
	programCounter uint16            // Pointeur du programme en execution
	stack          [stackSize]uint16 // prise en compte du flow de l'execution du programme
	stackPointer   byte              //Prise en compte du niveau de stack
	delayTimer     byte              // temps de délai entre chaque évènement
	screen         [screenSize]byte  // Taille de l'écran du Chip8
}

// fonction de chargement d'un programme Chip8 (ROM)
func (c *Chip8) loadProgram(filename string) error { //instantie programme par le biais d'un autre programme nommé par filename
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	} // recherche du fichier, si erreur retourner "erreur"

	for i, b := range data {
		c.memory[i+0x200] = b
	} // Si fichier trouvé, sauvegarder dans emplacement de mémoire: 0x200
	return nil // Ensuite retourner 'nil' si chargement effectué
}

// Fonction qui simule
func (c *Chip8) emulateCycle() {
	// Recherche de l'opcode dans la mémoire
	opcode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

	// Décodage et execution de l'opcode (Si plus de codes, ajouter cases)
	switch opcode & 0xF000 {
	case 0x0000:
		// Implementation 0x0NNN, 0x00E0, 0x00EE opcodes ici
	case 0x1000:
		// Implementation 0x1NNN opcode ici
	// Ajouter des cas pour plus d'opcodes
	default:
		fmt.Printf("Unknown opcode: 0x%X\n", opcode)
	}

	// Maj du timer (si nécessaire)

	// Passer à prochaine instruction
	c.programCounter += 2
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "") //display on screen
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
	if err := chip8.loadProgram("game/Tetris.ch8"); err != nil {
		log.Fatalf("Erreur lors du chargement de la ROM : %v", err)
	}

	// Lancez l'émulateur
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

