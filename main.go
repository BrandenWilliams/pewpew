package main

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth   = 640
	ScreenHight   = 480
	PlayerShipURL = "art/playership.png"
	EnemyShipURL  = "art/enemyship.png"
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
	hadFirstUpdate    bool
	playerImage       *ebiten.Image
	playerPixels      []byte
	enemyImage        *ebiten.Image
	enemyPixels       []byte
	bullets           []Bullet
	enemyBullets      []Bullet
	enemies           []Enemy
	x, y              float64
	firePressed       bool // Track fire key state
	enemySpawnTimer   float64
	initialSpawnDelay float64
}

// update player movement (spaceship)
func (g *Game) UpdatePlayer() {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if g.x > 0 {
			g.x -= 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if g.x < float64(ScreenWidth-g.playerImage.Bounds().Dx()-50) {
			g.x += 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if g.y > 0 {
			g.y -= 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if g.y < float64(ScreenHight-g.playerImage.Bounds().Dy()) {
			g.y += 2
		}
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

func (g *Game) CollisionCheck() {
	var newBullets []Bullet

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

func (g *Game) PlayerCollisionCheck() {
	var newEnemyBullets []Bullet
	// Get enemy sprite size
	ew, eh := g.playerImage.Size()

	for _, b := range g.enemyBullets {
		hit := false

		// First, do a bounding box check for early rejection
		if b.x+4 > g.x && b.x < g.x+float64(ew) &&
			b.y+2 > g.y && b.y+2 < g.y+float64(eh) {

			// Convert float positions to integers for pixel collision check
			bulletX := int(b.x)
			bulletY := int(b.y)
			playerX := int(g.x)
			playerY := int(g.y)

			// Call pixel-perfect collision detection
			if pixelsCollide(
				bulletX, bulletY, 4, 2,
				playerX, playerY, ew, eh, g.playerPixels,
			) {
				hit = true
				log.Println("you have been hit")
			}
		}

		// Keep the bullet if no hit occurred
		if !hit {
			newEnemyBullets = append(newEnemyBullets, b)
		}
	}

	g.enemyBullets = newEnemyBullets
}

func (g *Game) EnemyMovement() {
	// Move enemies left
	for i := range g.enemies {
		g.enemies[i].x -= 2 // Adjust speed as needed
	}
}

func (g *Game) EnemyBullets() {
	// Make each enemy fire a bullet every 2 seconds (adjust as needed)
	for _, e := range g.enemies {
		if rand.Float64() < 0.02 { // ~2% chance per frame
			g.enemyBullets = append(g.enemyBullets, Bullet{x: e.x, y: e.y + 16})
		}
	}

	// move bullets
	for i := range g.enemyBullets {
		g.enemyBullets[i].x -= 4
	}
}

func (g *Game) RemoveOffScreenObjects() {
	// **Remove off-screen player bullets** NEEDS TO BE REDONE
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
		if e.x < ScreenWidth {
			newEnemies = append(newEnemies, e)
		}
	}
	g.enemies = newEnemies

	// remove enemy bullets that are off screen
	newEnemyBullets := g.enemyBullets[:0]
	for _, eb := range g.enemyBullets {
		if eb.x > 0 {
			newEnemyBullets = append(newEnemyBullets, eb)
		}
	}
	g.enemyBullets = newEnemyBullets

}

// Update handles movement and shooting
func (g *Game) Update() error {
	if !g.hadFirstUpdate {
		g.extractPixels()
		g.hadFirstUpdate = true
	}

	// update player movement
	g.UpdatePlayer()

	// check for key press & update bullets
	g.BulletUpdate()

	// **Track the Space key state**
	g.firePressed = ebiten.IsKeyPressed(ebiten.KeySpace)

	// Enemy Collision Check
	g.CollisionCheck()

	// Check Player Collisions
	g.PlayerCollisionCheck()

	// Enemy spawn management
	g.EnemySpawn()

	// Move Enemys
	g.EnemyMovement()

	g.EnemyBullets()

	// remove off screen bullets and enemys
	g.RemoveOffScreenObjects()

	return nil
}

// long term make this support all enemys
func CreateEnemyLocation(enemyWidth, enemyHeight int) (el Enemy) {
	var newLocation Enemy
	// newLocation.x = float64(ScreenWidth - enemyWidth)
	newLocation.x = ScreenWidth - 1
	newLocation.y = float64(rand.Intn(ScreenHight - enemyHeight))
	el = newLocation
	log.Printf("el: x: %v y: %v", el.x, el.y)
	return
}

func (g *Game) EnemySpawn() {
	tps := ebiten.ActualTPS()

	// Ensure TPS is valid (avoid dividing by zero)
	if tps == 0 {
		tps = 60
	}

	// Start a delay before first wave
	if g.enemySpawnTimer < g.initialSpawnDelay {
		g.enemySpawnTimer += 1 / tps
		return
	}

	// Increment timer after the initial delay
	g.enemySpawnTimer += 1 / tps

	// Spawn a new enemy every 1.5 seconds
	spawnInterval := 1.5
	if g.enemySpawnTimer >= spawnInterval {
		g.enemies = append(g.enemies, CreateEnemyLocation(g.enemyImage.Bounds().Dx(), g.enemyImage.Bounds().Size().Y))
		g.enemySpawnTimer -= spawnInterval // Subtract instead of resetting to 0
	}
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

	// draw enemy bullets
	for _, eb := range g.enemyBullets {
		ebitenutil.DrawRect(screen, eb.x, eb.y, 4, 2, color.RGBA{255, 0, 0, 255})
	}

}

// Layout sets the game window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHight
}

func main() {
	// Load sprites

	// player sprite
	playerImg, _, err := ebitenutil.NewImageFromFile(PlayerShipURL)
	if err != nil {
		log.Fatal(err)
	}

	// enemy sprite
	enemyImg, _, err := ebitenutil.NewImageFromFile(EnemyShipURL)
	if err != nil {
		log.Fatal(err)
	}

	centerPlayer := (ScreenHight / 2) - playerImg.Bounds().Dy()/2

	game := &Game{
		playerImage: playerImg,
		enemyImage:  enemyImg,
		x:           50,
		y:           float64(centerPlayer),
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHight)
	ebiten.SetWindowTitle("pewpew v0.0.1")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
