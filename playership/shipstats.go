package playership

import (
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

	PlayerHealth int
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
			ps.X += 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if ps.Y > 0 {
			ps.Y -= 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if ps.Y < float64(ScreenHight-ps.PlayerImage.Bounds().Dy()) {
			ps.Y += 2
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
