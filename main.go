package main

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

type Chip8 struct {
	memory  [memorySize]byte   //mémoire du Chip8
	command [commandCount]byte // cmd V0 à VF
	//uint16 = non negative ints of 16 bits
	indexRegister  uint16            // Opérations pour la mémoire
	programCounter uint16            // Pointeur du programme en execution
	stack          [stackSize]uint16 // prise en compte du flow de l'execution du programme
	stackPointer   byte              //Prise en compte du niveau de stack
	delayTimer     byte              // temps de délai entre chaque évènement
	screen         [screenSize]byte  // Taille de l'écran du Chip8
}
