package insideship

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth     = 960
	ScreenHight     = 540
	GroundPlayerURL = "art/groundPlayer.png"
	InsideShipURL   = "art/InsideShip.png"
)

type InsideShip struct {
	X, Y float64

	PlayerImage  *ebiten.Image
	PlayerPixels []byte

	CurrentBackground *ebiten.Image

	Health int
}

func (is *InsideShip) CheckIfCanInteract() (newMode int, hadModeUpdate bool) {
	if is.X >= float64(ScreenWidth-is.PlayerImage.Bounds().Dx()-75) {
		// CHECK IF BUTTON IS PRESSED
		if ebiten.IsKeyPressed(ebiten.KeyE) {
			newMode = 1
			hadModeUpdate = true
		}
	}

	return
}

func (is *InsideShip) SpawnInsideShip() {
	is.Health = 10
	is.X = float64(is.PlayerImage.Bounds().Dx()/2 + 20)
	is.Y = float64(ScreenHight - is.PlayerImage.Bounds().Dy() - 15)
}

func (is *InsideShip) GetCurrentGroundPixels() {
	playerPoint := is.PlayerImage.Bounds().Size()
	is.PlayerPixels = make([]byte, 4*playerPoint.X*playerPoint.Y)
	is.PlayerImage.ReadPixels(is.PlayerPixels)
}

func (is *InsideShip) GetCurrentShipImage() {
	var err error
	is.PlayerImage, _, err = ebitenutil.NewImageFromFile(GroundPlayerURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (is *InsideShip) UpdateGroundLocation() {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if is.X > 0 {
			is.X -= 4
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if is.X < float64(ScreenWidth-is.PlayerImage.Bounds().Dx()-25) {
			is.X += 4
		}
	}

}
