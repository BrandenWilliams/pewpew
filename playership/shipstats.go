package playership

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth   = 960
	ScreenHight   = 540
	PlayerShipURL = "art/playership.png"
)

// Bullet represents the location of a single shot
type Bullet struct {
	X, Y float64
}

type PlayerShip struct {
	X, Y float64
	// player stuff to be moved
	PlayerImage  *ebiten.Image
	PlayerPixels []byte

	Bullets []Bullet

	MaxPlayerHealth     int
	CurrentPlayerHealth int
}

// update player movement (spaceship)
func (ps *PlayerShip) UpdateShipLocation() {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if ps.X > 0 {
			ps.X -= 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if ps.X < float64(ScreenWidth-ps.PlayerImage.Bounds().Dx()-50) {
			ps.X += 4
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if ps.Y > 0 {
			ps.Y -= 3
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if ps.Y < float64(ScreenHight-ps.PlayerImage.Bounds().Dy()) {
			ps.Y += 3
		}
	}
}

// TODO make speed ajustable for upgrades later
func (ps *PlayerShip) MoveShipBullets() {
	// **Update bullet positions (keep this separate!)**
	for i := range ps.Bullets {
		ps.Bullets[i].X += 5 // Move bullets up
	}
}

func (ps *PlayerShip) GetCurrentShipImage() {
	// player sprite
	var err error
	ps.PlayerImage, _, err = ebitenutil.NewImageFromFile(PlayerShipURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (ps *PlayerShip) GetCurrentShipPixels() {
	playerPoint := ps.PlayerImage.Bounds().Size()
	ps.PlayerPixels = make([]byte, 4*playerPoint.X*playerPoint.Y)
	ps.PlayerImage.ReadPixels(ps.PlayerPixels)
}

func (ps *PlayerShip) DrawShipHealth(screen *ebiten.Image) {
	barWidth := ps.PlayerImage.Bounds().Dx()
	barHeight := 5.0 // Height of the health bar

	// Calculate the current width based on health
	healthPercentage := float64(ps.CurrentPlayerHealth) / float64(ps.MaxPlayerHealth)
	//log.Printf("health %: %v, currentHealth %v, maxHealth %v", healthPercentage, e.CurrentHealth, e.MaxHealth)
	currentBarWidth := float64(ps.PlayerImage.Bounds().Dx()) * healthPercentage

	// Create a red background bar (full size)
	healthBarBg := ebiten.NewImage(int(barWidth), int(barHeight))
	healthBarBg.Fill(color.RGBA{255, 0, 0, 255}) // Red for missing health

	// Draw the background first
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(ps.X, ps.Y)
	screen.DrawImage(healthBarBg, opts)

	var healthBar *ebiten.Image
	if currentBarWidth > 0 {
		// Create a green foreground bar (based on current health)
		healthBar = ebiten.NewImage(int(currentBarWidth), int(barHeight))
		healthBar.Fill(color.RGBA{0, 255, 0, 255}) // Green for remaining health

		// Draw the foreground (current health)
		opts.GeoM.Reset()
		opts.GeoM.Translate(ps.X, ps.Y)
		screen.DrawImage(healthBar, opts)
	}
}
