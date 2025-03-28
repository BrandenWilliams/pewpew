package main

import (
	"image/color"
	"log"

	"github.com/BrandenWilliams/pewpew/enemies"
	"github.com/BrandenWilliams/pewpew/enemyprojectiles"
	"github.com/BrandenWilliams/pewpew/groundplayer"
	"github.com/BrandenWilliams/pewpew/insideship"
	"github.com/BrandenWilliams/pewpew/playership"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	ScreenWidth   = 960
	ScreenHight   = 540
	PlayerShipURL = "art/playership.png"
)

// Game holds the player, bullets, and player position
type Game struct {
	gameMode int

	firstUpdate          bool
	switchMode           bool
	hadShipFirstUpdate   bool
	hadGroundFirstUpdate bool
	hadInsideFirstUpdate bool

	isDead bool

	insideShip   insideship.InsideShip
	groundPlayer groundplayer.GroundPlayer
	playerShip   playership.PlayerShip

	allEnemyPixels   enemies.AllEP
	enemies          enemies.Enemies
	enemyProjectiles enemyprojectiles.EnemyProjectiles

	firePressed  bool // Track fire key state
	enterPressed bool // for respawn
}

func (g *Game) extractPixels() {
	// Get current Ship pixels later this will update for dmg
	g.playerShip.GetCurrentShipPixels()

	var err error
	if g.allEnemyPixels, err = g.enemies.GetAllEnemyPixels(); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) ShipBulletFire() {
	// **Shooting (Spacebar, one shot per press)**
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.firePressed {
		g.playerShip.Bullets = append(g.playerShip.Bullets, playership.Bullet{X: g.playerShip.X + 122, Y: g.playerShip.Y + 62})
	}
}

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
	var newBullets []playership.Bullet

	for _, b := range g.playerShip.Bullets {
		hit := false

		for i := 0; i < len(g.enemies.ES); i++ {
			e := g.enemies.ES[i]

			// Get enemy sprite size
			ew, eh := g.enemies.ES[i].EnemyImage.Size()

			// First, do a bounding box check for early rejection
			if b.X+4 > e.X && b.X < e.X+float64(ew) &&
				b.Y+2 > e.Y && b.Y < e.Y+float64(eh) {

				// Convert float positions to integers for pixel collision check
				bulletX := int(b.X)
				bulletY := int(b.Y)
				enemyX := int(e.X)
				enemyY := int(e.Y)

				// Call pixel-perfect collision detection
				if pixelsCollide(
					bulletX, bulletY, 4, 2,
					enemyX, enemyY, ew, eh, g.enemies.ES[i].EnemyPixels,
				) {
					hit = true

					damage := 1
					g.enemies.ES[i].CurrentHealth -= damage

					if g.enemies.ES[i].CurrentHealth <= 0 {
						g.enemies.ES = append(g.enemies.ES[:i], g.enemies.ES[i+1:]...)
					}

					break
				}
			}
		}

		// Keep the bullet if no hit occurred
		if !hit {
			newBullets = append(newBullets, b)
		}
	}

	g.playerShip.Bullets = newBullets
}

func (g *Game) PlayerCollisionCheck() {
	var newEnemyBullets []enemyprojectiles.Bullet
	// Get enemy sprite size
	ew, eh := g.playerShip.PlayerImage.Size()

	// must iterate through all enmies now that bullets are within enemys
	for _, b := range g.enemyProjectiles.EnemyBullets {
		hit := false

		// First, do a bounding box check for early rejection
		if b.X+4 > g.playerShip.X && b.X < g.playerShip.X+float64(ew) &&
			b.Y+2 > g.playerShip.Y && b.Y+2 < g.playerShip.Y+float64(eh) {

			// Convert float positions to integers for pixel collision check
			bulletX := int(b.X)
			bulletY := int(b.Y)
			playerX := int(g.playerShip.X)
			playerY := int(g.playerShip.Y)

			// Call pixel-perfect collision detection
			if pixelsCollide(
				bulletX, bulletY, 4, 2,
				playerX, playerY, ew, eh, g.playerShip.PlayerPixels,
			) {
				g.playerShip.CurrentPlayerHealth -= 1
				hit = true

				if g.playerShip.CurrentPlayerHealth <= 0 {
					g.isDead = true
				}
			}
		}

		// Keep the bullet if no hit occurred
		if !hit {
			newEnemyBullets = append(newEnemyBullets, b)
		}
	}

	g.enemyProjectiles.EnemyBullets = newEnemyBullets
}

