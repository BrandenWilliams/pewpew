package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Bullet represents a single shot
type Bullet struct {
	x, y float64
}

type Enemy struct {
	x, y float64
}

// Game holds the player, bullets, and player position
type Game struct {
	playerImage     *ebiten.Image
	enemyImage      *ebiten.Image
	bullets         []Bullet
	enemies         []Enemy
	x, y            float64
	firePressed     bool // Track fire key state
	enemySpawnTimer float64
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

	// **Shooting (Spacebar, one shot per press)**
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.firePressed {
		g.bullets = append(g.bullets, Bullet{x: g.x + 122, y: g.y + 62})
	}

	// **Update bullet positions (keep this separate!)**
	for i := range g.bullets {
		g.bullets[i].x += 5 // Move bullets up
	}

	newBullets := []Bullet{}
	newEnemies := []Enemy{}

	for _, b := range g.bullets {
		hit := false

		// Check each enemy against the bullet
		for i := 0; i < len(g.enemies); i++ {
			e := g.enemies[i]

			// Check collision (adjust sizes to match sprite)
			if b.x < e.x+20 && b.x+4 > e.x && b.y < e.y+20 && b.y+2 > e.y {
				hit = true                                            // Bullet hits this enemy
				g.enemies = append(g.enemies[:i], g.enemies[i+1:]...) // Remove the enemy
				break
			}
		}

		// If the bullet didn't hit anything, keep it
		if !hit {
			newBullets = append(newBullets, b)
		}
	}

	g.bullets = newBullets

	// **Keep only enemies that weren't hit**
	for _, e := range g.enemies {
		enemyHit := false
		for _, b := range g.bullets {
			if b.x+4 > e.x && b.x < e.x+20 && b.y+2 > e.y && b.y < e.y+20 {
				enemyHit = true
				break
			}
		}
		if !enemyHit {
			newEnemies = append(newEnemies, e)
		}
	}
	g.enemies = newEnemies

	// **Remove off-screen bullets**
	newBullets = g.bullets[:0]
	for _, b := range g.bullets {
		if b.y > 0 {
			newBullets = append(newBullets, b)
		}
	}
	g.bullets = newBullets

	// **Track the Space key state**
	g.firePressed = ebiten.IsKeyPressed(ebiten.KeySpace)

	// Track the current ticks per second
	tps := ebiten.ActualTPS()

	// Increment the timer based on elapsed time per frame
	g.enemySpawnTimer += 1 / tps

	// Spawn an enemy every 1.5 seconds
	if g.enemySpawnTimer >= 1.5 {
		g.enemies = append(g.enemies, Enemy{
			y: float64(50 + rand.Intn(480)), // Random Y position
			x: 480,                          // Spawn to the right of the screen
		})
		g.enemySpawnTimer = 0 // Reset the timer
	}

	// Move enemies left
	for i := range g.enemies {
		g.enemies[i].x -= 2 // Adjust speed as needed
	}

	// Remove enemies that move off the screen
	newEnemies = g.enemies[:0]
	for _, e := range g.enemies {
		if e.x < 480 { // Assuming screen height is 480
			newEnemies = append(newEnemies, e)
		}
	}
	g.enemies = newEnemies

	return nil
}

// Draw the player and bullets
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw player sprite
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(g.x, g.y)
	screen.DrawImage(g.playerImage, playerOp)

	// Draw enemy sprite
	for _, e := range g.enemies {
		enemyOp := &ebiten.DrawImageOptions{}
		enemyOp.GeoM.Translate(e.x, e.y)
		screen.DrawImage(g.enemyImage, enemyOp)
	}

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
	// Load sprites

	// player sprite
	playerImg, _, err := ebitenutil.NewImageFromFile("art/testship.png")
	if err != nil {
		log.Fatal(err)
	}

	// enemy sprite
	enemyImg, _, err := ebitenutil.NewImageFromFile("art/enemyship.png")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		playerImage: playerImg,
		enemyImage:  enemyImg,
		x:           320,
		y:           240,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("pewpew v0.0.1")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
