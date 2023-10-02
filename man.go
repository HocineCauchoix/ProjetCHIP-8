// package main

// import (
// 	"fmt"
// 	"image"
// 	"io/ioutil"
// 	"sync"
// 	"time"

// 	"log"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/inpututil"
// )

// // variables constantes
// const (
// 	//4 KB de RAM
// 	memorySize = 4096
// 	//Nb de Commandes = 16
// 	commandCount = 16
// 	//Taille de Stack
// 	stackSize = 16
// 	//Dimension de l'écran
// 	screenSize = 64 * 32
// )

// // structure du Chip8
// type Chip8 struct {
// 	memory  [memorySize]byte   //mémoire du Chip8
// 	command [commandCount]byte // cmd V0 à VF
// 	//uint16 = ints entiers non negatifs de 16 bits
// 	indexRegister  uint16            // Opérations pour la mémoire
// 	programCounter uint16            // Pointeur du programme en execution
// 	stack          [stackSize]uint16 // prise en compte du flow de l'execution du programme
// 	stackPointer   byte              //Prise en compte du niveau de stack
// 	delayTimer     byte              // temps de délai entre chaque évènement
// 	pc             uint16            // Program counter, set it to the initial memory offset
// 	stackFrame     int               // current stack frame. Starts at -1 and is set to 0 on first use
// 	I              uint16            // represents Index register
// 	screen         [screenSize]byte  // Taille de l'écran du Chip8
// 	delayTimerLock sync.Mutex        // lock for incrementing/setting/accessing the delay timer
// 	soundTimerLock sync.Mutex        // lock for incrementing/setting/accessing the sound timer
// }

// type Runtime struct {
// 	keys  []ebiten.Key // stores currently pressed keys
// 	image *image.RGBA  // screen buffer
// 	lock  sync.Mutex   // lock to protect image from concurrent read/writes.
// }

// var input = 0b01111110                     // same as hexadecimal 0x7E or decimal 126. Key is to look at the bits. 0111 1110
// var rightmostBitsMasked = input & 0xF0     // F0 == 11110000, by ANDing only the bits on the left are kept. Result => 01110000
// var leftmostBitsMasked = input & 0x0F      // 0F == 00001111, like above, only the bits on the right are kept. Result => 00001110
// var firstNibble = rightmostBitsMasked >> 4 // shift 4 steps to the right. Result => 00000111.
// var secondNibble = leftmostBitsMasked      // no need to shift here, the bits are already in the "rightmost" position.

// func (e *Chip8) startDelayTimer() {
// 	var tick = 1000 / 60
// 	for {
// 		time.Sleep(time.Millisecond * time.Duration(tick)) // sleep for approx 16ms
// 		e.delayTimerLock.Lock()                            // obtain lock
// 		if e.delayTimer > 0 {                              // if above 0, decrement by 1
// 			e.delayTimer--
// 		}
// 		e.delayTimerLock.Unlock() // release lock
// 	}
// }

// // fonction de chargement d'un programme Chip8 (ROM)
// func (c *Chip8) loadProgram(filename string) error { //instantie programme par le biais d'un autre programme nommé par filename
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return err
// 	} // recherche du fichier, si erreur retourner "erreur"

// 	for i, b := range data {
// 		c.memory[i+0x200] = b
// 	} // Si fichier trouvé, sauvegarder dans emplacement de mémoire: 0x200
// 	return nil // Ensuite retourner 'nil' si chargement effectué
// }

// // Fonction qui simule
// func (c *Chip8) emulateCycle() {
// 	// Recherche de l'opcode dans la mémoire
// 	opcode := uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

// 	// Décodage et execution de l'opcode (Si plus de codes, ajouter cases)
// 	switch opcode & 0xF000 {
// 	case 0x0000:
// 		// Implementation 0x0NNN, 0x00E0, 0x00EE opcodes ici
// 	case 0x1000:
// 		// Implementation 0x1NNN opcode ici
// 	// Ajouter des cas pour plus d'opcodes
// 	default:
// 		fmt.Printf("Unknown opcode: 0x%X\n", opcode)
// 	}

// 	// Maj du timer (si nécessaire)

// 	// Passer à prochaine instruction
// 	c.programCounter += 2
// }

// type Game struct{

// // Draw implements ebiten.Game.
// func (*Game) Draw(screen *v2.Image) {
// 	panic("unimplemented")
// }

// func (g *Game) Update() error {
// 	return nil
// }

// func (g *Runtime) Update() error {
// 	// store currently pressed keys, reusing the last update's slice of pressed keys.
// 	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
// 	return nil
// }

// func (g *Runtime) Draw(screen *ebiten.Image) {
// 	// Protected by the lock, call the screen's "WritePixels"
// 	// with the current []byte in g.image.Pix
// 	g.lock.Lock()
// 	screen.WritePixels(g.image.Pix)
// 	g.lock.Unlock()
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
// 	return 320, 240
// }
// }

// func main() {
// 	ebiten.SetWindowSize(640, 320)
// 	ebiten.SetWindowTitle("CHIP8 EMULATOR") // titre de la fenetre
// 	if err := ebiten.RunGame(&Game{}); err != nil {
// 		log.Fatal(err)
// 	}
// }