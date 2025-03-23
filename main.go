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

func pixelsCollide(x1, y1 int, pixels1 []byte, w1, h1 int,
	x2, y2 int, pixels2 []byte, w2, h2 int) bool {

	// Calculate the overlap region
	xStart := max(x1, x2)
	yStart := max(y1, y2)
	xEnd := min(x1+w1, x2+w2)
	yEnd := min(y1+h1, y2+h2)

	// If no overlap, bail early
	if xStart >= xEnd || yStart >= yEnd {
		return false
	}

	// Check overlapping pixels for transparency
	for y := yStart; y < yEnd; y++ {
		for x := xStart; x < xEnd; x++ {
			// Get pixel offsets for both images
			offset1 := ((y-y1)*w1 + (x - x1)) * 4
			offset2 := ((y-y2)*w2 + (x - x2)) * 4

			// Check alpha values (4th byte in RGBA) for non-transparent pixels
			if pixels1[offset1+3] > 0 && pixels2[offset2+3] > 0 {
				return true // Collision detected!
			}
		}
	}

	return false
}

func (g *Game) CollisionCheck(newBullets []Bullet, newEnemies []Enemy) {
	for _, b := range g.bullets {
		hit := false

		for i := 0; i < len(g.enemies); i++ {
			e := g.enemies[i]

			// Get player/enemy sprite sizes
			pw, ph := g.playerImage.Size()
			ew, eh := g.enemyImage.Size()

			// Pixel-perfect collision between bullet and enemy
			if pixelsCollide(
				int(b.x), int(b.y), g.playerPixels, pw, ph,
				int(e.x), int(e.y), g.enemyPixels, ew, eh,
			) {
				hit = true
				g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
				break
			}
		}

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
