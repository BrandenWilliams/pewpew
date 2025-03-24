package main

import (
	"image/color"
	"log"

	"github.com/BrandenWilliams/pewpew/enemies"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth   = 960
	ScreenHight   = 540
	PlayerShipURL = "art/playership.png"
	EnemyShipURL  = "art/enemyship.png"
)

// Bullet represents the location of a single shot
type Bullet struct {
	x, y float64
}

// Game holds the player, bullets, and player position
type Game struct {
	hadFirstUpdate bool
	playerImage    *ebiten.Image
	playerPixels   []byte

	bullets []Bullet

	enemies   enemies.Enemies
	enemyType enemies.EnemyType

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

	newEnemyType := g.enemyType.TypeOne()
	EnemyPoint := newEnemyType.Image.Bounds().Size()
	g.enemyType.EnemyPixels = make([]byte, 4*EnemyPoint.X*EnemyPoint.Y)
	newEnemyType.Image.ReadPixels(g.enemyType.EnemyPixels)
}

// CRASH ON COLLIDE

// bullet x, y, bullet width, bullet hieght, enemy X posision, enemy Y Posision, enemy width, enemy height and image pixels
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

		for i := 0; i < len(g.enemies.ES); i++ {
			e := g.enemies.ES[i]

			// Get enemy sprite size
			ew, eh := g.enemies.ES[i].EnemyImage.Size()

			// First, do a bounding box check for early rejection
			if b.x+4 > e.X && b.x < e.X+float64(ew) &&
				b.y+2 > e.Y && b.y < e.Y+float64(eh) {

				// Convert float positions to integers for pixel collision check
				bulletX := int(b.x)
				bulletY := int(b.y)
				enemyX := int(e.X)
				enemyY := int(e.Y)

				// Call pixel-perfect collision detection
				if pixelsCollide(
					bulletX, bulletY, 4, 2,
					enemyX, enemyY, ew, eh, g.enemies.ES[i].EnemyPixels,
				) {
					hit = true
					g.enemies.ES = append(g.enemies.ES[:i], g.enemies.ES[i+1:]...)
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
	var newEnemyBullets []enemies.Bullet
	// Get enemy sprite size
	ew, eh := g.playerImage.Size()

	// must iterate through all enmies now that bullets are within enemys
	for _, e := range g.enemies.ES {

		for _, b := range e.EnemyBullets {
			hit := false

			// First, do a bounding box check for early rejection
			if b.X+4 > g.x && b.X < g.x+float64(ew) &&
				b.Y+2 > g.y && b.Y+2 < g.y+float64(eh) {

				// Convert float positions to integers for pixel collision check
				bulletX := int(b.X)
				bulletY := int(b.Y)
				playerX := int(g.x)
				playerY := int(g.y)

				// Call pixel-perfect collision detection
				if pixelsCollide(
					bulletX, bulletY, 4, 2,
					playerX, playerY, ew, eh, g.playerPixels,
				) {
					hit = true
				}
			}

			// Keep the bullet if no hit occurred
			if !hit {
				newEnemyBullets = append(newEnemyBullets, b)
			}
		}
		e.EnemyBullets = newEnemyBullets
	}
}

func (g *Game) removeEnemyBullets(oeb []enemies.Bullet) (neb []enemies.Bullet) {
	// remove enemy bullets that are off screen
	newEnemyBullets := oeb[:0]
	for _, eb := range oeb {
		if eb.X > 0 {
			newEnemyBullets = append(newEnemyBullets, eb)
		}
	}

	neb = newEnemyBullets

	return
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
	newEnemies := g.enemies.ES[:0]
	for _, e := range g.enemies.ES {
		if e.X < ScreenWidth {
			e.EnemyBullets = g.removeEnemyBullets(e.EnemyBullets)
			newEnemies = append(newEnemies, e)
		}
	}
	g.enemies.ES = newEnemies

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
	g.enemies.EnemySpawn()

	// Move Enemys
	g.enemies.EnemyMovement()

	g.enemies.EnemyBullets()

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
	for _, e := range g.enemies.ES {
		enemyOp := &ebiten.DrawImageOptions{}
		enemyOp.GeoM.Translate(e.X, e.Y)
		screen.DrawImage(e.EnemyImage, enemyOp)
	}

	// Draw bullets
	for _, b := range g.bullets {
		ebitenutil.DrawRect(screen, b.x, b.y, 4, 2, color.RGBA{255, 255, 0, 255})
	}

	// draw enemy bullets
	for _, e := range g.enemies.ES {
		for _, eb := range e.EnemyBullets {
			ebitenutil.DrawRect(screen, eb.X, eb.Y, 4, 2, color.RGBA{255, 0, 0, 255})
		}
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

	centerPlayer := (ScreenHight / 2) - playerImg.Bounds().Dy()/2

	game := &Game{
		playerImage: playerImg,
		x:           50,
		y:           float64(centerPlayer),
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHight)
	ebiten.SetWindowTitle("pewpew v0.0.1")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
