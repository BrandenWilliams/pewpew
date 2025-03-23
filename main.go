package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Bullet represents the location of a single shot
type Bullet struct {
	x, y float64
}

// location of enemy
type Enemy struct {
	x, y float64
}

// Game holds the player, bullets, and player position
type Game struct {
	hadFirstUpdate  bool
	playerImage     *ebiten.Image
	playerPixels    []byte
	enemyImage      *ebiten.Image
	enemyPixels     []byte
	bullets         []Bullet
	enemies         []Enemy
	x, y            float64
	firePressed     bool // Track fire key state
	enemySpawnTimer float64
}

// update player movement (spaceship)
func (g *Game) UpdatePlayer() {
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
}

func (g *Game) BulletUpdate() {
	// **Shooting (Spacebar, one shot per press)**
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.firePressed {
		g.bullets = append(g.bullets, Bullet{x: g.x + 122, y: g.y + 62})
	}

	// **Update bullet positions (keep this separate!)**
	for i := range g.bullets {
		g.bullets[i].x += 5 // Move bullets up
	}
}

func (g *Game) extractPixels() {
	playerPoint := g.playerImage.Bounds().Size()
	g.playerPixels = make([]byte, 4*playerPoint.X*playerPoint.Y)
	g.playerImage.ReadPixels(g.playerPixels)

	EnemyPoint := g.enemyImage.Bounds().Size()
	g.enemyPixels = make([]byte, 4*EnemyPoint.X*EnemyPoint.Y)
	g.enemyImage.ReadPixels(g.enemyPixels)
}

func pixelsCollide(bulletX, bulletY, bulletWidth, bulletHeight int,
	enemyX, enemyY, enemyWidth, enemyHeight int, imgPixels []byte) bool {

	// Define the overlapping rectangle
	xStart := max(bulletX, enemyX)
	yStart := max(bulletY, enemyY)
	xEnd := min(bulletX+bulletWidth, enemyX+enemyWidth-1)
	yEnd := min(bulletY+bulletHeight, enemyY+enemyHeight-1)

	// If no overlap, skip checks
	if xStart >= xEnd || yStart >= yEnd {
		return false
	}

	// Rewrite / double check the rest
	// Loop through the overlapping pixels only
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			p1 := ((y-enemyY)*enemyWidth + (x - enemyX)) * 4

			// Check alpha channel for solidity (>= 128)
			if imgPixels[p1+3] >= 128 {
				return true
			}
		}
	}

	return false
}

/*
func pixelsCollide(bulletX, bulletY, bulletWidth1, bulletHeight1 int,

		enemyX, enemyY int, pixels2 []byte, enemyW, enemyH int) bool {



		return false
	}
*/
func (g *Game) CollisionCheck(newBullets []Bullet, newEnemies []Enemy) {
	for _, b := range g.bullets {
		hit := false

		for i := 0; i < len(g.enemies); i++ {
			e := g.enemies[i]

			// Get enemy sprite size
			ew, eh := g.enemyImage.Size()

			// First, do a bounding box check for early rejection
			if b.x+4 > e.x && b.x < e.x+float64(ew) &&
				b.y+2 > e.y && b.y < e.y+float64(eh) {

				// Convert float positions to integers for pixel collision check
				bulletX := int(b.x)
				bulletY := int(b.y)
				enemyX := int(e.x)
				enemyY := int(e.y)

				// Call pixel-perfect collision detection
				if pixelsCollide(
					bulletX, bulletY, 4, 2,
					enemyX, enemyY, ew, eh, g.enemyPixels,
				) {
					hit = true
					log.Printf("Pixel-perfect hit detected at Bullet(%v, %v) and Enemy(%v, %v)", b.x, b.y, e.x, e.y)
					g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
					break
				}
			}
		}

		// Keep the bullet if no hit occurred
		if !hit {
			newBullets = append(newBullets, b)
		}
	}

	g.bullets = newBullets
}

func (g *Game) EnemyMovement() {
	// Move enemies left
	for i := range g.enemies {
		g.enemies[i].x -= 2 // Adjust speed as needed
	}
}

func (g *Game) RemoveOffScreenObjects() {
	// **Remove off-screen bullets**
	newBullets := g.bullets[:0]
	for _, b := range g.bullets {
		if b.y > 0 {
			newBullets = append(newBullets, b)
		}
	}
	g.bullets = newBullets

	// Remove enemies that move off the screen
	newEnemies := g.enemies[:0]
	for _, e := range g.enemies {
		if e.x < 480 { // Assuming screen height is 480
			newEnemies = append(newEnemies, e)
		}
	}
	g.enemies = newEnemies
}

// Update handles movement and shooting
func (g *Game) Update() error {
	if !g.hadFirstUpdate {
		g.extractPixels()
		g.hadFirstUpdate = true
	}

	var (
		newBullets []Bullet
		newEnemies []Enemy
	)

	// update player movement
	g.UpdatePlayer()

	// check for key press & update bullets
	g.BulletUpdate()

	// **Track the Space key state**
	g.firePressed = ebiten.IsKeyPressed(ebiten.KeySpace)

	// Enemy Collision Check
	g.CollisionCheck(newBullets, newEnemies)

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

	// Move Enemys
	g.EnemyMovement()

	// remove off screen bullets and enemys
	g.RemoveOffScreenObjects()

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
	playerImg, _, err := ebitenutil.NewImageFromFile("art/playership.png")
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
