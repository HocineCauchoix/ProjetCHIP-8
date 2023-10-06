package chip8

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Chip8 Chip8
}

func (g *Game) Update() error { //Mise à jour du jeu

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

func (g *Game) Draw(screen *ebiten.Image) { //Affichage de l'écran
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) { //Dimension de la fenêtre
	return 64, 32
}
