package main

import (
	"log"
	"os"

	chip8 "ProjetCHIP-8/chip8"
	game "ProjetCHIP-8/chip8"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP8 EMULATOR")

	chip8 := &chip8.Chip8{}

	romFileName := "game/invaders.ch8"
	romData, err := os.ReadFile(romFileName)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier ROM : %v", err)
	}

	if err := chip8.LoadProgram(romData); err != nil {
		log.Fatalf("Erreur lors du chargement de la ROM : %v", err)
	}

	game := &game.Game{Chip8: *chip8}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
