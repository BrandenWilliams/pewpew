package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Bullet represents a single shot
type Bullet struct {
	x, y float64
}

// Game holds the player, bullets, and player position
type Game struct {
	playerImage *ebiten.Image
	bullets     []Bullet
	x, y        float64
}

// Update handles movement and shooting
func (g *Game) Update() error {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.x -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.x += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.y += 2
	}

	// Shooting (Spacebar)
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// Fire a bullet from the player's position
		g.bullets = append(g.bullets, Bullet{x: g.x + 128, y: g.y + 63})
	}

	// Update bullets
	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].x += 5 // Move bullet up
	}

	// Remove off-screen bullets
	newBullets := g.bullets[:0]
	for _, b := range g.bullets {
		if b.y > 0 {
			newBullets = append(newBullets, b)
		}
	}
	g.bullets = newBullets

	return nil
}

// Draw the player and bullets
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw player sprite
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.x, g.y)
	screen.DrawImage(g.playerImage, op)

	// Draw bullets
	for _, b := range g.bullets {
		ebitenutil.DrawRect(screen, b.x, b.y, 4, 2, color.RGBA{255, 255, 0, 255})
	}
}

// Layout sets the game window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// Load player sprite
	playerImg, _, err := ebitenutil.NewImageFromFile("art/testship.png")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		playerImage: playerImg,
		x:           320,
		y:           240,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Go Ebitengine Shooting Demo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
