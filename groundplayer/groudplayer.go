package groundplayer

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth     = 960
	ScreenHight     = 540
	GroundPlayerURL = "art/groundPlayer.png"
)

type GroundPlayer struct {
	X, Y float64

	PlayerImage  *ebiten.Image
	PlayerPixels []byte

	Health int
}

func (gp *GroundPlayer) SpawnGroundPlayer() {
	gp.Health = 10
	gp.X = float64(gp.PlayerImage.Bounds().Dx() / 2)
	gp.Y = float64(ScreenHight - gp.PlayerImage.Bounds().Dy() + 20)
}

func (gp *GroundPlayer) GetCurrentGroundPixels() {
	playerPoint := gp.PlayerImage.Bounds().Size()
	gp.PlayerPixels = make([]byte, 4*playerPoint.X*playerPoint.Y)
	gp.PlayerImage.ReadPixels(gp.PlayerPixels)
}

func (gp *GroundPlayer) GetCurrentShipImage() {
	// player sprite
	var err error
	gp.PlayerImage, _, err = ebitenutil.NewImageFromFile(GroundPlayerURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (gp *GroundPlayer) UpdateGroundLocation() {
	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if gp.X > 0 {
			gp.X -= 2
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if gp.X < float64(ScreenWidth-gp.PlayerImage.Bounds().Dx()-50) {
			gp.X += 2
		}
	}

	/*
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				if gp.Y > 0 {
					gp.Y -= 3
				}
			}

				if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			if gp.Y < float64(ScreenHight-gp.PlayerImage.Bounds().Dy()) {
				gp.Y += 3
			}
		}
	*/
}
