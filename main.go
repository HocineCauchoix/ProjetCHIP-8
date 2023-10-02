package main

import (
	
	"os"
	"log"
	"os/exec"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// Constantes pour la taille de l'écran et de la mémoire, etc.
	memorySize   = 4096
	commandCount = 16
	stackSize    = 16
	screenSize   = 64 * 32
)

// Structure du Chip8
type Chip8 struct {
	// Les champs de la structure Chip8 que vous avez déjà définis
}



// Fonction qui simule


type ClickableArea struct {
	X, Y, Width, Height int
}


type Game struct {
	chip8          Chip8
	gameState      string // Peut être "menu", "jeu", "quitter", etc.
	romLoaded      bool   // Indique si une ROM est chargée
	selectedOption int    // Option de menu sélectionnée (0: Menu, 1: Jouer, 2: Quitter)
	selectedColor  color.RGBA // Couleur pour le texte sélectionné
	menuOptions     []string
    clickableAreas []ClickableArea
}


func (g *Game) startNewGame() {
    if !g.romLoaded {
        // Vérifiez si le fichier jeu.go (ou jeu.ch8) existe dans le répertoire actuel.
        if _, err := os.Stat("./game/Tetris.ch8"); err == nil {
            // Exécutez le jeu Go en tant que processus distinct.
            cmd := exec.Command("go", "run", "./game/Tetris.ch8")
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr

            if err := cmd.Start(); err != nil {
                log.Fatal(err)
            }

            // Indiquez que le jeu est chargé et en cours d'exécution.
            g.romLoaded = true
        } else {
            log.Fatal("Le fichier jeu.go n'a pas été trouvé.")
        }
    }
}

func (g *Game) Update() error {
    if ebiten.IsKeyPressed(ebiten.KeyEscape) {
        // L'utilisateur appuie sur Échap pour revenir au menu.
        g.gameState = "menu"
    }

    if g.gameState == "menu" {
        if ebiten.IsKeyPressed(ebiten.KeyDown) {
            g.selectedOption = (g.selectedOption + 1) % 3
        } else if ebiten.IsKeyPressed(ebiten.KeyUp) {
            g.selectedOption = (g.selectedOption + 2) % 3
        } else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
            // L'utilisateur appuie sur Entrée pour sélectionner une option.
            if g.selectedOption == 0 {
                // Option "Menu" sélectionnée : Revenir au menu principal.
                g.gameState = "menu"
            } else if g.selectedOption == 1 {
                // Option "Jouer" sélectionnée : Lancer un nouveau jeu dans une nouvelle fenêtre.
                g.startNewGame()
            } else if g.selectedOption == 2 {
                // Option "Quitter" sélectionnée : Quitter le jeu.
                os.Exit(0)
            }
        }
    }

    // Gérez d'autres événements de clavier et la logique du jeu ici.

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Effacez l'écran en utilisant la couleur noire.
	screen.Fill(color.Black)

	switch g.gameState {
	case "menu":
		// Affichez "Hello, World!" en blanc en haut de l'écran.
		ebitenutil.DebugPrint(screen, "Hello, World!")

		// Affichez les options de menu ("Menu", "Jouer", "Quitter") centrées en blanc.
		menuOptions := []string{"Menu", "Jouer", "Quitter"}

		// Obtenez la largeur et la hauteur de l'écran.
		screenWidth, screenHeight := screen.Size()

		// Calculez la position verticale de départ pour centrer les options.
		startY := (screenHeight - len(menuOptions)*40) / 2

		// Définissez la taille de la police pour le centrage.
		fontSize := 40

		for i, option := range menuOptions {
			// Calculez la position horizontale pour centrer le texte.
			textWidth := len(option) * fontSize / 2
			x := (screenWidth - textWidth) / 2

			// Calculez la position verticale en fonction de la position de départ.
			y := startY + i*45

			ebitenutil.DebugPrintAt(screen, option, x, y)
		}

	case "jeu":
		// Affichez "Hello, World!" en blanc en haut de l'écran.
		ebitenutil.DebugPrint(screen, "Hello, World!")

		// Le reste de votre code pour dessiner le jeu va ici.

	default:
		// Affichez "Hello, World!" en blanc en haut de l'écran.
		ebitenutil.DebugPrint(screen, "Hello, World!")
	}
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("CHIP8 EMULATOR") // Titre de la fenêtre

	game := &Game{
		chip8:          Chip8{},
		gameState:      "menu",
		romLoaded:      false,
		selectedOption: 0,
	}

	// Exécutez le jeu en utilisant ebiten.RunGame.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}