func (g *Game) DespawnOffScreenObjects() {
	// depsawn player bullets
	g.removePlayerBullets()

	// Depsawn Enemys
	g.enemies.EnemyDespawn()

	// Depsawn Enemy Bullets
	g.enemyProjectiles.DespawnEnemyProjectiles()
}

// Remove player bullets
func (g *Game) removePlayerBullets() {
	var newBullets []playership.Bullet
	for _, b := range g.playerShip.Bullets {
		if b.X < ScreenWidth {
			newBullets = append(newBullets, b)
		}
	}
	g.playerShip.Bullets = newBullets
}

func (g *Game) SpawnEnemyProjectiles() {
	for _, enemySlice := range g.enemies.ES {
		g.enemyProjectiles.NewProjectile(enemySlice.X, enemySlice.Y, enemySlice.ProjectileType)
	}
}

func (g *Game) RespawnPlayer() {
	g.playerShip.CurrentPlayerHealth = 10
	g.playerShip.MaxPlayerHealth = 10
	g.playerShip.X = 50
	g.playerShip.Y = (ScreenHight / 2) - (float64(g.playerShip.PlayerImage.Bounds().Dy() / 2))
}

func (g *Game) DespawnAllBullets() {
	var newProjectiles enemyprojectiles.EnemyProjectiles
	g.enemyProjectiles = newProjectiles
}

func (g *Game) UpdateSpaceShipMode() error {
	if !g.hadShipFirstUpdate {
		g.extractPixels()
		g.hadShipFirstUpdate = true
		g.RespawnPlayer()
	}

	if !g.isDead {
		// TPS TEST KEEP FOR FUTURE DEBUGG
		// log.Printf("TPS: %v", ebiten.ActualTPS())

		// update player movement
		g.playerShip.UpdateShipLocation()

		// check for key press & update bullets
		g.ShipBulletFire()

		// **Track the Space key state**
		g.firePressed = ebiten.IsKeyPressed(ebiten.KeySpace)

		// move player bullets
		g.playerShip.MoveShipBullets()

		// Enemy Collision Check
		g.CollisionCheck()

		// Check Player Collisions
		g.PlayerCollisionCheck()

		// Enemy spawn management
		g.enemies.SpawnOneEnemy()
		// g.enemies.EnemySpawn()

		if len(g.enemies.ES) > 0 {
			// Move Enemys
			g.enemies.EnemiesMovement()
			// Spawn Projectiles
			g.SpawnEnemyProjectiles()
			// move projectiles
			g.enemyProjectiles.ManageEnemyProjectiles()
		}

		// remove off screen bullets and enemys
		g.DespawnOffScreenObjects()
	} else if g.isDead {
		g.enterPressed = ebiten.IsKeyPressed(ebiten.KeyEnter)

		if g.enterPressed {
			g.isDead = false
			g.RespawnPlayer()
			g.DespawnAllBullets()
			g.enemies.DespawnAllEnemies()
			return nil
		}
	}

	return nil
}

func (g *Game) ClearFirstUpdates() {
	g.hadGroundFirstUpdate = false
	g.hadInsideFirstUpdate = false
	g.hadShipFirstUpdate = false
}

func (g *Game) UpdateGroundMode() error {
	if !g.hadGroundFirstUpdate {
		g.groundPlayer.GetCurrentShipImage()
		g.groundPlayer.GetCurrentGroundPixels()
		g.groundPlayer.SpawnGroundPlayer()
		g.hadGroundFirstUpdate = true
	}

	// Update Ground Location
	g.groundPlayer.UpdateGroundLocation()

	return nil
}

func (g *Game) UpdateInsideShipMode() error {
	if !g.hadInsideFirstUpdate {
		g.insideShip.GetCurrentShipImage()
		g.insideShip.GetCurrentGroundPixels()
		g.insideShip.SpawnInsideShip()
		g.hadInsideFirstUpdate = true
	}

	// Update Ground Location
	g.insideShip.UpdateGroundLocation()

	g.gameMode, g.switchMode = g.insideShip.CheckIfCanInteract()

	if g.switchMode {
		g.ClearFirstUpdates()
	}

	return nil
}

// Update handles movement and shooting
func (g *Game) Update() error {
	if g.firstUpdate {
		g.gameMode = 0
		g.firstUpdate = true
	}

	if g.switchMode {
		g.switchMode = false
	}

	switch g.gameMode {
	case 0:
		g.UpdateInsideShipMode()
	case 1:
		g.UpdateGroundMode()
	case 2:
		g.UpdateSpaceShipMode()
	}

	return nil
}

func (g *Game) DrawGameOver(screen *ebiten.Image) {
	msg := "GAME OVER\n push enter to try again"
	text.Draw(screen, msg, basicfont.Face7x13, ScreenWidth/2-len(msg)*7, ScreenHight/2, color.White)
}

func (g *Game) SetInsideShipBackground(screen *ebiten.Image) {
	bgImage, _, err := ebitenutil.NewImageFromFile("art/insideShip.png")
	if err != nil {
		log.Fatal(err)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(0, 0)
	screen.DrawImage(bgImage, op) // Draw at (0,0)
}

func (g *Game) DrawInsideShip(screen *ebiten.Image) {
	g.SetInsideShipBackground(screen)

	// Draw player sprite
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(g.insideShip.X, g.insideShip.Y)
	screen.DrawImage(g.insideShip.PlayerImage, playerOp)
}

func (g *Game) setGroundModeBackground(level int, screen *ebiten.Image) {
	switch level {
	case 1:
		g.SetInsideShipBackground(screen)
	case 2:
	default:
	}
}

func (g *Game) DrawGroundMode(screen *ebiten.Image) {
	// g.setGroundModeBackground(1, screen)

	// Draw player sprite
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(g.groundPlayer.X, g.groundPlayer.Y)
	screen.DrawImage(g.groundPlayer.PlayerImage, playerOp)
}

func (g *Game) DrawShipMode(screen *ebiten.Image) {
	if g.isDead {
		g.DrawGameOver(screen)
		return
	}

	// Draw player sprite
	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(g.playerShip.X, g.playerShip.Y)
	screen.DrawImage(g.playerShip.PlayerImage, playerOp)

	g.playerShip.DrawShipHealth(screen)

	// Draw enemy sprite
	for _, e := range g.enemies.ES {
		enemyOp := &ebiten.DrawImageOptions{}
		enemyOp.GeoM.Translate(e.X, e.Y)
		screen.DrawImage(e.EnemyImage, enemyOp)
		g.enemies.DrawEnemyHealthBars(e, screen)
	}

	// Draw bullets
	for _, b := range g.playerShip.Bullets {
		ebitenutil.DrawRect(screen, b.X, b.Y, 4, 2, color.RGBA{255, 255, 0, 255})
	}

	// draw enemy bullets
	for _, eb := range g.enemyProjectiles.EnemyBullets {
		ebitenutil.DrawRect(screen, eb.X, eb.Y, 4, 2, color.RGBA{255, 0, 0, 255})
	}
}

// Draw the player and bullets
func (g *Game) Draw(screen *ebiten.Image) {
	if g.switchMode {
		return
	}

	switch g.gameMode {
	case 0:
		g.DrawInsideShip(screen)
	case 1:
		g.DrawGroundMode(screen)
	case 2:
		g.DrawShipMode(screen)
	}
}

// Layout sets the game window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHight
}

func main() {
	// player sprite
	playerImg, _, err := ebitenutil.NewImageFromFile(PlayerShipURL)
	if err != nil {
		log.Fatal(err)
	}

	centerPlayer := (ScreenHight / 2) - playerImg.Bounds().Dy()/2

	game := &Game{
		playerShip: playership.PlayerShip{
			PlayerImage: playerImg,
			X:           50,
			Y:           float64(centerPlayer),
		},
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHight)
	ebiten.SetWindowTitle("pewpew v0.0.1")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
